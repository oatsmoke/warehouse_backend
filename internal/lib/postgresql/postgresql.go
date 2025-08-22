package postgresql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Connect is a connection to the database
func Connect(ctx context.Context, dsn string) *pgxpool.Pool {
	connectDB, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatal(err)
	}

	err = connectDB.Ping(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return connectDB
}
