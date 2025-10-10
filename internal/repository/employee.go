package repository

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oatsmoke/warehouse_backend/internal/dto"
	"github.com/oatsmoke/warehouse_backend/internal/lib/list_filter"
	"github.com/oatsmoke/warehouse_backend/internal/lib/logger"
	"github.com/oatsmoke/warehouse_backend/internal/model"
)

type EmployeeRepository struct {
	postgresDB *pgxpool.Pool
}

func NewEmployeeRepository(postgresDB *pgxpool.Pool) *EmployeeRepository {
	return &EmployeeRepository{
		postgresDB: postgresDB,
	}
}

func (r *EmployeeRepository) Create(ctx context.Context, employee *model.Employee) (int64, error) {
	const query = `
		INSERT INTO employees (last_name, first_name, middle_name, phone)
		VALUES ($1, $2, $3, $4)
		RETURNING id;`

	var id int64
	if err := r.postgresDB.QueryRow(
		ctx,
		query,
		employee.LastName,
		employee.FirstName,
		employee.MiddleName,
		employee.Phone,
	).Scan(&id); err != nil {
		return 0, err
	}

	if id == 0 {
		return 0, logger.NoRowsAffected
	}

	return id, nil
}

func (r *EmployeeRepository) Read(ctx context.Context, id int64) (*model.Employee, error) {
	const query = `
		SELECT e.id, e.last_name, e.first_name, e.middle_name, e.phone, e.deleted_at,
		       d.id, d.title
		FROM employees e
		LEFT JOIN departments d ON d.id = e.department
		WHERE e.id = $1;`

	employee := model.NewEmployee()
	var (
		departmentID    sql.NullInt64
		departmentTitle sql.NullString
	)

	if err := r.postgresDB.QueryRow(ctx, query, id).Scan(
		&employee.ID,
		&employee.LastName,
		&employee.FirstName,
		&employee.MiddleName,
		&employee.Phone,
		&employee.DeletedAt,
		&departmentID,
		&departmentTitle,
	); err != nil {
		return nil, err
	}

	employee.Department.ID = validInt64(departmentID)
	employee.Department.Title = validString(departmentTitle)

	return employee, nil
}

func (r *EmployeeRepository) Update(ctx context.Context, employee *model.Employee) error {
	const query = `
		UPDATE employees
		SET last_name = $2, first_name = $3, middle_name = $4, phone = $5, department = $6
		WHERE id = $1;`

	var departmentID pgtype.Int8
	if employee.Department != nil && employee.Department.ID != 0 {
		departmentID.Int64 = employee.Department.ID
		departmentID.Valid = true
	}

	ct, err := r.postgresDB.Exec(
		ctx,
		query,
		employee.ID,
		employee.LastName,
		employee.FirstName,
		employee.MiddleName,
		employee.Phone,
		departmentID,
	)
	if err != nil {
		return err
	}

	if ct.RowsAffected() == 0 {
		return logger.NoRowsAffected
	}

	return nil
}

