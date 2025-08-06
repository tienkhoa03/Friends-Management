package entity

import "time"

type UserToken struct {
	Id           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId       int64     `gorm:"type:varchar(256);not null" json:"email"`
	RefreshToken string    `gorm:"type:varchar(256);not null" json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	ExpiresAt    time.Time `json:"expires_at"`
	IsRevoked    bool      `json:"is_revoked"`

	User *User `gorm:"foreignKey:UserId;references:Id"`
}
