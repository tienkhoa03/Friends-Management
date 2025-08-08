package service

import (
	"BE_Friends_Management/internal/domain/entity"
	auth "BE_Friends_Management/internal/repository/auth"
	users "BE_Friends_Management/internal/repository/users"
	"BE_Friends_Management/pkg/utils"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"

	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	repo     auth.AuthRepository
	userRepo users.UserRepository
}

func NewAuthService(repo auth.AuthRepository, userRepo users.UserRepository) AuthService {
	return &authService{repo: repo, userRepo: userRepo}
}

func (service *authService) RegisterUser(email, password string) (*entity.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := entity.User{Email: email, Password: string(hashedPassword)}
	newUser, err := service.userRepo.CreateUser(&user)
	if err != nil && strings.Contains(err.Error(), "duplicate key") {
		return nil, ErrAlreadyRegistered
	}
	if err != nil {
		return nil, err
	}
	return newUser, nil
}

func (service *authService) Login(email, password string) (string, string, error) {
	user, err := service.userRepo.GetUserByEmail(email)
	if err != nil {
		return "", "", ErrInvalidLoginRequest
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", "", ErrInvalidLoginRequest
	}
	accessTokenExpiredTime := time.Now().Add(utils.AccessTokenExpiredTime)
	accessToken, err := utils.GenerateAccessToken(user.Id, accessTokenExpiredTime)
	if err != nil {
		return "", "", err
	}
	refreshTokenExpiredTime := time.Now().Add(utils.RefreshTokenExpiredTime)
	refreshToken, err := utils.GenerateRefreshToken(user.Id, refreshTokenExpiredTime)
	if err != nil {
		return "", "", err
	}
	tokenRecord := &entity.UserToken{
		UserId:       user.Id,
		RefreshToken: refreshToken,
		ExpiresAt:    refreshTokenExpiredTime,
		IsRevoked:    false,
	}
	err = service.repo.CreateToken(tokenRecord)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func (service *authService) RefreshAccessToken(rawRefreshToken string) (string, error) {
	userToken, err := service.repo.FindByRefreshToken(rawRefreshToken)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", ErrInvalidRefreshToken
	}
	if err != nil {
		return "", err
	}
	if userToken.IsRevoked {
		return "", ErrRefreshTokenIsRevoked
	}
	claims, err := utils.ParseRefreshToken(rawRefreshToken)
	if errors.Is(err, utils.ErrInvalidRefreshToken) {
		return "", ErrInvalidRefreshToken
	}
	if errors.Is(err, utils.ErrInvalidSigningMethod) {
		return "", ErrInvalidSigningMethod
	}
	if err != nil {
		return "", err
	}
	if claims.ExpiresAt.Time.Before(time.Now()) {
		return "", ErrRefreshTokenExpires
	}
	accessToken, err := utils.GenerateAccessToken(userToken.UserId, time.Now().Add(utils.AccessTokenExpiredTime))
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

func (service *authService) Logout(rawRefreshToken string) error {
	userToken, err := service.repo.FindByRefreshToken(rawRefreshToken)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrInvalidRefreshToken
	}
	if err != nil {
		return err
	}
	if userToken.IsRevoked {
		return ErrRefreshTokenIsRevoked
	}
	claims, err := utils.ParseRefreshToken(rawRefreshToken)
	if errors.Is(err, utils.ErrInvalidRefreshToken) {
		return ErrInvalidRefreshToken
	}
	if errors.Is(err, utils.ErrInvalidSigningMethod) {
		return ErrInvalidSigningMethod
	}
	if err != nil {
		return err
	}
	if claims.ExpiresAt.Time.Before(time.Now()) {
		return ErrRefreshTokenExpires
	}
	err = service.repo.SetRefreshTokenIsRevoked(rawRefreshToken)
	if err != nil {
		return err
	}
	return nil
}
