package main

import (
	"context"
	"log/slog"
	"os/signal"
	"syscall"
	"time"

	queries "github.com/oatsmoke/warehouse_backend/internal/db"
	"github.com/oatsmoke/warehouse_backend/internal/handler"
	"github.com/oatsmoke/warehouse_backend/internal/lib/env"
	"github.com/oatsmoke/warehouse_backend/internal/lib/logger"
	"github.com/oatsmoke/warehouse_backend/internal/lib/postgresql"
	"github.com/oatsmoke/warehouse_backend/internal/lib/redis"
	"github.com/oatsmoke/warehouse_backend/internal/lib/server"
	"github.com/oatsmoke/warehouse_backend/internal/repository"
	"github.com/oatsmoke/warehouse_backend/internal/service"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger.Init(env.GetLogLevel())

	postgresDB := postgresql.Connect(ctx, env.GetPostgresDsn())
	defer postgresDB.Close()

	redisDB := redis.Connect(env.GetRedisDsn())
	defer redis.Disconnect(redisDB)

	newQ := queries.New(postgresDB)
	newR := repository.New(postgresDB, redisDB, newQ)
	newS := service.New(newR)
	newH := handler.New(newS)

	httpS := server.New(env.GetHttpPort(), newH)
	httpS.Run()

	<-ctx.Done()
	slog.Info("received shutdown signal")

	ctxStop, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	httpS.Stop(ctxStop)
}
