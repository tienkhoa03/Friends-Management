package handler

import (
	"BE_Friends_Management/constant"
	"BE_Friends_Management/internal/domain/dto"
	service "BE_Friends_Management/internal/service/friendship"
	"BE_Friends_Management/pkg"
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
	var friendship dto.CreateFriendshipRequest
	if err := c.ShouldBindJSON(&friendship); err != nil {
		log.Error("Happened error when mapping request from FE. Error: ", err)
		pkg.PanicExeption(constant.InvalidRequest, "Invalid request format.")
	}
	err := h.service.CreateFriendship(friendship.Friends[0], friendship.Friends[1])
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
