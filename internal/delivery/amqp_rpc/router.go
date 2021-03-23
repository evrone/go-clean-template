package amqprpc

import (
	"github.com/evrone/go-service-template/internal/service"
	"github.com/evrone/go-service-template/pkg/rabbitmq/rmq_rpc/server"
)

type router struct {
	translationService service.Translation
	routerMap          map[string]server.CallHandler
}

func NewRouter(translationService service.Translation) map[string]server.CallHandler {
	r := &router{
		translationService: translationService,
		routerMap:          make(map[string]server.CallHandler),
	}

	r.translationRoutes()

	return r.routerMap
}
