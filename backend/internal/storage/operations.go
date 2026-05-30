package storage

import (
	"context"
	"fmt"
	"io"
	"path"
	"strings"
)

// --- Avatar Operations ---

// UploadAvatar uploads a user's avatar image, deleting any existing avatar first.
func (c *Client) UploadAvatar(ctx context.Context, userKey, ext string, body io.Reader, contentType string) error {
	// Delete existing avatars with any extension
	if err := c.DeleteAvatars(ctx, userKey); err != nil {
		return fmt.Errorf("failed to delete existing avatars: %w", err)
	}

	key := c.resolver.Avatar(userKey, ext)
	return c.Upload(ctx, key, body, contentType)
}

// DeleteAvatars removes all avatar files for a user (handles different extensions).
func (c *Client) DeleteAvatars(ctx context.Context, userKey string) error {
	prefix := c.resolver.AvatarPrefix(userKey)
	keys, err := c.List(ctx, prefix)
	if err != nil {
		return err
	}

	for _, key := range keys {
		if err := c.Delete(ctx, key); err != nil {
			return err
		}
	}
	return nil
}

// GetAvatarKey finds the avatar key for a user (since extension may vary).
func (c *Client) GetAvatarKey(ctx context.Context, userKey string) (string, error) {
	prefix := c.resolver.AvatarPrefix(userKey)
	keys, err := c.List(ctx, prefix)
	if err != nil {
		return "", err
	}

	if len(keys) == 0 {
		return "", nil
	}

	return keys[0], nil
}

// --- Release Draft Operations ---

// UploadDraftCover uploads a cover image for a release draft.
func (c *Client) UploadDraftCover(ctx context.Context, artistKey, releaseId, ext string, body io.Reader, contentType string) error {
	key := c.resolver.ReleaseDraftCover(artistKey, releaseId, ext)
	return c.Upload(ctx, key, body, contentType)
}

// UploadDraftTrack uploads a track file to a release draft.
func (c *Client) UploadDraftTrack(ctx context.Context, artistKey, releaseId, format, filename string, body io.Reader, contentType string) error {
	key := c.resolver.ReleaseDraftTrack(artistKey, releaseId, format, filename)
	return c.Upload(ctx, key, body, contentType)
}

// DeleteDraftFile deletes a specific file from a release draft.
func (c *Client) DeleteDraftFile(ctx context.Context, artistKey, releaseId, storagePath string) error {
	return c.Delete(ctx, storagePath)
}

// ListDraftFiles returns all files in a release draft.
func (c *Client) ListDraftFiles(ctx context.Context, artistKey, releaseId string) ([]string, error) {
	prefix := c.resolver.ReleaseDraftPrefix(artistKey, releaseId)
	return c.List(ctx, prefix)
}

// --- Release Publish Operations ---

// PublishRelease moves all content from the draft folder to the release root.
func (c *Client) PublishRelease(ctx context.Context, artistKey, releaseId string) error {
	draftPrefix := c.resolver.ReleaseDraftPrefix(artistKey, releaseId)
	releasePrefix := c.resolver.ReleasePrefix(artistKey, releaseId)

	// List all objects in the draft folder
	draftKeys, err := c.List(ctx, draftPrefix)
	if err != nil {
		return fmt.Errorf("failed to list draft files: %w", err)
	}

	if len(draftKeys) == 0 {
		return fmt.Errorf("no draft files found for release %s", releaseId)
	}

	// Move each file from draft to published location
	for _, draftKey := range draftKeys {
		// Calculate the published key by removing "draft/" from the path
		relativePath := strings.TrimPrefix(draftKey, draftPrefix)
		publishedKey := releasePrefix + relativePath

		if err := c.Move(ctx, draftKey, publishedKey); err != nil {
			return fmt.Errorf("failed to move %s to %s: %w", draftKey, publishedKey, err)
		}
	}

	return nil
}

