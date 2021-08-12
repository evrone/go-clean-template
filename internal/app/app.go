// Package app configures and runs application.
package app

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"github.com/evrone/go-service-template/config"
	amqprpc "github.com/evrone/go-service-template/internal/delivery/amqp_rpc"
	v1 "github.com/evrone/go-service-template/internal/delivery/http/v1"
	"github.com/evrone/go-service-template/internal/service"
	"github.com/evrone/go-service-template/internal/service/repo"
	"github.com/evrone/go-service-template/internal/service/webapi"
	"github.com/evrone/go-service-template/pkg/httpserver"
	"github.com/evrone/go-service-template/pkg/logger"
	"github.com/evrone/go-service-template/pkg/postgres"
	"github.com/evrone/go-service-template/pkg/rabbitmq/rmq_rpc/server"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	// Repository
	pg, err := postgres.NewPostgres(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		logger.Fatal(err, "app - Run - postgres.NewPostgres")
	}
	defer pg.Close()

	// Service
	translationService := service.NewTranslationService(
		repo.NewTranslationRepo(pg),
		webapi.NewTranslationWebAPI(),
	)

	// RabbitMQ RPC Server
	rmqRouter := amqprpc.NewRouter(translationService)

	rmqServer, err := server.NewServer(cfg.RMQ.URL, cfg.RMQ.ServerExchange, rmqRouter)
	if err != nil {
		logger.Fatal(err, "app - Run - rmqServer - server.NewServer")
	}

	// HTTP Server
	handler := gin.New()
	v1.NewRouter(handler, translationService)
	httpServer := httpserver.NewServer(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		logger.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		logger.Error(err, "app - Run - httpServer.Notify")
	case err = <-rmqServer.Notify():
		logger.Error(err, "app - Run - rmqServer.Notify")
	}

	// Shutdown
	if err := httpServer.Shutdown(); err != nil {
		logger.Error(err, "app - Run - httpServer.Shutdown")
	}

	if err := rmqServer.Shutdown(); err != nil {
		logger.Error(err, "app - Run - rmqServer.Shutdown")
	}
}
