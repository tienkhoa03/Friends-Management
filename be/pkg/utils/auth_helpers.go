package utils

import (
	"BE_Friends_Management/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	AccessTokenExpiredTime  = time.Hour
	RefreshTokenExpiredTime = 10 * 24 * time.Hour
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

func ParseAccessToken(rawAccessToken string) (*Claims, error) {
	accessToken, err := jwt.ParseWithClaims(rawAccessToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSigningMethod
		}
		return []byte(config.AccessSecret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := accessToken.Claims.(*Claims)
	if !ok || !accessToken.Valid {
		return nil, ErrInvalidAccessToken
	}
	return claims, nil
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
		return nil, ErrInvalidRefreshToken
	}
	return claims, nil
}
