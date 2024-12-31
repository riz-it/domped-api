package dto

// Request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RegisterRequest struct {
	FullName string `json:"full_name" validate:"required"`
	Phone    string `json:"phone" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type EmailVerificationRequest struct {
	ReferenceID string `json:"reference_id" validate:"required"`
	OTP         string `json:"otp" validate:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// Response
type LoginResponse struct {
	User  CredentialData `json:"user"`
	Token TokenData      `json:"token"`
}

type RegisterResponse struct {
	ReferenceID string `json:"reference_id"`
}

// Data
type CredentialData struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

type TokenData struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
