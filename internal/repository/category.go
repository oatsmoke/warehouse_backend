package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oatsmoke/warehouse_backend/internal/model"
)

type CategoryRepository struct {
	DB *pgxpool.Pool
}

func NewCategoryRepository(db *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{DB: db}
}

// Create is category create
func (r *CategoryRepository) Create(ctx context.Context, title string) error {
	const query = `
		INSERT INTO categories (title) 
		VALUES ($1);`

	if _, err := r.DB.Exec(ctx, query, title); err != nil {
		return err
	}

	return nil
}

// Update is category update
func (r *CategoryRepository) Update(ctx context.Context, id int64, title string) error {
	const query = `
		UPDATE categories 
		SET title = $2
		WHERE id = $1;`

	if _, err := r.DB.Exec(ctx, query, id, title); err != nil {
		return err
	}

	return nil
}

// Delete is category delete
func (r *CategoryRepository) Delete(ctx context.Context, id int64) error {
	const query = `
		UPDATE categories 
		SET deleted = true
       	WHERE id = $1;`

	if _, err := r.DB.Exec(ctx, query, id); err != nil {
		return err
	}

	return nil
}

// Restore is category restore
func (r *CategoryRepository) Restore(ctx context.Context, id int64) error {
	const query = `
		UPDATE categories 
		SET deleted = false
       	WHERE id = $1;`

	if _, err := r.DB.Exec(ctx, query, id); err != nil {
		return err
	}

	return nil
}

// GetAll is to get all categories
func (r *CategoryRepository) GetAll(ctx context.Context, deleted bool) ([]*model.Category, error) {
	var categories []*model.Category
	query := ""

	if deleted {
		query = `
			SELECT id, title, deleted
			FROM categories
			ORDER BY title;`
	} else {
		query = `
			SELECT id, title, deleted
			FROM categories
			WHERE deleted = false
			ORDER BY title;`
	}

	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		category := new(model.Category)
		if err := rows.Scan(
			&category.ID,
			&category.Title,
			&category.Deleted,
		); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

// GetById is to get category by id
func (r *CategoryRepository) GetById(ctx context.Context, category *model.Category) (*model.Category, error) {
	const query = `
		SELECT title, deleted
		FROM categories 
		WHERE id = $1;`

	if err := r.DB.QueryRow(ctx, query, category.ID).Scan(
		&category.Title,
		&category.Deleted,
	); err != nil {
		return nil, err
	}

	return category, nil
}

//func (r *CategoryRepository) FindByTitle(ctx context.Context, title string) (int64, error) {
//	category := new(model.Category)
//
//	query := `
//			SELECT id
//			FROM categories
//			WHERE title = $1;`
//
//	if err := r.DB.QueryRow(ctx, query, title).Scan(&category.ID); err != nil {
//		return 0, err
//	}
//
//	return category.ID, nil
//}
