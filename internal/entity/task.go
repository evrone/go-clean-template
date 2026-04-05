package entity

import (
	"slices"
	"time"
)

// TaskStatus -.
type TaskStatus string // @name entity.TaskStatus

const (
	TaskStatusTodo       TaskStatus = "todo"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusDone       TaskStatus = "done"
)

// Task -.
type Task struct {
	ID          string     `json:"id"          example:"550e8400-e29b-41d4-a716-446655440000"`
	UserID      string     `json:"user_id"     example:"550e8400-e29b-41d4-a716-446655440000"`
	Title       string     `json:"title"       example:"My task"`
	Description string     `json:"description" example:"Task description"`
	Status      TaskStatus `json:"status"      example:"todo"`
	CreatedAt   time.Time  `json:"created_at"  example:"2026-01-01T00:00:00Z"`
	UpdatedAt   time.Time  `json:"updated_at"  example:"2026-01-01T00:00:00Z"`
} // @name entity.Task

// Valid reports whether s is a known task status.
func (s TaskStatus) Valid() bool {
	switch s {
	case TaskStatusTodo, TaskStatusInProgress, TaskStatusDone:
		return true
	default:
		return false
	}
}

// Transition validates and applies a status transition.
func (t *Task) Transition(newStatus TaskStatus) error {
	validTransitions := map[TaskStatus][]TaskStatus{
		TaskStatusTodo:       {TaskStatusInProgress},
		TaskStatusInProgress: {TaskStatusDone, TaskStatusTodo},
		TaskStatusDone:       {},
	}

	allowed, ok := validTransitions[t.Status]
	if !ok {
		return ErrInvalidTransition
	}

	if slices.Contains(allowed, newStatus) {
		t.Status = newStatus

		return nil
	}

	return ErrInvalidTransition
}
