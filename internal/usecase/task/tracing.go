package task

import (
	"context"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/usecase"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const _tracerName = "github.com/evrone/go-clean-template/internal/usecase/task"

// tracedUseCase wraps a Task usecase with OpenTelemetry spans, closing the
// gap between transport spans (HTTP/gRPC/AMQP/NATS) and repository spans.
type tracedUseCase struct {
	next usecase.Task
}

// newTraced wraps a Task usecase with tracing spans.
func newTraced(next usecase.Task) usecase.Task {
	return &tracedUseCase{next: next}
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

func (t *tracedUseCase) Create(ctx context.Context, userID, title, description string) (entity.Task, error) {
	ctx, span := startSpan(ctx, "TaskUseCase.Create", attribute.String("user.id", userID))

	result, err := t.next.Create(ctx, userID, title, description)
	endSpan(span, err)

	return result, err
}

func (t *tracedUseCase) Get(ctx context.Context, userID, taskID string) (entity.Task, error) {
	ctx, span := startSpan(
		ctx, "TaskUseCase.Get",
		attribute.String("user.id", userID),
		attribute.String("task.id", taskID),
	)

	result, err := t.next.Get(ctx, userID, taskID)
	endSpan(span, err)

	return result, err
}

func (t *tracedUseCase) List(
	ctx context.Context, userID string, status *entity.TaskStatus, limit, offset int,
) ([]entity.Task, int, error) {
	attrs := []attribute.KeyValue{
		attribute.String("user.id", userID),
		attribute.Int("task.limit", limit),
		attribute.Int("task.offset", offset),
	}
	if status != nil {
		attrs = append(attrs, attribute.String("task.status", string(*status)))
	}

	ctx, span := startSpan(ctx, "TaskUseCase.List", attrs...)

	tasks, total, err := t.next.List(ctx, userID, status, limit, offset)
	endSpan(span, err)

	return tasks, total, err
}

func (t *tracedUseCase) Update(ctx context.Context, userID, taskID, title, description string) (entity.Task, error) {
	ctx, span := startSpan(
		ctx, "TaskUseCase.Update",
		attribute.String("user.id", userID),
		attribute.String("task.id", taskID),
	)

	result, err := t.next.Update(ctx, userID, taskID, title, description)
	endSpan(span, err)

	return result, err
}

func (t *tracedUseCase) Transition(
	ctx context.Context, userID, taskID string, newStatus entity.TaskStatus,
) (entity.Task, error) {
	ctx, span := startSpan(
		ctx, "TaskUseCase.Transition",
		attribute.String("user.id", userID),
		attribute.String("task.id", taskID),
		attribute.String("task.new_status", string(newStatus)),
	)

	result, err := t.next.Transition(ctx, userID, taskID, newStatus)
	endSpan(span, err)

	return result, err
}

func (t *tracedUseCase) Delete(ctx context.Context, userID, taskID string) error {
	ctx, span := startSpan(
		ctx, "TaskUseCase.Delete",
		attribute.String("user.id", userID),
		attribute.String("task.id", taskID),
	)

	err := t.next.Delete(ctx, userID, taskID)
	endSpan(span, err)

	return err
}
