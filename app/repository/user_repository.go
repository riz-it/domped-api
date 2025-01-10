package repository

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"riz.it/domped/app/domain"
)

type UserRepository struct {
	Repository[domain.UserEntity]
	Log *logrus.Logger
}

func NewUser(log *logrus.Logger) *UserRepository {
	return &UserRepository{
		Log: log,
	}
}

func (u *UserRepository) FindByEmail(db *gorm.DB, user *domain.UserEntity, email string) error {
	return db.Model(&domain.UserEntity{}).Where("email = ?", email).First(&user).Error
}

func (u *UserRepository) CountByEmail(db *gorm.DB, email string) (count int64, err error) {
	err = db.Model(&domain.UserEntity{}).Where("email = ?", email).Count(&count).Error
	return count, err
}
