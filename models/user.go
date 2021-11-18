package models

type User struct {
	Username string `gorm:"size:100;not null;unique"  json:"username" binding:"required,gte=0"`
	Password string `gorm:"not null;" json:"password"`
}
