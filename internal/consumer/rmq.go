package consumer

import (
	"context"

	"github.com/evrone/go-service-template/domain"
)

type RmqConsumer struct {
	connect string
	useCase domain.EntityUseCase
}

func NewRmqConsumer(connect string, usecase domain.EntityUseCase) *RmqConsumer {
	return &RmqConsumer{connect, usecase}
}

func (m *RmqConsumer) Start() {
	entity := &domain.Entity{Msg: "41"}
	_ = m.useCase.Do(context.Background(), *entity)
}
