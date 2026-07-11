package task

import (
	"context"
	"math"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const _tracerName = "github.com/evrone/go-clean-template/internal/repo/persistent/task"

// tracedRepo wraps a TaskRepo with OpenTelemetry spans on top of the
// low-level pgx query spans, giving a semantic "TaskRepo.<Method>" view.
type tracedRepo struct {
	next repo.TaskRepo
}

func newTraced(next repo.TaskRepo) repo.TaskRepo {
	return &tracedRepo{next: next}
}

func startSpan(ctx context.Context, name string, attrs ...attribute.KeyValue) (context.Context, trace.Span) {
	return otel.Tracer(_tracerName).Start(ctx, name, trace.WithAttributes(attrs...))
}

func endSpan(span trace.Span, err error) {
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	span.End()
}

func safeUint64ToInt64(v uint64) int64 {
	if v > math.MaxInt64 {
		return math.MaxInt64
	}

	return int64(v)
}

func (r *tracedRepo) Store(ctx context.Context, task *entity.Task) error {
	ctx, span := startSpan(
		ctx, "TaskRepo.Store",
		attribute.String("user.id", task.UserID),
		attribute.String("task.id", task.ID),
	)

	err := r.next.Store(ctx, task)
	endSpan(span, err)

	return err
}

func (r *tracedRepo) GetByID(ctx context.Context, userID, taskID string) (entity.Task, error) {
	ctx, span := startSpan(
		ctx, "TaskRepo.GetByID",
		attribute.String("user.id", userID),
		attribute.String("task.id", taskID),
	)

	result, err := r.next.GetByID(ctx, userID, taskID)
	endSpan(span, err)

	return result, err
}

func (r *tracedRepo) List(ctx context.Context, userID string, filter repo.TaskFilter) ([]entity.Task, int, error) {
	attrs := []attribute.KeyValue{
		attribute.String("user.id", userID),
		attribute.Int64("task.limit", safeUint64ToInt64(filter.Limit)),
		attribute.Int64("task.offset", safeUint64ToInt64(filter.Offset)),
	}
	if filter.Status != nil {
		attrs = append(attrs, attribute.String("task.status", string(*filter.Status)))
	}

	ctx, span := startSpan(ctx, "TaskRepo.List", attrs...)

	tasks, total, err := r.next.List(ctx, userID, filter)
	endSpan(span, err)

	return tasks, total, err
}

func (r *tracedRepo) Update(ctx context.Context, task *entity.Task) error {
	ctx, span := startSpan(
		ctx, "TaskRepo.Update",
		attribute.String("user.id", task.UserID),
		attribute.String("task.id", task.ID),
	)

	err := r.next.Update(ctx, task)
	endSpan(span, err)

	return err
}

func (r *tracedRepo) Delete(ctx context.Context, userID, taskID string) error {
	ctx, span := startSpan(
		ctx, "TaskRepo.Delete",
		attribute.String("user.id", userID),
		attribute.String("task.id", taskID),
	)

	err := r.next.Delete(ctx, userID, taskID)
	endSpan(span, err)

	return err
}
