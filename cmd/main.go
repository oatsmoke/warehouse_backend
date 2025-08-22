package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/oatsmoke/warehouse_backend/internal/handler"
	"github.com/oatsmoke/warehouse_backend/internal/lib/env"
	"github.com/oatsmoke/warehouse_backend/internal/lib/logger"
	"github.com/oatsmoke/warehouse_backend/internal/lib/postgresql"
	"github.com/oatsmoke/warehouse_backend/internal/lib/server"
	"github.com/oatsmoke/warehouse_backend/internal/repository"
	"github.com/oatsmoke/warehouse_backend/internal/service"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger.Init(env.GetLogLevel())

	dbPostgres := postgresql.Connect(ctx, env.GetPostgresDsn())
	defer dbPostgres.Close()

	newR := repository.New(dbPostgres)
	newS := service.New(newR)
	newH := handler.New(newS)

	httpS := server.New(ctx, env.GetHttpPort(), newH)
	httpS.Run()
	defer httpS.Stop()

	<-ctx.Done()
}
