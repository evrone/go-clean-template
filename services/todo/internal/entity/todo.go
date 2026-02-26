// Package entity defines main entities for business logic (services), data base mapping and
// HTTP response objects if suitable. Each logic group entities in own file.
package entity

import "time"

// Todo -.
type Todo struct {
	ID          int       `json:"id"          example:"1"`
	Title       string    `json:"title"       example:"Buy groceries"`
	Description string    `json:"description" example:"Milk, eggs, bread"`
	Status      string    `json:"status"      example:"todo"`
	Priority    string    `json:"priority"    example:"medium"`
	DueDate     time.Time `json:"due_date"    example:"2026-03-01T00:00:00Z"`
	CreatedAt   time.Time `json:"created_at"  example:"2026-02-26T00:00:00Z"`
	UpdatedAt   time.Time `json:"updated_at"  example:"2026-02-26T00:00:00Z"`
}
