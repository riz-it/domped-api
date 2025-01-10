package domain

import (
	"context"
	"time"

	"gorm.io/gorm"
	"riz.it/domped/app/dto"
)

type NotificationEntity struct {
	ID        int64     `gorm:"column:id;primaryKey"`
	UserID    int64     `gorm:"column:user_id"`
	Status    int8      `gorm:"column:status"`
	Title     string    `gorm:"column:title"`
	Body      string    `gorm:"column:body"`
	IsRead    bool      `gorm:"column:is_read"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`

	// Relation
	User UserEntity `gorm:"foreignKey:UserID;reference:ID"`
}

func (NotificationEntity) TableName() string {
	return "public.notifications"
}

// Interface
type NotificationRepository interface {
	Create(db *gorm.DB, n *NotificationEntity) error
	FindAll(db *gorm.DB, n *[]NotificationEntity) error
	FindByID(db *gorm.DB, n *NotificationEntity, id int64) error
	Update(db *gorm.DB, n *NotificationEntity) error
	Delete(db *gorm.DB, n *NotificationEntity) error

	// Custom functions
	FindByUserID(db *gorm.DB, notifications *[]NotificationEntity, userID int64) error
}

type NotificationUseCase interface {
	FindByUserID(ctx context.Context, userID int64) (*[]dto.NotificationData, error)
}
