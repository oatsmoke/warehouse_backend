package service

import (
	"context"
	"errors"
	"warehouse_backend/internal/model"
	"warehouse_backend/internal/repository"
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

func (s *ContractService) Create(ctx context.Context, number, address string) error {
	if _, err := s.repositoryContract.FindByNumber(ctx, number); err == nil {
		return errors.New("number already exists")
	}

	return s.repositoryContract.Create(ctx, number, address)
}

func (s *ContractService) GetById(ctx context.Context, id int64) (*model.Contract, error) {
	return s.repositoryContract.GetById(ctx, id)
}

func (s *ContractService) GetAll(ctx context.Context) ([]*model.Contract, error) {
	return s.repositoryContract.GetAll(ctx)
}

func (s *ContractService) Update(ctx context.Context, id int64, number, address string) error {
	findId, err := s.repositoryContract.FindByNumber(ctx, number)
	if findId != id && err == nil {
		return errors.New("number already exists")
	}

	return s.repositoryContract.Update(ctx, id, number, address)
}

func (s *ContractService) Delete(ctx context.Context, id int64) error {
	equipments, err := s.repositoryEquipment.GetByLocationContract(ctx, id)
	if err != nil {
		return err
	}

	if len(equipments) > 0 {
		return errors.New("have equipment")
	}

	return s.repositoryContract.Delete(ctx, id)
}
