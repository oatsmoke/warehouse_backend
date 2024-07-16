package service

import (
	"context"
	"errors"
	"strings"
	"time"
	"warehouse_backend/internal/model"
	"warehouse_backend/internal/repository"
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

func (s *EquipmentService) Create(ctx context.Context, date int64, company int64, serialNumber string, profile int64, userId int64) (int64, error) {
	serialNumber = strings.ToUpper(serialNumber)
	if _, err := s.repositoryEquipment.FindBySerialNumber(ctx, serialNumber); err == nil {
		return 0, errors.New("serial number already exists")
	}
	return s.repositoryEquipment.Create(ctx, date, company, serialNumber, profile, userId)
}

func (s *EquipmentService) GetById(ctx context.Context, id int64) (*model.Location, error) {
	return s.repositoryEquipment.GetById(ctx, id)
}

func (s *EquipmentService) GetByIds(ctx context.Context, ids []int64) ([]*model.Location, error) {
	var equipments []*model.Location
	for _, id := range ids {
		equipment, err := s.repositoryEquipment.GetById(ctx, id)
		if err != nil {
			return nil, err
		}
		equipments = append(equipments, equipment)
	}
	return equipments, nil
}

func (s *EquipmentService) GetByLocation(ctx context.Context, toDepartment, toEmployee, toContract int64) ([]*model.Location, error) {
	switch {
	case toDepartment == 0 && toEmployee == 0 && toContract == 0:
		return s.repositoryEquipment.GetByLocationStorage(ctx)
	case toDepartment != 0:
		return s.repositoryEquipment.GetByLocationDepartment(ctx, toDepartment)
	case toEmployee != 0:
		return s.repositoryEquipment.GetByLocationEmployee(ctx, toEmployee)
	case toContract != 0:
		return s.repositoryEquipment.GetByLocationContract(ctx, toContract)
	}
	return nil, nil
}

func (s *EquipmentService) GetAll(ctx context.Context) ([]*model.Equipment, error) {
	return s.repositoryEquipment.GetAll(ctx)
}

func (s *EquipmentService) Update(ctx context.Context, id int64, serialNumber string, profileId int64) error {
	serialNumber = strings.ToUpper(serialNumber)
	findId, err := s.repositoryEquipment.FindBySerialNumber(ctx, serialNumber)
	if findId != id && err == nil {
		return errors.New("serial number already exists")
	}

	return s.repositoryEquipment.Update(ctx, id, serialNumber, profileId)
}

func (s *EquipmentService) Delete(ctx context.Context, id int64) error {
	return s.repositoryEquipment.Delete(ctx, id)
}

func (s *EquipmentService) ReportByCategory(ctx context.Context, departmentId int64, date int64) (*model.Report, error) {
	report := new(model.Report)
	fromDate := time.Unix(date, 0)
	toDate := time.Unix(date, 0).AddDate(0, 1, 0)
	categories, err := s.repositoryCategory.GetAll(ctx, false)
	if err != nil {
		return nil, err
	}
	report.Categories = categories
	departments := make(map[int64]*model.Department)
	leftover := make(map[int64][]*model.Location)
	total := make(map[int64][]*model.Location)
	fromStorage := make(map[int64][]*model.Location)
	toStorage := make(map[int64][]*model.Location)
	fromContract := make(map[int64][]*model.Location)
	toContract := make(map[int64][]*model.Location)
	fromDepartment := make(map[int64]map[int64][]*model.Location)
	toDepartment := make(map[int64]map[int64][]*model.Location)
	for _, category := range categories {
		equipment, err := s.repositoryEquipment.RemainderByCategory(ctx, category.ID, departmentId, fromDate)
		if err != nil {
			return nil, err
		}
		leftover[category.ID] = equipment

		equipment, err = s.repositoryEquipment.RemainderByCategory(ctx, category.ID, departmentId, toDate)
		if err != nil {
			return nil, err
		}
		total[category.ID] = equipment

		equipment, err = s.repositoryEquipment.TransferByCategory(ctx, category.ID, departmentId, fromDate, toDate, "STORAGE_TO_DEPARTMENT")
		if err != nil {
			return nil, err
		}
		fromStorage[category.ID] = equipment

		equipment, err = s.repositoryEquipment.TransferByCategory(ctx, category.ID, departmentId, fromDate, toDate, "DEPARTMENT_TO_STORAGE")
		if err != nil {
			return nil, err
		}
		toStorage[category.ID] = equipment

		equipment, err = s.repositoryEquipment.TransferByCategory(ctx, category.ID, departmentId, fromDate, toDate, "CONTRACT_TO_DEPARTMENT")
		if err != nil {
			return nil, err
		}
		fromContract[category.ID] = equipment

		equipment, err = s.repositoryEquipment.TransferByCategory(ctx, category.ID, departmentId, fromDate, toDate, "DEPARTMENT_TO_CONTRACT")
		if err != nil {
			return nil, err
		}
		toContract[category.ID] = equipment

		equipment, err = s.repositoryEquipment.FromDepartmentTransferByCategory(ctx, category.ID, departmentId, fromDate, toDate)
		if err != nil {
			return nil, err
		}

		locationFromDepartment := make(map[int64][]*model.Location)
		for _, row := range equipment {
			departments[row.FromDepartment.ID] = row.FromDepartment
			locationFromDepartment[row.FromDepartment.ID] = append(locationFromDepartment[row.FromDepartment.ID], row)
		}
		fromDepartment[category.ID] = locationFromDepartment
		equipmentFrom, err := s.repositoryEquipment.ToDepartmentTransferByCategory(ctx, category.ID, departmentId, fromDate, toDate)
		if err != nil {
			return nil, err
		}

		locationToDepartment := make(map[int64][]*model.Location)
		for _, row := range equipmentFrom {
			departments[row.ToDepartment.ID] = row.ToDepartment
			locationToDepartment[row.ToDepartment.ID] = append(locationToDepartment[row.ToDepartment.ID], row)
		}
		toDepartment[category.ID] = locationToDepartment
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
