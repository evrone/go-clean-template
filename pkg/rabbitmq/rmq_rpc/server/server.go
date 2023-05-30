package server

import (
	"encoding/json"
	"fmt"
	"github.com/evrone/go-clean-template/config"
	"time"

	"github.com/streadway/amqp"

	"github.com/evrone/go-clean-template/pkg/logger"
	rmqrpc "github.com/evrone/go-clean-template/pkg/rabbitmq/rmq_rpc"
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

	logger *logger.Logger
}

// New -.
func New(config *config.Config, log *logger.Logger, amqpRpcRouter map[string]CallHandler) *Server {

	cfg := rmqrpc.Config{
		URL:      config.RMQ.URL,
		WaitTime: _defaultWaitTime,
		Attempts: _defaultAttempts,
	}

	s := &Server{
		conn:    rmqrpc.New(config.RMQ.ServerExchange, cfg),
		error:   make(chan error),
		stop:    make(chan struct{}),
		router:  amqpRpcRouter,
		timeout: _defaultTimeout,
		logger:  log,
	}

	err := s.conn.AttemptConnect()
	if err != nil {
		panic(fmt.Errorf("rmq_rpc server - NewServer - s.conn.AttemptConnect: %w", err))
	}

	go s.consumer()

	return s
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

		s.logger.Error(err, "rmq_rpc server - server - serveCall - callHandler")

		return
	}

	body, err := json.Marshal(response)
	if err != nil {
		s.logger.Error(err, "rmq_rpc server - server - serveCall - json.Marshal")
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
		s.logger.Error(err, "rmq_rpc server - server - publish - s.conn.Channel.Publish")
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
		return fmt.Errorf("rmq_rpc server - server - Shutdown - s.Connection.Close: %w", err)
	}

	return nil
}
