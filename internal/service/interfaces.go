// Package service implements application business logic. Each logic group in own file.
package service

import (
	"context"

	"github.com/evrone/go-clean-template/internal/domain"
)

type (
	Translation interface {
		Translate(context.Context, domain.Translation) (domain.Translation, error)
		History(context.Context) ([]domain.Translation, error)
	}

	TranslationRepo interface {
		Store(context.Context, domain.Translation) error
		GetHistory(context.Context) ([]domain.Translation, error)
	}

	TranslationWebAPI interface {
		Translate(domain.Translation) (domain.Translation, error)
	}
)
