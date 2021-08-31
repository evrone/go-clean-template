include .env.example
export

compose-up:
	docker-compose up --build -d postgres rabbitmq && docker-compose logs -f
.PHONY: compose-up

compose-up-integration-test:
	docker-compose up --build --abort-on-container-exit --exit-code-from integration
.PHONY: compose-up-integration-test

compose-down:
	docker-compose down --remove-orphans
.PHONY: compose-down

swag-v1:
	swag init -g internal/controller/http/v1/router.go
.PHONY: swag-v1

run: swag-v1
	go mod tidy && go mod download && \
	DISABLE_SWAGGER_HTTP_HANDLER='' GIN_MODE=debug CGO_ENABLED=0 go run -tags migrate ./cmd/app
.PHONY: run

docker-rm-volume:
	docker volume rm go-clean-template_pg-data
.PHONY: docker-rm-volume

linter-golangci:
	golangci-lint run
.PHONY: linter-golangci

linter-hadolint:
	git ls-files --exclude='Dockerfile*' --ignored | xargs hadolint
.PHONY: linter-hadolint

linter-dotenv:
	dotenv-linter
.PHONY: linter-dotenv

test:
	go test -v -cover -race ./internal/...
.PHONY: test

integration-test:
	go clean -testcache && go test -v ./integration-test/...
.PHONY: integration-test

mock:
	mockery --all -r --case snake
.PHONY: mock

migrate-create:
	migrate create -ext sql -dir migrations 'migrate_name'
.PHONY: migrate-create

migrate-up:
	migrate -path migrations -database '$(PG_URL)?sslmode=disable' up
.PHONY: migrate-up
