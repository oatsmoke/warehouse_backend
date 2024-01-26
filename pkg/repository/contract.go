package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"warehouse_backend/pkg/model"
)

type ContractRepository struct {
	db *pgxpool.Pool
}

func NewContractRepository(db *pgxpool.Pool) *ContractRepository {
	return &ContractRepository{db: db}
}

func (r *ContractRepository) Create(number, address string) error {
	query := `
			INSERT INTO contracts (number, address) 
			VALUES ($1,$2);`
	_, err := r.db.Exec(context.Background(), query, number, address)
	if err != nil {
		return err
	}
	return nil
}

func (r *ContractRepository) GetById(id int) (model.Contract, error) {
	var contract model.Contract
	query := `
			SELECT number, address
			FROM contracts 
			WHERE id = $1 AND is_deleted = false;`
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&contract.Number,
		&contract.Address)
	if err != nil {
		return model.Contract{}, err
	}
	return contract, err
}

func (r *ContractRepository) GetAll() ([]model.Contract, error) {
	var contracts []model.Contract
	var contract model.Contract
	query := `
			SELECT contracts.id, contracts.number, contracts.address
			FROM contracts
			WHERE is_deleted = false
			ORDER BY contracts.number;`
	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(
			&contract.Id,
			&contract.Number,
			&contract.Address)
		if err != nil {
			return nil, err
		}
		contracts = append(contracts, contract)
	}
	return contracts, err
}

func (r *ContractRepository) FindByNumber(number string) (int, error) {
	var contract model.Contract
	query := `
			SELECT id 
			FROM contracts 
			WHERE number = $1;`
	err := r.db.QueryRow(context.Background(), query, number).Scan(&contract.Id)
	if err != nil {
		return 0, err
	}
	return contract.Id, nil
}

func (r *ContractRepository) Update(id int, number, address string) error {
	query := `
			UPDATE contracts 
			SET number = $2, address = $3
			WHERE id = $1;`
	_, err := r.db.Exec(context.Background(), query, id, number, address)
	if err != nil {
		return err
	}
	return nil
}

func (r *ContractRepository) Delete(id int) error {
	query := `
			UPDATE contracts 
			SET is_deleted = true
       		WHERE id = $1;`
	_, err := r.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}
	return nil
}
