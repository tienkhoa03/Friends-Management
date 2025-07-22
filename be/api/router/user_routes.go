package api

import (
	"BE_Friends_Management/api/handler"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func registerUserRoutes(api *gin.RouterGroup, h *handler.UserHandler, db *gorm.DB) {
	api.GET("/users", h.GetAllUser)
	api.GET("/users/:id", h.GetUserById)
	api.POST("/users", h.CreateUser)
	api.DELETE("/users/:id", h.DeleteUserById)
	api.PUT("users/:id", h.UpdateUser)
}
