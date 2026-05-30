package auth

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	sessionCookieName              = "session"
	sessionCookieMaxAge            = 60 * 60 * 24 * 30
	defaultProductionCookieDomain  = "silkwave.io"
	sessionCookiePath              = "/"
	sessionCookieHTTPOnly          = true
	sessionCookieProductionSecure  = true
	sessionCookieDevelopmentSecure = false
)

func SetSessionCookie(c *gin.Context, token string) {
	setSessionCookie(c, token, sessionCookieMaxAge)
}

func ClearSessionCookie(c *gin.Context) {
	setSessionCookie(c, "", -1)
}

func setSessionCookie(c *gin.Context, value string, maxAge int) {
	domain, secure, sameSite := sessionCookieOptions()

	c.SetSameSite(sameSite)
	c.SetCookie(sessionCookieName, value, maxAge, sessionCookiePath, domain, secure, sessionCookieHTTPOnly)
}

func sessionCookieOptions() (string, bool, http.SameSite) {
	if isDevelopmentEnvironment() {
		return "", sessionCookieDevelopmentSecure, http.SameSiteLaxMode
	}

	return productionCookieDomain(), sessionCookieProductionSecure, http.SameSiteNoneMode
}

func isDevelopmentEnvironment() bool {
	env := normalizedEnv("GO_ENV")
	if env == "" {
		env = normalizedEnv("ENVIRONMENT")
	}

	return env == "" || env == "development" || env == "dev" || env == "local"
}

func normalizedEnv(key string) string {
	return strings.ToLower(strings.TrimSpace(os.Getenv(key)))
}

func productionCookieDomain() string {
	if domain := strings.TrimSpace(os.Getenv("COOKIE_DOMAIN")); domain != "" {
		return domain
	}

	return defaultProductionCookieDomain
}
