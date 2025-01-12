package domain

import (
	"context"
	"time"

	"gorm.io/gorm"
	"riz.it/domped/app/dto"
)

type PinRecoveryEntity struct {
	ID        int64     `gorm:"column:id;primaryKey"`
	PinCode   string    `gorm:"column:pin_code"`
	WalletID  int64     `gorm:"column:wallet_id"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`

	// Relation
	Wallet WalletEntity `gorm:"foreignKey:WalletID;reference:ID"`
}

func (PinRecoveryEntity) TableName() string {
	return "public.pin_recoveries"
}

// Interface
type PinRecoveryRepository interface {
	Create(db *gorm.DB, pr *PinRecoveryEntity) error
	// FindAll(db *gorm.DB, pr *[]PinRecoveryEntity) error
	FindByID(db *gorm.DB, pr *PinRecoveryEntity, id int64) error
	// Update(db *gorm.DB, pr *PinRecoveryEntity) error
	// Delete(db *gorm.DB, pr *PinRecoveryEntity) error
}

type PinRecoveryUseCase interface {
	SetupWalletPIN(ctx context.Context, req *dto.SetupWalletPINRequest, userID int64) error
}
