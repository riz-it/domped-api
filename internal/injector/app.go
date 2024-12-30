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
	repository.NewUser,
	wire.Bind(new(domain.UserRepository), new(*repository.UserRepository)),
	usecase.NewAuthUseCase,
	controller.NewAuthController,
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
		delivery.NewRouter,
		util.NewJWTUtil,
		authSet,
		middlewareSet,
	)
	return nil
}
