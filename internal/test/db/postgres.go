package db

import (
	"context"
	"fmt"
	"github.com/evrone/go-clean-template/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/testcontainers/testcontainers-go"
	postgres2 "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"os"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	_ "github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file"
)

func MustStartPostgresContainer(err error, ctx context.Context, cfg *config.Config) {
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

	host, err := container.Host(ctx)
	realPort, err := container.MappedPort(ctx, "5432")

	cfg.PG.URL = fmt.Sprintf("postgres://user:pass@%v:%v/postgres?sslmode=disable", host, realPort.Port())

	//connStr, err := container.ConnectionString(ctx, "sslmode=disable", "application_name=test")
	//cfg.PG.URL = connStr

}

func ExecuteMigrate(pgConnectionUrl string) {
	cwd := mustGetCwd()

	migrationDirectoryUri := fmt.Sprintf("file://%s/migrations", cwd)
	m, err := migrate.New(
		migrationDirectoryUri,
		pgConnectionUrl,
	)

	if err != nil {
		panic(err)
	}
	if err := m.Up(); err != nil {
		panic(err)
	}
}

func mustGetCwd() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return cwd
}
