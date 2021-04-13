# Step 1: Modules caching
FROM golang:1.16.2-alpine3.13 as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Step 2: Tests
FROM golang:1.16.2-alpine3.13
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app

RUN go env -w CGO_ENABLED=0
RUN go env -w GOOS=linux
RUN go env -w GOARCH=amd64

CMD ["go", "test", "-v", "./integration-test/..."]