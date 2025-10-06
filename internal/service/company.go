package service

import (
	"context"
	"fmt"

	"github.com/oatsmoke/warehouse_backend/internal/dto"
	"github.com/oatsmoke/warehouse_backend/internal/lib/logger"
	"github.com/oatsmoke/warehouse_backend/internal/model"
	"github.com/oatsmoke/warehouse_backend/internal/repository"
)

type CompanyService struct {
	companyRepository repository.Company
}

func NewCompanyService(companyRepository repository.Company) *CompanyService {
	return &CompanyService{
		companyRepository: companyRepository,
	}
}

func (s *CompanyService) Create(ctx context.Context, company *model.Company) error {
	id, err := s.companyRepository.Create(ctx, company)
	if err != nil {
		return err
	}

	logger.InfoInConsole(fmt.Sprintf("company with id %d created", id))
	return nil
}

func (s *CompanyService) Read(ctx context.Context, id int64) (*model.Company, error) {
	read, err := s.companyRepository.Read(ctx, id)
	if err != nil {
		return nil, err
	}

	logger.InfoInConsole(fmt.Sprintf("company with id %d read", id))
	return read, nil
}

func (s *CompanyService) Update(ctx context.Context, company *model.Company) error {
	if err := s.companyRepository.Update(ctx, company); err != nil {
		return err
	}

	logger.InfoInConsole(fmt.Sprintf("company with id %d updated", company.ID))
	return nil
}

func (s *CompanyService) Delete(ctx context.Context, id int64) error {
	if err := s.companyRepository.Delete(ctx, id); err != nil {
		return err
	}

	logger.InfoInConsole(fmt.Sprintf("company with id %d deleted", id))
	return nil
}

func (s *CompanyService) Restore(ctx context.Context, id int64) error {
	if err := s.companyRepository.Restore(ctx, id); err != nil {
		return err
	}

	logger.InfoInConsole(fmt.Sprintf("company with id %d restored", id))
	return nil
}

func (s *CompanyService) List(ctx context.Context, qp *dto.QueryParams) ([]*model.Company, error) {
	list, err := s.companyRepository.List(ctx, qp)
	if err != nil {
		return nil, err
	}

	logger.InfoInConsole(fmt.Sprintf("%d company listed", len(list)))
	return list, nil
}
