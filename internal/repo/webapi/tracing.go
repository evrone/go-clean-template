package webapi

import (
	"context"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const _tracerName = "github.com/evrone/go-clean-template/internal/repo/webapi"

// tracedTranslationWebAPI wraps a TranslationWebAPI with OpenTelemetry spans
// for the outbound call to the external translation service.
type tracedTranslationWebAPI struct {
	next repo.TranslationWebAPI
}

func newTraced(next repo.TranslationWebAPI) repo.TranslationWebAPI {
	return &tracedTranslationWebAPI{next: next}
}

func (t *tracedTranslationWebAPI) Translate(
	ctx context.Context, translation entity.Translation,
) (entity.Translation, error) {
	ctx, span := otel.Tracer(_tracerName).Start(
		ctx, "TranslationWebAPI.Translate",
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(
			attribute.String("translation.source", translation.Source),
			attribute.String("translation.destination", translation.Destination),
		),
	)

	result, err := t.next.Translate(ctx, translation)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	span.End()

	return result, err
}
