package repo

import (
	"context"
	"fmt"

	"github.com/evrone/go-service-template/internal/domain"
	"github.com/evrone/go-service-template/pkg/postgres"
)

type TranslationRepo struct {
	*postgres.Postgres
}

func NewTranslationRepo(pg *postgres.Postgres) *TranslationRepo {
	return &TranslationRepo{pg}
}

func (r *TranslationRepo) GetHistory(ctx context.Context) ([]domain.Translation, error) {
	sql, _, err := r.Builder.
		Select("source, destination, original, translation").
		From("history").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("TranslationRepo - GetHistory - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("TranslationRepo - GetHistory - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	entities := make([]domain.Translation, 0, 64)

	for rows.Next() {
		e := domain.Translation{}

		err = rows.Scan(&e.Source, &e.Destination, &e.Original, &e.Translation)
		if err != nil {
			return nil, fmt.Errorf("TranslationRepo - GetHistory - rows.Scan: %w", err)
		}

		entities = append(entities, e)
	}

	return entities, nil
}

func (r *TranslationRepo) Store(ctx context.Context, entity domain.Translation) error {
	sql, args, err := r.Builder.
		Insert("history").
		Columns("source, destination, original, translation").
		Values(entity.Source, entity.Destination, entity.Original, entity.Translation).
		ToSql()
	if err != nil {
		return fmt.Errorf("TranslationRepo - Store - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("TranslationRepo - Store - r.Pool.Exec: %w", err)
	}

	return nil
}
