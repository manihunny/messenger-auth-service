package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Name     *string `json:"name" gorm:"type:varchar(255);not null"` // Имя (никнейм) пользователя
	Email    *string `json:"email" gorm:"type:varchar(255);not null"` // Электронная почта пользователя
	Password *string `json:"password" gorm:"type:varchar(255);not null"` // Пароль пользователя
}

func (User) TableName() string {
	return "users"
}
