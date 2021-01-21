package rmq

import (
	"context"
	"fmt"
	"github.com/evrone/go-service-template/domain"
)

type RabbitMQ struct {
	Usecase domain.EntityUsecase
}

func NewRabbitMQ(usecase domain.EntityUsecase) *RabbitMQ {
	return &RabbitMQ{usecase}
}

func (m *RabbitMQ) Start(ID int) {
	ctx := context.Background()
	msg, _ := m.Usecase.Get(ctx, ID)
	fmt.Println(msg)
}
