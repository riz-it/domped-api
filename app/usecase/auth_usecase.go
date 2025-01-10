package usecase

import (
	"context"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"riz.it/domped/app/domain"
	"riz.it/domped/app/dto"
	"riz.it/domped/app/util"
)

type AuthUseCase struct {
	DB               *gorm.DB
	Log              *logrus.Logger
	UserRepository   domain.UserRepository
	JWT              domain.JWT
	Validate         *validator.Validate
	Redis            *redis.Client
	Email            domain.Email
	WalletRepository domain.WalletRepository
}

func NewAuthUseCase(db *gorm.DB, log *logrus.Logger, userRepository domain.UserRepository, walletRepository domain.WalletRepository, jwt domain.JWT, validate *validator.Validate, redis *redis.Client, email domain.Email) domain.AuthUseCase {
	return &AuthUseCase{
		DB:               db,
		Log:              log,
		UserRepository:   userRepository,
		WalletRepository: walletRepository,
		JWT:              jwt,
		Validate:         validate,
		Redis:            redis,
		Email:            email,
	}
}

// Login implements domain.AuthUseCase.
func (a *AuthUseCase) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	// Set a timeout for the database context
	c, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Validate the login request
	if validationErrors := util.Validate(a.Validate, req); len(validationErrors) > 0 {
		// Return validation errors if request data is invalid
		return nil, domain.NewError(fiber.StatusBadRequest, "The provided data is invalid", validationErrors)
	}

	// Start a new database transaction
	tx := a.DB.WithContext(c).Begin()
	defer tx.Rollback() // Ensure rollback if an error occurs

	// Find the user by email
	user := new(domain.UserEntity)
	if err := a.UserRepository.FindByEmail(tx, user, req.Email); err != nil {
		// Return an error if the user is not found
		return nil, domain.NewError(fiber.StatusUnauthorized, "Invalid email or password")
	}

	// Verify the password is correct
	if !util.VerifyPassword(user.Password, req.Password) {
		// Return an error if the password is invalid
		return nil, domain.NewError(fiber.StatusUnauthorized, "Invalid email or password")
	}

	// Ensure the user's account is active
	if !user.IsActive {
		// Return an error if the account is not active
		return nil, domain.NewError(fiber.StatusUnauthorized, "Account has not been verified yet")
	}

	// Generate the access and refresh tokens
	accessToken, refreshToken, err := a.JWT.GenerateToken(user.ID)
	if err != nil {
		a.Log.WithError(err).Error("Failed to generate tokens")
		return nil, domain.NewError(fiber.StatusInternalServerError)
	}

	// Store the refresh token in the user's record
	user.HashedRt = refreshToken
	if err := a.UserRepository.Update(tx, user); err != nil {
		a.Log.WithError(err).Warnf("Failed to save user data: %+v", err)
		return nil, domain.NewError(fiber.StatusInternalServerError)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		a.Log.WithError(err).Warnf("Failed to commit transaction: %+v", err)
		return nil, domain.NewError(fiber.StatusInternalServerError)
	}

	// Return the successful login response
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
	// Set a timeout for the registration process
	c, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Validate the incoming request data
	if validationErrors := util.Validate(a.Validate, req); len(validationErrors) > 0 {
		return nil, domain.NewError(fiber.StatusBadRequest, "Invalid data provided", validationErrors)
	}

	// Start a new transaction to ensure atomicity
	tx := a.DB.WithContext(c).Begin()
	defer tx.Rollback()

	// Check if the email is already in use
	count, err := a.UserRepository.CountByEmail(tx, req.Email)
	if err != nil {
		a.Log.WithError(err).Warnf("Failed to count email: %+v", err)
		return nil, domain.NewError(fiber.StatusInternalServerError)
	}
	if count > 0 {
		return nil, domain.NewError(fiber.StatusConflict, "The email address is already in use")
	}

	// Hash the password for storage
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		a.Log.WithError(err).Warnf("Failed to hash password: %+v", err)
		return nil, domain.NewError(fiber.StatusInternalServerError)
	}

	// Create a new user entity
	user := &domain.UserEntity{
		Password: hashedPassword,
		FullName: req.FullName,
		Email:    req.Email,
	}

	// Save the new user to the database
	if err := a.UserRepository.Create(tx, user); err != nil {
		a.Log.WithError(err).Warnf("Failed to create user: %+v", err)
		return nil, domain.NewError(fiber.StatusInternalServerError)
	}

	// Generate a unique OTP reference ID and code
	otpReferenceId := util.GenerateUUID() + "-" + strconv.FormatInt(user.ID, 10)
	otpCode := util.GenerateRandomCode(4)
	ttl := time.Minute * 5

	// Store the OTP code in Redis with a TTL (Time to Live)
	insertOtp := a.Redis.Set(ctx, otpReferenceId, otpCode, ttl)
	if err := insertOtp.Err(); err != nil {
		a.Log.WithError(err).Warnf("Failed to store OTP in Redis: %+v", err)
		return nil, domain.NewError(fiber.StatusInternalServerError)
	}

	// Send the OTP code to the user's email
	if err := a.Email.Send(user.Email, "Your OTP Code", "Your OTP code is: <b>"+otpCode+"</b>"); err != nil {
		a.Log.WithError(err).Warnf("Failed to send OTP email: %+v", err)
		return nil, domain.NewError(fiber.StatusInternalServerError)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		a.Log.WithError(err).Warnf("Failed to commit transaction: %+v", err)
		return nil, domain.NewError(fiber.StatusInternalServerError)
	}

	// Return the OTP reference ID as part of the registration response
	return &dto.RegisterResponse{
		ReferenceID: otpReferenceId,
	}, nil
}

