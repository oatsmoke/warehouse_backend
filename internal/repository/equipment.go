package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oatsmoke/warehouse_backend/internal/dto"
	"github.com/oatsmoke/warehouse_backend/internal/lib/list_filter"
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
		return 0, err
	}

	if id == 0 {
		return 0, logger.NoRowsAffected
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

	equipment := newEquipment()
	if err := r.postgresDB.QueryRow(ctx, query, id).Scan(
		&equipment.ID,
		&equipment.SerialNumber,
		&equipment.DeletedAt,
		&equipment.Profile.ID,
		&equipment.Profile.Title,
		&equipment.Profile.Category.ID,
		&equipment.Profile.Category.Title,
	); err != nil {
		return nil, err
	}

	return equipment, nil
}

func (r *EquipmentRepository) Update(ctx context.Context, equipment *model.Equipment) error {
	const query = `
		UPDATE equipments 
		SET serial_number = $2, profile = $3
		WHERE id = $1;`

	ct, err := r.postgresDB.Exec(ctx, query, equipment.ID, equipment.SerialNumber, equipment.Profile.ID)
	if err != nil {
		return err
	}

	if ct.RowsAffected() == 0 {
		return logger.NoRowsAffected
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
		return err
	}

	if ct.RowsAffected() == 0 {
		return logger.NoRowsAffected
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
		return err
	}

	if ct.RowsAffected() == 0 {
		return logger.NoRowsAffected
	}

	return nil
}

func (r *EquipmentRepository) List(ctx context.Context, qp *dto.QueryParams) ([]*model.Equipment, error) {
	str, args := list_filter.BuildQuery(qp, []string{"e.serial_number", "p.title", "c.title"}, "e")
	query := `
		SELECT e.id, e.serial_number, e.deleted_at,
		       p.id, p.title, 
		       c.id, c.title
		FROM equipments e
		LEFT JOIN profiles p ON p.id = e.profile
		LEFT JOIN categories c ON c.id = p.category
		` + str
	fmt.Println("[Q]: ", query)
	fmt.Printf("[A]: %+v\n", args)
	rows, err := r.postgresDB.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var equipments []*model.Equipment
	for rows.Next() {
		equipment := newEquipment()
		if err := rows.Scan(
			&equipment.ID,
			&equipment.SerialNumber,
			&equipment.DeletedAt,
			&equipment.Profile.ID,
			&equipment.Profile.Title,
			&equipment.Profile.Category.ID,
			&equipment.Profile.Category.Title,
		); err != nil {
			return nil, err
		}
		equipments = append(equipments, equipment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return equipments, nil
}

func (r *EquipmentRepository) FindBySerialNumber(ctx context.Context, value string) ([]*model.Equipment, error) {
	const query = `
		SELECT equipments.id, equipments.serial_number,
		       profiles.title
		FROM equipments
		LEFT JOIN profiles ON profiles.id = equipments.profile
		WHERE LOWER(serial_number) LIKE LOWER($1);`

	rows, err := r.postgresDB.Query(ctx, query, value)
	if err != nil {
		return nil, err
	}

	var equipments []*model.Equipment
	for rows.Next() {
		equipment := newEquipment()
		if err := rows.Scan(
			&equipment.ID,
			&equipment.SerialNumber,
			&equipment.Profile.Title,
		); err != nil {
			return nil, err
		}
		equipments = append(equipments, equipment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return equipments, nil
}

func (r *EquipmentRepository) GetByIds(ctx context.Context, ids []int64) ([]*model.Equipment, error) {
	const query = `
		SELECT e.id, e.serial_number, p.deleted_at,
		       p.id, p.title
		FROM equipments e
		LEFT JOIN profiles p ON p.id = e.profile
		WHERE e.id = ANY($1);`

	rows, err := r.postgresDB.Query(ctx, query, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var equipments []*model.Equipment
	for rows.Next() {
		equipment := newEquipment()
		if err := rows.Scan(
			&equipment.ID,
			&equipment.SerialNumber,
			&equipment.DeletedAt,
			&equipment.Profile.ID,
			&equipment.Profile.Title,
		); err != nil {
			return nil, err
		}
		equipments = append(equipments, equipment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return equipments, nil
}

func newEquipment() *model.Equipment {
	return &model.Equipment{
		Profile: &model.Profile{
			Category: &model.Category{},
		},
	}
}
