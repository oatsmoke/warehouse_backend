package service

import (
	"errors"
	"warehouse_backend/pkg/model"
	"warehouse_backend/pkg/repository"
)

type CategoryService struct {
	repositoryCategory repository.Category
	repositoryProfile  repository.Profile
}

func NewCategoryService(repositoryCategory repository.Category,
	repositoryProfile repository.Profile) *CategoryService {
	return &CategoryService{
		repositoryCategory: repositoryCategory,
		repositoryProfile:  repositoryProfile,
	}
}

func (s *CategoryService) Create(title string) error {
	if _, err := s.repositoryCategory.FindByTitle(title); err == nil {
		return errors.New("title already exists")
	}
	return s.repositoryCategory.Create(title)
}

func (s *CategoryService) GetById(id int) (model.Category, error) {
	return s.repositoryCategory.GetById(id)
}

func (s *CategoryService) GetAll() ([]model.Category, error) {
	return s.repositoryCategory.GetAll()
}

func (s *CategoryService) Update(id int, title string) error {
	if _, err := s.repositoryCategory.FindByTitle(title); err == nil {
		return errors.New("title already exists")
	}
	return s.repositoryCategory.Update(id, title)
}

func (s *CategoryService) Delete(id int) error {
	profiles, err := s.repositoryProfile.GetByCategory(id)
	if err != nil {
		return err
	}
	if len(profiles) > 0 {
		return errors.New("used in profile")
	}
	return s.repositoryCategory.Delete(id)
}
