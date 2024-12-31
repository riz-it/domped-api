//go:build wireinject
// +build wireinject

package injector

import (
	"github.com/google/wire"
	"riz.it/domped/internal/config"
	"riz.it/domped/internal/delivery"
	"riz.it/domped/internal/delivery/controller"
	"riz.it/domped/internal/delivery/middleware"
	"riz.it/domped/internal/domain"
	"riz.it/domped/internal/repository"
	"riz.it/domped/internal/usecase"
	"riz.it/domped/internal/util"
)

var authSet = wire.NewSet(
	usecase.NewAuthUseCase,
	controller.NewAuthController,
)

var userSet = wire.NewSet(
	repository.NewUser,
	wire.Bind(new(domain.UserRepository), new(*repository.UserRepository)),
)

var walletSet = wire.NewSet(
	repository.NewWallet,
	wire.Bind(new(domain.WalletRepository), new(*repository.WalletRepository)),
)

var transactionSet = wire.NewSet(
	repository.NewTransaction,
	wire.Bind(new(domain.TransactionRepository), new(*repository.TransactionRepository)),
	usecase.NewTransactionUseCase,
	controller.NewTransactionController,
)

var middlewareSet = wire.NewSet(
	middleware.NewAuthMiddleware,
)

func InitializedApp() *config.App {
	wire.Build(
		config.Get,
		config.NewLogger,
		config.NewDatabase,
		config.NewValidator,
		config.NewFiber,
		config.NewApp,
		config.NewRedisClient,
		delivery.NewRouter,
		util.NewJWTUtil,
		util.NewEmailUtil,
		authSet,
		userSet,
		walletSet,
		transactionSet,
		middlewareSet,
	)
	return nil
}
