package service

import (
	"context"

	"github.com/oatsmoke/warehouse_backend/internal/lib/logger"
	"github.com/oatsmoke/warehouse_backend/internal/model"
	"github.com/oatsmoke/warehouse_backend/internal/repository"
)

type ContractService struct {
	ContractRepository repository.Contract
}

func NewContractService(contractRepository repository.Contract) *ContractService {
	return &ContractService{
		ContractRepository: contractRepository,
	}
}

// Create is contract create
func (s *ContractService) Create(ctx context.Context, number, address string) error {
	if err := s.ContractRepository.Create(ctx, number, address); err != nil {
		return logger.Err(err, "")
	}

	return nil
}

// Update is contract create
func (s *ContractService) Update(ctx context.Context, id int64, number, address string) error {
	if err := s.ContractRepository.Update(ctx, id, number, address); err != nil {
		return logger.Err(err, "")
	}

	return nil
}

// Delete is contract delete
func (s *ContractService) Delete(ctx context.Context, id int64) error {
	if err := s.ContractRepository.Delete(ctx, id); err != nil {
		return logger.Err(err, "")
	}

	return nil
}

// Restore is contract restore
func (s *ContractService) Restore(ctx context.Context, id int64) error {
	if err := s.ContractRepository.Restore(ctx, id); err != nil {
		return logger.Err(err, "")
	}

	return nil
}

// GetAll is to get all contracts
func (s *ContractService) GetAll(ctx context.Context, deleted bool) ([]*model.Contract, error) {
	res, err := s.ContractRepository.GetAll(ctx, deleted)
	if err != nil {
		return nil, logger.Err(err, "")
	}

	return res, nil
}

// GetById is to get contract by id
func (s *ContractService) GetById(ctx context.Context, id int64) (*model.Contract, error) {
	contract := &model.Contract{
		ID: id,
	}

	res, err := s.ContractRepository.GetById(ctx, contract)
	if err != nil {
		return nil, logger.Err(err, "")
	}

	return res, nil
}
