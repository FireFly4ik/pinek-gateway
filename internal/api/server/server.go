package server

import (
	"fmt"
	"gateway/internal/api/handler"
	"gateway/internal/config"
	"github.com/rs/zerolog/log"

	"github.com/gin-gonic/gin"
)

type APIServer struct {
	Port    string
	EnvConf *config.Config
	Handler *handler.Handler
}

func (s *APIServer) Run() {
	if s.EnvConf.ProductionType == "prod" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	if err := s.Handler.InitRoutes().Run(fmt.Sprintf(":%v", s.Port)); err != nil {
		log.Fatal().
			Err(err).
			Msg("failed to run API server")
	}
}
