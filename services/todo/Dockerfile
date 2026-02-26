# Build context must be the repo root (to access common_packages/)
# docker build -f services/todo/Dockerfile .

# Step 1: Modules caching
FROM golang:1.26-alpine3.23 AS modules

COPY go.work go.work.sum /app/
COPY common_packages/pkg/go.mod common_packages/pkg/go.sum /app/common_packages/pkg/
COPY common_packages/middleware/go.mod common_packages/middleware/go.sum /app/common_packages/middleware/
COPY services/todo/go.mod services/todo/go.sum /app/services/todo/

WORKDIR /app

RUN GONOSUMDB=evrone.local go mod download

# Step 2: Builder
FROM golang:1.26-alpine3.23 AS builder

COPY --from=modules /go/pkg /go/pkg
COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -tags migrate -o /bin/app ./services/todo/cmd/todo

# Step 3: Final
FROM scratch

COPY --from=builder /app/services/todo/migrations /migrations
COPY --from=builder /bin/app /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

CMD ["/app"]
