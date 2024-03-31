package service

import (
	"errors"
	"strings"
	"time"
	"warehouse_backend/pkg/model"
	"warehouse_backend/pkg/repository"
)

type EquipmentService struct {
	repositoryEquipment repository.Equipment
	repositoryCategory  repository.Category
}

func NewEquipmentService(repositoryEquipment repository.Equipment, repositoryCategory repository.Category) *EquipmentService {
	return &EquipmentService{
		repositoryEquipment: repositoryEquipment,
		repositoryCategory:  repositoryCategory,
	}
}

func (s *EquipmentService) Create(date int64, company int, serialNumber string, profile int, userId int) (int, error) {
	serialNumber = strings.ToUpper(serialNumber)
	if _, err := s.repositoryEquipment.FindBySerialNumber(serialNumber); err == nil {
		return 0, errors.New("serial number already exists")
	}
	return s.repositoryEquipment.Create(date, company, serialNumber, profile, userId)
}

func (s *EquipmentService) GetById(id int) (model.Location, error) {
	return s.repositoryEquipment.GetById(id)
}

func (s *EquipmentService) GetByIds(ids []int) ([]model.Location, error) {
	var equipments []model.Location
	for _, id := range ids {
		equipment, err := s.repositoryEquipment.GetById(id)
		if err != nil {
			return nil, err
		}
		equipments = append(equipments, equipment)
	}
	return equipments, nil
}

func (s *EquipmentService) GetByLocation(toDepartment, toEmployee, toContract int) ([]model.Location, error) {
	switch {
	case toDepartment == 0 && toEmployee == 0 && toContract == 0:
		return s.repositoryEquipment.GetByLocationStorage()
	case toDepartment != 0 && toEmployee == 0 && toContract == 0:
		return s.repositoryEquipment.GetByLocationDepartment(toDepartment)
	case toDepartment == 0 && toEmployee != 0 && toContract == 0:
		return s.repositoryEquipment.GetByLocationEmployee(toEmployee)
	case toDepartment == 0 && toEmployee == 0 && toContract != 0:
		return s.repositoryEquipment.GetByLocationContract(toContract)
	}
	return []model.Location{}, nil
}

func (s *EquipmentService) GetAll() ([]model.Equipment, error) {
	return s.repositoryEquipment.GetAll()
}

func (s *EquipmentService) Update(id int, serialNumber string, profile int) error {
	serialNumber = strings.ToUpper(serialNumber)
	findId, err := s.repositoryEquipment.FindBySerialNumber(serialNumber)
	if findId != id && err == nil {
		return errors.New("serial number already exists")
	}
	return s.repositoryEquipment.Update(id, serialNumber, profile)
}

func (s *EquipmentService) Delete(id int) error {
	return s.repositoryEquipment.Delete(id)
}

func (s *EquipmentService) ReportByCategory(departmentId int, date int64) (model.Report, error) {
	var report model.Report
	fromDate := time.Unix(date, 0)
	toDate := time.Unix(date, 0).AddDate(0, 1, 0)
	categories, err := s.repositoryCategory.GetAll()
	if err != nil {
		return model.Report{}, err
	}
	report.Categories = categories
	departments := make(map[int]model.Department)
	leftover := make(map[int][]model.Location)
	total := make(map[int][]model.Location)
	fromStorage := make(map[int][]model.Location)
	toStorage := make(map[int][]model.Location)
	fromContract := make(map[int][]model.Location)
	toContract := make(map[int][]model.Location)
	fromDepartment := make(map[int]map[int][]model.Location)
	toDepartment := make(map[int]map[int][]model.Location)
	for _, category := range categories {
		equipment, err := s.repositoryEquipment.RemainderByCategory(category.Id, departmentId, fromDate)
		if err != nil {
			return model.Report{}, err
		}
		leftover[category.Id] = equipment
		equipment, err = s.repositoryEquipment.RemainderByCategory(category.Id, departmentId, toDate)
		if err != nil {
			return model.Report{}, err
		}
		total[category.Id] = equipment
		equipment, err = s.repositoryEquipment.TransferByCategory(category.Id, departmentId, fromDate, toDate, "STORAGE_TO_DEPARTMENT")
		if err != nil {
			return model.Report{}, err
		}
		fromStorage[category.Id] = equipment
		equipment, err = s.repositoryEquipment.TransferByCategory(category.Id, departmentId, fromDate, toDate, "DEPARTMENT_TO_STORAGE")
		if err != nil {
			return model.Report{}, err
		}
		toStorage[category.Id] = equipment
		equipment, err = s.repositoryEquipment.TransferByCategory(category.Id, departmentId, fromDate, toDate, "CONTRACT_TO_DEPARTMENT")
		if err != nil {
			return model.Report{}, err
		}
		fromContract[category.Id] = equipment
		equipment, err = s.repositoryEquipment.TransferByCategory(category.Id, departmentId, fromDate, toDate, "DEPARTMENT_TO_CONTRACT")
		if err != nil {
			return model.Report{}, err
		}
		toContract[category.Id] = equipment
		equipment, err = s.repositoryEquipment.FromDepartmentTransferByCategory(category.Id, departmentId, fromDate, toDate)
		if err != nil {
			return model.Report{}, err
		}
		locationFromDepartment := make(map[int][]model.Location)
		for _, row := range equipment {
			departments[row.FromDepartment.Id] = row.FromDepartment
			locationFromDepartment[row.FromDepartment.Id] = append(locationFromDepartment[row.FromDepartment.Id], row)
		}
		fromDepartment[category.Id] = locationFromDepartment
		equipmentFrom, err := s.repositoryEquipment.ToDepartmentTransferByCategory(category.Id, departmentId, fromDate, toDate)
		if err != nil {
			return model.Report{}, err
		}
		locationToDepartment := make(map[int][]model.Location)
		for _, row := range equipmentFrom {
			departments[row.ToDepartment.Id] = row.ToDepartment
			locationToDepartment[row.ToDepartment.Id] = append(locationToDepartment[row.ToDepartment.Id], row)
		}
		toDepartment[category.Id] = locationToDepartment
	}
	for _, department := range departments {
		report.Departments = append(report.Departments, department)
	}
	report.Leftover = leftover
	report.Total = total
	report.FromStorage = fromStorage
	report.ToStorage = toStorage
	report.FromContract = fromContract
	report.ToContract = toContract
	report.FromDepartment = fromDepartment
	report.ToDepartment = toDepartment
	return report, nil
}
