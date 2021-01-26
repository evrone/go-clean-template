package repository

import (
	"context"

	"github.com/evrone/go-service-template/pkg/postgres"

	"github.com/evrone/go-service-template/domain"
)

type postgresEntityRepository struct {
	postgres.Postgres
	tableName string
}

func NewPostgresEntityRepository(pg postgres.Postgres, tableName string) domain.EntityRepository {
	return &postgresEntityRepository{pg, tableName}
}

func (r *postgresEntityRepository) Get(ctx context.Context, entity domain.Entity) (domain.Entity, error) {
	return entity, nil
}
