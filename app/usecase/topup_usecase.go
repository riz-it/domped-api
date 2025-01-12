package usecase

import (
	"context"
	"errors"
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
	err = t.WalletRepository.FindByUserID(tx, wallet, topup.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.NewError(fiber.StatusNotFound, "Wallet not found")
		}

		return domain.NewError(fiber.StatusInternalServerError)
	}

	topup.Status = 1
	err = t.TopUpRepository.Update(tx, topup)
	if err != nil {
		return domain.NewError(fiber.StatusInternalServerError)
	}

	wallet.Balance += topup.Amount
	err = t.WalletRepository.Update(tx, wallet)
	if err != nil {
		return domain.NewError(fiber.StatusInternalServerError)
	}

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

	// t.notificationAfterTopUp(tx, *wallet, topup.Amount)

	if err := tx.Commit().Error; err != nil {
		t.Log.WithError(err).Warnf("Failed to commit topup transaction: %+v", err)
		return domain.NewError(fiber.StatusInternalServerError)
	}

	return nil
}

// func (t *TopUpUseCase) notificationAfterTopUp(tx *gorm.DB, wallet domain.WalletEntity, amount int64) {

// 	formattedAmount := util.CurrencyFormat(float64(amount))

// 	notification := domain.NotificationEntity{
// 		UserID: wallet.UserID,
// 		Title:  "TopUp Berhasil",
// 		Body:   fmt.Sprintf("TopUp senilai %s berhasil dilakukan.", formattedAmount),
// 		IsRead: false,
// 		Status: 1,
// 	}

// 	if err := t.NotificationRepository.Create(tx, &notification); err != nil {
// 		t.Log.WithError(err).Error("Failed to create sender notification")
// 		tx.Rollback()
// 		return
// 	}

// 	if err := tx.Commit().Error; err != nil {
// 		t.Log.WithError(err).Error("Failed to commit notification transaction")
// 		tx.Rollback()
// 		return
// 	}

// }
