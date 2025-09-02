// Package repository provides data access implementation for working with categories
package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oatsmoke/warehouse_backend/internal/model"
)

// CategoryRepository provides methods for managing categories in a PostgreSQL database.
type CategoryRepository struct {
	PostgresDB *pgxpool.Pool
}

// NewCategoryRepository creates a new CategoryRepository with the given database pool.
// postgresDB: PostgreSQL connection pool.
// Returns a pointer to CategoryRepository.
func NewCategoryRepository(postgresDB *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{
		PostgresDB: postgresDB,
	}
}

// Create adds a new category with the specified title.
// ctx: request context.
// title: category name.
// Returns an error if the operation fails.
func (r *CategoryRepository) Create(ctx context.Context, title string) error {
	const query = `
		INSERT INTO categories (title) 
		VALUES ($1);`

	if _, err := r.PostgresDB.Exec(ctx, query, title); err != nil {
		return err
	}

	return nil
}

// Update changes the title of an existing category by its ID.
// ctx: request context.
// id: category ID.
// title: new category name.
// Returns an error if the operation fails.
func (r *CategoryRepository) Update(ctx context.Context, id int64, title string) error {
	const query = `
		UPDATE categories 
		SET title = $2
		WHERE id = $1;`

	if _, err := r.PostgresDB.Exec(ctx, query, id, title); err != nil {
		return err
	}

	return nil
}

// Delete marks the category as deleted by its ID (soft delete).
// ctx: request context.
// id: category ID.
// Returns an error if the operation fails.
func (r *CategoryRepository) Delete(ctx context.Context, id int64) error {
	const query = `
		UPDATE categories 
		SET deleted = true
       	WHERE id = $1;`

	if _, err := r.PostgresDB.Exec(ctx, query, id); err != nil {
		return err
	}

	return nil
}

// Restore unmarks the category as deleted by its ID.
// ctx: request context.
// id: category ID.
// Returns an error if the operation fails.
func (r *CategoryRepository) Restore(ctx context.Context, id int64) error {
	const query = `
		UPDATE categories 
		SET deleted = false
       	WHERE id = $1;`

	if _, err := r.PostgresDB.Exec(ctx, query, id); err != nil {
		return err
	}

	return nil
}

// GetAll returns a list of all categories.
// ctx: request context.
// deleted: if true, includes deleted categories.
// Returns a slice of Category pointers and an error if the operation fails.
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

	rows, err := r.PostgresDB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

// GetById returns a category by its ID.
// ctx: request context.
// id: category ID.
// Returns a pointer to Category and an error if the operation fails.
func (r *CategoryRepository) GetById(ctx context.Context, id int64) (*model.Category, error) {
	category := new(model.Category)

	const query = `
		SELECT id, title, deleted
		FROM categories 
		WHERE id = $1;`

	if err := r.PostgresDB.QueryRow(ctx, query, id).Scan(
		&category.ID,
		&category.Title,
		&category.Deleted,
	); err != nil {
		return nil, err
	}

	return category, nil
}
