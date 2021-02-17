// Package v2 is example version 2 HTTP API
package v2

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter(handler *gin.Engine) {
	api := handler.Group("/api/test")
	{
		api.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, "Hello from v2") })
	}
}
