package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"warehouse_backend/internal/model"
)

type DepartmentRepository struct {
	DB *pgxpool.Pool
}

func NewDepartmentRepository(db *pgxpool.Pool) *DepartmentRepository {
	return &DepartmentRepository{DB: db}
}

// Create is department create
func (r *DepartmentRepository) Create(ctx context.Context, title string) error {
	const query = `
		INSERT INTO departments (title) 
		VALUES ($1);`

	if _, err := r.DB.Exec(ctx, query, title); err != nil {
		return err
	}

	return nil
}

// Update is department update
func (r *DepartmentRepository) Update(ctx context.Context, id int64, title string) error {
	const query = `
		UPDATE departments 
		SET title = $2
		WHERE id = $1;`

	if _, err := r.DB.Exec(ctx, query, id, title); err != nil {
		return err
	}

	return nil
}

// Delete is department delete
func (r *DepartmentRepository) Delete(ctx context.Context, id int64) error {
	const query = `
		UPDATE departments 
		SET deleted = true
       	WHERE id = $1;`

	if _, err := r.DB.Exec(ctx, query, id); err != nil {
		return err
	}

	return nil
}

// Restore is department restore
func (r *DepartmentRepository) Restore(ctx context.Context, id int64) error {
	const query = `
		UPDATE departments 
		SET deleted = false
       	WHERE id = $1;`

	if _, err := r.DB.Exec(ctx, query, id); err != nil {
		return err
	}

	return nil
}

// GetAll is department get all
func (r *DepartmentRepository) GetAll(ctx context.Context, deleted bool) ([]*model.Department, error) {
	var departments []*model.Department
	query := ""

	if deleted {
		query = `
			SELECT id, title, deleted
			FROM departments
			ORDER BY title;`
	} else {
		query = `
			SELECT id, title, deleted
			FROM departments
			WHERE deleted = false
			ORDER BY title;`
	}

	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		department := new(model.Department)
		if err := rows.Scan(
			&department.ID,
			&department.Title,
			&department.Deleted,
		); err != nil {
			return nil, err
		}
		departments = append(departments, department)
	}

	return departments, err
}

// GetById is department get by id
func (r *DepartmentRepository) GetById(ctx context.Context, department *model.Department) (*model.Department, error) {
	const query = `
		SELECT title, deleted
		FROM departments 
		WHERE id = $1;`

	if err := r.DB.QueryRow(ctx, query, department.ID).Scan(
		&department.Title,
		&department.Deleted,
	); err != nil {
		return nil, err
	}
	return department, nil
}

// GetAllButOne is department get all but one
func (r *DepartmentRepository) GetAllButOne(ctx context.Context, id, employeeId int64) ([]*model.Department, error) {
	var departments []*model.Department

	const query = `
		SELECT departments.id, departments.title
		FROM departments
        LEFT JOIN employees on departments.id = employees.department
		WHERE departments.id != $1
  		AND departments.deleted = false
  		AND employees.id = $2
		AND departments.id = employees.department
		ORDER BY title;`

	rows, err := r.DB.Query(ctx, query, id, employeeId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		department := new(model.Department)
		if err := rows.Scan(
			&department.ID,
			&department.Title,
		); err != nil {
			return nil, err
		}
		departments = append(departments, department)
	}

	return departments, nil
}

// GetAllButOneForAdmin is department get all but one for admin
func (r *DepartmentRepository) GetAllButOneForAdmin(ctx context.Context, id int64) ([]*model.Department, error) {
	var departments []*model.Department

	const query = `
		SELECT departments.id, departments.title
		FROM departments
		WHERE departments.id != $1
  		AND departments.deleted = false
		ORDER BY title;`

	rows, err := r.DB.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		department := new(model.Department)
		if err := rows.Scan(
			&department.ID,
			&department.Title,
		); err != nil {
			return nil, err
		}
		departments = append(departments, department)
	}

	return departments, nil
}

//func (r *DepartmentRepository) FinDByTitle(ctx context.Context, title string) (int64, error) {
//	department := new(model.Department)
//
//	query := `
//			SELECT id
//			FROM departments
//			WHERE title = $1;`
//
//	if err := r.DB.QueryRow(ctx, query, title).Scan(&department.ID); err != nil {
//		return 0, err
//	}
//
//	return department.ID, nil
//}
