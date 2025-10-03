package postgresql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oatsmoke/warehouse_backend/internal/lib/env"
)

func Connect(ctx context.Context, dsn string) *pgxpool.Pool {
	connectDB, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err := connectDB.Ping(ctx); err != nil {
		log.Fatal(err)
	}

	return connectDB
}

func ConnectTest() *pgxpool.Pool {
	return Connect(context.Background(), env.GetTestPostgresDsn())
}
