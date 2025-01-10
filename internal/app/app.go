package app

import (
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/redis/go-redis/v9"
	"messenger-auth/config"
	"messenger-auth/internal/repositories"
	"messenger-auth/internal/services"
)

type App struct {
	Log           *slog.Logger
	Router        *gin.Engine
	Database      *gorm.DB
	RedisDatabase *redis.Client

	UserRepo repositories.UserRepository

	UserService services.UserService

	Config *config.Config
}

func NewApp(log *slog.Logger, config *config.Config) *App {
	return &App{
		Log:    log,
		Config: config,
	}
}

func (a *App) Initialize() {
	a.setupRouter()
	a.setupDatabase()
	a.setupRepositories()
	a.setupServices()
	a.setupHandlersAndRoutes()
}
func (a *App) Run() {
	if err := a.Router.Run(a.Config.ServiceAddress); err != nil {
		a.Log.Error("Failed to run user service", slog.String("error", err.Error()))
	}
}

func (a *App) setupRouter() {
	r := gin.Default()
	r.Use(gin.Recovery())
	a.Router = r
}

func (a *App) setupDatabase() {
	db, err := repositories.InitPostgres(a.Config)
	if err != nil {
		a.Log.Error("Failed to connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	a.Database = db

	redisDB, err := repositories.InitRedis(a.Config)
	if err != nil {
		a.Log.Error("Failed to connect to Redis database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	a.RedisDatabase = redisDB
}

func (a *App) setupRepositories() {
	repoLogger := a.Log.With(slog.String("service", "user"), slog.String("module", "repository"))
	if a.Config.RedisEnabled == "true" {
		a.UserRepo = repositories.NewUserRepoWithRedis(
			&repositories.UserRepoPostgres{DB: a.Database, Log: repoLogger},
			a.RedisDatabase,
			repoLogger,
		)
	} else {
		a.UserRepo = repositories.NewUserRepoPostgres(
			a.Database,
			repoLogger,
		)
	}
	
}

func (a *App) setupServices() {
	a.UserService = services.NewUserServiceGORM(a.UserRepo, a.Log.With(slog.String("service", "user"), slog.String("module", "service")))
}

func (a *App) setupHandlersAndRoutes() {
	SetupHandlers(a.Router, a.UserService, a.Log.With(slog.String("service", "user"), slog.String("module", "transport")))
}
