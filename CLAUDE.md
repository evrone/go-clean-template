# go-clean-template

Go clean-architecture reference service. Three domains (`user`, `task`, `translation`) exposed
over four transports (REST/Fiber, gRPC, RabbitMQ RPC, NATS RPC) from one shared use-case layer.
Module path: `github.com/evrone/go-clean-template`. The Go version is declared in `go.mod`.

## Commands — drive everything through the Makefile

The `Makefile` loads `.env` (falling back to `.env.example`) and exports it into every target, and
several targets carry flags that matter. Running the underlying tool directly means running it with
different settings than CI. Use the target, not the tool.

| Task | Use this | Not this |
| --- | --- | --- |
| Install tool binaries (`swag`, `mockgen`, `migrate`, linters) | `make bin-deps` | `go install ...` |
| Unit tests | `make test` | `go test ./...` — the target adds `-race -covermode atomic` and scopes to `./internal/... ./pkg/...` |
| Integration tests | `make compose-up-integration-test` | `make integration-test` — see below |
| Lint | `make linter-golangci` | `golangci-lint run` |
| Format | `make format` | `gofmt` — the target runs `go fix`, `gofumpt`, and `gci` with the repo's import grouping |
| Regenerate mocks | `make mock` | `mockgen ...` |
| Regenerate Swagger | `make swag-v1` | `swag init` — the target passes `--parseDependency -g internal/controller/restapi/router.go` |
| Regenerate protobuf | `make proto-v1` | `protoc ...` |
| Tidy / verify modules | `make deps` | `go mod tidy` |
| Vulnerability scan | `make deps-audit` | `govulncheck ./...` |
| Start dependencies (Postgres, RabbitMQ, NATS) | `make compose-up` | `docker compose up` |
| Start the whole stack including the app | `make compose-up-all` | `docker compose up` |
| Tear down | `make compose-down` | `docker compose down` |
| Run the app locally | `make run` | `go run ./cmd/app` — the target regenerates docs and builds with `-tags migrate` |
| Create a migration | `make migrate-create <name>` | `migrate create ...` |
| Apply migrations | `make migrate-up` | `migrate -path ... up` |
| Full check before pushing | `make pre-commit` | running the steps by hand |

`make help` lists every target.

Three traps in these targets:

- **`make integration-test` is not the one you want.** It runs `go test ./integration-test/...` on
  the host, where the suite cannot resolve the container hostnames it needs, so it always fails.
  `make compose-up-integration-test` is the real entry point.
- **`make migrate-create <name>` prints an error after it succeeds.** The target reads the name via
  `$(word 2,$(MAKECMDGOALS))`, so `make` then tries to build `<name>` as a target and reports
  `No rule to make target`. The migration files are already created; ignore that line.
- **`make run` and `make pre-commit` depend on `swag-v1` and `proto-v1`**, so they need `swag` and
  `protoc` on `PATH`. Run `make bin-deps` first (`protoc` itself is not installed by it).

Never claim a change is done without `make format`, `make linter-golangci` and `make test` passing.

## Dependency rule

```
cmd/app → internal/app → internal/controller/*  ─┐
                       → internal/repo/*        ─┤→ internal/usecase (interfaces)
                                                 └→ internal/entity
```

- `internal/entity` — domain types and sentinel errors. Imports nothing from this module.
- `internal/usecase/contracts.go` — the interfaces controllers call. `internal/usecase/<domain>/`
  implements them.
- `internal/repo/contracts.go` — the interfaces use cases call. `internal/repo/persistent/<domain>/`
  (Postgres) and `internal/repo/webapi/` (outbound HTTP) implement them.
- `internal/controller/<transport>/v1/` — one package per transport, each with its own
  `request/` and `response/` DTOs. Controllers never import each other.
- `pkg/` — transport-agnostic infrastructure (servers, logger, jwt, postgres, tracing). Must not
  import `internal/`.

