package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/brandon-v-seeters/go-silk-wave/internal/database"
	"github.com/brandon-v-seeters/go-silk-wave/internal/models"
	"github.com/brandon-v-seeters/go-silk-wave/internal/slug"
)

func (r *ReleaseRepository) ListDraftsByArtistKey(ctx context.Context, artistKey string) ([]models.ReleaseWithArtist, error) {
	q := /*aql*/ `
		FOR release IN Releases
			FILTER release.artistKey == @artistKey
			FILTER release.status == @draft
			RETURN release
	`

	drafts, err := database.Query[models.ReleaseWithArtist](ctx, r.db, q, map[string]interface{}{
		"artistKey": artistKey,
		"draft":     string(models.ReleaseStatusDraft),
	})
	if err != nil {
		return nil, fmt.Errorf("list drafts by artist key: %w", err)
	}

	return drafts, nil
}

func (r *ReleaseRepository) CreateDraftWithTracks(ctx context.Context, release models.Release, tracks []models.Track) (string, error) {
	var releaseKey string

	if err := slug.EnsureUnique(ctx, r, release.ArtistKey, release.Slug); err != nil {
		if errors.Is(err, slug.ErrNotUnique) {
			return "", ErrDuplicateReleaseSlug
		}

		return "", err
	}

	err := database.WithTransaction(ctx, r.db,
		[]string{releasesCollection},
		[]string{releasesCollection, tracksCollection},
		func(tx *database.Transaction) error {
			releaseKeyPtr, err := database.TxQueryOne[string](
				ctx,
				tx,
				/*aql*/ `
					INSERT @release
					IN Releases
					RETURN NEW._key
				`,
				map[string]interface{}{"release": release},
			)
			if err != nil {
				return fmt.Errorf("create draft release: %w", err)
			}

			releaseKey = *releaseKeyPtr
			for i := range tracks {
				tracks[i].ReleaseKey = releaseKey
			}

			if len(tracks) == 0 {
				return nil
			}

			err = database.TxExec(
				ctx,
				tx,
				/*aql*/ `
					FOR track IN @tracks
						INSERT track IN Tracks
				`,
				map[string]interface{}{"tracks": tracks},
			)
			if err != nil {
				return fmt.Errorf("create draft tracks: %w", err)
			}

			return nil
		})
	if err != nil {
		return "", fmt.Errorf("create draft with tracks: %w", err)
	}

	return releaseKey, nil
}

func (r *ReleaseRepository) GetDraftWithTracksByKey(ctx context.Context, releaseKey string) (*models.ReleaseWithTracks, error) {
	q := /*aql*/ `
		FOR draft IN Releases
			FILTER draft._key == @key
			FILTER draft.status == @draft
			LET tracks = (
				FOR track IN Tracks
					FILTER track.releaseKey == draft._key
					SORT track.order ASC
					RETURN track
			)
			LIMIT 1
			RETURN MERGE(draft, { tracks: tracks })
	`

	draft, err := database.QueryOne[models.ReleaseWithTracks](ctx, r.db, q, map[string]interface{}{
		"draft": string(models.ReleaseStatusDraft),
		"key":   releaseKey,
	})
	if err != nil {
		return nil, fmt.Errorf("get draft with tracks by key: %w", err)
	}

	return draft, nil
}
