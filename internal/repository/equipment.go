package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
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

func (r *EquipmentRepository) Create(ctx context.Context, equipment *model.Equipment) (int64, error) {
	const query = `
		INSERT INTO equipments (serial_number, profile) 
		VALUES ($1, $2)
		RETURNING id;`

	var id int64
	if err := r.postgresDB.QueryRow(ctx, query, equipment.SerialNumber, equipment.Profile.ID).Scan(&id); err != nil {
		return 0, logger.Error(logger.MsgFailedToInsert, err)
	}

	if id == 0 {
		return 0, logger.Error(logger.MsgFailedToInsert, logger.ErrNoRowsAffected)
	}

	return id, nil
}

func (r *EquipmentRepository) Read(ctx context.Context, id int64) (*model.Equipment, error) {
	const query = `
		SELECT e.id, e.serial_number, p.deleted_at,
		       p.id, p.title,
			   c.id, c.title
		FROM equipments e
		LEFT JOIN profiles p ON p.id = e.profile
		LEFT JOIN categories c ON c.id = p.category
		WHERE e.id = $1;`

	equipment := model.NewEquipment()
	if err := r.postgresDB.QueryRow(ctx, query, id).Scan(
		&equipment.ID,
		&equipment.SerialNumber,
		&equipment.DeletedAt,
		&equipment.Profile.ID,
		&equipment.Profile.Title,
		&equipment.Profile.Category.ID,
		&equipment.Profile.Category.Title,
	); err != nil {
		return nil, logger.Error(logger.MsgFailedToScan, err)
	}

	return equipment, nil
}

func (r *EquipmentRepository) Update(ctx context.Context, equipment *model.Equipment) error {
	const query = `
		UPDATE equipments 
		SET serial_number = $2, profile = $3
		WHERE id = $1 AND (serial_number != $2 OR profile != $3);`

	ct, err := r.postgresDB.Exec(ctx, query, equipment.ID, equipment.SerialNumber, equipment.Profile.ID)
	if err != nil {
		return logger.Error(logger.MsgFailedToUpdate, err)
	}

	if ct.RowsAffected() == 0 {
		return logger.Error(logger.MsgFailedToUpdate, logger.ErrNoRowsAffected)
	}

	return nil
}

func (r *EquipmentRepository) Delete(ctx context.Context, id int64) error {
	const query = `
		UPDATE equipments 
		SET deleted_at = now()
       	WHERE id = $1 AND deleted_at IS NULL;`

	ct, err := r.postgresDB.Exec(ctx, query, id)
	if err != nil {
		return logger.Error(logger.MsgFailedToDelete, err)
	}

	if ct.RowsAffected() == 0 {
		return logger.Error(logger.MsgFailedToDelete, logger.ErrNoRowsAffected)
	}

	return nil
}

func (r *EquipmentRepository) Restore(ctx context.Context, id int64) error {
	const query = `
		UPDATE equipments 
		SET deleted_at = NULL 
       	WHERE id = $1 AND deleted_at IS NOT NULL;`

	ct, err := r.postgresDB.Exec(ctx, query, id)
	if err != nil {
		return logger.Error(logger.MsgFailedToRestore, err)
	}

	if ct.RowsAffected() == 0 {
		return logger.Error(logger.MsgFailedToRestore, logger.ErrNoRowsAffected)
	}

	return nil
}

func (r *EquipmentRepository) List(ctx context.Context, qp *dto.QueryParams) ([]*model.Equipment, int, error) {
	const query = `
		SELECT e.id, e.serial_number, e.deleted_at,
		       p.id, p.title,
		       c.id, c.title, COUNT(*) OVER() AS total
		FROM equipments e
		LEFT JOIN profiles p ON p.id = e.profile
		LEFT JOIN categories c ON c.id = p.category
		WHERE ($1 = true OR e.deleted_at IS NULL)
		AND ($2 = '' OR (e.serial_number || ' ' || p.title || ' ' || c.title) ILIKE '%' || $2 || '%')
		AND (array_length($3::bigint[], 1) IS NULL OR e.id = ANY ($3))
		ORDER BY CASE WHEN $4 = 'id' AND $5 = 'asc' THEN e.id::text END,
				 CASE WHEN $4 = 'id' AND $5 = 'desc' THEN e.id::text END DESC,
				 CASE WHEN $4 = 'serial_number' AND $5 = 'asc' THEN e.serial_number END,
				 CASE WHEN $4 = 'serial_number' AND $5 = 'desc' THEN e.serial_number END DESC,
				 CASE WHEN $4 = 'p_title' AND $5 = 'asc' THEN p.title END,
				 CASE WHEN $4 = 'p_title' AND $5 = 'desc' THEN p.title END DESC,
				 CASE WHEN $4 = 'c_title' AND $5 = 'asc' THEN c.title END,
				 CASE WHEN $4 = 'c_title' AND $5 = 'desc' THEN c.title END DESC
		LIMIT $6 OFFSET $7;`

	rows, err := r.postgresDB.Query(
		ctx,
		query,
		qp.WithDeleted,
		qp.Search,
		qp.Ids,
		qp.SortBy,
		qp.Order,
		qp.Limit,
		qp.Offset,
	)
	if err != nil {
		return nil, 0, logger.Error(logger.MsgFailedToSelect, err)
	}
	defer rows.Close()

	var equipments []*model.Equipment
	var total int
	for rows.Next() {
		equipment := model.NewEquipment()
		if err := rows.Scan(
			&equipment.ID,
			&equipment.SerialNumber,
			&equipment.DeletedAt,
			&equipment.Profile.ID,
			&equipment.Profile.Title,
			&equipment.Profile.Category.ID,
			&equipment.Profile.Category.Title,
			&total,
		); err != nil {
			return nil, 0, logger.Error(logger.MsgFailedToScan, err)
		}
		equipments = append(equipments, equipment)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, logger.Error(logger.MsgFailedToIterateOverRows, err)
	}

	return equipments, total, nil
}
