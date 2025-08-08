package api

import (
	"BE_Friends_Management/api/handler"
	"BE_Friends_Management/api/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func registerUserRoutes(api *gin.RouterGroup, h *handler.UserHandler, db *gorm.DB) {
	api.Use(middleware.ValidateAccessToken())
	api.GET("/users", middleware.RequireAnyRole([]string{"admin"}), h.GetAllUser)
	api.GET("/users/:id", middleware.RequireAnyRole([]string{"admin"}), h.GetUserById)
	api.DELETE("/users/:id", middleware.RequireAnyRole([]string{"admin"}), h.DeleteUserById)
	api.PUT("users/:id", middleware.RequireAnyRole([]string{"admin"}), h.UpdateUser)
}
