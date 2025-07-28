package entity

import "time"

type Friendship struct {
	UserId1   int64 `gorm:"primaryKey" json:"user_id_1"`
	UserId2   int64 `gorm:"primaryKey" json:"user_id_2"`
	CreatedAt time.Time

	User1 *User `gorm:"foreignKey:UserId1;references:Id"`
	User2 *User `gorm:"foreignKey:UserId2;references:Id"`
}
