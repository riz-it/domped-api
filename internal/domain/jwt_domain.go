package domain

type JWT interface {
	GenerateToken(userID int64) (string, string, error)
	ValidateAccessToken(tokenString string) (int64, error)
	ValidateRefreshToken(tokenString string) (int64, error)
}
