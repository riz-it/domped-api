package domain

import (
	"time"

	"gorm.io/gorm"
)

// Entity
type WalletEntity struct {
	ID           int64     `gorm:"column:id;primaryKey"`
	UserID       int64     `gorm:"column:user_id"`
	WalletNumber string    `gorm:"column:wallet_number"`
	WalletPin    string    `gorm:"column:wallet_pin"`
	Balance      int64     `gorm:"column:balance"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`

	// Relation
	Transaction []TransactionEntity `gorm:"foreignKey:WalletID;reference:ID"`
	User        *UserEntity         `gorm:"foreignKey:UserID;reference:ID"`
}

func (WalletEntity) TableName() string {
	return "public.wallets"
}

// Interface
type WalletRepository interface {
	Create(db *gorm.DB, wallet *WalletEntity) error
	FindAll(db *gorm.DB, wallets *[]WalletEntity) error
	FindByID(db *gorm.DB, wallet *WalletEntity, id int64) error
	Update(db *gorm.DB, wallet *WalletEntity) error
	Delete(db *gorm.DB, wallet *WalletEntity) error

	// Custom functions
	FindByUserID(db *gorm.DB, user *WalletEntity, userID int64) error
	FindByWalletNumber(db *gorm.DB, wallet *WalletEntity, walletNumber string) error
}
