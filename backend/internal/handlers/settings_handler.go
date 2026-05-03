package handlers

import (
	"net/http"

	"github.com/brandon-v-seeters/go-silk-wave/internal/auth"
	"github.com/brandon-v-seeters/go-silk-wave/internal/database"
	"github.com/brandon-v-seeters/go-silk-wave/internal/logger"
	"github.com/brandon-v-seeters/go-silk-wave/internal/middleware"
	"github.com/brandon-v-seeters/go-silk-wave/internal/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SettingsHandler struct {
	passwordService *auth.PasswordService
	authService     *auth.AuthService
	db              *database.ArangoDB
}

func NewSettingsHandler(passwordService *auth.PasswordService, db *database.ArangoDB) *SettingsHandler {
	return &SettingsHandler{passwordService: passwordService, db: db}
}

func (h *SettingsHandler) UpdateEmail(c *gin.Context) {
	userKey, ok := middleware.GetUserKey(c)
	if !ok {
		logger.Warn("UpdateEmail called without authentication")
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Not authenticated"})
		return
	}

	var req models.UpdateEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Debug("UpdateEmail invalid request", zap.Error(err), zap.String("userKey", userKey))
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	userPasswordQuery := `
	FOR user IN Users FILTER user._key == @userKey RETURN user.password
	`

	userPassword, err := database.QueryOne[string](c.Request.Context(), h.db, userPasswordQuery, map[string]interface{}{"userKey": userKey})
	if err != nil {
		logger.Error("Failed to fetch user password", err, zap.String("userKey", userKey))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to verify password"})
		return
	}

	if !h.passwordService.CheckPassword(req.Password, *userPassword) {
		logger.Warn("Invalid password attempt for email update", zap.String("userKey", userKey))
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid password"})
		return
	}

	query := `
	UPDATE @userKey WITH { email: @email } IN Users RETURN NEW.email
	`

	bindVars := map[string]interface{}{
		"userKey": userKey,
		"email":   req.Email,
	}

	result, err := database.QueryOne[string](c.Request.Context(), h.db, query, bindVars)
	if err != nil {
		logger.Error("Failed to update email", err, zap.String("userKey", userKey), zap.String("newEmail", req.Email))

		if database.IsUniqueConstraintError(err) {
			dbErr := database.ParseError(err)
			c.JSON(http.StatusConflict, gin.H{"message": dbErr.Message, "field": dbErr.Field})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update email"})
		return
	}

	logger.Info("Email updated successfully", zap.String("userKey", userKey), zap.String("newEmail", req.Email))
	c.JSON(http.StatusOK, gin.H{"email": *result})
}

func (h *SettingsHandler) DeleteAccount(c *gin.Context) {
	var req models.DeleteAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Debug("DeleteAccount invalid request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	userKey, ok := middleware.GetUserKey(c)
	if !ok {
		logger.Warn("DeleteAccount called without authentication")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	query := `
	FOR user IN Users FILTER user._key == @userKey RETURN user.password
	`
	userVars := map[string]interface{}{
		"userKey": userKey,
	}

	userPassword, err := database.QueryOne[string](c.Request.Context(), h.db, query, userVars)
	if err != nil {
		logger.Error("Failed to fetch user password for account deletion", err, zap.String("userKey", userKey))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get user password"})
		return
	}

	if !h.passwordService.CheckPassword(req.Password, *userPassword) {
		logger.Warn("Invalid password attempt for account deletion", zap.String("userKey", userKey))
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid password"})
		return
	}

	deleteAccountQuery := `
	REMOVE {_key: @userKey} IN Users RETURN true
	`

	bindVars := map[string]interface{}{
		"userKey": userKey,
	}

	_, err = database.QueryOne[bool](c.Request.Context(), h.db, deleteAccountQuery, bindVars)
	logger.Debug("User deleted", zap.String("userKey", userKey))
	if err != nil {
		logger.Error("Failed to delete account", err, zap.String("userKey", userKey))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete account"})
		return
	}

	err = h.authService.Logout(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to logout."})
		return
	}

	logger.Info("Account deleted successfully", zap.String("userKey", userKey))

	c.JSON(http.StatusOK, gin.H{"message": "Account deleted successfully"})
}

func (h *SettingsHandler) UpdatePassword(c *gin.Context) {
	var req models.UpdatePasswordRequest

	userKey, _ := middleware.GetUserKey(c)

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Debug("UpdatePassword invalid request", zap.Error(err), zap.String("userKey", userKey))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userPasswordQuery := /*aql*/ `
		FOR user IN Users 
		FILTER user._key == @userKey 
		RETURN user.password
	`
	bindVars := map[string]interface{}{
		"userKey": userKey,
	}

	userPassword, err := database.QueryOne[string](c.Request.Context(), h.db, userPasswordQuery, bindVars)
	if err != nil {
		logger.Error("Failed to fetch user password", err, zap.String("userKey", userKey))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to verify password"})
		return
	}

	if !h.passwordService.CheckPassword(req.OldPassword, *userPassword) {
		logger.Warn("Invalid password attempt for password update", zap.String("userKey", userKey))
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid password"})
		return
	}

	hashedPassword, err := h.passwordService.HashPassword(req.NewPassword)
	if err != nil {
		logger.Error("Failed to hash new password", err, zap.String("userKey", userKey))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to hash new password"})
		return
	}

	updatePasswordQuery := /*aql*/ `
		UPDATE @userKey WITH { password: @password } 
		IN Users 
		RETURN NEW.password
	`

	updatePasswordVars := map[string]interface{}{
		"userKey":  userKey,
		"password": hashedPassword,
	}

	_, err = database.QueryOne[string](c.Request.Context(), h.db, updatePasswordQuery, updatePasswordVars)
	if err != nil {
		logger.Error("Failed to update password", err, zap.String("userKey", userKey))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}
