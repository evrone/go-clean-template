package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/evrone/go-service-template/internal/router"
	"github.com/evrone/go-service-template/pkg/server"

	"github.com/evrone/go-service-template/entity"
	"github.com/evrone/go-service-template/entity/publisher"
	"github.com/evrone/go-service-template/entity/repository"
	"github.com/evrone/go-service-template/internal/consumer"
	"github.com/evrone/go-service-template/pkg/postgres"
)

func main() {
	time.Sleep(time.Second * 3)

	conf := NewConfig()
	db := postgres.NewPostgres(conf.PgURL, conf.PgPoolMax)
	entityRepository := repository.NewPostgresEntityRepository(db, conf.PgTableName)

	connectRmq := "RabbitMQ"
	rmqPublisher := publisher.NewRmqPublisher(connectRmq)

	entityUseCase := entity.NewUseCase(entityRepository, rmqPublisher)
	rmqConsumer := consumer.NewRmqConsumer(connectRmq, entityUseCase)
	rmqConsumer.Start()

	probeRouter := router.NewProbeRouter()
	probeServer := server.NewServer(probeRouter, conf.AppProbePort)
	probeServer.Start()

	// Graceful shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Println("main - Interrupt signal", s.String()) // TODO
	case err := <-probeServer.Notify():
		log.Println("main - probeServer.Notify error", err.Error()) // TODO
	}

	err := probeServer.Stop()
	if err != nil {
		log.Println("main - probeServer.Stop error") // TODO
	}
}
