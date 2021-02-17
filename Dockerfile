# Step 1: Modules caching
FROM golang:1.15.7-alpine3.13 as modules
COPY go.mod go.sum /modules/
RUN cd /modules && go mod download


# Step 2: Builder
FROM golang:1.15.7-alpine3.13 as builder
RUN apk add --no-cache ca-certificates

COPY --from=modules /go/pkg /go/pkg

COPY . /app
WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -tags migrate -o /bin/app ./cmd/app


# Step 3: Final
FROM scratch
COPY --from=builder /app/migrations /migrations
COPY --from=builder /bin/app /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
CMD ["/app"]