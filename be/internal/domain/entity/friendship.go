package entity

type Friendship struct {
	UserId1 int64 `gorm:"primaryKey" json:"user_id_1"`
	UserId2 int64 `gorm:"primaryKey" json:"user_id_2"`
}
