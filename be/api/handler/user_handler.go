package handler

import (
	"BE_Friends_Management/constant"
	service "BE_Friends_Management/internal/service/users"
	"BE_Friends_Management/pkg"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// User godoc
// @Summary      Get all user
// @Description   Get all user
// @Tags         Users Management
// @Accept       json
// @Produce      json
// @param Authorization header string true "Authorization"
// @Router       /api/users [GET]
// @Success      200   {object}  dto.ApiResponseSuccessStruct
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func (h *UserHandler) GetAllUser(c *gin.Context) {
	defer pkg.PanicHandler(c)
	users, err := h.service.GetAllUser()
	if err != nil {
		log.Error("Happened error when getting all users. Error: ", err)
		pkg.PanicExeption(constant.UnknownError, "Happened error when getting all users")
	}
	c.JSON(http.StatusOK, pkg.BuildResponseSuccess(constant.Success, users))
}

// User godoc
// @Summary      Get user by ID
// @Description  Get user by ID
// @Tags         Users Management
// @Accept       json
// @Produce      json
// @Param 		 id path string true "User ID"
// @param Authorization header string true "Authorization"
// @Router       /api/users/{id} [GET]
// @Success      200   {object}  dto.ApiResponseSuccessStruct
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func (h *UserHandler) GetUserById(c *gin.Context) {
	defer pkg.PanicHandler(c)
	userIdStr := c.Param("id")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		log.Error("Happened error when converting userId to int64. Error: ", err)
		pkg.PanicExeption(constant.InvalidRequest, "Happened error when converting userId to int64")
	}
	user, err := h.service.GetUserById(userId)
	if err != nil {
		log.Error("Happened error when getting the user by ID. Error: ", err)
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			pkg.PanicExeption(constant.DataNotFound, err.Error())
		default:
			pkg.PanicExeption(constant.UnknownError, "Happened error when getting the user by ID")
		}
	}
	c.JSON(http.StatusOK, pkg.BuildResponseSuccess(constant.Success, user))
}

// User godoc
// @Summary      Delete user
// @Description  Delete user
// @Tags         Users Management
// @Accept       json
// @Produce      json
// @Param 		 id path string true "User ID"
// @param Authorization header string true "Authorization"
// @Router       /api/users/{id} [DELETE]
// @Success      200   {object}  dto.ApiResponseSuccessStruct
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func (h *UserHandler) DeleteUserById(c *gin.Context) {
	defer pkg.PanicHandler(c)
	userIdStr := c.Param("id")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		log.Error("Happened error when converting userId to int64. Error: ", err)
		pkg.PanicExeption(constant.InvalidRequest, "Happened error when converting userId to int64")
	}
	err = h.service.DeleteUserById(userId)
	if err != nil {
		log.Error("Happened error when deleting a user. Error: ", err)
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			pkg.PanicExeption(constant.DataNotFound, err.Error())
		default:
			pkg.PanicExeption(constant.UnknownError, "Happened error when deleting a user")
		}
	}
	c.JSON(http.StatusOK, pkg.BuildResponseSuccessNoData())
}

// User godoc
// @Summary      Update user
// @Description  Update user
// @Tags         Users Management
// @Accept multipart/form-data
// @Produce      json
// @Param		id	path		string				true	"id"
// @Param 		 email formData string true "User's New Email"
// @Param 		 password formData string true "User's New Password"
// @param Authorization header string true "Authorization"
// @Router       /api/users/{id} [PUT]
// @Success      200   {object}  dto.ApiResponseSuccessStruct
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func (h *UserHandler) UpdateUser(c *gin.Context) {
	defer pkg.PanicHandler(c)
	userIdStr := c.Param("id")
	email := c.PostForm("email")
	password := c.PostForm("password")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		log.Error("Happened error when converting userId to int64. Error: ", err)
		pkg.PanicExeption(constant.InvalidRequest, "Happened error when converting userId to int64")
	}
	updatedUser, err := h.service.UpdateUser(userId, email, password)
	if err != nil {
		log.Error("Happened error when updating user. Error: ", err)
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			pkg.PanicExeption(constant.DataNotFound, err.Error())
		default:
			pkg.PanicExeption(constant.UnknownError, "Happened error when updating a user")
		}
	}
	c.JSON(http.StatusOK, pkg.BuildResponseSuccess(constant.Success, updatedUser))
}
