package request

import "time"

// CreateTodo -.
type CreateTodo struct {
	Title       string    `json:"title"       validate:"required"                      example:"Buy groceries"`
	Description string    `json:"description"                                          example:"Milk, eggs, bread"`
	Priority    string    `json:"priority"    validate:"required,oneof=low medium high" example:"medium"`
	DueDate     time.Time `json:"due_date"                                             example:"2026-03-01T00:00:00Z"`
}

// UpdateTodo -.
type UpdateTodo struct {
	Title       string    `json:"title"       example:"Buy groceries"`
	Description string    `json:"description" example:"Milk, eggs, bread"`
	Status      string    `json:"status"      validate:"omitempty,oneof=todo in_progress done"    example:"in_progress"`
	Priority    string    `json:"priority"    validate:"omitempty,oneof=low medium high"           example:"high"`
	DueDate     time.Time `json:"due_date"    example:"2026-03-01T00:00:00Z"`
}
