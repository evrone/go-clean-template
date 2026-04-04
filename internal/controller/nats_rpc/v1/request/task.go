package request

// CreateTask -.
type CreateTask struct {
	Title       string `json:"title"       validate:"required,max=255"`
	Description string `json:"description" validate:"max=1000"`
}

// GetTask -.
type GetTask struct {
	ID string `json:"id" validate:"required"`
}

// ListTasks -.
type ListTasks struct {
	Status string `json:"status" validate:"omitempty,oneof=todo in_progress done"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

// UpdateTask -.
type UpdateTask struct {
	ID          string `json:"id"          validate:"required"`
	Title       string `json:"title"       validate:"required,max=255"`
	Description string `json:"description" validate:"max=1000"`
}

// TransitionTask -.
type TransitionTask struct {
	ID     string `json:"id"     validate:"required"`
	Status string `json:"status" validate:"required,oneof=todo in_progress done"`
}

// DeleteTask -.
type DeleteTask struct {
	ID string `json:"id" validate:"required"`
}
