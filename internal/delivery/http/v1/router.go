package v1

import (
	"github.com/evrone/go-service-template/internal/service"
	"github.com/gin-gonic/gin"
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
