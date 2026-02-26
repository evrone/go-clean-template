package v1

import (
	"github.com/evrone/todo-svc/internal/usecase"
	"evrone.local/common-pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// NewRoutes -.
func NewRoutes(apiV1Group fiber.Router, t usecase.Todo, l logger.Interface) {
	r := &V1{todo: t, l: l, v: validator.New(validator.WithRequiredStructEnabled())}

	todoGroup := apiV1Group.Group("/todo")

	{
		todoGroup.Post("/", r.create)
		todoGroup.Get("/", r.list)
		todoGroup.Get("/:id", r.getByID)
		todoGroup.Put("/:id", r.update)
		todoGroup.Delete("/:id", r.delete)
	}
}
