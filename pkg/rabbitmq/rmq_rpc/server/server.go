// Package server implements RabbitMQ RPC server.
package server

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/evrone/go-clean-template/pkg/logger"
	rmqrpc "github.com/evrone/go-clean-template/pkg/rabbitmq/rmq_rpc"
	"github.com/goccy/go-json"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const (
	_defaultWaitTime = 5 * time.Second
	_defaultAttempts = 10
	_defaultTimeout  = 2 * time.Second
)

const _tracerName = "rmq-rpc.server"

// CallHandler -.
type CallHandler func(context.Context, *amqp.Delivery) (any, error)

// Server -.
type Server struct {
	conn   *rmqrpc.Connection
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

	return s, nil
}

// Start starts the RabbitMQ RPC server and blocks until context is canceled.
func (s *Server) Start(ctx context.Context) error {
	s.logger.Info("rmq_rpc server - Server - Started")

	// Wait for shutdown signal
	go func() {
		<-ctx.Done()
		s.logger.Info("rmq_rpc server - Server - Shutting down...")
	}()

	// Handle messages until context is canceled
	err := s.handleMessages(ctx)

	// Close connection
	if closeErr := s.conn.Connection.Close(); closeErr != nil {
		s.logger.Error(closeErr, "rmq_rpc server - Server - Shutdown - s.Connection.Close")

		if err == nil {
			err = closeErr
		}
	}

	s.logger.Info("rmq_rpc server - Server - Shutdown complete")

	if errors.Is(err, context.Canceled) {
		return ctx.Err()
	}

	return err
}

func (s *Server) handleMessages(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case d, opened := <-s.conn.Delivery:
			if !opened {
				err := s.reconnect()
				if err != nil {
					return err
				}

				break
			}

			s.serveCall(ctx, &d)
		}
	}
}

func (s *Server) reconnect() error {
	return s.conn.AttemptConnect()
}

func (s *Server) serveCall(ctx context.Context, d *amqp.Delivery) {
	defer s.ack(d, false)

	ctx = otel.GetTextMapPropagator().Extract(ctx, rmqrpc.TableCarrier(d.Headers))

	ctx, span := otel.Tracer(_tracerName).Start(
		ctx, "rmq_rpc.process "+d.Type,
		trace.WithSpanKind(trace.SpanKindConsumer),
		trace.WithAttributes(attribute.String("rpc.method", d.Type)),
	)
	defer span.End()

	callHandler, ok := s.router[d.Type]
	if !ok {
		span.SetStatus(codes.Error, rmqrpc.ErrBadHandler.Error())
		s.publish(d, nil, rmqrpc.ErrBadHandler.Error())

		return
	}

	response, err := callHandler(ctx, d)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
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

func (s *Server) ack(d *amqp.Delivery, multiple bool) {
	err := d.Ack(multiple)
	if err != nil {
		s.logger.Error(err, "rmq_rpc server - Server - ack - d.Ack")
	}
}

func (s *Server) publish(d *amqp.Delivery, body []byte, status string) {
	err := s.conn.Channel.Publish(
		d.ReplyTo,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: d.CorrelationId,
			Type:          status,
			Body:          body,
		},
	)
	if err != nil {
		s.logger.Error(err, "rmq_rpc server - Server - publish - s.conn.Channel.Publish")
	}
}
