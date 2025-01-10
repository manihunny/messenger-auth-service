package services

import (
	"messenger-auth/internal/dto"
	"messenger-auth/internal/models"
	"messenger-auth/internal/repositories"
)

type UserService interface {
	CreateUser(userDTO *dto.UserData) error
	UpdateUser(id uint, userDTO *dto.UserData) error
	DeleteUser(id uint) error
	GetUserByID(id uint) (*models.User, error)
	GetUsers() ([]models.User, error)
	GetRepo() repositories.UserRepository
}
