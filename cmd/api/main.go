package main

import (
	"context"
	"fmt"
	handlerPKG "gateway/internal/api/handler"
	serverPKG "gateway/internal/api/server"
	"gateway/internal/config"
	"gateway/internal/consul"
	"gateway/internal/logger"
	"gateway/pkg/middleware"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "gateway/docs"
)

// @title Pinek Gateway API
// @version 1.0
// @description API Gateway for Pinek services

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host 82.202.138.76:8080
// @BasePath /api/v1/

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description JWT token in the format: Bearer {token}
func main() {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("No .env file found")
	}
	envConf := config.NewEnvConfig()
	config.PrintConfigWithHiddenSecrets(envConf)

	logger.Setup(envConf.ProductionType)

	consulProvider := consul.NewProvider(envConf)
	log.Info().Msg("service registered in Consul")

	rsaPubKey := middleware.LoadRSAPublicKey()

	prometheusRegistry := prometheus.NewRegistry()
	metrics := middleware.NewMetrics(prometheusRegistry)
	log.Info().Msg("prometheus metrics initialized")

	handler := handlerPKG.NewHandler(
		envConf,
		consulProvider,
		rsaPubKey,
		metrics,
	)
	server := &serverPKG.APIServer{
		Port:    envConf.Port,
		EnvConf: envConf,
		Handler: handler,
	}

	go server.Run()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	select {
	case s := <-sig:
		log.Info().Msg(fmt.Sprintf("signal received: %s — starting graceful shutdown", s))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	server.Shutdown(ctx)
	consulProvider.DeregisterService()

	log.Info().Msg("gateway shutdown gracefully")
}
