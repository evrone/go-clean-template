include .env
export

run: swag
	go mod download && GIN_MODE=debug CGO_ENABLED=0 go run ./cmd/app

compose-up:
	docker-compose up --build -d --remove-orphans && docker-compose logs -f

compose-down:
	docker-compose down --remove-orphans

test:
	go test -v -cover -race ./...

mock:
	mockery --all -r --case snake

migrate-create:
	migrate create -ext sql -dir migrations 'migrate_name'

migrate:
	migrate -path migrations -database '$(PG_URL)?sslmode=disable' up

swag:
	swag init -g internal/app/app.go

.PHONY: run, compose-up, compose-down, test, mock, migrate-create, migrate, swag