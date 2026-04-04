package v1

import (
	"context"
	"errors"

	v1 "github.com/evrone/go-clean-template/docs/proto/v1"
	grpcmw "github.com/evrone/go-clean-template/internal/controller/grpc/middleware"
	"github.com/evrone/go-clean-template/internal/controller/grpc/v1/response"
	"github.com/evrone/go-clean-template/internal/entity"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Register -.
func (c *AuthController) Register(ctx context.Context, req *v1.RegisterRequest) (*v1.RegisterResponse, error) {
	user, err := c.u.Register(ctx, req.GetUsername(), req.GetEmail(), req.GetPassword())
	if err != nil {
		c.l.Error(err, "grpc - v1 - Register")

		if errors.Is(err, entity.ErrUserAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}

		return nil, status.Error(codes.Internal, "internal server error")
	}

	return response.NewRegisterResponse(&user), nil
}

// Login -.
func (c *AuthController) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginResponse, error) {
	token, err := c.u.Login(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		c.l.Error(err, "grpc - v1 - Login")

		if errors.Is(err, entity.ErrInvalidCredentials) {
			return nil, status.Error(codes.Unauthenticated, "invalid credentials")
		}

		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &v1.LoginResponse{Token: token}, nil
}

// GetProfile -.
func (c *AuthController) GetProfile(ctx context.Context, _ *v1.GetProfileRequest) (*v1.GetProfileResponse, error) {
	userID, ok := grpcmw.UserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	user, err := c.u.GetUser(ctx, userID)
	if err != nil {
		c.l.Error(err, "grpc - v1 - GetProfile")

		if errors.Is(err, entity.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}

		return nil, status.Error(codes.Internal, "internal server error")
	}

	return response.NewGetProfileResponse(&user), nil
}
