package controllers

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"messenger-auth/internal/dto"
	"messenger-auth/internal/services"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Service services.UserService
	Log     *slog.Logger
}

func NewUserHandler(userService services.UserService, log *slog.Logger) *UserHandler {
	return &UserHandler{Service: userService, Log: log}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var userDTO dto.UserData
	if err := c.ShouldBindJSON(&userDTO); err != nil {
		h.Log.Error(fmt.Sprintf("Failed to create user. Error: %s", err.Error()), slog.String("method", c.Request.Method), slog.Int("code", http.StatusBadRequest), slog.String("url", c.Request.URL.Path), slog.String("client", c.ClientIP()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Service.CreateUser(&userDTO); err != nil {
		h.Log.Error(fmt.Sprintf("Failed to create user. Error: %s", err.Error()), slog.String("method", c.Request.Method), slog.Int("code", http.StatusInternalServerError), slog.String("url", c.Request.URL.Path), slog.String("client", c.ClientIP()), slog.Any("user_data", userDTO))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.Log.Info("User was created", slog.String("method", c.Request.Method), slog.Int("code", http.StatusCreated), slog.String("url", c.Request.URL.Path), slog.String("client", c.ClientIP()), slog.Any("user_data", userDTO))

	c.JSON(http.StatusCreated, userDTO)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil || id == 0{
		h.Log.Error("Failed to delete user. Error: Invalid user ID", slog.String("method", c.Request.Method), slog.Int("code", http.StatusBadRequest), slog.String("url", c.Request.URL.Path), slog.String("client", c.ClientIP()), slog.String("user_id", idStr))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err = h.Service.DeleteUser(uint(id)); err != nil {
		h.Log.Error(fmt.Sprintf("Failed to delete user. Error: %s", err.Error()), slog.String("method", c.Request.Method), slog.Int("code", http.StatusInternalServerError), slog.String("url", c.Request.URL.Path), slog.String("client", c.ClientIP()), slog.Uint64("user_id", id))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.Log.Info("User was deleted", slog.String("method", c.Request.Method), slog.Int("code", http.StatusOK), slog.String("url", c.Request.URL.Path), slog.String("client", c.ClientIP()), slog.Uint64("user_id", id))
	
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.Log.Error("Failed to receive user. Error: Invalid user ID", slog.String("method", c.Request.Method), slog.Int("code", http.StatusBadRequest), slog.String("url", c.Request.URL.Path), slog.String("client", c.ClientIP()), slog.String("user_id", idStr))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.Service.GetUserByID(uint(id))
	if err != nil {
		h.Log.Error(fmt.Sprintf("Failed to receive user. Error: %s", err.Error()), slog.String("method", c.Request.Method), slog.Int("code", http.StatusInternalServerError), slog.String("url", c.Request.URL.Path), slog.String("client", c.ClientIP()), slog.Uint64("user_id", id))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.Log.Info("User was received", slog.String("method", c.Request.Method), slog.Int("code", http.StatusOK), slog.String("url", c.Request.URL.Path), slog.String("client", c.ClientIP()), slog.Uint64("user_id", id))

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	users, err := h.Service.GetUsers()
	if err != nil {
		h.Log.Error(fmt.Sprintf("Failed to receive users. Error: %s", err.Error()), slog.String("method", c.Request.Method), slog.Int("code", http.StatusInternalServerError), slog.String("url", c.Request.URL.Path), slog.String("client", c.ClientIP()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.Log.Info("Users was received", slog.String("method", c.Request.Method), slog.Int("code", http.StatusOK), slog.String("url", c.Request.URL.Path), slog.String("client", c.ClientIP()))

	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.Log.Error("Failed to update user. Error: Invalid user ID", slog.String("method", c.Request.Method), slog.Int("code", http.StatusBadRequest), slog.String("url", c.Request.URL.Path), slog.String("client", c.ClientIP()), slog.String("user_id", idStr))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var userDTO dto.UserData
	if err := c.ShouldBindJSON(&userDTO); err != nil {
		h.Log.Error(fmt.Sprintf("Failed to update user. Error: %s", err.Error()), slog.String("method", c.Request.Method), slog.Int("code", http.StatusBadRequest), slog.String("url", c.Request.URL.Path), slog.String("client", c.ClientIP()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Service.UpdateUser(uint(id), &userDTO); err != nil {
		h.Log.Error(fmt.Sprintf("Failed to update user. Error: %s", err.Error()), slog.String("method", c.Request.Method), slog.Int("code", http.StatusInternalServerError), slog.String("url", c.Request.URL.Path), slog.String("client", c.ClientIP()), slog.Uint64("user_id", id), slog.Any("user_data", userDTO))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.Log.Info("User was updated", slog.String("method", c.Request.Method), slog.Int("code", http.StatusOK), slog.String("url", c.Request.URL.Path), slog.String("client", c.ClientIP()), slog.Uint64("user_id", id), slog.Any("user_data", userDTO))

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully", "id": id})
}
