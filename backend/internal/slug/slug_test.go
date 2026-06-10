package slug

import (
	"context"
	"errors"
	"testing"

	"github.com/brandon-v-seeters/go-silk-wave/internal/models"
)

func TestSlugify(t *testing.T) {
	got := Slugify("Live at the Roxy!")
	if got != "live-at-the-roxy" {
		t.Fatalf("expected live-at-the-roxy, got %q", got)
	}
}

func TestEnsureUnique(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		exists    bool
		storeErr  error
		expectErr error
	}{
		{name: "available"},
		{name: "duplicate", exists: true, expectErr: ErrNotUnique},
		{name: "store error", storeErr: errSlugStore, expectErr: errSlugStore},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &fakeReleaseStore{exists: tt.exists, existsErr: tt.storeErr}
			err := EnsureUnique(context.Background(), store, "artist-key", "live-at-the-roxy")

			if tt.expectErr == nil && err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}
			if tt.expectErr != nil && !errors.Is(err, tt.expectErr) {
				t.Fatalf("expected %v, got %v", tt.expectErr, err)
			}
		})
	}
}

func TestResolveReleaseDelegatesToStore(t *testing.T) {
	t.Parallel()

	store := &fakeReleaseStore{publicRelease: &models.PublicRelease{Release: models.Release{Id: "release-id"}}}
	release, err := ResolveRelease(context.Background(), store, "artist-slug", "release-slug")
	if err != nil {
		t.Fatalf("ResolveRelease returned error: %v", err)
	}
	if release.Id != "release-id" {
		t.Fatalf("expected release-id, got %q", release.Id)
	}
	if store.artistSlug != "artist-slug" || store.releaseSlug != "release-slug" {
		t.Fatalf("store saw %q/%q", store.artistSlug, store.releaseSlug)
	}
}

func TestResolveReleaseBySlugDelegatesToStore(t *testing.T) {
	t.Parallel()

	store := &fakeReleaseStore{publicRelease: &models.PublicRelease{Release: models.Release{Id: "release-id"}}}
	release, err := ResolveReleaseBySlug(context.Background(), store, "release-slug")
	if err != nil {
		t.Fatalf("ResolveReleaseBySlug returned error: %v", err)
	}
	if release.Id != "release-id" {
		t.Fatalf("expected release-id, got %q", release.Id)
	}
	if store.releaseSlug != "release-slug" {
		t.Fatalf("store saw %q", store.releaseSlug)
	}
}

func TestResolveByIdDelegatesToStore(t *testing.T) {
	t.Parallel()

	store := &fakeReleaseStore{release: &models.Release{Id: "release-id"}}
	release, err := ResolveById(context.Background(), store, "release-id")
	if err != nil {
		t.Fatalf("ResolveById returned error: %v", err)
	}
	if release.Id != "release-id" {
		t.Fatalf("expected release-id, got %q", release.Id)
	}
}

var errSlugStore = errors.New("slug store error")

type fakeReleaseStore struct {
	artistSlug    string
	releaseSlug   string
	exists        bool
	existsErr     error
	publicRelease *models.PublicRelease
	release       *models.Release
}

func (f *fakeReleaseStore) GetPublishedByArtistSlugAndReleaseSlug(ctx context.Context, artistSlug, releaseSlug string) (*models.PublicRelease, error) {
	f.artistSlug = artistSlug
	f.releaseSlug = releaseSlug
	if f.publicRelease == nil {
		return nil, errSlugStore
	}

	return f.publicRelease, nil
}

func (f *fakeReleaseStore) GetPublishedByReleaseSlug(ctx context.Context, releaseSlug string) (*models.PublicRelease, error) {
	f.releaseSlug = releaseSlug
	if f.publicRelease == nil {
		return nil, errSlugStore
	}

	return f.publicRelease, nil
}

func (f *fakeReleaseStore) GetReleaseById(ctx context.Context, releaseID string) (*models.Release, error) {
	if f.release == nil {
		return nil, errSlugStore
	}

	return f.release, nil
}

func (f *fakeReleaseStore) ReleaseSlugExists(ctx context.Context, artistKey, value string) (bool, error) {
	return f.exists, f.existsErr
}
