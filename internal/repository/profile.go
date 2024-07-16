package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"warehouse_backend/internal/model"
)

type ProfileRepository struct {
	db *pgxpool.Pool
}

func NewProfileRepository(db *pgxpool.Pool) *ProfileRepository {
	return &ProfileRepository{db: db}
}

func (r *ProfileRepository) Create(ctx context.Context, title string, category int64) error {
	query := `
			INSERT INTO profiles (title, category) 
			VALUES ($1, $2);`

	if _, err := r.db.Exec(ctx, query, title, category); err != nil {
		return err
	}

	return nil
}

func (r *ProfileRepository) GetByCategory(ctx context.Context, id int64) ([]*model.Profile, error) {
	var profiles []*model.Profile
	profile := new(model.Profile)

	query := `
			SELECT id, title
			FROM profiles
			WHERE category = $1;`

	rows, err := r.db.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		if err := rows.Scan(&profile.ID, &profile.Title); err != nil {
			return nil, err
		}
		profiles = append(profiles, profile)
	}

	return profiles, nil
}

func (r *ProfileRepository) GetById(ctx context.Context, id int64) (*model.Profile, error) {
	profile := new(model.Profile)

	query := `
			SELECT title, category
			FROM profiles
			WHERE id = $1;`

	err := r.db.QueryRow(ctx, query, id).Scan(&profile)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (r *ProfileRepository) GetAll(ctx context.Context) ([]*model.Profile, error) {
	var profiles []*model.Profile
	profile := new(model.Profile)

	query := `
			SELECT profiles.id, profiles.title, 
			       categories.id, categories.title
			FROM profiles
			LEFT JOIN categories ON categories.id = profiles.category
			ORDER BY profiles.title;`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		if err := rows.Scan(
			&profile.ID,
			&profile.Title,
			&profile.Category.ID,
			&profile.Category.Title); err != nil {
			return nil, err
		}
		profiles = append(profiles, profile)
	}

	return profiles, nil
}

func (r *ProfileRepository) FindByTitle(ctx context.Context, title string) (int64, error) {
	profile := new(model.Profile)

	query := `
			SELECT id 
			FROM profiles 
			WHERE title = $1;`

	if err := r.db.QueryRow(ctx, query, title).Scan(&profile.ID); err != nil {
		return 0, err
	}

	return profile.ID, nil
}

func (r *ProfileRepository) Update(ctx context.Context, id int64, title string, category int64) error {
	query := `
			UPDATE profiles 
			SET title = $2, category = $3
			WHERE id = $1;`

	if _, err := r.db.Exec(ctx, query, id, title, category); err != nil {
		return err
	}

	return nil
}

func (r *ProfileRepository) Delete(ctx context.Context, id int64) error {
	query := `
			DELETE FROM profiles 
       		WHERE id = $1;`

	if _, err := r.db.Exec(ctx, query, id); err != nil {
		return err
	}

	return nil
}
