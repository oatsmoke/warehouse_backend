package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
	"warehouse_backend/internal/model"
)

type EmployeeRepository struct {
	DB *pgxpool.Pool
}

func NewEmployeeRepository(db *pgxpool.Pool) *EmployeeRepository {
	return &EmployeeRepository{DB: db}
}

// Create is employee create
func (r *EmployeeRepository) Create(ctx context.Context, name, phone, email string) error {
	date := time.Now()

	const query = `
		INSERT INTO employees (name, phone, email, password, hash, registration_date, authorization_date)
		VALUES ($1, $2, $3, $4, $5, $6, $7);`

	if _, err := r.DB.Exec(ctx, query, name, phone, email, "", "", date, date); err != nil {
		return err
	}

	return nil
}

// Update is employee update
func (r *EmployeeRepository) Update(ctx context.Context, id int64, name, phone, email string) error {
	const query = `
		UPDATE employees 
		SET name = $2, phone = $3, email = $4 
		WHERE id = $1;`

	if _, err := r.DB.Exec(ctx, query, id, name, phone, email); err != nil {
		return err
	}

	return nil
}

// Delete is employee delete
func (r *EmployeeRepository) Delete(ctx context.Context, id int64) error {
	const query = `
		UPDATE employees 
		SET password = $2, hash = $3, activate = false, deleted = true
		WHERE id = $1;`

	if _, err := r.DB.Exec(ctx, query, id, "", ""); err != nil {
		return err
	}

	return nil
}

// Restore is employee restore
func (r *EmployeeRepository) Restore(ctx context.Context, id int64) error {
	const query = `
		UPDATE employees 
		SET deleted = false
		WHERE id = $1;`

	if _, err := r.DB.Exec(ctx, query, id); err != nil {
		return err
	}

	return nil
}

// GetAll is employee get all
func (r *EmployeeRepository) GetAll(ctx context.Context, deleted bool) ([]*model.Employee, error) {
	var employees []*model.Employee
	query := ""

	if deleted {
		query = `
			SELECT id, name, phone, email, registration_date, authorization_date, activate
			FROM employees 
			WHERE hidden = false 
			ORDER BY name;`
	} else {
		query = `
			SELECT id, name, phone, email, registration_date, authorization_date, activate
			FROM employees 
			WHERE hidden = false AND deleted = false 
			ORDER BY name;`
	}

	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		employee := new(model.Employee)
		if err := rows.Scan(
			&employee.ID,
			&employee.Name,
			&employee.Phone,
			&employee.Email,
			&employee.RegistrationDate,
			&employee.AuthorizationDate,
			&employee.Activate,
		); err != nil {
			return nil, err
		}
		employees = append(employees, employee)
	}

	return employees, err
}

// GetAllShort is employee get all short
func (r *EmployeeRepository) GetAllShort(ctx context.Context, deleted bool) ([]*model.Employee, error) {
	var employees []*model.Employee

	const query = `
		SELECT id, name
		FROM employees 
		WHERE hidden = false AND deleted = false 
		ORDER BY name;`

	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		employee := new(model.Employee)
		if err := rows.Scan(
			&employee.ID,
			&employee.Name); err != nil {
			return nil, err
		}
		employees = append(employees, employee)
	}

	return employees, err
}

// GetAllButOne is employee get all but one
func (r *EmployeeRepository) GetAllButOne(ctx context.Context, id int64, deleted bool) ([]*model.Employee, error) {
	var employees []*model.Employee
	query := ""

	if deleted {
		query = `
			SELECT id, name, phone, email, registration_date, authorization_date, activate, role, deleted
			FROM employees 
			WHERE hidden = false AND id != $1 
			ORDER BY name;`
	} else {
		query = `
			SELECT id, name, phone, email, registration_date, authorization_date, activate, role, deleted
			FROM employees 
			WHERE hidden = false AND deleted = false AND id != $1 
			ORDER BY name;`
	}

	rows, err := r.DB.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		employee := new(model.Employee)
		if err := rows.Scan(
			&employee.ID,
			&employee.Name,
			&employee.Phone,
			&employee.Email,
			&employee.RegistrationDate,
			&employee.AuthorizationDate,
			&employee.Activate,
			&employee.Role,
			&employee.Deleted,
		); err != nil {
			return nil, err
		}
		employees = append(employees, employee)
	}

	return employees, nil
}

