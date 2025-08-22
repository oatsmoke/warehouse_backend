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
	ctx        context.Context
	httpServer *http.Server
}

func New(ctx context.Context, port string, handlers *handler.Handler) *Server {
	return &Server{
		ctx: ctx,
		httpServer: &http.Server{
			Addr:    port,
			Handler: handlers.InitRoutes(),
		},
	}
}

// Run is server startup
func (s *Server) Run() {
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	logger.InfoInConsole(fmt.Sprintf("http server start on %s", s.httpServer.Addr))
}

// Stop is server stop
func (s *Server) Stop() {
	if err := s.httpServer.Shutdown(s.ctx); err != nil {
		log.Println(err)
	}

	logger.InfoInConsole("http server stop")
}
