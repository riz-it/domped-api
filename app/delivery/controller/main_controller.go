package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"riz.it/domped/app/dto"
)

type MainController struct {
	Log *logrus.Logger
}

func NewMainController(log *logrus.Logger) *MainController {
	return &MainController{
		Log: log,
	}
}

func (c *MainController) Main(ctx *fiber.Ctx) error {

	// Return the login response as a JSON object
	return ctx.JSON(&dto.ApiResponse[string]{
		Status:  true,
		Message: "API Domped Digital",
	})
}
