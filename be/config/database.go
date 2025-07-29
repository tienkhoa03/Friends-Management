package config

import (
	"log"

	entity "BE_Friends_Management/internal/domain/entity"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var users = []entity.User{
	{Email: "user1@gmail.com"},
	{Email: "user2@gmail.com"},
	{Email: "user3@gmail.com"},
	{Email: "user4@gmail.com"},
	{Email: "user5@gmail.com"},
	{Email: "user6@gmail.com"},
	{Email: "user7@gmail.com"},
	{Email: "user8@gmail.com"},
	{Email: "user9@gmail.com"},
	{Email: "user10@gmail.com"},
}

func ConnectToDB() *gorm.DB {
	db, err := gorm.Open(postgres.Open(DB_DNS), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database. Error:", err)
	}
	err = db.AutoMigrate(&entity.User{}, &entity.Friendship{}, &entity.Subscription{}, &entity.BlockRelationship{})
	if err != nil {
		log.Fatal("Error migrate to database. Error:", err)
	}
	for _, user := range users {
		var existing entity.User
		db.Where("email = ?", user.Email).FirstOrCreate(&existing, user)
	}
	return db
}
