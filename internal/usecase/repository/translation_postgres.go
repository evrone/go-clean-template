package repository

import (
	"context"
	"fmt"

	"github.com/evrone/go-clean-template/internal/domain/translation/entity"
	"github.com/evrone/go-clean-template/pkg/postgres"
)

const _defaultEntityCap = 64

// TranslationRepository -.
type TranslationRepository struct {
	*postgres.Postgres
}

// New -.
func New(pg *postgres.Postgres) *TranslationRepository {
	return &TranslationRepository{pg}
}

// GetHistory -.
func (r *TranslationRepository) GetHistory(ctx context.Context) ([]entity.Translation, error) {
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

	entities := make([]entity.Translation, 0, _defaultEntityCap)

	for rows.Next() {
		e := entity.Translation{}

		err = rows.Scan(&e.Source, &e.Destination, &e.Original, &e.Translation)
		if err != nil {
			return nil, fmt.Errorf("TranslationRepository - GetHistory - rows.Scan: %w", err)
		}

		entities = append(entities, e)
	}

	return entities, nil
}

// Store -.
func (r *TranslationRepository) Store(ctx context.Context, t entity.Translation) error {
	sql, args, err := r.Builder.
		Insert("history").
		Columns("source, destination, original, translation").
		Values(t.Source, t.Destination, t.Original, t.Translation).
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
