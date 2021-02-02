package domain

import (
	"context"
)

type Entity struct {
	Original    string `json:"original" binding:"required"`
	Translation string `json:"translation"`
}

type EntityUseCase interface {
	DoTranslate(entity Entity) (Entity, error)
	History() ([]Entity, error)
}

type EntityTranslator interface {
	Translate(entity Entity) (Entity, error)
}

type EntityRepository interface {
	Store(ctx context.Context, entity Entity) error
	GetHistory(context.Context) ([]Entity, error)
}
