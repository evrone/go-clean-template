//go:build system
// +build system

package app

import (
	"context"
	"github.com/evrone/go-clean-template/config"
	"github.com/evrone/go-clean-template/internal/test/db"
	"github.com/evrone/go-clean-template/pkg/logger"
	"github.com/evrone/go-clean-template/pkg/postgres"
	"github.com/evrone/go-clean-template/pkg/rabbitmq/rmq_rpc/client"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var httpEngine *gin.Engine
var cfg *config.Config

func init() {
	httpEngine, cfg = given()
}

func TestApp(t *testing.T) {

	t.Run("When calling the health endpoint, Then return 200", func(t *testing.T) {
		w := sendRequest("GET", "/healthz", httpEngine, nil)

		require.Equal(t, 200, w.Code)
		require.Equal(t, "", w.Body.String())
	})

	t.Run("When calling the do-translate endpoint providing all required information, Then return 200", func(t *testing.T) {
		body := `{
			"destination": "en",
			"original": "текст для перевода",
			"source": "auto"
		}`

		w := sendRequest("POST", "/v1/translation/do-translate", httpEngine, strings.NewReader(body))

		require.Equal(t, 200, w.Code)
		require.JSONEq(t, `{
			"source":"auto",
			"destination":"en",
			"original":"текст для перевода",
			"translation":"text to translate"
		}`, w.Body.String())
		require.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	})

	t.Run("When calling the do-translate endpoint missing source, Then return 400", func(t *testing.T) {
		body := `{
			"destination": "en",
			"original": "текст для перевода",
		}`

		w := sendRequest("POST", "/v1/translation/do-translate", httpEngine, strings.NewReader(body))

		require.Equal(t, 400, w.Code)
		require.JSONEq(t, `{"error":"invalid request body"}`, w.Body.String())
		require.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	})

	t.Run("When calling the history endpoint, Then return 200 containing history entries", func(t *testing.T) {

		w := sendRequest("GET", "/v1/translation/history", httpEngine, nil)

		require.Equal(t, 200, w.Code)
		require.Contains(t, w.Body.String(), `{"history":[{`)
	})

	t.Run("When calling the history endpoint using RabbitMQ RPC Client, Then returns history entries", func(t *testing.T) {

		rmqClient, err := client.New(cfg.RMQ.URL, cfg.RMQ.ServerExchange, cfg.RMQ.ClientExchange)
		if err != nil {
			panic(err)
		}
		if err != nil {
			t.Fatal("RabbitMQ RPC Client - init error - client.NewOrGetSingleton")
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

		for i := 0; i < 10; i++ {
			var history historyResponse

			err = rmqClient.RemoteCall("getHistory", nil, &history)
			if err != nil {
				t.Fatal("RabbitMQ RPC Client - remote call error - rmqClient.RemoteCall", err)
			}

			if history.History[0].Original != "текст для перевода" {
				t.Fatal("Original != текст для перевода")
			}
		}
	})
}

func given() (*gin.Engine, *config.Config) {
	ctx := context.Background()

	cfg := config.NewConfig()
	log := logger.New(cfg)

	db.MustStartPostgresContainer(ctx, cfg)
	db.MustStartRMQContainer(ctx, cfg)

	pg := setupPostgresClient(cfg)
	db.ExecuteMigrate(cfg.PG.URL, log)

	httpEngine := mustSetupHttpEngine(cfg, pg, log)

	return httpEngine, cfg
}

func mustSetupHttpEngine(config *config.Config, pg *postgres.Postgres, logger *logger.Logger) *gin.Engine {
	_, httpEngine := setupHttpEngine(config, logger)
	return httpEngine
}

func sendRequest(method string, url string, httpEngine *gin.Engine, body io.Reader) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, url, body)
	httpEngine.ServeHTTP(w, req)
	return w
}
