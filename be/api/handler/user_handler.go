package handler

import (
	"BE_Friends_Management/constant"
	service "BE_Friends_Management/internal/service/users"
	"BE_Friends_Management/pkg"
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
// @Tags         Users
// @Accept       json
// @Produce      json
// @Router       /api/users [GET]
// @Success      200   {object}  dto.ApiResponseSuccessStruct
func (h *UserHandler) GetAllUser(c *gin.Context) {
	defer pkg.PanicHandler(c)
	users, err := h.service.GetAllUser()
	if err != nil {
		log.Error("Happened error when getting all users. Error: ", err)
		pkg.PanicExeption(constant.UnknownError, "Happened error when getting all users")
	}
	c.JSON(http.StatusOK, pkg.BuildReponseSuccess(constant.Success, users))
}

// User godoc
// @Summary      Get user by ID
// @Description  Get user by ID
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param 		 id path string true "User ID"
// @Router       /api/users/{id} [GET]
// @Success      200   {object}  dto.ApiResponseSuccessStruct
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
		pkg.PanicExeption(constant.UnknownError, "Happened error when getting the user by ID")
	}
	c.JSON(http.StatusOK, pkg.BuildReponseSuccess(constant.Success, user))
}

// User godoc
// @Summary      Create new user
// @Description  Create new user
// @Tags         Users
// @Accept multipart/form-data
// @Produce      json
// @Param 		 email formData string true "User Email"
// @Router       /api/users [POST]
// @Success      200   {object}  dto.ApiResponseSuccessStruct
func (h *UserHandler) CreateUser(c *gin.Context) {
	defer pkg.PanicHandler(c)
	userEmail := c.PostForm("email")
	newUser, err := h.service.CreateUser(userEmail)
	if err != nil {
		log.Error("Happened error when creating new user. Error: ", err)
		pkg.PanicExeption(constant.UnknownError, "Happened error when creating new user")
	}
	c.JSON(http.StatusOK, pkg.BuildReponseSuccess(constant.Success, newUser))
}

// User godoc
// @Summary      Delete user
// @Description  Delete user
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param 		 id path string true "User ID"
// @Router       /api/users/{id} [DELETE]
// @Success      200   {object}  dto.ApiResponseSuccessStruct
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
		pkg.PanicExeption(constant.UnknownError, "Happened error when deleting a user")
	}
	c.JSON(http.StatusOK, pkg.BuildReponseSuccessNoData())
}

// User godoc
// @Summary      Update user
// @Description  Update user
// @Tags         Users
// @Accept multipart/form-data
// @Produce      json
// @Param		id	path		string				true	"id"
// @Param 		 email formData string true "User's New Email"
// @Router       /api/users/{id} [PUT]
// @Success      200   {object}  dto.ApiResponseSuccessStruct
func (h *UserHandler) UpdateUser(c *gin.Context) {
	defer pkg.PanicHandler(c)
	userIdStr := c.Param("id")
	email := c.PostForm("email")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		log.Error("Happened error when converting userId to int64. Error: ", err)
		pkg.PanicExeption(constant.InvalidRequest, "Happened error when converting userId to int64")
	}
	updatedUser, err := h.service.UpdateUser(userId, email)
	if err != nil {
		log.Error("Happened error when updating user. Error: ", err)
		pkg.PanicExeption(constant.UnknownError, "Happened error when updating a user")
	}
	c.JSON(http.StatusOK, pkg.BuildReponseSuccess(constant.Success, updatedUser))
}
