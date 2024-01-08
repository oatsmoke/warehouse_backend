package repository

import (
	"context"
	"github.com/jackc/pgx/pgxpool"
	"warehouse_backend/pkg/model"
)

type DepartmentRepository struct {
	db *pgxpool.Pool
}

func NewDepartmentRepository(db *pgxpool.Pool) *DepartmentRepository {
	return &DepartmentRepository{db: db}
}

func (r *DepartmentRepository) Create(title string) error {
	query := `
			INSERT INTO departments (title) 
			VALUES ($1);`
	_, err := r.db.Exec(context.Background(), query, title)
	if err != nil {
		return err
	}
	return nil
}

func (r *DepartmentRepository) GetById(id int) (model.Department, error) {
	var department model.Department
	query := `
			SELECT id, title
			FROM departments 
			WHERE id = $1;`
	err := r.db.QueryRow(context.Background(), query, id).Scan(&department.Id, &department.Title)
	if err != nil {
		return model.Department{}, err
	}
	return department, err
}

func (r *DepartmentRepository) GetAll() ([]model.Department, error) {
	var departments []model.Department
	var department model.Department
	query := `
			SELECT id, title
			FROM departments
			WHERE is_deleted = false
			ORDER BY title;`
	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&department.Id, &department.Title)
		if err != nil {
			return nil, err
		}
		departments = append(departments, department)
	}
	return departments, err
}

func (r *DepartmentRepository) GetAllButOne(id, employeeId int) ([]model.Department, error) {
	var departments []model.Department
	var department model.Department
	query := `
			SELECT departments.id, departments.title
			FROM departments
         	LEFT JOIN employees on departments.id = employees.department
			WHERE departments.id != $1
  			AND departments.is_deleted = false
  			AND employees.id = $2
			AND departments.id = employees.department
			ORDER BY title;`
	rows, err := r.db.Query(context.Background(), query, id, employeeId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&department.Id, &department.Title)
		if err != nil {
			return nil, err
		}
		departments = append(departments, department)
	}
	return departments, err
}

func (r *DepartmentRepository) GetAllButOneForAdmin(id int) ([]model.Department, error) {
	var departments []model.Department
	var department model.Department
	query := `
			SELECT departments.id, departments.title
			FROM departments
			WHERE departments.id != $1
  			AND departments.is_deleted = false
			ORDER BY title;`
	rows, err := r.db.Query(context.Background(), query, id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&department.Id, &department.Title)
		if err != nil {
			return nil, err
		}
		departments = append(departments, department)
	}
	return departments, err
}

func (r *DepartmentRepository) FindByTitle(title string) (int, error) {
	var department model.Department
	query := `
			SELECT id 
			FROM departments 
			WHERE title = $1;`
	err := r.db.QueryRow(context.Background(), query, title).Scan(&department.Id)
	if err != nil {
		return 0, err
	}
	return department.Id, nil
}

func (r *DepartmentRepository) Update(id int, title string) error {
	query := `
			UPDATE departments 
			SET title = $2
			WHERE id = $1;`
	_, err := r.db.Exec(context.Background(), query, id, title)
	if err != nil {
		return err
	}
	return nil
}

func (r *DepartmentRepository) Delete(id int) error {
	query := `
			UPDATE departments 
			SET is_deleted = true
       		WHERE id = $1;`
	_, err := r.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}
	return nil
}
