package repository

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"riz.it/domped/app/domain"
)

type TopUpRepository struct {
	Repository[domain.TopUpEntity]
	Log *logrus.Logger
}

func NewTopUp(log *logrus.Logger) *TopUpRepository {
	return &TopUpRepository{
		Log: log,
	}
}

func (u *TopUpRepository) FindByUUID(db *gorm.DB, topup *domain.TopUpEntity, orderID string) error {
	return db.Model(&domain.TopUpEntity{}).Where("id = ?", orderID).First(&topup).Error
}
