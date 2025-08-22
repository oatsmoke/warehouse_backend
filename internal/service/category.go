package service

import (
	"context"

	"github.com/oatsmoke/warehouse_backend/internal/lib/logger"
	"github.com/oatsmoke/warehouse_backend/internal/model"
	"github.com/oatsmoke/warehouse_backend/internal/repository"
)

type CategoryService struct {
	CategoryRepository repository.Category
}

func NewCategoryService(categoryRepository repository.Category) *CategoryService {
	return &CategoryService{
		CategoryRepository: categoryRepository,
	}
}

// Create is category create
func (s *CategoryService) Create(ctx context.Context, title string) error {
	if err := s.CategoryRepository.Create(ctx, title); err != nil {
		return logger.Err(err, "")
	}

	return nil
}

// Update is category update
func (s *CategoryService) Update(ctx context.Context, id int64, title string) error {
	if err := s.CategoryRepository.Update(ctx, id, title); err != nil {
		return logger.Err(err, "")
	}

	return nil
}

// Delete is category delete
func (s *CategoryService) Delete(ctx context.Context, id int64) error {
	if err := s.CategoryRepository.Delete(ctx, id); err != nil {
		return logger.Err(err, "")
	}

	return nil
}

// Restore is category restore
func (s *CategoryService) Restore(ctx context.Context, id int64) error {
	if err := s.CategoryRepository.Restore(ctx, id); err != nil {
		return logger.Err(err, "")
	}

	return nil
}

// GetAll is to get all categories
func (s *CategoryService) GetAll(ctx context.Context, deleted bool) ([]*model.Category, error) {
	res, err := s.CategoryRepository.GetAll(ctx, deleted)
	if err != nil {
		return nil, logger.Err(err, "")
	}

	return res, nil
}

// GetById is to get category by id
func (s *CategoryService) GetById(ctx context.Context, id int64) (*model.Category, error) {
	category := &model.Category{
		ID: id,
	}

	res, err := s.CategoryRepository.GetById(ctx, category)
	if err != nil {
		return nil, logger.Err(err, "")
	}

	return res, nil
}
