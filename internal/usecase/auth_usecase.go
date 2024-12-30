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
	// Create a context with a timeout of 10 seconds
	c, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Validate the input request
	if validationErrors := util.Validate(a.Validate, req); len(validationErrors) > 0 {
		return nil, domain.NewError(fiber.StatusBadRequest, "The provided data is invalid", validationErrors)
	}

	// Initialize a database transaction with the context
	tx := a.DB.WithContext(c).Begin()
	user := new(domain.UserEntity)

	// Check if a user exists with the given email
	if err := a.UserRepository.FindByEmail(tx, user, req.Email); err != nil {
		return nil, domain.NewError(fiber.StatusUnauthorized, "Invalid email or password")
	}

	// Verify the provided password
	if !util.VerifyPassword(user.Password, req.Password) {
		return nil, domain.NewError(fiber.StatusUnauthorized, "Invalid email or password")
	}

	// Generate access and refresh tokens
	accessToken, refreshToken, err := a.JWT.GenerateToken(user.ID)
	if err != nil {
		a.Log.WithError(err).Error("Failed to generate tokens")
		return nil, domain.NewError(fiber.StatusInternalServerError)
	}

	// Store the hashed refresh token in the database
	user.HashedRt = refreshToken
	if err := a.UserRepository.Update(tx, user); err != nil {
		a.Log.WithError(err).Warnf("Failed to save user data: %+v", err)
		return nil, domain.NewError(fiber.StatusInternalServerError)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, domain.NewError(fiber.StatusInternalServerError)
	}

	// Return the login response containing user credentials and tokens
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
func (a *AuthUseCase) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	panic("unimplemented")
}

// Logout implements domain.AuthUseCase.
func (a *AuthUseCase) Logout(ctx context.Context, userID uint) error {
	panic("unimplemented")
}

// Refresh implements domain.AuthUseCase.
func (a *AuthUseCase) Refresh(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.LoginResponse, error) {
	panic("unimplemented")
}
