package service

import (
	"errors"
	"strings"
	"warehouse_backend/pkg/model"
	"warehouse_backend/pkg/repository"
)

type EquipmentService struct {
	repositoryEquipment repository.Equipment
}

func NewEquipmentService(repositoryEquipment repository.Equipment) *EquipmentService {
	return &EquipmentService{repositoryEquipment: repositoryEquipment}
}

func (s *EquipmentService) Create(date int64, serialNumber string, profile int, userId int) (int, error) {
	serialNumber = strings.ToUpper(serialNumber)
	if _, err := s.repositoryEquipment.FindBySerialNumber(serialNumber); err == nil {
		return 0, errors.New("serial number already exists")
	}
	return s.repositoryEquipment.Create(date, serialNumber, profile, userId)
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
