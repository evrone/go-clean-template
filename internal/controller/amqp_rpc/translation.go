package amqprpc

import (
	"context"
	"fmt"

	"github.com/streadway/amqp"

	"github.com/evrone/go-clean-template/internal/domain/translation/entity"
	"github.com/evrone/go-clean-template/internal/usecase"
	"github.com/evrone/go-clean-template/pkg/rabbitmq/rmq_rpc/server"
)

type translationRoutes struct {
	translationUseCase *usecase.TranslationUseCase
}

func newTranslationRoutes(routes map[string]server.CallHandler, t *usecase.TranslationUseCase) {
	r := &translationRoutes{t}
	{
		routes["getHistory"] = r.getHistory()
	}
}

type historyResponse struct {
	History []entity.Translation `json:"history"`
}

func (r *translationRoutes) getHistory() server.CallHandler {
	return func(d *amqp.Delivery) (interface{}, error) {
		translations, err := r.translationUseCase.History(context.Background())
		if err != nil {
			return nil, fmt.Errorf("amqp_rpc - translationRoutes - getHistory - r.translationUseCase.History: %w", err)
		}

		response := historyResponse{translations}

		return response, nil
	}
}
