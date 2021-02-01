package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/evrone/go-service-template/business-logic/entity"
	"github.com/evrone/go-service-template/business-logic/entity/repository"
	"github.com/evrone/go-service-template/business-logic/entity/translator"
	"github.com/evrone/go-service-template/entrypoints/http/api/v1"
	"github.com/evrone/go-service-template/infrastructure/postgres"

	"github.com/evrone/go-service-template/infrastructure/httpserver"

	"github.com/evrone/go-service-template/entrypoints/http/probe"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	conf := NewConfig()
	pg := postgres.NewPostgres(conf.PgURL, conf.PgPoolMax, conf.PgConnAttempts)
	pgRepository := repository.NewPostgresEntityRepository(pg)

	translateApi := translator.NewGoogleTranslator()

	entityUseCase := entity.NewUseCase(pgRepository, translateApi)

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
		log.Println("main - Interrupt signal", s.String()) // TODO
	case err := <-probeServer.Notify():
		log.Println("main - probeServer.Notify error", err.Error()) // TODO
	case err := <-apiServer.Notify():
		log.Println("main - apiServer.Notify error", err.Error()) // TODO
	}

	err := probeServer.Stop()
	if err != nil {
		log.Println("main - probeServer.Stop error") // TODO
	}

	err = apiServer.Stop()
	if err != nil {
		log.Println("main - apiServer.Stop error") // TODO
	}

	pg.Close()
}
