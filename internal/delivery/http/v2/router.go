package v2

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type router struct{}

func NewRouter(handler *gin.Engine) {
	api := handler.Group("/api/test")
	{
		api.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, "Hello from v2") })
	}
}
