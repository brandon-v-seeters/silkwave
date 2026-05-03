package handlers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/brandon-v-seeters/go-silk-wave/internal/auth"
	"github.com/brandon-v-seeters/go-silk-wave/internal/database"
	"github.com/brandon-v-seeters/go-silk-wave/internal/logger"
	"github.com/brandon-v-seeters/go-silk-wave/internal/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthHandler struct {
	db              *database.ArangoDB
	jwtService      *auth.JWTService
	passwordService *auth.PasswordService
	authService     *auth.AuthService
}

func NewAuthHandler(db *database.ArangoDB, jwtService *auth.JWTService, passwordService *auth.PasswordService, authService *auth.AuthService) *AuthHandler {
	return &AuthHandler{
		db:              db,
		jwtService:      jwtService,
		passwordService: passwordService,
		authService:     authService,
	}
}

type PasswordAndKey struct {
	Password string `json:"password"`
	Key      string `json:"key"`
}

// Login handles POST /api/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Missing required parameters."})
		return
	}

	if req.Email == "" || req.Password == "" {
		logger.Warn("Missing required parameters", zap.String("email", req.Email), zap.String("password", req.Password))
		c.JSON(http.StatusBadRequest, gin.H{"message": "Missing required parameters."})
		return
	}

	query := `
	FOR user IN Users FILTER user.email == @email RETURN { 
			password: user.password,
			key: user._key
		}
	`

	// Find user by email
	user, err := database.QueryOne[PasswordAndKey](c.Request.Context(), h.db, query, map[string]interface{}{"email": req.Email})

	if err != nil || user.Password == "" || user.Key == "" {
		logger.Error("Failed to find user", err, zap.String("email", req.Email))
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Incorrect login details."})
		return
	}

	// Check password
	if !h.passwordService.CheckPassword(req.Password, user.Password) {
		logger.Warn("Incorrect login details", zap.String("email", req.Email))
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Incorrect login details."})
		return
	}

	// Generate token
	token, err := h.jwtService.GenerateToken(user.Key)
	if err != nil {
		logger.Error("Failed to generate session", err, zap.String("email", req.Email))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate session."})
		return
	}

	// Set cookie
	h.setSessionCookie(c, token)

	c.JSON(http.StatusOK, gin.H{"message": "Login successful."})
}

// Register handles POST /api/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Missing required parameters."})
		return
	}

	if req.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "An email is required."})
		return
	}

	if req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "A password is required."})
		return
	}

	// Hash password
	hashedPassword, err := h.passwordService.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to process password."})
		return
	}

	// Generate a temporary username from email hash
	tempUsername, _ := h.passwordService.HashPassword(req.Email)

	// Create user
	user := models.UserWithCredentials{
		User: models.User{
			Username:          tempUsername[:20], // Truncate for reasonable length
			Email:             req.Email,
			InvalidArtistName: true,
			ValidatedEmail:    false,
			Role:              models.UserRoleUser,
			Settings: models.UserSettings{
				ReceiveEmails: true,
			},
			CreatedAt: time.Now(),
		},
		Password:              hashedPassword,
		EmailUnsubscribeToken: "",
		Socials:               make(map[string]models.SocialWithCredentials),
	}

	fmt.Println(user)

	// Insert user
	result, err := database.QueryOne[string](c.Request.Context(), h.db,
		"INSERT @user INTO Users RETURN NEW._key",
		map[string]interface{}{"user": user})

	if err != nil {
		// Check for unique constraint violation
		if isUniqueConstraintError(err) {
			logger.Error("Unique constraint violation", err, zap.String("email", req.Email))
			c.JSON(http.StatusConflict, gin.H{"message": "Email already exists."})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"message": "User could not be created."})
		return
	}

	// Generate token
	token, err := h.jwtService.GenerateToken(*result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate session."})
		return
	}

	// Set cookie
	h.setSessionCookie(c, token)

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully."})
}

// Logout handles POST /api/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	err := h.authService.Logout(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to logout."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully."})
}

func (h *AuthHandler) setSessionCookie(c *gin.Context, token string) {
	// 30 days expiry
	maxAge := 60 * 60 * 24 * 30
	domain := "silkwave.io"
	secure := true
	c.SetSameSite(http.SameSiteNoneMode)

	if os.Getenv("ENVIRONMENT") == "development" {
		domain = "localhost"
		secure = false
		c.SetSameSite(http.SameSiteLaxMode)
	}

	c.SetCookie("session", token, maxAge, "/", domain, secure, true)
}

func isUniqueConstraintError(err error) bool {
	return err != nil && (contains(err.Error(), "unique constraint") || contains(err.Error(), "duplicate"))
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
