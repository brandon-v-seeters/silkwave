package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/brandon-v-seeters/go-silk-wave/internal/auth"
	"github.com/brandon-v-seeters/go-silk-wave/internal/logger"
	"github.com/brandon-v-seeters/go-silk-wave/internal/models"
	"github.com/brandon-v-seeters/go-silk-wave/internal/repository"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthHandler struct {
	userRepo        *repository.UserRepository
	jwtService      *auth.JWTService
	passwordService *auth.PasswordService
	authService     *auth.AuthService
}

func NewAuthHandler(userRepo *repository.UserRepository, jwtService *auth.JWTService, passwordService *auth.PasswordService, authService *auth.AuthService) *AuthHandler {
	return &AuthHandler{
		userRepo:        userRepo,
		jwtService:      jwtService,
		passwordService: passwordService,
		authService:     authService,
	}
}

// Login handles POST /api/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "invalid_request", "Missing required parameters.")
		return
	}

	user, err := h.userRepo.GetCredentialsByEmail(c.Request.Context(), req.Email)
	if err != nil {
		if errors.Is(err, repository.ErrUserCredentialsNotFound) {
			logger.Warn("Incorrect login details", zap.String("email", req.Email))
			respondError(c, http.StatusUnauthorized, "unauthorized", "Incorrect login details.")
			return
		}

		logger.Error("Failed to get user credentials", err, zap.String("email", req.Email))
		respondError(c, http.StatusInternalServerError, "internal_error", "Failed to login.")
		return
	}

	if !h.passwordService.CheckPassword(req.Password, user.Password) {
		logger.Warn("Incorrect login details", zap.String("email", req.Email))
		respondError(c, http.StatusUnauthorized, "unauthorized", "Incorrect login details.")
		return
	}

	token, err := h.jwtService.GenerateToken(user.Key)
	if err != nil {
		logger.Error("Failed to generate session", err, zap.String("email", req.Email))
		respondError(c, http.StatusInternalServerError, "internal_error", "Failed to generate session.")
		return
	}

	auth.SetSessionCookie(c, token)

	respondOK(c, nil, "Login successful.")
}

// Register handles POST /api/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "invalid_request", "Missing required parameters.")
		return
	}

	hashedPassword, err := h.passwordService.HashPassword(req.Password)
	if err != nil {
		logger.Error("Failed to process password", err, zap.String("email", req.Email))
		respondError(c, http.StatusInternalServerError, "internal_error", "Failed to process password.")
		return
	}

	tempUsername, err := h.passwordService.HashPassword(req.Email)
	if err != nil {
		logger.Error("Failed to generate temporary username", err, zap.String("email", req.Email))
		respondError(c, http.StatusInternalServerError, "internal_error", "User could not be created.")
		return
	}

	user := &models.UserWithCredentials{
		User: models.User{
			Username:          tempUsername[:20], // Truncate for reasonable length
			Email:             req.Email,
			InvalidArtistName: true,
			ValidatedEmail:    false,
			Role:              models.UserRoleUser,
			Settings: models.UserSettings{
				ReceiveEmails: true,
			},
		},
		Password:              hashedPassword,
		EmailUnsubscribeToken: "",
		Socials:               make(map[string]models.SocialWithCredentials),
	}

	created, err := h.userRepo.CreateWithCredentials(c.Request.Context(), user)
	if err != nil {
		if isUniqueConstraintError(err) {
			logger.Warn("Email already exists", zap.String("email", req.Email))
			respondError(c, http.StatusConflict, "email_taken", "Email already exists.")
			return
		}

		logger.Error("Failed to create user", err, zap.String("email", req.Email))
		respondError(c, http.StatusBadRequest, "invalid_request", "User could not be created.")
		return
	}

	token, err := h.jwtService.GenerateToken(created.Key)
	if err != nil {
		logger.Error("Failed to generate session", err, zap.String("email", req.Email))
		respondError(c, http.StatusInternalServerError, "internal_error", "Failed to generate session.")
		return
	}

	auth.SetSessionCookie(c, token)

	respondOK(c, nil, "User created successfully.")
}

// Logout handles POST /api/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	err := h.authService.Logout(c)
	if err != nil {
		logger.Error("Failed to logout", err)
		respondError(c, http.StatusInternalServerError, "internal_error", "Failed to logout.")
		return
	}

	respondOK(c, nil, "Logged out successfully.")
}

func isUniqueConstraintError(err error) bool {
	return err != nil && (strings.Contains(err.Error(), "unique constraint") || strings.Contains(err.Error(), "duplicate"))
}
