package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"riz.it/domped/app/domain"
	"riz.it/domped/app/dto"
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
	// Parse the login request from the request body
	request := new(dto.LoginRequest)
	if err := ctx.BodyParser(request); err != nil {
		// Return a bad request error if parsing fails
		return fiber.ErrBadRequest
	}

	// Call the Login use case to authenticate the user
	response, err := c.AuthUseCase.Login(ctx.UserContext(), request)
	if err != nil {
		// Return the error from the use case
		return err
	}

	// Return the login response as a JSON object
	return ctx.JSON(&dto.ApiResponse[*dto.LoginResponse]{
		Status:  true,
		Message: "Login successful",
		Data:    &response,
	})
}

func (c *AuthController) Register(ctx *fiber.Ctx) error {
	// Parse the registration request from the request body
	request := new(dto.RegisterRequest)
	if err := ctx.BodyParser(request); err != nil {
		// Return a bad request error if parsing fails
		return fiber.ErrBadRequest
	}

	// Call the Register use case to register the user
	response, err := c.AuthUseCase.Register(ctx.UserContext(), request)
	if err != nil {
		// Return the error from the use case
		return err
	}

	// Return the registration response as a JSON object
	return ctx.JSON(&dto.ApiResponse[*dto.RegisterResponse]{
		Status:  true,
		Message: "Registration successful",
		Data:    &response,
	})
}

func (c *AuthController) Logout(ctx *fiber.Ctx) error {
	// Extract user ID from the context
	userID := ctx.Locals("userId").(int64)

	// Call the Logout use case to log the user out
	err := c.AuthUseCase.Logout(ctx.UserContext(), userID)
	if err != nil {
		// Return the error from the use case
		return err
	}

	// Return the logout response as a JSON object
	return ctx.JSON(&dto.ApiResponse[*dto.LoginResponse]{
		Status:  true,
		Message: "Logout successful",
	})
}

func (c *AuthController) Refresh(ctx *fiber.Ctx) error {
	// Parse the refresh token request from the request body
	request := new(dto.RefreshTokenRequest)
	if err := ctx.BodyParser(request); err != nil {
		// Return a bad request error if parsing fails
		return fiber.ErrBadRequest
	}

	// Call the Refresh use case to refresh the tokens
	response, err := c.AuthUseCase.Refresh(ctx.UserContext(), request)
	if err != nil {
		// Return the error from the use case
		return err
	}

	// Return the refreshed token response as a JSON object
	return ctx.JSON(&dto.ApiResponse[*dto.LoginResponse]{
		Status:  true,
		Message: "Token refreshed successfully",
		Data:    &response,
	})
}

func (c *AuthController) EmailVerification(ctx *fiber.Ctx) error {
	// Parse the email verification request from the request body
	request := new(dto.EmailVerificationRequest)
	if err := ctx.BodyParser(request); err != nil {
		// Return a bad request error if parsing fails
		return fiber.ErrBadRequest
	}

	// Call the EmailVerification use case to verify the email
	response, err := c.AuthUseCase.EmailVerification(ctx.UserContext(), request)
	if err != nil {
		// Return the error from the use case
		return err
	}

	// Return the email verification response as a JSON object
	return ctx.JSON(&dto.ApiResponse[*dto.LoginResponse]{
		Status:  true,
		Message: "Email verified successfully",
		Data:    &response,
	})
}
