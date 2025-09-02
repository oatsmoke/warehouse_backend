// Package service implements business logic for working with categories.
package service

import (
	"context"

	"github.com/oatsmoke/warehouse_backend/internal/lib/logger"
	"github.com/oatsmoke/warehouse_backend/internal/model"
	"github.com/oatsmoke/warehouse_backend/internal/repository"
)

// CategoryService provides methods for managing categories.
type CategoryService struct {
	CategoryRepository repository.Category
}

// NewCategoryService creates a new CategoryService.
// categoryRepository: repository implementation for category operations.
// Returns a pointer to CategoryService.
func NewCategoryService(categoryRepository repository.Category) *CategoryService {
	return &CategoryService{
		CategoryRepository: categoryRepository,
	}
}

// Create adds a new category with the specified title.
// ctx: request context.
// title: category name.
// Returns an error if the operation fails.
func (s *CategoryService) Create(ctx context.Context, title string) error {
	if err := s.CategoryRepository.Create(ctx, title); err != nil {
		return logger.Err(err, "")
	}

	return nil
}

// Update changes the title of an existing category by its ID.
// ctx: request context.
// id: category ID.
// title: new category name.
// Returns an error if the operation fails.
func (s *CategoryService) Update(ctx context.Context, id int64, title string) error {
	if err := s.CategoryRepository.Update(ctx, id, title); err != nil {
		return logger.Err(err, "")
	}

	return nil
}

// Delete marks the category as deleted by its ID (soft delete).
// ctx: request context.
// id: category ID.
// Returns an error if the operation fails.
func (s *CategoryService) Delete(ctx context.Context, id int64) error {
	if err := s.CategoryRepository.Delete(ctx, id); err != nil {
		return logger.Err(err, "")
	}

	return nil
}

// Restore unmarks the category as deleted by its ID.
// ctx: request context.
// id: category ID.
// Returns an error if the operation fails.
func (s *CategoryService) Restore(ctx context.Context, id int64) error {
	if err := s.CategoryRepository.Restore(ctx, id); err != nil {
		return logger.Err(err, "")
	}

	return nil
}

// GetAll returns a list of all categories.
// ctx: request context.
// deleted: if true, includes deleted categories.
// Returns a slice of Category pointers and an error if the operation fails.
func (s *CategoryService) GetAll(ctx context.Context, deleted bool) ([]*model.Category, error) {
	res, err := s.CategoryRepository.GetAll(ctx, deleted)
	if err != nil {
		return nil, logger.Err(err, "")
	}

	return res, nil
}

// GetById returns a category by its ID.
// ctx: request context.
// id: category ID.
// Returns a pointer to Category and an error if the operation fails.
func (s *CategoryService) GetById(ctx context.Context, id int64) (*model.Category, error) {
	res, err := s.CategoryRepository.GetById(ctx, id)
	if err != nil {
		return nil, logger.Err(err, "")
	}

	return res, nil
}
