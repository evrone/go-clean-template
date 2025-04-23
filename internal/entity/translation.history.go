// Package entity defines main entities for business logic (services), data base mapping and
// HTTP response objects if suitable. Each logic group entities in own file.
package entity

// TranslationHistory -.
type TranslationHistory struct {
	History []Translation `json:"history"`
}
