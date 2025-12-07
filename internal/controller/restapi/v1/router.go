package v1

import (
	"github.com/evrone/go-clean-template/internal/usecase"
	"github.com/evrone/go-clean-template/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// NewTranslationRoutes -.
func NewTranslationRoutes(apiV1Group fiber.Router, t usecase.Translation, l logger.Interface) {
	r := &V1{t: t, l: l, v: validator.New(validator.WithRequiredStructEnabled())}

	translationGroup := apiV1Group.Group("/translation")

	{
		translationGroup.Get("/history", r.history)
		translationGroup.Post("/do-translate", r.doTranslate)
	}
}
