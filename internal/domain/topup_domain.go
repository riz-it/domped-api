package domain

import (
	"context"
	"time"

	"gorm.io/gorm"
	"riz.it/domped/internal/dto"
)

type TopUpEntity struct {
	ID        string    `gorm:"column:id;primaryKey;type:uuid"`
	UserID    int64     `gorm:"column:user_id;"`
	Amount    int64     `gorm:"column:amount"`
	Status    int8      `gorm:"column:status"`
	SnapURL   string    `gorm:"column:snap_url"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`

	// Relation
	User UserEntity `gorm:"foreignKey:UserID;reference:ID"`
}

func (TopUpEntity) TableName() string {
	return "public.topup"
}

// Interface
type TopUpRepository interface {
	Create(db *gorm.DB, topup *TopUpEntity) error
	FindAll(db *gorm.DB, topups *[]TopUpEntity) error
	FindByID(db *gorm.DB, topup *TopUpEntity, id int64) error
	FindByUUID(db *gorm.DB, topup *TopUpEntity, id string) error
	Update(db *gorm.DB, topup *TopUpEntity) error
	Delete(db *gorm.DB, topup *TopUpEntity) error

	// Custom functions
}

type TopUpUseCase interface {
	InitializeTopUp(ctx context.Context, req *dto.TopUpRequest, userID int64) (*dto.TopUpResponse, error)
	TopUpConfirmed(ctx context.Context, id string) error
}
