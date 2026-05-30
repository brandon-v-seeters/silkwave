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

func TestReleaseRepositoryListPublishedExcludesArchived(t *testing.T) {
	ctx := context.Background()
	db := setupReleaseRepositoryTestDB(t, ctx)
	repo := NewReleaseRepository(db)

	artistKey := createArtist(t, ctx, db, "artist-list")
	createRelease(t, ctx, db, models.Release{
		Id:        "published-list",
		Slug:      "published-list",
		Title:     "Published List",
		ArtistKey: artistKey,
		Status:    models.ReleaseStatusPublished,
		PublishAt: time.Now(),
	})
	createRelease(t, ctx, db, models.Release{
		Id:        "archived-list",
		Slug:      "archived-list",
		Title:     "Archived List",
		ArtistKey: artistKey,
		Status:    models.ReleaseStatusArchived,
		PublishAt: time.Now().Add(time.Hour),
	})

	releases, err := repo.ListPublished(ctx, ListPublishedReleasesParams{Limit: 20})
	if err != nil {
		t.Fatalf("ListPublished returned error: %v", err)
	}

	if len(releases) != 1 {
		t.Fatalf("expected 1 listed release, got %d", len(releases))
	}
	if releases[0].Slug != "published-list" {
		t.Fatalf("expected published release, got %q", releases[0].Slug)
	}
}

func TestReleaseRepositoryGetPublishedByArtistSlugAndReleaseSlugScopesToArtist(t *testing.T) {
	ctx := context.Background()
	db := setupReleaseRepositoryTestDB(t, ctx)
	repo := NewReleaseRepository(db)

	firstArtistKey := createArtist(t, ctx, db, "first-artist")
	secondArtistKey := createArtist(t, ctx, db, "second-artist")
	createRelease(t, ctx, db, models.Release{
		Id:        "first-shared-slug",
		Slug:      "shared-slug",
		Title:     "First Shared Slug",
		ArtistKey: firstArtistKey,
		Status:    models.ReleaseStatusPublished,
		PublishAt: time.Now(),
	})
	createRelease(t, ctx, db, models.Release{
		Id:        "second-shared-slug",
		Slug:      "shared-slug",
		Title:     "Second Shared Slug",
		ArtistKey: secondArtistKey,
		Status:    models.ReleaseStatusPublished,
		PublishAt: time.Now(),
	})

	firstRelease, err := repo.GetPublishedByArtistSlugAndReleaseSlug(ctx, "first-artist", "shared-slug")
	if err != nil {
		t.Fatalf("GetPublishedByArtistSlugAndReleaseSlug for first artist returned error: %v", err)
	}
	if firstRelease.Id != "first-shared-slug" {
		t.Fatalf("expected first artist release, got %q", firstRelease.Id)
	}

	secondRelease, err := repo.GetPublishedByArtistSlugAndReleaseSlug(ctx, "second-artist", "shared-slug")
	if err != nil {
		t.Fatalf("GetPublishedByArtistSlugAndReleaseSlug for second artist returned error: %v", err)
	}
	if secondRelease.Id != "second-shared-slug" {
		t.Fatalf("expected second artist release, got %q", secondRelease.Id)
	}
}

func TestReleaseRepositoryGetPublishedByArtistSlugAndReleaseSlugNotFoundForArchived(t *testing.T) {
	ctx := context.Background()
	db := setupReleaseRepositoryTestDB(t, ctx)
	repo := NewReleaseRepository(db)

	artistKey := createArtist(t, ctx, db, "artist-archived")
	createRelease(t, ctx, db, models.Release{
		Id:        "archived-detail",
		Slug:      "archived-detail",
		Title:     "Archived Detail",
		ArtistKey: artistKey,
		Status:    models.ReleaseStatusArchived,
		PublishAt: time.Now(),
	})

	_, err := repo.GetPublishedByArtistSlugAndReleaseSlug(ctx, "artist-archived", "archived-detail")
	if err == nil {
		t.Fatal("expected archived release lookup to fail")
	}
}

