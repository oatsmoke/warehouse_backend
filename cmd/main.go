package main

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"warehouse_backend/pkg/handler"
	"warehouse_backend/pkg/repository"
	"warehouse_backend/pkg/service"
)

type Config struct {
	Port   string               `json:"port"`
	DB     repository.ConfigDB  `json:"db"`
	Client handler.ConfigClient `json:"client"`
}

type Server struct {
	httpServer *http.Server
}

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	config, err := initConfig()
	if err != nil {
		logrus.Fatalf("Ошибка инициальзации конфига:%v", err)
	}
	db, err := repository.NewPostgresDB(config.DB)
	if err != nil {
		logrus.Fatalf("Ошибка инициальзации базы данных:%v", err)
	}
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)
	srv := new(Server)
	if err := srv.Run(config.Port, handlers.InitRoutes(config.Client)); err != nil {
		logrus.Fatalf("Ошибка запуска сервера:%v", err)
	}
}

func (s *Server) Run(port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}
	return s.httpServer.ListenAndServe()
}

func initConfig() (Config, error) {
	configData, err := os.ReadFile("config.json")
	if err != nil {
		return Config{}, err
	}
	config := Config{}
	err = json.Unmarshal(configData, &config)
	if err != nil {
		return Config{}, err
	}
	return config, err
}
