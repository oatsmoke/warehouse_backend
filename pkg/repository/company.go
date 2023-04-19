package repository

import (
	"context"
	"github.com/jackc/pgx/pgxpool"
	"warehouse_backend/pkg/model"
)

type CompanyRepository struct {
	db *pgxpool.Pool
}

func NewCompanyRepository(db *pgxpool.Pool) *CompanyRepository {
	return &CompanyRepository{db: db}
}

func (r *CompanyRepository) Create(title string) error {
	query := `
			INSERT INTO companies (title) 
			VALUES ($1);`
	_, err := r.db.Exec(context.Background(), query, title)
	if err != nil {
		return err
	}
	return nil
}

func (r *CompanyRepository) GetById(id int) (model.Company, error) {
	var company model.Company
	query := `
			SELECT title
			FROM companies 
			WHERE id = $1;`
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&company.Title)
	if err != nil {
		return model.Company{}, err
	}
	return company, err
}

func (r *CompanyRepository) GetAll() ([]model.Company, error) {
	var companies []model.Company
	var company model.Company
	query := `
			SELECT id, title
			FROM companies
			ORDER BY title;`
	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(
			&company.Id,
			&company.Title)
		if err != nil {
			return nil, err
		}
		companies = append(companies, company)
	}
	return companies, err
}

func (r *CompanyRepository) FindByTitle(title string) (int, error) {
	var company model.Company
	query := `
			SELECT id 
			FROM companies 
			WHERE title = $1;`
	err := r.db.QueryRow(context.Background(), query, title).Scan(&company.Id)
	if err != nil {
		return 0, err
	}
	return company.Id, nil
}

func (r *CompanyRepository) Update(id int, title string) error {
	query := `
			UPDATE companies 
			SET title = $2
			WHERE id = $1;`
	_, err := r.db.Exec(context.Background(), query, id, title)
	if err != nil {
		return err
	}
	return nil
}

func (r *CompanyRepository) Delete(id int) error {
	query := `
			DELETE FROM companies 
       		WHERE id = $1;`
	_, err := r.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}
	return nil
}
