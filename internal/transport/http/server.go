package http

import (
	"fmt"
	"gexabyte/internal/config"
	"gexabyte/internal/service"
	"log/slog"
	"net/http"
	"time"
)

type Server struct {
	http.Server

	service *service.Manager
	logger  *slog.Logger
}

func New(cfg *config.Config, logger *slog.Logger, service *service.Manager) *Server {
	return &Server{
		Server: http.Server{
			Addr:           fmt.Sprintf(":%s", cfg.Server.Port),
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},

		service: service,
		logger:  logger,
	}
}

func (s *Server) Start() error {
	s.setupRouter()
	return s.ListenAndServe()
}
