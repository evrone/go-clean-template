package v1

import (
	"github.com/gin-gonic/gin"
)

type response struct {
	Error string `json:"error" example:"message"`
}

func errorResponse(c *gin.Context, code int, err error, msg string) {
	c.AbortWithStatusJSON(code, response{msg})
}
