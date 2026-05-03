package storage

import (
	"path"
	"path/filepath"
	"strings"
)

// KeyResolver handles object key resolution for artist content storage in R2/S3.
// Object key structure:
//
//	artist_content/{artistKey}/releases/{releaseHash}/
//	  draft/                    <- exists only while unpublished
//	    cover.jpg
//	    wavs/
//	      01-track-title.wav
//	    mp3s/
//	      01-track-title.mp3
//	  cover.jpg                 <- after publishing (draft removed)
//	  wavs/
//	  mp3s/
//
//	avatars/{userKey}/
//	  avatar.png
type KeyResolver struct{}

// NewKeyResolver creates a new KeyResolver.
func NewKeyResolver() *KeyResolver {
	return &KeyResolver{}
}

// --- Avatar ---

// Avatar returns the object key for a user's avatar image.
// Key: avatars/{userKey}/avatar.{ext}
func (r *KeyResolver) Avatar(userKey, ext string) string {
	return path.Join("avatars", userKey, "avatar."+ext)
}

// AvatarPrefix returns the prefix for avatar files (for listing/deleting old avatars).
// Key: avatars/{userKey}/avatar
func (r *KeyResolver) AvatarPrefix(userKey string) string {
	return path.Join("avatars", userKey, "avatar")
}

// --- Artist Content Root ---

// ArtistContentPrefix returns the prefix for all of an artist's content.
// Key: artist_content/{artistKey}/
func (r *KeyResolver) ArtistContentPrefix(artistKey string) string {
	return path.Join("artist_content", artistKey) + "/"
}

// --- Release (Published) ---

// ReleasePrefix returns the prefix for a specific release.
// Key: artist_content/{artistKey}/releases/{releaseHash}/
func (r *KeyResolver) ReleasePrefix(artistKey, releaseHash string) string {
	return path.Join("artist_content", artistKey, "releases", releaseHash) + "/"
}

// ReleaseCover returns the key for a published release's cover image.
// Key: artist_content/{artistKey}/releases/{releaseHash}/cover.jpg
func (r *KeyResolver) ReleaseCover(artistKey, releaseHash, ext string) string {
	return path.Join("artist_content", artistKey, "releases", releaseHash, "cover."+ext)
}

// ReleaseCoverPrefix returns the prefix for cover files (to find cover with any extension).
// Key: artist_content/{artistKey}/releases/{releaseHash}/cover
func (r *KeyResolver) ReleaseCoverPrefix(artistKey, releaseHash string) string {
	return path.Join("artist_content", artistKey, "releases", releaseHash, "cover")
}

// ReleaseWAVsPrefix returns the prefix for WAV files in a published release.
// Key: artist_content/{artistKey}/releases/{releaseHash}/wavs/
func (r *KeyResolver) ReleaseWAVsPrefix(artistKey, releaseHash string) string {
	return path.Join("artist_content", artistKey, "releases", releaseHash, "wavs") + "/"
}

// ReleaseMP3sPrefix returns the prefix for MP3 files in a published release.
// Key: artist_content/{artistKey}/releases/{releaseHash}/mp3s/
func (r *KeyResolver) ReleaseMP3sPrefix(artistKey, releaseHash string) string {
	return path.Join("artist_content", artistKey, "releases", releaseHash, "mp3s") + "/"
}

// ReleaseTrack returns the key for a specific track file in a published release.
// Key: artist_content/{artistKey}/releases/{releaseHash}/{format}/{filename}
func (r *KeyResolver) ReleaseTrack(artistKey, releaseHash, format, filename string) string {
	return path.Join("artist_content", artistKey, "releases", releaseHash, format, filename)
}

// --- Release Draft ---

// ReleaseDraftPrefix returns the prefix for a draft release.
// Key: artist_content/{artistKey}/releases/{releaseHash}/draft/
func (r *KeyResolver) ReleaseDraftPrefix(artistKey, releaseHash string) string {
	return path.Join("artist_content", artistKey, "releases", releaseHash, "draft") + "/"
}

