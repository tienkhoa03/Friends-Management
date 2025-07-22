package api

import (
	"BE_Friends_Management/api/handler"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, userHandler *handler.UserHandler, db *gorm.DB) {
	api := r.Group("/api")
	registerUserRoutes(api, userHandler, db)
}
