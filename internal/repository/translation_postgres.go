package repository

import (
	"context"
	"fmt"

	"github.com/evrone/go-service-template/internal/domain"
	"github.com/evrone/go-service-template/pkg/postgres"
)

type TranslationRepository struct {
	*postgres.Postgres
}

func NewTranslationRepository(pg *postgres.Postgres) *TranslationRepository {
	return &TranslationRepository{pg}
}

func (r *TranslationRepository) GetHistory(ctx context.Context) ([]domain.Translation, error) {
	sql, _, err := r.Builder.
		Select("source, destination, original, translation").
		From("history").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("TranslationRepository - GetHistory - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("TranslationRepository - GetHistory - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	entities := make([]domain.Translation, 0, 64)

	for rows.Next() {
		e := domain.Translation{}

		err = rows.Scan(&e.Source, &e.Destination, &e.Original, &e.Translation)
		if err != nil {
			return nil, fmt.Errorf("TranslationRepository - GetHistory - rows.Scan: %w", err)
		}

		entities = append(entities, e)
	}

	return entities, nil
}

func (r *TranslationRepository) Store(ctx context.Context, entity domain.Translation) error {
	sql, args, err := r.Builder.
		Insert("history").
		Columns("source, destination, original, translation").
		Values(entity.Source, entity.Destination, entity.Original, entity.Translation).
		ToSql()
	if err != nil {
		return fmt.Errorf("TranslationRepository - Store - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("TranslationRepository - Store - r.Pool.Exec: %w", err)
	}

	return nil
}
