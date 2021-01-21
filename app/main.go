package main

import (
	"github.com/evrone/go-service-template/entity"
	"github.com/evrone/go-service-template/entity/publisher"
	"github.com/evrone/go-service-template/entity/repository"
	"github.com/evrone/go-service-template/internal/consumer"
)

func main() {
	connectDB := "OpenPostgresWithConfig"
	entityRepository := repository.NewPostgresEntityRepository(connectDB)

	connectRmq := "RabbitMQ"
	rmqPublisher := publisher.NewRmqPublisher(connectRmq)

	entityUseCase := entity.NewUseCase(entityRepository, rmqPublisher)
	rmqConsumer := consumer.NewRmqConsumer(connectRmq, entityUseCase)
	rmqConsumer.Start()
}
