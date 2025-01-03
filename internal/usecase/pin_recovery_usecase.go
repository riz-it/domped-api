package usecase

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"riz.it/domped/internal/domain"
	"riz.it/domped/internal/dto"
)

type PinRecoveryUsecase struct {
	DB               *gorm.DB
	Log              *logrus.Logger
	WalletRepository domain.WalletRepository
	Validate         *validator.Validate
}

func NewPinRecoveryUsecase(db *gorm.DB, log *logrus.Logger, walletRepository domain.WalletRepository, validate *validator.Validate) domain.PinRecoveryUseCase {
	return &PinRecoveryUsecase{
		DB:               db,
		Log:              log,
		WalletRepository: walletRepository,
		Validate:         validate,
	}
}

// SetupWalletPIN implements domain.PinRecoveryUseCase.
func (p *PinRecoveryUsecase) SetupWalletPIN(ctx context.Context, req *dto.SetupWalletPINRequest) error {
	panic("unimplemented")
}
