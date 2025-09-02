package repository

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oatsmoke/warehouse_backend/internal/lib/env"
	"github.com/oatsmoke/warehouse_backend/internal/model"
	"github.com/redis/go-redis/v9"
)

type claim struct {
	RegisteredClaims *jwt.RegisteredClaims `json:"registered_claims"`
	Revoked          bool                  `json:"revoked"`
}

type AuthRepository struct {
	PostgresDB *pgxpool.Pool
	RedisDB    *redis.Client
}

func NewAuthRepository(postgresDB *pgxpool.Pool, redisDB *redis.Client) *AuthRepository {
	return &AuthRepository{
		PostgresDB: postgresDB,
		RedisDB:    redisDB,
	}
}

// FindByPhone is to find by phone number
func (r *AuthRepository) FindByPhone(ctx context.Context, user *model.Employee) (*model.Employee, error) {
	const query = `
		SELECT id, password 
		FROM employees 
		WHERE phone = $1;`

	if err := r.PostgresDB.QueryRow(ctx, query, user.Phone).Scan(
		&user.ID,
		&user.Password,
	); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *AuthRepository) Set(ctx context.Context, claims *jwt.RegisteredClaims, revoked bool) error {
	value := &claim{
		RegisteredClaims: claims,
		Revoked:          revoked,
	}

	marshalClaims, err := json.Marshal(value)
	if err != nil {
		return err
	}

	refreshTTL, err := strconv.Atoi(env.GetRefreshTtl())
	if err != nil {
		return err
	}

	return r.RedisDB.Set(ctx, claims.ID, marshalClaims, time.Duration(refreshTTL)*time.Second).Err()
}

func (r *AuthRepository) Get(ctx context.Context, key string) (bool, error) {
	res, err := r.RedisDB.Get(ctx, key).Result()
	if err != nil {
		return false, err
	}

	value := new(claim)
	if err := json.Unmarshal([]byte(res), &value); err != nil {
		return false, err
	}

	return value.Revoked, nil
}

// SetHash is to set hash
//func (r *AuthRepository) SetHash(ctx context.Context, id int64, hash string) error {
//	const query = `
//		UPDATE employees
//		SET hash = $2, authorization_date = $3
//		WHERE id = $1;`
//
//	if _, err := r.DB.Exec(ctx, query, id, hash, time.Now()); err != nil {
//		return err
//	}
//
//	return nil
//}

// FindByHash is to find by hash
//func (r *AuthRepository) FindByHash(ctx context.Context, user *model.Employee) (*model.Employee, error) {
//	const query = `
//		SELECT id
//		FROM employees
//		WHERE hash = $1;`
//
//	if err := r.DB.QueryRow(ctx, query, user.Hash).Scan(&user.ID); err != nil {
//		return nil, err
//	}
//
//	return user, nil
//}
