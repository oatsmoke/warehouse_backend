package repository

import (
	"context"

	queries "github.com/oatsmoke/warehouse_backend/internal/db"
	"github.com/oatsmoke/warehouse_backend/internal/dto"
	"github.com/oatsmoke/warehouse_backend/internal/lib/logger"
	"github.com/oatsmoke/warehouse_backend/internal/model"
)

type DepartmentRepository struct {
	queries queries.Querier
}

func NewDepartmentRepository(queries queries.Querier) *DepartmentRepository {
	return &DepartmentRepository{
		queries: queries,
	}
}

func (r *DepartmentRepository) Create(ctx context.Context, department *model.Department) (int64, error) {
	req, err := r.queries.CreateDepartment(ctx, department.Title)
	if err != nil {
		return 0, logger.Error(logger.MsgFailedToInsert, err)
	}

	return req.ID, nil
}

func (r *DepartmentRepository) Read(ctx context.Context, id int64) (*model.Department, error) {
	req, err := r.queries.ReadDepartment(ctx, id)
	if err != nil {
		return nil, logger.Error(logger.MsgFailedToScan, err)
	}

	department := &model.Department{
		ID:        req.ID,
		Title:     req.Title,
		DeletedAt: validTime(req.DeletedAt),
	}

	return department, nil
}

func (r *DepartmentRepository) Update(ctx context.Context, department *model.Department) error {
	ct, err := r.queries.UpdateDepartment(ctx, &queries.UpdateDepartmentParams{
		ID:    department.ID,
		Title: department.Title,
	})
	if err != nil {
		return logger.Error(logger.MsgFailedToUpdate, err)
	}

	if ct.RowsAffected() == 0 {
		return logger.Error(logger.MsgFailedToUpdate, logger.ErrNoRowsAffected)
	}

	return nil
}

func (r *DepartmentRepository) Delete(ctx context.Context, id int64) error {
	ct, err := r.queries.DeleteDepartment(ctx, id)
	if err != nil {
		return logger.Error(logger.MsgFailedToDelete, err)
	}

	if ct.RowsAffected() == 0 {
		return logger.Error(logger.MsgFailedToDelete, logger.ErrNoRowsAffected)
	}

	return nil
}

func (r *DepartmentRepository) Restore(ctx context.Context, id int64) error {
	ct, err := r.queries.RestoreDepartment(ctx, id)
	if err != nil {
		return logger.Error(logger.MsgFailedToRestore, err)
	}

	if ct.RowsAffected() == 0 {
		return logger.Error(logger.MsgFailedToRestore, logger.ErrNoRowsAffected)
	}

	return nil
}

func (r *DepartmentRepository) List(ctx context.Context, qp *dto.QueryParams) ([]*model.Department, int64, error) {
	req, err := r.queries.ListDepartment(ctx, &queries.ListDepartmentParams{
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
		return []*model.Department{}, 0, nil
	}

	list := make([]*model.Department, len(req))
	for i, item := range req {
		department := &model.Department{
			ID:        item.ID,
			Title:     item.Title,
			DeletedAt: validTime(item.DeletedAt),
		}
		list[i] = department
	}

	return list, req[0].Total, nil
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
