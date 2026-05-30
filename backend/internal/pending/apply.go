package pending

import "github.com/brandon-v-seeters/go-silk-wave/internal/models"

func applyReleaseFields(release *models.Release, edit models.PendingReleaseEdit) {
	if edit.Title != nil {
		release.Title = *edit.Title
	}
	if edit.Slug != nil {
		release.Slug = *edit.Slug
	}
	if edit.Description != nil {
		release.Description = *edit.Description
	}
	if edit.ReleaseType != nil {
		release.ReleaseType = *edit.ReleaseType
	}
	if edit.Pricing != nil {
		release.Pricing = *edit.Pricing
	}
	if edit.Metadata != nil {
		release.Metadata = *edit.Metadata
	}
	if edit.Credits != nil {
		release.Credits = *edit.Credits
	}
	if edit.Assets != nil {
		release.Assets = *edit.Assets
	}
	if edit.Schedule != nil {
		release.Schedule = *edit.Schedule
	}
	if edit.PublishAt != nil {
		release.PublishAt = *edit.PublishAt
	}
	if edit.Cover != nil {
		release.Cover = *edit.Cover
	}
	if edit.License != nil {
		release.License = *edit.License
	}
	if edit.CustomLicense != nil {
		release.CustomLicense = edit.CustomLicense
	}
	if edit.IsExplicit != nil {
		release.IsExplicit = *edit.IsExplicit
	}
	if edit.DownloadEnabled != nil {
		release.DownloadEnabled = *edit.DownloadEnabled
	}
	if edit.StreamingEnabled != nil {
		release.StreamingEnabled = *edit.StreamingEnabled
	}
}

func applyTrackFields(tracks []models.Track, edits []models.PendingTrackEdit) {
	editsByKey := make(map[string]models.PendingTrackEdit, len(edits))
	for _, edit := range edits {
		editsByKey[edit.TrackKey] = edit
	}

	for i := range tracks {
		edit, ok := editsByKey[tracks[i].Key]
		if !ok {
			continue
		}

		applyTrackEdit(&tracks[i], edit)
	}
}

func applyTrackEdit(track *models.Track, edit models.PendingTrackEdit) {
	if edit.Title != nil {
		track.Title = *edit.Title
	}
	if edit.Order != nil {
		track.Order = *edit.Order
	}
	if edit.Duration != nil {
		track.Duration = *edit.Duration
	}
	if edit.DurationDisplay != nil {
		track.DurationDisplay = *edit.DurationDisplay
	}
	if edit.Metadata != nil {
		track.Metadata = *edit.Metadata
	}
	if edit.FeaturedArtists != nil {
		track.FeaturedArtists = *edit.FeaturedArtists
	}
	if edit.Writers != nil {
		track.Writers = *edit.Writers
	}
	if edit.IsExplicit != nil {
		track.IsExplicit = *edit.IsExplicit
	}
	if edit.StreamingEnabled != nil {
		track.StreamingEnabled = *edit.StreamingEnabled
	}
	if edit.DownloadEnabled != nil {
		track.DownloadEnabled = *edit.DownloadEnabled
	}
}
