package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	queries "github.com/oatsmoke/warehouse_backend/internal/db"
	"github.com/oatsmoke/warehouse_backend/internal/dto"
	"github.com/oatsmoke/warehouse_backend/internal/lib/logger"
	"github.com/oatsmoke/warehouse_backend/internal/model"
)

type EquipmentRepository struct {
	postgresDB *pgxpool.Pool
}

func NewEquipmentRepository(postgresDB *pgxpool.Pool) *EquipmentRepository {
	return &EquipmentRepository{
		postgresDB: postgresDB,
	}
}

func (r *EquipmentRepository) Create(ctx context.Context, equipment *queries.CreateEquipmentParams, location *queries.AddToStorageParams) (int64, error) {
	tx, err := r.postgresDB.Begin(ctx)
	if err != nil {
		return 0, logger.Error("", err)
	}
	defer tx.Rollback(ctx)

	q := queries.New(tx)

	e, err := q.CreateEquipment(ctx, equipment)
	if err != nil {
		return 0, logger.Error(logger.MsgFailedToInsert, err)
	}

	location.EquipmentID = e.ID

	ct, err := q.AddToStorage(ctx, location)
	if err != nil {
		return 0, logger.Error(logger.MsgFailedToInsert, err)
	}

	if ct.RowsAffected() == 0 {
		return 0, logger.Error(logger.MsgFailedToInsert, logger.ErrNoRowsAffected)
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, logger.Error("", err)
	}

	return e.ID, nil
}

func (r *EquipmentRepository) Read(ctx context.Context, id int64) (*model.Equipment, error) {
	res, err := queries.New(r.postgresDB).ReadEquipment(ctx, id)
	if err != nil {
		return nil, logger.Error(logger.MsgFailedToScan, err)
	}

	equipment := &model.Equipment{
		ID: res.ID,
		Company: &model.Company{
			ID:    res.CompanyID,
			Title: res.CompanyTitle,
		},
		Profile: &model.Profile{
			ID:    res.ProfileID,
			Title: res.ProfileTitle,
			Category: &model.Category{
				ID:    res.CategoryID,
				Title: res.CategoryTitle,
			},
		},
		SerialNumber: res.SerialNumber,
		DeletedAt:    validTime(res.DeletedAt),
	}

	return equipment, nil
}

func (r *EquipmentRepository) Update(ctx context.Context, equipment *model.Equipment) error {
	ct, err := queries.New(r.postgresDB).UpdateEquipment(ctx, &queries.UpdateEquipmentParams{
		ID:           equipment.ID,
		CompanyID:    equipment.Company.ID,
		ProfileID:    equipment.Profile.ID,
		SerialNumber: equipment.SerialNumber,
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
	ct, err := queries.New(r.postgresDB).DeleteEquipment(ctx, id)
	if err != nil {
		return logger.Error(logger.MsgFailedToDelete, err)
	}

	if ct.RowsAffected() == 0 {
		return logger.Error(logger.MsgFailedToDelete, logger.ErrNoRowsAffected)
	}

	return nil
}

func (r *EquipmentRepository) Restore(ctx context.Context, id int64) error {
	ct, err := queries.New(r.postgresDB).RestoreEquipment(ctx, id)
	if err != nil {
		return logger.Error(logger.MsgFailedToRestore, err)
	}

	if ct.RowsAffected() == 0 {
		return logger.Error(logger.MsgFailedToRestore, logger.ErrNoRowsAffected)
	}

	return nil
}

func (r *EquipmentRepository) List(ctx context.Context, qp *dto.QueryParams) ([]*model.Equipment, int64, error) {
	req, err := queries.New(r.postgresDB).ListEquipment(ctx, &queries.ListEquipmentParams{
		WithDeleted:      qp.WithDeleted,
		Search:           qp.Search,
		Ids:              qp.IDs,
		SortColumn:       qp.SortColumn,
		SortOrder:        qp.SortOrder,
		PaginationLimit:  qp.PaginationLimit,
		PaginationOffset: qp.PaginationOffset,
	})
	if err != nil {
		return nil, 0, logger.Error(logger.MsgFailedToSelect, err)
	}

	if len(req) < 1 {
		return []*model.Equipment{}, 0, nil
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
