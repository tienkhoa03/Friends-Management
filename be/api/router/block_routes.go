package api

import (
	"BE_Friends_Management/api/handler"
	"BE_Friends_Management/api/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func registerBlockRoutes(api *gin.RouterGroup, h *handler.BlockRelationshipHandler, db *gorm.DB) {
	api.Use(middleware.ValidateAccessToken())
	api.POST("/block", h.CreateBlockRelationship)
}
