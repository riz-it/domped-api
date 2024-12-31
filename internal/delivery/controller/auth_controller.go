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

func (c *AuthController) Register(ctx *fiber.Ctx) error {

	request := new(dto.RegisterRequest)
	if err := ctx.BodyParser(request); err != nil {
		return fiber.ErrBadRequest
	}

	response, err := c.AuthUseCase.Register(ctx.UserContext(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(&dto.ApiResponse[*dto.LoginResponse]{
		Status:  true,
		Message: "Register successfully",
		Data:    &response,
	})
}

func (c *AuthController) Logout(ctx *fiber.Ctx) error {

	userID := ctx.Locals("userId").(int64)
	err := c.AuthUseCase.Logout(ctx.UserContext(), userID)
	if err != nil {
		return err
	}

	return ctx.JSON(&dto.ApiResponse[*dto.LoginResponse]{
		Status:  true,
		Message: "Logout successfully",
	})
}

func (c *AuthController) Refresh(ctx *fiber.Ctx) error {

	request := new(dto.RefreshTokenRequest)
	if err := ctx.BodyParser(request); err != nil {
		return fiber.ErrBadRequest
	}

	response, err := c.AuthUseCase.Refresh(ctx.UserContext(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(&dto.ApiResponse[*dto.LoginResponse]{
		Status:  true,
		Message: "Token refreshed successfully",
		Data:    &response,
	})
}
