package service

import (
	"context"

	"github.com/oatsmoke/warehouse_backend/internal/model"
	"github.com/oatsmoke/warehouse_backend/internal/repository"
)

type CategoryService struct {
	categoryRepository repository.Category
}

func NewCategoryService(categoryRepository repository.Category) *CategoryService {
	return &CategoryService{
		categoryRepository: categoryRepository,
	}
}

func (s *CategoryService) Create(ctx context.Context, category *model.Category) error {
	return s.categoryRepository.Create(ctx, category)
}

func (s *CategoryService) Read(ctx context.Context, id int64) (*model.Category, error) {
	return s.categoryRepository.Read(ctx, id)
}

func (s *CategoryService) Update(ctx context.Context, category *model.Category) error {
	return s.categoryRepository.Update(ctx, category)
}

func (s *CategoryService) Delete(ctx context.Context, id int64) error {
	return s.categoryRepository.Delete(ctx, id)
}

func (s *CategoryService) Restore(ctx context.Context, id int64) error {
	return s.categoryRepository.Restore(ctx, id)
}

func (s *CategoryService) List(ctx context.Context, withDeleted bool) ([]*model.Category, error) {
	return s.categoryRepository.List(ctx, withDeleted)
}
