package usecase

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"riz.it/domped/internal/domain"
	"riz.it/domped/internal/dto"
	"riz.it/domped/internal/util"
)

type AuthUseCase struct {
	DB             *gorm.DB
	UserRepository domain.UserRepository
	Log            *logrus.Logger
	JWT            domain.JWT
	Validate       *validator.Validate
}

func NewAuthUseCase(db *gorm.DB, log *logrus.Logger, userRepository domain.UserRepository, jwt domain.JWT, validate *validator.Validate) domain.AuthUseCase {
	return &AuthUseCase{
		DB:             db,
		Log:            log,
		UserRepository: userRepository,
		JWT:            jwt,
		Validate:       validate,
	}
}

// Login implements domain.AuthUseCase.
func (a *AuthUseCase) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	c, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if validationErrors := util.Validate(a.Validate, req); len(validationErrors) > 0 {
		return nil, domain.NewError(fiber.StatusBadRequest, "The provided data is invalid", validationErrors)
	}

	tx := a.DB.WithContext(c).Begin()
	defer tx.Rollback()
	user := new(domain.UserEntity)
	if err := a.UserRepository.FindByEmail(tx, user, req.Email); err != nil {
		return nil, domain.NewError(fiber.StatusUnauthorized, "Invalid email or password")
	}
	if !util.VerifyPassword(user.Password, req.Password) {
		return nil, domain.NewError(fiber.StatusUnauthorized, "Invalid email or password")
	}

	accessToken, refreshToken, err := a.JWT.GenerateToken(user.ID)
	if err != nil {
		a.Log.WithError(err).Error("Failed to generate tokens")
		return nil, domain.NewError(fiber.StatusInternalServerError)
	}
	user.HashedRt = refreshToken
	if err := a.UserRepository.Update(tx, user); err != nil {
		a.Log.WithError(err).Warnf("Failed to save user data: %+v", err)
		return nil, domain.NewError(fiber.StatusInternalServerError)
	}
	if err := tx.Commit().Error; err != nil {
		return nil, domain.NewError(fiber.StatusInternalServerError)
	}

	return &dto.LoginResponse{
		User: dto.CredentialData{
			FullName: user.FullName,
			Email:    user.Email,
		},
		Token: dto.TokenData{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}

// Register implements domain.AuthUseCase.
func (a *AuthUseCase) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.LoginResponse, error) {
	c, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if validationErrors := util.Validate(a.Validate, req); len(validationErrors) > 0 {
		return nil, domain.NewError(fiber.StatusBadRequest, "The provided data is invalid", validationErrors)
	}

	tx := a.DB.WithContext(c).Begin()
	defer tx.Rollback()
	count, err := a.UserRepository.CountByEmail(tx, req.Email)
	if err != nil {
		return nil, domain.NewError(fiber.StatusInternalServerError)
	}
	if count > 0 {
		return nil, domain.NewError(fiber.StatusConflict, "Email already in use")
	}

	password, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, domain.NewError(fiber.StatusInternalServerError)
	}

	user := &domain.UserEntity{
		Password: string(password),
		FullName: req.FullName,
		Email:    req.Email,
		IsActive: true,
	}

	if err := a.UserRepository.Create(tx, user); err != nil {
		return nil, domain.NewError(fiber.StatusInternalServerError)
	}

	accessToken, refreshToken, err := a.JWT.GenerateToken(user.ID)
	if err != nil {
		a.Log.WithError(err).Error("Failed to generate tokens")
		return nil, domain.NewError(fiber.StatusInternalServerError)
	}
	user.HashedRt = refreshToken
	if err := a.UserRepository.Update(tx, user); err != nil {
		a.Log.WithError(err).Warnf("Failed to save user data: %+v", err)
		return nil, domain.NewError(fiber.StatusInternalServerError)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, domain.NewError(fiber.StatusInternalServerError)
	}

	return &dto.LoginResponse{
		User: dto.CredentialData{
			FullName: user.FullName,
			Email:    user.Email,
		},
		Token: dto.TokenData{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}

// Logout implements domain.AuthUseCase.
func (a *AuthUseCase) Logout(ctx context.Context, userID int64) error {
	c, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	tx := a.DB.WithContext(c).Begin()
	defer tx.Rollback()
	user := new(domain.UserEntity)
	if err := a.UserRepository.FindByID(tx, user, userID); err != nil {
		return domain.NewError(fiber.StatusNotFound, "User not found")
	}

	if user.HashedRt == "" {
		return domain.NewError(fiber.StatusUnauthorized, "User is not authorized")
	}

	user.HashedRt = ""
	if err := a.UserRepository.Update(tx, user); err != nil {
		a.Log.WithError(err).Warnf("Failed to save user data: %+v", err)
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		return domain.NewError(fiber.StatusInternalServerError)
	}

	return nil
}

// Refresh implements domain.AuthUseCase.
func (a *AuthUseCase) Refresh(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.LoginResponse, error) {
	c, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if validationErrors := util.Validate(a.Validate, req); len(validationErrors) > 0 {
		return nil, domain.NewError(fiber.StatusBadRequest, "The provided data is invalid", validationErrors)
	}

	userID, err := a.JWT.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, domain.NewError(fiber.StatusBadRequest, "Invalid refresh token")
	}

	tx := a.DB.WithContext(c).Begin()
	defer tx.Rollback()
	user := new(domain.UserEntity)
	if err := a.UserRepository.FindByID(tx, user, userID); err != nil {
		return nil, domain.NewError(fiber.StatusNotFound, "User not found")
	}
	if user.HashedRt != req.RefreshToken {
		return nil, domain.NewError(fiber.StatusUnauthorized, "Invalid refresh token")
	}

	accessToken, refreshToken, err := a.JWT.GenerateToken(user.ID)
	if err != nil {
		a.Log.WithError(err).Error("Failed to generate tokens")
		return nil, domain.NewError(fiber.StatusInternalServerError)
	}
	user.HashedRt = refreshToken
	if err := a.UserRepository.Update(tx, user); err != nil {
		a.Log.WithError(err).Warnf("Failed to save user data: %+v", err)
		return nil, domain.NewError(fiber.StatusInternalServerError)
	}
	if err := tx.Commit().Error; err != nil {
		return nil, domain.NewError(fiber.StatusInternalServerError)
	}

	return &dto.LoginResponse{
		User: dto.CredentialData{
			FullName: user.FullName,
			Email:    user.Email,
		},
		Token: dto.TokenData{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}
