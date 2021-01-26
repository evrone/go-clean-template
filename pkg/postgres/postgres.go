package postgres

import (
	"context"
	"log"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Postgres struct {
	Pool    *pgxpool.Pool
	Builder squirrel.StatementBuilderType
}

func NewPostgres(url string, maxPoolSize int) Postgres {
	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Fatalf("postgres connect error: %s", err)
	}

	poolConfig.MaxConns = int32(maxPoolSize)

	pool, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatalf("postgres connect error: %s", err)
	}

	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	return Postgres{
		Pool:    pool,
		Builder: builder,
	}
}

func (p Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
