package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oatsmoke/warehouse_backend/internal/dto"
	"github.com/oatsmoke/warehouse_backend/internal/lib/logger"
	"github.com/oatsmoke/warehouse_backend/internal/model"
)

type ProfileRepository struct {
	postgresDB *pgxpool.Pool
}

func NewProfileRepository(postgresDB *pgxpool.Pool) *ProfileRepository {
	return &ProfileRepository{
		postgresDB: postgresDB,
	}
}

func (r *ProfileRepository) Create(ctx context.Context, profile *model.Profile) (int64, error) {
	const query = `
		INSERT INTO profiles (title, category) 
		VALUES ($1, $2)
		RETURNING id;`

	var id int64
	if err := r.postgresDB.QueryRow(
		ctx,
		query,
		profile.Title,
		profile.Category.ID,
	).Scan(&id); err != nil {
		return 0, logger.Error(logger.MsgFailedToInsert, err)
	}

	if id == 0 {
		return 0, logger.Error(logger.MsgFailedToInsert, logger.ErrNoRowsAffected)
	}

	return id, nil
}

func (r *ProfileRepository) Read(ctx context.Context, id int64) (*model.Profile, error) {
	const query = `
		SELECT p.id, p.title, p.deleted_at,
		       c.id, c.title
		FROM profiles p
		LEFT JOIN categories c ON c.id = p.category
		WHERE p.id = $1;`

	profile := model.NewProfile()
	if err := r.postgresDB.QueryRow(ctx, query, id).Scan(
		&profile.ID,
		&profile.Title,
		&profile.DeletedAt,
		&profile.Category.ID,
		&profile.Category.Title,
	); err != nil {
		return nil, logger.Error(logger.MsgFailedToScan, err)
	}

	return profile, nil
}

func (r *ProfileRepository) Update(ctx context.Context, profile *model.Profile) error {
	const query = `
		UPDATE profiles 
		SET title = $2, category = $3
		WHERE id = $1 AND (title != $2 OR category != $3);`

	ct, err := r.postgresDB.Exec(
		ctx,
		query,
		profile.ID,
		profile.Title,
		profile.Category.ID)
	if err != nil {
		return logger.Error(logger.MsgFailedToUpdate, err)
	}

	if ct.RowsAffected() == 0 {
		return logger.Error(logger.MsgFailedToUpdate, logger.ErrNoRowsAffected)
	}

	return nil
}

func (r *ProfileRepository) Delete(ctx context.Context, id int64) error {
	const query = `
		UPDATE profiles 
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

func (r *ProfileRepository) Restore(ctx context.Context, id int64) error {
	const query = `
		UPDATE profiles 
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

func (r *ProfileRepository) List(ctx context.Context, qp *dto.QueryParams) ([]*model.Profile, int, error) {
	const query = `
		SELECT p.id, p.title, p.deleted_at,
		       c.id, c.title, COUNT(*) OVER() AS total
		FROM profiles p
		LEFT JOIN categories c ON c.id = p.category
		WHERE ($1 = true OR p.deleted_at IS NULL)
		  AND ($2 = '' OR (p.title || ' ' || c.title) ILIKE '%' || $2 || '%')
		  AND (array_length($3::bigint[], 1) IS NULL OR p.id = ANY ($3))
		ORDER BY CASE WHEN $4 = 'id' AND $5 = 'asc' THEN p.id::text END,
		         CASE WHEN $4 = 'id' AND $5 = 'desc' THEN p.id::text END DESC,
		         CASE WHEN $4 = 'title' AND $5 = 'asc' THEN p.title END,
		         CASE WHEN $4 = 'title' AND $5 = 'desc' THEN p.title END DESC,
		 		 CASE WHEN $4 = 'c_title' AND $5 = 'asc' THEN c.title END,
		         CASE WHEN $4 = 'c_title' AND $5 = 'desc' THEN c.title END DESC
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

	var profiles []*model.Profile
	var total int
	for rows.Next() {
		profile := model.NewProfile()
		if err := rows.Scan(
			&profile.ID,
			&profile.Title,
			&profile.DeletedAt,
			&profile.Category.ID,
			&profile.Category.Title,
			&total,
		); err != nil {
			return nil, 0, logger.Error(logger.MsgFailedToScan, err)
		}
		profiles = append(profiles, profile)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, logger.Error(logger.MsgFailedToIterateOverRows, err)
	}

	return profiles, total, nil
}
