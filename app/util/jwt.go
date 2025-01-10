package util

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"riz.it/domped/app/config"
	"riz.it/domped/app/domain"
)

type JWTUtil struct {
	Config *config.Config
}

func NewJWTUtil(config *config.Config) domain.JWT {
	return &JWTUtil{
		Config: config,
	}
}

// GenerateTokens implements domain.JWTHelper.
func (j *JWTUtil) GenerateToken(userID int64) (string, string, error) {
	iss := j.Config.Server.Host
	aud := j.Config.Server.Host
	accessExpTime, _ := strconv.Atoi(j.Config.Jwt.AccessTokenExp)
	refreshExpTime, _ := strconv.Atoi(j.Config.Jwt.RefreshTokenExp)
	exp := time.Now().Add(time.Hour * time.Duration(accessExpTime)).Unix()
	expRefresh := time.Now().Add(time.Hour * time.Duration(refreshExpTime)).Unix()
	sub := userID

	accessKey := j.Config.Jwt.AccessTokenKey
	refreshKey := j.Config.Jwt.RefreshTokenKey

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss": iss,
			"aud": aud,
			"exp": exp,
			"sub": sub,
		})

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss": iss,
			"aud": aud,
			"exp": expRefresh,
			"sub": sub,
		})

	signedAccessToken, err := accessToken.SignedString([]byte(accessKey))
	if err != nil {
		return "", "", err
	}

	signedRefreshToken, err := refreshToken.SignedString([]byte(refreshKey))
	if err != nil {
		return "", "", err
	}

	return signedAccessToken, signedRefreshToken, nil

}

// ValidateAccessToken implements domain.JWTHelper.
func (j *JWTUtil) ValidateAccessToken(tokenString string) (int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrInvalidKey
		}
		return []byte(j.Config.Jwt.AccessTokenKey), nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Konversi sub ke float64, lalu ubah ke int64
		if sub, ok := claims["sub"].(float64); ok {
			userID := int64(sub)
			return userID, nil
		}
		return 0, jwt.ErrInvalidKey
	}

	return 0, jwt.ErrInvalidKey
}

// ValidateRefreshToken implements domain.JWTHelper.
func (j *JWTUtil) ValidateRefreshToken(tokenString string) (int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrInvalidKey
		}
		return []byte(j.Config.Jwt.RefreshTokenKey), nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Konversi sub ke float64, lalu ubah ke int64
		if sub, ok := claims["sub"].(float64); ok {
			userID := int64(sub)
			return userID, nil
		}
		return 0, jwt.ErrInvalidKey
	}

	return 0, jwt.ErrInvalidKey
}
