package response

import "github.com/evrone/go-clean-template/internal/entity"

// TaskList -.
type TaskList struct {
	Tasks []entity.Task `json:"tasks"`
	Total int           `json:"total"`
}

// DeleteStatus -.
type DeleteStatus struct {
	Status string `json:"status"`
}
