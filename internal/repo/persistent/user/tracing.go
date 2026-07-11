package user

import (
	"context"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const _tracerName = "github.com/evrone/go-clean-template/internal/repo/persistent/user"

// tracedRepo wraps a UserRepo with OpenTelemetry spans.
type tracedRepo struct {
	next repo.UserRepo
}

func newTraced(next repo.UserRepo) repo.UserRepo {
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

func (r *tracedRepo) Store(ctx context.Context, user *entity.User) error {
	ctx, span := startSpan(ctx, "UserRepo.Store", attribute.String("user.id", user.ID))

	err := r.next.Store(ctx, user)
	endSpan(span, err)

	return err
}

func (r *tracedRepo) GetByID(ctx context.Context, id string) (entity.User, error) {
	ctx, span := startSpan(ctx, "UserRepo.GetByID", attribute.String("user.id", id))

	result, err := r.next.GetByID(ctx, id)
	endSpan(span, err)

	return result, err
}

func (r *tracedRepo) GetByEmail(ctx context.Context, email string) (entity.User, error) {
	ctx, span := startSpan(ctx, "UserRepo.GetByEmail", attribute.String("user.email", email))

	result, err := r.next.GetByEmail(ctx, email)
	endSpan(span, err)

	return result, err
}
