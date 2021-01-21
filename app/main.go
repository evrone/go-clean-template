package main

import (
	"github.com/evrone/go-service-template/entity"
	"github.com/evrone/go-service-template/entity/repository/postgres"
	"github.com/evrone/go-service-template/internal/ebus/rmq"
)

func main() {
	connectDB := "OpenPostgresWithConfig"
	entityRepository := postgres.NewEntityRepository(connectDB)
	entityUsecase := entity.NewUsecase(entityRepository)

	ebus := rmq.NewRabbitMQ(entityUsecase)
	ebus.Start(41)
}
