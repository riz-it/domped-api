package dto

// Request

// Response
type UserData struct {
	ID              int64  `json:"id"`
	FullName        string `json:"full_name"`
	Phone           string `json:"phone"`
	Email           string `json:"email"`
	EmailVerifiedAt string `json:"email_verified_at"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
	IsActive        bool   `json:"is_active"`
}
