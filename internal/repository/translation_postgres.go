package repository

import (
	"context"
	"fmt"

	"github.com/evrone/go-service-template/internal/domain"

	"github.com/evrone/go-service-template/pkg/postgres"
)

type translationRepository struct {
	*postgres.Postgres
}

func NewTranslationRepository(pg *postgres.Postgres) Translation {
	return &translationRepository{pg}
}

func (p *translationRepository) GetHistory(ctx context.Context) ([]domain.Translation, error) {
	sql, _, err := p.Builder.
		Select("source, destination, original, translation").
		From("history").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("translationRepository - GetHistory - p.Builder: %w", err)
	}

	rows, err := p.Pool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("translationRepository - GetHistory - p.Pool.Query: %w", err)
	}
	defer rows.Close()

	entities := make([]domain.Translation, 0, 64)
	for rows.Next() {
		e := domain.Translation{}
		err = rows.Scan(&e.Source, &e.Destination, &e.Original, &e.Translation)
		if err != nil {
			return nil, fmt.Errorf("translationRepository - GetHistory - rows.Scan: %w", err)
		}
		entities = append(entities, e)
	}

	return entities, nil
}

func (p *translationRepository) Store(ctx context.Context, entity domain.Translation) error {
	sql, args, err := p.Builder.
		Insert("history").
		Columns("source, destination, original, translation").
		Values(entity.Source, entity.Destination, entity.Original, entity.Translation).
		ToSql()
	if err != nil {
		return fmt.Errorf("translationRepository - Store - p.Builder: %w", err)
	}

	_, err = p.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("translationRepository - Store - p.Pool.Exec: %w", err)
	}

	return nil
}
