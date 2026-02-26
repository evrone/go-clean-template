// Package app configures and runs the ToDo microservice.
package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"evrone.local/common-pkg/httpserver"
	"evrone.local/common-pkg/logger"
	"evrone.local/common-pkg/postgres"
	"github.com/evrone/todo-svc/config"
	"github.com/evrone/todo-svc/internal/controller/restapi"
	"github.com/evrone/todo-svc/internal/repo/persistent"
	"github.com/evrone/todo-svc/internal/usecase/todo"
)

// Run creates objects via constructors and runs the ToDo microservice.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	// Repository
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	// Use-Case
	todoUseCase := todo.New(persistent.NewTodoRepo(pg))

	// HTTP Server
	httpServer := httpserver.New(l, httpserver.Port(cfg.HTTP.Port), httpserver.Prefork(cfg.HTTP.UsePreforkMode))
	restapi.NewRouter(httpServer.App, cfg, todoUseCase, l)

	// Start server
	httpServer.Start()

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: %s", s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