Inner layers never import outer ones. Wiring happens exactly once, in `internal/app/app.go`
(`initUseCases`, `initServers`).

## Where code goes

Decide by asking what the code knows about:

| The code knows about… | It belongs in | Shape |
| --- | --- | --- |
| nothing but the domain | `internal/entity/<name>.go` | struct + methods (`Task.Transition`, `TaskStatus.Valid`) |
| a domain rule that spans repositories | `internal/usecase/<domain>/<domain>.go` | method on `UseCase` |
| SQL, a table, a driver error code | `internal/repo/persistent/<domain>/<domain>.go` | method on `Repo` |
| an outbound HTTP API | `internal/repo/webapi/<name>.go` | method on the webapi struct |
| an HTTP status, a gRPC code, a message envelope | `internal/controller/<transport>/v1/<domain>.go` | handler |
| a server, pool, client or middleware with no domain knowledge | `pkg/<name>/` | reusable package |

Layout rules:

- One file per domain per layer, named after the domain (`task.go`, `user.go`, `translation.go`),
  plus a `tracing.go` in every `usecase/<domain>/` and `repo/persistent/<domain>/` package.
- Package name matches the directory name. `internal/repo/persistent/translation/` is the one
  violation in the tree — it declares `package persistent`, which is why `internal/app/app.go`
  imports the three repositories under aliases. Follow `task/` and `user/`, not `translation/`.
- Constructors are always `New`, return the *interface* from `contracts.go`, and take their
  dependencies as arguments. No package-level state — `gochecknoglobals` and `gochecknoinits` are on.
- Unexported package constants are prefixed with an underscore (`_defaultEntityCap`,
  `_tracerName`). `mnd` rejects bare numeric literals, so add a constant rather than inlining one.

## The handler contract

Every controller method, on every transport, does the same six things in the same order. Copy the
shape from a neighbouring handler in the same transport.

1. **Get the caller.** REST: `ctx.Locals("userID").(string)`. gRPC: `grpcmw.UserIDFromContext(ctx)`.
   AMQP / NATS: `extractUserID(d, r.j)`. A failure here is an auth error, not a 500.
2. **Decode into a transport-local DTO** from `<transport>/v1/request/`. Never decode into an
   `entity` type, and never pass a `request.*` type into a use case.
3. **Validate the DTO** with the injected `*validator.Validate` (`r.v.Struct(req)`). Structural
   validation (required, min, max, oneof) lives here. Domain validation — anything that needs to
   know the rules, like a status transition — lives on the entity or the use case and must not be
   duplicated in the controller.
4. **Call the use case** with the request context: `ctx.UserContext()` in Fiber (**not**
   `ctx.Context()`, which drops the trace), the handler's `ctx` everywhere else. `context.Context`
   is the first parameter of every method that crosses a layer.
5. **Map the error.** `errors.Is` against the `entity.Err*` sentinels → a transport status. Log the
   wrapped error, return a generic message to the caller.
6. **Encode the response.** Return an `entity` type directly only when its JSON shape is already
   the API shape — `entity.User` is safe this way because `PasswordHash` is tagged `json:"-"`.
   Anything with a wrapper, a projection or a different field set gets a type in
   `<transport>/v1/response/` (`response.TaskList`, `response.Token`). gRPC always converts through
   `response.New*Response`.

AMQP and NATS handlers are closures returning that transport's `server.CallHandler` —
`func(ctx, *amqp.Delivery) (any, error)` and `func(ctx, *nats.Msg) (any, error)`. Each transport
has its own `extractUserID` in `v1/auth.go` taking its own message type. Handlers return a plain
wrapped `error`; the RPC server turns it into a generic `ErrInternalServer` / `ErrBadHandler`
reply. Do not hand-build an error envelope there.

## Conventions that will bite you

