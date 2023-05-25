// Package app configures and runs application.
package app

import (
	"context"
	"fmt"
	openapi "github.com/evrone/go-clean-template/internal/interfaces/rest/v1/go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"github.com/evrone/go-clean-template/config"
	amqprpc "github.com/evrone/go-clean-template/internal/controller/amqp_rpc"
	"github.com/evrone/go-clean-template/internal/usecase"
	"github.com/evrone/go-clean-template/pkg/httpserver"
	"github.com/evrone/go-clean-template/pkg/logger"
	"github.com/evrone/go-clean-template/pkg/postgres"
	"github.com/evrone/go-clean-template/pkg/rabbitmq/rmq_rpc/server"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	log := logger.New(cfg.Log.Level)

	rmqServer, httpEngine := setupHttpEngine(cfg, log)
	httpServer := httpserver.New(httpEngine, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	var err error
	select {
	case s := <-interrupt:
		log.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		log.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	case err = <-rmqServer.Notify():
		log.Error(fmt.Errorf("app - Run - rmqServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		log.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}

	err = rmqServer.Shutdown()
	if err != nil {
		log.Error(fmt.Errorf("app - Run - rmqServer.Shutdown: %w", err))
	}
}

func setupHttpEngine(cfg *config.Config, log *logger.Logger) (*server.Server, *gin.Engine) {
	translationUseCase, rmqRouter := setupRabbitMqRouter()

	var ctx = context.Background()
	context.WithValue(ctx, "translationUseCase", translationUseCase)

	rmqServer, err := server.New(cfg.RMQ.URL, cfg.RMQ.ServerExchange, rmqRouter, log)
	if err != nil {
		panic(fmt.Errorf("app - Run - rmqServer - server.New: %w", err))
	}

	// HTTP Server
	router := openapi.NewRouter()
	setupRouter(router)

	return rmqServer, router
}

func setupRouter(handler *gin.Engine) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// K8s probe
	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Prometheus metrics
	handler.GET("/metrics", gin.WrapH(promhttp.Handler()))
}

func setupPostgresClient(cfg *config.Config) *postgres.Postgres {
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		panic(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	return pg
}

func setupRabbitMqRouter() (*usecase.TranslationUseCase, map[string]server.CallHandler) {
	// Use case
	translationUseCase := usecase.New()

	// RabbitMQ RPC Server
	rmqRouter := amqprpc.NewRouter(translationUseCase)
	return translationUseCase, rmqRouter
}
