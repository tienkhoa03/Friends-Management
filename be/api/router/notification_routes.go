package api

import (
	"BE_Friends_Management/api/handler"
	"BE_Friends_Management/api/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func registerNotificationRoutes(api *gin.RouterGroup, h *handler.NotificationHandler, db *gorm.DB) {
	api.Use(middleware.ValidateAccessToken())
	api.POST("/update-recipients", middleware.RequireAnyRole([]string{"admin", "user"}), h.GetUpdateRecipients)
}
