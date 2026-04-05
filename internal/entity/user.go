package entity

import "time"

// User -.
type User struct {
	ID           string    `json:"id"         example:"550e8400-e29b-41d4-a716-446655440000"`
	Username     string    `json:"username"    example:"johndoe"`
	Email        string    `json:"email"       example:"john@example.com"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"  example:"2026-01-01T00:00:00Z"`
	UpdatedAt    time.Time `json:"updated_at"  example:"2026-01-01T00:00:00Z"`
} // @name entity.User
