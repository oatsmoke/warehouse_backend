package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	queries "github.com/oatsmoke/warehouse_backend/internal/db"
	"github.com/oatsmoke/warehouse_backend/internal/dto"
	"github.com/oatsmoke/warehouse_backend/internal/lib/logger"
	"github.com/oatsmoke/warehouse_backend/internal/model"
)

type EmployeeRepository struct {
	queries queries.Querier
}

func NewEmployeeRepository(queries queries.Querier) *EmployeeRepository {
	return &EmployeeRepository{
		queries: queries,
	}
}

func (r *EmployeeRepository) Create(ctx context.Context, employee *model.Employee) (int64, error) {
	req, err := r.queries.CreateEmployee(ctx, &queries.CreateEmployeeParams{
		LastName:   employee.LastName,
		FirstName:  employee.FirstName,
		MiddleName: employee.MiddleName,
		Phone:      employee.Phone,
	})
	if err != nil {
		return 0, logger.Error(logger.MsgFailedToInsert, err)
	}

	return req.ID, nil
}

func (r *EmployeeRepository) Read(ctx context.Context, id int64) (*model.Employee, error) {
	req, err := r.queries.ReadEmployee(ctx, id)
	if err != nil {
		return nil, logger.Error(logger.MsgFailedToScan, err)
	}

	employee := &model.Employee{
		ID:         req.ID,
		LastName:   req.LastName,
		FirstName:  req.FirstName,
		MiddleName: req.MiddleName,
		Phone:      req.Phone,
		Department: &model.Department{
			ID:    validInt64(req.DepartmentID),
			Title: validString(req.DepartmentTitle),
		},
		DeletedAt: validTime(req.DeletedAt),
	}

	return employee, nil
}

func (r *EmployeeRepository) Update(ctx context.Context, employee *model.Employee) error {
	ct, err := r.queries.UpdateEmployee(ctx, &queries.UpdateEmployeeParams{
		ID:         employee.ID,
		LastName:   employee.LastName,
		FirstName:  employee.FirstName,
		MiddleName: employee.MiddleName,
		Phone:      employee.Phone,
	})
	if err != nil {
		return logger.Error(logger.MsgFailedToUpdate, err)
	}

	if ct.RowsAffected() == 0 {
		return logger.Error(logger.MsgFailedToUpdate, logger.ErrNoRowsAffected)
	}

	return nil
}

func (r *EmployeeRepository) Delete(ctx context.Context, id int64) error {
	ct, err := r.queries.DeleteEmployee(ctx, id)
	if err != nil {
		return logger.Error(logger.MsgFailedToDelete, err)
	}

	if ct.RowsAffected() == 0 {
		return logger.Error(logger.MsgFailedToDelete, logger.ErrNoRowsAffected)
	}

	return nil
}

func (r *EmployeeRepository) Restore(ctx context.Context, id int64) error {
	ct, err := r.queries.RestoreEmployee(ctx, id)
	if err != nil {
		return logger.Error(logger.MsgFailedToRestore, err)
	}

	if ct.RowsAffected() == 0 {
		return logger.Error(logger.MsgFailedToRestore, logger.ErrNoRowsAffected)
	}

	return nil
}

func (r *EmployeeRepository) List(ctx context.Context, qp *dto.QueryParams) ([]*model.Employee, int64, error) {
	req, err := r.queries.ListEmployee(ctx, &queries.ListEmployeeParams{
		WithDeleted:      qp.WithDeleted,
		Search:           qp.Search,
		Ids:              qp.IDs,
		SortColumn:       qp.SortColumn,
		SortOrder:        qp.SortOrder,
		PaginationLimit:  qp.PaginationLimit,
		PaginationOffset: qp.PaginationOffset,
	})
	if err != nil {
		return nil, 0, logger.Error(logger.MsgFailedToSelect, err)
	}

	if len(req) < 1 {
		return []*model.Employee{}, 0, nil
	}

	list := make([]*model.Employee, len(req))
	for i, item := range req {
		employee := &model.Employee{
			ID:         item.ID,
			LastName:   item.LastName,
			FirstName:  item.FirstName,
			MiddleName: item.MiddleName,
			Phone:      item.Phone,
			Department: &model.Department{
				ID:    validInt64(item.DepartmentID),
				Title: validString(item.DepartmentTitle),
			},
			DeletedAt: validTime(item.DeletedAt),
		}
		list[i] = employee
	}

	return list, req[0].Total, nil
}

func (r *EmployeeRepository) SetDepartment(ctx context.Context, id, departmentID int64) error {
	var d pgtype.Int8
	if departmentID != 0 {
		d = pgtype.Int8{
			Int64: departmentID,
			Valid: true,
		}
	}

	ct, err := r.queries.SetDepartmentEmployee(ctx, &queries.SetDepartmentEmployeeParams{
		ID:           id,
		DepartmentID: d,
	})
	if err != nil {
		return logger.Error(logger.MsgFailedToUpdate, err)
	}

	if ct.RowsAffected() == 0 {
		return logger.Error(logger.MsgFailedToUpdate, logger.ErrNoRowsAffected)
	}

	return nil
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
