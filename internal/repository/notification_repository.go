package repository

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"riz.it/domped/internal/domain"
)

type NotificationRepository struct {
	Repository[domain.NotificationEntity]
	Log *logrus.Logger
}

func NewNotification(log *logrus.Logger) *NotificationRepository {
	return &NotificationRepository{
		Log: log,
	}
}

func (u *NotificationRepository) FindByUserID(db *gorm.DB, notifications *[]domain.NotificationEntity, userID int64) error {
	return db.Where("user_id = ?", userID).Find(notifications).Error
}
