package main

import (
	"log"

	"BE_Friends_Management/api/handler"
	api "BE_Friends_Management/api/router"
	"BE_Friends_Management/config"
	"BE_Friends_Management/internal/repository"
	"BE_Friends_Management/internal/service"

	"BE_Friends_Management/cmd/server/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	config.LoadEnv()
	db := config.ConnectToDB()
	repos := repository.NewRepository(db)

	services := service.NewService(repos)
	handlers := handler.NewHandlers(services)
	docs.SwaggerInfo.Title = "API Friends Management"
	docs.SwaggerInfo.Description = "API Friends Management"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = config.BASE_URL_BACKEND_FOR_SWAGGER
	docs.SwaggerInfo.BasePath = ""
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	r := gin.Default()
	api.SetupRoutes(r, handlers, db)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := r.Run(config.Port); err != nil {
		log.Fatal("failed to run server:", err)
	}
}
