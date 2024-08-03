package service

import (
	"context"
	"strings"
	"time"
	"warehouse_backend/internal/model"
	"warehouse_backend/internal/repository"
)

type LocationService struct {
	repositoryLocation repository.Location
	repositoryReplace  repository.Replace
}

func NewLocationService(repositoryLocation repository.Location,
	repositoryReplace repository.Replace) *LocationService {
	return &LocationService{repositoryLocation: repositoryLocation,
		repositoryReplace: repositoryReplace,
	}
}

func (s *LocationService) TransferTo(ctx context.Context, id int64, requests []*model.RequestLocation) error {
	var code string
	replace := []int64{0, 0}
	for _, request := range requests {
		nowLocation, _ := s.repositoryLocation.GetLocationNow(ctx, request.EquipmentId)
		switch {
		case request.ToDepartment == 0 && request.ToEmployee == 0 && request.ToContract == 0:
			code = strings.ToUpper(request.ThisLocation + "_to_storage")
			transferId, err := s.repositoryLocation.TransferToStorage(ctx, request.Date, code, request.EquipmentId, id, request.Company, nowLocation)
			if err != nil {
				return err
			}
			if request.Way == "replace" && replace[0] == 0 {
				replace[0] = transferId
			} else if request.Way == "replace" && replace[1] == 0 {
				replace[1] = transferId
			}
		case request.ToDepartment != 0 && request.ToEmployee == 0 && request.ToContract == 0:
			if request.InDepartment {
				code = strings.ToUpper(request.ThisLocation + "_to_" + request.Where + "_in_department")
			} else {
				code = strings.ToUpper(request.ThisLocation + "_to_" + request.Where)
			}
			transferId, err := s.repositoryLocation.TransferToDepartment(ctx, request.Date, code, request.EquipmentId, id, request.Company, request.ToDepartment, nowLocation)
			if err != nil {
				return err
			}
			if request.Way == "replace" && replace[0] == 0 {
				replace[0] = transferId
			} else if request.Way == "replace" && replace[1] == 0 {
				replace[1] = transferId
			}
		case request.ToDepartment == 0 && request.ToEmployee != 0 && request.ToContract == 0:
			code = strings.ToUpper(request.ThisLocation + "_to_" + request.Where)
			transferId, err := s.repositoryLocation.TransferToEmployee(ctx, request.Date, code, request.EquipmentId, id, request.Company, request.ToEmployee, nowLocation)
			if err != nil {
				return err
			}
			if request.Way == "replace" && replace[0] == 0 {
				replace[0] = transferId
			} else if request.Way == "replace" && replace[1] == 0 {
				replace[1] = transferId
			}
		case request.ToDepartment != 0 && request.ToEmployee != 0 && request.ToContract == 0:
			if request.InDepartment {
				code = strings.ToUpper(request.ThisLocation + "_to_" + request.Where + "_in_department")
			} else {
				code = strings.ToUpper(request.ThisLocation + "_to_" + request.Where)
			}
			transferId, err := s.repositoryLocation.TransferToEmployeeInDepartment(ctx, request.Date, code, request.EquipmentId, id, request.Company, request.ToDepartment, request.ToEmployee, nowLocation)
			if err != nil {
				return err
			}
			if request.Way == "replace" && replace[0] == 0 {
				replace[0] = transferId
			} else if request.Way == "replace" && replace[1] == 0 {
				replace[1] = transferId
			}
		case request.ToDepartment == 0 && request.ToEmployee == 0 && request.ToContract != 0:
			code = strings.ToUpper(request.ThisLocation + "_to_" + request.Where)
			transferId, err := s.repositoryLocation.TransferToContract(ctx, request.Date, code, request.EquipmentId, id, request.Company, request.ToContract, request.TransferType, request.Price, nowLocation)
			if err != nil {
				return err
			}
			if request.Way == "replace" && replace[0] == 0 {
				replace[0] = transferId
			} else if request.Way == "replace" && replace[1] == 0 {
				replace[1] = transferId
			}
		}
	}
	if requests[0].Way == "replace" {
		err := s.repositoryReplace.Create(ctx, replace)
		if err != nil {
			return err
		}
	}
	return nil
}
func (s *LocationService) GetHistory(ctx context.Context, id int64) ([]*model.Location, error) {
	return s.repositoryLocation.GetHistory(ctx, id)
}

