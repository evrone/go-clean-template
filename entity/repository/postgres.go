package repository

import (
	"context"

	"github.com/evrone/go-service-template/domain"
)

type postgresEntityRepository struct {
	connect string
}

func NewPostgresEntityRepository(connect string) domain.EntityRepository {
	return &postgresEntityRepository{connect}
}

func (r *postgresEntityRepository) Get(ctx context.Context, entity domain.Entity) (domain.Entity, error) {
	return entity, nil
}
