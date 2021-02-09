package v2

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type router struct{}

func NewAPIRouter(handler *gin.Engine) {
	api := handler.Group("/api/v2")
	{
		api.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, "Hello from v2") })
	}
}
