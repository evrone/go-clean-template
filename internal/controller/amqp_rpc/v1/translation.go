package v1

import (
	"context"
	"fmt"

	"github.com/evrone/go-clean-template/internal/controller/amqp_rpc/v1/request"
	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/pkg/rabbitmq/rmq_rpc/server"
	"github.com/goccy/go-json"
	amqp "github.com/rabbitmq/amqp091-go"
)

func (r *V1) getHistory() server.CallHandler {
	return func(d *amqp.Delivery) (any, error) {
		userID, _, err := extractUserID(d, r.j)
		if err != nil {
			return nil, fmt.Errorf("amqp_rpc - V1 - getHistory - auth: %w", err)
		}

		translationHistory, err := r.t.History(context.Background(), userID)
		if err != nil {
			r.l.Error(err, "amqp_rpc - V1 - getHistory")

			return nil, fmt.Errorf("amqp_rpc - V1 - getHistory: %w", err)
		}

		return translationHistory, nil
	}
}

func (r *V1) translate() server.CallHandler {
	return func(d *amqp.Delivery) (any, error) {
		userID, data, err := extractUserID(d, r.j)
		if err != nil {
			return nil, fmt.Errorf("amqp_rpc - V1 - translate - auth: %w", err)
		}

		var req request.Translate

		err = json.Unmarshal(data, &req)
		if err != nil {
			return nil, fmt.Errorf("amqp_rpc - V1 - translate - json.Unmarshal: %w", err)
		}

		if err = r.v.Struct(req); err != nil {
			return nil, fmt.Errorf("amqp_rpc - V1 - translate - validation: %w", err)
		}

		translation, err := r.t.Translate(context.Background(), userID, entity.Translation{
			Source:      req.Source,
			Destination: req.Destination,
			Original:    req.Original,
		})
		if err != nil {
			r.l.Error(err, "amqp_rpc - V1 - translate")

			return nil, fmt.Errorf("amqp_rpc - V1 - translate: %w", err)
		}

		return translation, nil
	}
}
