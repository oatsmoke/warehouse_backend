package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oatsmoke/warehouse_backend/internal/dto"
	"github.com/oatsmoke/warehouse_backend/internal/lib/list_filter"
	"github.com/oatsmoke/warehouse_backend/internal/lib/logger"
	"github.com/oatsmoke/warehouse_backend/internal/model"
)

type CompanyRepository struct {
	postgresDB *pgxpool.Pool
}

func NewCompanyRepository(postgresDB *pgxpool.Pool) *CompanyRepository {
	return &CompanyRepository{
		postgresDB: postgresDB,
	}
}

func (r *CompanyRepository) Create(ctx context.Context, company *model.Company) (int64, error) {
	const query = `
		INSERT INTO companies (title) 
		VALUES ($1)
		RETURNING id;`

	var id int64
	if err := r.postgresDB.QueryRow(ctx, query, company.Title).Scan(&id); err != nil {
		return 0, err
	}

	if id == 0 {
		return 0, logger.NoRowsAffected
	}

	return id, nil
}

func (r *CompanyRepository) Read(ctx context.Context, id int64) (*model.Company, error) {
	const query = `
		SELECT id, title, deleted_at
		FROM companies 
		WHERE id = $1;`

	company := new(model.Company)
	if err := r.postgresDB.QueryRow(ctx, query, id).Scan(
		&company.ID,
		&company.Title,
		&company.DeletedAt,
	); err != nil {
		return nil, err
	}

	return company, nil
}

func (r *CompanyRepository) Update(ctx context.Context, company *model.Company) error {
	const query = `
		UPDATE companies 
		SET title = $2
		WHERE id = $1;`

	ct, err := r.postgresDB.Exec(ctx, query, company.ID, company.Title)
	if err != nil {
		return err
	}

	if ct.RowsAffected() == 0 {
		return logger.NoRowsAffected
	}

	return nil
}

func (r *CompanyRepository) Delete(ctx context.Context, id int64) error {
	const query = `
		UPDATE companies 
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

func (r *CompanyRepository) Restore(ctx context.Context, id int64) error {
	const query = `
		UPDATE companies 
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

func (r *CompanyRepository) List(ctx context.Context, qp *dto.QueryParams) ([]*model.Company, error) {
	fields := []string{"title"}
	str, args := list_filter.BuildQuery(qp, fields, "c")

	query := `
		SELECT id, title, deleted_at
		FROM companies c
		` + str

	rows, err := r.postgresDB.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var companies []*model.Company
	for rows.Next() {
		company := new(model.Company)
		if err := rows.Scan(
			&company.ID,
			&company.Title,
			&company.DeletedAt,
		); err != nil {
			return nil, err
		}
		companies = append(companies, company)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return companies, nil
}
