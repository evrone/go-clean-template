package postgres

import (
	"context"
	"time"

	"github.com/evrone/go-service-template/internal/business-logic/domain"

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
		domain.Logger.Fatal(err, "postgres connect error")
	}

	poolConfig.MaxConns = int32(maxPoolSize)

	var errConn error
	var pool *pgxpool.Pool

	for connAttempts > 0 {
		pool, errConn = pgxpool.ConnectConfig(context.Background(), poolConfig)
		if errConn == nil {
			break
		}
		domain.Logger.Debug("postgres is trying to connect",
			domain.Field{Key: "attempts left", Val: connAttempts},
		)

		time.Sleep(time.Second * 2)

		connAttempts--
	}

	if errConn != nil {
		domain.Logger.Fatal(errConn, "postgres connect error")
	}

	domain.Logger.Info("postgres connected")

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
