package services

import (
	"fmt"
	"log/slog"

	"messenger-auth/internal/dto"
	"messenger-auth/internal/models"
	"messenger-auth/internal/repositories"
)

type UserServiceGORM struct {
	Repo repositories.UserRepository
	Log  *slog.Logger
}

func NewUserServiceGORM(repo repositories.UserRepository, logger *slog.Logger) UserService {
	return &UserServiceGORM{Repo: repo, Log: logger}
}

func (userService *UserServiceGORM) CreateUser(userDTO *dto.UserData) error {
	var user models.User
	// Маппинг данных из DTO в модель
	if err := userDTO.Map(&user); err != nil {
		userService.Log.Error(fmt.Sprintf("Failed to create user. Error: %s", err.Error()), slog.Any("user_data", user))
		return err
	}
	if err := userService.Repo.Create(&user); err != nil {
		userService.Log.Error(fmt.Sprintf("Failed to create user. Error: %s", err.Error()), slog.Any("user_data", user))
		return err
	}
	userService.Log.Debug("User was created", slog.Any("user_data", user))
	return nil
}

func (userService *UserServiceGORM) UpdateUser(id uint, userDTO *dto.UserData) error {
	user, err := userService.Repo.FindByID(id)
	if err != nil {
		userService.Log.Error(fmt.Sprintf("Failed to update user. Error: %s", err.Error()), slog.Any("user_data", user))
		return err
	}
	// Маппинг данных из DTO в модель
	if err = userDTO.Map(user); err != nil {
		userService.Log.Error(fmt.Sprintf("Failed to update user. Error: %s", err.Error()), slog.Any("user_data", user))
		return err
	}
	if err = userService.Repo.Update(user); err != nil {
		userService.Log.Error(fmt.Sprintf("Failed to update user. Error: %s", err.Error()), slog.Any("user_data", user))
		return err
	}
	userService.Log.Debug("User was updated", slog.Any("user_data", user))
	return nil
}

func (userService *UserServiceGORM) DeleteUser(id uint) error {
	if err := userService.Repo.Delete(id); err != nil {
		userService.Log.Error(fmt.Sprintf("Failed to delete user. Error: %s", err.Error()), slog.Any("user_id", id))
		return err
	}
	userService.Log.Debug("User was deleted", slog.Any("user_id", id))
	return nil
}

func (userService *UserServiceGORM) GetUserByID(id uint) (*models.User, error) {
	user, err := userService.Repo.FindByID(id)
	if err != nil {
		userService.Log.Error(fmt.Sprintf("Failed to get user. Error: %s", err.Error()), slog.Any("user_id", id))
	}
	userService.Log.Debug("User was received from DB", slog.Uint64("user_id", uint64(id)))
	return user, err
}

func (userService *UserServiceGORM) GetUsers() ([]models.User, error) {
	users, err := userService.Repo.FindAll()
	if err != nil {
		userService.Log.Error(fmt.Sprintf("Failed to get all users. Error: %s", err.Error()))
	}
	userService.Log.Debug("Users was received from DB")
	return users, err
}

func (s *UserServiceGORM) GetRepo() repositories.UserRepository {
	return s.Repo
}
