//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package internal

import (
	"github.com/evrone/go-clean-template/config"
	"github.com/evrone/go-clean-template/internal/application"
	"github.com/evrone/go-clean-template/internal/infrastructure/googleapi"
	"github.com/evrone/go-clean-template/internal/infrastructure/repository"
	amqprpc "github.com/evrone/go-clean-template/internal/interfaces/amqp_rpc"
	"github.com/evrone/go-clean-template/pkg/logger"
	"github.com/evrone/go-clean-template/pkg/postgres"
	"github.com/evrone/go-clean-template/pkg/rabbitmq/rmq_rpc/server"
	"github.com/google/wire"
)

var deps = []interface{}{}

var providerSet wire.ProviderSet = wire.NewSet(
	postgres.NewOrGetSingleton,
	repository.New,
	googleapi.New,
	application.New,
	logger.New,
	amqprpc.NewRouter,
	server.New,
)

func InitializeConfig() *config.Config {
	wire.Build(config.NewConfig)
	return &config.Config{}
}

func InitializePostgresConnection() *postgres.Postgres {
	wire.Build(providerSet, config.NewConfig)
	return &postgres.Postgres{}
}

func InitializeTranslationRepository() *repository.TranslationRepository {
	wire.Build(providerSet, config.NewConfig)
	return &repository.TranslationRepository{}
}

func InitializeTranslationWebAPI() *googleapi.GoogleTranslator {
	wire.Build(providerSet)
	return &googleapi.GoogleTranslator{}
}

func InitializeTranslationUseCase() *application.TranslationUseCase {
	wire.Build(providerSet, config.NewConfig)
	return &application.TranslationUseCase{}
}

func InitializeLogger() *logger.Logger {
	wire.Build(providerSet, config.NewConfig)
	return &logger.Logger{}
}

func InitializeNewRmqRpcServer() *server.Server {
	wire.Build(providerSet, config.NewConfig)
	return &server.Server{}
}

func InitializeNewRmqRpcServerWithConfig(config *config.Config) *server.Server {
	wire.Build(providerSet)
	return &server.Server{}
}
