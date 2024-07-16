package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
	"warehouse_backend/internal/model"
)

type AuthRepository struct {
	DB *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{DB: db}
}

// FindByPhone is to find by phone number
func (r *AuthRepository) FindByPhone(ctx context.Context, user *model.Employee) (*model.Employee, error) {
	const query = `
		SELECT id, password 
		FROM employees 
		WHERE phone = $1;`

	if err := r.DB.QueryRow(ctx, query, user.Phone).Scan(
		&user.ID,
		&user.Password,
	); err != nil {
		return nil, err
	}

	return user, nil
}

// SetHash is to set hash
func (r *AuthRepository) SetHash(ctx context.Context, id int64, hash string) error {
	const query = `
		UPDATE employees 
		SET hash = $2, authorization_date = $3
		WHERE id = $1;`

	if _, err := r.DB.Exec(ctx, query, id, hash, time.Now()); err != nil {
		return err
	}

	return nil
}

// FindByHash is to find by hash
func (r *AuthRepository) FindByHash(ctx context.Context, user *model.Employee) (*model.Employee, error) {
	const query = `
		SELECT id 
		FROM employees 
		WHERE hash = $1;`

	if err := r.DB.QueryRow(ctx, query, user.Hash).Scan(&user.ID); err != nil {
		return nil, err
	}

	return user, nil
}
