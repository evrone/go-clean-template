package publisher

import (
	"context"
	"fmt"

	"github.com/evrone/go-service-template/domain"
)

type rmqPublisher struct {
	connect string
}

func NewRmqPublisher(connect string) domain.EntityPublisher {
	return &rmqPublisher{connect}
}

func (p *rmqPublisher) Publish(ctx context.Context, entity domain.Entity) error {
	fmt.Println(entity.Msg)

	return nil
}
