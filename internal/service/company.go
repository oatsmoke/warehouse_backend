package service

import (
	"context"

	"github.com/oatsmoke/warehouse_backend/internal/lib/logger"
	"github.com/oatsmoke/warehouse_backend/internal/model"
	"github.com/oatsmoke/warehouse_backend/internal/repository"
)

type CompanyService struct {
	CompanyRepository repository.Company
}

func NewCompanyService(companyRepository repository.Company) *CompanyService {
	return &CompanyService{
		CompanyRepository: companyRepository,
	}
}

// Create is company create
func (s *CompanyService) Create(ctx context.Context, title string) error {
	if err := s.CompanyRepository.Create(ctx, title); err != nil {
		return logger.Err(err, "")
	}

	return nil
}

// Update is a company update
func (s *CompanyService) Update(ctx context.Context, id int64, title string) error {
	if err := s.CompanyRepository.Update(ctx, id, title); err != nil {
		return logger.Err(err, "")
	}

	return nil
}

// Delete is a company delete
func (s *CompanyService) Delete(ctx context.Context, id int64) error {
	if err := s.CompanyRepository.Delete(ctx, id); err != nil {
		return logger.Err(err, "")
	}

	return nil
}

// Restore is a company restore
func (s *CompanyService) Restore(ctx context.Context, id int64) error {
	if err := s.CompanyRepository.Restore(ctx, id); err != nil {
		return logger.Err(err, "")
	}

	return nil
}

// GetAll is to get all companies
func (s *CompanyService) GetAll(ctx context.Context, deleted bool) ([]*model.Company, error) {
	res, err := s.CompanyRepository.GetAll(ctx, deleted)
	if err != nil {
		return nil, logger.Err(err, "")
	}

	return res, nil
}

// GetById is to get company by id
func (s *CompanyService) GetById(ctx context.Context, id int64) (*model.Company, error) {
	company := &model.Company{
		ID: id,
	}

	res, err := s.CompanyRepository.GetById(ctx, company)
	if err != nil {
		return nil, logger.Err(err, "")
	}

	return res, nil
}
