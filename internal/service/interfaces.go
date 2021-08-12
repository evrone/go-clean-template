// Package service implements application business logic. Each logic group in own file.
package service

import (
	"context"

	"github.com/evrone/go-clean-template/internal/domain"
)

type Translation interface {
	Translate(domain.Translation) (domain.Translation, error)
	History() ([]domain.Translation, error)
}

type TranslationRepo interface {
	Store(context.Context, domain.Translation) error
	GetHistory(context.Context) ([]domain.Translation, error)
}

type TranslationWebAPI interface {
	Translate(domain.Translation) (domain.Translation, error)
}
