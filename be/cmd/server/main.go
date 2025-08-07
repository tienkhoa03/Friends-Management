package main

import (
	"log"

	"BE_Friends_Management/api/handler"
	api "BE_Friends_Management/api/router"
	"BE_Friends_Management/cmd/server/docs"
	"BE_Friends_Management/config"
	"BE_Friends_Management/internal/repository"
	"BE_Friends_Management/internal/service"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Friends Management API
// @version         1.0
// @description     Friends Management API
// @BasePath
// @schemes         http https

func main() {
	config.LoadEnv()
	db := config.ConnectToDB()
	repos := repository.NewRepository(db)

	services := service.NewService(repos)
	handlers := handler.NewHandlers(services)

	r := gin.Default()
	api.SetupRoutes(r, handlers, db)
	docs.SwaggerInfo.Host = config.BASE_URL_BACKEND_FOR_SWAGGER
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := r.Run(config.Port); err != nil {
		log.Fatal("failed to run server:", err)
	}
}
