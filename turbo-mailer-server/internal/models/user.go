package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"unique;type:varchar(255)"`
	Password string `gorm:"not null;type:varchar(512)"`
}
