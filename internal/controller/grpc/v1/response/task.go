package response

import (
	"math"

	v1 "github.com/evrone/go-clean-template/docs/proto/v1"
	"github.com/evrone/go-clean-template/internal/entity"
)

// NewTaskResponse -.
func NewTaskResponse(task *entity.Task) *v1.TaskResponse {
	return &v1.TaskResponse{
		Id:          task.ID,
		UserId:      task.UserID,
		Title:       task.Title,
		Description: task.Description,
		Status:      string(task.Status),
		CreatedAt:   task.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   task.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

// NewListTasksResponse -.
func NewListTasksResponse(tasks []entity.Task, total int) *v1.ListTasksResponse {
	pbTasks := make([]*v1.TaskResponse, len(tasks))
	for i := range tasks {
		pbTasks[i] = NewTaskResponse(&tasks[i])
	}

	if total > math.MaxInt32 {
		total = math.MaxInt32
	}

	return &v1.ListTasksResponse{
		Tasks: pbTasks,
		Total: int32(total),
	}
}
