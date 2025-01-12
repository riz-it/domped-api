package controller

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"riz.it/domped/app/domain"
	"riz.it/domped/app/dto"
)

type TopUpController struct {
	TopUpUseCase domain.TopUpUseCase
	Log          *logrus.Logger
	MidtransUtil domain.Midtrans
}

func NewTopUpController(topUpUseCase domain.TopUpUseCase, log *logrus.Logger, midtransUtil domain.Midtrans) *TopUpController {
	return &TopUpController{
		TopUpUseCase: topUpUseCase,
		Log:          log,
		MidtransUtil: midtransUtil,
	}
}

func (t *TopUpController) Initialize(ctx *fiber.Ctx) error {
	// Extract user ID from the context
	userID := ctx.Locals("userId").(int64)

	// Parse the refresh token request from the request body
	request := new(dto.TopUpRequest)
	if err := ctx.BodyParser(request); err != nil {
		// Return a bad request error if parsing fails
		return fiber.ErrBadRequest
	}

	// Call the Refresh use case to refresh the tokens
	response, err := t.TopUpUseCase.InitializeTopUp(ctx.UserContext(), request, userID)
	if err != nil {
		// Return the error from the use case
		return err
	}

	// Return the refreshed token response as a JSON object
	return ctx.JSON(&dto.ApiResponse[*dto.TopUpResponse]{
		Status:  true,
		Message: "Waiting for payment",
		Data:    &response,
	})
}

func (t *TopUpController) Verify(ctx *fiber.Ctx) error {
	// Extract user ID from the context
	var payload map[string]interface{}
	// Parse the refresh token request from the request body
	fmt.Println(payload)
	if err := ctx.BodyParser(payload); err != nil {
		// Return a bad request error if parsing fails
		return fiber.ErrBadRequest
	}

	orderId, exists := payload["order_id"].(string)
	if !exists {
		return fiber.ErrBadRequest
	}

	// Call the Refresh use case to refresh the tokens
	_, err := t.MidtransUtil.VerifyPayment(ctx.Context(), orderId)
	if err != nil {
		// Return the error from the use case
		return err
	}

	if err = t.TopUpUseCase.TopUpConfirmed(ctx.Context(), orderId); err != nil {
		// Return the error from the use case
		return err
	}

	// Return the refreshed token response as a JSON object
	return ctx.JSON(&dto.ApiResponse[string]{
		Status:  true,
		Message: "TopUp successfully",
	})
}
