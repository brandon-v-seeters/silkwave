package models

import "time"

type PendingReleaseEdit struct {
	Title            *string            `json:"title,omitempty"`
	Slug             *string            `json:"slug,omitempty"`
	Description      *string            `json:"description,omitempty"`
	ReleaseType      *ReleaseType       `json:"releaseType,omitempty"`
	Pricing          *ReleasePricing    `json:"pricing,omitempty"`
	Metadata         *ReleaseMetadata   `json:"metadata,omitempty"`
	Credits          *ReleaseCredits    `json:"credits,omitempty"`
	Assets           *ReleaseAssets     `json:"assets,omitempty"`
	Schedule         *ReleaseSchedule   `json:"schedule,omitempty"`
	PublishAt        *time.Time         `json:"publishAt,omitempty"`
	Cover            *string            `json:"cover,omitempty"`
	License          *LicenseType       `json:"license,omitempty"`
	CustomLicense    *string            `json:"customLicense,omitempty"`
	IsExplicit       *bool              `json:"isExplicit,omitempty"`
	DownloadEnabled  *bool              `json:"downloadEnabled,omitempty"`
	StreamingEnabled *bool              `json:"streamingEnabled,omitempty"`
	TrackCount       *int               `json:"trackCount,omitempty"`
	TrackList        *[]string          `json:"trackList,omitempty"`
	Tracks           []PendingTrackEdit `json:"tracks,omitempty"`
}

type PendingTrackEdit struct {
	TrackKey         string            `json:"_key"`
	Title            *string           `json:"title,omitempty"`
	Order            *int              `json:"order,omitempty"`
	Duration         *int              `json:"duration,omitempty"`
	DurationDisplay  *string           `json:"durationDisplay,omitempty"`
	Metadata         *TrackMetadata    `json:"metadata,omitempty"`
	FeaturedArtists  *[]FeaturedArtist `json:"featuredArtists,omitempty"`
	Writers          *[]string         `json:"writers,omitempty"`
	IsExplicit       *bool             `json:"isExplicit,omitempty"`
	StreamingEnabled *bool             `json:"streamingEnabled,omitempty"`
	DownloadEnabled  *bool             `json:"downloadEnabled,omitempty"`
	Files            *TrackFiles       `json:"files,omitempty"`
	Remove           *bool             `json:"remove,omitempty"`
}
