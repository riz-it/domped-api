package config

import (
	"github.com/gofiber/fiber/v2"
	"riz.it/domped/app/delivery"
)

type App struct {
	Fiber  *fiber.App
	Config *Config
}

func NewApp(
	fiber *delivery.RouterConfig,
	config *Config,
) *App {
	return &App{
		Fiber:  fiber.App,
		Config: config,
	}
}