// GetById is employee get by id
func (r *EmployeeRepository) GetById(ctx context.Context, employee *model.Employee) (*model.Employee, error) {
	const query = `
		SELECT name, phone, email, registration_date, authorization_date, activate, role
		FROM employees 
		WHERE id = $1;`

	if err := r.DB.QueryRow(ctx, query, employee.ID).Scan(
		&employee.Name,
		&employee.Phone,
		&employee.Email,
		&employee.RegistrationDate,
		&employee.AuthorizationDate,
		&employee.Activate,
		&employee.Role,
	); err != nil {
		return nil, err
	}

	return employee, nil
}

// GetFree is employee get free
func (r *EmployeeRepository) GetFree(ctx context.Context) ([]*model.Employee, error) {
	var employees []*model.Employee

	const query = `
		SELECT id, name
		FROM employees
		WHERE hidden = false AND deleted = false AND department IS NULL
		ORDER BY name;`

	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		employee := new(model.Employee)
		if err := rows.Scan(
			&employee.ID,
			&employee.Name,
		); err != nil {
			return nil, err
		}
		employees = append(employees, employee)
	}

	return employees, nil
}

// GetByDepartment is employee get by department
func (r *EmployeeRepository) GetByDepartment(ctx context.Context, ids []int64, departmentId int64) ([]*model.Employee, error) {
	var employees []*model.Employee

	const query = `
		SELECT id, name 
		FROM employees 
		WHERE department = $1
		ORDER BY name;`
	//AND id = ANY($1)
	rows, err := r.DB.Query(ctx, query, departmentId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		employee := new(model.Employee)
		if err := rows.Scan(&employee.ID, &employee.Name); err != nil {
			return nil, err
		}
		employees = append(employees, employee)
	}

	return employees, nil
}

// AddToDepartment is employee add to department
func (r *EmployeeRepository) AddToDepartment(ctx context.Context, id, departmentId int64) error {
	const query = `
		UPDATE employees 
		SET department = $2 
		WHERE id = $1;`

	if _, err := r.DB.Exec(ctx, query, id, departmentId); err != nil {
		return err
	}

	return nil
}

// RemoveFromDepartment is employee remove from department
func (r *EmployeeRepository) RemoveFromDepartment(ctx context.Context, id int64) error {
	const query = `
		UPDATE employees 
		SET department = NULL 
		WHERE id = $1;`

	if _, err := r.DB.Exec(ctx, query, id); err != nil {
		return err
	}

	return nil
}

// Activate is employee activate
func (r *EmployeeRepository) Activate(ctx context.Context, id int64, password string) error {
	const query = `
		UPDATE employees 
		SET password = $2, activate = true
		WHERE id = $1;`

	if _, err := r.DB.Exec(ctx, query, id, password); err != nil {
		return err
	}

	return nil
}

// Deactivate is employee deactivate
func (r *EmployeeRepository) Deactivate(ctx context.Context, id int64) error {
	const query = `
		UPDATE employees 
		SET password = $2, hash = $3, activate = false
		WHERE id = $1;`

	if _, err := r.DB.Exec(ctx, query, id, "", ""); err != nil {
		return err
	}

	return nil
}

// ResetPassword is employee reset password
func (r *EmployeeRepository) ResetPassword(ctx context.Context, id int64, password string) error {
	const query = `
		UPDATE employees 
		SET password = $2
		WHERE id = $1;`

	if _, err := r.DB.Exec(ctx, query, id, password); err != nil {
		return err
	}

	return nil
}

// ChangeRole is employee change role
func (r *EmployeeRepository) ChangeRole(ctx context.Context, id int64, role string) error {
	const query = `
		UPDATE employees 
		SET role = $2
		WHERE id = $1;`

	if _, err := r.DB.Exec(ctx, query, id, role); err != nil {
		return err
	}

	return nil
}
