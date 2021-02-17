// Package postgres implements postgres connection.
package postgres

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/evrone/go-service-template/pkg/logger"
)

type Postgres struct {
	Pool    *pgxpool.Pool
	Builder squirrel.StatementBuilderType
}

func NewPostgres(url string, maxPoolSize, connAttempts int) *Postgres {
	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		logger.Fatal(err, "postgres connect error")
	}

	poolConfig.MaxConns = int32(maxPoolSize)

	var (
		errConn error
		pool    *pgxpool.Pool
	)

	for connAttempts > 0 {
		pool, errConn = pgxpool.ConnectConfig(context.Background(), poolConfig)
		if errConn == nil {
			break
		}

		logger.Info("postgres is trying to connect",
			logger.Field{Key: "attempts left", Val: connAttempts},
		)

		time.Sleep(time.Second)

		connAttempts--
	}

	if errConn != nil {
		logger.Fatal(errConn, "postgres connect error")
	}

	logger.Info("postgres connected")

	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	return &Postgres{
		Pool:    pool,
		Builder: builder,
	}
}

func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
