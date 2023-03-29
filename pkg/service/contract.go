package service

import (
	"errors"
	"warehouse_backend/pkg/model"
	"warehouse_backend/pkg/repository"
)

type ContractService struct {
	repositoryContract  repository.Contract
	repositoryEquipment repository.Equipment
}

func NewContractService(repositoryContract repository.Contract,
	repositoryEquipment repository.Equipment) *ContractService {
	return &ContractService{repositoryContract: repositoryContract,
		repositoryEquipment: repositoryEquipment,
	}
}

func (s *ContractService) Create(number, address string) error {
	if _, err := s.repositoryContract.FindByNumber(number); err == nil {
		return errors.New("number already exists")
	}
	return s.repositoryContract.Create(number, address)
}

func (s *ContractService) GetById(id int) (model.Contract, error) {
	if id != 0 {
		return s.repositoryContract.GetById(id)
	}
	return model.Contract{}, nil
}

func (s *ContractService) GetAll() ([]model.Contract, error) {
	return s.repositoryContract.GetAll()
}

func (s *ContractService) Update(id int, number, address string) error {
	findId, err := s.repositoryContract.FindByNumber(number)
	if findId != id && err == nil {
		return errors.New("number already exists")
	}
	return s.repositoryContract.Update(id, number, address)
}

func (s *ContractService) Delete(id int) error {
	equipments, err := s.repositoryEquipment.GetByLocationContract(id)
	if err != nil {
		return err
	}
	if len(equipments) > 0 {
		return errors.New("have equipment")
	}
	return s.repositoryContract.Delete(id)
}
