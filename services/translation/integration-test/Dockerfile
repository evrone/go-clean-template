# Build context must be the repo root (to access common_packages/)
# Step 1: Modules caching
FROM golang:1.26-alpine3.23 AS modules

COPY go.work go.work.sum /app/
COPY common_packages/pkg/go.mod common_packages/pkg/go.sum /app/common_packages/pkg/
COPY common_packages/middleware/go.mod common_packages/middleware/go.sum /app/common_packages/middleware/
COPY services/translation/go.mod services/translation/go.sum /app/services/translation/

WORKDIR /app

RUN GONOSUMDB=evrone.local go mod download

# Step 2: Tests
FROM golang:1.26-alpine3.23

COPY --from=modules /go/pkg /go/pkg
COPY . /app

WORKDIR /app

RUN go env -w CGO_ENABLED=0 \
    && go env -w GOOS=linux  \
    && go env -w GOARCH=amd64

CMD ["go", "test", "-v", "./services/translation/integration-test/..."]
