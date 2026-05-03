package middleware

import (
	"net/http"
	"strings"

	"github.com/brandon-v-seeters/go-silk-wave/internal/auth"
	"github.com/brandon-v-seeters/go-silk-wave/internal/database"
	"github.com/gin-gonic/gin"
)

const (
	UserKeyContext = "userKey"
)

// AuthMiddleware validates JWT tokens from cookies or Authorization header
func AuthMiddleware(jwtService *auth.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		// First try to get token from cookie
		if cookie, err := c.Cookie("session"); err == nil {
			tokenString = cookie
		}

		// Fallback to Authorization header
		if tokenString == "" {
			authHeader := c.GetHeader("Authorization")
			if strings.HasPrefix(authHeader, "Token ") {
				tokenString = strings.TrimPrefix(authHeader, "Token ")
			} else if strings.HasPrefix(authHeader, "Bearer ") {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No session token"})
			c.Abort()
			return
		}

		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired session token"})
			c.Abort()
			return
		}

		// Set user key in context for handlers to use
		c.Set(UserKeyContext, claims.Key)
		c.Next()
	}
}

// OptionalAuthMiddleware extracts user info if token exists, but doesn't require it
func OptionalAuthMiddleware(jwtService *auth.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		if cookie, err := c.Cookie("session"); err == nil {
			tokenString = cookie
		}

		if tokenString == "" {
			authHeader := c.GetHeader("Authorization")
			if strings.HasPrefix(authHeader, "Token ") {
				tokenString = strings.TrimPrefix(authHeader, "Token ")
			} else if strings.HasPrefix(authHeader, "Bearer ") {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		if tokenString != "" {
			if claims, err := jwtService.ValidateToken(tokenString); err == nil {
				c.Set(UserKeyContext, claims.Key)
			}
		}

		c.Next()
	}
}

// GetUserKey retrieves the user key from context
func GetUserKey(c *gin.Context) (string, bool) {
	userKey, exists := c.Get(UserKeyContext)
	if !exists {
		return "", false
	}
	return userKey.(string), true
}

// MustGetUserKey retrieves the user key from context for protected routes.
// Only use this in handlers behind AuthMiddleware - it panics if userKey is missing.
func MustGetUserKey(c *gin.Context) string {
	userKey, _ := c.Get(UserKeyContext)
	return userKey.(string)
}

func UserManagesArtist(c *gin.Context, db *database.ArangoDB, artistKey string) bool {
	userKey, _ := GetUserKey(c)
	if userKey == "" {
		return false
	}

	query := /*aql*/ `
		FOR a IN Artists
		FILTER a.userKey == @userKey && a._key == @artistKey
		RETURN a._key
	`
	bindVars := map[string]interface{}{
		"userKey":   userKey,
		"artistKey": artistKey,
	}

	result, err := database.QueryOne[string](c.Request.Context(), db, query, bindVars)
	if err != nil {
		return false
	}

	return result != nil
}

// CheckToken validates the token and sets user key in context if not already set
// Returns the user key and any error that occurred
func CheckToken(c *gin.Context, jwtService *auth.JWTService) (string, error) {
	// Check if user is already set in context
	if userKey, exists := GetUserKey(c); exists && userKey != "" {
		return userKey, nil
	}

	// Extract token from cookie or header
	var tokenString string

	if cookie, err := c.Cookie("session"); err == nil {
		tokenString = cookie
	}

	if tokenString == "" {
		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Token ") {
			tokenString = strings.TrimPrefix(authHeader, "Token ")
		} else if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		}
	}

	if tokenString == "" {
		return "", auth.ErrMissingToken
	}

	// Validate token
	claims, err := jwtService.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	// Set user key in context
	c.Set(UserKeyContext, claims.Key)

	return claims.Key, nil
}
