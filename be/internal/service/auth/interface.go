package service

import (
	"BE_Friends_Management/internal/domain/entity"
	"errors"
)

var (
	ErrAlreadyRegistered     = errors.New("email has already been registed")
	ErrInvalidLoginRequest   = errors.New("email or password is incorrect")
	ErrInvalidRefreshToken   = errors.New("invalid refresh token")
	ErrRefreshTokenIsRevoked = errors.New("refresh token is revoked")
	ErrRefreshTokenExpires   = errors.New("refresh token has expired")
	ErrInvalidSigningMethod  = errors.New("unexpected signing method")
)

//go:generate mockgen -source=interface.go -destination=../mock/mock_auth_service.go

type AuthService interface {
	RegisterUser(email, password string) (*entity.User, error)
	Login(email, password string) (string, string, error)
	RefreshAccessToken(rawRefreshToken string) (string, error)
	Logout(rawRefreshToken string) error
}
