package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"warehouse_backend/internal/model"
)

type ContractRepository struct {
	DB *pgxpool.Pool
}

func NewContractRepository(db *pgxpool.Pool) *ContractRepository {
	return &ContractRepository{DB: db}
}

// Create is contract create
func (r *ContractRepository) Create(ctx context.Context, number, address string) error {
	const query = `
		INSERT INTO contracts (number, address) 
		VALUES ($1, $2);`

	if _, err := r.DB.Exec(ctx, query, number, address); err != nil {
		return err
	}

	return nil
}

// Update is contract update
func (r *ContractRepository) Update(ctx context.Context, id int64, number, address string) error {
	const query = `
		UPDATE contracts 
		SET number = $2, address = $3
		WHERE id = $1;`

	if _, err := r.DB.Exec(ctx, query, id, number, address); err != nil {
		return err
	}

	return nil
}

// Delete is contract delete
func (r *ContractRepository) Delete(ctx context.Context, id int64) error {
	const query = `
		UPDATE contracts 
		SET deleted = true
       	WHERE id = $1;`

	if _, err := r.DB.Exec(ctx, query, id); err != nil {
		return err
	}

	return nil
}

// Restore is contract restore
func (r *ContractRepository) Restore(ctx context.Context, id int64) error {
	const query = `
		UPDATE contracts 
		SET deleted = false
       	WHERE id = $1;`

	if _, err := r.DB.Exec(ctx, query, id); err != nil {
		return err
	}

	return nil
}

// GetAll is to get all contracts
func (r *ContractRepository) GetAll(ctx context.Context, deleted bool) ([]*model.Contract, error) {
	var contracts []*model.Contract
	query := ""

	if deleted {
		query = `
			SELECT id, number, address
			FROM contracts
			ORDER BY number;`
	} else {
		query = `
			SELECT id, number, address
			FROM contracts
			WHERE deleted = false
			ORDER BY number;`
	}

	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		contract := new(model.Contract)
		if err := rows.Scan(
			&contract.ID,
			&contract.Number,
			&contract.Address,
		); err != nil {
			return nil, err
		}
		contracts = append(contracts, contract)
	}

	return contracts, nil
}

// GetById is to get contract by id
func (r *ContractRepository) GetById(ctx context.Context, contract *model.Contract) (*model.Contract, error) {
	const query = `
		SELECT number, address
		FROM contracts 
		WHERE id = $1;`

	if err := r.DB.QueryRow(ctx, query, contract.ID).Scan(
		&contract.Number,
		&contract.Address,
	); err != nil {
		return nil, err
	}

	return contract, nil
}

//func (r *ContractRepository) FinDByNumber(ctx context.Context, number string) (int64, error) {
//	contract := new(model.Contract)
//
//	query := `
//			SELECT id
//			FROM contracts
//			WHERE number = $1;`
//
//	if err := r.DB.QueryRow(ctx, query, number).Scan(&contract.ID); err != nil {
//		return 0, err
//	}
//
//	return contract.ID, nil
//}
