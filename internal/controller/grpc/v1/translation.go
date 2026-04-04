package v1

import (
	"context"

	v1 "github.com/evrone/go-clean-template/docs/proto/v1"
	grpcmw "github.com/evrone/go-clean-template/internal/controller/grpc/middleware"
	"github.com/evrone/go-clean-template/internal/controller/grpc/v1/response"
	"github.com/evrone/go-clean-template/internal/entity"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetHistory -.
func (c *TranslationController) GetHistory(ctx context.Context, _ *v1.GetHistoryRequest) (*v1.GetHistoryResponse, error) {
	userID, ok := grpcmw.UserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	translationHistory, err := c.t.History(ctx, userID)
	if err != nil {
		c.l.Error(err, "grpc - v1 - GetHistory")

		return nil, status.Error(codes.Internal, "internal server error")
	}

	return response.NewTranslationHistory(translationHistory), nil
}

// Translate -.
func (c *TranslationController) Translate(ctx context.Context, req *v1.TranslateRequest) (*v1.TranslateResponse, error) {
	userID, ok := grpcmw.UserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	t, err := c.t.Translate(ctx, userID, entity.Translation{
		Source:      req.GetSource(),
		Destination: req.GetDestination(),
		Original:    req.GetOriginal(),
	})
	if err != nil {
		c.l.Error(err, "grpc - v1 - Translate")

		return nil, status.Error(codes.Internal, "translation service problems")
	}

	return &v1.TranslateResponse{
		Source:      t.Source,
		Destination: t.Destination,
		Original:    t.Original,
		Translation: t.Translation,
	}, nil
}
