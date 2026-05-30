package repository

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/brandon-v-seeters/go-silk-wave/internal/config"
	"github.com/brandon-v-seeters/go-silk-wave/internal/database"
	"github.com/brandon-v-seeters/go-silk-wave/internal/models"
)

func TestUserRepositoryGetCredentialsByEmail(t *testing.T) {
	ctx := context.Background()
	db := setupUserRepositoryTestDB(t, ctx)
	repo := NewUserRepository(db)

	seedUserWithCredentials(t, ctx, db, models.UserWithCredentials{
		User: models.User{
			Username:       "login-user",
			Email:          "login@example.com",
			ValidatedEmail: true,
			Role:           models.UserRoleUser,
			Settings: models.UserSettings{
				ReceiveEmails: true,
			},
			CreatedAt: time.Now(),
		},
		Password:              "hashed-password",
		EmailUnsubscribeToken: "unsubscribe-token",
		Socials:               map[string]models.SocialWithCredentials{},
	})

	credentials, err := repo.GetCredentialsByEmail(ctx, "login@example.com")
	if err != nil {
		t.Fatalf("GetCredentialsByEmail returned error: %v", err)
	}

	if credentials.Key == "" {
		t.Fatal("expected user key")
	}
	if credentials.Password != "hashed-password" {
		t.Fatalf("expected password %q, got %q", "hashed-password", credentials.Password)
	}
}

func TestUserRepositoryGetCredentialsByEmailNotFound(t *testing.T) {
	ctx := context.Background()
	db := setupUserRepositoryTestDB(t, ctx)
	repo := NewUserRepository(db)

	_, err := repo.GetCredentialsByEmail(ctx, "missing@example.com")
	if !errors.Is(err, ErrUserCredentialsNotFound) {
		t.Fatalf("expected ErrUserCredentialsNotFound, got %v", err)
	}
}

func TestUserRepositoryGetByKeyLoadsArtistThroughUserArtistEdge(t *testing.T) {
	ctx := context.Background()
	db := setupUserRepositoryTestDB(t, ctx)
	repo := NewUserRepository(db)

	userKey := seedUserWithCredentials(t, ctx, db, models.UserWithCredentials{
		User: models.User{
			Username:       "artist-edge-user",
			Email:          "artist-edge@example.com",
			ValidatedEmail: true,
			Role:           models.UserRoleUser,
			Settings: models.UserSettings{
				ReceiveEmails: true,
			},
			CreatedAt: time.Now(),
		},
		Password:              "hashed-password",
		EmailUnsubscribeToken: "unsubscribe-token",
		Socials:               map[string]models.SocialWithCredentials{},
	})
	artistKey := seedArtist(t, ctx, db, models.Artist{
		DocumentMeta: models.DocumentMeta{Key: "artist_edge_loaded"},
		Id:           "artist-edge-loaded",
		Name:         "Artist Edge Loaded",
		Slug:         "artist-edge-loaded",
		CreatedAt:    time.Now(),
		PublishedAt:  time.Now(),
	})
	seedUserArtist(t, ctx, db, models.UserArtist{
		Edge: models.Edge{
			Key:       "user_artist_edge_loaded",
			From:      "Users/" + userKey,
			To:        "Artists/" + artistKey,
			CreatedAt: time.Now(),
		},
	})

	user, err := repo.GetByKey(ctx, userKey)
	if err != nil {
		t.Fatalf("GetByKey returned error: %v", err)
	}
	if user.Artist == nil {
		t.Fatal("expected GetByKey to load artist through UsersArtists edge")
	}
	if user.Artist.Key != artistKey {
		t.Fatalf("expected artist key %q, got %q", artistKey, user.Artist.Key)
	}
	if user.Artist.Slug != "artist-edge-loaded" {
		t.Fatalf("expected artist slug %q, got %q", "artist-edge-loaded", user.Artist.Slug)
	}
}

func setupUserRepositoryTestDB(t *testing.T, ctx context.Context) *database.ArangoDB {
	t.Helper()

	if os.Getenv("SILKWAVE_ARANGO_INTEGRATION") != "1" {
		t.Skip("set SILKWAVE_ARANGO_INTEGRATION=1 to run ArangoDB repository integration tests")
	}

	cfg := &config.Config{
		ArangoEndpoint: os.Getenv("ARANGO_ENDPOINT"),
		ArangoDatabase: fmt.Sprintf("silkwave_user_repo_test_%d", time.Now().UnixNano()),
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

func seedUserWithCredentials(t *testing.T, ctx context.Context, db *database.ArangoDB, user models.UserWithCredentials) string {
	t.Helper()

	col, err := db.GetCollection(ctx, usersCollection)
	if err != nil {
		t.Fatalf("get Users collection: %v", err)
	}

	meta, err := col.CreateDocument(ctx, user)
	if err != nil {
		t.Fatalf("create user: %v", err)
	}

	return meta.Key
}

func seedArtist(t *testing.T, ctx context.Context, db *database.ArangoDB, artist models.Artist) string {
	t.Helper()

	col, err := db.GetCollection(ctx, "Artists")
	if err != nil {
		t.Fatalf("get Artists collection: %v", err)
	}

	meta, err := col.CreateDocument(ctx, artist)
	if err != nil {
		t.Fatalf("create artist: %v", err)
	}

	return meta.Key
}

func seedUserArtist(t *testing.T, ctx context.Context, db *database.ArangoDB, userArtist models.UserArtist) string {
	t.Helper()

	col, err := db.GetCollection(ctx, "UsersArtists")
	if err != nil {
		t.Fatalf("get UsersArtists collection: %v", err)
	}

	meta, err := col.CreateDocument(ctx, userArtist)
	if err != nil {
		t.Fatalf("create user artist edge: %v", err)
	}

	return meta.Key
}
