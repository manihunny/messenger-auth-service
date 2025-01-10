package repositories

import (
	"github.com/jinzhu/gorm"
	"messenger-auth/internal/models"
)

type UserRepository interface {
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(id uint) error
	FindByID(id uint) (*models.User, error)
	FindAll() ([]models.User, error)
	GetDB() *gorm.DB
}
