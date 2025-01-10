package repository

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"riz.it/domped/app/domain"
)

type WalletRepository struct {
	Repository[domain.WalletEntity]
	Log *logrus.Logger
}

func NewWallet(log *logrus.Logger) *WalletRepository {
	return &WalletRepository{
		Log: log,
	}
}

func (u *WalletRepository) FindByUserID(db *gorm.DB, wallet *domain.WalletEntity, userID int64) error {
	return db.Model(&domain.WalletEntity{}).Where("user_id = ?", userID).First(&wallet).Error
}

func (u *WalletRepository) FindByWalletNumber(db *gorm.DB, wallet *domain.WalletEntity, walletNumber string) error {
	return db.Model(&domain.WalletEntity{}).Where("wallet_number = ?", walletNumber).First(&wallet).Error
}

func (u *WalletRepository) CountByWalletNumber(db *gorm.DB, walletNumber string) (count int64, err error) {
	err = db.Model(&domain.WalletEntity{}).Where("wallet_number = ?", walletNumber).Count(&count).Error
	return count, err
}
