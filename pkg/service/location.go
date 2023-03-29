package service

import (
	"strings"
	"warehouse_backend/pkg/model"
	"warehouse_backend/pkg/repository"
)

type LocationService struct {
	repositoryLocation repository.Location
}

func NewLocationService(repositoryLocation repository.Location) *LocationService {
	return &LocationService{repositoryLocation: repositoryLocation}
}

func (s *LocationService) TransferTo(id int, requests []model.RequestLocation) error {
	var code string
	for _, request := range requests {
		switch {
		case request.ToDepartment == 0 && request.ToEmployee == 0 && request.ToContract == 0:
			code = strings.ToUpper(request.ThisLocation + "_to_storage")
			err := s.repositoryLocation.TransferToStorage(request.Date, code, request.EquipmentId, id)
			if err != nil {
				return err
			}
		case request.ToDepartment != 0 && request.ToEmployee == 0 && request.ToContract == 0:
			if request.InDepartment {
				code = strings.ToUpper(request.ThisLocation + "_to_" + request.Where + "_in_department")
			} else {
				code = strings.ToUpper(request.ThisLocation + "_to_" + request.Where)
			}
			err := s.repositoryLocation.TransferToDepartment(request.Date, code, request.EquipmentId, id, request.ToDepartment)
			if err != nil {
				return err
			}
		case request.ToDepartment == 0 && request.ToEmployee != 0 && request.ToContract == 0:
			code = strings.ToUpper(request.ThisLocation + "_to_" + request.Where)
			err := s.repositoryLocation.TransferToEmployee(request.Date, code, request.EquipmentId, id, request.ToEmployee)
			if err != nil {
				return err
			}
		case request.ToDepartment != 0 && request.ToEmployee != 0 && request.ToContract == 0:
			code = strings.ToUpper(request.ThisLocation + "_to_" + request.Where + "_in_department")
			err := s.repositoryLocation.TransferToEmployeeInDepartment(request.Date, code, request.EquipmentId, id, request.ToDepartment, request.ToEmployee)
			if err != nil {
				return err
			}
		case request.ToDepartment == 0 && request.ToEmployee == 0 && request.ToContract != 0:
			code = strings.ToUpper(request.ThisLocation + "_to_" + request.Where)
			err := s.repositoryLocation.TransferToContract(request.Date, code, request.EquipmentId, id, request.ToContract)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
func (s *LocationService) GetHistory(id int) ([]model.Location, error) {
	return s.repositoryLocation.GetHistory(id)
}

func (s *LocationService) Delete(id int) error {
	return s.repositoryLocation.Delete(id)
}
