package repository

import (
	"context"

	queries "github.com/oatsmoke/warehouse_backend/internal/db"
	"github.com/oatsmoke/warehouse_backend/internal/dto"
	"github.com/oatsmoke/warehouse_backend/internal/lib/logger"
	"github.com/oatsmoke/warehouse_backend/internal/model"
)

type CompanyRepository struct {
	queries queries.Querier
}

func NewCompanyRepository(queries queries.Querier) *CompanyRepository {
	return &CompanyRepository{
		queries: queries,
	}
}

func (r *CompanyRepository) Create(ctx context.Context, company *model.Company) (int64, error) {
	req, err := r.queries.CreateCompany(ctx, company.Title)
	if err != nil {
		return 0, logger.Error(logger.MsgFailedToInsert, err)
	}

	return req.ID, nil
}

func (r *CompanyRepository) Read(ctx context.Context, id int64) (*model.Company, error) {
	req, err := r.queries.ReadCompany(ctx, id)
	if err != nil {
		return nil, logger.Error(logger.MsgFailedToScan, err)
	}

	company := &model.Company{
		ID:        req.ID,
		Title:     req.Title,
		DeletedAt: validTime(req.DeletedAt),
	}

	return company, nil
}

func (r *CompanyRepository) Update(ctx context.Context, company *model.Company) error {
	ct, err := r.queries.UpdateCompany(ctx, &queries.UpdateCompanyParams{
		ID:    company.ID,
		Title: company.Title,
	})
	if err != nil {
		return logger.Error(logger.MsgFailedToUpdate, err)
	}

	if ct.RowsAffected() == 0 {
		return logger.Error(logger.MsgFailedToUpdate, logger.ErrNoRowsAffected)
	}

	return nil
}

func (r *CompanyRepository) Delete(ctx context.Context, id int64) error {
	ct, err := r.queries.DeleteCompany(ctx, id)
	if err != nil {
		return logger.Error(logger.MsgFailedToDelete, err)
	}

	if ct.RowsAffected() == 0 {
		return logger.Error(logger.MsgFailedToDelete, logger.ErrNoRowsAffected)
	}

	return nil
}

func (r *CompanyRepository) Restore(ctx context.Context, id int64) error {
	ct, err := r.queries.RestoreCompany(ctx, id)
	if err != nil {
		return logger.Error(logger.MsgFailedToRestore, err)
	}

	if ct.RowsAffected() == 0 {
		return logger.Error(logger.MsgFailedToRestore, logger.ErrNoRowsAffected)
	}

	return nil
}

func (r *CompanyRepository) List(ctx context.Context, qp *dto.QueryParams) ([]*model.Company, int64, error) {
	req, err := r.queries.ListCompany(ctx, &queries.ListCompanyParams{
		WithDeleted:      qp.WithDeleted,
		Search:           qp.Search,
		Ids:              qp.Ids,
		SortColumn:       qp.SortColumn,
		SortOrder:        qp.SortOrder,
		PaginationLimit:  qp.PaginationLimit,
		PaginationOffset: qp.PaginationOffset,
	})
	if err != nil {
		return nil, 0, logger.Error(logger.MsgFailedToSelect, err)
	}

	if len(req) < 1 {
		return nil, 0, nil
	}

	list := make([]*model.Company, len(req))
	for i, item := range req {
		company := &model.Company{
			ID:        item.ID,
			Title:     item.Title,
			DeletedAt: validTime(item.DeletedAt),
		}
		list[i] = company
	}

	return list, req[0].Total, nil
}
