package server

import (
	"fmt"
	"log"
	"warehouse_backend/internal/handler"
	"warehouse_backend/internal/lib/logger"
)

type Server struct {
	Handlers *handler.Handler
}

func NewServer(handlers *handler.Handler) *Server {
	return &Server{Handlers: handlers}
}

// Run is server startup
func (s *Server) Run(address string) {
	logger.InfoInConsole(fmt.Sprintf("http server start on %s", address), "")
	if err := s.Handlers.InitRoutes().Run(address); err != nil {
		log.Fatal(err)
	}
}
