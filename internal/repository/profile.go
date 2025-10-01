package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
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
	if err := r.postgresDB.QueryRow(ctx, query, profile.Title, profile.Category.ID).Scan(&id); err != nil {
		return 0, err
	}

	if id == 0 {
		return 0, logger.NoRowsAffected
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

	profile := newProfile()
	if err := r.postgresDB.QueryRow(ctx, query, id).Scan(
		&profile.ID,
		&profile.Title,
		&profile.DeletedAt,
		&profile.Category.ID,
		&profile.Category.Title,
	); err != nil {
		return nil, err
	}

	return profile, nil
}

func (r *ProfileRepository) Update(ctx context.Context, profile *model.Profile) error {
	const query = `
		UPDATE profiles 
		SET title = $2, category = $3
		WHERE id = $1;`

	ct, err := r.postgresDB.Exec(ctx, query, profile.ID, profile.Title, profile.Category.ID)
	if err != nil {
		return err
	}

	if ct.RowsAffected() == 0 {
		return logger.NoRowsAffected
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
		return err
	}

	if ct.RowsAffected() == 0 {
		return logger.NoRowsAffected
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
		return err
	}

	if ct.RowsAffected() == 0 {
		return logger.NoRowsAffected
	}

	return nil
}

func (r *ProfileRepository) List(ctx context.Context, withDeleted bool) ([]*model.Profile, error) {
	const query = `
		SELECT p.id, p.title, p.deleted_at,
		       c.id, c.title
		FROM profiles p
		LEFT JOIN categories c ON c.id = p.category
		WHERE $1 OR p.deleted_at IS NULL
		ORDER BY p.title;`

	rows, err := r.postgresDB.Query(ctx, query, withDeleted)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []*model.Profile
	for rows.Next() {
		profile := newProfile()
		if err := rows.Scan(
			&profile.ID,
			&profile.Title,
			&profile.DeletedAt,
			&profile.Category.ID,
			&profile.Category.Title,
		); err != nil {
			return nil, err
		}
		profiles = append(profiles, profile)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return profiles, nil
}

func newProfile() *model.Profile {
	return &model.Profile{
		Category: &model.Category{},
	}
}
