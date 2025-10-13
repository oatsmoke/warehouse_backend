package repository

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/oatsmoke/warehouse_backend/internal/lib/env"
	"github.com/oatsmoke/warehouse_backend/internal/model"
	"github.com/redis/go-redis/v9"
)

type AuthRepository struct {
	RedisDB *redis.Client
}

func NewAuthRepository(redisDB *redis.Client) *AuthRepository {
	return &AuthRepository{
		RedisDB: redisDB,
	}
}

func (r *AuthRepository) Get(ctx context.Context, key string) (bool, error) {
	res, err := r.RedisDB.Get(ctx, key).Result()
	if err != nil {
		return false, err
	}

	value := new(model.AuthClaims)
	if err := json.Unmarshal([]byte(res), &value); err != nil {
		return false, err
	}

	return value.Revoked, nil
}

func (r *AuthRepository) Set(ctx context.Context, claims *jwt.RegisteredClaims, revoked bool) error {
	value := &model.AuthClaims{
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
