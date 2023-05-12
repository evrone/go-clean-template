package db

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/evrone/go-clean-template/config"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"time"
)

func MustStartRMQContainer(ctx context.Context, cfg *config.Config) {

	port, err := nat.NewPort("", "5672")
	if err != nil {
		panic(fmt.Errorf("failed to build port: %v", err))
	}

	timeout := 5 * time.Minute // Default timeout
	tag := "3.11.15"

	req := testcontainers.ContainerRequest{
		Image:        fmt.Sprintf("rabbitmq:%s", tag),
		ExposedPorts: []string{string(port)},
		WaitingFor:   wait.ForListeningPort(port).WithStartupTimeout(timeout),
	}

	rmqContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(fmt.Errorf("failed to start container: %v", err))
	}

	host, err := rmqContainer.Host(ctx)
	if err != nil {
		panic(fmt.Errorf("failed to get container host: %v", err))
	}

	realPort, err := rmqContainer.MappedPort(ctx, port)
	if err != nil {
		panic(fmt.Errorf("failed to get exposed container port: %v", err))
	}

	cfg.RMQ.URL = fmt.Sprintf("amqp://%s:%s", host, realPort.Port())
}
