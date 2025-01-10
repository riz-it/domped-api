package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"riz.it/domped/app/domain"
	"riz.it/domped/app/dto"
)

type TransactionController struct {
	TransactionUseCase domain.TransactionUseCase
	Log                *logrus.Logger
}

func NewTransactionController(transactionUseCase domain.TransactionUseCase, log *logrus.Logger) *TransactionController {
	return &TransactionController{
		TransactionUseCase: transactionUseCase,
		Log:                log,
	}
}

func (t *TransactionController) Inquiry(ctx *fiber.Ctx) error {
	// Extract user ID from the context
	userID := ctx.Locals("userId").(int64)

	// Parse the refresh token request from the request body
	request := new(dto.TransferInquiryRequest)
	if err := ctx.BodyParser(request); err != nil {
		// Return a bad request error if parsing fails
		return fiber.ErrBadRequest
	}

	// Call the Refresh use case to refresh the tokens
	response, err := t.TransactionUseCase.TransferInquiry(ctx.UserContext(), request, userID)
	if err != nil {
		// Return the error from the use case
		return err
	}

	// Return the refreshed token response as a JSON object
	return ctx.JSON(&dto.ApiResponse[*dto.TransferInquiryResponse]{
		Status:  true,
		Message: "Transfer in progress",
		Data:    &response,
	})
}

func (t *TransactionController) Execute(ctx *fiber.Ctx) error {
	// Extract user ID from the context
	userID := ctx.Locals("userId").(int64)

	// Parse the refresh token request from the request body
	request := new(dto.TransferExecuteRequest)
	if err := ctx.BodyParser(request); err != nil {
		// Return a bad request error if parsing fails
		return fiber.ErrBadRequest
	}

	// Call the Refresh use case to refresh the tokens
	response, err := t.TransactionUseCase.TransferExecute(ctx.UserContext(), request, userID)
	if err != nil {
		// Return the error from the use case
		return err
	}

	// Return the refreshed token response as a JSON object
	return ctx.JSON(&dto.ApiResponse[*dto.TransferExecuteResponse]{
		Status:  true,
		Message: "Transfer in successfully completed",
		Data:    &response,
	})
}
