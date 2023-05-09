package main

import (
	"github.com/evrone/go-clean-template/internal/app"
	"log"

	"github.com/evrone/go-clean-template/config"
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
