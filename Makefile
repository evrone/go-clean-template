include .env
export

.PHONY: compose-up
.PHONY: compose-up-integration-test
.PHONY: compose-down
.PHONY: swag
.PHONY: run
.PHONY: docker-rm-volume
.PHONY: linter-golangci
.PHONY: linter-hadolint
.PHONY: linter-dotenv
.PHONY: test
.PHONY: integration-test
.PHONY: mock
.PHONY: migrate-create
.PHONY: migrate-up

compose-up:
	docker-compose up --build -d postgres rabbitmq && docker-compose logs -f

compose-up-integration-test:
	docker-compose up --build --abort-on-container-exit --exit-code-from integration

compose-down:
	docker-compose down --remove-orphans

swag:
	swag init -g internal/app/app.go

#run: swag
#	go mod tidy && go mod download && GIN_MODE=debug CGO_ENABLED=0 go run -tags migrate ./cmd/app

run:
	go mod tidy && go mod download && GIN_MODE=release CGO_ENABLED=0 go run -tags migrate ./cmd/app

docker-rm-volume:
	docker volume rm go-service-template_pg-data

linter-golangci:
	golangci-lint run

linter-hadolint:
	git ls-files --exclude='Dockerfile*' --ignored | xargs hadolint

linter-dotenv:
	dotenv-linter

test:
	go test -v -cover -race ./internal/...

integration-test:
	go clean -testcache && HOST=localhost:8080 go test -v ./integration-test/...

mock:
	mockery --all -r --case snake

migrate-create:
	migrate create -ext sql -dir migrations 'migrate_name'

migrate-up:
	migrate -path migrations -database '$(PG_URL)?sslmode=disable' up
