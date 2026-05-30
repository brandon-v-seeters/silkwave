package slug

import (
	"context"
	"errors"
	"fmt"

	"github.com/avelino/slugify"
	"github.com/brandon-v-seeters/go-silk-wave/internal/models"
)

var ErrNotUnique = errors.New("release slug already exists for artist")

type ReleaseStore interface {
	GetPublishedByArtistSlugAndReleaseSlug(ctx context.Context, artistSlug, releaseSlug string) (*models.PublicRelease, error)
	GetReleaseById(ctx context.Context, releaseID string) (*models.Release, error)
	ReleaseSlugExists(ctx context.Context, artistKey, value string) (bool, error)
}

func ResolveRelease(ctx context.Context, store ReleaseStore, artistSlug, releaseSlug string) (*models.PublicRelease, error) {
	release, err := store.GetPublishedByArtistSlugAndReleaseSlug(ctx, artistSlug, releaseSlug)
	if err != nil {
		return nil, fmt.Errorf("resolve release slug: %w", err)
	}

	return release, nil
}

func ResolveById(ctx context.Context, store ReleaseStore, releaseID string) (*models.Release, error) {
	release, err := store.GetReleaseById(ctx, releaseID)
	if err != nil {
		return nil, fmt.Errorf("resolve release id: %w", err)
	}

	return release, nil
}

func EnsureUnique(ctx context.Context, store ReleaseStore, artistKey, value string) error {
	exists, err := store.ReleaseSlugExists(ctx, artistKey, value)
	if err != nil {
		return fmt.Errorf("check release slug uniqueness: %w", err)
	}
	if exists {
		return ErrNotUnique
	}

	return nil
}

func Slugify(title string) string {
	return slugify.Slugify(title)
}
