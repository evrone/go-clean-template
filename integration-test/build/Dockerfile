# Step 1: Modules caching
FROM golang:1.15.7-alpine3.13 as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Step 2: Builder
FROM golang:1.15.7-alpine3.13 as builder
RUN apk --no-cache add bash=5.1.0-r0
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app
CMD ["./integration-test/build/wait.sh"]