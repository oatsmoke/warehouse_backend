package server

import (
	"log"
	"warehouse_backend/internal/handler"
	"warehouse_backend/internal/lib/config"
)

type Server struct {
	Handlers *handler.Handler
}

func NewServer(handlers *handler.Handler) *Server {
	return &Server{Handlers: handlers}
}

// Run is server startup
func (s *Server) Run(client *config.Client, address string) {
	if err := s.Handlers.InitRoutes(client).Run(address); err != nil {
		log.Panic(err)
	}
}
