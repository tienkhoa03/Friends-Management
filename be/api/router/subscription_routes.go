package api

import (
	"BE_Friends_Management/api/handler"
	"BE_Friends_Management/api/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func registerSubscriptionRoutes(api *gin.RouterGroup, h *handler.SubscriptionHandler, db *gorm.DB) {
	api.Use(middleware.ValidateAccessToken())
	api.POST("/subscription", middleware.RequireAnyRole([]string{"user"}), h.CreateSubscription)
}
