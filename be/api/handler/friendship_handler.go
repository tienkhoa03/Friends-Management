package handler

import (
	"BE_Friends_Management/constant"
	"BE_Friends_Management/internal/domain/dto"
	service "BE_Friends_Management/internal/service/friendship"
	"BE_Friends_Management/pkg"
	"BE_Friends_Management/pkg/utils"
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

// User godoc
// @Summary      Create new friendship
// @Description  Create new friendship
// @Tags         Friendship
// @Accept 		json
// @Produce      json
// @Param 		 request body dto.CreateFriendshipRequest true "List of 2 friend emails"  example({"friends": ["andy@example.com", "john@example.com"]})
// @Router       /api/friendship [POST]
// @Success      200   {object}  dto.ApiResponseSuccessNoData
func (h *FriendshipHandler) CreateFriendship(c *gin.Context) {
	defer pkg.PanicHandler(c)
	var request dto.CreateFriendshipRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Error("Happened error when mapping request. Error: ", err)
		pkg.PanicExeption(constant.InvalidRequest, "Invalid request format.")
	}
	err := h.service.CreateFriendship(request.Friends[0], request.Friends[1])
	if err != nil {
		log.Error("Happened error when creating new friendship. Error: ", err)
		switch err {
		case service.ErrInvalidRequest:
			pkg.PanicExeption(constant.InvalidRequest, err.Error())
		case service.ErrUserNotFound:
			pkg.PanicExeption(constant.DataNotFound, err.Error())
		case service.ErrAlreadyFriend:
			pkg.PanicExeption(constant.Conflict, err.Error())
		default:
			pkg.PanicExeption(constant.UnknownError, "Happened error when creating new friendship.")
		}
	}
	c.JSON(http.StatusOK, pkg.BuildReponseSuccessNoData())
}

// User godoc
// @Summary      Retrieve friends list for an email address
// @Description  Retrieve friends list for an email address
// @Tags         Friendship
// @Accept 		json
// @Produce      json
// @Param 		 email query string true "Email address"
// @Router       /api/friendship/friends [GET]
// @Success      200   {object}  dto.ApiResponseSuccessWithFriendsList
func (h *FriendshipHandler) RetrieveFriendsList(c *gin.Context) {
	defer pkg.PanicHandler(c)
	requestEmail := c.Query("email")
	if requestEmail == "" {
		log.Error("Happened error when mapping request. Error: received no email input.")
		pkg.PanicExeption(constant.InvalidRequest, "Invalid request format.")
		return
	}
	friends, err := h.service.RetrieveFriendsList(requestEmail)
	if err != nil {
		log.Error("Happened error when retrieving friends list. Error: ", err)
		switch err {
		case service.ErrUserNotFound:
			pkg.PanicExeption(constant.DataNotFound, err.Error())
		default:
			pkg.PanicExeption(constant.UnknownError, "Happened error when retrieving friends list.")
		}
	}
	emails := utils.ConvertUsersToEmails(friends)
	count := h.service.CountFriends(friends)
	c.JSON(http.StatusOK, pkg.BuildReponseSuccessWithFriendsList(emails, count))
}

// User godoc
// @Summary      Retrieve common friends list between two email addresses
// @Description  Retrieve common friends list between two email addresses
// @Tags         Friendship
// @Accept 		json
// @Produce      json
// @Param 		 email1 query string true "Email address of user 1"
// @Param 		 email2 query string true "Email address of user 2"
// @Router       /api/friendship/common-friends [GET]
// @Success      200   {object}  dto.ApiResponseSuccessWithFriendsList
func (h *FriendshipHandler) RetrieveCommonFriends(c *gin.Context) {
	defer pkg.PanicHandler(c)
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
	friends, err := h.service.RetrieveCommonFriends(requestEmail1, requestEmail2)
	if err != nil {
		log.Error("Happened error when retrieving common friends list. Error: ", err)
		switch err {
		case service.ErrInvalidRequest:
			pkg.PanicExeption(constant.InvalidRequest, err.Error())
		case service.ErrUserNotFound:
			pkg.PanicExeption(constant.DataNotFound, err.Error())
		default:
			pkg.PanicExeption(constant.UnknownError, "Happened error when retrieving common friends list.")
		}
	}
	emails := utils.ConvertUsersToEmails(friends)
	count := h.service.CountFriends(friends)
	c.JSON(http.StatusOK, pkg.BuildReponseSuccessWithFriendsList(emails, count))
}
