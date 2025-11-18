package repository

import (
	"context"

	queries "github.com/oatsmoke/warehouse_backend/internal/db"
	"github.com/oatsmoke/warehouse_backend/internal/dto"
	"github.com/oatsmoke/warehouse_backend/internal/lib/logger"
	"github.com/oatsmoke/warehouse_backend/internal/model"
)

type EquipmentRepository struct {
	queries queries.Querier
}

func NewEquipmentRepository(queries queries.Querier) *EquipmentRepository {
	return &EquipmentRepository{
		queries: queries,
	}
}

func (r *EquipmentRepository) Create(ctx context.Context, equipment *model.Equipment) (int64, error) {
	req, err := r.queries.CreateEquipment(ctx, &queries.CreateEquipmentParams{
		SerialNumber: equipment.SerialNumber,
		ProfileID:    equipment.Profile.ID,
	})
	if err != nil {
		return 0, logger.Error(logger.MsgFailedToInsert, err)
	}

	return req.ID, nil
}

func (r *EquipmentRepository) Read(ctx context.Context, id int64) (*model.Equipment, error) {
	req, err := r.queries.ReadEquipment(ctx, id)
	if err != nil {
		return nil, logger.Error(logger.MsgFailedToScan, err)
	}

	equipment := &model.Equipment{
		ID:           req.ID,
		SerialNumber: req.SerialNumber,
		Profile: &model.Profile{
			ID:    req.ProfileID,
			Title: req.ProfileTitle,
			Category: &model.Category{
				ID:    req.CategoryID,
				Title: req.CategoryTitle,
			},
		},
		DeletedAt: validTime(req.DeletedAt),
	}

	return equipment, nil
}

func (r *EquipmentRepository) Update(ctx context.Context, equipment *model.Equipment) error {
	ct, err := r.queries.UpdateEquipment(ctx, &queries.UpdateEquipmentParams{
		ID:           equipment.ID,
		SerialNumber: equipment.SerialNumber,
		ProfileID:    equipment.Profile.ID,
	})
	if err != nil {
		return logger.Error(logger.MsgFailedToUpdate, err)
	}

	if ct.RowsAffected() == 0 {
		return logger.Error(logger.MsgFailedToUpdate, logger.ErrNoRowsAffected)
	}

	return nil
}

func (r *EquipmentRepository) Delete(ctx context.Context, id int64) error {
	ct, err := r.queries.DeleteEquipment(ctx, id)
	if err != nil {
		return logger.Error(logger.MsgFailedToDelete, err)
	}

	if ct.RowsAffected() == 0 {
		return logger.Error(logger.MsgFailedToDelete, logger.ErrNoRowsAffected)
	}

	return nil
}

func (r *EquipmentRepository) Restore(ctx context.Context, id int64) error {
	ct, err := r.queries.RestoreEquipment(ctx, id)
	if err != nil {
		return logger.Error(logger.MsgFailedToRestore, err)
	}

	if ct.RowsAffected() == 0 {
		return logger.Error(logger.MsgFailedToRestore, logger.ErrNoRowsAffected)
	}

	return nil
}

func (r *EquipmentRepository) List(ctx context.Context, qp *dto.QueryParams) ([]*model.Equipment, int64, error) {
	req, err := r.queries.ListEquipment(ctx, &queries.ListEquipmentParams{
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

	list := make([]*model.Equipment, len(req))
	for i, item := range req {
		equipment := &model.Equipment{
			ID:           item.ID,
			SerialNumber: item.SerialNumber,
			Profile: &model.Profile{
				ID:    item.ProfileID,
				Title: item.ProfileTitle,
				Category: &model.Category{
					ID:    item.CategoryID,
					Title: item.CategoryTitle,
				},
			},
			DeletedAt: validTime(item.DeletedAt),
		}
		list[i] = equipment
	}

	return list, req[0].Total, nil
}
