package client

import (
	"context"
	"errors"
	"fmt"
	"time"

	natsrpc "github.com/evrone/go-clean-template/pkg/nats/nats_rpc"
	"github.com/goccy/go-json"
	"github.com/nats-io/nats.go"
)

const (
	_defaultWaitTime = 5 * time.Second
	_defaultAttempts = 10
	_defaultTimeout  = 2 * time.Second
)

// Client -.
type Client struct {
	subject    string
	connection *nats.Conn

	timeout time.Duration
}

// New -.
func New(
	url string,
	serverSubject string,
	opts ...Option,
) (*Client, error) {
	connection, err := nats.Connect(
		url,
		nats.ReconnectWait(_defaultWaitTime),
		nats.MaxReconnects(_defaultAttempts),
		nats.Timeout(_defaultWaitTime),
	)
	if err != nil {
		return nil, fmt.Errorf("nats_rpc client - NewClient - nats.Connect: %w", err)
	}

	c := &Client{
		subject:    serverSubject,
		connection: connection,
		timeout:    _defaultTimeout,
	}

	// Custom options
	for _, opt := range opts {
		opt(c)
	}

	c.connection = connection

	return c, nil
}

// Shutdown -.
func (c *Client) Shutdown() error {
	c.connection.Close()

	return nil
}

// RemoteCall -.
func (c *Client) RemoteCall(handler string, request, response interface{}) error {
	var (
		requestBody []byte
		err         error
	)

	if request != nil {
		requestBody, err = json.Marshal(request)
		if err != nil {
			return err
		}
	}

	requestMessage := nats.Msg{
		Subject: c.subject,
		Header: nats.Header{
			"Handler": []string{handler},
		},
		Data: requestBody,
	}

	message, err := c.connection.RequestMsg(&requestMessage, c.timeout)
	if errors.Is(err, context.DeadlineExceeded) {
		return natsrpc.ErrTimeout
	}

	if err != nil {
		return fmt.Errorf("nats_rpc client - Client - RemoteCall - c.connection.Conn.Request: %w", err)
	}

	switch message.Header.Get("Status") {
	case natsrpc.Success:
		err = json.Unmarshal(message.Data, &response)
		if err != nil {
			return fmt.Errorf("nats_rpc client - Client - RemoteCall - json.Unmarshal: %w", err)
		}
	case natsrpc.ErrBadHandler.Error():
		return natsrpc.ErrBadHandler
	case natsrpc.ErrInternalServer.Error():
		return natsrpc.ErrInternalServer
	}

	return nil
}
