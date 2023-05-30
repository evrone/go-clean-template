// The build tag makes sure the stub is not built in the final build.
//go:build wireinject
// +build wireinject

package internal

import (
	"github.com/evrone/go-clean-template/config"
	amqprpc "github.com/evrone/go-clean-template/internal/controller/amqp_rpc"
	"github.com/evrone/go-clean-template/internal/usecase"
	"github.com/evrone/go-clean-template/internal/usecase/repository"
	"github.com/evrone/go-clean-template/internal/usecase/webapi"
	"github.com/evrone/go-clean-template/pkg/logger"
	"github.com/evrone/go-clean-template/pkg/postgres"
	"github.com/evrone/go-clean-template/pkg/rabbitmq/rmq_rpc/server"
	"github.com/google/wire"
)

var deps = []interface{}{}

var providerSet wire.ProviderSet = wire.NewSet(
	postgres.NewOrGetSingleton,
	repository.New,
	webapi.New,
	usecase.New,
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

func InitializeTranslationWebAPI() *webapi.TranslationWebAPI {
	wire.Build(providerSet)
	return &webapi.TranslationWebAPI{}
}

func InitializeTranslationUseCase() *usecase.TranslationUseCase {
	wire.Build(providerSet, config.NewConfig)
	return &usecase.TranslationUseCase{}
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
