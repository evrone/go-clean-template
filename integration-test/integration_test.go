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

	"github.com/evrone/go-clean-template/pkg/rabbitmq/rmq_rpc/client"
	"github.com/goccy/go-json"
)

const (
	// Attempts connection
	host           = "app:8080"
	healthPath     = "http://" + host + "/healthz"
	attempts       = 20
	requestTimeout = 5 * time.Second

	// HTTP REST
	basePath = "http://" + host + "/v1"

	// RabbitMQ RPC
	rmqURL            = "amqp://guest:guest@rabbitmq:5672/"
	rpcServerExchange = "rpc_server"
	rpcClientExchange = "rpc_client"
	requests          = 10
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
		log.Fatalf("Integration tests: host %s is not available: %s", host, err)
	}

	log.Printf("Integration tests: host %s is available", host)

	code := m.Run()
	os.Exit(code)
}

// HTTP POST: /translation/do-translate.
func TestHTTPDoTranslate(t *testing.T) {
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
			url := basePath + "/translation/do-translate"
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

// HTTP GET: /translation/history.
func TestHTTPHistory(t *testing.T) {
	url := basePath + "/translation/history"
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

// RabbitMQ RPC Client: getHistory.
func TestRMQClientRPC(t *testing.T) {
	rmqClient, err := client.New(rmqURL, rpcServerExchange, rpcClientExchange)
	if err != nil {
		t.Fatal("RabbitMQ RPC Client - init error - client.New")
	}

	defer func() {
		err = rmqClient.Shutdown()
		if err != nil {
			t.Fatal("RabbitMQ RPC Client - shutdown error - rmqClient.RemoteCall", err)
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

		err = rmqClient.RemoteCall("getHistory", nil, &history)
		if err != nil {
			t.Fatal("RabbitMQ RPC Client - remote call error - rmqClient.RemoteCall", err)
		}

		if history.History[0].Original != "текст для перевода" {
			t.Fatal("Original != текст для перевода")
		}
	}
}
