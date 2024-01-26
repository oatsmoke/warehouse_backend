package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"warehouse_backend/pkg/model"
)

type CategoryRepository struct {
	db *pgxpool.Pool
}

func NewCategoryRepository(db *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) Create(title string) error {
	query := `
			INSERT INTO categories (title) 
			VALUES ($1);`
	_, err := r.db.Exec(context.Background(), query, title)
	if err != nil {
		return err
	}
	return nil
}

func (r *CategoryRepository) GetById(id int) (model.Category, error) {
	var category model.Category
	query := `
			SELECT title
			FROM categories 
			WHERE id = $1;`
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&category.Title)
	if err != nil {
		return model.Category{}, err
	}
	return category, err
}

func (r *CategoryRepository) GetAll() ([]model.Category, error) {
	var categories []model.Category
	var category model.Category
	query := `
			SELECT id, title
			FROM categories
			ORDER BY title;`
	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(
			&category.Id,
			&category.Title)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, err
}

func (r *CategoryRepository) FindByTitle(title string) (int, error) {
	var category model.Category
	query := `
			SELECT id 
			FROM categories 
			WHERE title = $1;`
	err := r.db.QueryRow(context.Background(), query, title).Scan(&category.Id)
	if err != nil {
		return 0, err
	}
	return category.Id, nil
}

func (r *CategoryRepository) Update(id int, title string) error {
	query := `
			UPDATE categories 
			SET title = $2
			WHERE id = $1;`
	_, err := r.db.Exec(context.Background(), query, id, title)
	if err != nil {
		return err
	}
	return nil
}

func (r *CategoryRepository) Delete(id int) error {
	query := `
			DELETE FROM categories 
       		WHERE id = $1;`
	_, err := r.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}
	return nil
}
