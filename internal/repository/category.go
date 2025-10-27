package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oatsmoke/warehouse_backend/internal/dto"
	"github.com/oatsmoke/warehouse_backend/internal/lib/list_filter"
	"github.com/oatsmoke/warehouse_backend/internal/lib/logger"
	"github.com/oatsmoke/warehouse_backend/internal/model"
)

type CategoryRepository struct {
	postgresDB *pgxpool.Pool
}

func NewCategoryRepository(postgresDB *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{
		postgresDB: postgresDB,
	}
}

func (r *CategoryRepository) Create(ctx context.Context, category *model.Category) (int64, error) {
	const query = `
		INSERT INTO categories (title) 
		VALUES ($1)
		RETURNING id;`

	var id int64
	if err := r.postgresDB.QueryRow(ctx, query, category.Title).Scan(&id); err != nil {
		return 0, logger.Error(logger.MsgFailedToInsert, err)
	}

	if id == 0 {
		return 0, logger.Error(logger.MsgFailedToInsert, logger.ErrNoRowsAffected)
	}

	return id, nil
}

func (r *CategoryRepository) Read(ctx context.Context, id int64) (*model.Category, error) {
	const query = `
		SELECT id, title, deleted_at
		FROM categories 
		WHERE id = $1;`

	category := new(model.Category)
	if err := r.postgresDB.QueryRow(ctx, query, id).Scan(
		&category.ID,
		&category.Title,
		&category.DeletedAt,
	); err != nil {
		return nil, logger.Error(logger.MsgFailedToScan, err)
	}

	return category, nil
}

func (r *CategoryRepository) Update(ctx context.Context, category *model.Category) error {
	const query = `
		UPDATE categories 
		SET title = $2
		WHERE id = $1 AND title != $2;`

	ct, err := r.postgresDB.Exec(ctx, query, category.ID, category.Title)
	if err != nil {
		return logger.Error(logger.MsgFailedToUpdate, err)
	}

	if ct.RowsAffected() == 0 {
		return logger.Error(logger.MsgFailedToUpdate, logger.ErrNoRowsAffected)
	}

	return nil
}

func (r *CategoryRepository) Delete(ctx context.Context, id int64) error {
	const query = `
		UPDATE categories 
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

func (r *CategoryRepository) Restore(ctx context.Context, id int64) error {
	const query = `
		UPDATE categories 
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

func (r *CategoryRepository) List(ctx context.Context, qp *dto.QueryParams) ([]*model.Category, int, error) {
	fields := []string{"title"}
	str, args := list_filter.BuildQuery(qp, fields, "c")

	query := `
		SELECT id, title, deleted_at, COUNT(*) OVER() AS total
		FROM categories c
		` + str

	rows, err := r.postgresDB.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, logger.Error(logger.MsgFailedToSelect, err)
	}
	defer rows.Close()

	var categories []*model.Category
	var total int
	for rows.Next() {
		category := new(model.Category)
		if err := rows.Scan(
			&category.ID,
			&category.Title,
			&category.DeletedAt,
			&total,
		); err != nil {
			return nil, 0, logger.Error(logger.MsgFailedToScan, err)
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, logger.Error(logger.MsgFailedToIterateOverRows, err)
	}

	return categories, total, nil
}