// Logout implements domain.AuthUseCase.
func (a *AuthUseCase) Logout(ctx context.Context, userID int64) error {
	// Set a timeout for the logout process
	c, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Start a new transaction to ensure atomicity
	tx := a.DB.WithContext(c).Begin()
	defer tx.Rollback()

	// Retrieve the user based on userID
	user := new(domain.UserEntity)
	if err := a.UserRepository.FindByID(tx, user, userID); err != nil {
		// If user is not found, return a 'Not Found' error
		return domain.NewError(fiber.StatusNotFound, "User not found")
	}

	// Check if the user has a valid refresh token
	if user.HashedRt == "" {
		// If no refresh token exists, return an 'Unauthorized' error
		return domain.NewError(fiber.StatusUnauthorized, "User is not authorized")
	}

	// Remove the refresh token (effectively logging the user out)
	user.HashedRt = ""

	// Update the user record in the database to remove the refresh token
	if err := a.UserRepository.Update(tx, user); err != nil {
		a.Log.WithError(err).Warnf("Failed to save user data: %+v", err)
		return domain.NewError(fiber.StatusInternalServerError)
	}

	// Commit the transaction to persist the changes
	if err := tx.Commit().Error; err != nil {
		a.Log.WithError(err).Warnf("Failed to commit transaction: %+v", err)
		return domain.NewError(fiber.StatusInternalServerError)
	}

	// Successful logout, return nil indicating no errors
	return nil
}

// Refresh implements domain.AuthUseCase.
func (a *AuthUseCase) Refresh(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.LoginResponse, error) {
	// Set a timeout for the refresh process
	c, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Validate the incoming request data
	if validationErrors := util.Validate(a.Validate, req); len(validationErrors) > 0 {
		// Return an error if validation fails
		return nil, domain.NewError(fiber.StatusBadRequest, "Invalid request data", validationErrors)
	}

	// Validate the provided refresh token
	userID, err := a.JWT.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		// Return an error if the refresh token is invalid
		return nil, domain.NewError(fiber.StatusBadRequest, "Invalid refresh token")
	}

	// Start a new transaction to ensure atomicity
	tx := a.DB.WithContext(c).Begin()
	defer tx.Rollback()

	// Retrieve the user based on userID
	user := new(domain.UserEntity)
	if err := a.UserRepository.FindByID(tx, user, userID); err != nil {
		// If the user is not found, return a 'Not Found' error
		return nil, domain.NewError(fiber.StatusNotFound, "User not found")
	}

	// Ensure that the stored refresh token matches the provided one
	if user.HashedRt != req.RefreshToken {
		// If refresh token does not match, return an 'Unauthorized' error
		return nil, domain.NewError(fiber.StatusUnauthorized, "Invalid refresh token")
	}

	// Generate new access and refresh tokens for the user
	accessToken, refreshToken, err := a.JWT.GenerateToken(user.ID)
	if err != nil {
		a.Log.WithError(err).Error("Failed to generate new tokens")
		return nil, domain.NewError(fiber.StatusInternalServerError)
	}

	// Update the user record with the new refresh token
	user.HashedRt = refreshToken
	if err := a.UserRepository.Update(tx, user); err != nil {
		a.Log.WithError(err).Warnf("Failed to update user refresh token: %+v", err)
		return nil, domain.NewError(fiber.StatusInternalServerError)
	}

	// Commit the transaction to persist the changes
	if err := tx.Commit().Error; err != nil {
		a.Log.WithError(err).Warnf("Failed to commit transaction: %+v", err)
		return nil, domain.NewError(fiber.StatusInternalServerError)
	}

	// Return the new tokens and user details
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

