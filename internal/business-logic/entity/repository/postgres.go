package repository

import (
	"context"

	"github.com/evrone/go-service-template/pkg/postgres"

	"github.com/evrone/go-service-template/internal/business-logic/domain"
)

type postgresEntityRepository struct {
	postgres.Postgres
}

func NewPostgresEntityRepository(pg postgres.Postgres) domain.EntityRepository {
	return &postgresEntityRepository{pg}
}

func (p *postgresEntityRepository) GetHistory(ctx context.Context) ([]domain.Entity, error) {
	sql, _, err := p.Builder.
		Select("original, translation").
		From("history").
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := p.Pool.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entities := make([]domain.Entity, 0, 64)
	for rows.Next() {
		e := domain.Entity{}
		err = rows.Scan(&e.Original, &e.Translation)
		if err != nil {
			return nil, err
		}
		entities = append(entities, e)
	}

	return entities, nil
}

func (p *postgresEntityRepository) Store(ctx context.Context, entity domain.Entity) error {
	sql, args, err := p.Builder.
		Insert("history").
		Columns("original, translation").
		Values(entity.Original, entity.Translation).
		ToSql()
	if err != nil {
		return err
	}

	_, err = p.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}
