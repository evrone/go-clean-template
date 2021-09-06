package integration_test

import (
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	. "github.com/Eun/go-hit"

	"github.com/evrone/go-clean-template/pkg/rabbitmq/rmq_rpc/client"
)

const (
	// Attempts connection
	host       = "app:8080"
	healthPath = "http://" + host + "/healthz"
	attempts   = 20

	// HTTP REST
	basePath = "http://" + host + "/v1"

	// RabbitMQ RPC
	rmqURL            = "amqp://guest:guest@rabbitmq:5672/"
	rpcServerExchange = "rpc_server"
	rpcClientExchange = "rpc_client"
	requests          = 10
)

func TestMain(m *testing.M) {
	err := healthCheck(attempts)
	if err != nil {
		log.Fatalf("Integration tests: host %s is not available: %s", host, err)
	}

	log.Printf("Integration tests: host %s is available", host)

	code := m.Run()
	os.Exit(code)
}

func healthCheck(attempts int) error {
	var err error

	for attempts > 0 {
		err = Do(Get(healthPath), Expect().Status().Equal(http.StatusOK))
		if err == nil {
			return nil
		}

		log.Printf("Integration tests: url %s is not available, attempts left: %d", healthPath, attempts)

		time.Sleep(time.Second)

		attempts--
	}

	return err
}

// HTTP POST: /translation/do-translate.
func TestHTTPDoTranslate(t *testing.T) {
	body := `{
		"destination": "en",
		"original": "текст для перевода",
		"source": "auto"
	}`
	Test(t,
		Description("DoTranslate Success"),
		Post(basePath+"/translation/do-translate"),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Body().String(body),
		Expect().Status().Equal(http.StatusOK),
		Expect().Body().JSON().JQ(".translation").Equal("text for translation"),
	)

	body = `{
		"destination": "en",
		"original": "текст для перевода"
	}`
	Test(t,
		Description("DoTranslate Fail"),
		Post(basePath+"/translation/do-translate"),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Body().String(body),
		Expect().Status().Equal(http.StatusBadRequest),
		Expect().Body().JSON().JQ(".error").Equal("invalid request body"),
	)
}

// HTTP GET: /translation/history.
func TestHTTPHistory(t *testing.T) {
	Test(t,
		Description("History Success"),
		Get(basePath+"/translation/history"),
		Expect().Status().Equal(http.StatusOK),
		Expect().Body().String().Contains(`{"history":[{`),
	)
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
