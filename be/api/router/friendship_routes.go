package api

import (
	"BE_Friends_Management/api/handler"
	"BE_Friends_Management/api/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func registerFriendshipRoutes(api *gin.RouterGroup, h *handler.FriendshipHandler, db *gorm.DB) {
	api.Use(middleware.ValidateAccessToken())
	api.POST("/friendship", middleware.RequireAnyRole([]string{"user"}), h.CreateFriendship)
	api.GET("/friendship/friends", middleware.RequireAnyRole([]string{"admin", "user"}), h.RetrieveFriendsList)
	api.GET("/friendship/common-friends", middleware.RequireAnyRole([]string{"admin", "user"}), h.RetrieveCommonFriends)
}
