package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"warehouse_backend/internal/model"
)

type ContractRepository struct {
	db *pgxpool.Pool
}

func NewContractRepository(db *pgxpool.Pool) *ContractRepository {
	return &ContractRepository{db: db}
}

func (r *ContractRepository) Create(ctx context.Context, number, address string) error {
	query := `
			INSERT INTO contracts (number, address) 
			VALUES ($1,$2);`

	if _, err := r.db.Exec(ctx, query, number, address); err != nil {
		return err
	}

	return nil
}

func (r *ContractRepository) GetById(ctx context.Context, id int64) (*model.Contract, error) {
	contract := new(model.Contract)

	query := `
			SELECT number, address
			FROM contracts 
			WHERE id = $1 AND is_deleted = false;`

	if err := r.db.QueryRow(ctx, query, id).Scan(&contract.Number, &contract.Address); err != nil {
		return nil, err
	}

	return contract, nil
}

func (r *ContractRepository) GetAll(ctx context.Context) ([]*model.Contract, error) {
	var contracts []*model.Contract
	contract := new(model.Contract)

	query := `
			SELECT id, number, address
			FROM contracts
			WHERE is_deleted = false
			ORDER BY number;`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		if err := rows.Scan(&contract.ID, &contract.Number, &contract.Address); err != nil {
			return nil, err
		}
		contracts = append(contracts, contract)
	}

	return contracts, nil
}

func (r *ContractRepository) FindByNumber(ctx context.Context, number string) (int64, error) {
	contract := new(model.Contract)

	query := `
			SELECT id 
			FROM contracts 
			WHERE number = $1;`

	if err := r.db.QueryRow(ctx, query, number).Scan(&contract.ID); err != nil {
		return 0, err
	}

	return contract.ID, nil
}

func (r *ContractRepository) Update(ctx context.Context, id int64, number, address string) error {
	query := `
			UPDATE contracts 
			SET number = $2, address = $3
			WHERE id = $1;`

	if _, err := r.db.Exec(ctx, query, id, number, address); err != nil {
		return err
	}

	return nil
}

func (r *ContractRepository) Delete(ctx context.Context, id int64) error {
	query := `
			UPDATE contracts 
			SET is_deleted = true
       		WHERE id = $1;`

	if _, err := r.db.Exec(ctx, query, id); err != nil {
		return err
	}

	return nil
}
