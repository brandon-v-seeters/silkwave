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
	ReleaseDate time.Time `json:"releaseDate" binding:"required"`
	Cover       string    `json:"cover,omitempty"`
	TrackCount  int       `json:"trackCount,omitempty"`
	Hash        string    `json:"hash,omitempty"`
	Genres      []string  `json:"genres,omitempty"`
	Tracks      []Track   `json:"tracks,omitempty"`
}

// UpdateReleaseRequest represents the request body for updating a release
type UpdateReleaseRequest struct {
	Title       string     `json:"title,omitempty"`
	Slug        string     `json:"slug,omitempty"`
	ArtistKey   string     `json:"artistKey,omitempty"`
	Description string     `json:"description,omitempty"`
	ReleaseDate *time.Time `json:"releaseDate,omitempty"`
	Cover       string     `json:"cover,omitempty"`
	TrackCount  *int       `json:"trackCount,omitempty"`
	Genres      []string   `json:"genres,omitempty"`
	Tracks      []Track    `json:"tracks,omitempty"`
}

type ConfirmDraftReleaseRequest struct {
	ArtistKey   string `json:"artistKey" binding:"required"`
	ReleaseHash string `json:"releaseHash" binding:"required"`
}

type ReleaseStatus string

const (
	ReleaseStatusDraft     ReleaseStatus = "draft"
	ReleaseStatusScheduled ReleaseStatus = "scheduled"
	ReleaseStatusPublished ReleaseStatus = "published"
	ReleaseStatusArchived  ReleaseStatus = "archived"
	ReleaseStatusDeleted   ReleaseStatus = "deleted"
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
	BasePath string         `json:"basePath"` // artist_content/{artistKey}/{hash}/
}

type ReleaseSchedule struct {
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
	ReleaseDate  *time.Time `json:"releaseDate"`
	PublishedAt  *time.Time `json:"publishedAt"`
	PreOrderDate *time.Time `json:"preOrderDate"`
}

type Release struct {
	DocumentMeta `tstype:",extends"`

	// Identifiers
	Hash string `json:"hash"` // Public URL identifier
	Slug string `json:"slug"` // SEO-friendly URL

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
	IsExplicit       bool      `json:"isExplicit"`
	IsFeatured       bool      `json:"isFeatured"`
	DownloadEnabled  bool      `json:"downloadEnabled"`
	StreamingEnabled bool      `json:"streamingEnabled"`
	ReleaseDate      time.Time `json:"releaseDate" binding:"required"`
	Cover            string    `json:"cover,omitempty"`
	Genres           []string  `json:"genres,omitempty"`
	Published        bool      `json:"published"`
	IsUploaded       bool      `json:"isUploaded"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

type CreateDraftRequest struct {
	ArtistKey string           `json:"artistKey"`
	Title     string           `json:"title"`
	Genres    []string         `json:"genres"`
	Tracks    []TrackRequest   `json:"tracks"`
	CoverArt  *CoverArtRequest `json:"coverArt"`
}

type TrackRequest struct {
	Title    string `json:"title"`
	FileName string `json:"fileName"`
	FileType string `json:"fileType"`
	FileSize int64  `json:"fileSize"`
	Duration string `json:"duration"`
	Order    int    `json:"order"`
}

type CoverArtRequest struct {
	FileName string `json:"fileName"`
	FileType string `json:"fileType"`
	FileSize int64  `json:"fileSize"`
}

type CreateDraftResponse struct {
	ReleaseHash   string           `json:"releaseHash"`
	ReleaseKey    string           `json:"releaseKey"`
	ArtistKey     string           `json:"artistKey"`
	PresignedUrls PresignedUrlsDTO `json:"presignedUrls"`
}

type PresignedUrlsDTO struct {
	CoverArt *string       `json:"coverArt"`
	Tracks   []TrackUrlDTO `json:"tracks"`
}

type TrackUrlDTO struct {
	Hash         string `json:"hash"`
	FileName     string `json:"fileName"`
	PresignedUrl string `json:"presignedUrl"`
	StoragePath  string `json:"storagePath"`
}

type ConfirmUploadsRequest struct {
	Tracks       []ConfirmedTrack `json:"tracks"`
	CoverArtPath *string          `json:"coverArtPath"`
}

type ConfirmedTrack struct {
	TrackID     string `json:"trackId"`
	StoragePath string `json:"storagePath"`
}

// ReleaseWithArtist represents a release with artist data joined
type ReleaseWithArtist struct {
	Release
	Artist *Artist `json:"artist,omitempty"`
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
