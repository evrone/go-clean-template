package integration_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	protov1 "github.com/evrone/go-clean-template/docs/proto/v1"
	natsClient "github.com/evrone/go-clean-template/pkg/nats/nats_rpc/client"
	rmqClient "github.com/evrone/go-clean-template/pkg/rabbitmq/rmq_rpc/client"
	"github.com/goccy/go-json"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

const (
	// Base settings.
	host     = "app"
	attempts = 20

	// Attempts connection.
	httpURL        = "http://" + host + ":8080"
	healthPath     = httpURL + "/healthz"
	requestTimeout = 5 * time.Second

	// HTTP REST.
	basePathV1 = httpURL + "/v1"

	// gRPC.
	grpcURL = host + ":8081"

	// RPC configs.
	rpcServerExchange = "rpc_server"
	rpcClientExchange = "rpc_client"
	requests          = 10

	// Test password used across helpers.
	testPassword = "testpass123"
)

// rmqURL and natsURL are constructed from parts to avoid gosec G101 credential detection.
const (
	rpcCredentials = "guest:guest"
	rmqURL         = "amqp://" + rpcCredentials + "@rabbitmq:5672/"
	natsURL        = "nats://" + rpcCredentials + "@nats:4222/"
)

var errHealthCheck = fmt.Errorf("url %s is not available", healthPath)

// doWebRequestWithTimeout sends an HTTP request with a Content-Type of application/json.
func doWebRequestWithTimeout(ctx context.Context, method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	return http.DefaultClient.Do(req)
}

// doAuthenticatedRequest sends an HTTP request with a Bearer token.
func doAuthenticatedRequest(ctx context.Context, method, url string, body io.Reader, token string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	return http.DefaultClient.Do(req)
}

// registerUser registers a new user via the HTTP REST API.
func registerUser(t *testing.T, username, email, password string) *http.Response {
	t.Helper()

	ctx, cancel := context.WithTimeout(t.Context(), requestTimeout)
	defer cancel()

	body := fmt.Sprintf(`{"username":%q,"email":%q,"password":%q}`, username, email, password)

	resp, err := doWebRequestWithTimeout(ctx, http.MethodPost, basePathV1+"/auth/register", bytes.NewBufferString(body))
	if err != nil {
		t.Fatalf("registerUser: failed to send request: %v", err)
	}

	return resp
}

