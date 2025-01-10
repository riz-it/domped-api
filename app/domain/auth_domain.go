package domain

import (
	"context"

	"riz.it/domped/app/dto"
)

// Interface
type AuthUseCase interface {
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error)
	EmailVerification(ctx context.Context, req *dto.EmailVerificationRequest) (*dto.LoginResponse, error)
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error)
	Logout(ctx context.Context, userID int64) error
	Refresh(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.LoginResponse, error)
}
