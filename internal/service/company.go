package service

import (
	"context"
	"warehouse_backend/internal/lib/logger"
	"warehouse_backend/internal/model"
	"warehouse_backend/internal/repository"
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
	const fn = "service.Company.Create"

	if err := s.CompanyRepository.Create(ctx, title); err != nil {
		return logger.Err(err, "", fn)
	}

	return nil
}

// Update is a company update
func (s *CompanyService) Update(ctx context.Context, id int64, title string) error {
	const fn = "service.Company.Update"

	if err := s.CompanyRepository.Update(ctx, id, title); err != nil {
		return logger.Err(err, "", fn)
	}

	return nil
}

// Delete is a company delete
func (s *CompanyService) Delete(ctx context.Context, id int64) error {
	const fn = "service.Company.Delete"

	if err := s.CompanyRepository.Delete(ctx, id); err != nil {
		return logger.Err(err, "", fn)
	}

	return nil
}

// Restore is a company restore
func (s *CompanyService) Restore(ctx context.Context, id int64) error {
	const fn = "service.Company.Restore"

	if err := s.CompanyRepository.Restore(ctx, id); err != nil {
		return logger.Err(err, "", fn)
	}

	return nil
}

// GetAll is to get all companies
func (s *CompanyService) GetAll(ctx context.Context, deleted bool) ([]*model.Company, error) {
	const fn = "service.Company.GetAll"

	res, err := s.CompanyRepository.GetAll(ctx, deleted)
	if err != nil {
		return nil, logger.Err(err, "", fn)
	}

	return res, nil
}

// GetById is to get company by id
func (s *CompanyService) GetById(ctx context.Context, id int64) (*model.Company, error) {
	const fn = "service.Company.GetById"

	company := &model.Company{
		ID: id,
	}

	res, err := s.CompanyRepository.GetById(ctx, company)
	if err != nil {
		return nil, logger.Err(err, "", fn)
	}

	return res, nil
}
