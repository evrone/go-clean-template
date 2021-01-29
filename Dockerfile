# Step 1: Modules
FROM golang:1.15 as modules
COPY go.mod go.sum /modules/
RUN cd /modules && go mod download


# Step 2: Builder
FROM golang:1.15 as builder
RUN update-ca-certificates

COPY --from=modules /go/pkg /go/pkg

RUN useradd -u 2000 appuser

RUN mkdir -p /app
COPY . /app
WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /bin/app ./app


# Step 3: Final
FROM scratch
COPY --from=builder /etc/passwd /etc/passwd
USER appuser

COPY --from=builder /bin/app /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
CMD ["/app"]
