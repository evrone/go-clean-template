up:
	docker-compose up --build -d --remove-orphans && docker-compose logs -f

down:
	docker-compose down --remove-orphans

test:
	go test -v -cover -race ./...

mock:
	mockery --all

migrate:
	migrate -path migrations -database 'postgres://user:pass@localhost:5432/postgres?sslmode=disable' up

create:
	migrate create -ext sql -dir migrations -seq 'migrate_name'

.PHONY: build, up, down, test, mock, migrate, create