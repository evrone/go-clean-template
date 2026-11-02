![Go Clean Template](docs/img/logo.svg)

# Go Clean template

[🇹🇼 繁體中文](README_zh-TW.md)
[🇨🇳 中文](README_CN.md)
[🇷🇺 RU](README_RU.md)

適用於 Golang 服務的 Clean Architecture（整潔架構）範本

[![Release](https://img.shields.io/github/v/release/evrone/go-clean-template.svg)](https://github.com/evrone/go-clean-template/releases/)
[![License](https://img.shields.io/badge/License-MIT-success)](https://github.com/evrone/go-clean-template/blob/master/LICENSE)
[![codecov](https://codecov.io/gh/evrone/go-clean-template/branch/master/graph/badge.svg?token=XE3E0X3EVQ)](https://codecov.io/gh/evrone/go-clean-template)

[![Web Framework](https://img.shields.io/badge/Fiber-Web%20Framework-blue)](https://github.com/gofiber/fiber)
[![API Documentation](https://img.shields.io/badge/Swagger-API%20Documentation-blue)](https://github.com/swaggo/swag)
[![Validation](https://img.shields.io/badge/Validator-Data%20Integrity-blue)](https://github.com/go-playground/validator)
[![JSON Handling](https://img.shields.io/badge/Go--JSON-Fast%20Serialization-blue)](https://github.com/goccy/go-json)
[![Query Builder](https://img.shields.io/badge/Squirrel-SQL%20Query%20Builder-blue)](https://github.com/Masterminds/squirrel)
[![Database Migrations](https://img.shields.io/badge/Migrations-Seamless%20Schema%20Updates-blue)](https://github.com/golang-migrate/migrate)
[![Logging](https://img.shields.io/badge/ZeroLog-Structured%20Logging-blue)](https://github.com/rs/zerolog)
[![Metrics](https://img.shields.io/badge/Prometheus-Metrics%20Integration-blue)](https://github.com/ansrivas/fiberprometheus)
[![Tracing](https://img.shields.io/badge/OpenTelemetry-Distributed%20Tracing-blue)](https://opentelemetry.io/)
[![Testing](https://img.shields.io/badge/Testify-Testing%20Framework-blue)](https://github.com/stretchr/testify)
[![Mocking](https://img.shields.io/badge/Mock-Mocking%20Library-blue)](https://go.uber.org/mock)

## 概述

這個範本的目的是示範：

- 如何組織專案，避免它變成義大利麵式程式碼（spaghetti code）
- 商業邏輯該放在哪裡，才能保持獨立、乾淨且易於擴充
- 微服務逐漸變大時，如何不失去掌控

遵循 Robert Martin（又稱 Uncle Bob）所提出的原則。

[Go-clean-template](https://evrone.com/go-clean-template?utm_source=github&utm_campaign=go-clean-template) 由
[Evrone](https://evrone.com/?utm_source=github&utm_campaign=go-clean-template) 建立並維護。

這個範本實作了四種伺服器：

- AMQP RPC（以 RabbitMQ 作為[傳輸層](https://github.com/rabbitmq/amqp091-go)，
  採用[請求-回覆模式（Request-Reply pattern）](https://www.enterpriseintegrationpatterns.com/patterns/messaging/RequestReply.html)）
- MQ RPC（以 NATS 作為[傳輸層](https://github.com/nats-io/nats.go)，
  採用[請求-回覆模式（Request-Reply pattern）](https://www.enterpriseintegrationpatterns.com/patterns/messaging/RequestReply.html)）
- gRPC（以 protobuf 為基礎的 [gRPC](https://grpc.io/) 框架）
- REST API（[Fiber](https://github.com/gofiber/fiber) 框架）

範本包含三個領域，用來示範多服務架構：

- **使用者驗證** — 註冊、登入、以 JWT 為基礎的授權
- **任務管理** — 含狀態轉換（todo、in_progress、done）的 CRUD 操作
- **翻譯** — 文字翻譯，並記錄翻譯歷史

所有領域都可透過四種傳輸方式使用（REST、gRPC、AMQP RPC、NATS RPC）。

## 目錄

- [領域](#領域)
- [快速開始](#快速開始)
- [可觀測性](#可觀測性)
- [專案結構](#專案結構)
- [相依性注入](#相依性注入)
- [整潔架構](#整潔架構)

## 領域

範本內含三個完整實作的領域，每個領域都可透過四種傳輸方式使用（REST、gRPC、AMQP RPC、NATS RPC）。

### 使用者驗證

註冊、登入，以及以 JWT 為基礎的授權。

| 操作         | REST                     | gRPC                     |
|--------------|--------------------------|--------------------------|
| 註冊         | `POST /v1/auth/register` | `AuthService/Register`   |
| 登入         | `POST /v1/auth/login`    | `AuthService/Login`      |
| 取得個人資料 | `GET /v1/user/profile`   | `AuthService/GetProfile` |

- 密碼以 bcrypt 雜湊
- JWT 權杖的有效期限可設定
- 所有傳輸方式都有驗證中介軟體（middleware）

### 任務管理

含狀態機（state machine）的 CRUD 操作。

| 操作     | REST                         | gRPC                         |
|----------|------------------------------|------------------------------|
| 建立     | `POST /v1/tasks`             | `TaskService/CreateTask`     |
| 列表     | `GET /v1/tasks`              | `TaskService/ListTasks`      |
| 取得     | `GET /v1/tasks/:id`          | `TaskService/GetTask`        |
| 更新     | `PUT /v1/tasks/:id`          | `TaskService/UpdateTask`     |
| 狀態轉換 | `PATCH /v1/tasks/:id/status` | `TaskService/TransitionTask` |
| 刪除     | `DELETE /v1/tasks/:id`       | `TaskService/DeleteTask`     |

- 狀態轉換：`todo` → `in_progress` → `done`（以及 `in_progress` → `todo`）
- 以 `limit`/`offset` 分頁，並可選擇性地依狀態篩選
- 任務僅限於已驗證的使用者本人

### 翻譯

透過外部 API 進行文字翻譯，並記錄翻譯歷史。

| 操作 | REST                                | gRPC                                    |
|------|-------------------------------------|-----------------------------------------|
| 翻譯 | `POST /v1/translation/do-translate` | `TranslationHistoryService/DoTranslate` |
| 歷史 | `GET /v1/translation/history`       | `TranslationHistoryService/ShowHistory` |

## 快速開始

### 本機開發

```sh
# Postgres, RabbitMQ, NATS
make compose-up
# Run app with migrations
make run
```

### 整合測試（可在 CI 中執行）

```sh
# DB, app + migrations, integration tests
make compose-up-integration-test
```

### 搭配反向代理的完整 Docker 環境

```sh
make compose-up-all 
```

檢查各項服務：

- AMQP RPC：
  - URL：`amqp://guest:guest@127.0.0.1:5672/`
  - Client Exchange：`rpc_client`
  - Server Exchange：`rpc_server`
- NATS RPC：
  - URL：`nats://guest:guest@127.0.0.1:4222/`
  - Server Exchange：`rpc_server`
- REST API：
  - http://app.lvh.me/healthz | http://127.0.0.1:8080/healthz
  - http://app.lvh.me/metrics | http://127.0.0.1:8080/metrics
  - http://app.lvh.me/swagger | http://127.0.0.1:8080/swagger
- gRPC：
  - URL：`tcp://grpc.lvh.me:8081` | `tcp://127.0.0.1:8081`
  - [v1/auth.proto](docs/proto/v1/auth.proto)
  - [v1/task.proto](docs/proto/v1/task.proto)
  - [v1/translation.history.proto](docs/proto/v1/translation.history.proto)
- PostgreSQL：
  - `postgres://user:myAwEsOm3pa55@w0rd@127.0.0.1:5432/db`
- RabbitMQ：
  - http://rabbitmq.lvh.me | http://127.0.0.1:15672
  - 帳號密碼：`guest` / `guest`
- NATS 監控：
  - http://nats.lvh.me | http://127.0.0.1:8222/
  - 帳號密碼：`guest` / `guest`
- Jaeger（追蹤介面）：
  - http://jaeger.lvh.me | http://127.0.0.1:16686

## 可觀測性

分散式追蹤（distributed tracing）由 [OpenTelemetry](https://opentelemetry.io/) 提供。Span 會透過 OTLP/gRPC 匯出到
collector — 在 Docker 環境中是 [Jaeger](https://www.jaegertracing.io/)。

- **Context 傳遞** — 使用 W3C `traceparent` + `baggage`，讓同一條 trace 能橫跨四種傳輸方式。REST 使用
  [otelfiber](https://github.com/gofiber/contrib/tree/main/otelfiber) 中介軟體，gRPC 使用
  [otelgrpc](https://github.com/open-telemetry/opentelemetry-go-contrib) stats handler，而 AMQP RPC / NATS RPC 則透過自訂的 carrier，
  將 trace context 放在訊息標頭中傳遞（`pkg/rabbitmq/rmq_rpc/otel_carrier.go`、
  `pkg/nats/nats_rpc/otel_carrier.go`）。
- **已插入追蹤的層** — use case 與 repository 都以追蹤用的裝飾器（decorator）包裝（`*/tracing.go`），因此一條 trace
  能呈現完整路徑：controller → usecase → repo / webapi。
- **停用時為 no-op** — 設定 `TRACING_ENABLED=false` 時會安裝 no-op 的 tracer provider，因此追蹤程式碼可無條件執行，
  且沒有額外負擔。
- **取樣** — `TRACING_SAMPLE_RATE` 是以父 span 為準（parent-based）的比例。本機請使用 `1.0`；正式環境請調低（`0.1` / `0.01`）。
- 內部路由（`/healthz`、`/metrics`、`/swagger`）不納入追蹤 — 只有帶版本的 API 群組會被插入追蹤。

設定（請參閱 [.env.example](.env.example)）：

| 變數                    | 預設值             | 說明                              |
|-------------------------|--------------------|-----------------------------------|
| `TRACING_ENABLED`       | `false`            | 啟用 OpenTelemetry 追蹤           |
| `TRACING_OTLP_ENDPOINT` | `localhost:4317`   | OTLP/gRPC collector 端點          |
| `TRACING_OTLP_INSECURE` | `true`             | 停用 OTLP exporter 的 TLS         |
| `TRACING_SAMPLE_RATE`   | `0.1`              | 以父 span 為準的取樣比例          |

## 專案結構

### `cmd/app/main.go`

初始化設定與日誌（logger）。接著 main 函式會「接續」到
`internal/app/app.go`。

### `config`

Twelve-factor app 把設定存放在環境變數中（常簡稱為 `env vars` 或 `env`）。環境變數很容易在不同部署之間更換，
完全不必修改程式碼；與設定檔不同，它們被意外提交到程式碼儲存庫的機率很低；而且相較於自訂的設定檔，
或是 Java System Properties 這類其他設定機制，它們是不受語言與作業系統限制的通用標準。

設定：[config.go](config/config.go)

範例：[.env.example](.env.example)

[docker-compose.yml](docker-compose.yml) 使用 `env` 變數來設定各項服務。

### `docs`

Swagger 文件，由 [swag](https://github.com/swaggo/swag) 函式庫自動產生。
你不需要自行修改任何內容。

#### `docs/proto`

Protobuf 檔案，用來產生 gRPC 服務的 Go 程式碼。
這些 proto 檔案也會用來產生 gRPC 服務的文件。
你不需要自行修改任何內容。

### `integration-test`

整合測試。
它們會以獨立的容器啟動，與應用程式容器並排執行。

### `internal/app`

`app.go` 檔案中永遠有一個 _Run_ 函式，它「接續」_main_ 函式。

所有主要的物件都在這裡建立。
相依性注入透過「New ...」建構函式完成（請參閱相依性注入）。
這個技巧讓我們能依照[相依性注入](#相依性注入)原則來分層。
這使得商業邏輯獨立於其他層。

接著，我們啟動伺服器，並在 _select_ 中等待訊號，以便平順地結束。
如果 `app.go` 開始變得龐大，可以把它拆成多個檔案。

需要注入的相依物件很多時，可以使用 [wire](https://github.com/google/wire)。

`migrate.go` 檔案用於資料庫的自動遷移（migration）。
只有在指定了 _migrate_ 標籤（tag）的參數時才會被包含進來。
例如：

```sh
go run -tags migrate ./cmd/app
```

### `internal/controller`

伺服器處理常式（handler）層（MVC 的 controller）。範本示範了 4 種伺服器：

- AMQP RPC（以 RabbitMQ 作為傳輸層）
- NATS RPC（以 NATS 作為傳輸層）
- gRPC（以 protobuf 為基礎的 [gRPC](https://grpc.io/) 框架）
- REST API（[Fiber](https://github.com/gofiber/fiber) 框架）

各伺服器的路由器（router）採用相同的寫法：

- 處理常式依應用領域（共同的基礎）分組
- 每個群組各自建立一個路由器結構，其方法負責處理各個路徑
- 商業邏輯的結構會被注入到路由器結構中，由處理常式呼叫

#### `internal/controller/amqp_rpc`

簡單的 RPC 版本管理。
若要新增 v2，需要加入內容相同的 `amqp_rpc/v2` 資料夾，
並在 `internal/controller/amqp_rpc/router.go` 檔案中加入這幾行：

```go
routes := make(map[string]server.CallHandler)

{
    v1.NewRoutes(routes, t, u, tk, j, l)
}

{
    v2.NewTranslationRoutes(routes, t, l)
}
```

#### `internal/controller/grpc`

簡單的 gRPC 版本管理。
若要新增 v2，需要加入內容相同的 `grpc/v2` 資料夾。
同時也要在 `docs/proto` 的 proto 檔案中加入 `v2` 資料夾。
並在 `internal/controller/grpc/router.go` 檔案中加入這幾行：

```go
{
    v1.NewAuthRoutes(app, u, l)
    v1.NewTaskRoutes(app, tk, l)
    v1.NewTranslationRoutes(app, t, l)
}

{
    v2.NewAuthRoutes(app, u, l)
    v2.NewTaskRoutes(app, tk, l)
    v2.NewTranslationRoutes(app, t, l)
}

reflection.Register(app)
```

#### `internal/controller/nats_rpc`

簡單的 RPC 版本管理。
若要新增 v2，需要加入內容相同的 `nats_rpc/v2` 資料夾。
並在 `internal/controller/nats_rpc/router.go` 檔案中加入這幾行：

```go
routes := make(map[string]server.CallHandler)

{
    v1.NewRoutes(routes, t, u, tk, j, l)
}

{
    v2.NewTranslationRoutes(routes, t, l)
}
```

#### `internal/controller/restapi`

簡單的 REST 版本管理。
若要新增 v2，需要加入內容相同的 `restapi/v2` 資料夾。
並在 `internal/controller/restapi/router.go` 檔案中加入這幾行：

```go
apiV1Group := app.Group("/v1")
{
	v1.NewRoutes(apiV1Group, t, u, tk, jwtManager, l)
}
apiV2Group := app.Group("/v2")
{
	v2.NewRoutes(apiV2Group, t, u, tk, jwtManager, l)
}
```

你可以用任何其他 HTTP 框架來取代 [Fiber](https://github.com/gofiber/fiber)。

在 `router.go` 以及各處理常式方法的上方，有用來產生 Swagger 文件的註解，
由 [swag](https://github.com/swaggo/swag) 使用。

### `internal/entity`

商業邏輯的實體（entity，即模型），可以在任何一層使用。
也可以包含方法，例如用於驗證的方法。

### `internal/usecase`

商業邏輯。

- 方法依應用領域（共同的基礎）分組
- 每個群組都有自己的結構
- 一個檔案對應一個結構

Repository、webapi、rpc，以及其他商業邏輯結構，都會被注入到商業邏輯的結構中
（請參閱[相依性注入](#相依性注入)）。

#### `internal/repo/persistent`

Repository 是商業邏輯所操作的抽象儲存空間（資料庫）。

#### `internal/repo/webapi`

它是商業邏輯所操作的抽象 Web API。
例如，它可能是商業邏輯透過 REST API 存取的另一個微服務。
套件名稱會依用途而變。

### `pkg/rabbitmq`

RabbitMQ RPC 模式：

- RabbitMQ 內部沒有路由
- 使用 fanout exchange，並綁定 1 個獨佔佇列（exclusive queue），這是效能最好的設定
- 連線中斷時會自動重新連線

### `pkg/tracing`

OpenTelemetry 設定：以 OTLP/gRPC exporter 以及 W3C trace-context + baggage
propagator 來設定全域的 `TracerProvider`。停用追蹤時會安裝 no-op provider，讓追蹤程式碼可無條件執行。
請參閱[可觀測性](#可觀測性)。

## 相依性注入

為了讓商業邏輯不依賴外部套件，我們使用相依性注入（Dependency Injection）。

例如，透過 New 建構函式，我們把相依物件注入到商業邏輯的結構中。
這讓商業邏輯保持獨立（且可移植）。
我們可以替換介面的實作，而不必修改 `usecase` 套件。

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

這也讓我們能自動產生 mock（例如使用 [go.uber.org/mock](https://go.uber.org/mock)），
並輕鬆撰寫單元測試。

> 我們不綁定特定的實作，這樣就能隨時把某個元件換成另一個。
> 只要新元件實作了該介面，商業邏輯就完全不需要修改。

## 整潔架構

### 核心概念

程式設計師通常要等到大部分程式碼寫完之後，才會領悟應用程式的最佳架構。

> 好的架構能讓決策盡可能延後再做。

### 主要原則

相依反轉（Dependency Inversion，即 SOLID 中的那個 D）是相依性注入的原則。
相依的方向是由外層指向內層。
因此，商業邏輯與實體能夠獨立於系統的其他部分。

於是，應用程式被分為內層與外層兩層：

1. **商業邏輯**（Go 標準函式庫）。
2. **工具**（資料庫、伺服器、訊息代理（message broker）、任何其他套件與框架）。

![Clean Architecture](docs/img/layers-1.png)

**內層**是商業邏輯，應該保持乾淨。它應該：

- 不從外層匯入任何套件。
- 只使用標準函式庫提供的功能。
- 透過介面（interface）呼叫外層（！）。

商業邏輯對 Postgres 或特定的 Web API 一無所知。
商業邏輯只有一個介面，用來操作_抽象_的資料庫或_抽象_的 Web API。

**外層**有其他限制：

- 這一層的所有元件彼此並不知道對方的存在。要如何從一個工具呼叫另一個？不能直接呼叫，
  只能透過商業邏輯的內層。
- 所有對內層的呼叫都透過介面進行（！）。
- 資料以對商業邏輯方便的格式傳遞（`internal/entity`）。

例如，你需要從 HTTP（controller）存取資料庫。
HTTP 與資料庫都位於外層，也就是說它們彼此互不相識。
它們之間的溝通是透過 `usecase`（商業邏輯）進行的：

```
    HTTP > usecase
           usecase > repository (Postgres)
           usecase < repository (Postgres)
    HTTP < usecase
```

符號 > 與 < 表示透過介面跨越層與層的邊界。
圖片中呈現的是同樣的內容：

![Example](docs/img/example-http-db.png)

或是更複雜的商業邏輯：

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

### 分層

![Example](docs/img/layers-2.png)

### Clean Architecture 術語

- **Entities（實體）** 是商業邏輯所操作的結構。
  它們位於 `internal/entity` 資料夾。
  以 MVC 的說法，實體就是 model。
- **Use Cases（使用案例）** 是位於 `internal/usecase` 的商業邏輯。

與商業邏輯直接互動的那一層，通常稱為_基礎設施（infrastructure）_層。
它們可以是 repository `internal/repo/persistent`、外部 webapi `internal/repo/webapi`、任何 pkg，以及其他
微服務。
在這個範本中，_基礎設施_套件位於 `internal/repo` 之內。

入口點（entry point）要怎麼命名，可以依你的喜好選擇。可用的選項有：

- controller（我們採用的名稱）
- delivery
- transport
- gateways
- entrypoints
- primary
- input

### 額外的層

經典版的 [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
是為了建構大型單體式（monolithic）應用程式而設計，共有 4 層。

在原始版本中，外層又被再分成兩層，這兩層彼此之間同樣採用相依反轉（指向內側），
並透過介面溝通。

在邏輯複雜的情況下，內層同樣可以再分成兩層（並將介面分離出來）。

---

複雜的工具可以再拆分成額外的層。
不過，只有在真正必要時才應該增加層。

### 其他做法

除了 Clean Architecture 之外，_洋蔥式架構（Onion architecture）_與_六角形架構（Hexagonal，又稱 Ports and adapters）_也與它相似。
兩者都以相依反轉原則為基礎。
_Ports and adapters_ 與 _Clean Architecture_ 非常接近，差異主要在於用語。

## 類似的專案

- [https://github.com/bxcodec/go-clean-arch](https://github.com/bxcodec/go-clean-arch)
- [https://github.com/zhashkevych/courses-backend](https://github.com/zhashkevych/courses-backend)

## 實用連結

- [The Clean Architecture 文章](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Twelve factors](https://12factor.net/ru/)
