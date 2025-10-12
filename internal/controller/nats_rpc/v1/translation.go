package v1

import (
	"context"
	"fmt"

	"github.com/evrone/go-clean-template/pkg/nats/nats_rpc/server"
	"github.com/nats-io/nats.go"
)

func (r *V1) getHistory() server.CallHandler {
	return func(_ *nats.Msg) (interface{}, error) {
		translationHistory, err := r.t.History(context.Background())
		if err != nil {
			r.l.Error(err, "amqp_rpc - V1 - getHistory")

			return nil, fmt.Errorf("amqp_rpc - V1 - getHistory: %w", err)
		}

		return translationHistory, nil
	}
}
