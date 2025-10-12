package integration_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	protov1 "github.com/evrone/go-clean-template/docs/proto/v1"
	natsClient "github.com/evrone/go-clean-template/pkg/nats/nats_rpc/client"
	rmqClient "github.com/evrone/go-clean-template/pkg/rabbitmq/rmq_rpc/client"
	"github.com/goccy/go-json"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	// Base settings
	host     = "app"
	attempts = 20

	// Attempts connection
	httpURL        = "http://" + host + ":8080"
	healthPath     = httpURL + "/healthz"
	requestTimeout = 5 * time.Second

	// HTTP REST
	basePathV1 = httpURL + "/v1"

	// gRPC
	grpcURL = host + ":8081"

	// RPC configs
	rpcServerExchange = "rpc_server"
	rpcClientExchange = "rpc_client"
	requests          = 10

	// RabbitMQ RPC
	rmqURL = "amqp://guest:guest@rabbitmq:5672/"

	// RabbitMQ RPC
	natsURL = "nats://guest:guest@nats:4222/"

	// Test data
	expectedOriginal = "текст для перевода"
)

var errHealthCheck = fmt.Errorf("url %s is not available", healthPath)

func doWebRequestWithTimeout(ctx context.Context, method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	return http.DefaultClient.Do(req)
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

// HTTP POST: /v1/translation/do-translate.
func TestHTTPDoTranslateV1(t *testing.T) {
	tests := []struct {
		description string
		body        string
		expected    int
	}{
		{
			description: "DoTranslate Success",
			body: `{
				"destination": "en",
				"original": "текст для перевода",
				"source": "auto"
			}`,
			expected: http.StatusOK,
		},
		{
			description: "DoTranslate Success",
			body: `{
				"destination": "en",
				"original": "Текст для перевода",
				"source": "ru"
			}`,
			expected: http.StatusOK,
		},
		{
			description: "DoTranslate Fail",
			body: `{
				"destination": "en",
				"original": "текст для перевода"
			}`,
			expected: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			url := basePathV1 + "/translation/do-translate"
			ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)

			defer cancel()

			resp, err := doWebRequestWithTimeout(ctx, http.MethodPost, url, bytes.NewBuffer([]byte(tt.body)))
			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}

			defer resp.Body.Close()

			if resp.StatusCode != tt.expected {
				t.Errorf("Expected status %d, got %d", tt.expected, resp.StatusCode)
			}
		})
	}
}

// HTTP GET: /v1/translation/history.
func TestHTTPHistoryV1(t *testing.T) {
	url := basePathV1 + "/translation/history"
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)

	defer cancel()

	resp, err := doWebRequestWithTimeout(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var body struct {
		History []struct {
			Source      string `json:"source"`
			Destination string `json:"destination"`
			Original    string `json:"original"`
			Translation string `json:"translation"`
		} `json:"history"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	if len(body.History) == 0 {
		t.Error("Expected non-empty history")
	}
}

// gRPC Client V1: GetHistory.
func TestClientGRPCV1(t *testing.T) {
	grpcConn, err := grpc.NewClient(grpcURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal("gRPC Client - init error - grpc.NewClient", err)
	}

	defer func() {
		err = grpcConn.Close()
		if err != nil {
			t.Fatal("gRPC Client - shutdown error - grpcClientV1.GetHistory", err)
		}
	}()

	grpcClientV1 := protov1.NewTranslationClient(grpcConn)

	for i := 0; i < requests; i++ {
		history, err := grpcClientV1.GetHistory(t.Context(), &protov1.GetHistoryRequest{})
		if err != nil {
			t.Fatal("gRPC Client - remote call error - grpcClientV1.GetHistory", err)
		}

		if len(history.History) == 0 {
			t.Fatal("History slice is empty, expected at least one entry")
		}

		if history.History[0].Original != expectedOriginal {
			t.Fatalf("Original mismatch: expected %q, got %q", expectedOriginal, history.History[0].Original)
		}
	}
}

// RabbitMQ RPC Client V1: getHistory.
func TestClientRMQRPCV1(t *testing.T) { //nolint: dupl,gocritic,nolintlint
	client, err := rmqClient.New(rmqURL, rpcServerExchange, rpcClientExchange)
	if err != nil {
		t.Fatal("RabbitMQ RPC Client - init error - rmqClient.New", err)
	}

	defer func() {
		err = client.Shutdown()
		if err != nil {
			t.Fatal("RabbitMQ RPC Client - shutdown error - client.RemoteCall", err)
		}
	}()

	type Translation struct {
		Source      string `json:"source"`
		Destination string `json:"destination"`
		Original    string `json:"original"`
		Translation string `json:"translation"`
	}

	type historyResponse struct {
		History []Translation `json:"history"`
	}

	for i := 0; i < requests; i++ {
		var history historyResponse

		err = client.RemoteCall("v1.getHistory", nil, &history)
		if err != nil {
			t.Fatal("RabbitMQ RPC Client - remote call error - client.RemoteCall", err)
		}

		if len(history.History) == 0 {
			t.Fatal("History slice is empty, expected at least one entry")
		}

		if history.History[0].Original != expectedOriginal {
			t.Fatalf("Original mismatch: expected %q, got %q", expectedOriginal, history.History[0].Original)
		}
	}
}

// NATS RPC Client V1: getHistory.
func TestClientNATSRPCV1(t *testing.T) { //nolint: dupl,gocritic,nolintlint
	client, err := natsClient.New(natsURL, rpcServerExchange)
	if err != nil {
		t.Fatal("NATS RPC Client - init error - natsClient.New", err)
	}

	defer func() {
		err = client.Shutdown()
		if err != nil {
			t.Fatal("NATS RPC Client - shutdown error - rmqClient.RemoteCall", err)
		}
	}()

	type Translation struct {
		Source      string `json:"source"`
		Destination string `json:"destination"`
		Original    string `json:"original"`
		Translation string `json:"translation"`
	}

	type historyResponse struct {
		History []Translation `json:"history"`
	}

	for i := 0; i < requests; i++ {
		var history historyResponse

		err = client.RemoteCall("v1.getHistory", nil, &history)
		if err != nil {
			t.Fatal("NATS RPC Client - remote call error - rmqClient.RemoteCall", err)
		}

		if len(history.History) == 0 {
			t.Fatal("History slice is empty, expected at least one entry")
		}

		if history.History[0].Original != expectedOriginal {
			t.Fatalf("Original mismatch: expected %q, got %q", expectedOriginal, history.History[0].Original)
		}
	}
}
