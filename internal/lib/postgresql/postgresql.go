package postgresql

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"warehouse_backend/internal/lib/config"
)

// Connect is a connection to the database
func Connect(ctx context.Context, cfg *config.DB) *pgxpool.Pool {
	connectStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.UserName, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	connectDB, err := pgxpool.New(ctx, connectStr)
	if err != nil {
		log.Panic(err)
	}

	err = connectDB.Ping(ctx)
	if err != nil {
		log.Panic(err)
	}

	return connectDB
}
