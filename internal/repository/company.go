package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oatsmoke/warehouse_backend/internal/model"
)

type CompanyRepository struct {
	DB *pgxpool.Pool
}

func NewCompanyRepository(db *pgxpool.Pool) *CompanyRepository {
	return &CompanyRepository{DB: db}
}

// Create is company create
func (r *CompanyRepository) Create(ctx context.Context, title string) error {
	const query = `
		INSERT INTO companies (title) 
		VALUES ($1);`

	if _, err := r.DB.Exec(ctx, query, title); err != nil {
		return err
	}

	return nil
}

// Update is company update
func (r *CompanyRepository) Update(ctx context.Context, id int64, title string) error {
	query := `
			UPDATE companies 
			SET title = $2
			WHERE id = $1;`

	if _, err := r.DB.Exec(ctx, query, id, title); err != nil {
		return err
	}

	return nil
}

// Delete is company delete
func (r *CompanyRepository) Delete(ctx context.Context, id int64) error {
	const query = `
		UPDATE companies 
		SET deleted = true
       	WHERE id = $1;`

	if _, err := r.DB.Exec(ctx, query, id); err != nil {
		return err
	}

	return nil
}

// Restore is company restore
func (r *CompanyRepository) Restore(ctx context.Context, id int64) error {
	const query = `
		UPDATE companies 
		SET deleted = false
       	WHERE id = $1;`

	if _, err := r.DB.Exec(ctx, query, id); err != nil {
		return err
	}

	return nil
}

// GetAll is to get all companies
func (r *CompanyRepository) GetAll(ctx context.Context, deleted bool) ([]*model.Company, error) {
	var companies []*model.Company
	query := ""

	if deleted {
		query = `
			SELECT id, title, deleted
			FROM companies
			ORDER BY title;`
	} else {
		query = `
			SELECT id, title, deleted
			FROM companies
			WHERE deleted = false
			ORDER BY title;`
	}

	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		company := new(model.Company)
		if err := rows.Scan(
			&company.ID,
			&company.Title,
			&company.Deleted,
		); err != nil {
			return nil, err
		}
		companies = append(companies, company)
	}

	return companies, nil
}

// GetById is to get company by id
func (r *CompanyRepository) GetById(ctx context.Context, company *model.Company) (*model.Company, error) {
	const query = `
		SELECT title, deleted
		FROM companies 
		WHERE id = $1;`

	if err := r.DB.QueryRow(ctx, query, company.ID).Scan(
		&company.Title,
		&company.Deleted,
	); err != nil {
		return nil, err
	}

	return company, nil
}

//func (r *CompanyRepository) FinDByTitle(ctx context.Context, title string) (int64, error) {
//	company := new(model.Company)
//
//	query := `
//			SELECT id
//			FROM companies
//			WHERE title = $1;`
//
//	if err := r.DB.QueryRow(ctx, query, title).Scan(&company.ID); err != nil {
//		return 0, err
//	}
//
//	return company.ID, nil
//}
