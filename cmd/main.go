package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"warehouse_backend/internal/handler"
	"warehouse_backend/internal/lib/env"
	"warehouse_backend/internal/lib/logger"
	"warehouse_backend/internal/lib/postgresql"
	"warehouse_backend/internal/lib/server"
	"warehouse_backend/internal/repository"
	"warehouse_backend/internal/service"
)

func main() {
	ctx := context.Background()
	logger.Init(env.GetLogLevel())
	dbPostgres := postgresql.Connect(ctx, env.GetPostgresDsn())
	defer dbPostgres.Close()
	newR := repository.NewRepository(dbPostgres)
	newS := service.NewService(newR)
	newH := handler.NewHandler(newS)
	httpS := server.NewServer(newH)
	port := fmt.Sprintf(":%s", env.GetHttpPort())
	go httpS.Run(port)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
}
