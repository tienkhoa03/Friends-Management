package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	Port                         string
	AccessSecret                 string
	RefreshSecret                string
	PasswordSecret               string
	SmtpPasswd                   string
	BASE_URL_BACKEND             string
	DB_DNS                       string
	BASE_URL_BACKEND_FOR_SWAGGER string
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, continuing with environment variables")
	}
	Port = ":" + os.Getenv("PORT")
	AccessSecret = os.Getenv("AccessSecret")
	RefreshSecret = os.Getenv("refreshSecret")
	PasswordSecret = os.Getenv("PasswordSecret")
	SmtpPasswd = os.Getenv("SMTP_PASSWORD")
	BASE_URL_BACKEND_FOR_SWAGGER = os.Getenv("BASE_URL_BACKEND_FOR_SWAGGER")
	BASE_URL_BACKEND = os.Getenv("BASE_URL_BACKEND")
	DB_DNS = os.Getenv("DATABASE_URL")
}
