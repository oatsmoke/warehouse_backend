package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oatsmoke/warehouse_backend/internal/lib/logger"
	"github.com/oatsmoke/warehouse_backend/internal/model"
)

type ContractRepository struct {
	postgresDB *pgxpool.Pool
}

func NewContractRepository(postgresDB *pgxpool.Pool) *ContractRepository {
	return &ContractRepository{
		postgresDB: postgresDB,
	}
}

func (r *ContractRepository) Create(ctx context.Context, contract *model.Contract) (int64, error) {
	const query = `
		INSERT INTO contracts (number, address) 
		VALUES ($1, $2)
		RETURNING id;`

	var id int64
	if err := r.postgresDB.QueryRow(ctx, query, contract.Number, contract.Address).Scan(&id); err != nil {
		return 0, err
	}

	if id == 0 {
		return 0, logger.NoRowsAffected
	}

	return id, nil
}

func (r *ContractRepository) Read(ctx context.Context, id int64) (*model.Contract, error) {
	const query = `
		SELECT id, number, address, deleted_at
		FROM contracts 
		WHERE id = $1;`

	contract := new(model.Contract)
	if err := r.postgresDB.QueryRow(ctx, query, id).Scan(
		&contract.ID,
		&contract.Number,
		&contract.Address,
		&contract.DeletedAt,
	); err != nil {
		return nil, err
	}

	return contract, nil
}

func (r *ContractRepository) Update(ctx context.Context, contract *model.Contract) error {
	const query = `
		UPDATE contracts 
		SET number = $2, address = $3
		WHERE id = $1;`

	ct, err := r.postgresDB.Exec(ctx, query, contract.ID, contract.Number, contract.Address)
	if err != nil {
		return err
	}

	if ct.RowsAffected() == 0 {
		return logger.NoRowsAffected
	}

	return nil
}

func (r *ContractRepository) Delete(ctx context.Context, id int64) error {
	const query = `
		UPDATE contracts 
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

func (r *ContractRepository) Restore(ctx context.Context, id int64) error {
	const query = `
		UPDATE contracts 
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

func (r *ContractRepository) List(ctx context.Context, withDeleted bool) ([]*model.Contract, error) {
	const query = `
		SELECT id, number, address, deleted_at
		FROM contracts
		WHERE $1 OR deleted_at IS NULL
		ORDER BY number;`

	rows, err := r.postgresDB.Query(ctx, query, withDeleted)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contracts []*model.Contract
	for rows.Next() {
		contract := new(model.Contract)
		if err := rows.Scan(
			&contract.ID,
			&contract.Number,
			&contract.Address,
			&contract.DeletedAt,
		); err != nil {
			return nil, err
		}
		contracts = append(contracts, contract)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return contracts, nil
}
