package api

import (
	"BE_Friends_Management/api/handler"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func registerNotificationRoutes(api *gin.RouterGroup, h *handler.NotificationHandler, db *gorm.DB) {
	api.POST("/update-recipients", h.GetUpdateRecipients)
}
