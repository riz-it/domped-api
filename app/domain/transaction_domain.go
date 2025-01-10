package domain

import (
	"context"
	"time"

	"gorm.io/gorm"
	"riz.it/domped/app/dto"
)

// Entity
type TransactionEntity struct {
	ID              int64     `gorm:"column:id;primaryKey"`
	WalletID        int64     `gorm:"column:wallet_id"`
	SofNumber       string    `gorm:"column:sof_number"`
	DofNumber       string    `gorm:"column:dof_number"`
	Amount          int64     `gorm:"column:amount"`
	TransactionType string    `gorm:"column:transaction_type"`
	TransactionAt   time.Time `gorm:"column:transaction_at;autoCreateTime"`

	// Relation
	Wallet WalletEntity `gorm:"foreignKey:WalletID;reference:ID"`
}

func (TransactionEntity) TableName() string {
	return "public.transactions"
}

// Interface
type TransactionRepository interface {
	Create(db *gorm.DB, transaction *TransactionEntity) error
	FindAll(db *gorm.DB, transactions *[]TransactionEntity) error
	FindByID(db *gorm.DB, transaction *TransactionEntity, id int64) error
	Update(db *gorm.DB, transaction *TransactionEntity) error
	Delete(db *gorm.DB, transaction *TransactionEntity) error

	// Custom functions
}

type TransactionUseCase interface {
	TransferInquiry(ctx context.Context, req *dto.TransferInquiryRequest, userID int64) (*dto.TransferInquiryResponse, error)
	TransferExecute(ctx context.Context, req *dto.TransferExecuteRequest, userID int64) (*dto.TransferExecuteResponse, error)
}
