package handler

import (
	"BE_Friends_Management/constant"
	"BE_Friends_Management/internal/domain/dto"
	service "BE_Friends_Management/internal/service/subscription"
	"BE_Friends_Management/pkg"
	"BE_Friends_Management/pkg/utils"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type SubscriptionHandler struct {
	service service.SubscriptionService
}

func NewSubscriptionHandler(service service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{service: service}
}

// Subscription godoc
// @Summary      Create new subscription
// @Description  Create new subscription
// @Tags         Subscription
// @Accept 		json
// @Produce      json
// @Param 		 request body dto.CreateSubscriptionRequest true "Requestor's email and target's email"
// @param Authorization header string true "Authorization"
// @Router       /api/subscription [POST]
// @Success      200   {object}  dto.ApiResponseSuccessNoData
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func (h *SubscriptionHandler) CreateSubscription(c *gin.Context) {
	defer pkg.PanicHandler(c)
	authUserId := utils.GetAuthUserID(c)
	var request dto.CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Error("Happened error when mapping request. Error: ", err)
		pkg.PanicExeption(constant.InvalidRequest, "Invalid request format.")
	}
	err := h.service.CreateSubscription(authUserId, request.Requestor, request.Target)
	if err != nil {
		log.Error("Happened error when creating new subscription. Error: ", err)
		switch {
		case errors.Is(err, service.ErrInvalidRequest):
			pkg.PanicExeption(constant.InvalidRequest, err.Error())
		case errors.Is(err, service.ErrUserNotFound):
			pkg.PanicExeption(constant.DataNotFound, err.Error())
		case errors.Is(err, service.ErrAlreadySubscribed):
			pkg.PanicExeption(constant.Conflict, err.Error())
		case errors.Is(err, service.ErrIsBlocked):
			pkg.PanicExeption(constant.StatusForbidden, err.Error())
		case errors.Is(err, service.ErrNotPermitted):
			pkg.PanicExeption(constant.StatusForbidden, err.Error())
		default:
			pkg.PanicExeption(constant.UnknownError, "Happened error when creating new subscription.")
		}
	}
	c.JSON(http.StatusOK, pkg.BuildResponseSuccessNoData())
}
