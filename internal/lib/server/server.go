package server

import (
	"log"
	"warehouse_backend/internal/handler"
)

type Server struct {
	Handlers *handler.Handler
}

func NewServer(handlers *handler.Handler) *Server {
	return &Server{Handlers: handlers}
}

// Run is server startup
func (s *Server) Run(address string) {
	if err := s.Handlers.InitRoutes().Run(address); err != nil {
		log.Fatal(err)
	}
}
