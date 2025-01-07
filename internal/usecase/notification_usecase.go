package usecase

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"riz.it/domped/internal/domain"
	"riz.it/domped/internal/dto"
)

type NotificationUseCase struct {
	DB                     *gorm.DB
	Log                    *logrus.Logger
	NotificationRepository domain.NotificationRepository
	Validate               *validator.Validate
}

func NewNotificationUseCase(db *gorm.DB, log *logrus.Logger, notificationRepository domain.NotificationRepository, validate *validator.Validate) domain.NotificationUseCase {
	return &NotificationUseCase{
		DB:                     db,
		Log:                    log,
		NotificationRepository: notificationRepository,
		Validate:               validate,
	}
}

// FindByUserID implements domain.NotificationUseCase.
func (n *NotificationUseCase) FindByUserID(ctx context.Context, userID int64) (*[]dto.NotificationData, error) {
	// Set a timeout for the process
	c, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Start a new transaction to ensure atomicity
	tx := n.DB.WithContext(c)

	// Retrieve the notifications for the given userID
	notifications := new([]domain.NotificationEntity)
	if err := n.NotificationRepository.FindByUserID(tx, notifications, userID); err != nil {
		// If an error occurs while fetching notifications, return the error
		return nil, domain.NewError(fiber.StatusNotFound, "Notifications not found for the given user")
	}

	// If no notifications found, return an empty array instead of nil
	if len(*notifications) == 0 {
		return &[]dto.NotificationData{}, nil
	}

	var result []dto.NotificationData
	for _, v := range *notifications {
		result = append(result, dto.NotificationData{
			ID:        v.ID,
			Title:     v.Title,
			Body:      v.Body,
			IsRead:    v.IsRead,
			Status:    v.Status,
			CreatedAt: v.CreatedAt.String(),
			UpdatedAt: v.UpdatedAt.String(),
		})
	}

	if result == nil {
		result = make([]dto.NotificationData, 0)
	}

	return &result, nil
}
