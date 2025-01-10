package repositories

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/redis/go-redis/v9"
	"messenger-auth/internal/models"
)

// Структура для работы с Postgres
type UserRepoPostgres struct {
	DB  *gorm.DB
	Log *slog.Logger
}

func NewUserRepoPostgres(db *gorm.DB, logger *slog.Logger) UserRepository {
	return &UserRepoPostgres{DB: db, Log: logger}
}

func (repo *UserRepoPostgres) Create(user *models.User) error {
	if err := repo.DB.Create(user).Error; err != nil {
		repo.Log.Error(fmt.Sprintf("Failed to create user in DB. Error: %s", err.Error()), slog.Any("user_data", user))
		return err
	}
	repo.Log.Debug("User was created in DB", slog.Any("user_data", user))
	return nil
}

func (repo *UserRepoPostgres) Update(user *models.User) error {
	if err := repo.DB.Save(user).Error; err != nil {
		repo.Log.Error(fmt.Sprintf("Failed to update user in DB. Error: %s", err.Error()), slog.Any("user_data", user))
		return err
	}
	repo.Log.Debug("User was updated in DB", slog.Any("user_data", user))
	return nil
}

func (repo *UserRepoPostgres) Delete(id uint) error {
	if err := repo.DB.Delete(&models.User{}, id).Error; err != nil {
		repo.Log.Error(fmt.Sprintf("Failed to delete user from DB. Error: %s", err.Error()), slog.Uint64("user_id", uint64(id)))
		return err
	}
	repo.Log.Debug("User was deleted from DB", slog.Uint64("user_id", uint64(id)))
	return nil
}

func (repo *UserRepoPostgres) FindByID(id uint) (*models.User, error) {
	var user models.User
	if err := repo.DB.First(&user, id).Error; err != nil {
		repo.Log.Error(fmt.Sprintf("Failed to get user. Error: %s", err.Error()), slog.Uint64("user_id", uint64(id)))
	}
	repo.Log.Debug("User was received from DB", slog.Uint64("user_id", uint64(id)))
	return &user, nil
}

func (repo *UserRepoPostgres) FindAll() ([]models.User, error) {
	var users []models.User
	if err := repo.DB.Find(&users).Error; err != nil {
		repo.Log.Error(fmt.Sprintf("Failed to get users. Error: %s", err.Error()))
	}
	repo.Log.Debug("Users was received from DB")
	return users, nil
}

// GetDB дает доступ к полю DB
func (repo *UserRepoPostgres) GetDB() *gorm.DB {
	return repo.DB
}

// Структура для работы с Redis (надстройка для Postgres)
type UserRepoRedis struct {
	DBRepo  *UserRepoPostgres
	RedisDB *redis.Client
	Log     *slog.Logger
}

func NewUserRepoWithRedis(dbRepo *UserRepoPostgres, redisDB *redis.Client, logger *slog.Logger) UserRepository {
	return &UserRepoRedis{DBRepo: dbRepo, RedisDB: redisDB, Log: logger}
}

func (repo *UserRepoRedis) Create(user *models.User) error {
	// Create record in DB
	err := repo.DBRepo.Create(user)
	if err != nil {
		return err
	}

	// Create record in Redis
	if err := setWithMarshal(repo.RedisDB, "user_"+strconv.FormatUint(uint64(user.ID), 10), user); err != nil {
		repo.Log.Warn("Couldn't create record in Redis", slog.String("error", err.Error()))
	}
	return nil
}

func (repo *UserRepoRedis) Update(user *models.User) error {
	// Update record in DB
	err := repo.DBRepo.Update(user)
	if err != nil {
		return err
	}

	// Update record in Redis
	if err := setWithMarshal(repo.RedisDB, "user_"+strconv.FormatUint(uint64(user.ID), 10), user); err != nil {
		repo.Log.Warn("Couldn't update record in Redis", slog.String("error", err.Error()))
	}
	return nil
}

func (repo *UserRepoRedis) Delete(id uint) error {
	// Delete record from DB
	err := repo.DBRepo.Delete(id)
	if err != nil {
		return err
	}

	// Delete record from Redis
	ctx := context.Background()
	if err := repo.RedisDB.Del(ctx, "user_"+strconv.FormatUint(uint64(id), 10)).Err(); err != nil {
		repo.Log.Warn("Couldn't delete record in Redis", slog.String("error", err.Error()))
	}
	return nil
}

func (repo *UserRepoRedis) FindByID(id uint) (*models.User, error) {
	var user models.User

	// Try get record from Redis, otherwise get from DB and create record in Redis
	err := getWithUnmarshal(repo.RedisDB, "user_"+strconv.FormatUint(uint64(id), 10), &user)
	if err != nil {
		user, err := repo.DBRepo.FindByID(id)
		if err != nil {
			return user, err
		}

		err = setWithMarshal(repo.RedisDB, "user_"+strconv.FormatUint(uint64(user.ID), 10), user)
		if err != nil {
			repo.Log.Warn("Couldn't create record in Redis", slog.String("error", err.Error()))
		}
	}
	return &user, err
}

func (repo *UserRepoRedis) FindAll() ([]models.User, error) {
	return repo.DBRepo.FindAll()
}

// GetDB дает доступ к полю DB
func (repo *UserRepoRedis) GetDB() *gorm.DB {
	return repo.DBRepo.DB
}
