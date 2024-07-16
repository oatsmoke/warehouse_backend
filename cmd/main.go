package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"warehouse_backend/internal/handler"
	"warehouse_backend/internal/lib/config"
	"warehouse_backend/internal/lib/logger"
	"warehouse_backend/internal/lib/postgresql"
	"warehouse_backend/internal/lib/server"
	"warehouse_backend/internal/repository"
	"warehouse_backend/internal/service"
)

func main() {
	ctx := context.Background()
	initConfig := config.Init()
	if err := os.Setenv("signingKey", "12345678"); err != nil {
		log.Panic(err)
	}
	if err := os.Setenv("tokenTTL", initConfig.TokenTTL); err != nil {
		log.Panic(err)
	}
	logger.Init(initConfig.Logger)
	dbPostgres := postgresql.Connect(ctx, initConfig.DB)
	defer dbPostgres.Close()
	newR := repository.NewRepository(dbPostgres)
	newS := service.NewService(newR)
	newH := handler.NewHandler(newS)
	httpS := server.NewServer(newH)
	address := fmt.Sprintf(":%s", initConfig.Port)
	go httpS.Run(initConfig.Client, address)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
}
