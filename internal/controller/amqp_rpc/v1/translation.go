package v1

import (
	"context"
	"fmt"

	"github.com/evrone/go-clean-template/pkg/rabbitmq/rmq_rpc/server"
	amqp "github.com/rabbitmq/amqp091-go"
)

func (r *V1) getHistory() server.CallHandler {
	return func(_ *amqp.Delivery) (interface{}, error) {
		translationHistory, err := r.t.History(context.Background())
		if err != nil {
			r.l.Error(err, "amqp_rpc - V1 - getHistory")

			return nil, fmt.Errorf("amqp_rpc - V1 - getHistory: %w", err)
		}

		return translationHistory, nil
	}
}
