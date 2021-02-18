include .env
export

swag:
	swag init -g internal/app/app.go

run: swag
	go mod tidy && go mod download && GIN_MODE=debug CGO_ENABLED=0 go run -tags migrate ./cmd/app

compose-up-db:
	docker-compose up --build -d --remove-orphans db && docker-compose logs -f

compose-up:
	docker-compose up --build -d --remove-orphans && docker-compose logs -f

compose-down:
	docker-compose down --remove-orphans

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

test-integration:
	go clean -testcache && HOST=localhost:8080 go test -v ./integration-test/...

mock:
	mockery --all -r --case snake

migrate-create:
	migrate create -ext sql -dir migrations 'migrate_name'

migrate:
	migrate -path migrations -database '$(PG_URL)?sslmode=disable' up

.PHONY: swag, run, compose-up-db, compose-up, compose-down, docker-rm-pg-data, linter-golangci, linter-hadolint
.PHONY: linter-dotenv, test, test-integration, mock, migrate-create, migrate