package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"warehouse_backend/internal/model"
)

type EquipmentRepository struct {
	DB *pgxpool.Pool
}

func NewEquipmentRepository(db *pgxpool.Pool) *EquipmentRepository {
	return &EquipmentRepository{DB: db}
}

// Create is equipment create
func (r *EquipmentRepository) Create(ctx context.Context, serialNumber string, profileId int64) (int64, error) {
	const query = `
		INSERT INTO equipments (serial_number, profile) 
		VALUES ($1, $2)
		RETURNING id;`

	var equipmentId int64
	if err := r.DB.QueryRow(ctx, query, serialNumber, profileId).Scan(&equipmentId); err != nil {
		return 0, err
	}

	return equipmentId, nil
}

// Update is equipment update
func (r *EquipmentRepository) Update(ctx context.Context, id int64, serialNumber string, profileId int64) error {
	const query = `
		UPDATE equipments 
		SET serial_number = $2, profile = $3
		WHERE id = $1;`

	if _, err := r.DB.Exec(ctx, query, id, serialNumber, profileId); err != nil {
		return err
	}

	return nil
}

// Delete is equipment delete
func (r *EquipmentRepository) Delete(ctx context.Context, id int64) error {
	query := `
			UPDATE equipments 
			SET deleted = true
       		WHERE id = $1;`

	if _, err := r.DB.Exec(ctx, query, id); err != nil {
		return err
	}

	return nil
}

// Restore is equipment restore
func (r *EquipmentRepository) Restore(ctx context.Context, id int64) error {
	query := `
			UPDATE equipments 
			SET deleted = false
       		WHERE id = $1;`

	if _, err := r.DB.Exec(ctx, query, id); err != nil {
		return err
	}

	return nil
}

// GetAll is equipment get all
func (r *EquipmentRepository) GetAll(ctx context.Context) ([]*model.Equipment, error) {
	var equipments []*model.Equipment

	const query = `
		SELECT equipments.id, equipments.serial_number, 
		       profiles.id, profiles.title, 
		       categories.id, categories.title
		FROM equipments
		LEFT JOIN profiles ON profiles.id = equipments.profile
		LEFT JOIN categories ON categories.id = profiles.category
		WHERE equipments.deleted = FALSE
		ORDER BY profiles.title;`

	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		equipment := new(model.Equipment)
		if err := rows.Scan(
			&equipment.ID,
			&equipment.SerialNumber,
			&equipment.Profile.ID,
			&equipment.Profile.Title,
			&equipment.Profile.Category.ID,
			&equipment.Profile.Category.Title,
		); err != nil {
			return nil, err
		}
		equipments = append(equipments, equipment)
	}

	return equipments, nil
}

// GetByIds is equipment get by ids
func (r *EquipmentRepository) GetByIds(ctx context.Context, ids []int64) ([]*model.Equipment, error) {
	var equipments []*model.Equipment

	const query = `
		SELECT equipments.id, equipments.serial_number,
		       profiles.id, profiles.title
		FROM equipments
		LEFT JOIN profiles ON profiles.id = equipments.profile
		WHERE equipments.id = ANY($1);`

	rows, err := r.DB.Query(ctx, query, ids)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		equipment := new(model.Equipment)
		equipment.Profile = new(model.Profile)
		if err := rows.Scan(
			&equipment.ID,
			&equipment.SerialNumber,
			&equipment.Profile.ID,
			&equipment.Profile.Title,
		); err != nil {
			return nil, err
		}
		equipments = append(equipments, equipment)
	}

	return equipments, nil
}

//// GetByProfile is equipment get by profile
//func (r *EquipmentRepository) GetByProfile(ctx context.Context, id int64) ([]*model.Equipment, error) {
//	var equipments []*model.Equipment
//
//	const query = `
//		SELECT equipments.id, equipments.serial_number
//		FROM equipments
//		LEFT JOIN profiles ON profiles.id = equipments.profile
//		WHERE profiles.id = $1;`
//
//	rows, err := r.DB.Query(ctx, query, id)
//	if err != nil {
//		return nil, err
//	}
//
//	for rows.Next() {
//		equipment := new(model.Equipment)
//		if err := rows.Scan(
//			&equipment.ID,
//			&equipment.SerialNumber,
//		); err != nil {
//			return nil, err
//		}
//		equipments = append(equipments, equipment)
//	}
//
//	return equipments, nil
//}
//
//func (r *EquipmentRepository) GetBySerialNumber(ctx context.Context, equipment *model.Equipment) (*model.Equipment, error) {
//	const query = `
//		SELECT id
//		FROM equipments
//		WHERE serial_number = $1;`
//
//	if err := r.DB.QueryRow(ctx, query, equipment.SerialNumber).Scan(&equipment.ID); err != nil {
//		return nil, err
//	}
//
//	return equipment, nil
//}