func TestReleaseRepositoryCreateDraftWithTracks(t *testing.T) {
	ctx := context.Background()
	db := setupReleaseRepositoryTestDB(t, ctx)
	repo := NewReleaseRepository(db)

	artistKey := createArtist(t, ctx, db, "artist-draft")
	release := models.Release{
		Id:        "draft-create",
		Slug:      "draft-create",
		Title:     "Draft Create",
		ArtistKey: artistKey,
		Status:    models.ReleaseStatusDraft,
	}
	tracks := []models.Track{{
		Id:    "track-one",
		Title: "Track One",
		Order: 1,
	}}

	releaseKey, err := repo.CreateDraftWithTracks(ctx, release, tracks)
	if err != nil {
		t.Fatalf("CreateDraftWithTracks returned error: %v", err)
	}
	if releaseKey == "" {
		t.Fatal("expected release key")
	}

	draft, err := repo.GetDraftWithTracksByKey(ctx, releaseKey)
	if err != nil {
		t.Fatalf("GetDraftWithTracksByKey returned error: %v", err)
	}
	if len(draft.Tracks) != 1 {
		t.Fatalf("expected 1 track, got %d", len(draft.Tracks))
	}
	if draft.Tracks[0].ReleaseKey != releaseKey {
		t.Fatalf("expected track releaseKey %q, got %q", releaseKey, draft.Tracks[0].ReleaseKey)
	}
}

func TestReleaseRepositoryCreateDraftWithTracksDuplicateSlug(t *testing.T) {
	ctx := context.Background()
	db := setupReleaseRepositoryTestDB(t, ctx)
	repo := NewReleaseRepository(db)

	artistKey := createArtist(t, ctx, db, "artist-duplicate")
	release := models.Release{
		Id:        "draft-duplicate-one",
		Slug:      "duplicate-slug",
		Title:     "Draft Duplicate One",
		ArtistKey: artistKey,
		Status:    models.ReleaseStatusDraft,
	}

	if _, err := repo.CreateDraftWithTracks(ctx, release, nil); err != nil {
		t.Fatalf("first CreateDraftWithTracks returned error: %v", err)
	}

	release.Id = "draft-duplicate-two"
	_, err := repo.CreateDraftWithTracks(ctx, release, nil)
	if !errors.Is(err, ErrDuplicateReleaseSlug) {
		t.Fatalf("expected ErrDuplicateReleaseSlug, got %v", err)
	}
}

func TestReleaseRepositoryGetReleaseByIdNotFound(t *testing.T) {
	ctx := context.Background()
	db := setupReleaseRepositoryTestDB(t, ctx)
	repo := NewReleaseRepository(db)

	_, err := repo.GetReleaseById(ctx, "missing")
	if err == nil {
		t.Fatal("expected missing release lookup to fail")
	}
}

func setupReleaseRepositoryTestDB(t *testing.T, ctx context.Context) *database.ArangoDB {
	t.Helper()

	if os.Getenv("SILKWAVE_ARANGO_INTEGRATION") != "1" {
		t.Skip("set SILKWAVE_ARANGO_INTEGRATION=1 to run ArangoDB repository integration tests")
	}

	cfg := &config.Config{
		ArangoEndpoint: os.Getenv("ARANGO_ENDPOINT"),
		ArangoDatabase: fmt.Sprintf("silkwave_release_repo_test_%d", time.Now().UnixNano()),
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

func createArtist(t *testing.T, ctx context.Context, db *database.ArangoDB, slug string) string {
	t.Helper()

	col, err := db.GetCollection(ctx, "Artists")
	if err != nil {
		t.Fatalf("get Artists collection: %v", err)
	}

	meta, err := col.CreateDocument(ctx, models.Artist{
		Id:        slug,
		Name:      slug,
		Slug:      slug,
		CreatedAt: time.Now(),
	})
	if err != nil {
		t.Fatalf("create artist: %v", err)
	}

	return meta.Key
}

func createRelease(t *testing.T, ctx context.Context, db *database.ArangoDB, release models.Release) string {
	t.Helper()

	col, err := db.GetCollection(ctx, releasesCollection)
	if err != nil {
		t.Fatalf("get Releases collection: %v", err)
	}

	meta, err := col.CreateDocument(ctx, release)
	if err != nil {
		t.Fatalf("create release: %v", err)
	}

	return meta.Key
}
