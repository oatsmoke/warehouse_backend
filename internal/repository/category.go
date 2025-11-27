package repository

import (
	"context"

	queries "github.com/oatsmoke/warehouse_backend/internal/db"
	"github.com/oatsmoke/warehouse_backend/internal/dto"
	"github.com/oatsmoke/warehouse_backend/internal/lib/logger"
	"github.com/oatsmoke/warehouse_backend/internal/model"
)

type CategoryRepository struct {
	queries queries.Querier
}

func NewCategoryRepository(queries queries.Querier) *CategoryRepository {
	return &CategoryRepository{
		queries: queries,
	}
}

func (r *CategoryRepository) Create(ctx context.Context, category *model.Category) (int64, error) {
	req, err := r.queries.CreateCategory(ctx, category.Title)
	if err != nil {
		return 0, logger.Error(logger.MsgFailedToInsert, err)
	}

	return req.ID, nil
}

func (r *CategoryRepository) Read(ctx context.Context, id int64) (*model.Category, error) {
	req, err := r.queries.ReadCategory(ctx, id)
	if err != nil {
		return nil, logger.Error(logger.MsgFailedToScan, err)
	}

	category := &model.Category{
		ID:        req.ID,
		Title:     req.Title,
		DeletedAt: validTime(req.DeletedAt),
	}

	return category, nil
}

func (r *CategoryRepository) Update(ctx context.Context, category *model.Category) error {
	ct, err := r.queries.UpdateCategory(ctx, &queries.UpdateCategoryParams{
		ID:    category.ID,
		Title: category.Title,
	})
	if err != nil {
		return logger.Error(logger.MsgFailedToUpdate, err)
	}

	if ct.RowsAffected() == 0 {
		return logger.Error(logger.MsgFailedToUpdate, logger.ErrNoRowsAffected)
	}

	return nil
}

func (r *CategoryRepository) Delete(ctx context.Context, id int64) error {
	ct, err := r.queries.DeleteCategory(ctx, id)
	if err != nil {
		return logger.Error(logger.MsgFailedToDelete, err)
	}

	if ct.RowsAffected() == 0 {
		return logger.Error(logger.MsgFailedToDelete, logger.ErrNoRowsAffected)
	}

	return nil
}

func (r *CategoryRepository) Restore(ctx context.Context, id int64) error {
	ct, err := r.queries.RestoreCategory(ctx, id)
	if err != nil {
		return logger.Error(logger.MsgFailedToRestore, err)
	}

	if ct.RowsAffected() == 0 {
		return logger.Error(logger.MsgFailedToRestore, logger.ErrNoRowsAffected)
	}

	return nil
}

func (r *CategoryRepository) List(ctx context.Context, qp *dto.QueryParams) ([]*model.Category, int64, error) {
	req, err := r.queries.ListCategory(ctx, &queries.ListCategoryParams{
		WithDeleted:      qp.WithDeleted,
		Search:           qp.Search,
		Ids:              qp.Ids,
		SortColumn:       qp.SortColumn,
		SortOrder:        qp.SortOrder,
		PaginationLimit:  qp.PaginationLimit,
		PaginationOffset: qp.PaginationOffset,
	})
	if err != nil {
		return nil, 0, logger.Error(logger.MsgFailedToSelect, err)
	}

	if len(req) < 1 {
		return []*model.Category{}, 0, nil
	}

	list := make([]*model.Category, len(req))
	for i, item := range req {
		category := &model.Category{
			ID:        item.ID,
			Title:     item.Title,
			DeletedAt: validTime(item.DeletedAt),
		}
		list[i] = category
	}

	return list, req[0].Total, nil
}
