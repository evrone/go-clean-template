package app

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/evrone/go-service-template/internal/business-logic/domain"

	"github.com/evrone/go-service-template/pkg/logger"

	"github.com/evrone/go-service-template/internal/business-logic/entity"
	"github.com/evrone/go-service-template/internal/business-logic/entity/repository"
	"github.com/evrone/go-service-template/internal/business-logic/entity/translator"
	"github.com/evrone/go-service-template/internal/entrypoints/http/api/v1"
	"github.com/evrone/go-service-template/pkg/postgres"

	"github.com/evrone/go-service-template/pkg/httpserver"

	"github.com/evrone/go-service-template/internal/entrypoints/http/probe"
)

func Run() {
	conf := NewConfig()

	zap := logger.NewZapLogger(conf.ZapLogLevel)
	defer zap.Close()
	rollbar := logger.NewRollbarLogger(conf.RollbarAccessToken, conf.RollbarEnvironment)
	defer rollbar.Close()
	domain.Logger = logger.NewAppLogger(zap, rollbar)

	pg := postgres.NewPostgres(conf.PgURL, conf.PgPoolMax, conf.PgConnAttempts)
	pgRepository := repository.NewPostgresEntityRepository(pg)
	defer pg.Close()

	googleTranslateAPI := translator.NewGoogleTranslator()

	entityUseCase := entity.NewUseCase(pgRepository, googleTranslateAPI)

	apiRouter := v1.NewApiRouter(entityUseCase)
	apiServer := httpserver.NewServer(apiRouter, conf.HttpApiPort)
	apiServer.Start()

	probeRouter := probe.NewHealthRouter()
	probeServer := httpserver.NewServer(probeRouter, conf.HttpProbePort)
	probeServer.Start()

	// Graceful shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		domain.Logger.Info("main - signal: " + s.String())
	case err := <-probeServer.Notify():
		domain.Logger.Error(err, "main - probeServer.Notify")
	case err := <-apiServer.Notify():
		domain.Logger.Error(err, "main - apiServer.Notify")
	}

	err := probeServer.Stop()
	if err != nil {
		domain.Logger.Error(err, "main - probeServer.Stop")
	}

	err = apiServer.Stop()
	if err != nil {
		domain.Logger.Error(err, "main - apiServer.Stop")
	}
}
