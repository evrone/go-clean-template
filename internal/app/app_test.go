package app

import (
	"context"
	"github.com/evrone/go-clean-template/config"
	"github.com/evrone/go-clean-template/internal/test/db"
	"github.com/evrone/go-clean-template/pkg/logger"
	"github.com/evrone/go-clean-template/pkg/postgres"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var httpEngine *gin.Engine

func init() {
	httpEngine = given()
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
	})

}

func given() *gin.Engine {
	ctx := context.Background()

	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}
	log := logger.New(cfg.Log.Level)

	db.MustStartPostgresContainer(err, ctx, cfg)
	db.MustStartRMQContainer(ctx, cfg)

	pg := setupPostgresClient(cfg)
	db.ExecuteMigrate(cfg.PG.URL)

	httpEngine := mustSetupHttpEngine(cfg, pg, log)

	return httpEngine
}

func mustSetupHttpEngine(config *config.Config, pg *postgres.Postgres, logger *logger.Logger) *gin.Engine {
	_, httpEngine := setupHttpEngine(config, pg, logger)
	return httpEngine
}

func sendRequest(method string, url string, httpEngine *gin.Engine, body io.Reader) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, url, body)
	httpEngine.ServeHTTP(w, req)
	return w
}
