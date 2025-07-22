package repository

import (
	"BE_Friends_Management/internal/domain/entity"

	"gorm.io/gorm"
)

type PostgreSQLUserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &PostgreSQLUserRepository{db: db}
}

func (r *PostgreSQLUserRepository) GetAllUser() []*entity.User {
	var users = []*entity.User{}
	result := r.db.Model(entity.User{}).Find(&users)
	if result.Error != nil {
		return nil
	}
	return users
}

func (r *PostgreSQLUserRepository) GetUserById(userId int64) (*entity.User, error) {
	var user = entity.User{}
	result := r.db.Model(entity.User{}).Where("id = ?", userId).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *PostgreSQLUserRepository) GetUserByEmail(email string) (*entity.User, error) {
	var user = entity.User{}
	result := r.db.Model(entity.User{}).Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *PostgreSQLUserRepository) CreateUser(user *entity.User) (*entity.User, error) {
	result := r.db.Create(user)
	return user, result.Error
}

func (r *PostgreSQLUserRepository) DeleteUserById(userId int64) error {
	result := r.db.Model(entity.User{}).Where("id = ?", userId).Delete(entity.User{})
	return result.Error
}

func (r *PostgreSQLUserRepository) UpdateUser(user *entity.User) (*entity.User, error) {
	result := r.db.Model(entity.User{}).Where("id = ?", user.Id).Updates(user)
	if result.Error != nil {
		return nil, result.Error
	}
	var updatedUser = entity.User{}
	result = r.db.First(&updatedUser, user.Id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &updatedUser, nil

}
