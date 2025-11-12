package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oatsmoke/warehouse_backend/internal/dto"
	"github.com/oatsmoke/warehouse_backend/internal/lib/logger"
	"github.com/oatsmoke/warehouse_backend/internal/model"
)

type CompanyRepository struct {
	postgresDB *pgxpool.Pool
}

func NewCompanyRepository(postgresDB *pgxpool.Pool) *CompanyRepository {
	return &CompanyRepository{
		postgresDB: postgresDB,
	}
}

func (r *CompanyRepository) Create(ctx context.Context, company *model.Company) (int64, error) {
	const query = `
		INSERT INTO companies (title) 
		VALUES ($1)
		RETURNING id;`

	var id int64
	if err := r.postgresDB.QueryRow(ctx, query, company.Title).Scan(&id); err != nil {
		return 0, logger.Error(logger.MsgFailedToInsert, err)
	}

	if id == 0 {
		return 0, logger.Error(logger.MsgFailedToMarshal, logger.ErrNoRowsAffected)
	}

	return id, nil
}

func (r *CompanyRepository) Read(ctx context.Context, id int64) (*model.Company, error) {
	const query = `
		SELECT id, title, deleted_at
		FROM companies 
		WHERE id = $1;`

	company := new(model.Company)
	if err := r.postgresDB.QueryRow(ctx, query, id).Scan(
		&company.ID,
		&company.Title,
		&company.DeletedAt,
	); err != nil {
		return nil, logger.Error(logger.MsgFailedToScan, err)
	}

	return company, nil
}

func (r *CompanyRepository) Update(ctx context.Context, company *model.Company) error {
	const query = `
		UPDATE companies 
		SET title = $2
		WHERE id = $1 AND title != $2;`

	ct, err := r.postgresDB.Exec(ctx, query, company.ID, company.Title)
	if err != nil {
		return err
	}

	if ct.RowsAffected() == 0 {
		return logger.Error(logger.MsgFailedToUpdate, logger.ErrNoRowsAffected)
	}

	return nil
}

func (r *CompanyRepository) Delete(ctx context.Context, id int64) error {
	const query = `
		UPDATE companies 
		SET deleted_at = now()
       	WHERE id = $1 AND deleted_at IS NULL;`

	ct, err := r.postgresDB.Exec(ctx, query, id)
	if err != nil {
		return logger.Error(logger.MsgFailedToDelete, err)
	}

	if ct.RowsAffected() == 0 {
		return logger.Error(logger.MsgFailedToDelete, logger.ErrNoRowsAffected)
	}

	return nil
}

func (r *CompanyRepository) Restore(ctx context.Context, id int64) error {
	const query = `
		UPDATE companies 
		SET deleted_at = NULL
       	WHERE id = $1 AND deleted_at IS NOT NULL;`

	ct, err := r.postgresDB.Exec(ctx, query, id)
	if err != nil {
		return logger.Error(logger.MsgFailedToRestore, err)
	}

	if ct.RowsAffected() == 0 {
		return logger.Error(logger.MsgFailedToRestore, logger.ErrNoRowsAffected)
	}

	return nil
}

func (r *CompanyRepository) List(ctx context.Context, qp *dto.QueryParams) ([]*model.Company, int, error) {
	const query = `
		SELECT id, title, deleted_at, count(*) OVER () AS total
		FROM companies
		WHERE ($1 = true OR deleted_at IS NULL)
		  AND ($2 = '' OR title ILIKE '%' || $2 || '%')
		  AND (array_length($3::bigint[], 1) IS NULL OR id = ANY ($3))
		ORDER BY CASE WHEN $4 = 'id' AND $5 = 'asc' THEN id::text END,
		         CASE WHEN $4 = 'id' AND $5 = 'desc' THEN id::text END DESC,
		         CASE WHEN $4 = 'title' AND $5 = 'asc' THEN title END,
		         CASE WHEN $4 = 'title' AND $5 = 'desc' THEN title END DESC
		LIMIT $6 OFFSET $7;`

	rows, err := r.postgresDB.Query(
		ctx,
		query,
		qp.WithDeleted,
		qp.Search,
		qp.Ids,
		qp.SortBy,
		qp.Order,
		qp.Limit,
		qp.Offset,
	)
	if err != nil {
		return nil, 0, logger.Error(logger.MsgFailedToSelect, err)
	}
	defer rows.Close()

	var companies []*model.Company
	var total int
	for rows.Next() {
		company := new(model.Company)
		if err := rows.Scan(
			&company.ID,
			&company.Title,
			&company.DeletedAt,
			&total,
		); err != nil {
			return nil, 0, logger.Error(logger.MsgFailedToScan, err)
		}
		companies = append(companies, company)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, logger.Error(logger.MsgFailedToIterateOverRows, err)
	}

	return companies, total, nil
}