// UnpublishRelease moves all published content back to the draft folder.
func (c *Client) UnpublishRelease(ctx context.Context, artistKey, releaseId string) error {
	releasePrefix := c.resolver.ReleasePrefix(artistKey, releaseId)
	draftPrefix := c.resolver.ReleaseDraftPrefix(artistKey, releaseId)

	// List all objects in the release folder (excluding draft/)
	allKeys, err := c.List(ctx, releasePrefix)
	if err != nil {
		return fmt.Errorf("failed to list release files: %w", err)
	}

	// Filter out keys that are already in draft/
	var publishedKeys []string
	for _, key := range allKeys {
		if !strings.HasPrefix(key, draftPrefix) {
			publishedKeys = append(publishedKeys, key)
		}
	}

	if len(publishedKeys) == 0 {
		return fmt.Errorf("no published files found for release %s", releaseId)
	}

	// Move each file from published to draft location
	for _, publishedKey := range publishedKeys {
		relativePath := strings.TrimPrefix(publishedKey, releasePrefix)
		draftKey := draftPrefix + relativePath

		if err := c.Move(ctx, publishedKey, draftKey); err != nil {
			return fmt.Errorf("failed to move %s to %s: %w", publishedKey, draftKey, err)
		}
	}

	return nil
}

// IsDraft checks if a release has draft content.
func (c *Client) IsDraft(ctx context.Context, artistKey, releaseId string) (bool, error) {
	draftPrefix := c.resolver.ReleaseDraftPrefix(artistKey, releaseId)
	keys, err := c.List(ctx, draftPrefix)
	if err != nil {
		return false, err
	}
	return len(keys) > 0, nil
}

// VerifyDraftUploads checks if all expected files exist in the draft folder.
// Returns a list of missing file paths.
func (c *Client) VerifyDraftUploads(ctx context.Context, artistKey, releaseId string, expectedPaths []string) ([]string, error) {
	var missingFiles []string

	for _, expectedPath := range expectedPaths {
		exists, err := c.Exists(ctx, expectedPath)
		if err != nil {
			return nil, fmt.Errorf("failed to check file %s: %w", expectedPath, err)
		}
		if !exists {
			missingFiles = append(missingFiles, expectedPath)
		}
	}

	return missingFiles, nil
}

// --- Release Operations ---

// UploadReleaseCover uploads a cover image for a published release.
func (c *Client) UploadReleaseCover(ctx context.Context, artistKey, releaseId, ext string, body io.Reader, contentType string) error {
	key := c.resolver.ReleaseCover(artistKey, releaseId, ext)
	return c.Upload(ctx, key, body, contentType)
}

// DeleteRelease removes all content for a release (both draft and published).
func (c *Client) DeleteRelease(ctx context.Context, artistKey, releaseId string) error {
	prefix := c.resolver.ReleasePrefix(artistKey, releaseId)
	return c.DeletePrefix(ctx, prefix)
}

// ListReleaseFiles returns all files for a release (both draft and published).
func (c *Client) ListReleaseFiles(ctx context.Context, artistKey, releaseId string) ([]string, error) {
	prefix := c.resolver.ReleasePrefix(artistKey, releaseId)
	return c.List(ctx, prefix)
}

// --- Artist Operations ---

// DeleteArtistContent removes all content for an artist.
func (c *Client) DeleteArtistContent(ctx context.Context, artistKey string) error {
	prefix := c.resolver.ArtistContentPrefix(artistKey)
	return c.DeletePrefix(ctx, prefix)
}

// ListArtistContent returns all files for an artist.
func (c *Client) ListArtistContent(ctx context.Context, artistKey string) ([]string, error) {
	prefix := c.resolver.ArtistContentPrefix(artistKey)
	return c.List(ctx, prefix)
}

// --- Presigned URL Helpers ---

// GetAvatarURL generates a presigned URL for downloading a user's avatar.
func (c *Client) GetAvatarURL(ctx context.Context, userKey string, expireSeconds int64) (string, error) {
	key, err := c.GetAvatarKey(ctx, userKey)
	if err != nil {
		return "", err
	}
	if key == "" {
		return "", nil
	}
	return c.GetPresignedURL(ctx, key, expireSeconds)
}

// GetDraftCoverURL generates a presigned URL for a draft release's cover.
func (c *Client) GetDraftCoverURL(ctx context.Context, artistKey, releaseId string, expireSeconds int64) (string, error) {
	// Find cover with any extension
	prefix := c.resolver.ReleaseDraftCoverPrefix(artistKey, releaseId)
	keys, err := c.List(ctx, prefix)
	if err != nil {
		return "", err
	}
	if len(keys) == 0 {
		return "", nil
	}
	return c.GetPresignedURL(ctx, keys[0], expireSeconds)
}

