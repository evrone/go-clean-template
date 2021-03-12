package amqprpc

import (
	"github.com/evrone/go-service-template/internal/service"
	"github.com/evrone/go-service-template/pkg/rmq"
)

type router struct {
	translationService service.Translation
	routerMap          map[string]rmq.CallHandler
}

func NewRouter(translationService service.Translation) map[string]rmq.CallHandler {
	r := &router{
		translationService: translationService,
		routerMap:          make(map[string]rmq.CallHandler),
	}

	r.translationRoutes()

	return r.routerMap
}
