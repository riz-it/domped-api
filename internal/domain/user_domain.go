package domain

import (
	"time"

	"gorm.io/gorm"
)

// Entity
type UserEntity struct {
	ID              int64      `gorm:"column:id;primaryKey"`
	FullName        string     `gorm:"column:full_name"`
	Email           string     `gorm:"column:email"`
	Phone           string     `gorm:"column:phone"`
	Password        string     `gorm:"column:password"`
	HashedRt        string     `gorm:"column:hashed_rt"`
	IsActive        bool       `gorm:"column:is_active"`
	EmailVerifiedAt *time.Time `gorm:"column:email_verified_at"`
	CreatedAt       time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt       time.Time  `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`

	// Relation
	Wallet WalletEntity `gorm:"foreignKey:UserID;reference:ID"`
}

func (UserEntity) TableName() string {
	return "public.users"
}

// Interface
type UserRepository interface {
	Create(db *gorm.DB, user *UserEntity) error
	FindAll(db *gorm.DB, users *[]UserEntity) error
	FindByID(db *gorm.DB, user *UserEntity, id int64) error
	Update(db *gorm.DB, user *UserEntity) error
	Delete(db *gorm.DB, user *UserEntity) error

	// Custom functions
	FindByEmail(db *gorm.DB, user *UserEntity, email string) error
	CountByEmail(db *gorm.DB, email string) (count int64, err error)
}