// GetReleaseCoverURL generates a presigned URL for a published release's cover.
func (c *Client) GetReleaseCoverURL(ctx context.Context, artistKey, releaseId string, expireSeconds int64) (string, error) {
	prefix := c.resolver.ReleaseCoverPrefix(artistKey, releaseId)
	keys, err := c.List(ctx, prefix)
	if err != nil {
		return "", err
	}
	if len(keys) == 0 {
		return "", nil
	}
	return c.GetPresignedURL(ctx, keys[0], expireSeconds)
}

// GetDraftTrackURL generates a presigned URL for a track in a draft release.
func (c *Client) GetDraftTrackURL(ctx context.Context, artistKey, releaseId, storagePath string, expireSeconds int64) (string, error) {
	return c.GetPresignedURL(ctx, storagePath, expireSeconds)
}

// GetReleaseTrackURL generates a presigned URL for a track in a published release.
func (c *Client) GetReleaseTrackURL(ctx context.Context, artistKey, releaseId, filename string, expireSeconds int64) (string, error) {
	ext := strings.ToLower(path.Ext(filename))
	format := "mp3s"
	if ext == ".wav" {
		format = "wavs"
	} else if ext == ".flac" {
		format = "flacs"
	}
	key := c.resolver.ReleaseTrack(artistKey, releaseId, format, filename)
	return c.GetPresignedURL(ctx, key, expireSeconds)
}

// GetPublishedReleaseObjectURL generates a presigned URL for a stored object path on a published release.
func (c *Client) GetPublishedReleaseObjectURL(ctx context.Context, artistKey, releaseId, storagePath string, expireSeconds int64) (string, error) {
	if storagePath == "" {
		return "", nil
	}

	key := c.resolver.DraftToPublishedKey(storagePath)
	releasePrefix := c.resolver.ReleasePrefix(artistKey, releaseId)
	if !strings.HasPrefix(key, releasePrefix) || key == releasePrefix {
		return "", fmt.Errorf("object path %q is outside release %s", storagePath, releaseId)
	}
	if strings.Contains(key, "/draft/") {
		return "", fmt.Errorf("object path %q still points to draft content", storagePath)
	}

	return c.GetPresignedURL(ctx, key, expireSeconds)
}

// GetPublishedTrackURL generates a presigned URL for a stored track path on a published release.
func (c *Client) GetPublishedTrackURL(ctx context.Context, artistKey, releaseId, storagePath string, expireSeconds int64) (string, error) {
	return c.GetPublishedReleaseObjectURL(ctx, artistKey, releaseId, storagePath, expireSeconds)
}

// GetUploadAvatarURL generates a presigned URL for uploading an avatar.
func (c *Client) GetUploadAvatarURL(ctx context.Context, userKey, ext, contentType string, expireSeconds int64) (string, error) {
	key := c.resolver.Avatar(userKey, ext)
	return c.GetPresignedUploadURL(ctx, key, contentType, expireSeconds)
}

// GetUploadDraftCoverURL generates a presigned URL for uploading a draft cover.
func (c *Client) GetUploadDraftCoverURL(ctx context.Context, artistKey, releaseId, ext, contentType string, expireSeconds int64) (string, error) {
	key := c.resolver.ReleaseDraftCover(artistKey, releaseId, ext)
	return c.GetPresignedUploadURL(ctx, key, contentType, expireSeconds)
}

// GetUploadDraftTrackURL generates a presigned URL for uploading a draft track.
func (c *Client) GetUploadDraftTrackURL(ctx context.Context, artistKey, releaseId string, fileId, fileType string, expireSeconds int64) (string, string, error) {
	format := c.resolver.GetFormatFolder(fileType)
	storagePath := c.resolver.ReleaseDraftTrack(artistKey, releaseId, format, fileId)

	url, err := c.GetPresignedUploadURL(ctx, storagePath, fileType, expireSeconds)
	if err != nil {
		return "", "", err
	}

	return url, storagePath, nil
}
