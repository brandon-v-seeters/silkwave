package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/brandon-v-seeters/go-silk-wave/internal/database"
	"github.com/brandon-v-seeters/go-silk-wave/internal/models"
)

const (
	releasesCollection = "Releases"
	tracksCollection   = "Tracks"
)

var ErrDuplicateReleaseSlug = errors.New("release slug already exists for artist")

type ListPublishedReleasesParams struct {
	Limit       int
	Offset      int
	SearchQuery string
}

type ReleaseRepository struct {
	db *database.ArangoDB
}

func NewReleaseRepository(db *database.ArangoDB) *ReleaseRepository {
	return &ReleaseRepository{db: db}
}

func (r *ReleaseRepository) ListPublished(ctx context.Context, params ListPublishedReleasesParams) ([]models.ReleaseWithArtist, error) {
	bindVars := map[string]interface{}{
		"limit":     params.Limit,
		"now":       time.Now(),
		"offset":    params.Offset,
		"published": string(models.ReleaseStatusPublished),
	}

	if params.SearchQuery != "" {
		q := /*aql*/ `
			FOR release IN Releases
				FILTER release.status == @published
				FILTER release.publishAt == null OR release.publishAt <= @now
				FILTER (
					LOWER(release.title) LIKE LOWER(CONCAT("%", @searchQuery, "%")) OR
					LOWER(release.description) LIKE LOWER(CONCAT("%", @searchQuery, "%"))
				)
				LET artist = FIRST(
					FOR a IN Artists
						FILTER a._key == release.artistKey
						RETURN a
				)
				SORT release.publishAt DESC
				LIMIT @offset, @limit
				RETURN MERGE(release, { artist: artist })
		`
		bindVars["searchQuery"] = params.SearchQuery

		releases, err := database.Query[models.ReleaseWithArtist](ctx, r.db, q, bindVars)
		if err != nil {
			return nil, fmt.Errorf("list published releases: %w", err)
		}

		return releases, nil
	}

	q := /*aql*/ `
		FOR release IN Releases
			FILTER release.status == @published
			FILTER release.publishAt == null OR release.publishAt <= @now
			LET artist = FIRST(
				FOR a IN Artists
					FILTER a._key == release.artistKey
					RETURN a
			)
			SORT release.publishAt DESC
			LIMIT @offset, @limit
			RETURN MERGE(release, { artist: artist })
	`

	releases, err := database.Query[models.ReleaseWithArtist](ctx, r.db, q, bindVars)
	if err != nil {
		return nil, fmt.Errorf("list published releases: %w", err)
	}

	return releases, nil
}
