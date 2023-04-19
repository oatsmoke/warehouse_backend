package service

import (
	"errors"
	"warehouse_backend/pkg/model"
	"warehouse_backend/pkg/repository"
)

type CompanyService struct {
	repositoryCompany repository.Company
}

func NewCompanyService(repositoryCompany repository.Company) *CompanyService {
	return &CompanyService{
		repositoryCompany: repositoryCompany,
	}
}

func (s *CompanyService) Create(title string) error {
	if _, err := s.repositoryCompany.FindByTitle(title); err == nil {
		return errors.New("title already exists")
	}
	return s.repositoryCompany.Create(title)
}

func (s *CompanyService) GetById(id int) (model.Company, error) {
	return s.repositoryCompany.GetById(id)
}

func (s *CompanyService) GetAll() ([]model.Company, error) {
	return s.repositoryCompany.GetAll()
}

func (s *CompanyService) Update(id int, title string) error {
	if _, err := s.repositoryCompany.FindByTitle(title); err == nil {
		return errors.New("title already exists")
	}
	return s.repositoryCompany.Update(id, title)
}

func (s *CompanyService) Delete(id int) error {
	//profiles, err := s.repositoryProfile.GetByCategory(id)
	//if err != nil {
	//	return err
	//}
	//if len(profiles) > 0 {
	//	return errors.New("used in profile")
	//}
	return s.repositoryCompany.Delete(id)
}
