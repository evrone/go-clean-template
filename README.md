![Go Clean Template](docs/img/logo.svg)

# Go Clean template

[🇨🇳 中文](README_CN.md)
[🇷🇺 RU](README_RU.md)

Clean Architecture template for Golang services

[![Release](https://img.shields.io/github/v/release/evrone/go-clean-template.svg)](https://github.com/evrone/go-clean-template/releases/)
[![License](https://img.shields.io/badge/License-MIT-success)](https://github.com/evrone/go-clean-template/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/evrone/go-clean-template)](https://goreportcard.com/report/github.com/evrone/go-clean-template)
[![codecov](https://codecov.io/gh/evrone/go-clean-template/branch/master/graph/badge.svg?token=XE3E0X3EVQ)](https://codecov.io/gh/evrone/go-clean-template)

[![Web Framework](https://img.shields.io/badge/Fiber-Web%20Framework-blue)](https://github.com/gofiber/fiber)
[![API Documentation](https://img.shields.io/badge/Swagger-API%20Documentation-blue)](https://github.com/swaggo/swag)
[![Validation](https://img.shields.io/badge/Validator-Data%20Integrity-blue)](https://github.com/go-playground/validator)
[![JSON Handling](https://img.shields.io/badge/Go--JSON-Fast%20Serialization-blue)](https://github.com/goccy/go-json)
[![Query Builder](https://img.shields.io/badge/Squirrel-SQL%20Query%20Builder-blue)](https://github.com/Masterminds/squirrel)
[![Database Migrations](https://img.shields.io/badge/Migrations-Seamless%20Schema%20Updates-blue)](https://github.com/golang-migrate/migrate)
[![Logging](https://img.shields.io/badge/ZeroLog-Structured%20Logging-blue)](https://github.com/rs/zerolog)
[![Metrics](https://img.shields.io/badge/Prometheus-Metrics%20Integration-blue)](https://github.com/ansrivas/fiberprometheus)
[![Testing](https://img.shields.io/badge/Testify-Testing%20Framework-blue)](https://github.com/stretchr/testify)
[![Mocking](https://img.shields.io/badge/Mock-Mocking%20Library-blue)](https://go.uber.org/mock)

## Overview

The purpose of the template is to show:

- how to organize a project and prevent it from turning into spaghetti code
- where to store business logic so that it remains independent, clean, and extensible
- how not to lose control when a microservice grows

Using the principles of Robert Martin (aka Uncle Bob).

[Go-clean-template](https://evrone.com/go-clean-template?utm_source=github&utm_campaign=go-clean-template) is created &
supported by [Evrone](https://evrone.com/?utm_source=github&utm_campaign=go-clean-template).

This template implements three types of servers:

- AMQP RPC (based on RabbitMQ as [transport](https://github.com/rabbitmq/amqp091-go)
  and [Request-Reply pattern](https://www.enterpriseintegrationpatterns.com/patterns/messaging/RequestReply.html))
- MQ RPC (based on NATS as [transport](https://github.com/nats-io/nats.go)
  and [Request-Reply pattern](https://www.enterpriseintegrationpatterns.com/patterns/messaging/RequestReply.html))
- gRPC ([gRPC](https://grpc.io/) framework based on protobuf)
- REST API ([Fiber](https://github.com/gofiber/fiber) framework)

## Content

- [Quick start](#quick-start)
- [Project structure](#project-structure)
- [Dependency Injection](#dependency-injection)
- [Clean Architecture](#clean-architecture)

## Quick start

### Local development

```sh
# Postgres, RabbitMQ, NATS
make compose-up
# Run app with migrations
make run
```

### Integration tests (can be run in CI)

```sh
# DB, app + migrations, integration tests
make compose-up-integration-test
```

### Full docker stack with reverse proxy

```sh
make compose-up-all 
```

Check services:

- AMQP RPC:
  - URL: `amqp://guest:guest@127.0.0.1:5672/`
  - Client Exchange: `rpc_client`
  - Server Exchange: `rpc_server`
- NATS RPC:
  - URL: `nats://guest:guest@127.0.0.1:4222/`
  - Server Exchange: `rpc_server`
- REST API:
  - http://app.lvh.me/healthz | http://127.0.0.1:8080/healthz
  - http://app.lvh.me/metrics | http://127.0.0.1:8080/metrics
  - http://app.lvh.me/swagger | http://127.0.0.1:8080/swagger
- gRPC:
  - URL: `tcp://grpc.lvh.me:8081` | `tcp://127.0.0.1:8081`
  - [v1/translation.history.proto](docs/proto/v1/translation.history.proto)
- PostgreSQL:
  - `postgres://user:myAwEsOm3pa55@w0rd@127.0.0.1:5432/db`
- RabbitMQ:
  - http://rabbitmq.lvh.me | http://127.0.0.1:15672
  - Credentials: `guest` / `guest`
- NATS monitoring:
  - http://nats.lvh.me | http://127.0.0.1:8222/
  - Credentials: `guest` / `guest`

## Project structure

### `cmd/app/main.go`

Configuration and logger initialization. Then the main function "continues" in
`internal/app/app.go`.

### `config`

The twelve-factor app stores config in environment variables (often shortened to `env vars` or `env`). Env vars are easy
to change between deploys without changing any code; unlike config files, there is little chance of them being checked
into the code repo accidentally; and unlike custom config files, or other config mechanisms such as Java System
Properties, they are a language- and OS-agnostic standard.

Config: [config.go](config/config.go)

Example: [.env.example](.env.example)

[docker-compose.yml](docker-compose.yml) uses `env` variables to configure services.

### `docs`

Swagger documentation. Auto-generated by [swag](https://github.com/swaggo/swag) library.
You don't need to correct anything by yourself.

#### `docs/proto`

Protobuf files. They are used to generate Go code for gRPC services.
The proto files are also used to generate documentation for gRPC services.
You don't need to correct anything by yourself.

### `integration-test`

Integration tests.
They are launched as a separate container, next to the application container.

### `internal/app`

There is always one _Run_ function in the `app.go` file, which "continues" the _main_ function.

This is where all the main objects are created.
Dependency injection occurs through the "New ..." constructors (see Dependency Injection).
This technique allows us to layer the application using the [Dependency Injection](#dependency-injection) principle.
This makes the business logic independent from other layers.

Next, we start the server and wait for signals in _select_ for graceful completion.
If `app.go` starts to grow, you can split it into multiple files.

For a large number of injections, [wire](https://github.com/google/wire) can be used.

The `migrate.go` file is used for database auto migrations.
It is included if an argument with the _migrate_ tag is specified.
For example:

```sh
go run -tags migrate ./cmd/app
```

### `internal/controller`

Server handler layer (MVC controllers). The template shows 3 servers:

- AMQP RPC (based on RabbitMQ as transport)
- gRPC ([gRPC](https://grpc.io/) framework based on protobuf)
- REST API ([Fiber](https://github.com/gofiber/fiber) framework)

Server routers are written in the same style:

- Handlers are grouped by area of application (by a common basis)
- For each group, its own router structure is created, the methods of which process paths
- The structure of the business logic is injected into the router structure, which will be called by the handlers

#### `internal/controller/amqp_rpc`

Simple RPC versioning.
For v2, we will need to add the `amqp_rpc/v2` folder with the same content.
And in the file `internal/controller/amqp_rpc/router.go` add the line:

```go
routes := make(map[string]server.CallHandler)

{
    v1.NewTranslationRoutes(routes, t, l)
}

{
    v2.NewTranslationRoutes(routes, t, l)
}
```

#### `internal/controller/grpc`

Simple gRPC versioning.
For v2, we will need to add the `grpc/v2` folder with the same content.
Also add the `v2` folder to the proto files in `docs/proto`.
And in the file `internal/controller/grpc/router.go` add the line:

```go
{
    v1.NewTranslationRoutes(app, t, l)
}

{
    v2.NewTranslationRoutes(app, t, l)
}

reflection.Register(app)
```

#### `internal/controller/http`

Simple REST versioning.
For v2, we will need to add the `http/v2` folder with the same content.
And in the file `internal/controller/http/router.go` add the line:

```go
apiV1Group := app.Group("/v1")
{
	v1.NewTranslationRoutes(apiV1Group, t, l)
}
apiV2Group := app.Group("/v2")
{
	v2.NewTranslationRoutes(apiV2Group, t, l)
}
```

Instead of [Fiber](https://github.com/gofiber/fiber), you can use any other http framework.

In `router.go` and above the handler methods, there are comments for generating swagger documentation
using [swag](https://github.com/swaggo/swag).

### `internal/entity`

Entities of business logic (models) can be used in any layer.
There can also be methods, for example, for validation.

### `internal/usecase`

Business logic.

- Methods are grouped by area of application (on a common basis)
- Each group has its own structure
- One file - one structure

Repositories, webapi, rpc, and other business logic structures are injected into business logic structures
(see [Dependency Injection](#dependency-injection)).

#### `internal/repo/persistent`

A repository is an abstract storage (database) that business logic works with.

#### `internal/repo/webapi`

It is an abstract web API that business logic works with.
For example, it could be another microservice that business logic accesses via the REST API.
The package name changes depending on the purpose.

### `pkg/rabbitmq`

RabbitMQ RPC pattern:

- There is no routing inside RabbitMQ
- Exchange fanout is used, to which 1 exclusive queue is bound, this is the most productive config
- Reconnect on the loss of connection

## Dependency Injection

In order to remove the dependence of business logic on external packages, dependency injection is used.

For example, through the New constructor, we inject the dependency into the structure of the business logic.
This makes the business logic independent (and portable).
We can override the implementation of the interface without making changes to the `usecase` package.

```go
package usecase

import (
// Nothing!
)

type Repository interface {
	Get()
}

type UseCase struct {
	repo Repository
}

func New(r Repository) *UseCase {
	return &UseCase{
		repo: r,
	}
}

func (uc *UseCase) Do() {
	uc.repo.Get()
}
```

It will also allow us to do auto-generation of mocks (for example with [mockery](https://github.com/vektra/mockery)) and
easily write unit tests.

> We are not tied to specific implementations in order to always be able to change one component to another.
> If the new component implements the interface, nothing needs to be changed in the business logic.

## Clean Architecture

### Key idea

Programmers realize the optimal architecture for an application after most of the code has been written.

> A good architecture allows decisions to be delayed to as late as possible.

### The main principle

Dependency Inversion (the same one from SOLID) is the principle of dependency injection.
The direction of dependencies goes from the outer layer to the inner layer.
Due to this, business logic and entities remain independent from other parts of the system.

So, the application is divided into 2 layers, internal and external:

1. **Business logic** (Go standard library).
2. **Tools** (databases, servers, message brokers, any other packages and frameworks).

![Clean Architecture](docs/img/layers-1.png)

**The inner layer** with business logic should be clean. It should:

- Not have package imports from the outer layer.
- Use only the capabilities of the standard library.
- Make calls to the outer layer through the interface (!).

The business logic doesn't know anything about Postgres or a specific web API.
Business logic has an interface for working with an _abstract_ database or _abstract_ web API.

**The outer layer** has other limitations:

- All components of this layer are unaware of each other's existence. How to call another from one tool? Not directly,
  only through the inner layer of business logic.
- All calls to the inner layer are made through the interface (!).
- Data is transferred in a format that is convenient for business logic (`internal/entity`).

For example, you need to access the database from HTTP (controller).
Both HTTP and database are in the outer layer, which means they know nothing about each other.
The communication between them is carried out through `usecase` (business logic):

```
    HTTP > usecase
           usecase > repository (Postgres)
           usecase < repository (Postgres)
    HTTP < usecase
```

The symbols > and < show the intersection of layer boundaries through Interfaces.
The same is shown in the picture:

![Example](docs/img/example-http-db.png)

Or more complex business logic:

```
    HTTP > usecase
           usecase > repository
           usecase < repository
           usecase > webapi
           usecase < webapi
           usecase > RPC
           usecase < RPC
           usecase > repository
           usecase < repository
    HTTP < usecase
```

### Layers

![Example](docs/img/layers-2.png)

### Clean Architecture Terminology

- **Entities** are structures that business logic operates on.
  They are located in the `internal/entity` folder.
  In MVC terms, entities are models.
- **Use Cases** is business logic located in `internal/usecase`.

The layer with which business logic directly interacts is usually called the _infrastructure_ layer.
These can be repositories `internal/usecase/repo`, external webapi `internal/usecase/webapi`, any pkg, and other
microservices.
In the template, the _infrastructure_ packages are located inside `internal/usecase`.

You can choose how to call the entry points as you wish. The options are:

- controller (in our case)
- delivery
- transport
- gateways
- entrypoints
- primary
- input

### Additional layers

The classic version
of [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) was designed for
building large monolithic applications and has 4 layers.

In the original version, the outer layer is divided into two more, which also have an inversion of dependencies
to each other (directed inward) and communicate through interfaces.

The inner layer is also divided into two (with separation of interfaces), in the case of complex logic.

---

Complex tools can be divided into additional layers.
However, you should add layers only if really necessary.

### Alternative approaches

In addition to Clean architecture, _Onion architecture_ and _Hexagonal_ (_Ports and adapters_) are similar to it.
Both are based on the principle of Dependency Inversion.
_Ports and adapters_ are very close to _Clean Architecture_, the differences are mainly in terminology.

## Similar projects

- [https://github.com/bxcodec/go-clean-arch](https://github.com/bxcodec/go-clean-arch)
- [https://github.com/zhashkevych/courses-backend](https://github.com/zhashkevych/courses-backend)

## Useful links

- [The Clean Architecture article](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Twelve factors](https://12factor.net/ru/)
