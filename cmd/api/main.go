package main

import (
	handlerPKG "gateway/internal/api/handler"
	serverPKG "gateway/internal/api/server"
	"gateway/internal/config"
	"gateway/internal/consul"
	"gateway/internal/logger"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found")
	}
	envConf := config.NewEnvConfig()
	config.PrintConfigWithHiddenSecrets(envConf)

	logger.Setup(envConf.ProductionType)

	consulProvider := consul.NewProvider(envConf)

	handler := handlerPKG.NewHandler(envConf, consulProvider)
	server := &serverPKG.APIServer{
		Port:    envConf.Port,
		EnvConf: envConf,
		Handler: handler,
	}

	server.Run()

	//todo: graceful shutdown (cp, server)
}
