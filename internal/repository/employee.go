package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
	"warehouse_backend/internal/model"
)

type EmployeeRepository struct {
	db *pgxpool.Pool
}

func NewEmployeeRepository(db *pgxpool.Pool) *EmployeeRepository {
	return &EmployeeRepository{db: db}
}

func (r *EmployeeRepository) Create(ctx context.Context, name, phone, email string) error {
	date := time.Now()

	query := `
			INSERT INTO employees (name, phone, email, password, hash, registration_date, authorization_date)
			VALUES ($1, $2, $3, $4, $5, $6, $7);`

	if _, err := r.db.Exec(ctx, query, name, phone, email, "", "", date, date); err != nil {
		return err
	}

	return nil
}

func (r *EmployeeRepository) GetById(ctx context.Context, id int64) (*model.Employee, error) {
	employee := new(model.Employee)
	query := `
			SELECT name, phone, email, registration_date, authorization_date, activate, department, role
			FROM employees 
			WHERE id = $1;`

	if err := r.db.QueryRow(ctx, query, id).Scan(
		&employee.Name,
		&employee.Phone,
		&employee.Email,
		&employee.RegistrationDate,
		&employee.AuthorizationDate,
		&employee.Activate,
		&employee.Department,
		&employee.Role); err != nil {
		return nil, err
	}

	return employee, nil
}

func (r *EmployeeRepository) GetByDepartment(ctx context.Context, ids []int64, departmentId int64) ([]*model.Employee, error) {
	var employees []*model.Employee
	employee := new(model.Employee)

	query := `
			SELECT id, name 
			FROM employees 
			WHERE department = $2 AND id IN $2 
			ORDER BY name;`

	rows, err := r.db.Query(ctx, query, ids, departmentId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		if err := rows.Scan(&employee.ID, &employee.Name); err != nil {
			return nil, err
		}
		employees = append(employees, employee)
	}

	return employees, nil
}

func (r *EmployeeRepository) GetAll(ctx context.Context) ([]*model.Employee, error) {
	var employees []*model.Employee
	employee := new(model.Employee)

	query := `
			SELECT id, name, phone, email, registration_date, authorization_date, activate
			FROM employees 
			WHERE hidden = false AND is_deleted = false 
			ORDER BY name;`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		if err := rows.Scan(
			&employee.ID,
			&employee.Name,
			&employee.Phone,
			&employee.Email,
			&employee.RegistrationDate,
			&employee.AuthorizationDate,
			&employee.Activate); err != nil {
			return nil, err
		}
		employees = append(employees, employee)
	}

	return employees, err
}

func (r *EmployeeRepository) GetFree(ctx context.Context) ([]*model.Employee, error) {
	var employees []*model.Employee
	employee := new(model.Employee)

	query := `
			SELECT id, name
			FROM employees
			WHERE hidden = false AND is_deleted = false AND department IS NULL
			ORDER BY name;`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		if err := rows.Scan(&employee.ID, &employee.Name); err != nil {
			return nil, err
		}
		employees = append(employees, employee)
	}

	return employees, nil
}

func (r *EmployeeRepository) GetAllButOne(ctx context.Context, id int64) ([]*model.Employee, error) {
	var employees []*model.Employee
	employee := new(model.Employee)

	query := `
			SELECT id, name, phone, email, registration_date, authorization_date, activate, role
			FROM employees
			WHERE hidden = false AND id != $1
			ORDER BY name;`

	rows, err := r.db.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		if err := rows.Scan(
			&employee.ID,
			&employee.Name,
			&employee.Phone,
			&employee.Email,
			&employee.RegistrationDate,
			&employee.AuthorizationDate,
			&employee.Activate,
			&employee.Role); err != nil {
			return nil, err
		}
		employees = append(employees, employee)
	}

	return employees, nil
}

func (r *EmployeeRepository) AddToDepartment(ctx context.Context, id, departmentId int64) error {
	query := `
			UPDATE employees 
			SET department = $2 
			WHERE id = $1;`

	if _, err := r.db.Exec(ctx, query, id, departmentId); err != nil {
		return err
	}

	return nil
}

func (r *EmployeeRepository) RemoveFromDepartment(ctx context.Context, id int64) error {
	query := `
			UPDATE employees 
			SET department = NULL 
			WHERE id = $1;`

	if _, err := r.db.Exec(ctx, query, id); err != nil {
		return err
	}

	return nil
}

func (r *EmployeeRepository) Update(ctx context.Context, id int64, name, phone, email string) error {
	query := `
			UPDATE employees 
			SET name = $2, phone = $3, email = $4 
			WHERE id = $1;`

	if _, err := r.db.Exec(ctx, query, id, name, phone, email); err != nil {
		return err
	}

	return nil
}

func (r *EmployeeRepository) Delete(ctx context.Context, id int64) error {
	query := `
			UPDATE employees 
			SET password = $2, hash = $3, activate = false, is_deleted = true
			WHERE id = $1;`

	if _, err := r.db.Exec(ctx, query, id, "", ""); err != nil {
		return err
	}

	return nil
}

func (r *EmployeeRepository) Activate(ctx context.Context, id int64, password string) error {
	query := `
			UPDATE employees 
			SET password = $2, activate = true
			WHERE id = $1;`

	if _, err := r.db.Exec(ctx, query, id, password); err != nil {
		return err
	}

	return nil
}

func (r *EmployeeRepository) Deactivate(ctx context.Context, id int64) error {
	query := `
			UPDATE employees 
			SET password = $2, hash = $3, activate = false
			WHERE id = $1;`

	if _, err := r.db.Exec(ctx, query, id, "", ""); err != nil {
		return err
	}

	return nil
}

func (r *EmployeeRepository) ResetPassword(ctx context.Context, id int64, password string) error {
	query := `
			UPDATE employees 
			SET password = $2
			WHERE id = $1;`

	if _, err := r.db.Exec(ctx, query, id, password); err != nil {
		return err
	}

	return nil
}

func (r *EmployeeRepository) ChangeRole(ctx context.Context, id int64, role string) error {
	query := `
			UPDATE employees 
			SET role = $2
			WHERE id = $1;`

	if _, err := r.db.Exec(ctx, query, id, role); err != nil {
		return err
	}

	return nil
}
