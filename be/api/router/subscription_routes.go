package api

import (
	"BE_Friends_Management/api/handler"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func registerSubscriptionRoutes(api *gin.RouterGroup, h *handler.SubscriptionHandler, db *gorm.DB) {
	api.POST("/subscription", h.CreateSubscription)
}
