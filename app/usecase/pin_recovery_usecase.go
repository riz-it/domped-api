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

type PinRecoveryUseCase struct {
	DB                    *gorm.DB
	Log                   *logrus.Logger
	WalletRepository      domain.WalletRepository
	PinRecoveryRepository domain.PinRecoveryRepository
	Validate              *validator.Validate
}

func NewPinRecoveryUseCase(db *gorm.DB, log *logrus.Logger, walletRepository domain.WalletRepository, pinRecoveryRepository domain.PinRecoveryRepository, validate *validator.Validate) domain.PinRecoveryUseCase {
	return &PinRecoveryUseCase{
		DB:                    db,
		Log:                   log,
		WalletRepository:      walletRepository,
		PinRecoveryRepository: pinRecoveryRepository,
		Validate:              validate,
	}
}

// SetupWalletPIN implements domain.PinRecoveryUseCase.
func (p *PinRecoveryUseCase) SetupWalletPIN(ctx context.Context, req *dto.SetupWalletPINRequest) error {
	// Set a timeout for the registration process
	c, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Validate the incoming request data
	if validationErrors := util.Validate(p.Validate, req); len(validationErrors) > 0 {
		return domain.NewError(fiber.StatusBadRequest, "Invalid data provided", validationErrors)
	}

	if req.PinCode != req.PinCodeConfirmation {
		return domain.NewError(fiber.StatusBadRequest, "Pin code does not match")
	}

	// Start a new transaction to ensure atomicity
	tx := p.DB.WithContext(c).Begin()
	defer tx.Rollback()

	// Check if the wallet is exist
	wallet := new(domain.WalletEntity)
	err := p.WalletRepository.FindByID(tx, wallet, req.WalletID)
	if err != nil {
		// Cek apakah error adalah "not found"
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.NewError(fiber.StatusNotFound, "Wallet does not exist")
		}
		p.Log.WithError(err).Warnf("Failed to find wallet: %+v", err)
		return domain.NewError(fiber.StatusInternalServerError)
	}
	if wallet == (&domain.WalletEntity{}) {
		return domain.NewError(fiber.StatusNotFound, "Wallet not found")
	}

	// Hash the password for storage
	hashedPin, err := util.HashPassword(req.PinCode)
	if err != nil {
		p.Log.WithError(err).Warnf("Failed to hash pin code: %+v", err)
		return domain.NewError(fiber.StatusInternalServerError)
	}

	// Create a new user entity
	pinRecovery := &domain.PinRecoveryEntity{
		WalletID: wallet.ID,
		PinCode:  hashedPin,
	}

	// Save the new user to the database
	if err := p.PinRecoveryRepository.Create(tx, pinRecovery); err != nil {
		p.Log.WithError(err).Warnf("Failed to create pin recovery: %+v", err)
		return domain.NewError(fiber.StatusInternalServerError)
	}

	// Update pin code
	wallet.WalletPin = hashedPin
	if err := p.WalletRepository.Update(tx, wallet); err != nil {
		p.Log.WithError(err).Warnf("Failed to update wallet pin: %+v", err)
		return domain.NewError(fiber.StatusInternalServerError)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		p.Log.WithError(err).Warnf("Failed to commit transaction: %+v", err)
		return domain.NewError(fiber.StatusInternalServerError)
	}

	// Return nil error
	return nil
}
