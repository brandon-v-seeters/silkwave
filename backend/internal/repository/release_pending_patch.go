package repository

import (
	"time"

	"github.com/brandon-v-seeters/go-silk-wave/internal/models"
)

func releasePatch(release models.Release) map[string]interface{} {
	return map[string]interface{}{
		"title":            release.Title,
		"slug":             release.Slug,
		"description":      release.Description,
		"releaseType":      release.ReleaseType,
		"pricing":          release.Pricing,
		"metadata":         release.Metadata,
		"credits":          release.Credits,
		"assets":           release.Assets,
		"schedule":         release.Schedule,
		"publishAt":        release.PublishAt,
		"cover":            release.Cover,
		"license":          release.License,
		"customLicense":    release.CustomLicense,
		"isExplicit":       release.IsExplicit,
		"downloadEnabled":  release.DownloadEnabled,
		"streamingEnabled": release.StreamingEnabled,
		"pending":          nil,
		"updatedAt":        release.UpdatedAt,
	}
}

func trackPatches(tracks []models.Track) []map[string]interface{} {
	now := time.Now()
	patches := make([]map[string]interface{}, 0, len(tracks))
	for _, track := range tracks {
		patches = append(patches, map[string]interface{}{
			"_key":             track.Key,
			"title":            track.Title,
			"order":            track.Order,
			"duration":         track.Duration,
			"durationDisplay":  track.DurationDisplay,
			"metadata":         track.Metadata,
			"featuredArtists":  track.FeaturedArtists,
			"writers":          track.Writers,
			"isExplicit":       track.IsExplicit,
			"streamingEnabled": track.StreamingEnabled,
			"downloadEnabled":  track.DownloadEnabled,
			"updatedAt":        now,
		})
	}

	return patches
}
