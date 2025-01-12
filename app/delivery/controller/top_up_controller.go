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
	if err := ctx.BodyParser(&payload); err != nil {
		// Log error when parsing the body fails
		fmt.Printf("Error parsing request body: %v\n", err)
		return fiber.ErrBadRequest
	}

	// Extract order_id from the payload
	orderId, exists := payload["order_id"].(string)
	if !exists {
		// Log missing order_id field in the request
		fmt.Println("Error: order_id not found in payload")
		return fiber.ErrBadRequest
	}

	// Log the order_id received for debugging
	fmt.Printf("Received order_id: %s\n", orderId)

	// Call the VerifyPayment use case to verify the payment status
	success, err := t.MidtransUtil.VerifyPayment(ctx.Context(), orderId)
	if err != nil {
		// Log the error from VerifyPayment
		fmt.Printf("Error verifying payment for order_id %s: %v\n", orderId, err)
		return err
	}

	if success {
		// Log success of payment verification
		fmt.Printf("Payment verified successfully for order_id: %s\n", orderId)

		// Call the TopUpConfirmed use case to confirm the top-up transaction
		if err = t.TopUpUseCase.TopUpConfirmed(ctx.Context(), orderId); err != nil {
			// Log error if TopUpConfirmed fails
			fmt.Printf("Error confirming top-up for order_id %s: %v\n", orderId, err)
			return err
		}

		// Log success of the top-up confirmation
		fmt.Printf("Top-up confirmed successfully for order_id: %s\n", orderId)
	} else {
		// Log if payment verification failed
		fmt.Printf("Payment verification failed for order_id: %s\n", orderId)
	}

	// Return the top-up success response as a JSON object
	return ctx.JSON(&dto.ApiResponse[string]{
		Status:  true,
		Message: "TopUp successfully",
	})
}
