package models

import "time"

type OriginalFile struct {
	Path       string `json:"path"`
	Format     string `json:"format"`
	FileSize   int64  `json:"fileSize"`   // bytes
	Bitrate    *int   `json:"bitrate"`    // kbps
	SampleRate *int   `json:"sampleRate"` // Hz
}

type TranscodedFile struct {
	Format   AudioFormat `json:"format"`
	Path     string      `json:"path"`
	FileSize int64       `json:"fileSize"`
	Bitrate  int         `json:"bitrate"`
}

type TrackFiles struct {
	Original     OriginalFile     `json:"original"`
	Transcoded   []TranscodedFile `json:"transcoded"`
	WaveformPath *string          `json:"waveformPath"`
	PreviewPath  *string          `json:"previewPath"` // 30s preview
}

type TrackMetadata struct {
	BPM    *int     `json:"bpm"`
	Key    *string  `json:"key"`    // "Am", "C#", "7A"
	Genres []string `json:"genres"` // Inherits from album if empty
	Tags   []string `json:"tags"`
	ISRC   *string  `json:"isrc"` // International Standard Recording Code
}

type FeaturedArtist struct {
	ArtistKey string `json:"artistKey"`
	Name      string `json:"name"`
}

type Track struct {
	DocumentMeta `tstype:",extends"`

	// Identifiers
	Hash string `json:"hash"`

	// Relations
	ReleaseKey string `json:"releaseKey"`

	// Core info
	Title           string `json:"title"`
	Order           int    `json:"order"`
	Duration        int    `json:"duration"`        // seconds
	DurationDisplay string `json:"durationDisplay"` // "3:45"

	// Files
	Files TrackFiles `json:"files"`

	// Metadata
	Metadata TrackMetadata `json:"metadata"`

	// Credits
	FeaturedArtists []FeaturedArtist `json:"featuredArtists"`
	Writers         []string         `json:"writers"` // Inherits from album if empty

	// Status
	Uploaded            bool `json:"uploaded"`
	TranscodingComplete bool `json:"transcodingComplete"`

	// Timestamps
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	// Flags
	IsExplicit       bool `json:"isExplicit"`
	StreamingEnabled bool `json:"streamingEnabled"`
	DownloadEnabled  bool `json:"downloadEnabled"`

	// Stats (consider separate collection for frequent updates)
	PlayCount     int `json:"playCount"`
	DownloadCount int `json:"downloadCount"`
}
