package usecase

import (
	"context"
	"errors"
	"fmt"
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
	DB                     *gorm.DB
	Log                    *logrus.Logger
	NotificationRepository domain.NotificationRepository
	MidtransUtil           domain.Midtrans
	TopUpRepository        domain.TopUpRepository
	WalletRepository       domain.WalletRepository
	Validate               *validator.Validate
}

func NewTopUpUseCase(db *gorm.DB, log *logrus.Logger, notificationRepository domain.NotificationRepository, midtransUtil domain.Midtrans, topUpRepository domain.TopUpRepository, walletRepository domain.WalletRepository, validate *validator.Validate) domain.TopUpUseCase {
	return &TopUpUseCase{
		Log:                    log,
		NotificationRepository: notificationRepository,
		MidtransUtil:           midtransUtil,
		TopUpRepository:        topUpRepository,
		WalletRepository:       walletRepository,
		Validate:               validate,
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
		ID:     util.GenerateUUID(),
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
func (t *TopUpUseCase) TopUpConfirmed(c context.Context, id string) error {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	tx := t.DB.WithContext(ctx)

	topup := new(domain.TopUpEntity)
	err := t.TopUpRepository.FindByUUID(tx, topup, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.NewError(fiber.StatusNotFound, "Top-up not found")
		}

		return domain.NewError(fiber.StatusInternalServerError)
	}

	wallet := new(domain.WalletEntity)
	err = t.WalletRepository.FindByID(tx, wallet, topup.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.NewError(fiber.StatusNotFound, "Wallet not found")
		}

		return domain.NewError(fiber.StatusInternalServerError)
	}

	wallet.Balance += topup.Amount
	err = t.WalletRepository.Update(tx, wallet)
	if err != nil {
		return domain.NewError(fiber.StatusInternalServerError)
	}

	t.notificationAfterTopUp(c, *wallet, topup.Amount)

	if err := tx.Commit().Error; err != nil {
		t.Log.WithError(err).Warnf("Failed to commit transaction: %+v", err)
		return domain.NewError(fiber.StatusInternalServerError)
	}

	return nil
}

func (t *TopUpUseCase) notificationAfterTopUp(c context.Context, wallet domain.WalletEntity, amount int64) {
	tx := t.DB.WithContext(c).Begin()

	formattedAmount := util.CurrencyFormat(float64(amount))

	notification := domain.NotificationEntity{
		UserID: wallet.UserID,
		Title:  "TopUp Berhasil",
		Body:   fmt.Sprintf("TopUp senilai %s berhasil dilakukan.", formattedAmount),
		IsRead: false,
		Status: 1,
	}

	if err := t.NotificationRepository.Create(tx, &notification); err != nil {
		t.Log.WithError(err).Error("Failed to create sender notification")
		tx.Rollback()
		return
	}

	if err := tx.Commit().Error; err != nil {
		t.Log.WithError(err).Error("Failed to commit transaction")
		tx.Rollback()
		return
	}

}
