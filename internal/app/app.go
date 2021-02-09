package app

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"github.com/evrone/go-service-template/internal/delivery/http/v1"
	"github.com/evrone/go-service-template/internal/delivery/http/v2"

	"github.com/evrone/go-service-template/internal/domain"

	"github.com/evrone/go-service-template/pkg/logger"

	"github.com/evrone/go-service-template/internal/service"
	"github.com/evrone/go-service-template/internal/subservice/repository"
	"github.com/evrone/go-service-template/internal/subservice/webapi"
	"github.com/evrone/go-service-template/pkg/postgres"

	"github.com/evrone/go-service-template/pkg/httpserver"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title       Go Service Template API
// @version     1.0
// @description Using a translation service as an example

// @host        localhost:8080
// @BasePath    /api/v1/
func Run() {
	conf := NewConfig()

	// Logger
	zap := logger.NewZapLogger(conf.ZapLogLevel)
	defer zap.Close()
	rollbar := logger.NewRollbarLogger(conf.RollbarAccessToken, conf.RollbarEnvironment)
	defer rollbar.Close()
	domain.Logger = logger.NewAppLogger(zap, rollbar)

	// Repository
	pg := postgres.NewPostgres(conf.PgURL, conf.PgPoolMax, conf.PgConnAttempts)
	pgRepository := repository.NewPostgresEntityRepository(pg)
	defer pg.Close()

	// WebAPI
	googleTranslateAPI := webapi.NewGoogleTranslator()

	// Service
	entityUseCase := service.NewUseCase(pgRepository, googleTranslateAPI)

	// Router
	handler := gin.Default()
	handler.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler)) // Swagger
	handler.GET("/health", func(c *gin.Context) { c.Status(http.StatusOK) })   // K8s probe

	v1.NewAPIRouter(handler, entityUseCase)
	v2.NewAPIRouter(handler)

	server := httpserver.NewServer(handler, conf.HttpApiPort)
	server.Start()

	// Graceful shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		domain.Logger.Info("main - signal: " + s.String())
	case err := <-server.Notify():
		domain.Logger.Error(err, "main - server.Notify")
	}

	err := server.Stop()
	if err != nil {
		domain.Logger.Error(err, "main - server.Stop")
	}
}
