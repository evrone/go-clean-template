include .env
export

run:
	go mod download && GIN_MODE=debug CGO_ENABLED=0 go run ./cmd/app

up:
	docker-compose up --build -d --remove-orphans && docker-compose logs -f

down:
	docker-compose down --remove-orphans

test:
	go test -v -cover -race ./...

mock:
	mockery --all -r --case snake

migrate:
	migrate -path migrations -database '$(PG_URL)?sslmode=disable' up

create:
	migrate create -ext sql -dir migrations 'migrate_name'

.PHONY: run, up, down, test, mock, migrate, create