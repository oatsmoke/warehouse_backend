package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oatsmoke/warehouse_backend/internal/model"
)

type ReplaceRepository struct {
	DB *pgxpool.Pool
}

func NewReplaceRepository(db *pgxpool.Pool) *ReplaceRepository {
	return &ReplaceRepository{DB: db}
}

// Create is replace create
func (r *ReplaceRepository) Create(ctx context.Context, transferIds []int64) error {
	const query = `
		INSERT INTO replaces (transfer_from, transfer_to) 
		VALUES ($1, $2);`

	if _, err := r.DB.Exec(ctx, query, transferIds[0], transferIds[1]); err != nil {
		return err
	}

	return nil
}

// FindByLocationId is replaced get by location
func (r *ReplaceRepository) FindByLocationId(ctx context.Context, locationId int64) (*model.Replace, error) {
	replace := new(model.Replace)

	query := `
			SELECT id, transfer_from, transfer_to 
			FROM replaces 
			WHERE transfer_from = $1 OR transfer_to = $1;`

	if err := r.DB.QueryRow(ctx, query, locationId).Scan(
		&replace.ID,
		&replace.TransferFrom,
		&replace.TransferTo,
	); err != nil {
		return nil, err
	}

	return replace, nil
}
