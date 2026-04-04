package entity_test

import (
	"testing"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTask_Transition(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		from      entity.TaskStatus
		to        entity.TaskStatus
		wantErr   bool
		wantState entity.TaskStatus
	}{
		{"todo to in_progress", entity.TaskStatusTodo, entity.TaskStatusInProgress, false, entity.TaskStatusInProgress},
		{"in_progress to done", entity.TaskStatusInProgress, entity.TaskStatusDone, false, entity.TaskStatusDone},
		{"in_progress to todo", entity.TaskStatusInProgress, entity.TaskStatusTodo, false, entity.TaskStatusTodo},
		{"todo to done (invalid)", entity.TaskStatusTodo, entity.TaskStatusDone, true, entity.TaskStatusTodo},
		{"done to todo (invalid)", entity.TaskStatusDone, entity.TaskStatusTodo, true, entity.TaskStatusDone},
		{"done to in_progress (invalid)", entity.TaskStatusDone, entity.TaskStatusInProgress, true, entity.TaskStatusDone},
		{"unknown status (invalid)", entity.TaskStatus("unknown"), entity.TaskStatusTodo, true, entity.TaskStatus("unknown")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			task := entity.Task{Status: tt.from}
			err := task.Transition(tt.to)

			if tt.wantErr {
				require.ErrorIs(t, err, entity.ErrInvalidTransition)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.wantState, task.Status)
		})
	}
}
