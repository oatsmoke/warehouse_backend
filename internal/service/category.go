package service

import (
	"context"
	"fmt"

	"github.com/oatsmoke/warehouse_backend/internal/dto"
	"github.com/oatsmoke/warehouse_backend/internal/lib/logger"
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
	id, err := s.categoryRepository.Create(ctx, category)
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("category with id %d created", id))
	return nil
}

func (s *CategoryService) Read(ctx context.Context, id int64) (*model.Category, error) {
	read, err := s.categoryRepository.Read(ctx, id)
	if err != nil {
		return nil, err
	}

	logger.Info(fmt.Sprintf("category with id %d read", id))
	return read, err
}

func (s *CategoryService) Update(ctx context.Context, category *model.Category) error {
	if err := s.categoryRepository.Update(ctx, category); err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("category with id %d updated", category.ID))
	return nil
}

func (s *CategoryService) Delete(ctx context.Context, id int64) error {
	if err := s.categoryRepository.Delete(ctx, id); err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("category with id %d deleted", id))
	return nil
}

func (s *CategoryService) Restore(ctx context.Context, id int64) error {
	if err := s.categoryRepository.Restore(ctx, id); err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("category with id %d restored", id))
	return nil
}

func (s *CategoryService) List(ctx context.Context, qp *dto.QueryParams) (*dto.ListResponse[[]*model.Category], error) {
	list, total, err := s.categoryRepository.List(ctx, qp)
	if err != nil {
		return nil, err
	}

	logger.Info(fmt.Sprintf("%d category listed", len(list)))
	return &dto.ListResponse[[]*model.Category]{
		List:  list,
		Total: total,
	}, nil
}
