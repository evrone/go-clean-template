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

var providerSet wire.ProviderSet = wire.NewSet(
	postgres.NewOrGetSingleton,
	config.NewConfig,
	repository.New,
	webapi.New,
	usecase.New,
	logger.New,
	amqprpc.NewRouter,
	server.New,
)

func InitializePostgresConnection() *postgres.Postgres {
	wire.Build(providerSet)
	return &postgres.Postgres{}
}

func InitializeTranslationRepository() *repository.TranslationRepository {
	wire.Build(providerSet)
	return &repository.TranslationRepository{}
}

func InitializeTranslationWebAPI() *webapi.TranslationWebAPI {
	wire.Build(providerSet)
	return &webapi.TranslationWebAPI{}
}

func InitializeTranslationUseCase() *usecase.TranslationUseCase {
	wire.Build(providerSet)
	return &usecase.TranslationUseCase{}
}

func InitializeLogger() *logger.Logger {
	wire.Build(providerSet)
	return &logger.Logger{}
}

func InitializeNewAmqpRpcRouter() map[string]server.CallHandler {
	wire.Build(providerSet)
	return map[string]server.CallHandler{}
}

func InitializeNewRmqRpcServer() *server.Server {
	wire.Build(providerSet)
	return &server.Server{}
}
