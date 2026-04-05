package middleware_test

import (
	"context"
	"testing"
	"time"

	grpcmw "github.com/evrone/go-clean-template/internal/controller/grpc/middleware"
	"github.com/evrone/go-clean-template/pkg/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const respOK = "ok"

type ctxCapture struct {
	ctx context.Context
}

func (c *ctxCapture) handler(_ context.Context, _ any) (any, error) {
	return respOK, nil
}

func (c *ctxCapture) capturingHandler(ctx context.Context, _ any) (any, error) {
	c.ctx = ctx

	return respOK, nil
}

func newJWTManager(t *testing.T) *jwt.Manager {
	t.Helper()

	return jwt.New("test-secret", time.Hour)
}

func runSkipAuthTest(t *testing.T, method string) {
	t.Helper()

	jwtMgr := newJWTManager(t)
	interceptor := grpcmw.AuthInterceptor(jwtMgr)
	info := &grpc.UnaryServerInfo{FullMethod: method}

	called := false
	handler := func(_ context.Context, _ any) (any, error) {
		called = true

		return respOK, nil
	}

	resp, err := interceptor(t.Context(), nil, info, handler)

	require.NoError(t, err)
	assert.Equal(t, respOK, resp)
	assert.True(t, called)
}

func TestAuthInterceptor_SkipRegister(t *testing.T) {
	t.Parallel()
	runSkipAuthTest(t, "/grpc.v1.AuthService/Register")
}

func TestAuthInterceptor_SkipLogin(t *testing.T) {
	t.Parallel()
	runSkipAuthTest(t, "/grpc.v1.AuthService/Login")
}

func TestAuthInterceptor_MissingMetadata(t *testing.T) {
	t.Parallel()

	jwtMgr := newJWTManager(t)
	interceptor := grpcmw.AuthInterceptor(jwtMgr)
	info := &grpc.UnaryServerInfo{FullMethod: "/grpc.v1.TaskService/GetTask"}

	capture := &ctxCapture{}

	resp, err := interceptor(t.Context(), nil, info, capture.handler)

	assert.Nil(t, resp)
	require.Error(t, err)

	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.Unauthenticated, st.Code())
	assert.Contains(t, st.Message(), "missing metadata")
}

func TestAuthInterceptor_MissingAuthorizationToken(t *testing.T) {
	t.Parallel()

	jwtMgr := newJWTManager(t)
	interceptor := grpcmw.AuthInterceptor(jwtMgr)
	info := &grpc.UnaryServerInfo{FullMethod: "/grpc.v1.TaskService/GetTask"}

	md := metadata.New(map[string]string{"other-key": "value"})
	ctx := metadata.NewIncomingContext(t.Context(), md)

	capture := &ctxCapture{}

	resp, err := interceptor(ctx, nil, info, capture.handler)

	assert.Nil(t, resp)
	require.Error(t, err)

	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.Unauthenticated, st.Code())
	assert.Contains(t, st.Message(), "missing authorization token")
}

func TestAuthInterceptor_InvalidToken(t *testing.T) {
	t.Parallel()

	jwtMgr := newJWTManager(t)
	interceptor := grpcmw.AuthInterceptor(jwtMgr)
	info := &grpc.UnaryServerInfo{FullMethod: "/grpc.v1.TaskService/GetTask"}

	md := metadata.Pairs("authorization", "invalid-token")
	ctx := metadata.NewIncomingContext(t.Context(), md)

	capture := &ctxCapture{}

	resp, err := interceptor(ctx, nil, info, capture.handler)

	assert.Nil(t, resp)
	require.Error(t, err)

	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.Unauthenticated, st.Code())
	assert.Contains(t, st.Message(), "invalid or expired token")
}

func TestAuthInterceptor_ValidToken(t *testing.T) {
	t.Parallel()

	jwtMgr := newJWTManager(t)
	interceptor := grpcmw.AuthInterceptor(jwtMgr)
	info := &grpc.UnaryServerInfo{FullMethod: "/grpc.v1.TaskService/GetTask"}

	token, err := jwtMgr.GenerateToken("user-id-123")
	require.NoError(t, err)

	md := metadata.Pairs("authorization", token)
	ctx := metadata.NewIncomingContext(t.Context(), md)

	capture := &ctxCapture{}

	resp, err := interceptor(ctx, nil, info, capture.capturingHandler)

	require.NoError(t, err)
	assert.Equal(t, respOK, resp)

	userID, ok := grpcmw.UserIDFromContext(capture.ctx)
	assert.True(t, ok)
	assert.Equal(t, "user-id-123", userID)
}

func TestUserIDFromContext_WithValue(t *testing.T) {
	t.Parallel()

	jwtMgr := newJWTManager(t)
	interceptor := grpcmw.AuthInterceptor(jwtMgr)
	info := &grpc.UnaryServerInfo{FullMethod: "/grpc.v1.TaskService/GetTask"}

	token, err := jwtMgr.GenerateToken("user-42")
	require.NoError(t, err)

	md := metadata.Pairs("authorization", token)
	ctx := metadata.NewIncomingContext(t.Context(), md)

	capture := &ctxCapture{}

	_, err = interceptor(ctx, nil, info, capture.capturingHandler)
	require.NoError(t, err)

	userID, ok := grpcmw.UserIDFromContext(capture.ctx)
	assert.True(t, ok)
	assert.Equal(t, "user-42", userID)
}

func TestUserIDFromContext_WithoutValue(t *testing.T) {
	t.Parallel()

	userID, ok := grpcmw.UserIDFromContext(t.Context())
	assert.False(t, ok)
	assert.Empty(t, userID)
}
