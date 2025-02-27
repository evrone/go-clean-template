package server

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/evrone/go-clean-template/pkg/logger"
	rmqrpc "github.com/evrone/go-clean-template/pkg/rabbitmq/rmq_rpc"
	"github.com/streadway/amqp"
)

const (
	_defaultWaitTime = 5 * time.Second
	_defaultAttempts = 10
	_defaultTimeout  = 2 * time.Second
)

// CallHandler -.
type CallHandler func(*amqp.Delivery) (interface{}, error)

// Server -.
type Server struct {
	conn   *rmqrpc.Connection
	error  chan error
	stop   chan struct{}
	router map[string]CallHandler

	timeout time.Duration

	logger logger.Interface
}

// New -.
func New(url, serverExchange string, router map[string]CallHandler, l logger.Interface, opts ...Option) (*Server, error) {
	cfg := rmqrpc.Config{
		URL:      url,
		WaitTime: _defaultWaitTime,
		Attempts: _defaultAttempts,
	}

	s := &Server{
		conn:    rmqrpc.New(serverExchange, cfg),
		error:   make(chan error),
		stop:    make(chan struct{}),
		router:  router,
		timeout: _defaultTimeout,
		logger:  l,
	}

	// Custom options
	for _, opt := range opts {
		opt(s)
	}

	err := s.conn.AttemptConnect()
	if err != nil {
		return nil, fmt.Errorf("rmq_rpc server - NewServer - s.conn.AttemptConnect: %w", err)
	}

	go s.consumer()

	return s, nil
}

func (s *Server) consumer() {
	for {
		select {
		case <-s.stop:
			return
		case d, opened := <-s.conn.Delivery:
			if !opened {
				s.reconnect()

				return
			}

			_ = d.Ack(false) //nolint:errcheck // don't need this

			s.serveCall(&d)
		}
	}
}

func (s *Server) serveCall(d *amqp.Delivery) {
	callHandler, ok := s.router[d.Type]
	if !ok {
		s.publish(d, nil, rmqrpc.ErrBadHandler.Error())

		return
	}

	response, err := callHandler(d)
	if err != nil {
		s.publish(d, nil, rmqrpc.ErrInternalServer.Error())

		s.logger.Error(err, "rmq_rpc server - Server - serveCall - callHandler")

		return
	}

	body, err := json.Marshal(response)
	if err != nil {
		s.logger.Error(err, "rmq_rpc server - Server - serveCall - json.Marshal")
	}

	s.publish(d, body, rmqrpc.Success)
}

func (s *Server) publish(d *amqp.Delivery, body []byte, status string) {
	err := s.conn.Channel.Publish(d.ReplyTo, "", false, false,
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: d.CorrelationId,
			Type:          status,
			Body:          body,
		})
	if err != nil {
		s.logger.Error(err, "rmq_rpc server - Server - publish - s.conn.Channel.Publish")
	}
}

func (s *Server) reconnect() {
	close(s.stop)

	err := s.conn.AttemptConnect()
	if err != nil {
		s.error <- err
		close(s.error)

		return
	}

	s.stop = make(chan struct{})

	go s.consumer()
}

// Notify -.
func (s *Server) Notify() <-chan error {
	return s.error
}

// Shutdown -.
func (s *Server) Shutdown() error {
	select {
	case <-s.error:
		return nil
	default:
	}

	close(s.stop)
	time.Sleep(s.timeout)

	err := s.conn.Connection.Close()
	if err != nil {
		return fmt.Errorf("rmq_rpc server - Server - Shutdown - s.Connection.Close: %w", err)
	}

	return nil
}
