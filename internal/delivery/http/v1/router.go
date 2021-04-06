// Package v1 implements routing paths. Each services in own file.
package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// Swagger docs.
	_ "github.com/evrone/go-service-template/docs"
	"github.com/evrone/go-service-template/internal/service"
)

// Swagger spec:
// @title       Go Service Template API
// @description Using a translation service as an example
// @version     1.0
// @host        localhost:8080
// @BasePath    /api/v1

func NewRouter(handler *gin.Engine, translationService service.Translation) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// Swagger
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/swagger/*any", swaggerHandler)

	// K8s probe
	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Routers
	h := handler.Group("/api/v1")
	{
		newTranslationRoutes(h, translationService)
	}
}
