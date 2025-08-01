package entity

import "time"

type User struct {
	Id        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Email     string    `gorm:"type:varchar(256);not null;unique" json:"email"`
	CreatedAt time.Time `json:"created_at"`
}
