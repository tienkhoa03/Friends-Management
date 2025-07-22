package repository

import (
	user "BE_Friends_Management/internal/repository/users"

	"gorm.io/gorm"
)

type Repository struct {
	User user.UserRepository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		User: user.NewUserRepository(db),
	}
}
