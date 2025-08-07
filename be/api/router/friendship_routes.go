package api

import (
	"BE_Friends_Management/api/handler"
	"BE_Friends_Management/api/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func registerFriendshipRoutes(api *gin.RouterGroup, h *handler.FriendshipHandler, db *gorm.DB) {
	api.Use(middleware.ValidateAccessToken())
	api.POST("/friendship", h.CreateFriendship)
	api.GET("/friendship/friends", h.RetrieveFriendsList)
	api.GET("/friendship/common-friends", h.RetrieveCommonFriends)
}
