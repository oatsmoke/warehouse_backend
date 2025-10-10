package repository

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oatsmoke/warehouse_backend/internal/lib/logger"
	"github.com/oatsmoke/warehouse_backend/internal/model"
)

type UserRepository struct {
	postgresDB *pgxpool.Pool
}

func NewUserRepository(postgresDB *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		postgresDB: postgresDB,
	}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) (int64, error) {
	const query = `
		INSERT INTO users (username,  password_hash, email, role, employee)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;`

	var id int64
	var employeeID pgtype.Int8
	if user.Employee != nil && user.Employee.ID != 0 {
		employeeID = pgtype.Int8{
			Int64: user.Employee.ID,
			Valid: true,
		}
	}

	if err := r.postgresDB.QueryRow(
		ctx,
		query,
		user.Username,
		user.PasswordHash,
		user.Email,
		user.Role,
		employeeID,
	).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *UserRepository) Read(ctx context.Context, id int64) (*model.User, error) {
	const query = `
		SELECT u.id, u.username, u.email, u.role, u.enabled, u.last_login_at,
		       e.id, e.last_name, e.first_name, e.middle_name, e.phone,
		       d.id, d.title
		FROM users u
		LEFT JOIN public.employees e ON e.id = u.employee
		LEFT JOIN public.departments d on d.id = e.department
		WHERE u.id = $1;`

	user := model.NewUser()
	var (
		employeeID, employeeDepartmentID                                                                sql.NullInt64
		employeeLastName, employeeFirstName, employeeMiddleName, employeePhone, employeeDepartmentTitle sql.NullString
	)

	if err := r.postgresDB.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Role,
		&user.Enabled,
		&user.LastLoginAt,
		&employeeID,
		&employeeLastName,
		&employeeFirstName,
		&employeeMiddleName,
		&employeePhone,
		&employeeDepartmentID,
		&employeeDepartmentTitle,
	); err != nil {
		return nil, err
	}

	user.Employee.ID = validInt64(employeeID)
	user.Employee.LastName = validString(employeeLastName)
	user.Employee.FirstName = validString(employeeFirstName)
	user.Employee.MiddleName = validString(employeeMiddleName)
	user.Employee.Phone = validString(employeePhone)
	user.Employee.Department.ID = validInt64(employeeDepartmentID)
	user.Employee.Department.Title = validString(employeeDepartmentTitle)

	return user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	const query = `
		UPDATE users
		SET username = $2, email = $3, role = $4, employee = $5
		WHERE id = $1;`

	employeeID := new(pgtype.Int8)
	if user.Employee != nil {
		employeeID = &pgtype.Int8{
			Int64: user.Employee.ID,
			Valid: user.Employee.ID != 0,
		}
	}

	ct, err := r.postgresDB.Exec(
		ctx,
		query,
		user.ID,
		user.Username,
		user.Email,
		user.Role,
		employeeID,
	)
	if err != nil {
		return err
	}

	if ct.RowsAffected() == 0 {
		return logger.NoRowsAffected
	}

	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	const query = `
		DELETE FROM users
		WHERE id = $1;`

	ct, err := r.postgresDB.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if ct.RowsAffected() == 0 {
		return logger.NoRowsAffected
	}

	return nil
}

func (r *UserRepository) List(ctx context.Context) ([]*model.User, error) {
	const query = `
		SELECT u.id, u.username, e.email, e.role, u.enabled, u.last_login_at,
		       e.id, e.last_name, e.first_name, e.middle_name, e.phone,
		       d.id, d.title
		FROM users u
		LEFT JOIN public.employees e ON e.id = u.employee
		LEFT JOIN public.departments d on d.id = e.department
		ORDER BY u.id;`

	user := model.NewUser()
	var (
		employeeID, employeeDepartmentID                                                                sql.NullInt64
		employeeLastName, employeeFirstName, employeeMiddleName, employeePhone, employeeDepartmentTitle sql.NullString
	)

	rows, err := r.postgresDB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		user = model.NewUser()
		if err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.Role,
			&user.Enabled,
			&user.LastLoginAt,
			&employeeID,
			&employeeLastName,
			&employeeFirstName,
			&employeeMiddleName,
			&employeePhone,
			&employeeDepartmentID,
			&employeeDepartmentTitle,
		); err != nil {
			return nil, err
		}

		user.Employee.ID = validInt64(employeeID)
		user.Employee.LastName = validString(employeeLastName)
		user.Employee.FirstName = validString(employeeFirstName)
		user.Employee.MiddleName = validString(employeeMiddleName)
		user.Employee.Phone = validString(employeePhone)
		user.Employee.Department.ID = validInt64(employeeDepartmentID)
		user.Employee.Department.Title = validString(employeeDepartmentTitle)
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
