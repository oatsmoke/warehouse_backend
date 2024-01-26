package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"warehouse_backend/pkg/model"
)

type ProfileRepository struct {
	db *pgxpool.Pool
}

func NewProfileRepository(db *pgxpool.Pool) *ProfileRepository {
	return &ProfileRepository{db: db}
}

func (r *ProfileRepository) Create(title string, category int) error {
	query := `
			INSERT INTO profiles (title, category) 
			VALUES ($1, $2);`
	_, err := r.db.Exec(context.Background(), query, title, category)
	if err != nil {
		return err
	}
	return nil
}

func (r *ProfileRepository) GetByCategory(id int) ([]model.Profile, error) {
	var profiles []model.Profile
	var profile model.Profile
	query := `
			SELECT id, title
			FROM profiles
			WHERE category = $1;`
	rows, err := r.db.Query(context.Background(), query, id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(
			&profile.Id,
			&profile.Title)
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, profile)
	}
	return profiles, err
}

func (r *ProfileRepository) GetById(id int) (model.Profile, error) {
	var profile model.Profile
	query := `
			SELECT title, category
			FROM profiles
			WHERE id = $1;`
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&profile.Title,
		&profile.Category.Id)
	if err != nil {
		return model.Profile{}, err
	}
	return profile, err
}

func (r *ProfileRepository) GetAll() ([]model.Profile, error) {
	var profiles []model.Profile
	var profile model.Profile
	query := `
			SELECT profiles.id, profiles.title, categories.id, categories.title
			FROM profiles
			LEFT JOIN categories ON categories.id = profiles.category
			ORDER BY profiles.title;`
	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(
			&profile.Id,
			&profile.Title,
			&profile.Category.Id,
			&profile.Category.Title)
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, profile)
	}
	return profiles, err
}

func (r *ProfileRepository) FindByTitle(title string) (int, error) {
	var profile model.Profile
	query := `
			SELECT id 
			FROM profiles 
			WHERE title = $1;`
	err := r.db.QueryRow(context.Background(), query, title).Scan(&profile.Id)
	if err != nil {
		return 0, err
	}
	return profile.Id, nil
}

func (r *ProfileRepository) Update(id int, title string, category int) error {
	query := `
			UPDATE profiles 
			SET title = $2, category = $3
			WHERE id = $1;`
	_, err := r.db.Exec(context.Background(), query, id, title, category)
	if err != nil {
		return err
	}
	return nil
}

func (r *ProfileRepository) Delete(id int) error {
	query := `
			DELETE FROM profiles 
       		WHERE id = $1;`
	_, err := r.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}
	return nil
}
