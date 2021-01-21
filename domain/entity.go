package domain

import "context"

type Entity struct{}

type EntityUsecase interface {
	Get(ctx context.Context, ID int) (string, error)
}

type EntityRepository interface {
	GetByID(ctx context.Context, ID int) (string, error)
}
