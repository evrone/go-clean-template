package app

import (
	"github.com/evrone/go-clean-template/config"
	"github.com/evrone/go-clean-template/pkg/logger"
	"github.com/evrone/go-clean-template/pkg/postgres"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestApp(t *testing.T) {

	t.Run("When calling the health endpoint, Then return 200", func(t *testing.T) {
		httpEngine := given()

		w := sendRequest("GET", "/healthz", httpEngine)

		require.Equal(t, 200, w.Code)
		require.Equal(t, "", w.Body.String())
	})

}

func given() *gin.Engine {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	log := logger.New(cfg.Log.Level)
	pg := setupPostgresClient(cfg, log)

	httpEngine := mustSetupHttpEngine(cfg, pg, log)

	return httpEngine
}

func mustSetupHttpEngine(config *config.Config, pg *postgres.Postgres, logger *logger.Logger) *gin.Engine {
	_, err, httpEngine := setupHttpEngine(config, pg, logger)
	if err != nil {
		panic(err)
	}
	return httpEngine
}

func sendRequest(method string, url string, httpEngine *gin.Engine) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, url, nil)
	httpEngine.ServeHTTP(w, req)
	return w
}
