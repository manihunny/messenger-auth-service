package app

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"messenger-auth/internal/controllers"
	"messenger-auth/internal/services"
)

func SetupHandlers(r *gin.Engine, userService services.UserService, log *slog.Logger) {
	userHandler := &controllers.UserHandler{Service: userService, Log: log}

	// Оповещение docker-compose о том, что контейнер готов к работе
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "OK",
		})
	})
	
	v1 := r.Group("/user/api/v1")
	{
		v1.POST("/", userHandler.CreateUser)
		v1.PUT("/:id", userHandler.UpdateUser)
		v1.DELETE("/:id", userHandler.DeleteUser)
		v1.GET("/:id", userHandler.GetUserByID)
		v1.GET("/", userHandler.GetUsers)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
