package handlers

import (
	"fmt"
	"net/http"

	"github.com/brandon-v-seeters/go-silk-wave/internal/database"
	"github.com/brandon-v-seeters/go-silk-wave/internal/middleware"
	"github.com/brandon-v-seeters/go-silk-wave/internal/models"
	"github.com/brandon-v-seeters/go-silk-wave/internal/repository"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	repo *repository.UserRepository
	db   *database.ArangoDB
}

func NewUserHandler(repo *repository.UserRepository, db *database.ArangoDB) *UserHandler {
	return &UserHandler{repo: repo, db: db}
}

// CreateUser handles POST /api/users (registration)
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req models.RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &models.User{
		Email:          req.Email,
		Role:           models.UserRoleUser,
		ValidatedEmail: false,
		Settings: models.UserSettings{
			ReceiveEmails: true,
		},
	}

	created, err := h.repo.Create(c.Request.Context(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, created)
}

// GetUser handles GET /api/users/:key
func (h *UserHandler) GetUser(c *gin.Context) {
	key := c.Param("key")

	user, err := h.repo.GetByKey(c.Request.Context(), key)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Println(user)

	c.JSON(http.StatusOK, user)
}

// GetAllUsers handles GET /api/users
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.repo.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

// // UpdateUser handles PUT /api/users/:key
// func (h *UserHandler) UpdateUser(c *gin.Context) {
// 	key := c.Param("key")

// 	var req models.UpdateUserRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Get existing user first
// 	existing, err := h.repo.GetByKey(c.Request.Context(), key)
// 	if err != nil {
// 		if err.Error() == "user not found" {
// 			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
// 			return
// 		}
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Update fields if provided
// 	if req.Username != "" {
// 		existing.Username = req.Username
// 	}
// 	if req.Email != "" {
// 		existing.Email = req.Email
// 	}
// 	if req.Settings != nil {
// 		existing.Settings = *req.Settings
// 	}

// 	updated, err := h.repo.Update(c.Request.Context(), key, existing)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, updated)
// }

// DeleteUser handles DELETE /api/users/:key
func (h *UserHandler) DeleteUser(c *gin.Context) {
	key := c.Param("key")

	err := h.repo.Delete(c.Request.Context(), key)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}

// UpdateMode handles POST /api/user/mode
func (h *UserHandler) UpdateMode(c *gin.Context) {
	userKey, ok := middleware.GetUserKey(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	var req struct {
		Mode string `json:"mode" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := database.QueryOne[string](c.Request.Context(), h.db,
		"UPDATE @userKey WITH { mode: @mode } IN Users RETURN NEW.mode",
		map[string]interface{}{
			"userKey": userKey,
			"mode":    req.Mode,
		})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update mode"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mode": req.Mode})
}

// GetCurrentUser handles GET /api/user (get logged in user)
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	fmt.Println("GetCurrentUser")
	userKey, ok := middleware.GetUserKey(c)
	fmt.Println("userKey", userKey)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	user, err := h.repo.GetByKey(c.Request.Context(), userKey)
	fmt.Println("user", user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}
