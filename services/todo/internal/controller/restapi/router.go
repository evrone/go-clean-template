package restapi

import (
	"net/http"

	"github.com/evrone/todo-svc/config"
	"evrone.local/common-middleware"
	v1 "github.com/evrone/todo-svc/internal/controller/restapi/v1"
	"github.com/evrone/todo-svc/internal/usecase"
	"evrone.local/common-pkg/logger"
	"github.com/gofiber/fiber/v2"
)

// NewRouter -.
// Swagger spec:
// @title       ToDo Service API
// @description CRUD REST microservice for managing todo items
// @version     1.0
// @host        localhost:8082
// @BasePath    /v1
func NewRouter(app *fiber.App, cfg *config.Config, t usecase.Todo, l logger.Interface) {
	// Middleware
	app.Use(middleware.Logger(l))
	app.Use(middleware.Recovery(l))

	// K8s probe
	app.Get("/healthz", func(ctx *fiber.Ctx) error { return ctx.SendStatus(http.StatusOK) })

	// Routers
	apiV1Group := app.Group("/v1")
	{
		v1.NewRoutes(apiV1Group, t, l)
	}
}
