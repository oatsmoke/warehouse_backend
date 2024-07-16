package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"warehouse_backend/internal/model"
)

type DepartmentRepository struct {
	db *pgxpool.Pool
}

func NewDepartmentRepository(db *pgxpool.Pool) *DepartmentRepository {
	return &DepartmentRepository{db: db}
}

func (r *DepartmentRepository) Create(ctx context.Context, title string) error {
	query := `
			INSERT INTO departments (title) 
			VALUES ($1);`

	if _, err := r.db.Exec(ctx, query, title); err != nil {
		return err
	}

	return nil
}

func (r *DepartmentRepository) GetById(ctx context.Context, id int64) (*model.Department, error) {
	department := new(model.Department)

	query := `
			SELECT id, title
			FROM departments 
			WHERE id = $1;`

	if err := r.db.QueryRow(ctx, query, id).Scan(&department.ID, &department.Title); err != nil {
		return nil, err
	}
	return department, nil
}

func (r *DepartmentRepository) GetAll(ctx context.Context) ([]*model.Department, error) {
	var departments []*model.Department
	department := new(model.Department)

	query := `
			SELECT id, title
			FROM departments
			WHERE is_deleted = false
			ORDER BY title;`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		if err := rows.Scan(&department.ID, &department.Title); err != nil {
			return nil, err
		}
		departments = append(departments, department)
	}

	return departments, err
}

func (r *DepartmentRepository) GetAllButOne(ctx context.Context, id, employeeId int64) ([]*model.Department, error) {
	var departments []*model.Department
	department := new(model.Department)

	query := `
			SELECT departments.id, departments.title
			FROM departments
         	LEFT JOIN employees on departments.id = employees.department
			WHERE departments.id != $1
  			AND departments.is_deleted = false
  			AND employees.id = $2
			AND departments.id = employees.department
			ORDER BY title;`

	rows, err := r.db.Query(ctx, query, id, employeeId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		if err := rows.Scan(&department.ID, &department.Title); err != nil {
			return nil, err
		}
		departments = append(departments, department)
	}

	return departments, nil
}

func (r *DepartmentRepository) GetAllButOneForAdmin(ctx context.Context, id int64) ([]*model.Department, error) {
	var departments []*model.Department
	department := new(model.Department)

	query := `
			SELECT departments.id, departments.title
			FROM departments
			WHERE departments.id != $1
  			AND departments.is_deleted = false
			ORDER BY title;`

	rows, err := r.db.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		if err := rows.Scan(&department.ID, &department.Title); err != nil {
			return nil, err
		}
		departments = append(departments, department)
	}

	return departments, nil
}

func (r *DepartmentRepository) FindByTitle(ctx context.Context, title string) (int64, error) {
	department := new(model.Department)

	query := `
			SELECT id 
			FROM departments 
			WHERE title = $1;`

	if err := r.db.QueryRow(ctx, query, title).Scan(&department.ID); err != nil {
		return 0, err
	}

	return department.ID, nil
}

func (r *DepartmentRepository) Update(ctx context.Context, id int64, title string) error {
	query := `
			UPDATE departments 
			SET title = $2
			WHERE id = $1;`

	if _, err := r.db.Exec(ctx, query, id, title); err != nil {
		return err
	}

	return nil
}

func (r *DepartmentRepository) Delete(ctx context.Context, id int64) error {
	query := `
			UPDATE departments 
			SET is_deleted = true
       		WHERE id = $1;`

	if _, err := r.db.Exec(ctx, query, id); err != nil {
		return err
	}

	return nil
}
