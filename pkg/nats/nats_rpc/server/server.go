// Package server implements NATS RPC server.
package server

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/evrone/go-clean-template/pkg/logger"
	natsrpc "github.com/evrone/go-clean-template/pkg/nats/nats_rpc"
	"github.com/goccy/go-json"
	"github.com/nats-io/nats.go"
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

const _tracerName = "nats-rpc.server"

// CallHandler -.
type CallHandler func(context.Context, *nats.Msg) (any, error)

// Server -.
type Server struct {
	subject      string
	connection   *nats.Conn
	subscription *nats.Subscription
	router       map[string]CallHandler

	timeout time.Duration

	logger logger.Interface
}

// New -.
func New(
	url,
	serverSubject string,
	router map[string]CallHandler,
	l logger.Interface,
	opts ...Option,
) (*Server, error) {
	connection, err := nats.Connect(
		url,
		nats.ReconnectWait(_defaultWaitTime),
		nats.MaxReconnects(_defaultAttempts),
		nats.Timeout(_defaultWaitTime),
	)
	if err != nil {
		return nil, fmt.Errorf("nats_rpc server - NewServer - nats.Connect: %w", err)
	}

	s := &Server{
		subject:    serverSubject,
		connection: connection,
		router:     router,
		timeout:    _defaultTimeout,
		logger:     l,
	}

	// Custom options
	for _, opt := range opts {
		opt(s)
	}

	return s, nil
}

// Start starts the NATS RPC server and blocks until context is canceled.
func (s *Server) Start(ctx context.Context) error {
	err := s.subscribe(ctx)
	if err != nil {
		return err
	}

	s.logger.Info("nats_rpc server - Server - Started")

	// Wait for shutdown signal
	<-ctx.Done()

	s.logger.Info("nats_rpc server - Server - Shutting down...")

	var shutdownErrors []error

	// Unsubscribe
	if s.subscription != nil {
		if err := s.subscription.Unsubscribe(); err != nil {
			s.logger.Error(err, "nats_rpc server - Server - Shutdown - s.subscription.Unsubscribe")
			shutdownErrors = append(shutdownErrors, err)
		}
	}

	// Close connection
	s.connection.Close()

	s.logger.Info("nats_rpc server - Server - Shutdown complete")

	if len(shutdownErrors) > 0 {
		return errors.Join(shutdownErrors...)
	}

	return ctx.Err()
}

func (s *Server) subscribe(ctx context.Context) error {
	subscription, err := s.connection.Subscribe(s.subject, func(msg *nats.Msg) {
		s.handleMessage(ctx, msg)
	})
	if err != nil {
		return fmt.Errorf("nats_rpc server - subscribe - s.connection.Subscribe: %w", err)
	}

	s.subscription = subscription

	return nil
}

func (s *Server) handleMessage(ctx context.Context, msg *nats.Msg) {
	handler := msg.Header.Get("Handler")

	ctx = otel.GetTextMapPropagator().Extract(ctx, natsrpc.HeaderCarrier(msg.Header))

	ctx, span := otel.Tracer(_tracerName).Start(
		ctx, "nats_rpc.process "+handler,
		trace.WithSpanKind(trace.SpanKindConsumer),
		trace.WithAttributes(attribute.String("rpc.method", handler)),
	)
	defer span.End()

	callHandler, ok := s.router[handler]
	if !ok {
		span.SetStatus(codes.Error, natsrpc.ErrBadHandler.Error())
		s.publish(msg, nil, natsrpc.ErrBadHandler.Error())

		return
	}

	response, err := callHandler(ctx, msg)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		s.publish(msg, nil, natsrpc.ErrInternalServer.Error())

		s.logger.Error(err, "nats_rpc server - Server - handleMessage - callHandler")

		return
	}

	body, err := json.Marshal(response)
	if err != nil {
		s.logger.Error(err, "nats_rpc server - Server - handleMessage - json.Marshal")

		s.publish(msg, nil, natsrpc.ErrInternalServer.Error())

		return
	}

	s.publish(msg, body, natsrpc.Success)
}

func (s *Server) publish(msg *nats.Msg, body []byte, status string) {
	respondMsg := nats.NewMsg(msg.Reply)
	respondMsg.Header.Set("Status", status)
	respondMsg.Data = body

	err := s.connection.PublishMsg(respondMsg)
	if err != nil {
		s.logger.Error(err, "nats_rpc server - Server - publish - msg.Respond")
	}
}
