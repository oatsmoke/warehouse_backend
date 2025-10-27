package service

import (
	"context"
	"fmt"

	"github.com/oatsmoke/warehouse_backend/internal/dto"
	"github.com/oatsmoke/warehouse_backend/internal/lib/logger"
	"github.com/oatsmoke/warehouse_backend/internal/model"
	"github.com/oatsmoke/warehouse_backend/internal/repository"
)

type ContractService struct {
	contractRepository repository.Contract
}

func NewContractService(contractRepository repository.Contract) *ContractService {
	return &ContractService{
		contractRepository: contractRepository,
	}
}

func (s *ContractService) Create(ctx context.Context, contract *model.Contract) error {
	id, err := s.contractRepository.Create(ctx, contract)
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("contract with id %d created", id))
	return nil
}

func (s *ContractService) Read(ctx context.Context, id int64) (*model.Contract, error) {
	read, err := s.contractRepository.Read(ctx, id)
	if err != nil {
		return nil, err
	}

	logger.Info(fmt.Sprintf("contract with id %d read", id))
	return read, nil
}

func (s *ContractService) Update(ctx context.Context, contract *model.Contract) error {
	if err := s.contractRepository.Update(ctx, contract); err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("contract with id %d updated", contract.ID))
	return nil
}

func (s *ContractService) Delete(ctx context.Context, id int64) error {
	if err := s.contractRepository.Delete(ctx, id); err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("contract with id %d deleted", id))
	return nil
}

func (s *ContractService) Restore(ctx context.Context, id int64) error {
	if err := s.contractRepository.Restore(ctx, id); err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("contract with id %d restored", id))
	return nil
}

func (s *ContractService) List(ctx context.Context, qp *dto.QueryParams) (*dto.ListResponse[[]*model.Contract], error) {
	list, total, err := s.contractRepository.List(ctx, qp)
	if err != nil {
		return nil, err
	}

	logger.Info(fmt.Sprintf("%d contract listed", len(list)))
	return &dto.ListResponse[[]*model.Contract]{
		List:  list,
		Total: total,
	}, nil
}
