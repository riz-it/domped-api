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
	"riz.it/domped/app/domain"
	"riz.it/domped/app/dto"
	"riz.it/domped/app/util"
)

type TopUpUseCase struct {
	DB                     *gorm.DB
	Log                    *logrus.Logger
	NotificationRepository domain.NotificationRepository
	MidtransUtil           domain.Midtrans
	TopUpRepository        domain.TopUpRepository
	WalletRepository       domain.WalletRepository
	TransactionRepository  domain.TransactionRepository
	Validate               *validator.Validate
}

func NewTopUpUseCase(db *gorm.DB, log *logrus.Logger, notificationRepository domain.NotificationRepository, midtransUtil domain.Midtrans, topUpRepository domain.TopUpRepository, walletRepository domain.WalletRepository, transactionRepository domain.TransactionRepository, validate *validator.Validate) domain.TopUpUseCase {
	return &TopUpUseCase{
		Log:                    log,
		DB:                     db,
		NotificationRepository: notificationRepository,
		MidtransUtil:           midtransUtil,
		TopUpRepository:        topUpRepository,
		WalletRepository:       walletRepository,
		TransactionRepository:  transactionRepository,
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

	tx := t.DB.WithContext(ctx)

	if err := t.TopUpRepository.Create(tx, topup); err != nil {
		t.Log.WithError(err).Error("Failed to create top-up")
		return nil, domain.NewError(fiber.StatusInternalServerError)
	}

	if err := t.MidtransUtil.GenerateSnapURL(ctx, topup); err != nil {
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

	t.Log.Info("Starting TopUpConfirmed process")
	t.Log.WithField("topup_id", id).Info("Fetching top-up details")

	tx := t.DB.WithContext(ctx)

	// Find top-up by UUID
	topup := new(domain.TopUpEntity)
	err := t.TopUpRepository.FindByUUID(tx, topup, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			t.Log.WithField("topup_id", id).Warn("Top-up not found")
			return domain.NewError(fiber.StatusNotFound, "Top-up not found")
		}
		t.Log.WithError(err).Error("Failed to fetch top-up details")
		return domain.NewError(fiber.StatusInternalServerError)
	}

	t.Log.WithField("user_id", topup.UserID).Info("Fetching wallet details")
	// Find wallet by user ID
	wallet := new(domain.WalletEntity)
	err = t.WalletRepository.FindByUserID(tx, wallet, topup.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			t.Log.WithField("user_id", topup.UserID).Warn("Wallet not found")
			return domain.NewError(fiber.StatusNotFound, "Wallet not found")
		}
		t.Log.WithError(err).Error("Failed to fetch wallet details")
		return domain.NewError(fiber.StatusInternalServerError)
	}

	// Log current values for debugging
	t.Log.WithFields(logrus.Fields{
		"topup_id":        topup.ID,
		"wallet_id":       wallet.ID,
		"current_balance": wallet.Balance,
		"topup_amount":    topup.Amount,
	}).Info("Retrieved top-up and wallet data")

	// Update top-up status
	t.Log.WithField("topup_id", id).Info("Updating top-up status")
	topup.Status = 1
	err = t.TopUpRepository.Update(tx, topup)
	if err != nil {
		t.Log.WithError(err).Error("Failed to update top-up status")
		return domain.NewError(fiber.StatusInternalServerError)
	}

	// Update wallet balance
	t.Log.WithField("user_id", wallet.UserID).Info("Updating wallet balance")
	wallet.Balance += topup.Amount
	err = t.WalletRepository.Update(tx, wallet)
	if err != nil {
		t.Log.WithError(err).Error("Failed to update wallet balance")
		return domain.NewError(fiber.StatusInternalServerError)
	}

	// Log updated wallet balance
	t.Log.WithFields(logrus.Fields{
		"wallet_id":   wallet.ID,
		"new_balance": wallet.Balance,
	}).Info("Wallet balance updated successfully")

	// Create transaction record
	t.Log.Info("Creating transaction record")
	transaction := &domain.TransactionEntity{
		WalletID:        wallet.ID,
		SofNumber:       "00",
		DofNumber:       wallet.WalletNumber,
		Amount:          topup.Amount,
		TransactionType: "D",
	}
	if err = t.TransactionRepository.Create(tx, transaction); err != nil {
		t.Log.WithError(err).Error("Failed to create transaction")
		return domain.NewError(fiber.StatusInternalServerError)
	}

	// Send notification
	t.Log.Info("Sending notification after top-up")
	t.notificationAfterTopUp(c, *wallet, topup.Amount)

	// Log before commit
	t.Log.WithFields(logrus.Fields{
		"topup_id":  topup.ID,
		"wallet_id": wallet.ID,
		"amount":    topup.Amount,
	}).Info("Ready to commit transaction")

	// Commit transaction
	t.Log.Info("Committing transaction")
	if err := tx.Commit().Error; err != nil {
		t.Log.WithError(err).Error("Failed to commit top-up transaction")
		return domain.NewError(fiber.StatusInternalServerError)
	}

	t.Log.Info("TopUpConfirmed process completed successfully")
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
		t.Log.WithError(err).Error("Failed to commit notification transaction")
		tx.Rollback()
		return
	}

}
