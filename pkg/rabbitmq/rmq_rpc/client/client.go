package client

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	rmqrpc "github.com/evrone/go-clean-template/pkg/rabbitmq/rmq_rpc"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/sync/errgroup"
)

// ErrConnectionClosed -.
var ErrConnectionClosed = errors.New("rmq_rpc client - Client - RemoteCall - Connection closed")

const (
	_defaultWaitTime = 5 * time.Second
	_defaultAttempts = 10
	_defaultTimeout  = 2 * time.Second
)

// Message -.
type Message struct {
	Queue         string
	Priority      uint8
	ContentType   string
	Body          []byte
	ReplyTo       string
	CorrelationID string
}

type pendingCall struct {
	done   chan struct{}
	status string
	body   []byte
}

// Client -.
type Client struct {
	ctx context.Context
	eg  *errgroup.Group

	conn           *rmqrpc.Connection
	serverExchange string
	notify         chan error
	stop           chan struct{}

	rw    sync.RWMutex
	calls map[string]*pendingCall

	timeout time.Duration
}

// New -.
func New(url, serverExchange, clientExchange string, opts ...Option) (*Client, error) {
	group, ctx := errgroup.WithContext(context.Background())
	group.SetLimit(1) // Run only one goroutine

	cfg := rmqrpc.Config{
		URL:      url,
		WaitTime: _defaultWaitTime,
		Attempts: _defaultAttempts,
	}

	c := &Client{
		ctx:            ctx,
		eg:             group,
		conn:           rmqrpc.New(clientExchange, cfg),
		serverExchange: serverExchange,
		notify:         make(chan error),
		stop:           make(chan struct{}),
		calls:          make(map[string]*pendingCall),
		timeout:        _defaultTimeout,
	}

	// Custom options
	for _, opt := range opts {
		opt(c)
	}

	err := c.conn.AttemptConnect()
	if err != nil {
		return nil, fmt.Errorf("rmq_rpc client - NewClient - c.conn.AttemptConnect: %w", err)
	}

	c.start()

	return c, nil
}

// Shutdown -.
func (c *Client) Shutdown() error {
	var shutdownErrors []error

	close(c.stop)

	// Wait for all goroutines to finish and get any error
	err := c.eg.Wait()
	if err != nil && !errors.Is(err, context.Canceled) {
		shutdownErrors = append(shutdownErrors, err)
	}

	// Close connection

	err = c.conn.Connection.Close()
	if err != nil {
		shutdownErrors = append(shutdownErrors, err)
	}

	return errors.Join(shutdownErrors...)
}

// RemoteCall -.
func (c *Client) RemoteCall(handler string, request, response interface{}) error {
	err := c.preRemoteCallWait()
	if err != nil {
		return fmt.Errorf("rmq_rpc client - Client - RemoteCall - c.preWait: %w", err)
	}

	corrID := uuid.New().String()

	err = c.publish(corrID, handler, request)
	if err != nil {
		return fmt.Errorf("rmq_rpc client - Client - RemoteCall - c.publish: %w", err)
	}

	call := &pendingCall{done: make(chan struct{})}

	c.addCall(corrID, call)
	defer c.deleteCall(corrID)

	err = c.remoteCallWait(call)
	if err != nil {
		return fmt.Errorf("rmq_rpc client - Client - RemoteCall - c.remoteCallWait: %w", err)
	}

	switch call.status {
	case rmqrpc.Success:
		err = json.Unmarshal(call.body, &response)
		if err != nil {
			return fmt.Errorf("rmq_rpc client - Client - RemoteCall - json.Unmarshal: %w", err)
		}

		return nil
	case rmqrpc.ErrBadHandler.Error():
		return rmqrpc.ErrBadHandler
	case rmqrpc.ErrInternalServer.Error():
		return rmqrpc.ErrInternalServer
	}

	return nil
}

func (c *Client) preRemoteCallWait() error {
	select {
	case <-c.ctx.Done():
		return c.ctx.Err()
	case err := <-c.notify:
		return fmt.Errorf("rmq_rpc client - Client - RemoteCall - c.notify: %w", err)
	case <-c.stop:
		return ErrConnectionClosed
	default:
	}

	return nil
}

func (c *Client) remoteCallWait(call *pendingCall) error {
	timeout := time.NewTimer(c.timeout)
	defer timeout.Stop()

	select {
	case <-c.ctx.Done():
		return c.ctx.Err()
	case err := <-c.notify:
		return fmt.Errorf("rmq_rpc client - Client - RemoteCall - c.notify: %w", err)
	case <-c.stop:
		return ErrConnectionClosed
	case <-timeout.C:
		return rmqrpc.ErrTimeout
	case <-call.done:
	}

	return nil
}

func (c *Client) start() {
	c.eg.Go(func() error {
		err := c.handleMessages()
		if err != nil {
			c.notify <- err

			close(c.notify)

			return err
		}

		return nil
	})
}

func (c *Client) handleMessages() error {
	for {
		select {
		case <-c.ctx.Done():
			return c.ctx.Err()
		case <-c.stop:
			return nil
		case d, opened := <-c.conn.Delivery:
			if !opened {
				err := c.reconnect()
				if err != nil {
					return err
				}

				break
			}

			c.serveCall(&d)
		}
	}
}

func (c *Client) reconnect() error {
	return c.conn.AttemptConnect()
}

func (c *Client) serveCall(d *amqp.Delivery) {
	defer c.ack(d, false)

	call, ok := c.getCall(d.CorrelationId)
	if !ok {
		return
	}

	call.status = d.Type
	call.body = d.Body
	close(call.done)
}

func (c *Client) getCall(corrID string) (*pendingCall, bool) {
	c.rw.RLock()
	defer c.rw.RUnlock()

	call, ok := c.calls[corrID]

	return call, ok
}

func (c *Client) addCall(corrID string, call *pendingCall) {
	c.rw.Lock()
	defer c.rw.Unlock()

	c.calls[corrID] = call
}

func (c *Client) deleteCall(corrID string) {
	c.rw.Lock()
	defer c.rw.Unlock()

	delete(c.calls, corrID)
}

func (c *Client) ack(d *amqp.Delivery, multiple bool) {
	d.Ack(multiple) //nolint:errcheck // we can't do anything with this error
}

func (c *Client) publish(corrID, handler string, request interface{}) error {
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

	err = c.conn.Channel.Publish(
		c.serverExchange,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: corrID,
			ReplyTo:       c.conn.ConsumerExchange,
			Type:          handler,
			Body:          requestBody,
		},
	)
	if err != nil {
		return fmt.Errorf("c.Channel.Publish: %w", err)
	}

	return nil
}