func (r *EmployeeRepository) Delete(ctx context.Context, id int64) error {
	const query = `
		UPDATE employees 
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

func (r *EmployeeRepository) Restore(ctx context.Context, id int64) error {
	const query = `
		UPDATE employees 
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

func (r *EmployeeRepository) List(ctx context.Context, qp *dto.QueryParams) ([]*model.Employee, error) {
	fields := []string{"e.last_name", "e.first_name", "e.middle_name", "e.phone", "e.email", "d.title"}
	str, args := list_filter.BuildQuery(qp, fields, "e")

	query := `
		SELECT e.id, e.last_name, e.first_name, e.middle_name, e.phone, e.deleted_at,
		       d.id, d.title
		FROM employees e
		LEFT JOIN public.departments d ON d.id = e.department
		` + str

	rows, err := r.postgresDB.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []*model.Employee
	for rows.Next() {
		employee := model.NewEmployee()
		var (
			departmentID    sql.NullInt64
			departmentTitle sql.NullString
		)

		if err := rows.Scan(
			&employee.ID,
			&employee.LastName,
			&employee.FirstName,
			&employee.MiddleName,
			&employee.Phone,
			&employee.DeletedAt,
			&departmentID,
			&departmentTitle,
		); err != nil {
			return nil, err
		}

		employee.Department.ID = validInt64(departmentID)
		employee.Department.Title = validString(departmentTitle)
		employees = append(employees, employee)
	}

	return employees, err
}

// GetAllShort is employee get all short
//func (r *EmployeeRepository) GetAllShort(ctx context.Context, deleted bool) ([]*model.Employee, error) {
//	var employees []*model.Employee
//
//	const query = `
//		SELECT id, name
//		FROM employees
//		WHERE hidden = false AND deleted = false
//		ORDER BY name;`
//
//	rows, err := r.DB.Query(ctx, query)
//	if err != nil {
//		return nil, err
//	}
//
//	for rows.Next() {
//		employee := new(model.Employee)
//		if err := rows.Scan(
//			&employee.ID,
//			&employee.Name); err != nil {
//			return nil, err
//		}
//		employees = append(employees, employee)
//	}
//
//	return employees, err
//}
//
//// GetAllButOne is employee get all but one
//func (r *EmployeeRepository) GetAllButOne(ctx context.Context, id int64, deleted bool) ([]*model.Employee, error) {
//	var employees []*model.Employee
//	query := ""
//
//	if deleted {
//		query = `
//			SELECT id, name, phone, email, registration_date, authorization_date, activate, role, deleted
//			FROM employees
//			WHERE hidden = false AND id != $1
//			ORDER BY name;`
//	} else {
//		query = `
//			SELECT id, name, phone, email, registration_date, authorization_date, activate, role, deleted
//			FROM employees
//			WHERE hidden = false AND deleted = false AND id != $1
//			ORDER BY name;`
//	}
//
//	rows, err := r.DB.Query(ctx, query, id)
//	if err != nil {
//		return nil, err
//	}
//
//	for rows.Next() {
//		employee := new(model.Employee)
//		if err := rows.Scan(
//			&employee.ID,
//			&employee.Name,
//			&employee.Phone,
//			&employee.Email,
//			&employee.RegistrationDate,
//			&employee.AuthorizationDate,
//			&employee.Activate,
//			&employee.Role,
//			&employee.Deleted,
//		); err != nil {
//			return nil, err
//		}
//		employees = append(employees, employee)
//	}
//
//	return employees, nil
//}
//
//// GetById is employee get by id
//
//// GetFree is employee get free
//func (r *EmployeeRepository) GetFree(ctx context.Context) ([]*model.Employee, error) {
//	var employees []*model.Employee
//
//	const query = `
//		SELECT id, name
//		FROM employees
//		WHERE hidden = false AND deleted = false AND department IS NULL
//		ORDER BY name;`
//
//	rows, err := r.DB.Query(ctx, query)
//	if err != nil {
//		return nil, err
//	}
//
//	for rows.Next() {
//		employee := new(model.Employee)
//		if err := rows.Scan(
//			&employee.ID,
//			&employee.Name,
//		); err != nil {
//			return nil, err
//		}
//		employees = append(employees, employee)
//	}
//
//	return employees, nil
//}
//
//// GetByDepartment is employee get by department
//func (r *EmployeeRepository) GetByDepartment(ctx context.Context, ids []int64, departmentId int64) ([]*model.Employee, error) {
//	var employees []*model.Employee
//
//	const query = `
//		SELECT id, name
//		FROM employees
//		WHERE department = $1
//		ORDER BY name;`
//	//AND id = ANY($1)
//	rows, err := r.DB.Query(ctx, query, departmentId)
//	if err != nil {
//		return nil, err
//	}
//
//	for rows.Next() {
//		employee := new(model.Employee)
//		if err := rows.Scan(&employee.ID, &employee.Name); err != nil {
//			return nil, err
//		}
//		employees = append(employees, employee)
//	}
//
//	return employees, nil
//}
//
//// AddToDepartment is employee add to department
//func (r *EmployeeRepository) AddToDepartment(ctx context.Context, id, departmentId int64) error {
//	const query = `
//		UPDATE employees
//		SET department = $2
//		WHERE id = $1;`
//
//	if _, err := r.DB.Exec(ctx, query, id, departmentId); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//// RemoveFromDepartment is employee remove from department
//func (r *EmployeeRepository) RemoveFromDepartment(ctx context.Context, id int64) error {
//	const query = `
//		UPDATE employees
//		SET department = NULL
//		WHERE id = $1;`
//
//	if _, err := r.DB.Exec(ctx, query, id); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//// Activate is employee activate
//func (r *EmployeeRepository) Activate(ctx context.Context, id int64, password string) error {
//	const query = `
//		UPDATE employees
//		SET password = $2, activate = true
//		WHERE id = $1;`
//
//	if _, err := r.DB.Exec(ctx, query, id, password); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//// Deactivate is employee deactivate
//func (r *EmployeeRepository) Deactivate(ctx context.Context, id int64) error {
//	const query = `
//		UPDATE employees
//		SET password = $2, hash = $3, activate = false
//		WHERE id = $1;`
//
//	if _, err := r.DB.Exec(ctx, query, id, "", ""); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//// ResetPassword is employee reset password
//func (r *EmployeeRepository) ResetPassword(ctx context.Context, id int64, password string) error {
//	const query = `
//		UPDATE employees
//		SET password = $2
//		WHERE id = $1;`
//
//	if _, err := r.DB.Exec(ctx, query, id, password); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//// ChangeRole is employee change role
//func (r *EmployeeRepository) ChangeRole(ctx context.Context, id int64, role string) error {
//	const query = `
//		UPDATE employees
//		SET role = $2
//		WHERE id = $1;`
//
//	if _, err := r.DB.Exec(ctx, query, id, role); err != nil {
//		return err
//	}
//
//	return nil
//}
