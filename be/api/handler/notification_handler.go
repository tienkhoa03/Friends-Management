package handler

import (
	"BE_Friends_Management/constant"
	"BE_Friends_Management/internal/domain/dto"
	service "BE_Friends_Management/internal/service/notification"
	"BE_Friends_Management/pkg"
	"BE_Friends_Management/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type NotificationHandler struct {
	service service.NotificationService
}

func NewNotificationHandler(service service.NotificationService) *NotificationHandler {
	return &NotificationHandler{service: service}
}

// User godoc
// @Summary      Get update recipients
// @Description  Get all email addresses that can receive updates from an email address.
// @Tags         Notification
// @Accept 		json
// @Produce      json
// @Param 		 request body dto.GetUpdateRecipientsRequest true "Sender email and update text"
// @Router       /api/update-recipients [POST]
// @Success      200   {object}  dto.ApiResponseSuccessNoData
func (h *NotificationHandler) GetUpdateRecipients(c *gin.Context) {
	defer pkg.PanicHandler(c)
	var request dto.GetUpdateRecipientsRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Error("Happened error when mapping request. Error: ", err)
		pkg.PanicExeption(constant.InvalidRequest, "Invalid request format.")
	}
	recipients, err := h.service.GetUpdateRecipients(request.Sender, request.Text)
	if err != nil {
		log.Error("Happened error when getting recipients. Error: ", err)
		switch err {
		case service.ErrUserNotFound:
			pkg.PanicExeption(constant.DataNotFound, err.Error())
		default:
			pkg.PanicExeption(constant.UnknownError, "Happened error when getting recipients.")
		}
	}
	recipientEmails := utils.ConvertUsersToEmails(recipients)
	c.JSON(http.StatusOK, pkg.BuildResponseSuccessWithRecipients(recipientEmails))
}
