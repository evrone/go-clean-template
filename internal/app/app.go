// Package app configures and runs application.
package app

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	amqprpc "github.com/evrone/go-service-template/internal/delivery/amqp_rpc"
	v1 "github.com/evrone/go-service-template/internal/delivery/http/v1"
	v2 "github.com/evrone/go-service-template/internal/delivery/http/v2"
	"github.com/evrone/go-service-template/internal/repository"
	"github.com/evrone/go-service-template/internal/service"
	"github.com/evrone/go-service-template/internal/webapi"
	"github.com/evrone/go-service-template/pkg/httpserver"
	"github.com/evrone/go-service-template/pkg/logger"
	"github.com/evrone/go-service-template/pkg/postgres"
	"github.com/evrone/go-service-template/pkg/rmq"
)

// @title       Go Service Template API
// @version     1.0
// @description Using a translation service as an example

// @host        localhost:8080
// @BasePath    /api/v1/

// Run like main, runs application.
func Run() { //nolint:funlen // it's ok
	conf := NewConfig()

	// Logger
	zap := logger.NewZapLogger(conf.ZapLogLevel)
	defer zap.Close()

	rollbar := logger.NewRollbarLogger(conf.RollbarAccessToken, conf.RollbarEnvironment)
	defer rollbar.Close()

	logger.NewAppLogger(zap, rollbar, conf.ServiceName, conf.ServiceVersion)

	// Repository
	postgresDB := postgres.NewPostgres(conf.PgURL, conf.PgPoolMax, conf.PgConnAttempts)
	defer postgresDB.Close()

	translationRepository := repository.NewTranslationRepository(postgresDB)

	// WebAPI
	translationWebAPI := webapi.NewTranslationWebAPI()

	// Service
	translationService := service.NewTranslationService(translationRepository, translationWebAPI)

	// RabbitMQ Client
	rmqClient := rmq.NewClient("rpc_client", "rpc_server")

	// RabbitMQ Server
	rmqRouter := amqprpc.NewRouter(translationService)
	rmqServer := rmq.NewServer(rmqRouter, "rpc_server")

	//nolint:gocritic // example
	// Example RabbitMQ - RemoteCall
	//go func() {
	//	type historyResponse struct {
	//		History []domain.Translation `json:"history"`
	//	}
	//
	//	for i := 0; i < 100; i++ {
	//		var history historyResponse
	//
	//		err := rmqClient.RemoteCall("getHistory", nil, &history)
	//		if err != nil {
	//			log.Println("Error!", err)
	//		}
	//	}
	//}()

	// REST
	handler := gin.New()
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())
	handler.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler)) // Swagger
	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })  // K8s probe

	v1.NewRouter(handler, translationService)
	v2.NewRouter(handler)

	httpServer := httpserver.NewServer(handler, conf.HTTPAPIPort)

	// Graceful shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		logger.Info("app - Run - signal: " + s.String())
	case err := <-httpServer.Notify():
		logger.Error(err, "app - Run - httpServer.Notify")
	case err := <-rmqClient.Notify():
		logger.Error(err, "app - Run - rmqClient.Notify")
	case err := <-rmqServer.Notify():
		logger.Error(err, "app - Run - rmqServer.Notify")
	}

	err := httpServer.Shutdown()
	if err != nil {
		logger.Error(err, "app - Run - httpServer.Shutdown")
	}

	err = rmqClient.Shutdown()
	if err != nil {
		logger.Error(err, "app - Run - rmqClient.Shutdown")
	}

	err = rmqServer.Shutdown()
	if err != nil {
		logger.Error(err, "app - Run - rmqServer.Shutdown")
	}
}
