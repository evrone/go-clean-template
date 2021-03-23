// Package postgres implements postgres connection.
package postgres

import (
	"context"
	"log"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

type Postgres struct {
	maxPoolSize  int
	connAttempts int
	Builder      squirrel.StatementBuilderType
	Pool         *pgxpool.Pool
}

func NewPostgres(url string, opts ...Option) (*Postgres, error) {
	pg := &Postgres{
		// Default
		maxPoolSize:  1,
		connAttempts: 10,
	}

	// Set options
	for _, opt := range opts {
		opt(pg)
	}

	pg.Builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, errors.Wrap(err, "postgres - NewPostgres - pgxpool.ParseConfig")
	}

	poolConfig.MaxConns = int32(pg.maxPoolSize)

	for pg.connAttempts > 0 {
		pg.Pool, err = pgxpool.ConnectConfig(context.Background(), poolConfig)
		if err == nil {
			break
		}

		log.Printf("Postgres is trying to connect, attempts left: %d", pg.connAttempts)

		time.Sleep(time.Second)

		pg.connAttempts--
	}

	if err != nil {
		return nil, errors.Wrap(err, "postgres - NewPostgres - connAttempts == 0")
	}

	return pg, nil
}

func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
