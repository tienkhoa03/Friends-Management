package utils

import "errors"

var (
	ErrInvalidSigningMethod = errors.New("unexpected signing method")
	ErrInvalidRefreshToken  = errors.New("invalid refresh token")
	ErrInvalidAccessToken   = errors.New("invalid access token")
)
