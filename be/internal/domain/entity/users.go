package entity

import "time"

type User struct {
	Id        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Email     string    `gorm:"type:varchar(256);not null;unique" json:"email"`
	Password  string    `gorm:"type:varchar(256);not null" json:"-"`
	Role      string    `gorm:"type:role_slug" json:"role"`
	CreatedAt time.Time `json:"created_at"`
}
