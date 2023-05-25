package db

import (
	"context"
	"fmt"
	"github.com/evrone/go-clean-template/config"
	"github.com/evrone/go-clean-template/pkg/logger"
	"github.com/golang-migrate/migrate/v4"
	"github.com/testcontainers/testcontainers-go"
	postgres2 "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"path/filepath"
	"runtime"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	_ "github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file"
)

func MustStartPostgresContainer(ctx context.Context, cfg *config.Config) {
	container, err := postgres2.RunContainer(ctx,
		testcontainers.WithImage("postgres:15.2"),
		postgres2.WithDatabase("postgres"),
		postgres2.WithUsername("user"),
		postgres2.WithPassword("pass"),
		testcontainers.WithWaitStrategy(wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(5*time.Second)),
		testcontainers.CustomizeRequest(testcontainers.GenericContainerRequest{
			ContainerRequest: testcontainers.ContainerRequest{
				Name: "postgres-container",
			},
			Reuse: true,
		}),
	)

	if err != nil {
		panic(err)
	}

	host, err := container.Host(ctx)
	realPort, err := container.MappedPort(ctx, "5432")

	cfg.PG.URL = fmt.Sprintf("postgres://user:pass@%v:%v/postgres?sslmode=disable", host, realPort.Port())
}

func ExecuteMigrate(pgConnectionUrl string, log *logger.Logger) {
	projectRoot := projectRoot()

	migrationDirectoryUri := fmt.Sprintf("file://%s/migrations", projectRoot)
	m, err := migrate.New(
		migrationDirectoryUri,
		pgConnectionUrl,
	)

	if err != nil {
		panic(err)
	}
	if err := m.Up(); err != nil {
		// errors if no migration need to be executed
		log.Info(fmt.Sprintf("MIGRATE: %s", err))
	}
}

func projectRoot() string {
	_, b, _, _ := runtime.Caller(0)
	projectRoot := filepath.Dir(b)

	return projectRoot + "/../../../"
}