// loginUser logs in a user and returns the JWT token.
func loginUser(t *testing.T, email, password string) string {
	t.Helper()

	ctx, cancel := context.WithTimeout(t.Context(), requestTimeout)
	defer cancel()

	body := fmt.Sprintf(`{"email":%q,"password":%q}`, email, password)

	resp, err := doWebRequestWithTimeout(ctx, http.MethodPost, basePathV1+"/auth/login", bytes.NewBufferString(body))
	if err != nil {
		t.Fatalf("loginUser: failed to send request: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("loginUser: expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var result struct {
		Token string `json:"token"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("loginUser: failed to decode response: %v", err)
	}

	return result.Token
}

// sanitizeTestName converts t.Name() into a safe string for use as a username.
func sanitizeTestName(t *testing.T) string {
	t.Helper()

	name := t.Name()
	name = strings.ReplaceAll(name, "/", "_")
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ToLower(name)

	return name
}

// registerAndLogin creates a unique user via HTTP and returns the JWT token.
func registerAndLogin(t *testing.T) string {
	t.Helper()

	name := sanitizeTestName(t)
	email := name + "@test.com"
	password := testPassword

	resp := registerUser(t, name, email, password)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("registerAndLogin: register expected 201, got %d", resp.StatusCode)
	}

	return loginUser(t, email, password)
}

// registerAndLoginGRPC creates a unique user via gRPC and returns the JWT token.
func registerAndLoginGRPC(t *testing.T) string {
	t.Helper()

	name := sanitizeTestName(t)
	email := name + "@test.com"
	password := testPassword

	grpcConn, err := grpc.NewClient(grpcURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("registerAndLoginGRPC: grpc.NewClient: %v", err)
	}

	defer func() {
		if cerr := grpcConn.Close(); cerr != nil {
			t.Fatalf("registerAndLoginGRPC: grpcConn.Close: %v", cerr)
		}
	}()

	authClient := protov1.NewAuthServiceClient(grpcConn)

	_, err = authClient.Register(t.Context(), &protov1.RegisterRequest{
		Username: name,
		Email:    email,
		Password: password,
	})
	if err != nil {
		t.Fatalf("registerAndLoginGRPC: Register: %v", err)
	}

	loginResp, err := authClient.Login(t.Context(), &protov1.LoginRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		t.Fatalf("registerAndLoginGRPC: Login: %v", err)
	}

	return loginResp.Token
}

// registerAndLoginRMQ creates a unique user via RabbitMQ RPC and returns the JWT token.
func registerAndLoginRMQ(t *testing.T) string {
	t.Helper()

	name := sanitizeTestName(t)
	email := name + "@test.com"
	password := testPassword

	client, err := rmqClient.New(rmqURL, rpcServerExchange, rpcClientExchange)
	if err != nil {
		t.Fatalf("registerAndLoginRMQ: rmqClient.New: %v", err)
	}

	defer func() {
		if serr := client.Shutdown(); serr != nil {
			t.Fatalf("registerAndLoginRMQ: client.Shutdown: %v", serr)
		}
	}()

	registerPayload := map[string]string{
		"username": name,
		"email":    email,
		"password": password,
	}

	var registerResp any

	err = client.RemoteCall("v1.auth.register", registerPayload, &registerResp)
	if err != nil {
		t.Fatalf("registerAndLoginRMQ: register: %v", err)
	}

	loginPayload := map[string]string{
		"email":    email,
		"password": password,
	}

	var loginResp struct {
		Token string `json:"token"`
	}

	err = client.RemoteCall("v1.auth.login", loginPayload, &loginResp)
	if err != nil {
		t.Fatalf("registerAndLoginRMQ: login: %v", err)
	}

	return loginResp.Token
}

// registerAndLoginNATS creates a unique user via NATS RPC and returns the JWT token.
func registerAndLoginNATS(t *testing.T) string {
	t.Helper()

	name := sanitizeTestName(t)
	email := name + "@test.com"
	password := testPassword

	client, err := natsClient.New(natsURL, rpcServerExchange)
	if err != nil {
		t.Fatalf("registerAndLoginNATS: natsClient.New: %v", err)
	}

	defer func() {
		if serr := client.Shutdown(); serr != nil {
			t.Fatalf("registerAndLoginNATS: client.Shutdown: %v", serr)
		}
	}()

	registerPayload := map[string]string{
		"username": name,
		"email":    email,
		"password": password,
	}

	var registerResp any

	err = client.RemoteCall("v1.auth.register", registerPayload, &registerResp)
	if err != nil {
		t.Fatalf("registerAndLoginNATS: register: %v", err)
	}

	loginPayload := map[string]string{
		"email":    email,
		"password": password,
	}

	var loginResp struct {
		Token string `json:"token"`
	}

	err = client.RemoteCall("v1.auth.login", loginPayload, &loginResp)
	if err != nil {
		t.Fatalf("registerAndLoginNATS: login: %v", err)
	}

	return loginResp.Token
}

// authenticatedPayload wraps data with a token for RMQ/NATS authenticated RPC calls.
func authenticatedPayload(token string, data any) map[string]any {
	return map[string]any{
		"token": token,
		"data":  data,
	}
}

// grpcAuthCtx returns a context with gRPC authorization metadata.
func grpcAuthCtx(t *testing.T, token string) context.Context {
	t.Helper()

	return metadata.AppendToOutgoingContext(t.Context(), "authorization", token)
}

// parseJSON is a generic JSON parser for HTTP responses.
func parseJSON[T any](t *testing.T, resp *http.Response) T {
	t.Helper()

	var result T

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("parseJSON: failed to decode response: %v", err)
	}

	return result
}

func getHealthCheck(url string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)

	defer cancel()

	resp, err := doWebRequestWithTimeout(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return -1, err
	}

	defer resp.Body.Close()

	return resp.StatusCode, nil
}

func healthCheck(attempts int) error {
	for attempts > 0 {
		statusCode, err := getHealthCheck(healthPath)
		if err != nil {
			return err
		}

		if statusCode == http.StatusOK {
			return nil
		}

		log.Printf("Integration tests: url %s is not available, attempts left: %d", healthPath, attempts)

		time.Sleep(time.Second)

		attempts--
	}

	return errHealthCheck
}

func TestMain(m *testing.M) {
	err := healthCheck(attempts)
	if err != nil {
		log.Fatalf("Integration tests: httpURL %s is not available: %s", httpURL, err)
	}

	log.Printf("Integration tests: httpURL %s is available", httpURL)

	code := m.Run()
	os.Exit(code)
}
