package models

import (
	"time"
)

// ReleaseType represents the type of release
type ReleaseType string

const (
	ReleaseTypeAlbum       ReleaseType = "album"
	ReleaseTypeEP          ReleaseType = "ep"
	ReleaseTypeSingle      ReleaseType = "single"
	ReleaseTypeCompilation ReleaseType = "compilation"
	ReleaseTypeRemix       ReleaseType = "remix"
	ReleaseTypeLive        ReleaseType = "live"
)

// CreateReleaseRequest represents the request body for creating a release
type CreateReleaseRequest struct {
	Key         string    `json:"key,omitempty"`
	Title       string    `json:"title" binding:"required"`
	Slug        string    `json:"slug" binding:"required"`
	ArtistKey   string    `json:"artistKey" binding:"required"`
	Description string    `json:"description,omitempty"`
	PublishAt   time.Time `json:"publishAt" binding:"required"`
	Cover       string    `json:"cover,omitempty"`
	TrackCount  int       `json:"trackCount,omitempty"`
	Id          string    `json:"id,omitempty"`
	Genres      []string  `json:"genres,omitempty"`
	Tracks      []Track   `json:"tracks,omitempty"`
}

// UpdateReleaseRequest represents the request body for updating a release
type UpdateReleaseRequest struct {
	Title       string     `json:"title,omitempty"`
	Slug        string     `json:"slug,omitempty"`
	ArtistKey   string     `json:"artistKey,omitempty"`
	Description string     `json:"description,omitempty"`
	PublishAt   *time.Time `json:"publishAt,omitempty"`
	Cover       string     `json:"cover,omitempty"`
	TrackCount  *int       `json:"trackCount,omitempty"`
	Genres      []string   `json:"genres,omitempty"`
	Tracks      []Track    `json:"tracks,omitempty"`
}

type ConfirmDraftReleaseRequest struct {
	ArtistKey string `json:"artistKey" binding:"required"`
	ReleaseId string `json:"releaseId" binding:"required"`
}

type ReleaseStatus string

const (
	ReleaseStatusDraft     ReleaseStatus = "draft"
	ReleaseStatusPublished ReleaseStatus = "published"
	ReleaseStatusArchived  ReleaseStatus = "archived"
)

type LicenseType string

const (
	LicenseAllRightsReserved LicenseType = "all_rights_reserved"
	LicenseCreativeCommons   LicenseType = "creative_commons"
	LicenseRoyaltyFree       LicenseType = "royalty_free"
	LicenseCustom            LicenseType = "custom"
)

type AudioFormat string

const (
	AudioFormatMP3_128 AudioFormat = "mp3_128"
	AudioFormatMP3_320 AudioFormat = "mp3_320"
	AudioFormatFLAC    AudioFormat = "flac"
	AudioFormatWAV     AudioFormat = "wav"
	AudioFormatAIFF    AudioFormat = "aiff"
)

// =============================================================================
// NESTED STRUCTS
// =============================================================================

type ReleasePricing struct {
	BasePrice    int  `json:"basePrice"`    // In cents (999 = $9.99)
	MinimumPrice *int `json:"minimumPrice"` // For "pay what you want"
	TrackPrice   *int `json:"trackPrice"`   // Individual track price
}

type ReleaseMetadata struct {
	Genres        []string  `json:"genres"`
	Tags          []string  `json:"tags"`
	Moods         []string  `json:"moods"`
	BPMRange      *BPMRange `json:"bpmRange"`
	Key           *string   `json:"key"`
	Label         *string   `json:"label"`
	CatalogNumber *string   `json:"catalogNumber"`
}

type BPMRange struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

type ArtistCredit struct {
	ArtistKey string `json:"artistKey"`
	Name      string `json:"name"`
	Role      string `json:"role"` // "primary", "featured", "remixer"
}

type Credit struct {
	Name string `json:"name"`
	Role string `json:"role"` // "Producer", "Mixing Engineer", etc.
}

type ReleaseCredits struct {
	Artists         []ArtistCredit `json:"artists"`
	Credits         []Credit       `json:"credits"`
	Writers         []string       `json:"writers"`
	CopyrightHolder string         `json:"copyrightHolder"`
	CopyrightYear   int            `json:"copyrightYear"`
}

type CoverArtAsset struct {
	Original  string  `json:"original"`
	Thumbnail *string `json:"thumbnail"` // 300x300
	Medium    *string `json:"medium"`    // 600x600
}

type ReleaseAssets struct {
	CoverArt *CoverArtAsset `json:"coverArt"`
	BasePath string         `json:"basePath"` // artist_content/{artistKey}/releases/{id}/
}

