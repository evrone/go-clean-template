package amqprpc

import (
	"github.com/evrone/go-clean-template/internal/service"
	"github.com/evrone/go-clean-template/pkg/rabbitmq/rmq_rpc/server"
)

func NewRouter(translationService service.Translation) map[string]server.CallHandler {
	routes := make(map[string]server.CallHandler)
	{
		newTranslationRoutes(routes, translationService)
	}

	return routes
}
