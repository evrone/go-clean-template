package repository

import (
	"context"

	"github.com/evrone/go-service-template/internal/domain"

	"github.com/evrone/go-service-template/pkg/postgres"
)

type postgresEntityRepository struct {
	postgres.Postgres
}

func NewPostgresEntityRepository(pg postgres.Postgres) Translation {
	return &postgresEntityRepository{pg}
}

func (p *postgresEntityRepository) GetHistory(ctx context.Context) ([]domain.Translation, error) {
	sql, _, err := p.Builder.
		Select("source, destination, original, translation").
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

	entities := make([]domain.Translation, 0, 64)
	for rows.Next() {
		e := domain.Translation{}
		err = rows.Scan(&e.Source, &e.Destination, &e.Original, &e.Translation)
		if err != nil {
			return nil, err
		}
		entities = append(entities, e)
	}

	return entities, nil
}

func (p *postgresEntityRepository) Store(ctx context.Context, entity domain.Translation) error {
	sql, args, err := p.Builder.
		Insert("history").
		Columns("source, destination, original, translation").
		Values(entity.Source, entity.Destination, entity.Original, entity.Translation).
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
