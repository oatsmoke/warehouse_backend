package service

import (
	"context"
	"strings"
	"time"
	"warehouse_backend/internal/lib/logger"
	"warehouse_backend/internal/model"
	"warehouse_backend/internal/repository"
)

type LocationService struct {
	LocationRepository repository.Location
	ReplaceRepository  repository.Replace
	CategoryRepository repository.Category
}

func NewLocationService(locationRepository repository.Location, replaceRepository repository.Replace, categoryRepository repository.Category) *LocationService {
	return &LocationService{
		LocationRepository: locationRepository,
		ReplaceRepository:  replaceRepository,
		CategoryRepository: categoryRepository,
	}
}

// AddToStorage is equipment add to storage
func (s *LocationService) AddToStorage(ctx context.Context, date string, equipmentId, employeeId, companyId int64) error {
	const fn = "service.Location.AddToStorage"

	if err := s.LocationRepository.AddToStorage(ctx, date, equipmentId, employeeId, companyId); err != nil {
		return logger.Err(err, "", fn)
	}

	return nil
}

// TransferTo is equipment transfer to
func (s *LocationService) TransferTo(ctx context.Context, EmployeeId int64, requests []*model.RequestLocation) error {
	const fn = "service.Location.TransferTo"

	var code string
	replace := []int64{0, 0}
	for _, request := range requests {
		nowLocation, _ := s.LocationRepository.GetLocationNow(ctx, request.EquipmentId)
		switch {

		case request.ToDepartment == 0 && request.ToEmployee == 0 && request.ToContract == 0:
			code = strings.ToUpper(request.ThisLocation + "_to_storage")
			transferId, err := s.LocationRepository.TransferToStorage(ctx, request.Date, code, request.EquipmentId, EmployeeId, request.Company, nowLocation)
			if err != nil {
				return logger.Err(err, "", fn)
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
			transferId, err := s.LocationRepository.TransferToDepartment(ctx, request.Date, code, request.EquipmentId, EmployeeId, request.Company, request.ToDepartment, nowLocation)
			if err != nil {
				return logger.Err(err, "", fn)
			}
			if request.Way == "replace" && replace[0] == 0 {
				replace[0] = transferId
			} else if request.Way == "replace" && replace[1] == 0 {
				replace[1] = transferId
			}

		case request.ToDepartment == 0 && request.ToEmployee != 0 && request.ToContract == 0:
			code = strings.ToUpper(request.ThisLocation + "_to_" + request.Where)
			transferId, err := s.LocationRepository.TransferToEmployee(ctx, request.Date, code, request.EquipmentId, EmployeeId, request.Company, request.ToEmployee, nowLocation)
			if err != nil {
				return logger.Err(err, "", fn)
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
			transferId, err := s.LocationRepository.TransferToEmployeeInDepartment(ctx, request.Date, code, request.EquipmentId, EmployeeId, request.Company, request.ToDepartment, request.ToEmployee, nowLocation)
			if err != nil {
				return logger.Err(err, "", fn)
			}
			if request.Way == "replace" && replace[0] == 0 {
				replace[0] = transferId
			} else if request.Way == "replace" && replace[1] == 0 {
				replace[1] = transferId
			}

		case request.ToDepartment == 0 && request.ToEmployee == 0 && request.ToContract != 0:
			code = strings.ToUpper(request.ThisLocation + "_to_" + request.Where)
			transferId, err := s.LocationRepository.TransferToContract(ctx, request.Date, code, request.EquipmentId, EmployeeId, request.Company, request.ToContract, request.TransferType, request.Price, nowLocation)
			if err != nil {
				return logger.Err(err, "", fn)
			}
			if request.Way == "replace" && replace[0] == 0 {
				replace[0] = transferId
			} else if request.Way == "replace" && replace[1] == 0 {
				replace[1] = transferId
			}
		}
	}

	if requests[0].Way == "replace" {
		err := s.ReplaceRepository.Create(ctx, replace)
		if err != nil {
			return logger.Err(err, "", fn)
		}
	}

	return nil
}

// Delete is equipment transfer to
func (s *LocationService) Delete(ctx context.Context, id int64) error {
	fn := "service.Location.Delete"

	replace, err := s.ReplaceRepository.FindByLocationId(ctx, id)
	if err != nil {
		return s.LocationRepository.Delete(ctx, id)
	}

	if err = s.LocationRepository.Delete(ctx, replace.TransferFrom); err != nil {
		return logger.Err(err, "", fn)
	}

	if err = s.LocationRepository.Delete(ctx, replace.TransferTo); err != nil {
		return logger.Err(err, "", fn)
	}

	return nil
}

// GetById is equipment get by id
func (s *LocationService) GetById(ctx context.Context, equipmentId int64) (*model.Location, error) {
	fn := "service.Location.GetById"

	res, err := s.LocationRepository.GetById(ctx, equipmentId)
	if err != nil {
		return nil, logger.Err(err, "", fn)
	}

	return res, nil
}

// GetByIds is equipment get by ids
func (s *LocationService) GetByIds(ctx context.Context, equipmentIds []int64) ([]*model.Location, error) {
	fn := "service.Location.GetByIds"
	var equipments []*model.Location

	for _, id := range equipmentIds {
		equipment, err := s.LocationRepository.GetById(ctx, id)
		if err != nil {
			return nil, logger.Err(err, "", fn)
		}
		equipments = append(equipments, equipment)
	}
	return equipments, nil
}

// GetHistory is equipment get history
func (s *LocationService) GetHistory(ctx context.Context, equipmentId int64) ([]*model.Location, error) {
	fn := "service.Location.GetHistory"

	res, err := s.LocationRepository.GetHistory(ctx, equipmentId)
	if err != nil {
		return nil, logger.Err(err, "", fn)
	}
	return res, nil
}

// GetByLocation is equipment get by location
func (s *LocationService) GetByLocation(ctx context.Context, toDepartmentId, toEmployeeId, toContractId int64) ([]*model.Location, error) {
	fn := "service.Location.GetByLocation"

	switch {
	case toDepartmentId != 0:
		res, err := s.LocationRepository.GetByLocationDepartment(ctx, toDepartmentId)
		if err != nil {
			return nil, logger.Err(err, "", fn)
		}
		return res, nil

	case toEmployeeId != 0:
		res, err := s.LocationRepository.GetByLocationEmployee(ctx, toDepartmentId)
		if err != nil {
			return nil, logger.Err(err, "", fn)
		}
		return res, nil

	case toContractId != 0:
		res, err := s.LocationRepository.GetByLocationContract(ctx, toDepartmentId)
		if err != nil {
			return nil, logger.Err(err, "", fn)
		}
		return res, nil

	default:
		res, err := s.LocationRepository.GetByLocationStorage(ctx)
		if err != nil {
			return nil, logger.Err(err, "", fn)
		}
		return res, nil
	}
}

// ReportByCategory is equipment report by category
func (s *LocationService) ReportByCategory(ctx context.Context, departmentId int64, date string) (*model.Report, error) {
	fn := "service.Location.ReportByCategory"

	report := new(model.Report)
	fromDate := date
	parseTime, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return nil, logger.Err(err, "", fn)
	}
	toDate := parseTime.AddDate(0, 1, 0).String()
	categories, err := s.CategoryRepository.GetAll(ctx, false)
	if err != nil {
		return nil, logger.Err(err, "", fn)
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
		equipment, err := s.LocationRepository.RemainderByCategory(ctx, category.ID, departmentId, fromDate)
		if err != nil {
			return nil, logger.Err(err, "", fn)
		}
		leftover[category.ID] = equipment

		equipment, err = s.LocationRepository.RemainderByCategory(ctx, category.ID, departmentId, toDate)
		if err != nil {
			return nil, logger.Err(err, "", fn)
		}
		total[category.ID] = equipment

		equipment, err = s.LocationRepository.TransferByCategory(ctx, category.ID, departmentId, fromDate, toDate, "STORAGE_TO_DEPARTMENT")
		if err != nil {
			return nil, logger.Err(err, "", fn)
		}
		fromStorage[category.ID] = equipment

		equipment, err = s.LocationRepository.TransferByCategory(ctx, category.ID, departmentId, fromDate, toDate, "DEPARTMENT_TO_STORAGE")
		if err != nil {
			return nil, logger.Err(err, "", fn)
		}
		toStorage[category.ID] = equipment

		equipment, err = s.LocationRepository.TransferByCategory(ctx, category.ID, departmentId, fromDate, toDate, "CONTRACT_TO_DEPARTMENT")
		if err != nil {
			return nil, logger.Err(err, "", fn)
		}
		fromContract[category.ID] = equipment

		equipment, err = s.LocationRepository.TransferByCategory(ctx, category.ID, departmentId, fromDate, toDate, "DEPARTMENT_TO_CONTRACT")
		if err != nil {
			return nil, logger.Err(err, "", fn)
		}
		toContract[category.ID] = equipment

		equipment, err = s.LocationRepository.FromDepartmentTransferByCategory(ctx, category.ID, departmentId, fromDate, toDate)
		if err != nil {
			return nil, logger.Err(err, "", fn)
		}

		locationFromDepartment := make(map[int64][]*model.Location)
		for _, row := range equipment {
			departments[row.FromDepartment.ID] = row.FromDepartment
			locationFromDepartment[row.FromDepartment.ID] = append(locationFromDepartment[row.FromDepartment.ID], row)
		}
		fromDepartment[category.ID] = locationFromDepartment
		equipmentFrom, err := s.LocationRepository.ToDepartmentTransferByCategory(ctx, category.ID, departmentId, fromDate, toDate)
		if err != nil {
			return nil, logger.Err(err, "", fn)
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
