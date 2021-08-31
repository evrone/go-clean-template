// Package service implements application business logic. Each logic group in own file.
package service

import (
	"context"

	"github.com/evrone/go-clean-template/internal/entity"
)

type (
	Translation interface {
		Translate(context.Context, entity.Translation) (entity.Translation, error)
		History(context.Context) ([]entity.Translation, error)
	}

	TranslationRepo interface {
		Store(context.Context, entity.Translation) error
		GetHistory(context.Context) ([]entity.Translation, error)
	}

	TranslationWebAPI interface {
		Translate(entity.Translation) (entity.Translation, error)
	}
)
