package main

import (
	"gexabyte/internal/config"
	"gexabyte/internal/repository"
	"gexabyte/internal/service"
	"gexabyte/internal/transport/http"
	"gexabyte/pkg/logger"
	"log"
	"time"

	binance_connector "github.com/binance/binance-connector-go"
)

func init() {
	time.Local = time.UTC
}

//	@title			Gexabyte
//	@version		1.0
//	@description	Gexabyte test assignment

// @BasePath	/api/v1
func main() {
	cfg := config.MustLoad()
	logger := logger.New(cfg.LogLevel)

	repo, err := repository.NewRepository(cfg)
	if err != nil {
		panic(err)
	}

	binanceClient := binance_connector.NewClient(cfg.Binance.ApiKey, cfg.Binance.SecretKey)

	service := service.New(cfg, logger, binanceClient, repo)

	server := http.New(cfg, logger, service)

	log.Fatalln(server.Start())
}
