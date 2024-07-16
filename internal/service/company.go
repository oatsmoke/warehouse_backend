package service

import (
	"context"
	"errors"
	"warehouse_backend/internal/model"
	"warehouse_backend/internal/repository"
)

type CompanyService struct {
	repositoryCompany repository.Company
}

func NewCompanyService(repositoryCompany repository.Company) *CompanyService {
	return &CompanyService{
		repositoryCompany: repositoryCompany,
	}
}

func (s *CompanyService) Create(ctx context.Context, title string) error {
	if _, err := s.repositoryCompany.FindByTitle(ctx, title); err == nil {
		return errors.New("title already exists")
	}

	return s.repositoryCompany.Create(ctx, title)
}

func (s *CompanyService) GetById(ctx context.Context, id int64) (*model.Company, error) {
	return s.repositoryCompany.GetById(ctx, id)
}

func (s *CompanyService) GetAll(ctx context.Context) ([]*model.Company, error) {
	return s.repositoryCompany.GetAll(ctx)
}

func (s *CompanyService) Update(ctx context.Context, id int64, title string) error {
	if _, err := s.repositoryCompany.FindByTitle(ctx, title); err == nil {
		return errors.New("title already exists")
	}

	return s.repositoryCompany.Update(ctx, id, title)
}

func (s *CompanyService) Delete(ctx context.Context, id int64) error {
	//profiles, err := s.repositoryProfile.GetByCategory(id)
	//if err != nil {
	//	return err
	//}
	//if len(profiles) > 0 {
	//	return errors.New("used in profile")
	//}
	return s.repositoryCompany.Delete(ctx, id)
}
