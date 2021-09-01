// Package app configures and runs application.
package app

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"github.com/evrone/go-clean-template/config"
	amqprpc "github.com/evrone/go-clean-template/internal/controller/amqp_rpc"
	v1 "github.com/evrone/go-clean-template/internal/controller/http/v1"
	"github.com/evrone/go-clean-template/internal/usecase"
	"github.com/evrone/go-clean-template/internal/usecase/repo"
	"github.com/evrone/go-clean-template/internal/usecase/webapi"
	"github.com/evrone/go-clean-template/pkg/httpserver"
	"github.com/evrone/go-clean-template/pkg/logger"
	"github.com/evrone/go-clean-template/pkg/postgres"
	"github.com/evrone/go-clean-template/pkg/rabbitmq/rmq_rpc/server"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	logger := logger.New(cfg.Log.Level)

	// Repository
	pg, err := postgres.NewPostgres(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		logger.Fatal(err, "app - Run - postgres.NewPostgres")
	}
	defer pg.Close()

	// Use case
	translationUseCase := usecase.New(
		repo.NewTranslationRepo(pg),
		webapi.NewTranslationWebAPI(),
	)

	// RabbitMQ RPC Server
	rmqRouter := amqprpc.NewRouter(translationUseCase)

	rmqServer, err := server.NewServer(cfg.RMQ.URL, cfg.RMQ.ServerExchange, rmqRouter, logger)
	if err != nil {
		logger.Fatal(err, "app - Run - rmqServer - server.NewServer")
	}

	// HTTP Server
	handler := gin.New()
	v1.NewRouter(handler, logger, translationUseCase)
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
	err = httpServer.Shutdown()
	if err != nil {
		logger.Error(err, "app - Run - httpServer.Shutdown")
	}

	err = rmqServer.Shutdown()
	if err != nil {
		logger.Error(err, "app - Run - rmqServer.Shutdown")
	}
}
