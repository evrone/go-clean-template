package v1

import (
	"github.com/evrone/todo-svc/internal/usecase"
	"evrone.local/common-pkg/logger"
	"github.com/go-playground/validator/v10"
)

// V1 -.
type V1 struct {
	todo usecase.Todo
	l    logger.Interface
	v    *validator.Validate
}
