package app

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/evrone/go-clean-template/config"
	"github.com/evrone/go-clean-template/pkg/logger"
	"github.com/evrone/go-clean-template/pkg/postgres"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	postgres2 "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
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

	ctx := context.Background()

	pg := startPostgresContainer(err, ctx, cfg, log)

	err = startRMQContainer(ctx, cfg)
	if err != nil {
		panic(err)
	}

	httpEngine := mustSetupHttpEngine(cfg, pg, log)

	return httpEngine
}

func startPostgresContainer(err error, ctx context.Context, cfg *config.Config, log *logger.Logger) *postgres.Postgres {
	container, err := postgres2.RunContainer(ctx,
		testcontainers.WithImage("postgres:15.2"),
		postgres2.WithDatabase("postgres"),
		postgres2.WithUsername("user"),
		postgres2.WithPassword("pass"),
		testcontainers.WithWaitStrategy(wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		panic(err)
	}

	connStr, err := container.ConnectionString(ctx, "sslmode=disable", "application_name=test")
	cfg.PG.URL = connStr

	pg := setupPostgresClient(cfg, log)
	return pg
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

// startRMQContainer ...
func startRMQContainer(ctx context.Context, cfg *config.Config) error {

	port, err := nat.NewPort("", "5672")
	if err != nil {
		return fmt.Errorf("failed to build port: %v", err)
	}

	timeout := 5 * time.Minute // Default timeout
	tag := "3.11.15"

	req := testcontainers.ContainerRequest{
		Image:        fmt.Sprintf("rabbitmq:%s", tag),
		ExposedPorts: []string{string(port)},
		WaitingFor:   wait.ForListeningPort(port).WithStartupTimeout(timeout),
		// WaitingFor:   wait.ForLog("Server startup complete").WithStartupTimeout(timeout),
	}

	//tc.MergeRequest(&req, &options.ContainerOptions.ContainerRequest)

	rmqContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return fmt.Errorf("failed to start container: %v", err)
	}

	host, err := rmqContainer.Host(ctx)
	if err != nil {
		return fmt.Errorf("failed to get container host: %v", err)
	}

	realPort, err := rmqContainer.MappedPort(ctx, port)
	if err != nil {
		return fmt.Errorf("failed to get exposed container port: %v", err)
	}

	cfg.RMQ.URL = fmt.Sprintf("amqp://%s:%s", host, realPort.Port())

	return nil
}
