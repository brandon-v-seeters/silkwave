package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/brandon-v-seeters/go-silk-wave/internal/database"
	"github.com/brandon-v-seeters/go-silk-wave/internal/models"
	"github.com/brandon-v-seeters/go-silk-wave/internal/pending"
)

var (
	ErrReleaseMustBePublished = errors.New("release must be published")
	ErrNoPendingEdit          = errors.New("release has no pending edit")
)

func (r *ReleaseRepository) GetReleaseWithTracksById(ctx context.Context, releaseId string) (*models.ReleaseWithTracks, error) {
	q := /*aql*/ `
		FOR release IN Releases
			FILTER release.id == @id
			LET tracks = (
				FOR track IN Tracks
					FILTER track.releaseKey == release._key
					SORT track.order ASC
					RETURN track
			)
			LIMIT 1
			RETURN MERGE(release, { tracks: tracks })
	`

	release, err := database.QueryOne[models.ReleaseWithTracks](ctx, r.db, q, map[string]interface{}{
		"id": releaseId,
	})
	if err != nil {
		return nil, fmt.Errorf("get release with tracks by id: %w", err)
	}

	return release, nil
}

func (r *ReleaseRepository) StagePendingEdit(ctx context.Context, releaseId string, edit models.PendingReleaseEdit) (*models.ReleaseWithTracks, error) {
	release, err := r.GetReleaseWithTracksById(ctx, releaseId)
	if err != nil {
		return nil, err
	}
	if release.Status != models.ReleaseStatusPublished {
		return nil, ErrReleaseMustBePublished
	}

	stagedRelease, err := pending.Stage(release.Release, release.Tracks, edit)
	if err != nil {
		return nil, err
	}

	if err := r.updatePendingEdit(ctx, releaseId, stagedRelease.Pending); err != nil {
		return nil, err
	}

	release.Pending = stagedRelease.Pending
	preview, err := pending.Apply(*release, edit)
	if err != nil {
		return nil, err
	}

	return &preview, nil
}

func (r *ReleaseRepository) DiscardPendingEdit(ctx context.Context, releaseId string) error {
	return r.updatePendingEdit(ctx, releaseId, nil)
}

func (r *ReleaseRepository) PendingPreview(ctx context.Context, releaseId string) (*models.ReleaseWithTracks, error) {
	release, err := r.GetReleaseWithTracksById(ctx, releaseId)
	if err != nil {
		return nil, err
	}
	if release.Pending == nil {
		return release, nil
	}

	preview, err := pending.Apply(*release, *release.Pending)
	if err != nil {
		return nil, err
	}

	return &preview, nil
}

func (r *ReleaseRepository) PublishPendingEdit(ctx context.Context, releaseId string) (*models.ReleaseWithTracks, error) {
	release, err := r.GetReleaseWithTracksById(ctx, releaseId)
	if err != nil {
		return nil, err
	}
	if release.Status != models.ReleaseStatusPublished {
		return nil, ErrReleaseMustBePublished
	}
	if release.Pending == nil {
		return nil, ErrNoPendingEdit
	}

	published, err := pending.PublishOnto(*release, *release.Pending)
	if err != nil {
		return nil, err
	}
	published.UpdatedAt = time.Now()

	if err := r.updateReleaseAndTracks(ctx, published); err != nil {
		return nil, err
	}

	return &published, nil
}

func (r *ReleaseRepository) updatePendingEdit(ctx context.Context, releaseId string, edit *models.PendingReleaseEdit) error {
	q := /*aql*/ `
		FOR release IN Releases
			FILTER release.id == @id
			UPDATE release WITH { pending: @pending, updatedAt: @now } IN Releases
			RETURN NEW._key
	`

	if _, err := database.QueryOne[string](ctx, r.db, q, map[string]interface{}{
		"id":      releaseId,
		"now":     time.Now(),
		"pending": edit,
	}); err != nil {
		return fmt.Errorf("update pending edit: %w", err)
	}

	return nil
}

func (r *ReleaseRepository) updateReleaseAndTracks(ctx context.Context, release models.ReleaseWithTracks) error {
	return database.WithTransaction(ctx, r.db,
		[]string{releasesCollection, tracksCollection},
		[]string{releasesCollection, tracksCollection},
		func(tx *database.Transaction) error {
			if err := updateReleaseInTransaction(ctx, tx, release.Release); err != nil {
				return err
			}
			return updateTracksInTransaction(ctx, tx, release.Tracks)
		})
}

func updateReleaseInTransaction(ctx context.Context, tx *database.Transaction, release models.Release) error {
	_, err := database.TxQueryOne[string](
		ctx,
		tx,
		/*aql*/ `
			FOR existing IN Releases
				FILTER existing._key == @key
				UPDATE existing WITH @patch IN Releases
				RETURN NEW._key
		`,
		map[string]interface{}{
			"key":   release.Key,
			"patch": releasePatch(release),
		},
	)
	if err != nil {
		return fmt.Errorf("publish pending release fields: %w", err)
	}

	return nil
}

func updateTracksInTransaction(ctx context.Context, tx *database.Transaction, tracks []models.Track) error {
	err := database.TxExec(
		ctx,
		tx,
		/*aql*/ `
			FOR patch IN @tracks
				FOR track IN Tracks
					FILTER track._key == patch._key
					UPDATE track WITH UNSET(patch, "_key") IN Tracks
		`,
		map[string]interface{}{"tracks": trackPatches(tracks)},
	)
	if err != nil {
		return fmt.Errorf("publish pending track fields: %w", err)
	}

	return nil
}
