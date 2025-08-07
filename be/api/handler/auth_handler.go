package handler

import (
	"BE_Friends_Management/constant"
	"BE_Friends_Management/internal/domain/dto"
	service "BE_Friends_Management/internal/service/auth"
	"BE_Friends_Management/pkg"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(service service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

// User godoc
// @Summary      Register new user
// @Description  Register new user
// @Tags         Auth
// @Accept multipart/form-data
// @Produce      json
// @Param 		 email formData string true "User Email"
// @Param 		 password formData string true "User Password"
// @Router       /api/auth/register [POST]
// @Success      200   {object}  dto.ApiResponseSuccessStruct
func (h *AuthHandler) RegisterUser(c *gin.Context) {
	defer pkg.PanicHandler(c)
	userEmail := c.PostForm("email")
	userPassword := c.PostForm("password")
	newUser, err := h.service.RegisterUser(userEmail, userPassword)
	if err != nil {
		log.Error("Happened error when registing new user. Error: ", err)
		switch {
		case errors.Is(err, service.ErrAlreadyRegistered):
			pkg.PanicExeption(constant.Conflict, err.Error())
		default:
			pkg.PanicExeption(constant.UnknownError, "Happened error when registing new user")
		}
	}
	c.JSON(http.StatusOK, pkg.BuildResponseSuccess(constant.Success, newUser))
}

// User godoc
// @Summary      Login
// @Description  Login
// @Tags         Auth
// @Accept 		 json
// @Produce      json
// @Param 		 request body dto.LoginRequest true "User's email and password"
// @Router       /api/auth/login [POST]
// @Success      200   {object}  dto.ApiResponseSuccessStruct
func (h *AuthHandler) Login(c *gin.Context) {
	defer pkg.PanicHandler(c)
	var request dto.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Error("Happened error when mapping request. Error: ", err)
		pkg.PanicExeption(constant.InvalidRequest, "Invalid request format.")
	}
	accessToken, refreshToken, err := h.service.Login(request.Email, request.Password)
	if err != nil {
		log.Error("Happened error when creating new friendship. Error: ", err)
		switch {
		case errors.Is(err, service.ErrInvalidLoginRequest):
			pkg.PanicExeption(constant.InvalidRequest, err.Error())
		default:
			pkg.PanicExeption(constant.UnknownError, "Happened error when creating new friendship.")
		}
	}
	c.JSON(http.StatusOK, pkg.BuildResponseSuccessWithTokens(accessToken, refreshToken))
}

// User godoc
// @Summary      Refresh Access Token
// @Description  Refresh Access Token
// @Tags         Auth
// @Accept 		 json
// @Produce      json
// @Param 		 request body dto.RefreshRequest true "User's refresh token"
// @Router       /api/auth/refresh [POST]
// @Success      200   {object}  dto.ApiResponseSuccessStruct
func (h *AuthHandler) RefreshAccessToken(c *gin.Context) {
	defer pkg.PanicHandler(c)
	var request dto.RefreshRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Error("Happened error when mapping request. Error: ", err)
		pkg.PanicExeption(constant.InvalidRequest, "Invalid request format.")
	}
	accessToken, err := h.service.RefreshAccessToken(request.RefreshToken)
	if err != nil {
		log.Error("Happened error when creating new friendship. Error: ", err)
		switch {
		case errors.Is(err, service.ErrInvalidRefreshToken):
			pkg.PanicExeption(constant.Unauthorized, err.Error())
		case errors.Is(err, service.ErrRefreshTokenIsRevoked):
			pkg.PanicExeption(constant.Unauthorized, err.Error())
		case errors.Is(err, service.ErrRefreshTokenExpires):
			pkg.PanicExeption(constant.Unauthorized, err.Error())
		case errors.Is(err, service.ErrInvalidSigningMethod):
			pkg.PanicExeption(constant.Unauthorized, err.Error())
		default:
			pkg.PanicExeption(constant.UnknownError, "Happened error when refreshing access token.")
		}
	}
	c.JSON(http.StatusOK, pkg.BuildResponseSuccessWithAccessToken(accessToken))
}

// User godoc
// @Summary      Logout
// @Description  Logout
// @Tags         Auth
// @Accept 		 json
// @Produce      json
// @Param 		 request body dto.LogoutRequest true "User's refresh token"
// @Router       /api/auth/logout [POST]
// @Success      200   {object}  dto.ApiResponseSuccessStruct
func (h *AuthHandler) Logout(c *gin.Context) {
	defer pkg.PanicHandler(c)
	var request dto.LogoutRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Error("Happened error when mapping request. Error: ", err)
		pkg.PanicExeption(constant.InvalidRequest, "Invalid request format.")
	}
	err := h.service.Logout(request.RefreshToken)
	if err != nil {
		log.Error("Happened error when creating new friendship. Error: ", err)
		switch {
		case errors.Is(err, service.ErrInvalidRefreshToken):
			pkg.PanicExeption(constant.Unauthorized, err.Error())
		case errors.Is(err, service.ErrRefreshTokenIsRevoked):
			pkg.PanicExeption(constant.Unauthorized, err.Error())
		case errors.Is(err, service.ErrRefreshTokenExpires):
			pkg.PanicExeption(constant.Unauthorized, err.Error())
		case errors.Is(err, service.ErrInvalidSigningMethod):
			pkg.PanicExeption(constant.Unauthorized, err.Error())
		default:
			pkg.PanicExeption(constant.UnknownError, "Happened error when refreshing access token.")
		}
	}
	c.JSON(http.StatusOK, pkg.BuildResponseSuccessNoData())
}
