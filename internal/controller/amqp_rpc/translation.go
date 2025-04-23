package amqprpc

import (
	"context"
	"fmt"

	"github.com/evrone/go-clean-template/internal/usecase"
	"github.com/evrone/go-clean-template/pkg/rabbitmq/rmq_rpc/server"
	amqp "github.com/rabbitmq/amqp091-go"
)

type translationRoutes struct {
	translationUseCase usecase.Translation
}

func newTranslationRoutes(routes map[string]server.CallHandler, t usecase.Translation) {
	r := &translationRoutes{t}
	{
		routes["getHistory"] = r.getHistory()
	}
}

func (r *translationRoutes) getHistory() server.CallHandler {
	return func(_ *amqp.Delivery) (interface{}, error) {
		translationHistory, err := r.translationUseCase.History(context.Background())
		if err != nil {
			return nil, fmt.Errorf("amqp_rpc - translationRoutes - getHistory - r.translationUseCase.History: %w", err)
		}

		return translationHistory, nil
	}
}
