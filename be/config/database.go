package config

import (
	"log"

	entity "BE_Friends_Management/internal/domain/entity"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var users = []entity.User{
	{
		Email:    "user1@gmail.com",
		Password: "$2a$10$uD2Sp/ceVMQs.Fxa9883Lejcy4QSiEsWFIihuosOkCqwQaCrs011.",
	},
	{
		Email:    "user2@gmail.com",
		Password: "$2a$10$Rkga1eAiQ4xSFSfIA.ZFyuraVz8lAE7/d.OsrVHb8Cd2J/KoVnkWu",
	},
	{
		Email:    "user3@gmail.com",
		Password: "$2a$10$AGvvpScnwlpreNybde2RYOu3YwXWR5upqH4CYgY4kyrR9IUOS/2SC",
	},
	{
		Email:    "user4@gmail.com",
		Password: "$2a$10$gPgRynYgAnJga.yDxY/E7OcjJFMFv4fsB3lL4lvnsvmpigYNMNJ2W",
	},
	{
		Email:    "user5@gmail.com",
		Password: "$2a$10$Uu4bpMgDh5BqgCoxNNMD6ePiPXYJHOdCmDGf9JO7LflS6rxVo29t6",
	},
	{
		Email:    "user6@gmail.com",
		Password: "$2a$10$HpKZlAE1EgXm2qSVUzDNY.Jl21nJdJoJF9N8Eo2h07WrFpKgd3hE6",
	},
	{
		Email:    "user7@gmail.com",
		Password: "$2a$10$FbSLfcYefGmoqUFZWxIF2.TPb3ujjSsHCKhSYMP86VpEYozx6JCr6",
	},
	{
		Email:    "user8@gmail.com",
		Password: "$2a$10$wVSYY0LmSYRXbEO3JyRWMu.JnNk.tsjCJgAMMSuWkm58eNMe2XmdW",
	},
	{
		Email:    "user9@gmail.com",
		Password: "$2a$10$5t/A3R/jOLUxA2EFuCS/oeZA27i2YZ4PLBQAQ8/CK456dYUpMRrCa",
	},
	{
		Email:    "user10@gmail.com",
		Password: "$2a$10$uD2Sp/ceVMQs.Fxa9883Lejcy4QSiEsWFIihuosOkCqwQaCrs011.",
	},
}

func ConnectToDB() *gorm.DB {
	db, err := gorm.Open(postgres.Open(DB_DNS), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database. Error:", err)
	}
	err = db.AutoMigrate(&entity.User{}, &entity.Friendship{}, &entity.Subscription{}, &entity.BlockRelationship{}, &entity.UserToken{})
	if err != nil {
		log.Fatal("Error migrate to database. Error:", err)
	}
	for _, user := range users {
		var existing entity.User
		db.Where("email = ?", user.Email).FirstOrCreate(&existing, user)
	}
	return db
}
