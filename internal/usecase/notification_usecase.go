package usecase

import (
	"context"

	"github.com/go-playground/validator/v10"
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

func NewNotificationRepository(db *gorm.DB, log *logrus.Logger, notificationRepository domain.NotificationRepository, validate *validator.Validate) domain.NotificationUseCase {
	return &NotificationUseCase{
		DB:                     db,
		Log:                    log,
		NotificationRepository: notificationRepository,
		Validate:               validate,
	}
}

// FindByUserID implements domain.NotificationUseCase.
func (n *NotificationUseCase) FindByUserID(ctx context.Context, userID int64) (*dto.NotificationData, error) {
	// Set a timeout for the logout process
	// c, cancel := context.WithTimeout(ctx, 10*time.Second)
	// defer cancel()

	// // Start a new transaction to ensure atomicity
	// tx := n.DB.WithContext(c)

	// // Retrieve the user based on userID
	// notifications := new([]domain.NotificationEntity)
	// if err := n.NotificationRepository.FindByUserID(tx, notifications, userID); err != nil {
	// 	// If user is not found, return a 'Not Found' error
	// 	return nil, domain.NewError(fiber.StatusNotFound, "User not found")
	// }
	// // Successful logout, return nil indicating no errors
	return nil, nil
}
