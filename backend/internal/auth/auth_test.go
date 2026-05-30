package auth

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestLogoutClearsHostOnlySessionCookieInDevelopment(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Setenv("GO_ENV", "development")
	t.Setenv("ENVIRONMENT", "")

	recorder, context := newTestContext()

	service := NewAuthService()
	if err := service.Logout(context); err != nil {
		t.Fatalf("Logout() error = %v", err)
	}

	cookie := singleSetCookieHeader(t, recorder)
	if strings.Contains(cookie, "Domain=") {
		t.Fatalf("development logout cookie must be host-only, got %q", cookie)
	}
	if strings.Contains(cookie, "Secure") {
		t.Fatalf("development logout cookie must not be secure-only, got %q", cookie)
	}
	if !strings.Contains(cookie, "SameSite=Lax") {
		t.Fatalf("development logout cookie must use SameSite=Lax, got %q", cookie)
	}
}

func TestSetSessionCookieUsesHostOnlyCookieInDevelopment(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Setenv("GO_ENV", "development")
	t.Setenv("ENVIRONMENT", "")

	recorder, context := newTestContext()

	SetSessionCookie(context, "token")

	cookie := singleSetCookieHeader(t, recorder)
	if strings.Contains(cookie, "Domain=") {
		t.Fatalf("development session cookie must be host-only, got %q", cookie)
	}
	if strings.Contains(cookie, "Secure") {
		t.Fatalf("development session cookie must not be secure-only, got %q", cookie)
	}
	if !strings.Contains(cookie, "SameSite=Lax") {
		t.Fatalf("development session cookie must use SameSite=Lax, got %q", cookie)
	}
}

func TestSetSessionCookieUsesProductionDomainOutsideDevelopment(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Setenv("GO_ENV", "production")
	t.Setenv("ENVIRONMENT", "")
	t.Setenv("COOKIE_DOMAIN", "auth.example.test")

	recorder, context := newTestContext()

	SetSessionCookie(context, "token")

	cookie := singleSetCookieHeader(t, recorder)
	if !strings.Contains(cookie, "Domain=auth.example.test") {
		t.Fatalf("production session cookie must use configured domain, got %q", cookie)
	}
	if !strings.Contains(cookie, "Secure") {
		t.Fatalf("production session cookie must be secure-only, got %q", cookie)
	}
	if !strings.Contains(cookie, "SameSite=None") {
		t.Fatalf("production session cookie must use SameSite=None, got %q", cookie)
	}
}

func newTestContext() (*httptest.ResponseRecorder, *gin.Context) {
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)

	return recorder, context
}

func singleSetCookieHeader(t *testing.T, recorder *httptest.ResponseRecorder) string {
	t.Helper()

	cookies := recorder.Header().Values("Set-Cookie")
	if len(cookies) != 1 {
		t.Fatalf("expected 1 Set-Cookie header, got %d: %v", len(cookies), cookies)
	}

	return cookies[0]
}
