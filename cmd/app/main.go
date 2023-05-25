package main

import (
	"github.com/evrone/go-clean-template/config"
	"github.com/evrone/go-clean-template/internal/app"
)

func main() {
	// Configuration
	cfg := config.NewConfig()

	// Run
	app.Run(cfg)
}
