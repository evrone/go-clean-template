package main

import (
	"fmt"
	"github.com/evrone/go-clean-template/internal"
	"github.com/evrone/go-clean-template/pkg/httpserver"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log := internal.InitializeLogger()
	cfg := internal.InitializeConfig()

	//servers
	rmqServer := internal.InitializeNewRmqRpcServer()
	httpServer, _ := httpserver.New(cfg)

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	var err error
	select {
	case s := <-interrupt:
		log.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		log.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	case err = <-rmqServer.Notify():
		log.Error(fmt.Errorf("app - Run - rmqServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		log.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}

	err = rmqServer.Shutdown()
	if err != nil {
		log.Error(fmt.Errorf("app - Run - rmqServer.Shutdown: %w", err))
	}

}
