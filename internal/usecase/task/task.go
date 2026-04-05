package task

import (
	"context"
	"fmt"
	"time"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
	"github.com/google/uuid"
)

// UseCase -.
type UseCase struct {
	repo repo.TaskRepo
}

// New -.
func New(r repo.TaskRepo) *UseCase {
	return &UseCase{repo: r}
}

// Create -.
func (uc *UseCase) Create(ctx context.Context, userID, title, description string) (entity.Task, error) {
	now := time.Now().UTC()

	task := entity.Task{
		ID:          uuid.New().String(),
		UserID:      userID,
		Title:       title,
		Description: description,
		Status:      entity.TaskStatusTodo,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	err := uc.repo.Store(ctx, &task)
	if err != nil {
		return entity.Task{}, fmt.Errorf("TaskUseCase - Create - uc.repo.Store: %w", err)
	}

	return task, nil
}

// Get -.
func (uc *UseCase) Get(ctx context.Context, userID, taskID string) (entity.Task, error) {
	task, err := uc.repo.GetByID(ctx, userID, taskID)
	if err != nil {
		return entity.Task{}, fmt.Errorf("TaskUseCase - Get - uc.repo.GetByID: %w", err)
	}

	return task, nil
}

// List -.
func (uc *UseCase) List(ctx context.Context, userID string, status *entity.TaskStatus, limit, offset int) ([]entity.Task, int, error) {
	if limit <= 0 {
		limit = 10
	}

	if offset < 0 {
		offset = 0
	}

	tasks, total, err := uc.repo.List(ctx, userID, repo.TaskFilter{
		Status: status,
		Limit:  uint64(limit),
		Offset: uint64(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("TaskUseCase - List - uc.repo.List: %w", err)
	}

	return tasks, total, nil
}

// Update -.
func (uc *UseCase) Update(ctx context.Context, userID, taskID, title, description string) (entity.Task, error) {
	now := time.Now().UTC()

	task, err := uc.repo.GetByID(ctx, userID, taskID)
	if err != nil {
		return entity.Task{}, fmt.Errorf("TaskUseCase - Update - uc.repo.GetByID: %w", err)
	}

	task.Title = title
	task.Description = description
	task.UpdatedAt = now

	err = uc.repo.Update(ctx, &task)
	if err != nil {
		return entity.Task{}, fmt.Errorf("TaskUseCase - Update - uc.repo.Update: %w", err)
	}

	return task, nil
}

// Transition -.
func (uc *UseCase) Transition(ctx context.Context, userID, taskID string, newStatus entity.TaskStatus) (entity.Task, error) {
	now := time.Now().UTC()

	task, err := uc.repo.GetByID(ctx, userID, taskID)
	if err != nil {
		return entity.Task{}, fmt.Errorf("TaskUseCase - Transition - uc.repo.GetByID: %w", err)
	}

	err = task.Transition(newStatus)
	if err != nil {
		return entity.Task{}, err
	}

	task.UpdatedAt = now

	err = uc.repo.Update(ctx, &task)
	if err != nil {
		return entity.Task{}, fmt.Errorf("TaskUseCase - Transition - uc.repo.Update: %w", err)
	}

	return task, nil
}

// Delete -.
func (uc *UseCase) Delete(ctx context.Context, userID, taskID string) error {
	err := uc.repo.Delete(ctx, userID, taskID)
	if err != nil {
		return fmt.Errorf("TaskUseCase - Delete - uc.repo.Delete: %w", err)
	}

	return nil
}
