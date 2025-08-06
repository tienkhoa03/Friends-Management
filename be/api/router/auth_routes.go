package api

import (
	"BE_Friends_Management/api/handler"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func registerAuthRoutes(api *gin.RouterGroup, h *handler.AuthHandler, db *gorm.DB) {
	api.POST("/auth/register", h.RegisterUser)
	api.POST("/auth/login", h.Login)
	api.POST("/auth/refresh", h.RefreshAccessToken)
}
