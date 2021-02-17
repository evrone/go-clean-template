// +build migrate

package app

import (
	"errors"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func init() {
	var err error
	var m *migrate.Migrate

	attempts := 10
	for attempts > 0 {
		m, err = migrate.New("file://migrations",
			"postgres://user:pass@db:5432/postgres?sslmode=disable")
		if err == nil {
			break
		}
		log.Printf("Migrate: postgres is trying to connect, attempts left: %d", attempts)
		time.Sleep(time.Second)
		attempts--
	}

	if err != nil {
		log.Fatalf("Migrate: postgres connect error: %s", err)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Migrate: up error: %s", err)
	}
	m.Close()
}
