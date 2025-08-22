package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oatsmoke/warehouse_backend/internal/model"
)

type ProfileRepository struct {
	DB *pgxpool.Pool
}

func NewProfileRepository(db *pgxpool.Pool) *ProfileRepository {
	return &ProfileRepository{DB: db}
}

// Create is profile create
func (r *ProfileRepository) Create(ctx context.Context, title string, categoryId int64) error {
	const query = `
		INSERT INTO profiles (title, category) 
		VALUES ($1, $2);`

	if _, err := r.DB.Exec(ctx, query, title, categoryId); err != nil {
		return err
	}

	return nil
}

// Update is profile update
func (r *ProfileRepository) Update(ctx context.Context, id int64, title string, categoryId int64) error {
	const query = `
		UPDATE profiles 
		SET title = $2, category = $3
		WHERE id = $1;`

	if _, err := r.DB.Exec(ctx, query, id, title, categoryId); err != nil {
		return err
	}

	return nil
}

// Delete is profile delete
func (r *ProfileRepository) Delete(ctx context.Context, id int64) error {
	const query = `
		UPDATE profiles 
		SET deleted = true
       	WHERE id = $1;`

	if _, err := r.DB.Exec(ctx, query, id); err != nil {
		return err
	}

	return nil
}

// Restore is profile restore
func (r *ProfileRepository) Restore(ctx context.Context, id int64) error {
	const query = `
		UPDATE profiles 
		SET deleted = false
       	WHERE id = $1;`

	if _, err := r.DB.Exec(ctx, query, id); err != nil {
		return err
	}

	return nil
}

// GetAll is to get all profiles
func (r *ProfileRepository) GetAll(ctx context.Context, deleted bool) ([]*model.Profile, error) {
	var profiles []*model.Profile
	query := ""

	if deleted {
		query = `
			SELECT profiles.id, profiles.title, profiles.deleted,
			       categories.id, categories.title
			FROM profiles
			LEFT JOIN categories ON categories.id = profiles.category
			ORDER BY profiles.title;`
	} else {
		query = `
			SELECT profiles.id, profiles.title, profiles.deleted,
			       categories.id, categories.title
			FROM profiles
			LEFT JOIN categories ON categories.id = profiles.category
			WHERE profiles.deleted = false
			ORDER BY profiles.title;`
	}

	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		profile := new(model.Profile)
		profile.Category = new(model.Category)
		if err := rows.Scan(
			&profile.ID,
			&profile.Title,
			&profile.Deleted,
			&profile.Category.ID,
			&profile.Category.Title,
		); err != nil {
			return nil, err
		}
		profiles = append(profiles, profile)
	}

	return profiles, nil
}

// GetById is to get profile by id
func (r *ProfileRepository) GetById(ctx context.Context, profile *model.Profile) (*model.Profile, error) {
	const query = `
		SELECT profiles.title, profiles.deleted,
		       categories.id, categories.title
		FROM profiles
		LEFT JOIN categories ON categories.id = profiles.category
		WHERE profiles.id = $1;`

	if err := r.DB.QueryRow(ctx, query, profile.ID).Scan(
		&profile.Title,
		&profile.Deleted,
		&profile.Category.ID,
		&profile.Category.Title,
	); err != nil {
		return nil, err
	}

	return profile, nil
}

//// GetByCategory is to get profile by category id
//func (r *ProfileRepository) GetByCategory(ctx context.Context, categoryId int64) ([]*model.Profile, error) {
//	var profiles []*model.Profile
//
//	const query = `
//		SELECT id, title, deleted
//		FROM profiles
//		WHERE category = $1;`
//
//	rows, err := r.DB.Query(ctx, query, categoryId)
//	if err != nil {
//		return nil, err
//	}
//
//	for rows.Next() {
//		profile := new(model.Profile)
//		if err := rows.Scan(
//			&profile.ID,
//			&profile.Title,
//		); err != nil {
//			return nil, err
//		}
//		profiles = append(profiles, profile)
//	}
//
//	return profiles, nil
//}

//func (r *ProfileRepository) FinDByTitle(ctx context.Context, title string) (int64, error) {
//	profile := new(model.Profile)
//
//	query := `
//			SELECT id
//			FROM profiles
//			WHERE title = $1;`
//
//	if err := r.DB.QueryRow(ctx, query, title).Scan(&profile.ID); err != nil {
//		return 0, err
//	}
//
//	return profile.ID, nil
//}
