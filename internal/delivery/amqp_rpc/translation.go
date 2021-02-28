package amqprpc

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"

	"github.com/evrone/go-service-template/internal/domain"
	"github.com/evrone/go-service-template/pkg/rmq"
)

func (r *router) translationRoutes() {
	r.routerMap["getHistory"] = r.getHistory()
}

type historyResponse struct {
	History []domain.Translation `json:"history"`
}

func (r *router) getHistory() rmq.CallHandler {
	return func(d *amqp.Delivery) ([]byte, error) {
		translations, err := r.translationService.History()
		if err != nil {
			return nil, errors.Wrap(err, "amqp_rpc - router - getHistory - r.translationService.History")
		}

		response, err := json.Marshal(historyResponse{translations})
		if err != nil {
			return nil, errors.Wrap(err, "amqp_rpc - router - getHistory - json.Marshal")
		}

		return response, nil
	}
}
