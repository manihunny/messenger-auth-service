package dto

import (
	"messenger-auth/internal/models"
	"messenger-auth/internal/utils"
)

// UserData представляет данные проекта, которые может установить пользователь
type UserData struct {
	Name     *string `json:"name" gorm:"type:varchar(255);not null"` // Имя (никнейм) пользователя
	Email    *string `json:"email" gorm:"type:varchar(255);not null"` // Электронная почта пользователя
}

// Parse достает данные из модели и вставляет их в UserData
func (dto *UserData) Parse(user *models.User) error {
	utils.CopyIfNotNil(&dto.Name, user.Name)
	utils.CopyIfNotNil(&dto.Email, user.Email)
	
	return nil
}

// Map обновляет данные модели значениями из UserData
func (dto *UserData) Map(user *models.User) error {
	utils.CopyIfNotNil(&dto.Name, user.Name)
	utils.CopyIfNotNil(&dto.Email, user.Email)

	return nil
}
