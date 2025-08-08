package handler

import (
	"BE_Friends_Management/constant"
	"BE_Friends_Management/internal/domain/dto"
	service "BE_Friends_Management/internal/service/friendship"
	"BE_Friends_Management/pkg"
	"BE_Friends_Management/pkg/utils"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type FriendshipHandler struct {
	service service.FriendshipService
}

func NewFriendshipHandler(service service.FriendshipService) *FriendshipHandler {
	return &FriendshipHandler{service: service}
}

// Friendship godoc
// @Summary      Create new friendship
// @Description  Create new friendship
// @Tags         Friendship
// @Accept 		json
// @Produce      json
// @Param 		 request body dto.CreateFriendshipRequest true "List of 2 friend emails"  example({"friends": ["andy@example.com", "john@example.com"]})
// @param Authorization header string true "Authorization"
// @Router       /api/friendship [POST]
// @Success      200   {object}  dto.ApiResponseSuccessNoData
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func (h *FriendshipHandler) CreateFriendship(c *gin.Context) {
	defer pkg.PanicHandler(c)
	authUserId := utils.GetAuthUserId(c)
	var request dto.CreateFriendshipRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Error("Happened error when mapping request. Error: ", err)
		pkg.PanicExeption(constant.InvalidRequest, "Invalid request format.")
	}
	err := h.service.CreateFriendship(authUserId, request.Friends[0], request.Friends[1])
	if err != nil {
		log.Error("Happened error when creating new friendship. Error: ", err)
		switch {
		case errors.Is(err, service.ErrInvalidRequest):
			pkg.PanicExeption(constant.InvalidRequest, err.Error())
		case errors.Is(err, service.ErrUserNotFound):
			pkg.PanicExeption(constant.DataNotFound, err.Error())
		case errors.Is(err, service.ErrAlreadyFriend):
			pkg.PanicExeption(constant.Conflict, err.Error())
		case errors.Is(err, service.ErrIsBlocked):
			pkg.PanicExeption(constant.StatusForbidden, err.Error())
		case errors.Is(err, service.ErrNotPermitted):
			pkg.PanicExeption(constant.StatusForbidden, err.Error())
		default:
			pkg.PanicExeption(constant.UnknownError, "Happened error when creating new friendship.")
		}
	}
	c.JSON(http.StatusOK, pkg.BuildResponseSuccessNoData())
}

// Friendship godoc
// @Summary      Retrieve friends list for an email address
// @Description  Retrieve friends list for an email address
// @Tags         Friendship
// @Accept 		json
// @Produce      json
// @Param 		 email query string true "Email address"
// @param Authorization header string true "Authorization"
// @Router       /api/friendship/friends [GET]
// @Success      200   {object}  dto.ApiResponseSuccessWithFriendsList
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func (h *FriendshipHandler) RetrieveFriendsList(c *gin.Context) {
	defer pkg.PanicHandler(c)
	authUserId := utils.GetAuthUserId(c)
	authUserRole := utils.GetAuthUserRole(c)
	requestEmail := c.Query("email")
	if requestEmail == "" {
		log.Error("Happened error when mapping request. Error: received no email input.")
		pkg.PanicExeption(constant.InvalidRequest, "Invalid request format.")
		return
	}
	friends, err := h.service.RetrieveFriendsList(authUserId, authUserRole, requestEmail)
	if err != nil {
		log.Error("Happened error when retrieving friends list. Error: ", err)
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			pkg.PanicExeption(constant.DataNotFound, err.Error())
		case errors.Is(err, service.ErrNotPermitted):
			pkg.PanicExeption(constant.StatusForbidden, err.Error())
		default:
			pkg.PanicExeption(constant.UnknownError, "Happened error when retrieving friends list.")
		}
	}
	emails := utils.ConvertUsersToEmails(friends)
	count := h.service.CountFriends(friends)
	c.JSON(http.StatusOK, pkg.BuildResponseSuccessWithFriendsList(emails, count))
}

// Friendship godoc
// @Summary      Retrieve common friends list between two email addresses
// @Description  Retrieve common friends list between two email addresses
// @Tags         Friendship
// @Accept 		json
// @Produce      json
// @Param 		 email1 query string true "Email address of user 1"
// @Param 		 email2 query string true "Email address of user 2"
// @param Authorization header string true "Authorization"
// @Router       /api/friendship/common-friends [GET]
// @Success      200   {object}  dto.ApiResponseSuccessWithFriendsList
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func (h *FriendshipHandler) RetrieveCommonFriends(c *gin.Context) {
	defer pkg.PanicHandler(c)
	authUserId := utils.GetAuthUserId(c)
	authUserRole := utils.GetAuthUserRole(c)
	requestEmail1 := c.Query("email1")
	if requestEmail1 == "" {
		log.Error("Happened error when mapping request. Error: received no email input.")
		pkg.PanicExeption(constant.InvalidRequest, "Invalid request format.")
		return
	}
	requestEmail2 := c.Query("email2")
	if requestEmail2 == "" {
		log.Error("Happened error when mapping request. Error: received no email input.")
		pkg.PanicExeption(constant.InvalidRequest, "Invalid request format.")
		return
	}
	friends, err := h.service.RetrieveCommonFriends(authUserId, authUserRole, requestEmail1, requestEmail2)
	if err != nil {
		log.Error("Happened error when retrieving common friends list. Error: ", err)
		switch {
		case errors.Is(err, service.ErrInvalidRequest):
			pkg.PanicExeption(constant.InvalidRequest, err.Error())
		case errors.Is(err, service.ErrUserNotFound):
			pkg.PanicExeption(constant.DataNotFound, err.Error())
		case errors.Is(err, service.ErrNotPermitted):
			pkg.PanicExeption(constant.StatusForbidden, err.Error())
		default:
			pkg.PanicExeption(constant.UnknownError, "Happened error when retrieving common friends list.")
		}
	}
	emails := utils.ConvertUsersToEmails(friends)
	count := h.service.CountFriends(friends)
	c.JSON(http.StatusOK, pkg.BuildResponseSuccessWithFriendsList(emails, count))
}
