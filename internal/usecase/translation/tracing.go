package translation

import (
	"context"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/usecase"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const _tracerName = "github.com/evrone/go-clean-template/internal/usecase/translation"

// tracedUseCase wraps a Translation usecase with OpenTelemetry spans, closing
// the gap between transport spans (HTTP/gRPC/AMQP/NATS) and repository spans.
type tracedUseCase struct {
	next usecase.Translation
}

// newTraced wraps a Translation usecase with tracing spans.
func newTraced(next usecase.Translation) usecase.Translation {
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

func (t *tracedUseCase) Translate(ctx context.Context, userID string, tr entity.Translation) (entity.Translation, error) {
	ctx, span := startSpan(
		ctx, "TranslationUseCase.Translate",
		attribute.String("user.id", userID),
		attribute.String("translation.source", tr.Source),
		attribute.String("translation.destination", tr.Destination),
	)

	result, err := t.next.Translate(ctx, userID, tr)
	endSpan(span, err)

	return result, err
}

func (t *tracedUseCase) History(ctx context.Context, userID string) (entity.TranslationHistory, error) {
	ctx, span := startSpan(ctx, "TranslationUseCase.History", attribute.String("user.id", userID))

	result, err := t.next.History(ctx, userID)
	endSpan(span, err)

	return result, err
}
