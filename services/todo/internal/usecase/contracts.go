// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"

	"github.com/evrone/todo-svc/internal/entity"
)

//go:generate mockgen -source=contracts.go -destination=./mocks_usecase_test.go -package=usecase_test

type (
	// Todo -.
	Todo interface {
		Create(context.Context, entity.Todo) (entity.Todo, error)
		GetByID(context.Context, int) (entity.Todo, error)
		List(context.Context) ([]entity.Todo, error)
		Update(context.Context, int, entity.Todo) (entity.Todo, error)
		Delete(context.Context, int) error
	}
)
