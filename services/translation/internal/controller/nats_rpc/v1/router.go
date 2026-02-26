package v1

import (
	"github.com/evrone/translation-svc/internal/usecase"
	"evrone.local/common-pkg/logger"
	"evrone.local/common-pkg/nats/nats_rpc/server"
	"github.com/go-playground/validator/v10"
)

// NewTranslationRoutes -.
func NewTranslationRoutes(routes map[string]server.CallHandler, t usecase.Translation, l logger.Interface) {
	r := &V1{t: t, l: l, v: validator.New(validator.WithRequiredStructEnabled())}

	{
		routes["v1.getHistory"] = r.getHistory()
	}
}
