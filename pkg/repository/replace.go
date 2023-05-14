package repository

import (
	"context"
	"github.com/jackc/pgx/pgxpool"
	"warehouse_backend/pkg/model"
)

type ReplaceRepository struct {
	db *pgxpool.Pool
}

func NewReplaceRepository(db *pgxpool.Pool) *ReplaceRepository {
	return &ReplaceRepository{db: db}
}

func (r *ReplaceRepository) Create(ids []int) error {
	query := `
			INSERT INTO replaces (transfer_from, transfer_to) 
			VALUES ($1, $2);`
	_, err := r.db.Exec(context.Background(), query, ids[0], ids[1])
	if err != nil {
		return err
	}
	return nil
}
func (r *ReplaceRepository) FindByLocationId(id int) (model.Replace, error) {
	var replace model.Replace
	query := `
			SELECT id, transfer_from,transfer_to 
			FROM replaces 
			WHERE transfer_from = $1 OR transfer_to = $1;`
	err := r.db.QueryRow(context.Background(), query, id).Scan(&replace.Id, &replace.TransferFrom, &replace.TransferTo)
	if err != nil {
		return model.Replace{}, err
	}
	return replace, nil
}