**Traced decorators.** Every repo and use-case constructor returns the *interface*, wrapped:
`func New(...) usecase.Task { return newTraced(&UseCase{...}) }`. The wrapper lives in
`tracing.go` next to each implementation. Adding a method to a contract means implementing it in
three places: the interface, the struct, and the `traced*` decorator. Skipping the decorator is a
compile error that reads as an unrelated interface mismatch.

**Mocks are generated, not written.** `internal/usecase/mocks_*_test.go` come from `go:generate`
directives in the two `contracts.go` files. After any contract change run `make mock` before
running tests.

**Error style.** Domain failures are sentinel errors in `internal/entity/errors.go`
(`ErrTaskNotFound`, `ErrTaskForbidden`, `ErrInvalidTransition`, …). Every other error is wrapped
with the layer path: `fmt.Errorf("TaskUseCase - Create - uc.repo.Store: %w", err)`. The `err113`
linter rejects `errors.New`/`fmt.Errorf` for new dynamic errors — add a sentinel to `entity`
instead. Controllers are the only place that translates a sentinel into a transport status.

**Config is environment-only.** `config/config.go` parses env vars with `caarlos0/env`. A new
setting needs the struct field, an entry in `.env.example`, and — if the container needs it — the
`environment` block in `docker-compose.yml`. Fields tagged `required` make the app fail at
startup.

**The linter is strict and non-obvious.** `.golangci.yml` enables ~45 linters. The ones that
reject otherwise-fine code: `wsl_v5` and `nlreturn` (blank line required before `return` and
around blocks), `funlen` (65 lines / 40 statements), `gocyclo` (10), `gocognit` (15), `mnd`,
`gochecknoglobals`, `exhaustive` (switches over enum types need every case or a `default`), `dupl`
(threshold 100 — copy-pasting a handler across transports trips it).

## Generated files — do not hand-edit

- `docs/docs.go`, `docs/swagger.json`, `docs/swagger.yaml` ← `make swag-v1`, driven by the
  annotation comments above REST handlers and on `NewRouter` in `internal/controller/restapi/router.go`.
- `docs/proto/v1/*.pb.go` ← `make proto-v1` from the `.proto` files beside them.
- `internal/usecase/mocks_*_test.go` ← `make mock`.

## READMEs must be updated as a set

`README.md` (English) is canonical. `README_RU.md`, `README_CN.md` and `README_zh-TW.md` are full translations that
mirror it **section for section, table row for table row** — same headings in the same order, same
tables, same code blocks, same links. Changing one and not the others is the most common
regression in this repo.

Update all four whenever a change touches what they document:

| Change | What to update in every README |
| --- | --- |
| New or changed REST / gRPC endpoint | the domain table under `Domains` (`Домены` / `领域` / `領域`) |
| New domain | the `Overview` bullet list, `Content` index, and a new `Domains` subsection |
| New or renamed env var | the config table under `Observability` if tracing-related; otherwise the `config` section |
| New transport, server, port or service URL | the `Quick start` service list |
| New `make` target used in the getting-started flow | the `Quick start` code blocks |
| Renamed or moved package under `internal/` or `pkg/` | the matching `Project structure` subsection |

Rules for the translations:

- Translate the prose; leave identifiers, paths, `make` targets, URLs, env var names, code blocks
  and route strings verbatim in English.
- Headings are translated, so the `Content` index anchors are localized too (`#домены`, `#领域`).
  Adding or renaming a section means updating that file's own anchors — an anchor copied from the
  English file will not resolve.
- Keep the heading level identical across the three files.
- If you genuinely cannot produce a translation, still add the section in English to the other two
  files so the structure stays aligned, and say so in the change description.

## Adding a use case to an existing domain

1. Add the method to `internal/usecase/contracts.go`.
2. Implement it in `internal/usecase/<domain>/<domain>.go` and mirror it in the same package's
   `tracing.go`.
3. If it needs new persistence, repeat the same three steps in `internal/repo/contracts.go`,
   `internal/repo/persistent/<domain>/`, and that package's `tracing.go`.
