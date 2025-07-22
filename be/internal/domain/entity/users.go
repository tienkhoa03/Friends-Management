package entity

type User struct {
	Id    int64  `gorm: "primaryKey;autoIncrement" json:"id"`
	Email string `gorm:"type:varchar(256);not null;unique" json:"email"`
}
