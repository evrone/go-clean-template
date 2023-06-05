// Package httpserver implements HTTP Server.
package httpserver

import (
	"context"
	"github.com/evrone/go-clean-template/config"
	openapi "github.com/evrone/go-clean-template/internal/interfaces/rest/v1/go"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net"
	"net/http"
	"time"
)

const (
	_defaultReadTimeout     = 5 * time.Second
	_defaultWriteTimeout    = 5 * time.Second
	_defaultAddr            = ":80"
	_defaultShutdownTimeout = 3 * time.Second
)

// Server -.
type Server struct {
	server          *http.Server
	notify          chan error
	shutdownTimeout time.Duration
}

func New(cfg *config.Config) (*Server, *gin.Engine) {

	router := openapi.NewRouter()
	setupMonitoringRoutes(router)

	server := prepareHttpServer(cfg, router)
	server.start()

	return server, router
}

func prepareHttpServer(cfg *config.Config, router *gin.Engine) *Server {
	httpServer := &http.Server{
		Handler:      router,
		ReadTimeout:  _defaultReadTimeout,
		WriteTimeout: _defaultWriteTimeout,
		Addr:         _defaultAddr,
	}
	httpServer.Addr = net.JoinHostPort("", cfg.HTTP.Port)

	s := &Server{
		server:          httpServer,
		notify:          make(chan error, 1),
		shutdownTimeout: _defaultShutdownTimeout,
	}
	return s
}

func setupMonitoringRoutes(handler *gin.Engine) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// K8s probe
	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Prometheus metrics
	handler.GET("/metrics", gin.WrapH(promhttp.Handler()))
}

func (s *Server) start() {
	go func() {
		s.notify <- s.server.ListenAndServe()
		close(s.notify)
	}()
}

// Notify -.
func (s *Server) Notify() <-chan error {
	return s.notify
}

// Shutdown -.
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}
