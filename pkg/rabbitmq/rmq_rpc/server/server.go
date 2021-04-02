package server

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"

	"github.com/evrone/go-service-template/pkg/logger"
	rmqrpc "github.com/evrone/go-service-template/pkg/rabbitmq/rmq_rpc"
)

const (
	defaultWaitTime = 5 * time.Second
	defaultAttempts = 10
	defaultTimeout  = 2 * time.Second
)

type CallHandler func(*amqp.Delivery) (interface{}, error)

type Server struct {
	conn   *rmqrpc.Connection
	error  chan error
	stop   chan struct{}
	router map[string]CallHandler

	timeout time.Duration
}

func NewServer(url, serverExchange string, router map[string]CallHandler, opts ...Option) (*Server, error) {
	cfg := rmqrpc.Config{
		URL:      url,
		WaitTime: defaultWaitTime,
		Attempts: defaultAttempts,
	}

	s := &Server{
		conn:    rmqrpc.NewConnection(serverExchange, cfg),
		error:   make(chan error),
		stop:    make(chan struct{}),
		router:  router,
		timeout: defaultTimeout,
	}

	// Custom options
	for _, opt := range opts {
		opt(s)
	}

	err := s.conn.AttemptConnect()
	if err != nil {
		return nil, errors.Wrap(err, "rmq_rpc server - NewServer - s.conn.AttemptConnect")
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

		logger.Error(err, "rmq_rpc server - Server - serveCall - callHandler")

		return
	}

	body, err := json.Marshal(response)
	if err != nil {
		logger.Error(err, "rmq_rpc server - Server - serveCall - json.Marshal")
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
		logger.Error(err, "rmq_rpc server - Server - publish - s.conn.Channel.Publish")
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

func (s *Server) Notify() <-chan error {
	return s.error
}

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
