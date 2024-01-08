package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/pgxpool"
	"time"
	"warehouse_backend/pkg/model"
)

type EmployeeRepository struct {
	db *pgxpool.Pool
}

func NewEmployeeRepository(db *pgxpool.Pool) *EmployeeRepository {
	return &EmployeeRepository{db: db}
}

func (r *EmployeeRepository) Create(name, phone, email string) error {
	date := time.Now()
	query := `
			INSERT INTO employees (name, phone, email, password, hash, registration_date, authorization_date)
			VALUES ($1, $2, $3, $4, $5, $6, $7);`
	_, err := r.db.Exec(context.Background(), query, name, phone, email, "", "", date, date)
	if err != nil {
		return err
	}
	return nil
}

func (r *EmployeeRepository) GetById(id int) (model.Employee, error) {
	var employee model.Employee
	var rDate, aDate time.Time
	var department interface{}
	query := `
			SELECT name, phone, email, registration_date, authorization_date, activate,department, role
			FROM employees 
			WHERE id = $1;`
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&employee.Name,
		&employee.Phone,
		&employee.Email,
		&rDate,
		&aDate,
		&employee.Activate,
		&department,
		&employee.Role)
	employee.RegistrationDate = rDate.Unix()
	employee.AuthorizationDate = aDate.Unix()
	employee.Department.Id = InterfaceToInt(department)
	if err != nil {
		return model.Employee{}, err
	}
	return employee, err
}

func (r *EmployeeRepository) GetByDepartment(ids []int, id int) ([]model.Employee, error) {
	var employees []model.Employee
	var employee model.Employee
	str := ""
	for _, id := range ids {
		str = fmt.Sprintf(" AND id != %d %s", id, str)
	}
	query := fmt.Sprintf("SELECT id, name FROM employees WHERE department = $1 %s ORDER BY name;", str)
	rows, err := r.db.Query(context.Background(), query, id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(
			&employee.Id,
			&employee.Name)
		if err != nil {
			return nil, err
		}
		employees = append(employees, employee)
	}
	return employees, err
}

func (r *EmployeeRepository) GetAll() ([]model.Employee, error) {
	var employees []model.Employee
	var employee model.Employee
	var rDate, aDate time.Time
	query := `
			SELECT id, name, phone, email, registration_date, authorization_date, activate
			FROM employees 
			WHERE hidden = false AND is_deleted = false 
			ORDER BY name;`
	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(
			&employee.Id,
			&employee.Name,
			&employee.Phone,
			&employee.Email,
			&rDate,
			&aDate,
			&employee.Activate)
		employee.RegistrationDate = rDate.Unix()
		employee.AuthorizationDate = aDate.Unix()
		if err != nil {
			return nil, err
		}
		employees = append(employees, employee)
	}
	return employees, err
}

func (r *EmployeeRepository) GetFree() ([]model.Employee, error) {
	var employees []model.Employee
	var employee model.Employee
	query := `
			SELECT id, name
			FROM employees
			WHERE hidden = false AND is_deleted = false AND department IS NULL
			ORDER BY name;`
	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(
			&employee.Id,
			&employee.Name)
		if err != nil {
			return nil, err
		}
		employees = append(employees, employee)
	}
	return employees, err
}

func (r *EmployeeRepository) GetAllButOne(id int) ([]model.Employee, error) {
	var employees []model.Employee
	var employee model.Employee
	var rDate, aDate time.Time
	query := `
			SELECT id, name, phone, email, registration_date, authorization_date, activate, role
			FROM employees
			WHERE hidden = false AND id != $1
			ORDER BY name;`
	rows, err := r.db.Query(context.Background(), query, id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(
			&employee.Id,
			&employee.Name,
			&employee.Phone,
			&employee.Email,
			&rDate,
			&aDate,
			&employee.Activate,
			&employee.Role)
		employee.RegistrationDate = rDate.Unix()
		employee.AuthorizationDate = aDate.Unix()
		if err != nil {
			return nil, err
		}
		employees = append(employees, employee)
	}
	return employees, err
}

func (r *EmployeeRepository) FindUser(login, password string) (int, error) {
	var user model.Employee
	query := `
			SELECT id 
			FROM employees 
			WHERE phone = $1 AND password = $2;`
	err := r.db.QueryRow(context.Background(), query, login, password).Scan(&user.Id)
	if err != nil {
		return 0, err
	}
	return user.Id, nil
}

func (r *EmployeeRepository) FindByPhone(phone string) (int, error) {
	var user model.Employee
	query := `
			SELECT id 
			FROM employees 
			WHERE phone = $1;`
	err := r.db.QueryRow(context.Background(), query, phone).Scan(&user.Id)
	if err != nil {
		return 0, err
	}
	return user.Id, nil
}

func (r *EmployeeRepository) FindByHash(hash string) (int, error) {
	var user model.Employee
	query := `
			SELECT id 
			FROM employees 
			WHERE hash = $1;`
	err := r.db.QueryRow(context.Background(), query, hash).Scan(&user.Id)
	if err != nil {
		return 0, err
	}
	return user.Id, nil
}

func (r *EmployeeRepository) SetHash(id int, hash string) error {
	date := time.Now()
	query := `
			UPDATE employees 
			SET hash = $2, authorization_date = $3
			WHERE id = $1;`
	_, err := r.db.Exec(context.Background(), query, id, hash, date)
	if err != nil {
		return err
	}
	return nil
}
func (r *EmployeeRepository) AddToDepartment(id, department int) error {
	query := `
			UPDATE employees 
			SET department = $2 
			WHERE id = $1;`
	_, err := r.db.Exec(context.Background(), query, id, department)
	if err != nil {
		return err
	}
	return nil
}

func (r *EmployeeRepository) RemoveFromDepartment(id int) error {
	query := `
			UPDATE employees 
			SET department = NULL 
			WHERE id = $1;`
	_, err := r.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *EmployeeRepository) Update(id int, name, phone, email string) error {
	query := `
			UPDATE employees 
			SET name = $2, phone = $3, email = $4 
			WHERE id = $1;`
	_, err := r.db.Exec(context.Background(), query, id, name, phone, email)
	if err != nil {
		return err
	}
	return nil
}

func (r *EmployeeRepository) Delete(id int) error {
	query := `
			UPDATE employees 
			SET password = $2, hash = $3, activate = false, is_deleted = true
			WHERE id = $1;`
	_, err := r.db.Exec(context.Background(), query, id, "", "")
	if err != nil {
		return err
	}
	return nil
}

func (r *EmployeeRepository) Activate(id int, password string) error {
	query := `
			UPDATE employees 
			SET password = $2, activate = true
			WHERE id = $1;`
	_, err := r.db.Exec(context.Background(), query, id, password)
	if err != nil {
		return err
	}
	return nil
}

func (r *EmployeeRepository) Deactivate(id int) error {
	query := `
			UPDATE employees 
			SET password = $2, hash = $3, activate = false
			WHERE id = $1;`
	_, err := r.db.Exec(context.Background(), query, id, "", "")
	if err != nil {
		return err
	}
	return nil
}

func (r *EmployeeRepository) ResetPassword(id int, password string) error {
	query := `
			UPDATE employees 
			SET password = $2
			WHERE id = $1;`
	_, err := r.db.Exec(context.Background(), query, id, password)
	if err != nil {
		return err
	}
	return nil
}

func (r *EmployeeRepository) ChangeRole(id int, role string) error {
	query := `
			UPDATE employees 
			SET role = $2
			WHERE id = $1;`
	_, err := r.db.Exec(context.Background(), query, id, role)
	if err != nil {
		return err
	}
	return nil
}
