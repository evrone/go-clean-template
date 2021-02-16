package repository

import (
	"context"

	"github.com/evrone/go-service-template/internal/domain"
)

type Translation interface {
	Store(ctx context.Context, entity domain.Translation) error
	GetHistory(context.Context) ([]domain.Translation, error)
}
