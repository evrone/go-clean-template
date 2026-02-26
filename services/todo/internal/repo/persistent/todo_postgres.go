package persistent

import (
	"context"
	"fmt"

	"github.com/evrone/todo-svc/internal/entity"
	"evrone.local/common-pkg/postgres"
)

const _defaultEntityCap = 64

// TodoRepo -.
type TodoRepo struct {
	*postgres.Postgres
}

// NewTodoRepo -.
func NewTodoRepo(pg *postgres.Postgres) *TodoRepo {
	return &TodoRepo{pg}
}

// Create -.
func (r *TodoRepo) Create(ctx context.Context, t entity.Todo) (entity.Todo, error) {
	sql, args, err := r.Builder.
		Insert("todos").
		Columns("title, description, status, priority, due_date").
		Values(t.Title, t.Description, t.Status, t.Priority, t.DueDate).
		Suffix("RETURNING id, created_at, updated_at").
		ToSql()
	if err != nil {
		return entity.Todo{}, fmt.Errorf("TodoRepo - Create - r.Builder: %w", err)
	}

	row := r.Pool.QueryRow(ctx, sql, args...)

	err = row.Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return entity.Todo{}, fmt.Errorf("TodoRepo - Create - row.Scan: %w", err)
	}

	return t, nil
}

// GetByID -.
func (r *TodoRepo) GetByID(ctx context.Context, id int) (entity.Todo, error) {
	sql, args, err := r.Builder.
		Select("id, title, description, status, priority, due_date, created_at, updated_at").
		From("todos").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		return entity.Todo{}, fmt.Errorf("TodoRepo - GetByID - r.Builder: %w", err)
	}

	var t entity.Todo

	row := r.Pool.QueryRow(ctx, sql, args...)

	err = row.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.DueDate, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return entity.Todo{}, fmt.Errorf("TodoRepo - GetByID - row.Scan: %w", err)
	}

	return t, nil
}

// List -.
func (r *TodoRepo) List(ctx context.Context) ([]entity.Todo, error) {
	sql, _, err := r.Builder.
		Select("id, title, description, status, priority, due_date, created_at, updated_at").
		From("todos").
		OrderBy("created_at DESC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("TodoRepo - List - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("TodoRepo - List - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	todos := make([]entity.Todo, 0, _defaultEntityCap)

	for rows.Next() {
		var t entity.Todo

		err = rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.DueDate, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("TodoRepo - List - rows.Scan: %w", err)
		}

		todos = append(todos, t)
	}

	return todos, nil
}

// Update -.
func (r *TodoRepo) Update(ctx context.Context, id int, t entity.Todo) (entity.Todo, error) {
	sql, args, err := r.Builder.
		Update("todos").
		Set("title", t.Title).
		Set("description", t.Description).
		Set("status", t.Status).
		Set("priority", t.Priority).
		Set("due_date", t.DueDate).
		Set("updated_at", "NOW()").
		Where("id = ?", id).
		Suffix("RETURNING id, title, description, status, priority, due_date, created_at, updated_at").
		ToSql()
	if err != nil {
		return entity.Todo{}, fmt.Errorf("TodoRepo - Update - r.Builder: %w", err)
	}

	row := r.Pool.QueryRow(ctx, sql, args...)

	var updated entity.Todo

	err = row.Scan(
		&updated.ID, &updated.Title, &updated.Description,
		&updated.Status, &updated.Priority, &updated.DueDate,
		&updated.CreatedAt, &updated.UpdatedAt,
	)
	if err != nil {
		return entity.Todo{}, fmt.Errorf("TodoRepo - Update - row.Scan: %w", err)
	}

	return updated, nil
}

// Delete -.
func (r *TodoRepo) Delete(ctx context.Context, id int) error {
	sql, args, err := r.Builder.
		Delete("todos").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		return fmt.Errorf("TodoRepo - Delete - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("TodoRepo - Delete - r.Pool.Exec: %w", err)
	}

	return nil
}
