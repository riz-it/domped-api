package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"riz.it/domped/app/domain"
	"riz.it/domped/app/dto"
	"riz.it/domped/app/util"
)

type TransactionUseCase struct {
	DB                     *gorm.DB
	Log                    *logrus.Logger
	WalletRepository       domain.WalletRepository
	TransactionRepository  domain.TransactionRepository
	NotificationRepository domain.NotificationRepository
	Validate               *validator.Validate
	Redis                  *redis.Client
}

func NewTransactionUseCase(db *gorm.DB, log *logrus.Logger, walletRepository domain.WalletRepository, transactionRepository domain.TransactionRepository, notificationRepository domain.NotificationRepository, validate *validator.Validate, redis *redis.Client) domain.TransactionUseCase {
	return &TransactionUseCase{
		DB:                     db,
		Log:                    log,
		WalletRepository:       walletRepository,
		TransactionRepository:  transactionRepository,
		NotificationRepository: notificationRepository,
		Validate:               validate,
		Redis:                  redis,
	}
}

// TransferInquiry implements domain.TransactionUseCase.
func (t *TransactionUseCase) TransferInquiry(ctx context.Context, req *dto.TransferInquiryRequest, userID int64) (*dto.TransferInquiryResponse, error) {
	// Set a timeout for the process
	c, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Validate the incoming request data
	if validationErrors := util.Validate(t.Validate, req); len(validationErrors) > 0 {
		return nil, domain.NewError(fiber.StatusBadRequest, "Invalid data provided", validationErrors)
	}

	// Retrieve source wallet based on userID
	wallet := new(domain.WalletEntity)
	if err := t.WalletRepository.FindByUserID(t.DB.WithContext(c), wallet, userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.NewError(fiber.StatusNotFound, "Source wallet not found")
		}
		t.Log.WithError(err).Warn("Failed to query source wallet")
		return nil, domain.NewError(fiber.StatusInternalServerError)
	}

	// Check if destination wallet is same as source wallet
	if req.AccountNumber == wallet.WalletNumber {
		return nil, domain.NewError(fiber.StatusBadRequest, "Invalid destination wallet")
	}

	// Check if pin wallet not setup yet
	if wallet.WalletPin == "" {
		return nil, domain.NewError(fiber.StatusBadRequest, "Please set your pin wallet")
	}

	// Retrieve destination wallet based on account number
	dofWallet := new(domain.WalletEntity)
	if err := t.WalletRepository.FindByWalletNumber(t.DB.WithContext(c), dofWallet, req.AccountNumber); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.NewError(fiber.StatusNotFound, "Destination wallet not found")
		}
		t.Log.WithError(err).Warn("Failed to query destination wallet")
		return nil, domain.NewError(fiber.StatusInternalServerError)
	}

	// Check if balance is sufficient
	if wallet.Balance < req.Amount {
		return nil, domain.NewError(fiber.StatusBadRequest, "Balance is insufficient")
	}

	// Generate inquiry key and serialize request
	inquiryKey := util.GenerateRandomString(32)
	ttl := time.Hour * 24

	inquiryData, err := json.Marshal(req)
	if err != nil {
		t.Log.WithError(err).Warn("Failed to serialize inquiry data")
		return nil, domain.NewError(fiber.StatusInternalServerError, "Failed to process inquiry")
	}

	// Store the inquiry in Redis
	insertInquiry := t.Redis.Set(ctx, inquiryKey, inquiryData, ttl)
	if err := insertInquiry.Err(); err != nil {
		t.Log.WithError(err).Warn("Failed to store inquiry in Redis")
		return nil, domain.NewError(fiber.StatusInternalServerError)
	}

	// Return the inquiry key
	return &dto.TransferInquiryResponse{
		InquiryKey: inquiryKey,
	}, nil
}

