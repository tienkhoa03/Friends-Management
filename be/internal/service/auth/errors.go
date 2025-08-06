package service

import "errors"

var (
	ErrAlreadyRegistered     = errors.New("email has already been registed")
	ErrInvalidLoginRequest   = errors.New("email or password is incorrect")
	ErrInvalidRefreshToken   = errors.New("invalid refresh token")
	ErrRefreshTokenIsRevoked = errors.New("refresh token is revoked")
	ErrRefreshTokenExpires   = errors.New("refresh token has expired")
	ErrInvalidSigningMethod  = errors.New("unexpected signing method")
)
