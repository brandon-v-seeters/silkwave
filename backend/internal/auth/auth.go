package auth

import "github.com/gin-gonic/gin"

type AuthService struct {
}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (a *AuthService) Logout(c *gin.Context) error {
	ClearSessionCookie(c)
	return nil
}
