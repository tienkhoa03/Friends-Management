package entity

import "time"

type Subscription struct {
	RequestorId int64 `gorm:"primaryKey" json:"requestor_id"`
	TargetId    int64 `gorm:"primaryKey" json:"target_id"`
	CreatedAt   time.Time

	Requestor *User `gorm:"foreignKey:RequestorId;references:Id"`
	Target    *User `gorm:"foreignKey:TargetId;references:Id"`
}
