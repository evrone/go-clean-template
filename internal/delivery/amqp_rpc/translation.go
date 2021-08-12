package amqprpc

import (
	"github.com/pkg/errors"
	"github.com/streadway/amqp"

	"github.com/evrone/go-clean-template/internal/domain"
	"github.com/evrone/go-clean-template/internal/service"
	"github.com/evrone/go-clean-template/pkg/rabbitmq/rmq_rpc/server"
)

type translationRoutes struct {
	translationService service.Translation
}

func newTranslationRoutes(routes map[string]server.CallHandler, ts service.Translation) {
	r := &translationRoutes{ts}
	{
		routes["getHistory"] = r.getHistory()
	}
}

type historyResponse struct {
	History []domain.Translation `json:"history"`
}

func (r *translationRoutes) getHistory() server.CallHandler {
	return func(d *amqp.Delivery) (interface{}, error) {
		translations, err := r.translationService.History()
		if err != nil {
			return nil, errors.Wrap(err, "amqp_rpc - translationRoutes - getHistory - r.translationService.History")
		}

		response := historyResponse{translations}

		return response, nil
	}
}
