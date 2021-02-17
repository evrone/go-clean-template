// Package v1 implements routing paths. Each services in own file.
package v1

import (
	"github.com/gin-gonic/gin"

	"github.com/evrone/go-service-template/internal/service"
)

type router struct {
	translationService service.Translation
}

func NewRouter(handler *gin.Engine, translationService service.Translation) {
	r := &router{translationService}

	api := handler.Group("/api/v1")
	{
		r.translationRouts(api)
	}
}
