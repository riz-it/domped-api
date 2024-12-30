package domain

import (
	"database/sql"

	"gorm.io/gorm"
)

// Entity
type UserEntity struct {
	ID              int64        `gorm:"column:id;primaryKey"`
	FullName        string       `gorm:"column:full_name"`
	Email           string       `gorm:"column:email"`
	Password        string       `gorm:"column:password"`
	HashedRt        string       `gorm:"column:hashed_rt"`
	IsActive        bool         `gorm:"column:is_active"`
	EmailVerifiedAt sql.NullTime `gorm:"column:email_verified_at"`
	CreatedAt       sql.NullTime `gorm:"created_at"`
	UpdatedAt       sql.NullTime `gorm:"updated_at"`
}

func (UserEntity) TableName() string {
	return "public.users"
}

// Interface
type UserRepository interface {
	FindByID(db *gorm.DB, user *UserEntity, id int64) error
	FindByEmail(db *gorm.DB, user *UserEntity, email string) error
	Update(db *gorm.DB, user *UserEntity) error
}
