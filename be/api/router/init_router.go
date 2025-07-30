package api

import (
	"BE_Friends_Management/api/handler"
	"BE_Friends_Management/api/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, handlers *handler.Handlers, db *gorm.DB) {
	r.Use(middleware.CORSMiddleware())
	api := r.Group("/api")
	registerUserRoutes(api, handlers.User, db)
	registerFriendshipRoutes(api, handlers.Friendship, db)
	registerSubscriptionRoutes(api, handlers.Subscription, db)
	registerBlockRoutes(api, handlers.BlockRelationship, db)
	registerNotificationRoutes(api, handlers.NotificationHandler, db)
}
