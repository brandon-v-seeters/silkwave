package auth

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type AuthService struct {
}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (a *AuthService) Logout(c *gin.Context) error {
	c.SetSameSite(http.SameSiteNoneMode)
	domain := "silkwave.io"
	secure := true

	if os.Getenv("ENVIRONMENT") == "development" {
		domain = "localhost"
		secure = false
		c.SetSameSite(http.SameSiteLaxMode)
	}

	c.SetCookie("session", "", -1, "/", domain, secure, true)
	return nil
}
