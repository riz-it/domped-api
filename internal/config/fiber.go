package config

import (
	"github.com/gofiber/fiber/v2"
	"riz.it/domped/internal/dto"
)

func NewFiber(config *Config) *fiber.App {
	var app = fiber.New(fiber.Config{
		AppName:      config.Server.Name,
		ErrorHandler: NewErrorHandler(),
	})

	return app
}

func NewErrorHandler() fiber.ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {
		if e, ok := err.(*dto.ApiResponse[any]); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(e)
		}

		return ctx.Status(err.(*fiber.Error).Code).JSON(&dto.ApiResponse[string]{
			Status:  false,
			Message: err.Error(),
		})
	}
}
