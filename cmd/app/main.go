package main

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"

	"github.com/evrone/go-service-template/config"
	"github.com/evrone/go-service-template/internal/app"
	"github.com/evrone/go-service-template/pkg/logger"
)

func main() {
	// Configuration
	var cfg config.Config

	err := cleanenv.ReadConfig("./config/config.yml", &cfg)
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Logger
	zap := logger.NewZapLogger(cfg.Log.ZapLevel)
	defer zap.Close()

	rollbar := logger.NewRollbarLogger(cfg.Log.RollbarToken, cfg.Log.RollbarEnv)
	defer rollbar.Close()

	logger.NewAppLogger(zap, rollbar, cfg.App.Name, cfg.App.Version)

	// Run
	app.Run(&cfg)
}
