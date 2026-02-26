// Package repo implements application outer layer logic. Each logic group in own file.
package repo

import (
	"context"

	"github.com/evrone/todo-svc/internal/entity"
)

//go:generate mockgen -source=contracts.go -destination=../usecase/mocks_repo_test.go -package=usecase_test

type (
	// TodoRepo -.
	TodoRepo interface {
		Create(context.Context, entity.Todo) (entity.Todo, error)
		GetByID(context.Context, int) (entity.Todo, error)
		List(context.Context) ([]entity.Todo, error)
		Update(context.Context, int, entity.Todo) (entity.Todo, error)
		Delete(context.Context, int) error
	}
)
