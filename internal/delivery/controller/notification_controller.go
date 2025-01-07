package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"riz.it/domped/internal/domain"
	"riz.it/domped/internal/dto"
)

type NotificationController struct {
	NotificationUseCase domain.NotificationUseCase
	Log                 *logrus.Logger
}

func NewNotificationController(notificationUseCase domain.NotificationUseCase, log *logrus.Logger) *NotificationController {
	return &NotificationController{
		NotificationUseCase: notificationUseCase,
		Log:                 log,
	}
}

func (n *NotificationController) GetUserNotifications(ctx *fiber.Ctx) error {

	// Extract user ID from the context
	userID := ctx.Locals("userId").(int64)

	// Call the notification use case to refresh the tokens
	result, err := n.NotificationUseCase.FindByUserID(ctx.UserContext(), userID)
	if err != nil {
		// Return the error from the use case
		return err
	}

	// Return the refreshed token response as a JSON object
	return ctx.JSON(&dto.ApiResponse[[]dto.NotificationData]{
		Status:  true,
		Message: "User notifications retrieved successfully",
		Data:    result,
	})
}
