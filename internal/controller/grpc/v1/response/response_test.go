package response_test

import (
	"testing"
	"time"

	"github.com/evrone/go-clean-template/internal/controller/grpc/v1/response"
	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type userResponseFields struct {
	id        string
	username  string
	email     string
	createdAt string
	updatedAt string
}

type userResponseGetter interface {
	GetId() string
	GetUsername() string
	GetEmail() string
	GetCreatedAt() string
	GetUpdatedAt() string
}

func assertUserResponseFields(t *testing.T, f *userResponseFields, got userResponseGetter) {
	t.Helper()

	require.NotNil(t, got)
	assert.Equal(t, f.id, got.GetId())
	assert.Equal(t, f.username, got.GetUsername())
	assert.Equal(t, f.email, got.GetEmail())
	assert.Equal(t, f.createdAt, got.GetCreatedAt())
	assert.Equal(t, f.updatedAt, got.GetUpdatedAt())
}

func TestNewRegisterResponse(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	user := &entity.User{
		ID:        "user-id-123",
		Username:  "testuser",
		Email:     "test@example.com",
		CreatedAt: now,
		UpdatedAt: now,
	}

	resp := response.NewRegisterResponse(user)

	assertUserResponseFields(t, &userResponseFields{
		id:        user.ID,
		username:  user.Username,
		email:     user.Email,
		createdAt: "2026-01-01T00:00:00Z",
		updatedAt: "2026-01-01T00:00:00Z",
	}, resp)
}

func TestNewGetProfileResponse(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 3, 15, 12, 30, 0, 0, time.UTC)
	user := &entity.User{
		ID:        "user-id-456",
		Username:  "anotheruser",
		Email:     "another@example.com",
		CreatedAt: now,
		UpdatedAt: now,
	}

	resp := response.NewGetProfileResponse(user)

	assertUserResponseFields(t, &userResponseFields{
		id:        user.ID,
		username:  user.Username,
		email:     user.Email,
		createdAt: "2026-03-15T12:30:00Z",
		updatedAt: "2026-03-15T12:30:00Z",
	}, resp)
}

func TestNewTaskResponse(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 2, 10, 8, 0, 0, 0, time.UTC)
	task := &entity.Task{
		ID:          "task-id-789",
		UserID:      "user-id-123",
		Title:       "My Task",
		Description: "Task description",
		Status:      entity.TaskStatusInProgress,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	resp := response.NewTaskResponse(task)

	require.NotNil(t, resp)
	assert.Equal(t, task.ID, resp.Id)
	assert.Equal(t, task.UserID, resp.UserId)
	assert.Equal(t, task.Title, resp.Title)
	assert.Equal(t, task.Description, resp.Description)
	assert.Equal(t, string(task.Status), resp.Status)
	assert.Equal(t, "2026-02-10T08:00:00Z", resp.CreatedAt)
	assert.Equal(t, "2026-02-10T08:00:00Z", resp.UpdatedAt)
}

func TestNewListTasksResponse(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	tasks := []entity.Task{
		{
			ID:        "task-1",
			UserID:    "user-id-123",
			Title:     "Task One",
			Status:    entity.TaskStatusTodo,
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        "task-2",
			UserID:    "user-id-123",
			Title:     "Task Two",
			Status:    entity.TaskStatusDone,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	resp := response.NewListTasksResponse(tasks, 2)

	require.NotNil(t, resp)
	assert.Len(t, resp.Tasks, 2)
	assert.Equal(t, int32(2), resp.Total)
	assert.Equal(t, "task-1", resp.Tasks[0].Id)
	assert.Equal(t, "task-2", resp.Tasks[1].Id)
}

func TestNewListTasksResponse_Empty(t *testing.T) {
	t.Parallel()

	resp := response.NewListTasksResponse([]entity.Task{}, 0)

	require.NotNil(t, resp)
	assert.Empty(t, resp.Tasks)
	assert.Equal(t, int32(0), resp.Total)
}

func TestNewTranslationHistory(t *testing.T) {
	t.Parallel()

	history := entity.TranslationHistory{
		History: []entity.Translation{
			{
				Source:      "en",
				Destination: "ru",
				Original:    "hello",
				Translation: "привет",
			},
			{
				Source:      "ru",
				Destination: "en",
				Original:    "мир",
				Translation: "world",
			},
		},
	}

	resp := response.NewTranslationHistory(history)

	require.NotNil(t, resp)
	require.Len(t, resp.History, 2)
	assert.Equal(t, "en", resp.History[0].Source)
	assert.Equal(t, "ru", resp.History[0].Destination)
	assert.Equal(t, "hello", resp.History[0].Original)
	assert.Equal(t, "привет", resp.History[0].Translation)
	assert.Equal(t, "ru", resp.History[1].Source)
	assert.Equal(t, "en", resp.History[1].Destination)
	assert.Equal(t, "мир", resp.History[1].Original)
	assert.Equal(t, "world", resp.History[1].Translation)
}
