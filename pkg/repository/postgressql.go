package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ConfigDB struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	DBName   string `json:"dbName"`
}

func NewPostgresDB(cfg ConfigDB) (*pgxpool.Pool, error) {
	connectStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.UserName, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	connectDB, err := pgxpool.New(context.Background(), connectStr)
	if err != nil {
		return nil, err
	}
	err = connectDB.Ping(context.Background())
	if err != nil {
		return nil, err
	}
	return connectDB, nil
}
