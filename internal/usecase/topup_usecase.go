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
	"riz.it/domped/internal/util"
)

type TopUpUseCase struct {
	DB                  *gorm.DB
	Log                 *logrus.Logger
	NotificationUseCase domain.NotificationUseCase
	MidtransUtil        domain.Midtrans
	TopUpRepository     domain.TopUpRepository
	Validate            *validator.Validate
}

func NewTopUpUseCase(db *gorm.DB, log *logrus.Logger, notificationUseCase domain.NotificationUseCase, midtransUtil domain.Midtrans, topUpRepository domain.TopUpRepository, validate *validator.Validate) domain.TopUpUseCase {
	return &TopUpUseCase{
		Log:                 log,
		NotificationUseCase: notificationUseCase,
		MidtransUtil:        midtransUtil,
		TopUpRepository:     topUpRepository,
		Validate:            validate,
	}
}

// InitializeTopUp implements domain.TopUpUseCase.
func (t *TopUpUseCase) InitializeTopUp(ctx context.Context, req *dto.TopUpRequest, userID int64) (*dto.TopUpResponse, error) {
	// Set a timeout for the registration process
	_, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Validate the incoming request data
	if validationErrors := util.Validate(t.Validate, req); len(validationErrors) > 0 {
		return nil, domain.NewError(fiber.StatusBadRequest, "Invalid data provided", validationErrors)
	}

	topup := &domain.TopUpEntity{
		UserID: userID,
		Amount: req.Amount,
		Status: 0,
	}

	err := t.MidtransUtil.GenerateSnapURL(ctx, topup)

	if err != nil {
		return nil, domain.NewError(fiber.StatusInternalServerError)
	}

	return &dto.TopUpResponse{
		SnapURL: topup.SnapURL,
	}, nil
}

// TopUpConfirmed implements domain.TopUpUseCase.
func (t *TopUpUseCase) TopUpConfirmed(ctx context.Context, id string) error {
	panic("unimplemented")
}
