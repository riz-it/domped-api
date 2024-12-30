package dto

// Request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
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
	User CredentialData `json:"user"`
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
