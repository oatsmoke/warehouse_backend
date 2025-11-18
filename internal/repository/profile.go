package repository

import (
	"context"

	queries "github.com/oatsmoke/warehouse_backend/internal/db"
	"github.com/oatsmoke/warehouse_backend/internal/dto"
	"github.com/oatsmoke/warehouse_backend/internal/lib/logger"
	"github.com/oatsmoke/warehouse_backend/internal/model"
)

type ProfileRepository struct {
	queries queries.Querier
}

func NewProfileRepository(queries queries.Querier) *ProfileRepository {
	return &ProfileRepository{
		queries: queries,
	}
}

func (r *ProfileRepository) Create(ctx context.Context, profile *model.Profile) (int64, error) {
	req, err := r.queries.CreateProfile(ctx, &queries.CreateProfileParams{
		Title:      profile.Title,
		CategoryID: profile.Category.ID,
	})
	if err != nil {
		return 0, logger.Error(logger.MsgFailedToInsert, err)
	}

	return req.ID, nil
}

func (r *ProfileRepository) Read(ctx context.Context, id int64) (*model.Profile, error) {
	req, err := r.queries.ReadProfile(ctx, id)
	if err != nil {
		return nil, logger.Error(logger.MsgFailedToScan, err)
	}

	profile := &model.Profile{
		ID:    req.ID,
		Title: req.Title,
		Category: &model.Category{
			ID:    req.CategoryID,
			Title: req.CategoryTitle,
		},
		DeletedAt: validTime(req.DeletedAt),
	}

	return profile, nil
}

func (r *ProfileRepository) Update(ctx context.Context, profile *model.Profile) error {
	ct, err := r.queries.UpdateProfile(ctx, &queries.UpdateProfileParams{
		ID:         profile.ID,
		Title:      profile.Title,
		CategoryID: profile.Category.ID,
	})
	if err != nil {
		return logger.Error(logger.MsgFailedToUpdate, err)
	}

	if ct.RowsAffected() == 0 {
		return logger.Error(logger.MsgFailedToUpdate, logger.ErrNoRowsAffected)
	}

	return nil
}

func (r *ProfileRepository) Delete(ctx context.Context, id int64) error {
	ct, err := r.queries.DeleteProfile(ctx, id)
	if err != nil {
		return logger.Error(logger.MsgFailedToDelete, err)
	}

	if ct.RowsAffected() == 0 {
		return logger.Error(logger.MsgFailedToDelete, logger.ErrNoRowsAffected)
	}

	return nil
}

func (r *ProfileRepository) Restore(ctx context.Context, id int64) error {
	ct, err := r.queries.RestoreProfile(ctx, id)
	if err != nil {
		return logger.Error(logger.MsgFailedToRestore, err)
	}

	if ct.RowsAffected() == 0 {
		return logger.Error(logger.MsgFailedToRestore, logger.ErrNoRowsAffected)
	}

	return nil
}

func (r *ProfileRepository) List(ctx context.Context, qp *dto.QueryParams) ([]*model.Profile, int64, error) {
	req, err := r.queries.ListProfile(ctx, &queries.ListProfileParams{
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
		return nil, 0, nil
	}

	list := make([]*model.Profile, len(req))
	for i, item := range req {
		profile := &model.Profile{
			ID:    item.ID,
			Title: item.Title,
			Category: &model.Category{
				ID:    item.CategoryID,
				Title: item.CategoryTitle,
			},
			DeletedAt: validTime(item.DeletedAt),
		}
		list[i] = profile
	}

	return list, req[0].Total, nil
}
