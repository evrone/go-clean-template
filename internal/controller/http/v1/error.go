package v1

import (
	"github.com/gin-gonic/gin"

	"github.com/evrone/go-clean-template/pkg/logger"
)

type response struct {
	Error string `json:"error" example:"message"`
}

func errorResponse(c *gin.Context, code int, err error, msg string) {
	logger.Error(err, "http - v1 - errorResponse",
		logger.Field{Key: "path", Val: c.FullPath()},
		logger.Field{Key: "request_method", Val: c.Request.Method},
		logger.Field{Key: "response_code", Val: code},
	)
	c.AbortWithStatusJSON(code, response{msg})
}
