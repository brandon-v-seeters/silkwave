package middleware

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/brandon-v-seeters/go-silk-wave/internal/config"
	"github.com/brandon-v-seeters/go-silk-wave/internal/database"
	"github.com/gin-gonic/gin"
)

func TestUserManagesArtistUsesUserArtistEdge(t *testing.T) {
	ctx := context.Background()
	db := setupMiddlewareTestDB(t, ctx)

	userKey := seedMiddlewareDocument(t, ctx, db, "Users", map[string]interface{}{
		"username":       "edge-owner",
		"email":          "edge-owner@example.com",
		"validatedEmail": true,
		"role":           "user",
		"settings":       map[string]interface{}{"receiveEmails": true},
		"createdAt":      time.Now(),
	})
	artistKey := seedMiddlewareDocument(t, ctx, db, "Artists", map[string]interface{}{
		"name":      "Edge Owner Artist",
		"slug":      "edge-owner-artist",
		"createdAt": time.Now(),
	})
	otherArtistKey := seedMiddlewareDocument(t, ctx, db, "Artists", map[string]interface{}{
		"name":      "Other Artist",
		"slug":      "other-artist",
		"createdAt": time.Now(),
	})
	seedMiddlewareDocument(t, ctx, db, "UsersArtists", map[string]interface{}{
		"_from":     "Users/" + userKey,
		"_to":       "Artists/" + artistKey,
		"createdAt": time.Now(),
	})

	requestContext := newMiddlewareTestContext(userKey)
	if !UserManagesArtist(requestContext, db, artistKey) {
		t.Fatal("expected user to manage artist through UsersArtists edge")
	}
	if UserManagesArtist(requestContext, db, otherArtistKey) {
		t.Fatal("expected user not to manage artist without UsersArtists edge")
	}
	if UserManagesArtist(newMiddlewareTestContext(""), db, artistKey) {
		t.Fatal("expected missing user context not to manage artist")
	}
}

func setupMiddlewareTestDB(t *testing.T, ctx context.Context) *database.ArangoDB {
	t.Helper()

	if os.Getenv("SILKWAVE_ARANGO_INTEGRATION") != "1" {
		t.Skip("set SILKWAVE_ARANGO_INTEGRATION=1 to run ArangoDB middleware integration tests")
	}

	cfg := &config.Config{
		ArangoEndpoint: os.Getenv("ARANGO_ENDPOINT"),
		ArangoDatabase: fmt.Sprintf("silkwave_middleware_test_%d", time.Now().UnixNano()),
		ArangoUsername: os.Getenv("ARANGO_USERNAME"),
		ArangoPassword: os.Getenv("ARANGO_PASSWORD"),
	}
	if cfg.ArangoEndpoint == "" {
		cfg.ArangoEndpoint = "http://localhost:8529"
	}
	if cfg.ArangoUsername == "" {
		cfg.ArangoUsername = "root"
	}

	db, err := database.NewArangoDB(cfg)
	if err != nil {
		t.Fatalf("NewArangoDB returned error: %v", err)
	}
	t.Cleanup(func() {
		if err := db.Database.Remove(context.Background()); err != nil {
			t.Logf("remove test database %q: %v", cfg.ArangoDatabase, err)
		}
	})

	if err := db.Migrate(ctx); err != nil {
		t.Fatalf("Migrate returned error: %v", err)
	}

	return db
}

func seedMiddlewareDocument(t *testing.T, ctx context.Context, db *database.ArangoDB, collection string, document map[string]interface{}) string {
	t.Helper()

	col, err := db.GetCollection(ctx, collection)
	if err != nil {
		t.Fatalf("get %s collection: %v", collection, err)
	}

	meta, err := col.CreateDocument(ctx, document)
	if err != nil {
		t.Fatalf("create %s document: %v", collection, err)
	}

	return meta.Key
}

func newMiddlewareTestContext(userKey string) *gin.Context {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	if userKey != "" {
		c.Set(UserKeyContext, userKey)
	}

	return c
}
