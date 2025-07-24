package api

import (
	"BE_Friends_Management/api/handler"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func registerFriendshipRoutes(api *gin.RouterGroup, h *handler.FriendshipHandler, db *gorm.DB) {
	api.POST("/friendship", h.CreateFriendship)
	api.GET("/friendship/friends", h.RetrieveFriendsList)
}
