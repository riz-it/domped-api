package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"riz.it/domped/internal/domain"
	"riz.it/domped/internal/dto"
)

type AuthController struct {
	AuthUseCase domain.AuthUseCase
	Log         *logrus.Logger
}

func NewAuthController(authUseCase domain.AuthUseCase, log *logrus.Logger) *AuthController {
	return &AuthController{
		AuthUseCase: authUseCase,
		Log:         log,
	}
}

func (c *AuthController) Login(ctx *fiber.Ctx) error {

	request := new(dto.LoginRequest)
	if err := ctx.BodyParser(request); err != nil {
		return fiber.ErrBadRequest
	}

	response, err := c.AuthUseCase.Login(ctx.UserContext(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(&dto.ApiResponse[*dto.LoginResponse]{
		Status:  true,
		Message: "Login successful",
		Data:    &response,
	})
}
