package persistent

import (
	"context"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const _tracerName = "github.com/evrone/go-clean-template/internal/repo/persistent/task"

// tracedRepo wraps a TranslationRepo with OpenTelemetry spans.
type tracedRepo struct {
	next repo.TranslationRepo
}

func newTraced(next repo.TranslationRepo) repo.TranslationRepo {
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

func (r *tracedRepo) Store(ctx context.Context, userID string, t entity.Translation) error {
	ctx, span := startSpan(
		ctx, "TranslationRepo.Store",
		attribute.String("user.id", userID),
		attribute.String("translation.source", t.Source),
		attribute.String("translation.destination", t.Destination),
	)

	err := r.next.Store(ctx, userID, t)
	endSpan(span, err)

	return err
}

func (r *tracedRepo) GetHistory(ctx context.Context, userID string) ([]entity.Translation, error) {
	ctx, span := startSpan(ctx, "TranslationRepo.GetHistory", attribute.String("user.id", userID))

	result, err := r.next.GetHistory(ctx, userID)
	endSpan(span, err)

	return result, err
}
