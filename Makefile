SHELL := /bin/bash

BASE_STACK     = docker compose -f docker-compose.yml
INT_TEST_STACK = $(BASE_STACK) -f docker-compose-integration-test.yml

# HELP =================================================================================================================
.PHONY: help
help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Docker Compose

compose-up: ## Start infrastructure only (db, rabbitmq, nats)
	$(BASE_STACK) up --build -d db rabbitmq nats && docker compose logs -f
.PHONY: compose-up

compose-up-all: ## Start all services (translation + todo + infra)
	$(BASE_STACK) up --build -d
.PHONY: compose-up-all

compose-up-integration-test: ## Run integration tests via docker compose
	$(INT_TEST_STACK) up --build --abort-on-container-exit --exit-code-from integration-test
.PHONY: compose-up-integration-test

compose-down: ## Stop and remove all containers
	$(INT_TEST_STACK) down --remove-orphans
.PHONY: compose-down

docker-rm-volume: ## Remove persistent docker volumes
	docker volume rm go-clean-template_db_data go-clean-template_rabbitmq_data go-clean-template_nats_data
.PHONY: docker-rm-volume

##@ Run (local — requires infra via compose-up)

run-translation: ## Run translation service locally with auto-migration
	cd services/translation && set -a && . ./.env && set +a && \
	CGO_ENABLED=0 go run -tags migrate ./cmd/app
.PHONY: run-translation

run-todo: ## Run todo service locally with auto-migration
	cd services/todo && set -a && . ./.env && set +a && \
	CGO_ENABLED=0 go run -tags migrate ./cmd/todo
.PHONY: run-todo

##@ Migrations

migrate-up-translation: ## Apply translation migrations (requires PG running)
	cd services/translation && set -a && . ./.env && set +a && \
	migrate -path migrations -database "$$PG_URL?sslmode=disable" up
.PHONY: migrate-up-translation

migrate-up-todo: ## Apply todo migrations (requires PG running)
	cd services/todo && set -a && . ./.env && set +a && \
	migrate -path migrations -database "$$PG_URL?sslmode=disable" up
.PHONY: migrate-up-todo

migrate-create-translation: ## Create translation migration: make migrate-create-translation NAME=add_column
	migrate create -ext sql -dir services/translation/migrations $(NAME)
.PHONY: migrate-create-translation

migrate-create-todo: ## Create todo migration: make migrate-create-todo NAME=add_column
	migrate create -ext sql -dir services/todo/migrations $(NAME)
.PHONY: migrate-create-todo

##@ Code Generation

swag-translation: ## Regenerate swagger docs for translation service
	cd services/translation && \
	~/go/bin/swag init -g internal/controller/restapi/router.go -o docs/
.PHONY: swag-translation

proto-translation: ## Regenerate gRPC stubs for translation service
	cd services/translation && \
	protoc --go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
		docs/proto/v1/*.proto
.PHONY: proto-translation

mock-translation: ## Regenerate mocks for translation service
	cd services/translation && \
	mockgen -source ./internal/repo/contracts.go -package usecase_test > ./internal/usecase/mocks_repo_test.go && \
	mockgen -source ./internal/usecase/contracts.go -package usecase_test > ./internal/usecase/mocks_usecase_test.go
.PHONY: mock-translation

mock-todo: ## Regenerate mocks for todo service
	cd services/todo && \
	mockgen -source ./internal/repo/contracts.go -package usecase_test > ./internal/usecase/mocks_repo_test.go && \
	mockgen -source ./internal/usecase/contracts.go -package usecase_test > ./internal/usecase/mocks_usecase_test.go
.PHONY: mock-todo

mock: mock-translation mock-todo ## Regenerate all mocks
.PHONY: mock

##@ Dependencies

deps: ## Run go mod tidy for all workspace modules
	cd common_packages/pkg        && go mod tidy
	cd common_packages/middleware && go mod tidy
	cd services/translation       && go mod tidy
	cd services/todo              && go mod tidy
.PHONY: deps

deps-audit: ## Check all modules for vulnerabilities
	cd services/translation && govulncheck ./...
	cd services/todo        && govulncheck ./...
.PHONY: deps-audit

##@ Testing

test-translation: ## Run translation unit tests
	cd services/translation && \
	go test -v -race -covermode atomic -coverprofile=coverage.txt ./internal/...
.PHONY: test-translation

test-todo: ## Run todo unit tests
	cd services/todo && \
	go test -v -race -covermode atomic -coverprofile=coverage.txt ./internal/...
.PHONY: test-todo

test: test-translation test-todo ## Run all unit tests
.PHONY: test

integration-test: ## Run translation integration tests locally
	cd services/translation && go clean -testcache && go test -v ./integration-test/...
.PHONY: integration-test

##@ Linting & Formatting

linter-golangci: ## Run golangci-lint across all services
	cd services/translation && golangci-lint run
	cd services/todo        && golangci-lint run
.PHONY: linter-golangci

linter-hadolint: ## Lint all Dockerfiles with hadolint
	git ls-files --exclude='Dockerfile*' --ignored | xargs hadolint
.PHONY: linter-hadolint

linter-dotenv: ## Lint .env.example files with dotenv-linter
	dotenv-linter services/translation/.env.example services/todo/.env.example
.PHONY: linter-dotenv

format: ## Format all Go code
	gofumpt -l -w .
	gci write . --skip-generated -s standard -s default
.PHONY: format

##@ Tools

bin-deps: ## Install required dev tools
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install go.uber.org/mock/mockgen@latest
	go install mvdan.cc/gofumpt@latest
	go install github.com/daixiang0/gci@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
.PHONY: bin-deps

pre-commit: swag-translation proto-translation mock format linter-golangci test ## Run all pre-commit checks
.PHONY: pre-commit
