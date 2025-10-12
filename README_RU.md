![Go Clean Template](docs/img/logo.svg)

# Go Чистая Архитектура

Шаблон Чистой Архитектуры для приложений на Golang

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

## Обзор

Цель этого шаблона - показать принципы Чистой Архитектуры Роберта Мартина (дядюшки Боба):

- как структурировать проект и не дать ему превратиться в спагетти-код
- где хранить бизнес-логику, чтобы она оставалась независимой, чистой и расширяемой
- как не потерять контроль при росте проекта

[Go-clean-template](https://evrone.com/go-clean-template?utm_source=github&utm_campaign=go-clean-template) создан и
поддерживается [Evrone](https://evrone.com/?utm_source=github&utm_campaign=go-clean-template).

Этот шаблон поддерживает три типа серверов:

- AMQP RPC (на основе RabbitMQ в качестве [транспорта](https://github.com/rabbitmq/amqp091-go)
  и [Request-Reply паттерна]((https://www.enterpriseintegrationpatterns.com/patterns/messaging/RequestReply.html)))
- NATS RPC (на основе NATS в качестве [транспорта](https://github.com/nats-io/nats.go)
  и [Request-Reply паттерна]((https://www.enterpriseintegrationpatterns.com/patterns/messaging/RequestReply.html))))
- gRPC ([gRPC](https://grpc.io/) фреймворк на основе protobuf)
- REST API ([Fiber](https://github.com/gofiber/fiber) фреймворк)

## Содержание

- [Быстрый старт](#быстрый-старт)
- [Структура проекта](#структура-проекта)
- [Внедрение зависимостей](#внедрение-зависимостей)
- [Чистая Архитектура](#чистая-архитектура)

## Быстрый старт

### Локальная разработка

```sh
# Postgres, RabbitMQ, NATS
make compose-up
# Запуск приложения и миграций
make run
```

### Интеграционные тесты (может быть использовано с CI)

```sh
# DB, app + migrations, integration tests
make compose-up-integration-test
```

### Весь docker stack с reverse proxy

```sh
make compose-up-all 
```

Проверьте сервисы:

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

## Структура проекта

### `cmd/app/main.go`

Инициализация конфигурации и логгера. Здесь вызывается основная часть приложения из `internal/app/app.go`.

### `config`

Приложение двенадцати факторов хранит конфигурацию в переменных окружения (часто сокращается до `env vars` или `env`).
Переменные окружения легко изменить между развёртываниями, не изменяя код; в отличие от файлов конфигурации, менее
вероятно случайно сохранить их в репозиторий кода; и в отличие от пользовательских конфигурационных файлов или других
механизмов конфигурации, таких как Java System Properties, они являются независимым от языка и операционной системы
стандартом.

Конфигурация: [config.go](config/config.go)

Пример: [.env.example](.env.example)

[docker-compose.yml](docker-compose.yml) использует переменные `env` для настройки сервисов.

### `docs`

Документация Swagger. Генерируется автоматически с помощью библиотеки [swag](https://github.com/swaggo/swag).
Вам не нужно ничего редактировать вручную.

#### `docs/proto`

Protobuf файлы. Они используются для генерации Go-кода для gRPC сервисов.
Protobuf файлы также используются для генерации документации для gRPC сервисов.
Вам не нужно ничего исправлять самостоятельно.

### `integration-test`

Интеграционные тесты.
Они запускаются в отдельном контейнере, рядом с контейнером приложения.

### `internal/app`

Здесь находится только одна функция _Run_. Она размещена в файле `app.go` и является логическим продолжением функции
_main_.

Здесь создаются все основные объекты.
[Внедрение зависимостей](#внедрение-зависимостей) происходит через конструктор "New ...".
Это позволяет слоировать приложение, делая бизнес-логику независимой от других слоев.

Далее запускается сервер и ожидается сигнал в _select_ для корректного завершения работы.
Если `app.go` стал слишком большим, вы можете разделить его на несколько файлов.

Если зависимостей много, то для удобства можно использовать [wire](https://github.com/google/wire).

Файл `migrate.go` используется для автоматической миграции базы данных.
Он включается в компиляцию только при указании тега _migrate_.
Пример:

```sh
go run -tags migrate ./cmd/app
```

### `internal/controller`

Слой хэндлеров сервера (MVC контроллеры). В шаблоне показана работа 3 серверов:

- AMQP RPC (на основе RabbitMQ в качестве транспорта)
- gRPC ([gRPC](https://grpc.io/) фреймворк на основе protobuf)
- REST API ([Fiber](https://github.com/gofiber/fiber) фреймворк)

Маршрутизаторы http сервера пишутся в едином стиле:

- Хэндлеры группируются по области применения (по общему критерию)
- Для каждой группы создается свой маршрутизатор
- Объект бизнес-логики передается в маршрутизатор, чтобы быть доступным внутри хэндлеров

#### `internal/controller/amqp_rpc`

Простое версионирование RPC.
Для версии v2 нужно будет добавить папку `amqp_rpc/v2` с таким же содержимым.
А в файле `internal/controller/amqp_rpc/router.go` добавить строку:

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

Простое версионирование gRPC.  
Для версии v2 нужно будет добавить папку `grpc/v2` с таким же содержимым.  
Также добавьте папку `v2` в proto-файлы в `docs/proto`.  
И в файле `internal/controller/grpc/router.go` добавьте строку:

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

Простое версионирование REST API.
Для создания версии v2 нужно создать папку `http/v2` с таким же содержимым.
Добавить в файл `internal/controller/http/router.go` строки:

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

Вместо [Fiber](https://github.com/gofiber/fiber) можно использовать любой другой http фреймворк.

В файле `router.go` над хэндлером написаны комментарии для генерации документации через
swagger [swag](https://github.com/swaggo/swag).

### `internal/entity`

Сущности бизнес-логики (модели). Могут быть использованы в любом слое.
Также они могут иметь методы, например, для валидации.

### `internal/usecase`

Бизнес-логика.

- Методы группируются по области применения (по общему критерию)
- У каждой группы своя отдельная структура
- Один файл - одна структура

Репозитории, webapi, rpc и другие структуры передаются в слой бизнес-логики в связующем файле `internal/app/app.go`
(смотрите [Внедрение зависимостей](#внедрение-зависимостей)).

#### `internal/repo/persistent`

Репозиторий — это абстрактное хранилище (база данных), с которым взаимодействует бизнес-логика.

#### `internal/repo/webapi`

Это абстрактное web API, с которым взаимодействует бизнес-логика.
Например, это может быть внешний микросервис, к которому бизнес-логика обращается через REST API.
Название пакета выбирается таким, чтобы соответствовать его назначению.

### `pkg/rabbitmq`

RabbitMQ RPC паттерн:

- Внутри RabbitMQ не используется маршрутизация
- Используется fanout-обмен, к которому привязана одна эксклюзивная очередь - это наиболее производительная конфигурация
- Переподключение при потере соединения

## Внедрение зависимостей

Для устранения зависимости бизнес-логики от внешних пакетов используется внедрение зависимостей.

Например, через конструктор "New" внедряется репозиторий в слой бизнес-логики.
Это делает бизнес-логику независимой и переносимой.
Мы можем переписать реализацию интерфейса репозитория, не внося изменения в пакет бизнес-логики `usecase`.

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

Благодаря разделению через интерфейсы можно генерировать моки (например,
используя [mockery](https://github.com/vektra/mockery)) и легко писать юнит-тесты.

> Мы не привязаны к конкретным реализациям и всегда можем заменить один компонент на другой.
> Если новый компонент реализует интерфейс, то в бизнес-логике ничего не нужно менять.

## Чистая Архитектура

### Ключевая идея

Программисты создают оптимальную архитектуру приложения после написания основной части кода.

> Хорошая архитектура позволяет откладывать изменения как можно дольше.

### Основной принцип

Инверсия зависимостей (та же, что и в SOLID) используется как принцип для внедрения зависимостей.
Зависимости направлены от внешнего слоя к внутреннему.
Благодаря этому бизнес-логика и сущности остаются независимыми от других частей системы.

Например, приложение можно разделить на два слоя - внутренний и внешний:

1. **Бизнес-логика** (например, стандартная библиотека Go).
2. **Инструменты** (базы данных, серверы, брокеры сообщений и другие библиотеки и фреймворки).

![Чистая архитектура](docs/img/layers-1.png)

**Внутренний слой** с бизнес-логикой должен быть чистым. Он обязан:

- Не импортировать пакеты из внешних слоев.
- Использовать только стандартную библиотеку.
- Взаимодействовать с внешними слоями через интерфейсы (!).

Бизнес-логика не должна ничего знать о Postgres или о реализации web API.
Бизнес-логика имеет интерфейс для взаимодействия с _абстрактной_ базой данных или _абстрактным_ web API.

**Внешний слой** имеет ограничения:

- Компоненты этого слоя не могут знать друг о друге и взаимодействовать напрямую. Обращение друг к другу происходит
  через внутренний слой - слой бизнес-логики.
- Вызовы во внутренний слой выполняются через интерфейсы (!).
- Данные передаются в формате, удобном для бизнес-логики (структуры хранятся в `internal/entity`).

Например, нужно обратиться к базе данных из HTTP хэндлера (в слое контроллер).
База данных и HTTP находятся во внешнем слое. Они не знают друг о друге ничего и не могут взаимодействовать напрямую.
Взаимодействие будет происходить через слой бизнес-логики `usecase`:

```
    HTTP > usecase
           usecase > repository (Postgres)
           usecase < repository (Postgres)
    HTTP < usecase
```

Символы > и < показывают пересечения слоев через интерфейсы и направления.
Это же показано на схеме:

![Пример](docs/img/example-http-db.png)

Пример более сложного пути данных:

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

### Слои

![Пример](docs/img/layers-2.png)

### Терминология в Чистой Архитектуре

- **Entities** (сущности) - это структуры, с которыми работает бизнес-логика.
  Они располагаются в папке `internal/entity`.
  В терминологии MVC сущности - это модели.
- **Use Cases** - это бизнес-логика. Располагается в папке `internal/usecase`.

Слой, с которым бизнес-логика взаимодействует напрямую, обычно называется _инфраструктурным_ слоем.
Это может быть репозиторий `internal/usecase/repo`, внешнее webapi `internal/usecase/webapi`, любой пакет или
микросервис.
В шаблоне пакеты _infrastructure_ размещены внутри `internal/usecase`.

Вы можете выбирать, как называть точки входа, по своему усмотрению. Варианты такие:

- controller (в нашем случае)
- delivery
- transport
- gateways
- entrypoints
- primary
- input

### Дополнительные слои

В классической версии [Чистой Архитектуры](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
для создания больших монолитных приложений предложено 4 слоя.

В исходной версии внешний слой делится на два, которые также имеют инверсию зависимостей в другие слои и взаимодействуют
через интерфейсы.

Внутренний слой также делится на два (с использованием интерфейсов) в случае сложной логики.

---

Сложные инструменты могут быть разделены на дополнительные слои.
Однако добавлять слои следует только в том случае, если это действительно необходимо.

### Другие подходы

Кроме Чистой Архитектуры есть и другие подходы:

- Луковая Архитектура
- Гексагональная (_Порты и адаптеры_ также похожа на неё)
  Они обе основаны на принципе инверсии зависимостей.
  _Порты и адаптеры_ очень похожи на _Чистую Архитектуру_. Различия в основном заключаются в терминологии.

## Похожие проекты

- [https://github.com/bxcodec/go-clean-arch](https://github.com/bxcodec/go-clean-arch)
- [https://github.com/zhashkevych/courses-backend](https://github.com/zhashkevych/courses-backend)

## Дополнительная информация

- [The Clean Architecture article](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Twelve factors](https://12factor.net/ru/)
