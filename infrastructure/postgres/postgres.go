package postgres

import (
	"context"
	"log"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Postgres struct {
	Pool    *pgxpool.Pool
	Builder squirrel.StatementBuilderType
}

func NewPostgres(url string, maxPoolSize, connAttempts int) Postgres {
	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Fatalf("postgres connect error: %s", err)
	}

	poolConfig.MaxConns = int32(maxPoolSize)

	var errConn error
	var pool *pgxpool.Pool

	for connAttempts > 0 {
		pool, errConn = pgxpool.ConnectConfig(context.Background(), poolConfig)
		if errConn == nil {
			break
		}

		log.Printf("postgres is trying to connect, attempts left: %d", connAttempts)

		time.Sleep(time.Second)

		connAttempts--
	}

	if errConn != nil {
		log.Fatalf("postgres connect error: %s", errConn)
	}

	log.Print("postgres connected")

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
