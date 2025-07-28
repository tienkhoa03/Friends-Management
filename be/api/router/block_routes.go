package api

import (
	"BE_Friends_Management/api/handler"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func registerBlockRoutes(api *gin.RouterGroup, h *handler.BlockRelationshipHandler, db *gorm.DB) {
	api.POST("/block", h.CreateBlockRelationship)
}