// ReleaseDraftCover returns the key for a draft release's cover image.
// Key: artist_content/{artistKey}/releases/{releaseHash}/draft/cover.{ext}
func (r *KeyResolver) ReleaseDraftCover(artistKey, releaseHash, ext string) string {
	return path.Join("artist_content", artistKey, "releases", releaseHash, "draft", "cover."+ext)
}

// ReleaseDraftCoverPrefix returns the prefix for draft cover files.
// Key: artist_content/{artistKey}/releases/{releaseHash}/draft/cover
func (r *KeyResolver) ReleaseDraftCoverPrefix(artistKey, releaseHash string) string {
	return path.Join("artist_content", artistKey, "releases", releaseHash, "draft", "cover")
}

// ReleaseDraftWAVsPrefix returns the prefix for WAV files in a draft release.
// Key: artist_content/{artistKey}/releases/{releaseHash}/draft/wavs/
func (r *KeyResolver) ReleaseDraftWAVsPrefix(artistKey, releaseHash string) string {
	return path.Join("artist_content", artistKey, "releases", releaseHash, "draft", "wavs") + "/"
}

// ReleaseDraftMP3sPrefix returns the prefix for MP3 files in a draft release.
// Key: artist_content/{artistKey}/releases/{releaseHash}/draft/mp3s/
func (r *KeyResolver) ReleaseDraftMP3sPrefix(artistKey, releaseHash string) string {
	return path.Join("artist_content", artistKey, "releases", releaseHash, "draft", "mp3s") + "/"
}

// ReleaseDraftTrack returns the key for a specific track file in a draft release.
// Key: artist_content/{artistKey}/releases/{releaseHash}/draft/{format}/{filename}
func (r *KeyResolver) ReleaseDraftTrack(artistKey, releaseHash, format, filename string) string {
	return path.Join("artist_content", artistKey, "releases", releaseHash, "draft", format, filename)
}

// --- Utility ---

// GetFormatFolder returns the folder name based on file type/mime type.
func (r *KeyResolver) GetFormatFolder(fileType string) string {
	fileType = strings.ToLower(fileType)
	switch {
	case strings.Contains(fileType, "wav"):
		return "wavs"
	case strings.Contains(fileType, "flac"):
		return "flacs"
	case strings.Contains(fileType, "aiff"):
		return "aiffs"
	default:
		return "mp3s"
	}
}

// SanitizeFileName cleans a filename for storage.
func (r *KeyResolver) SanitizeFileName(filename string) string {
	// Remove path, keep extension
	base := filepath.Base(filename)
	// Replace spaces with hyphens, lowercase
	clean := strings.ToLower(strings.ReplaceAll(base, " ", "-"))
	return clean
}

// BuildTrackStoragePath builds the full storage path for a track.
// Returns: artist_content/{artistKey}/releases/{releaseHash}/draft/{format}/{order:02d}-{sanitizedFilename}
func (r *KeyResolver) BuildTrackStoragePath(artistKey, releaseHash string, order int, fileName, fileType string) string {
	format := r.GetFormatFolder(fileType)
	cleanName := r.SanitizeFileName(fileName)
	storageName := strings.ToLower(cleanName)
	// Add order prefix if not already present
	if len(storageName) < 3 || storageName[2] != '-' {
		storageName = strings.ToLower(cleanName)
	}
	return r.ReleaseDraftTrack(artistKey, releaseHash, format, storageName)
}

// DraftToPublishedKey converts a draft object key to its published equivalent.
// Example: "artist_content/abc/releases/xyz/draft/cover.jpg" -> "artist_content/abc/releases/xyz/cover.jpg"
func (r *KeyResolver) DraftToPublishedKey(draftKey string) string {
	return strings.Replace(draftKey, "/draft/", "/", 1)
}

// PublishedToDraftKey converts a published object key to its draft equivalent.
// Example: "artist_content/abc/releases/xyz/cover.jpg" -> "artist_content/abc/releases/xyz/draft/cover.jpg"
func (r *KeyResolver) PublishedToDraftKey(publishedKey, artistKey, releaseHash string) string {
	releasePrefix := r.ReleasePrefix(artistKey, releaseHash)
	relativePath := strings.TrimPrefix(publishedKey, releasePrefix)
	return path.Join("artist_content", artistKey, "releases", releaseHash, "draft", relativePath)
}
