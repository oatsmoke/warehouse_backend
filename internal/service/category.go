package service

import (
	"context"
	"warehouse_backend/internal/lib/logger"
	"warehouse_backend/internal/model"
	"warehouse_backend/internal/repository"
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
	const fn = "service.Category.Create"

	if err := s.CategoryRepository.Create(ctx, title); err != nil {
		return logger.Err(err, "", fn)
	}

	return nil
}

// Update is category update
func (s *CategoryService) Update(ctx context.Context, id int64, title string) error {
	const fn = "service.Category.Update"

	if err := s.CategoryRepository.Update(ctx, id, title); err != nil {
		return logger.Err(err, "", fn)
	}

	return nil
}

// Delete is category delete
func (s *CategoryService) Delete(ctx context.Context, id int64) error {
	const fn = "service.Category.Delete"

	if err := s.CategoryRepository.Delete(ctx, id); err != nil {
		return logger.Err(err, "", fn)
	}

	return nil
}

// Restore is category restore
func (s *CategoryService) Restore(ctx context.Context, id int64) error {
	const fn = "service.Category.Restore"

	if err := s.CategoryRepository.Restore(ctx, id); err != nil {
		return logger.Err(err, "", fn)
	}

	return nil
}

// GetAll is to get all categories
func (s *CategoryService) GetAll(ctx context.Context, deleted bool) ([]*model.Category, error) {
	const fn = "service.Category.GetAll"

	categories, err := s.CategoryRepository.GetAll(ctx, deleted)
	if err != nil {
		return nil, logger.Err(err, "", fn)
	}

	return categories, nil
}

// GetById is to get category by id
func (s *CategoryService) GetById(ctx context.Context, id int64) (*model.Category, error) {
	const fn = "service.Category.GetById"

	category, err := s.CategoryRepository.GetById(ctx, &model.Category{ID: id})
	if err != nil {
		return nil, logger.Err(err, "", fn)
	}

	return category, nil
}
