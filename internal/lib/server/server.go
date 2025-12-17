package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/oatsmoke/warehouse_backend/internal/handler"
	"github.com/oatsmoke/warehouse_backend/internal/lib/logger"
)

type Server struct {
	httpServer *http.Server
}

func New(port string, handlers *handler.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    port,
			Handler: handlers.InitRoutes(),
		},
	}
}

func (s *Server) Run() {
	go func() {
		logger.Info(fmt.Sprintf("http server started on %s", s.httpServer.Addr))

		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()
}

func (s *Server) Stop(ctx context.Context) {
	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	logger.Info("http server stopped")
}
