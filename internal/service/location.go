package service

import (
	"context"
	"strings"
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
