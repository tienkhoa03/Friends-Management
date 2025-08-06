package utils

import (
	"BE_Friends_Management/config"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidSigningMethod  = errors.New("unexpected signing method")
	ErrInvalidRefreshRequest = errors.New("invalid refresh token")
)

type Claims struct {
	UserId int64 `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(userId int64, expiredTime time.Time) (string, error) {
	claims := &Claims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiredTime),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessString, err := accessToken.SignedString([]byte(config.AccessSecret))
	if err != nil {
		return "", err
	}
	return accessString, nil
}

func GenerateRefreshToken(userId int64, expiredTime time.Time) (string, error) {
	claims := &Claims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiredTime),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshString, err := refreshToken.SignedString([]byte(config.RefreshSecret))
	if err != nil {
		return "", err
	}
	return refreshString, nil
}

func ParseRefreshToken(rawRefreshToken string) (*Claims, error) {
	refreshToken, err := jwt.ParseWithClaims(rawRefreshToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSigningMethod
		}
		return []byte(config.RefreshSecret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := refreshToken.Claims.(*Claims)
	if !ok || !refreshToken.Valid {
		return nil, ErrInvalidRefreshRequest
	}
	return claims, nil
}