// TransferExecute implements domain.TransactionUseCase.
func (t *TransactionUseCase) TransferExecute(ctx context.Context, req *dto.TransferExecuteRequest, userID int64) (*dto.TransferExecuteResponse, error) {
	// Set a timeout for the process
	c, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Validate the incoming request data
	if validationErrors := util.Validate(t.Validate, req); len(validationErrors) > 0 {
		return nil, domain.NewError(fiber.StatusBadRequest, "Invalid data provided", validationErrors)
	}

	// Retrieve Inquiry from Redis
	getInquiry := t.Redis.Get(c, req.InquiryKey)
	if err := getInquiry.Err(); err != nil {
		// Return an error if the ReferenceID is invalid or the OTP retrieval fails
		return nil, domain.NewError(fiber.StatusBadRequest, "Invalid inquiry key")
	}

	data, err := getInquiry.Result()
	if err != nil {
		t.Log.WithError(err).Error("Failed to get inquiry response")
		return nil, domain.NewError(fiber.StatusInternalServerError)
	}

	inquiryData := new(dto.TransferInquiryRequest)
	if err := json.Unmarshal([]byte(data), &inquiryData); err != nil {
		t.Log.WithError(err).Error("Failed to deserialize inquiry data")
		return nil, domain.NewError(fiber.StatusInternalServerError)
	}

	// Start a transaction to ensure atomicity of user update
	tx := t.DB.WithContext(c).Begin()
	defer tx.Rollback()

	// Retrieve source wallet based on userID
	wallet := new(domain.WalletEntity)
	if err := t.WalletRepository.FindByUserID(tx, wallet, userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.NewError(fiber.StatusNotFound, "Source wallet not found")
		}
		t.Log.WithError(err).Warn("Failed to query source wallet")
		return nil, domain.NewError(fiber.StatusInternalServerError)
	}

	// Retrieve destination wallet based on account number
	dofWallet := new(domain.WalletEntity)
	if err := t.WalletRepository.FindByWalletNumber(tx, dofWallet, inquiryData.AccountNumber); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.NewError(fiber.StatusNotFound, "Destination wallet not found")
		}
		t.Log.WithError(err).Warn("Failed to query destination wallet")
		return nil, domain.NewError(fiber.StatusInternalServerError)
	}

	// Check if pin code is valid
	if !util.VerifyPassword(wallet.WalletPin, req.PinCode) {
		// Return an error if the password is invalid
		return nil, domain.NewError(fiber.StatusBadRequest, "Invalid pin code")
	}

	now := time.Now()

	// Transaction
	debitTransaction := domain.TransactionEntity{
		WalletID:        dofWallet.ID,
		SofNumber:       wallet.WalletNumber,
		DofNumber:       dofWallet.WalletNumber,
		TransactionType: "D",
		Amount:          inquiryData.Amount,
		TransactionAt:   now,
	}

	if err := t.TransactionRepository.Create(tx, &debitTransaction); err != nil {
		t.Log.WithError(err).Error("Failed to create debit transaction")
		return nil, domain.NewError(fiber.StatusInternalServerError)
	}

	creditTransaction := domain.TransactionEntity{
		WalletID:        wallet.ID,
		SofNumber:       wallet.WalletNumber,
		DofNumber:       dofWallet.WalletNumber,
		TransactionType: "C",
		Amount:          inquiryData.Amount,
		TransactionAt:   now,
	}

	if err := t.TransactionRepository.Create(tx, &creditTransaction); err != nil {
		t.Log.WithError(err).Error("Failed to create credit transaction")
		return nil, domain.NewError(fiber.StatusInternalServerError)
	}

	wallet.Balance -= inquiryData.Amount
	if err := t.WalletRepository.Update(tx, wallet); err != nil {
		t.Log.WithError(err).Error("Failed to update wallet balance")
		return nil, domain.NewError(fiber.StatusInternalServerError)
	}

	dofWallet.Balance += inquiryData.Amount
	if err := t.WalletRepository.Update(tx, dofWallet); err != nil {
		t.Log.WithError(err).Error("Failed to update wallet balance")
		return nil, domain.NewError(fiber.StatusInternalServerError)
	}

	// Delete the OTP from Redis as it is no longer needed
	delInquiry := t.Redis.Del(c, req.InquiryKey)
	if err := delInquiry.Err(); err != nil {
		t.Log.WithError(err).Warnf("Failed to delete OTP: %+v", err)
		return nil, domain.NewError(fiber.StatusInternalServerError)
	}

	// Commit the transaction to persist changes
	if err := tx.Commit().Error; err != nil {
		t.Log.WithError(err).Warnf("Failed to commit topup transaction: %+v", err)
		return nil, domain.NewError(fiber.StatusInternalServerError)
	}

	t.notificationAfterTransfer(c, *wallet, *dofWallet, inquiryData.Amount)

	return &dto.TransferExecuteResponse{
		InquiryKey: req.InquiryKey,
		Information: dto.TransferData{
			SofNumber:     wallet.WalletNumber,
			DofNumber:     dofWallet.WalletNumber,
			Amount:        inquiryData.Amount,
			TransactionAt: now.String(),
			Status:        "success",
		},
	}, nil

}
func (t *TransactionUseCase) notificationAfterTransfer(c context.Context, sofWallet domain.WalletEntity, dofWallet domain.WalletEntity, amount int64) {
	tx := t.DB.WithContext(c).Begin()

	formattedAmount := util.CurrencyFormat(float64(amount))

	notificationSender := domain.NotificationEntity{
		UserID: sofWallet.UserID,
		Title:  "Transfer Berhasil",
		Body:   fmt.Sprintf("Transfer senilai %s berhasil dilakukan.", formattedAmount),
		IsRead: false,
		Status: 1,
	}

	notificationReceiver := domain.NotificationEntity{
		UserID: dofWallet.UserID,
		Title:  "Dana Diterima",
		Body:   fmt.Sprintf("Dana senilai %s telah diterima.", formattedAmount),
		IsRead: false,
		Status: 1,
	}

	if err := t.NotificationRepository.Create(tx, &notificationSender); err != nil {
		t.Log.WithError(err).Error("Failed to create sender notification")
		tx.Rollback()
		return
	}

	if err := t.NotificationRepository.Create(tx, &notificationReceiver); err != nil {
		t.Log.WithError(err).Error("Failed to create receiver notification")
		tx.Rollback()
		return
	}

	if err := tx.Commit().Error; err != nil {
		t.Log.WithError(err).Error("Failed to commit notification transaction")
		tx.Rollback()
		return
	}

}
