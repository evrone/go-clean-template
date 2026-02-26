package main

import (
	"log"

	"github.com/evrone/todo-svc/config"
	"github.com/evrone/todo-svc/internal/app"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
