package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/oatsmoke/warehouse_backend/internal/handler"
	"github.com/oatsmoke/warehouse_backend/internal/lib/env"
	"github.com/oatsmoke/warehouse_backend/internal/lib/logger"
	"github.com/oatsmoke/warehouse_backend/internal/lib/postgresql"
	"github.com/oatsmoke/warehouse_backend/internal/lib/redis"
	"github.com/oatsmoke/warehouse_backend/internal/lib/server"
	"github.com/oatsmoke/warehouse_backend/internal/repository"
	"github.com/oatsmoke/warehouse_backend/internal/service"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger.Init(env.GetLogLevel())

	postgresDsn := env.GetPostgresDsn()
	postgresDB := postgresql.Connect(ctx, postgresDsn)
	defer postgresDB.Close()
	//migration.Run(ctx, postgresDsn)

	redisDB := redis.Connect()
	defer redisDB.Disconnect()

	newR := repository.New(postgresDB, redisDB.Conn)
	newS := service.New(newR)
	newH := handler.New(newS)

	httpS := server.New(ctx, env.GetHttpPort(), newH)
	httpS.Run()
	defer httpS.Stop()

	<-ctx.Done()
}
