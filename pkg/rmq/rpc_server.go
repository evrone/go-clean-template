package rmq

import (
	"fmt"
	"time"

	"github.com/streadway/amqp"

	"github.com/evrone/go-service-template/pkg/logger"
)

type CallHandler func(*amqp.Delivery) ([]byte, error)

type Server struct {
	conn    *Connection
	timeout time.Duration
	error   chan error
	stop    chan struct{}
	router  map[string]CallHandler
}

func NewServer(router map[string]CallHandler, serverExchange string) *Server {
	s := &Server{
		conn:    newConnection(serverExchange),
		timeout: 2 * time.Second, //nolint:gomnd // will be config
		error:   make(chan error),
		stop:    make(chan struct{}),
		router:  router,
	}

	go s.consumer()

	return s
}

func (s *Server) consumer() {
	for {
		select {
		case <-s.stop:
			return
		case d, opened := <-s.conn.delivery:
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
		s.publish(d, nil, ErrBadHandler.Error())

		return
	}

	body, err := callHandler(d)
	if err != nil {
		s.publish(d, nil, ErrInternalServer.Error())

		logger.Error(err, "rmq - Server - serveCall - callHandler")

		return
	}

	s.publish(d, body, Success)
}

func (s *Server) publish(d *amqp.Delivery, body []byte, status string) {
	err := s.conn.channel.Publish(d.ReplyTo, "", false, false,
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: d.CorrelationId,
			Type:          status,
			Body:          body,
		})
	if err != nil {
		logger.Error(err, "rmq - Server - publish - s.conn.channel.Publish")
	}
}

func (s *Server) reconnect() {
	close(s.stop)

	err := s.conn.attemptConnect()
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

	err := s.conn.connection.Close()
	if err != nil {
		return fmt.Errorf("rmq - Server - Shutdown - s.connection.Close: %w", err)
	}

	return nil
}
