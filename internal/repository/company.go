package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"warehouse_backend/internal/model"
)

type CompanyRepository struct {
	db *pgxpool.Pool
}

func NewCompanyRepository(db *pgxpool.Pool) *CompanyRepository {
	return &CompanyRepository{db: db}
}

func (r *CompanyRepository) Create(ctx context.Context, title string) error {
	query := `
			INSERT INTO companies (title) 
			VALUES ($1);`

	if _, err := r.db.Exec(ctx, query, title); err != nil {
		return err
	}

	return nil
}

func (r *CompanyRepository) GetById(ctx context.Context, id int64) (*model.Company, error) {
	company := new(model.Company)

	query := `
			SELECT title
			FROM companies 
			WHERE id = $1;`

	if err := r.db.QueryRow(ctx, query, id).Scan(&company.Title); err != nil {
		return nil, err
	}

	return company, nil
}

func (r *CompanyRepository) GetAll(ctx context.Context) ([]*model.Company, error) {
	var companies []*model.Company
	company := new(model.Company)

	query := `
			SELECT id, title
			FROM companies
			ORDER BY title;`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		if err := rows.Scan(&company.ID, &company.Title); err != nil {
			return nil, err
		}
		companies = append(companies, company)
	}

	return companies, nil
}

func (r *CompanyRepository) FindByTitle(ctx context.Context, title string) (int64, error) {
	company := new(model.Company)

	query := `
			SELECT id 
			FROM companies 
			WHERE title = $1;`

	if err := r.db.QueryRow(ctx, query, title).Scan(&company.ID); err != nil {
		return 0, err
	}

	return company.ID, nil
}

func (r *CompanyRepository) Update(ctx context.Context, id int64, title string) error {
	query := `
			UPDATE companies 
			SET title = $2
			WHERE id = $1;`

	if _, err := r.db.Exec(ctx, query, id, title); err != nil {
		return err
	}

	return nil
}

func (r *CompanyRepository) Delete(ctx context.Context, id int64) error {
	query := `
			DELETE FROM companies 
       		WHERE id = $1;`

	if _, err := r.db.Exec(ctx, query, id); err != nil {
		return err
	}

	return nil
}
