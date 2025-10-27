package server

import (
	"context"
	"gateway/internal/api/handler"
	"gateway/internal/config"
	"github.com/rs/zerolog/log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIServer struct {
	Port    string
	EnvConf *config.Config
	Handler *handler.Handler
	server  *http.Server
}

func (s *APIServer) Run() {
	if s.EnvConf.ProductionType == "prod" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	s.server = &http.Server{
		Handler: s.Handler.InitRoutes(),
		Addr:    ":" + s.EnvConf.Port,
	}

	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal().
			Err(err).
			Msg("failed to run API server")
	}
}

func (s *APIServer) Shutdown(ctx context.Context) {
	if err := s.server.Shutdown(ctx); err != nil {
		log.Error().
			Err(err).
			Msg("error during API server shutdown")
	}

	log.Info().Msg("API server shutdown gracefully")
}
