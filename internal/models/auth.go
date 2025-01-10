package models

import "time"

type RefreshToken struct {
	ID        uint      `json:"id" gorm:"primaryKey"`         // Уникальный идентификатор токена
	Token     string    `json:"token" gorm:"unique;not null"` // Токен (уникальный и обязательный)
	UserID    uint      `json:"userId" gorm:"not null"`       // Идентификатор пользователя
	ExpiresAt time.Time `json:"expiresAt" gorm:"not null"`    // Время истечения токена
}
