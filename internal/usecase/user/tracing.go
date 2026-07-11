package user

import (
	"context"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/usecase"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const _tracerName = "github.com/evrone/go-clean-template/internal/usecase/user"

// tracedUseCase wraps a User usecase with OpenTelemetry spans, closing the
// gap between transport spans (HTTP/gRPC/AMQP/NATS) and repository spans.
type tracedUseCase struct {
	next usecase.User
}

// newTraced wraps a User usecase with tracing spans.
func newTraced(next usecase.User) usecase.User {
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

func (u *tracedUseCase) Register(ctx context.Context, username, email, password string) (entity.User, error) {
	ctx, span := startSpan(ctx, "UserUseCase.Register", attribute.String("user.email", email))

	result, err := u.next.Register(ctx, username, email, password)
	endSpan(span, err)

	return result, err
}

func (u *tracedUseCase) Login(ctx context.Context, email, password string) (string, error) {
	ctx, span := startSpan(ctx, "UserUseCase.Login", attribute.String("user.email", email))

	result, err := u.next.Login(ctx, email, password)
	endSpan(span, err)

	return result, err
}

func (u *tracedUseCase) GetUser(ctx context.Context, userID string) (entity.User, error) {
	ctx, span := startSpan(ctx, "UserUseCase.GetUser", attribute.String("user.id", userID))

	result, err := u.next.GetUser(ctx, userID)
	endSpan(span, err)

	return result, err
}
