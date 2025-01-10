package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"riz.it/domped/app/domain"
	"riz.it/domped/app/dto"
)

type PinRecoveryController struct {
	PinRecoveryUseCase domain.PinRecoveryUseCase
	Log                *logrus.Logger
}

func NewPinRecoveryController(pinRecoveryUseCase domain.PinRecoveryUseCase, log *logrus.Logger) *PinRecoveryController {
	return &PinRecoveryController{
		PinRecoveryUseCase: pinRecoveryUseCase,
		Log:                log,
	}
}

func (p *PinRecoveryController) SetupPin(ctx *fiber.Ctx) error {

	// Parse the refresh token request from the request body
	request := new(dto.SetupWalletPINRequest)
	if err := ctx.BodyParser(request); err != nil {
		// Return a bad request error if parsing fails
		return fiber.ErrBadRequest
	}

	// Call the Refresh use case to refresh the tokens
	err := p.PinRecoveryUseCase.SetupWalletPIN(ctx.UserContext(), request)
	if err != nil {
		// Return the error from the use case
		return err
	}

	// Return the refreshed token response as a JSON object
	return ctx.JSON(&dto.ApiResponse[string]{
		Status:  true,
		Message: "Pin code successfully updated",
	})
}