type ReleaseSchedule struct {
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
	PublishAt    *time.Time `json:"publishAt"`
	PublishedAt  *time.Time `json:"publishedAt"`
	PreOrderDate *time.Time `json:"preOrderDate"`
}

type Release struct {
	DocumentMeta `tstype:",extends"`

	// Identifiers
	Id   string `json:"id"`   // Stable external UUID for authenticated API paths.
	Slug string `json:"slug"` // SEO-friendly public URL slug.

	// Core info
	Title       string        `json:"title"`
	Description string        `json:"description"` // Supports markdown
	ReleaseType ReleaseType   `json:"releaseType"`
	Status      ReleaseStatus `json:"status"`

	// Ownership
	ArtistKey string `json:"artistKey"`

	// Nested data
	Pricing  ReleasePricing  `json:"pricing"`
	Metadata ReleaseMetadata `json:"metadata"`
	Credits  ReleaseCredits  `json:"credits"`
	Assets   ReleaseAssets   `json:"assets"`
	Schedule ReleaseSchedule `json:"schedule"`

	// Computed (updated when tracks change)
	TotalDuration    int           `json:"totalDuration"` // seconds
	TrackCount       int           `json:"trackCount"`
	AvailableFormats []AudioFormat `json:"availableFormats"`
	TrackList        []string      `json:"trackList,omitempty"`

	// Licensing
	License       LicenseType `json:"license"`
	CustomLicense *string     `json:"customLicense"`

	// Flags
	IsExplicit       bool                `json:"isExplicit"`
	IsFeatured       bool                `json:"isFeatured"`
	DownloadEnabled  bool                `json:"downloadEnabled"`
	StreamingEnabled bool                `json:"streamingEnabled"`
	Pending          *PendingReleaseEdit `json:"pending,omitempty"`
	PublishAt        time.Time           `json:"publishAt" binding:"required"`
	Cover            string              `json:"cover,omitempty"`
	Genres           []string            `json:"genres,omitempty"`
	CreatedAt        time.Time           `json:"createdAt"`
	UpdatedAt        time.Time           `json:"updatedAt"`
}

type CreateDraftRequest struct {
	ArtistKey string           `json:"artistKey" binding:"required"`
	Title     string           `json:"title" binding:"required"`
	Genres    []string         `json:"genres"`
	Tracks    []TrackRequest   `json:"tracks" binding:"required,dive"`
	CoverArt  *CoverArtRequest `json:"coverArt"`
}

type TrackRequest struct {
	Title    string `json:"title" binding:"required"`
	FileName string `json:"fileName" binding:"required"`
	FileType string `json:"fileType" binding:"required"`
	FileSize int64  `json:"fileSize" binding:"required"`
	Duration string `json:"duration"`
	Order    int    `json:"order"`
}

type CoverArtRequest struct {
	FileName string `json:"fileName" binding:"required"`
	FileType string `json:"fileType" binding:"required"`
	FileSize int64  `json:"fileSize" binding:"required"`
}

type CreateDraftResponse struct {
	ReleaseId     string           `json:"releaseId"`
	ReleaseKey    string           `json:"releaseKey"`
	ArtistKey     string           `json:"artistKey"`
	PresignedUrls PresignedUrlsDTO `json:"presignedUrls"`
}

type PresignedUrlsDTO struct {
	CoverArt *string       `json:"coverArt"`
	Tracks   []TrackUrlDTO `json:"tracks"`
}

type TrackUrlDTO struct {
	Id           string `json:"id"`
	FileName     string `json:"fileName"`
	PresignedUrl string `json:"presignedUrl"`
	StoragePath  string `json:"storagePath"`
}

type ConfirmUploadsRequest struct {
	Tracks       []ConfirmedTrack `json:"tracks" binding:"required,dive"`
	CoverArtPath *string          `json:"coverArtPath"`
}

type ConfirmedTrack struct {
	TrackID     string `json:"trackId" binding:"required"`
	StoragePath string `json:"storagePath" binding:"required"`
}

// ReleaseWithArtist represents a release with artist data joined
type ReleaseWithArtist struct {
	Release
	Artist *Artist `json:"artist,omitempty"`
}

type PublicRelease struct {
	Release `tstype:",extends"`
	Artist  *Artist       `json:"artist,omitempty"`
	Tracks  []PublicTrack `json:"tracks"`
}

// ReleasesResponse represents the response for listing releases
type ReleasesResponse struct {
	Releases []ReleaseWithArtist `json:"releases"`
	Limit    int                 `json:"limit"`
	Offset   int                 `json:"offset"`
}

type ReleaseWithTracks struct {
	Release `tstype:",extends"`
	Tracks  []Track `json:"tracks"`
}
