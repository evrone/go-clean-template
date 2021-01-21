package domain

import "context"

type Entity struct {
	Msg string
}

type EntityUseCase interface {
	Do(ctx context.Context, entity Entity) error
}

type EntityRepository interface {
	Get(ctx context.Context, entity Entity) (Entity, error)
}

type EntityPublisher interface {
	Publish(ctx context.Context, entity Entity) error
}
