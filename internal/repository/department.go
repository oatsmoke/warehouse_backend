package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oatsmoke/warehouse_backend/internal/dto"
	"github.com/oatsmoke/warehouse_backend/internal/lib/list_filter"
	"github.com/oatsmoke/warehouse_backend/internal/lib/logger"
	"github.com/oatsmoke/warehouse_backend/internal/model"
)

type DepartmentRepository struct {
	postgresDB *pgxpool.Pool
}

func NewDepartmentRepository(postgresDB *pgxpool.Pool) *DepartmentRepository {
	return &DepartmentRepository{
		postgresDB: postgresDB,
	}
}

func (r *DepartmentRepository) Create(ctx context.Context, department *model.Department) (int64, error) {
	const query = `
		INSERT INTO departments (title) 
		VALUES ($1)
		RETURNING id;`

	var id int64
	if err := r.postgresDB.QueryRow(ctx, query, department.Title).Scan(&id); err != nil {
		return 0, err
	}

	if id == 0 {
		return 0, logger.NoRowsAffected
	}

	return id, nil
}

func (r *DepartmentRepository) Read(ctx context.Context, id int64) (*model.Department, error) {
	const query = `
		SELECT id, title, deleted_at
		FROM departments 
		WHERE id = $1;`

	department := new(model.Department)
	if err := r.postgresDB.QueryRow(ctx, query, id).Scan(
		&department.ID,
		&department.Title,
		&department.DeletedAt,
	); err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	return department, nil
}

func (r *DepartmentRepository) Update(ctx context.Context, department *model.Department) error {
	const query = `
		UPDATE departments 
		SET title = $2
		WHERE id = $1 AND title != $2;`

	ct, err := r.postgresDB.Exec(ctx, query, department.ID, department.Title)
	if err != nil {
		return err
	}

	if ct.RowsAffected() == 0 {
		return logger.NoRowsAffected
	}

	return nil
}

func (r *DepartmentRepository) Delete(ctx context.Context, id int64) error {
	const query = `
		UPDATE departments 
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

func (r *DepartmentRepository) Restore(ctx context.Context, id int64) error {
	const query = `
		UPDATE departments 
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

func (r *DepartmentRepository) List(ctx context.Context, qp *dto.QueryParams) ([]*model.Department, error) {
	fields := []string{"title"}
	str, args := list_filter.BuildQuery(qp, fields, "d")

	query := `
		SELECT id, title, deleted_at
		FROM departments d
		` + str

	rows, err := r.postgresDB.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var departments []*model.Department
	for rows.Next() {
		department := new(model.Department)
		if err := rows.Scan(
			&department.ID,
			&department.Title,
			&department.DeletedAt,
		); err != nil {
			return nil, err
		}
		departments = append(departments, department)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return departments, err
}

//func (r *DepartmentRepository) GetAllButOne(ctx context.Context, id, employeeId int64) ([]*model.Department, error) {
//	var departments []*model.Department
//
//	const query = `
//		SELECT departments.id, departments.title
//		FROM departments
//        LEFT JOIN employees on departments.id = employees.department
//		WHERE departments.id != $1
//  		AND departments.deleted = false
//  		AND employees.id = $2
//		AND departments.id = employees.department
//		ORDER BY title;`
//
//	rows, err := r.DB.Query(ctx, query, id, employeeId)
//	if err != nil {
//		return nil, err
//	}
//
//	for rows.Next() {
//		department := new(model.Department)
//		if err := rows.Scan(
//			&department.ID,
//			&department.Title,
//		); err != nil {
//			return nil, err
//		}
//		departments = append(departments, department)
//	}
//
//	return departments, nil
//}
//
//// GetAllButOneForAdmin is department get all but one for admin
//func (r *DepartmentRepository) GetAllButOneForAdmin(ctx context.Context, id int64) ([]*model.Department, error) {
//	var departments []*model.Department
//
//	const query = `
//		SELECT departments.id, departments.title
//		FROM departments
//		WHERE departments.id != $1
//  		AND departments.deleted = false
//		ORDER BY title;`
//
//	rows, err := r.DB.Query(ctx, query, id)
//	if err != nil {
//		return nil, err
//	}
//
//	for rows.Next() {
//		department := new(model.Department)
//		if err := rows.Scan(
//			&department.ID,
//			&department.Title,
//		); err != nil {
//			return nil, err
//		}
//		departments = append(departments, department)
//	}
//
//	return departments, nil
//}

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
