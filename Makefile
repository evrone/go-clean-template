include .env
export

LOCAL_BIN:=$(CURDIR)/bin
PATH:=$(LOCAL_BIN):$(PATH)

# HELP =================================================================================================================
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help

help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

compose-up: ### Run docker-compose
	docker-compose up --build -d postgres rabbitmq && docker-compose logs -f
.PHONY: compose-up

compose-down: ### Down docker-compose
	docker-compose down --remove-orphans
.PHONY: compose-down

swag-v1: ### swag init
	swag init -g internal/controller/http/v1/router.go
.PHONY: swag-v1

run: swag-v1 ### swag run
	go mod tidy && go mod download && \
	DISABLE_SWAGGER_HTTP_HANDLER='' GIN_MODE=debug CGO_ENABLED=0 go run ./cmd/app
.PHONY: run

docker-rm-volume: ### remove docker volume
	docker volume rm go-clean-template_pg-data
.PHONY: docker-rm-volume

linter-golangci: ### check by golangci linter
	golangci-lint run
.PHONY: linter-golangci

linter-hadolint: ### check by hadolint linter
	find . -name 'Dockerfile' | xargs hadolint
.PHONY: linter-hadolint

linter-yaml:
	yamllint . -s
.PHONY: linter-yaml

linter-dotenv: ### check by dotenv linter
	dotenv-linter
.PHONY: linter-dotenv

lint: linter-golangci linter-hadolint linter-yaml linter-dotenv ### run all linters
.PHONY: lint

test: ### run all tests including slow running system (e.g. system-tests)
	go test --tags=system -v -cover -covermode atomic -coverprofile=coverage.txt ./internal/... ./pkg/...
.PHONY: test

test-fast: ### run fast tests only
	go test -v -cover ./internal/... ./pkg/...
.PHONY: test-fast

mock: ### run mockgen
	mockgen -source ./internal/usecase/interfaces.go -package usecase_test > ./internal/usecase/mocks_test.go
.PHONY: mock

migrate-create:  ### create new migration
	migrate create -ext sql -dir migrations 'migrate_name'
.PHONY: migrate-create

migrate-up: ### migration up
	migrate -path migrations -database '$(PG_URL)?sslmode=disable' up
.PHONY: migrate-up

setup-mac: ### setup mac os dependencies to run all tasks
	brew install openapi-generator
.PHONY: setup-mac

generate: ### Generate server files based on OpenAPI specs
	openapi-generator generate -i docs/openapi.yaml -g go-gin-server  -o ./internal/interfaces/rest/v1 && rm -rf ./internal/interfaces/rest/v1/main.go ./internal/interfaces/rest/v1/go.mod ./internal/interfaces/rest/v1/Dockerfile ./internal/interfaces/rest/v1/go.sum ./internal/interfaces/rest/v1/api
.PHONY: generate