include .env
export

swag:
	swag init -g internal/app/app.go

run: swag
	go mod download && GIN_MODE=debug CGO_ENABLED=0 go run ./cmd/app

run-with-migrate:
	go mod download && CGO_ENABLED=0 go run -tags migrate ./cmd/app

compose-up-db:
	docker-compose up --build -d --remove-orphans db && docker-compose logs -f

compose-up:
	docker-compose up --build -d --remove-orphans && docker-compose logs -f

compose-down:
	docker-compose down --remove-orphans

lint:
	golangci-lint run

hadolint:
	git ls-files --exclude='Dockerfile*' --ignored | xargs hadolint

test:
	go test -v -cover -race ./internal/...

integration-test:
	go clean -testcache && HOST=localhost:8080 go test -v ./integration-test/...

mock:
	mockery --all -r --case snake

migrate-create:
	migrate create -ext sql -dir migrations 'migrate_name'

migrate:
	migrate -path migrations -database '$(PG_URL)?sslmode=disable' up

.PHONY: swag, run, run-with-migrate, compose-up-db, compose-up, compose-down
.PHONY: lint, hadolint, test, integration-test, mock, migrate-create, migrate