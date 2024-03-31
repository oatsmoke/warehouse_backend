package service

import (
	"strings"
	"warehouse_backend/pkg/model"
	"warehouse_backend/pkg/repository"
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

func (s *LocationService) TransferTo(id int, requests []model.RequestLocation) error {
	var code string
	replace := []int{0, 0}
	for _, request := range requests {
		nowLocation, _ := s.repositoryLocation.GetLocationNow(request.EquipmentId)
		switch {
		case request.ToDepartment == 0 && request.ToEmployee == 0 && request.ToContract == 0:
			code = strings.ToUpper(request.ThisLocation + "_to_storage")
			transferId, err := s.repositoryLocation.TransferToStorage(request.Date, code, request.EquipmentId, id, request.Company, nowLocation)
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
			transferId, err := s.repositoryLocation.TransferToDepartment(request.Date, code, request.EquipmentId, id, request.Company, request.ToDepartment, nowLocation)
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
			transferId, err := s.repositoryLocation.TransferToEmployee(request.Date, code, request.EquipmentId, id, request.Company, request.ToEmployee, nowLocation)
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
			transferId, err := s.repositoryLocation.TransferToEmployeeInDepartment(request.Date, code, request.EquipmentId, id, request.Company, request.ToDepartment, request.ToEmployee, nowLocation)
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
			transferId, err := s.repositoryLocation.TransferToContract(request.Date, code, request.EquipmentId, id, request.Company, request.ToContract, request.TransferType, request.Price, nowLocation)
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
		err := s.repositoryReplace.Create(replace)
		if err != nil {
			return err
		}
	}
	return nil
}
func (s *LocationService) GetHistory(id int) ([]model.Location, error) {
	return s.repositoryLocation.GetHistory(id)
}

func (s *LocationService) Delete(id int) error {
	replace, err := s.repositoryReplace.FindByLocationId(id)
	if err != nil {
		return s.repositoryLocation.Delete(id)
	}
	err = s.repositoryLocation.Delete(replace.TransferFrom)
	err = s.repositoryLocation.Delete(replace.TransferTo)
	return err
}
