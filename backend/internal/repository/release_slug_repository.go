package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/brandon-v-seeters/go-silk-wave/internal/database"
	"github.com/brandon-v-seeters/go-silk-wave/internal/models"
)

func (r *ReleaseRepository) GetPublishedByArtistSlugAndReleaseSlug(ctx context.Context, artistSlug, releaseSlug string) (*models.PublicRelease, error) {
	q := /*aql*/ `
		FOR artist IN Artists
			FILTER artist.slug == @artistSlug
			FOR release IN Releases
			FILTER release.artistKey == artist._key
			FILTER release.slug == @releaseSlug
			FILTER release.status == @published
			FILTER release.publishAt == null OR release.publishAt <= @now
			LET tracks = (
				FOR track IN Tracks
					FILTER track.releaseKey == release._key
					SORT track.order ASC
					RETURN track
			)
			LIMIT 1
			RETURN MERGE(release, { artist: artist, tracks: tracks })
	`

	release, err := database.QueryOne[models.PublicRelease](ctx, r.db, q, map[string]interface{}{
		"artistSlug":  artistSlug,
		"now":         time.Now(),
		"published":   string(models.ReleaseStatusPublished),
		"releaseSlug": releaseSlug,
	})
	if err != nil {
		return nil, fmt.Errorf("get published release by artist slug and release slug: %w", err)
	}

	return release, nil
}

func (r *ReleaseRepository) GetReleaseById(ctx context.Context, releaseId string) (*models.Release, error) {
	q := /*aql*/ `
		FOR release IN Releases
			FILTER release.id == @id
			LIMIT 1
			RETURN release
	`

	release, err := database.QueryOne[models.Release](ctx, r.db, q, map[string]interface{}{
		"id": releaseId,
	})
	if err != nil {
		return nil, fmt.Errorf("get release by id: %w", err)
	}

	return release, nil
}

func (r *ReleaseRepository) ReleaseSlugExists(ctx context.Context, artistKey, value string) (bool, error) {
	q := /*aql*/ `
		FOR release IN Releases
			FILTER release.artistKey == @artistKey
			FILTER release.slug == @slug
			LIMIT 1
			RETURN true
	`

	exists, err := database.Query[bool](ctx, r.db, q, map[string]interface{}{
		"artistKey": artistKey,
		"slug":      value,
	})
	if err != nil {
		return false, fmt.Errorf("release slug exists: %w", err)
	}

	return len(exists) > 0 && exists[0], nil
}