func (s *LocationService) Delete(ctx context.Context, id int64) error {
	replace, err := s.repositoryReplace.FindByLocationId(ctx, id)
	if err != nil {
		return s.repositoryLocation.Delete(ctx, id)
	}

	if err = s.repositoryLocation.Delete(ctx, replace.TransferFrom); err != nil {
		return err
	}

	if err = s.repositoryLocation.Delete(ctx, replace.TransferTo); err != nil {
		return err
	}

	return nil
}

func (s *EquipmentService) GetById(ctx context.Context, id int64) (*model.Location, error) {
	return s.EquipmentRepository.GetById(ctx, id)
}

func (s *EquipmentService) GetByIds(ctx context.Context, ids []int64) ([]*model.Location, error) {
	var equipments []*model.Location
	for _, id := range ids {
		equipment, err := s.EquipmentRepository.GetById(ctx, id)
		if err != nil {
			return nil, err
		}
		equipments = append(equipments, equipment)
	}
	return equipments, nil
}

func (s *EquipmentService) GetByLocation(ctx context.Context, toDepartmentId, toEmployeeId, toContractId int64) ([]*model.Location, error) {
	switch {
	case toDepartmentId != 0:
		return s.EquipmentRepository.GetByLocationDepartment(ctx, toDepartmentId)
	case toEmployeeId != 0:
		return s.EquipmentRepository.GetByLocationEmployee(ctx, toEmployeeId)
	case toContractId != 0:
		return s.EquipmentRepository.GetByLocationContract(ctx, toContractId)
	default:
		return s.EquipmentRepository.GetByLocationStorage(ctx)
	}
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
		equipment, err := s.EquipmentRepository.RemainderByCategory(ctx, category.ID, departmentId, fromDate)
		if err != nil {
			return nil, err
		}
		leftover[category.ID] = equipment

		equipment, err = s.EquipmentRepository.RemainderByCategory(ctx, category.ID, departmentId, toDate)
		if err != nil {
			return nil, err
		}
		total[category.ID] = equipment

		equipment, err = s.EquipmentRepository.TransferByCategory(ctx, category.ID, departmentId, fromDate, toDate, "STORAGE_TO_DEPARTMENT")
		if err != nil {
			return nil, err
		}
		fromStorage[category.ID] = equipment

		equipment, err = s.EquipmentRepository.TransferByCategory(ctx, category.ID, departmentId, fromDate, toDate, "DEPARTMENT_TO_STORAGE")
		if err != nil {
			return nil, err
		}
		toStorage[category.ID] = equipment

		equipment, err = s.EquipmentRepository.TransferByCategory(ctx, category.ID, departmentId, fromDate, toDate, "CONTRACT_TO_DEPARTMENT")
		if err != nil {
			return nil, err
		}
		fromContract[category.ID] = equipment

		equipment, err = s.EquipmentRepository.TransferByCategory(ctx, category.ID, departmentId, fromDate, toDate, "DEPARTMENT_TO_CONTRACT")
		if err != nil {
			return nil, err
		}
		toContract[category.ID] = equipment

		equipment, err = s.EquipmentRepository.FromDepartmentTransferByCategory(ctx, category.ID, departmentId, fromDate, toDate)
		if err != nil {
			return nil, err
		}

		locationFromDepartment := make(map[int64][]*model.Location)
		for _, row := range equipment {
			departments[row.FromDepartment.ID] = row.FromDepartment
			locationFromDepartment[row.FromDepartment.ID] = append(locationFromDepartment[row.FromDepartment.ID], row)
		}
		fromDepartment[category.ID] = locationFromDepartment
		equipmentFrom, err := s.EquipmentRepository.ToDepartmentTransferByCategory(ctx, category.ID, departmentId, fromDate, toDate)
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
