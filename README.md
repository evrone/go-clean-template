# Go Clean Template — Monorepo

Clean Architecture monorepo template for Golang microservices.

[![License](https://img.shields.io/badge/License-MIT-success)](https://github.com/evrone/go-clean-template/blob/master/LICENSE)
[![Web Framework](https://img.shields.io/badge/Fiber-Web%20Framework-blue)](https://github.com/gofiber/fiber)
[![API Documentation](https://img.shields.io/badge/Swagger-API%20Documentation-blue)](https://github.com/swaggo/swag)
[![Validation](https://img.shields.io/badge/Validator-Data%20Integrity-blue)](https://github.com/go-playground/validator)
[![Query Builder](https://img.shields.io/badge/Squirrel-SQL%20Query%20Builder-blue)](https://github.com/Masterminds/squirrel)
[![Database Migrations](https://img.shields.io/badge/Migrations-Seamless%20Schema%20Updates-blue)](https://github.com/golang-migrate/migrate)
[![Logging](https://img.shields.io/badge/ZeroLog-Structured%20Logging-blue)](https://github.com/rs/zerolog)
[![Metrics](https://img.shields.io/badge/Prometheus-Metrics%20Integration-blue)](https://github.com/ansrivas/fiberprometheus)
[![Testing](https://img.shields.io/badge/Testify-Testing%20Framework-blue)](https://github.com/stretchr/testify)
[![Mocking](https://img.shields.io/badge/Mock-Mocking%20Library-blue)](https://go.uber.org/mock)

## Overview

This template demonstrates how to:

- Organize a Go monorepo with multiple independent microservices
- Share common packages (logger, httpserver, postgres, middleware) across services without duplication
- Keep business logic independent, clean, and extensible using Clean Architecture (Uncle Bob)
- Prevent services from turning into spaghetti code as they grow

---

## Monorepo Structure

![Module Dependencies](docs/img/diagram-1.png)

```
.
├── go.work                        # Go workspace — ties all modules together
├── go.work.sum
├── common_packages/
│   ├── pkg/                       # evrone.local/common-pkg
│   │   ├── logger/                # Zerolog wrapper
│   │   ├── httpserver/            # Fiber HTTP server with graceful shutdown
│   │   ├── postgres/              # pgx connection pool
│   │   ├── grpcserver/            # gRPC server
│   │   ├── rabbitmq/              # RabbitMQ RPC (server + client)
│   │   └── nats/                  # NATS RPC (server + client)
│   └── middleware/                # evrone.local/common-middleware
│       ├── logger.go              # Request logging middleware (Fiber)
│       └── recovery.go            # Panic recovery middleware (Fiber)
├── services/
│   ├── translation/               # github.com/evrone/translation-svc
│   │   ├── cmd/app/               # Entry point
│   │   ├── config/                # Env-based config
│   │   ├── docs/                  # Swagger + proto generated files
│   │   ├── integration-test/      # Integration tests (run in Docker)
│   │   ├── migrations/            # SQL migrations
│   │   └── internal/
│   │       ├── app/               # Bootstrap (app.go + migrate.go)
│   │       ├── controller/        # REST, gRPC, AMQP RPC, NATS RPC
│   │       ├── entity/            # Business entities
│   │       ├── repo/              # Repository interfaces + implementations
│   │       └── usecase/           # Business logic
│   └── todo/                      # github.com/evrone/todo-svc
│       ├── cmd/todo/              # Entry point
│       ├── config/                # Env-based config
│       ├── migrations/            # SQL migrations
│       └── internal/
│           ├── app/               # Bootstrap (app.go + migrate.go)
│           ├── controller/restapi # REST API only
│           ├── entity/            # Business entities
│           ├── repo/              # Repository interfaces + implementations
│           └── usecase/           # Business logic
├── docker-compose.yml             # Full stack
├── docker-compose-integration-test.yml
├── nginx/nginx.conf               # Reverse proxy for all services
└── Makefile
```

---

## Quick Start

### 1. Install dev tools

```sh
make bin-deps
```

### 2. Start infrastructure (PostgreSQL, RabbitMQ, NATS)

```sh
make compose-up
```

### 3. Run services locally (each in a separate terminal)

```sh
# Translation service — http://localhost:8080, gRPC :8081
make run-translation

# Todo service — http://localhost:8082
make run-todo
```

> Both commands auto-run migrations on startup via `-tags migrate`.

### 4. Run full stack in Docker

```sh
make compose-up-all
```

### 5. Integration tests (CI-ready)

```sh
make compose-up-integration-test
```

---

## Service Endpoints

![Infrastructure Overview](docs/img/diagram-3.png)

### Translation service

| Type       | URL / Address                                      |
|------------|-----------------------------------------------------|
| REST API   | http://app.lvh.me or http://localhost:8080          |
| Health     | http://localhost:8080/healthz                       |
| Metrics    | http://localhost:8080/metrics                       |
| Swagger    | http://localhost:8080/swagger                       |
| gRPC       | `tcp://grpc.lvh.me:8081` / `localhost:8081`         |
| AMQP RPC   | `amqp://guest:guest@localhost:5672/`                |
| NATS RPC   | `nats://guest:guest@localhost:4222/`                |

### Todo service

| Type       | URL / Address                                      |
|------------|-----------------------------------------------------|
| REST API   | http://todo.lvh.me or http://localhost:8082          |
| Health     | http://localhost:8082/healthz                       |

### Infrastructure

| Service    | URL                                                 |
|------------|-----------------------------------------------------|
| PostgreSQL | `postgres://user:***@localhost:5432/db`             |
| RabbitMQ   | http://rabbitmq.lvh.me / http://localhost:15672     |
| NATS       | http://nats.lvh.me / http://localhost:8222          |

---

## Configuration

Each service is configured via environment variables following the [12-factor app](https://12factor.net/) principle.

```
services/translation/.env.example   ← copy to .env and fill in values
services/todo/.env.example          ← copy to .env and fill in values
```

`.env` files are gitignored — only `.env.example` is committed.

---

## Migrations

Migrations run automatically on service startup when built with `-tags migrate` (which `make run-*` does).

To run manually:

```sh
make migrate-up-translation
make migrate-up-todo
```

To create a new migration:

```sh
make migrate-create-todo NAME=add_priority_index
make migrate-create-translation NAME=add_column
```

---

## Make Targets

```sh
make help          # List all targets with descriptions
```

| Category      | Target                       | Description                              |
|---------------|------------------------------|------------------------------------------|
| Docker        | `compose-up`                 | Start infra only (db, rabbitmq, nats)    |
|               | `compose-up-all`             | Start everything including services      |
|               | `compose-down`               | Stop all containers                      |
| Run local     | `run-translation`            | Run translation service with migrations  |
|               | `run-todo`                   | Run todo service with migrations         |
| Migrations    | `migrate-up-translation`     | Apply translation DB migrations          |
|               | `migrate-up-todo`            | Apply todo DB migrations                 |
| Tests         | `test`                       | Run all unit tests                       |
|               | `test-translation`           | Run translation unit tests               |
|               | `test-todo`                  | Run todo unit tests                      |
|               | `integration-test`           | Run translation integration tests        |
| Code gen      | `swag-translation`           | Regenerate Swagger docs                  |
|               | `proto-translation`          | Regenerate gRPC stubs                    |
|               | `mock`                       | Regenerate all mocks                     |
| Deps          | `deps`                       | `go mod tidy` for all 4 modules          |

---

## Clean Architecture

![Clean Architecture Layers](docs/img/diagram-2.png)

Business logic is kept independent from frameworks, databases, and transports.

```
Controller (HTTP/gRPC/RMQ) → UseCase → Repository (Postgres)
```

Dependencies only point **inward**:
- Controllers call UseCases via interfaces
- UseCases call Repositories via interfaces
- Entities have no external dependencies

### Layers

```
internal/
├── entity/          # Pure Go structs — no external deps
├── usecase/         # Business logic — depends only on interfaces
│   ├── contracts.go # UseCase interfaces
│   └── {domain}/    # Implementation
├── repo/
│   ├── contracts.go         # Repository interfaces
│   └── persistent/          # Postgres implementations
└── controller/
    └── restapi/             # HTTP handlers — depend on UseCase interfaces
```

### Dependency Injection

```go
// UseCase only knows about its interface — not Postgres, not HTTP
type TodoRepo interface {
    Create(ctx context.Context, t entity.Todo) (entity.Todo, error)
    GetByID(ctx context.Context, id int) (entity.Todo, error)
    List(ctx context.Context) ([]entity.Todo, error)
    Update(ctx context.Context, id int, t entity.Todo) (entity.Todo, error)
    Delete(ctx context.Context, id int) error
}
```

---

## Adding a New Service

1. Create `services/{name}/` with the same structure as `services/todo/`
2. Add a `go.mod` with module path `github.com/evrone/{name}-svc`
3. Add `use ./services/{name}` to `go.work`
4. Add a `Dockerfile` at `services/{name}/Dockerfile` (build context = repo root)
5. Add the service to `docker-compose.yml` and `nginx/nginx.conf`

---

## Useful Links

- [The Clean Architecture — Uncle Bob](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Twelve-factor app](https://12factor.net/)
- [Go Workspaces](https://go.dev/blog/get-familiar-with-workspaces)
- [Fiber](https://github.com/gofiber/fiber)
- [golang-migrate](https://github.com/golang-migrate/migrate)