4. `make mock`, then extend `internal/usecase/<domain>_test.go` (table-driven, gomock, `t.Parallel()` —
   `paralleltest` enforces it).
5. Expose it on whichever transports need it:
   - REST: handler in `internal/controller/restapi/v1/<domain>.go` + route in that package's
     `router.go` + swagger annotations, then `make swag-v1`.
   - gRPC: `.proto` change → `make proto-v1` → controller method in
     `internal/controller/grpc/v1/`.
   - AMQP / NATS: handler returning a `CallHandler` + a `"v1.<domain>.<action>"` key in the
     transport's `v1/router.go`.

Transports are independent by design; adding a method to only one is a legitimate choice, not an
oversight — but say which ones you covered.

## Adding a whole domain

Work outward from the entity. Every step has a counterpart in the existing `task` domain — read it
before writing.

1. **Entity** — `internal/entity/<domain>.go` with the type and its invariants; new sentinel errors
   go in `internal/entity/errors.go`.
2. **Migration** — `make migrate-create create_<table>`, then fill both the `.up.sql` and the
   `.down.sql`. A missing down-migration is a review failure.
3. **Repo contract** — add the interface to `internal/repo/contracts.go`.
4. **Repo implementation** — `internal/repo/persistent/<domain>/<domain>.go` (`package <domain>`,
   Squirrel queries, driver errors mapped to sentinels) plus `tracing.go` next to it, and a `New`
   returning the interface wrapped in `newTraced`.
5. **Use-case contract** — add the interface to `internal/usecase/contracts.go`.
6. **Use-case implementation** — `internal/usecase/<domain>/<domain>.go` plus its `tracing.go`,
   same `New` pattern.
7. **`make mock`**, then `internal/usecase/<domain>_test.go` against the generated mocks.
8. **Wiring** — add a field to `useCases` in `internal/app/app.go`, construct repo and use case in
   `initUseCases`, and thread the use case through `initServers` into each router you are exposing.
9. **Transports** — for each one: DTOs in `<transport>/v1/request/` and `response/`, handlers in
   `<transport>/v1/<domain>.go`, registration in that transport's `v1/router.go`, and the use case
   added to the `V1` / controller struct in `<transport>/v1/controller.go`. gRPC also needs
   `docs/proto/v1/<domain>.proto` and `make proto-v1`.
10. **Docs** — swagger annotations on the REST handlers, `make swag-v1`, and the README updates
    listed above in all four READMEs.
11. **Integration test** — `integration-test/<domain>_test.go` exercising the domain over the
    transports you exposed.

## API versioning

Routes are grouped under a version package (`<transport>/v1/`). A v2 means a new sibling package
with the same layout, registered alongside v1 in the transport's top-level `router.go` — REST adds
an `app.Group("/v2")`, gRPC registers a second set of `NewXRoutes`, AMQP/NATS add more
`"v2.<domain>.<action>"` keys to the same route map. v1 stays untouched. `README.md` documents this
with worked examples per transport; follow it rather than inventing a scheme.

## Database

`migrations/` holds golang-migrate pairs (`<timestamp>_<name>.up.sql` / `.down.sql`). Create with
`make migrate-create <name>`, apply with `make migrate-up`. The app applies them itself at startup
only when built with the `migrate` build tag (`internal/app/migrate.go`) — that is what `make run`
and the Dockerfile do. Queries are built with Squirrel via the embedded `*postgres.Postgres`;
there is no ORM and no raw string concatenation.

## Tests

- Unit: `internal/...` and `pkg/...`, run by `make test` with `-race`. Use-case tests use the
  generated gomock mocks, are table-driven, and call `t.Parallel()` (`paralleltest` enforces it).
- Integration: `integration-test/` runs *inside* the docker network — it resolves the service as
  host `app` and talks to `rabbitmq` / `nats` by container name, so it fails on the host machine.
  Always run it through `make compose-up-integration-test`.