// EmailVerification implements domain.AuthUseCase.
func (a *AuthUseCase) EmailVerification(ctx context.Context, req *dto.EmailVerificationRequest) (*dto.LoginResponse, error) {
	// Set a timeout for the email verification process
	c, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Validate the incoming request data
	if validationErrors := util.Validate(a.Validate, req); len(validationErrors) > 0 {
		// Return an error if validation fails
		return nil, domain.NewError(fiber.StatusBadRequest, "Invalid request data", validationErrors)
	}

	// Retrieve OTP from Redis using the provided ReferenceID
	getOtp := a.Redis.Get(c, req.ReferenceID)
	if err := getOtp.Err(); err != nil {
		// Return an error if the ReferenceID is invalid or the OTP retrieval fails
		return nil, domain.NewError(fiber.StatusBadRequest, "Invalid reference ID")
	}

	otp, err := getOtp.Result()
	if err != nil {
		a.Log.WithError(err).Error("Failed to get OTP response")
		return nil, domain.NewError(fiber.StatusInternalServerError)
	}

	// Validate the provided OTP
	if otp != req.OTP {
		// Return an error if the OTP does not match
		return nil, domain.NewError(fiber.StatusBadRequest, "Invalid OTP code")
	}

	// Extract userID from the ReferenceID
	userID, err := util.ExtractIDFromReference(req.ReferenceID)
	if err != nil {
		a.Log.WithError(err).Warnf("Failed to extract reference id: %+v", err)
		return nil, domain.NewError(fiber.StatusInternalServerError)
	}

	// Start a transaction to ensure atomicity of user update
	tx := a.DB.WithContext(c).Begin()
	defer tx.Rollback()

	// Retrieve the user by ID from the repository
	user := new(domain.UserEntity)
	if err := a.UserRepository.FindByID(tx, user, userID); err != nil {
		// Return an error if user is not found
		return nil, domain.NewError(fiber.StatusNotFound, "User not found")
	}

	// Generate new access and refresh tokens for the user
	accessToken, refreshToken, err := a.JWT.GenerateToken(user.ID)
	if err != nil {
		a.Log.WithError(err).Error("Failed to generate tokens")
		return nil, domain.NewError(fiber.StatusInternalServerError)
	}

	// Update user data: mark the user as active and set the email verification timestamp
	user.IsActive = true
	now := time.Now()
	user.EmailVerifiedAt = &now
	user.HashedRt = refreshToken

	// Generate wallet number
	walletNumber, err := util.GenerateWalletNumber(8)
	if err != nil {
		a.Log.WithError(err).Error("Failed to generate wallet number")
		return nil, domain.NewError(fiber.StatusInternalServerError)
	}

	// Check if wallet number already exists
	count, err := a.WalletRepository.CountByWalletNumber(tx, walletNumber)
	if err != nil {
		a.Log.WithError(err).Warnf("Failed to count wallet: %+v", err)
		return nil, domain.NewError(fiber.StatusInternalServerError)
	}
	if count > 0 {
		walletNumber, _ = util.GenerateWalletNumber(8)
	}

	user.Wallet = domain.WalletEntity{
		Balance:      0,
		WalletNumber: walletNumber,
	}

	// Update the user record in the repository
	if err := a.UserRepository.Update(tx, user); err != nil {
		a.Log.WithError(err).Warnf("Failed to save user data: %+v", err)
		return nil, domain.NewError(fiber.StatusInternalServerError)
	}

	// Delete the OTP from Redis as it is no longer needed
	delOtp := a.Redis.Del(c, req.ReferenceID)
	if err := delOtp.Err(); err != nil {
		a.Log.WithError(err).Warnf("Failed to delete OTP: %+v", err)
		return nil, domain.NewError(fiber.StatusInternalServerError)
	}

	// Commit the transaction to persist changes
	if err := tx.Commit().Error; err != nil {
		a.Log.WithError(err).Warnf("Failed to commit transaction: %+v", err)
		return nil, domain.NewError(fiber.StatusInternalServerError)
	}

	// Return the new tokens and user details
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
