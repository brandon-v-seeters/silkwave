package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/brandon-v-seeters/go-silk-wave/internal/database"
	"github.com/brandon-v-seeters/go-silk-wave/internal/models"
)

func (r *ReleaseRepository) MarkUploadsComplete(ctx context.Context, releaseId string, trackIds []string) error {
	q := /*aql*/ `
		FOR release IN Releases
			FILTER release.id == @releaseId
			FOR track IN Tracks
				FILTER track.releaseKey == release._key
				FILTER track.id IN @ids
				UPDATE track WITH { uploaded: true, updatedAt: @now } IN Tracks
				RETURN NEW._key
	`

	if _, err := database.Query[string](ctx, r.db, q, map[string]interface{}{
		"ids":       trackIds,
		"now":       time.Now(),
		"releaseId": releaseId,
	}); err != nil {
		return fmt.Errorf("mark track uploads complete: %w", err)
	}

	return nil
}

func (r *ReleaseRepository) UploadsComplete(ctx context.Context, releaseKey string) (bool, error) {
	q := /*aql*/ `
		LET total = (
			FOR track IN Tracks
				FILTER track.releaseKey == @releaseKey
				RETURN 1
		)
		LET incomplete = (
			FOR track IN Tracks
				FILTER track.releaseKey == @releaseKey
				FILTER track.uploaded != true
				RETURN 1
		)
		RETURN LENGTH(total) > 0 && LENGTH(incomplete) == 0
	`

	complete, err := database.QueryOne[bool](ctx, r.db, q, map[string]interface{}{
		"releaseKey": releaseKey,
	})
	if err != nil {
		return false, fmt.Errorf("check uploads complete: %w", err)
	}
	if complete == nil {
		return false, nil
	}

	return *complete, nil
}

func (r *ReleaseRepository) UpdateReleaseStatus(ctx context.Context, releaseId string, status models.ReleaseStatus) error {
	q := /*aql*/ `
		FOR release IN Releases
			FILTER release.id == @id
			UPDATE release WITH { status: @status, updatedAt: @now } IN Releases
			RETURN NEW._key
	`

	_, err := database.QueryOne[string](ctx, r.db, q, map[string]interface{}{
		"id":     releaseId,
		"now":    time.Now(),
		"status": string(status),
	})
	if err != nil {
		return fmt.Errorf("update release status: %w", err)
	}

	return nil
}

func (r *ReleaseRepository) DeleteReleaseWithTracks(ctx context.Context, releaseKey, releaseId string) error {
	return database.WithTransaction(ctx, r.db,
		[]string{releasesCollection, tracksCollection},
		[]string{releasesCollection, tracksCollection},
		func(tx *database.Transaction) error {
			if err := database.TxExec(
				ctx,
				tx,
				/*aql*/ `
					FOR track IN Tracks
						FILTER track.releaseKey == @key
						REMOVE track IN Tracks
				`,
				map[string]interface{}{"key": releaseKey},
			); err != nil {
				return fmt.Errorf("delete release tracks: %w", err)
			}

			if err := database.TxExec(
				ctx,
				tx,
				/*aql*/ `
					FOR release IN Releases
						FILTER release.id == @id
						REMOVE release IN Releases
				`,
				map[string]interface{}{"id": releaseId},
			); err != nil {
				return fmt.Errorf("delete release: %w", err)
			}

			return nil
		})
}
