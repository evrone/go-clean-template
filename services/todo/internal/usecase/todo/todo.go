package todo

import (
	"context"
	"fmt"

	"github.com/evrone/todo-svc/internal/entity"
	"github.com/evrone/todo-svc/internal/repo"
)

// UseCase -.
type UseCase struct {
	repo repo.TodoRepo
}

// New -.
func New(r repo.TodoRepo) *UseCase {
	return &UseCase{repo: r}
}

// Create -.
func (uc *UseCase) Create(ctx context.Context, t entity.Todo) (entity.Todo, error) {
	if t.Status == "" {
		t.Status = "todo"
	}

	if t.Priority == "" {
		t.Priority = "medium"
	}

	todo, err := uc.repo.Create(ctx, t)
	if err != nil {
		return entity.Todo{}, fmt.Errorf("TodoUseCase - Create - uc.repo.Create: %w", err)
	}

	return todo, nil
}

// GetByID -.
func (uc *UseCase) GetByID(ctx context.Context, id int) (entity.Todo, error) {
	todo, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return entity.Todo{}, fmt.Errorf("TodoUseCase - GetByID - uc.repo.GetByID: %w", err)
	}

	return todo, nil
}

// List -.
func (uc *UseCase) List(ctx context.Context) ([]entity.Todo, error) {
	todos, err := uc.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("TodoUseCase - List - uc.repo.List: %w", err)
	}

	return todos, nil
}

// Update -.
func (uc *UseCase) Update(ctx context.Context, id int, t entity.Todo) (entity.Todo, error) {
	todo, err := uc.repo.Update(ctx, id, t)
	if err != nil {
		return entity.Todo{}, fmt.Errorf("TodoUseCase - Update - uc.repo.Update: %w", err)
	}

	return todo, nil
}

// Delete -.
func (uc *UseCase) Delete(ctx context.Context, id int) error {
	err := uc.repo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("TodoUseCase - Delete - uc.repo.Delete: %w", err)
	}

	return nil
}
