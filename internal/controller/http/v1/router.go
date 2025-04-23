package v1

import (
	"github.com/evrone/go-clean-template/internal/usecase"
	"github.com/evrone/go-clean-template/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type translationRoutes struct {
	t usecase.Translation
	l logger.Interface
	v *validator.Validate
}

func NewTranslationRouter(apiV1Group fiber.Router, t usecase.Translation, l logger.Interface) {
	r := &translationRoutes{t, l, validator.New(validator.WithRequiredStructEnabled())}

	translationGroup := apiV1Group.Group("/translation")
	{
		translationGroup.Get("/history", r.history)
		translationGroup.Post("/do-translate", r.doTranslate)
	}
}
