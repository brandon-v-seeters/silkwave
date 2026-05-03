package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/avelino/slugify"
	"github.com/brandon-v-seeters/go-silk-wave/internal/database"
	"github.com/brandon-v-seeters/go-silk-wave/internal/logger"
	"github.com/brandon-v-seeters/go-silk-wave/internal/middleware"
	"github.com/brandon-v-seeters/go-silk-wave/internal/models"
	"github.com/brandon-v-seeters/go-silk-wave/internal/storage"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const presignedURLExpireSeconds int64 = 3600 // 1 hour

type ReleaseHandler struct {
	db            *database.ArangoDB
	storageClient *storage.Client
}

func NewReleaseHandler(db *database.ArangoDB, storageClient *storage.Client) *ReleaseHandler {
	return &ReleaseHandler{db: db, storageClient: storageClient}
}

// =============================================================================
// Public Release Endpoints
// =============================================================================

// GetReleases handles GET /api/releases
func (h *ReleaseHandler) GetReleases(c *gin.Context) {
	// Parse query parameters
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	searchQuery := c.Query("q")

	var query string
	bindVars := map[string]interface{}{
		"limit":  limit,
		"offset": offset,
	}

	if searchQuery != "" {
		// Search by title or description
		query = `
			FOR release IN Releases
				FILTER release.published == true
				FILTER (
					LOWER(release.title) LIKE LOWER(CONCAT("%", @searchQuery, "%")) OR
					LOWER(release.description) LIKE LOWER(CONCAT("%", @searchQuery, "%"))
				)
				LET artist = FIRST(
					FOR a IN Artists
						FILTER a._key == release.artistKey
						RETURN a
				)
				SORT release.releaseDate DESC
				LIMIT @offset, @limit
				RETURN MERGE(release, { artist: artist })
		`
		bindVars["searchQuery"] = searchQuery
	} else {
		// Get latest releases
		query = `
			FOR release IN Releases
				FILTER release.published == true
				LET artist = FIRST(
					FOR a IN Artists
						FILTER a._key == release.artistKey
						RETURN a
				)
				SORT release.releaseDate DESC
				LIMIT @offset, @limit
				RETURN MERGE(release, { artist: artist })
		`
	}

	releases, err := database.Query[models.ReleaseWithArtist](c.Request.Context(), h.db, query, bindVars)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch releases"})
		return
	}

	// Return empty array instead of null
	if releases == nil {
		releases = []models.ReleaseWithArtist{}
	}

	c.JSON(http.StatusOK, models.ReleasesResponse{
		Releases: releases,
		Limit:    limit,
		Offset:   offset,
	})
}

func (h *ReleaseHandler) GetDraftsByArtistKey(c *gin.Context) {
	artistKey := c.Query("artistKey")

	if artistKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "artistKey is required"})
		return
	}

	if !middleware.UserManagesArtist(c, h.db, artistKey) {
		c.JSON(http.StatusForbidden, gin.H{"error": "user does not manage this artist"})
		return
	}

	query := /*aql*/ `
		FOR a IN ReleaseDrafts
		FILTER a.artistKey == @artistKey
		RETURN a
	`
	drafts, err := database.Query[models.ReleaseWithArtist](
		c.Request.Context(),
		h.db,
		query,
		map[string]interface{}{"artistKey": artistKey},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get drafts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"drafts": drafts})
}

// GetReleaseBySlug handles GET /api/releases/:slug
func (h *ReleaseHandler) GetReleaseBySlug(c *gin.Context) {
	slug := c.Param("slug")

	release, err := database.QueryOne[models.ReleaseWithArtist](c.Request.Context(), h.db, `
		FOR release IN Releases
			FILTER release.slug == @slug
			LET artist = FIRST(
				FOR a IN Artists
					FILTER a._key == release.artistKey
					RETURN a
			)
			RETURN MERGE(release, { artist: artist })
	`, map[string]interface{}{"slug": slug})

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Release not found"})
		return
	}

	c.JSON(http.StatusOK, release)
}

// =============================================================================
// Avatar Upload
// =============================================================================

// UploadAvatar handles avatar upload - placeholder for now
func (h *ReleaseHandler) UploadAvatar(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

// =============================================================================
// Draft Release Operations
// =============================================================================

// SaveDraftRelease handles POST /api/releases/draft
// Creates a draft release and returns presigned URLs for uploading files
func (h *ReleaseHandler) SaveDraftRelease(c *gin.Context) {
	var req models.CreateDraftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !middleware.UserManagesArtist(c, h.db, req.ArtistKey) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to manage this artist"})
		return
	}

	ctx := c.Request.Context()
	releaseHash := uuid.New().String()
	fmt.Println("releaseHash", releaseHash)

	presignedUrls, trackRecords, coverPath, err := h.generatePresignedURLs(ctx, req.ArtistKey, releaseHash, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	releaseKey, err := h.saveDraftToDatabase(ctx, req.ArtistKey, releaseHash, coverPath, req, trackRecords)
	if err != nil {
		logger.Error("failed to save draft to database", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.CreateDraftResponse{
		ReleaseHash:   releaseHash,
		ReleaseKey:    releaseKey,
		ArtistKey:     req.ArtistKey,
		PresignedUrls: presignedUrls,
	})
}

// generatePresignedURLs creates presigned URLs for cover and all tracks
func (h *ReleaseHandler) generatePresignedURLs(
	ctx context.Context,
	artistKey, releaseHash string,
	req models.CreateDraftRequest,
) (models.PresignedUrlsDTO, []models.Track, string, error) {

	resolver := h.storageClient.Resolver()
	presignedUrls := models.PresignedUrlsDTO{
		Tracks: make([]models.TrackUrlDTO, 0, len(req.Tracks)),
	}
	var trackRecords []models.Track
	var coverPath string

	// Cover art
	if req.CoverArt != nil {
		// Prefer filename extension (user's explicit choice), fallback to content-type
		ext := storage.GetExtensionFromFileName(req.CoverArt.FileName)
		if ext == "" {
			ext = storage.GetExtensionFromContentType(req.CoverArt.FileType)
		}

		coverPath = resolver.ReleaseDraftCover(artistKey, releaseHash, ext)
		coverPresignedUrl, err := h.storageClient.GetPresignedUploadURL(ctx, coverPath, req.CoverArt.FileType, presignedURLExpireSeconds)
		if err != nil {
			return presignedUrls, nil, "", fmt.Errorf("failed to generate cover URL: %w", err)
		}
		presignedUrls.CoverArt = &coverPresignedUrl
	}

	// Tracks
	for _, trackReq := range req.Tracks {
		fileHash := uuid.New().String()[:8]

		url, storagePath, err := h.storageClient.GetUploadDraftTrackURL(
			ctx, artistKey, releaseHash,
			fileHash, trackReq.FileType,
			presignedURLExpireSeconds,
		)
		if err != nil {
			return presignedUrls, nil, "", fmt.Errorf("failed to generate track URL: %w", err)
		}

		presignedUrls.Tracks = append(presignedUrls.Tracks, models.TrackUrlDTO{
			Hash:         fileHash,
			FileName:     trackReq.FileName,
			PresignedUrl: url,
			StoragePath:  storagePath,
		})

		trackRecords = append(trackRecords, models.Track{
			Hash:  fileHash,
			Title: trackReq.Title,
			Order: trackReq.Order,
			Files: models.TrackFiles{
				Original: models.OriginalFile{
					Path:     storagePath,
					Format:   trackReq.FileType,
					FileSize: trackReq.FileSize,
				},
			},
			Uploaded:  false,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
	}

	return presignedUrls, trackRecords, coverPath, nil
}

// saveDraftToDatabase saves the release and tracks in a transaction
func (h *ReleaseHandler) saveDraftToDatabase(
	ctx context.Context,
	artistKey, releaseHash, coverPath string,
	req models.CreateDraftRequest,
	trackRecords []models.Track,
) (string, error) {
	fmt.Println("saveDraftToDatabase")
	resolver := h.storageClient.Resolver()
	var releaseKey string

	err := database.WithTransaction(ctx, h.db,
		[]string{"ReleaseDrafts", "TrackDrafts"},
		[]string{"ReleaseDrafts", "TrackDrafts"},
		func(tx *database.Transaction) error {
			// Insert release
			release := models.Release{
				Hash:      releaseHash,
				ArtistKey: artistKey,
				Title:     req.Title,
				Slug:      slugify.Slugify(req.Title),
				Status:    models.ReleaseStatusDraft,
				Metadata:  models.ReleaseMetadata{Genres: req.Genres},
				Assets: models.ReleaseAssets{
					BasePath: resolver.ReleaseDraftPrefix(artistKey, releaseHash),
					CoverArt: &models.CoverArtAsset{Original: coverPath},
				},
				TrackCount: len(req.Tracks),
				IsUploaded: false,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			}

			releaseKeyPtr, err := database.TxQueryOne[string](
				tx.TransactionContext(ctx), tx,
				/*aql*/ `
				    UPSERT { _key: @key }
					INSERT @release
					UPDATE @release
					IN ReleaseDrafts 
					RETURN NEW._key
				`,
				map[string]interface{}{"release": release, "key": releaseKey},
			)

			if err != nil {
				return fmt.Errorf("failed to insert draft release: %w", err)
			}

			releaseKey = *releaseKeyPtr

			// Set releaseKey on all tracks
			for i := range trackRecords {
				trackRecords[i].ReleaseKey = releaseKey
			}

			// Insert all tracks
			err = database.TxExec(
				tx.TransactionContext(ctx), tx,
				/*aql*/ `
				FOR track IN @tracks 
				INSERT track IN TrackDrafts
				`,
				map[string]interface{}{"tracks": trackRecords},
			)
			if err != nil {
				return fmt.Errorf("failed to insert draft tracks: %w", err)
			}

			return nil
		})

	return releaseKey, err
}

// =============================================================================
// Confirm Draft Release
// =============================================================================

// ConfirmDraftRelease handles POST /api/releases/:releaseHash/confirm
// Verifies files are uploaded to R2 and updates status
func (h *ReleaseHandler) ConfirmDraftRelease(c *gin.Context) {
	var req models.ConfirmUploadsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	releaseHash := c.Param("releaseHash")
	if releaseHash == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "releaseHash is required"})
		return
	}

	ctx := c.Request.Context()

	// Get and verify release ownership
	release, err := h.getReleaseByHash(ctx, releaseHash)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Release not found"})
		return
	}

	if !middleware.UserManagesArtist(c, h.db, release.ArtistKey) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to confirm this release"})
		return
	}

	// Collect expected paths
	expectedPaths := make([]string, 0, len(req.Tracks)+1)
	if req.CoverArtPath != nil && *req.CoverArtPath != "" {
		expectedPaths = append(expectedPaths, *req.CoverArtPath)
	}
	for _, track := range req.Tracks {
		expectedPaths = append(expectedPaths, track.StoragePath)
	}

	// Verify files exist in R2
	missingFiles, err := h.storageClient.VerifyDraftUploads(ctx, release.ArtistKey, releaseHash, expectedPaths)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify uploads"})
		return
	}
	if len(missingFiles) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Some files were not uploaded", "missingFiles": missingFiles})
		return
	}

	// Update database
	trackKeys := make([]string, len(req.Tracks))
	for i, track := range req.Tracks {
		trackKeys[i] = track.TrackID
	}

	if err := h.markUploadsComplete(ctx, releaseHash, trackKeys); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "releaseHash": releaseHash})
}

// markUploadsComplete updates release and tracks as uploaded
func (h *ReleaseHandler) markUploadsComplete(ctx context.Context, releaseHash string, trackKeys []string) error {
	return database.WithTransaction(ctx, h.db,
		[]string{"ReleaseDrafts", "TrackDrafts"},
		[]string{"ReleaseDrafts", "TrackDrafts"},
		func(tx *database.Transaction) error {
			// Update release
			_, err := database.TxQueryOne[string](
				tx.TransactionContext(ctx), tx,
				`FOR a IN ReleaseDrafts FILTER a.hash == @hash UPDATE a WITH { isUploaded: true, updatedAt: @now } IN ReleaseDrafts RETURN NEW._key`,
				map[string]interface{}{"hash": releaseHash, "now": time.Now()},
			)
			if err != nil {
				return fmt.Errorf("failed to update release: %w", err)
			}

			// Update tracks
			err = database.TxExec(
				tx.TransactionContext(ctx), tx,
				`FOR t IN TrackDrafts FILTER t._key IN @keys UPDATE t WITH { uploaded: true, updatedAt: @now } IN TrackDrafts`,
				map[string]interface{}{"keys": trackKeys, "now": time.Now()},
			)
			if err != nil {
				return fmt.Errorf("failed to update tracks: %w", err)
			}

			return nil
		})
}

// =============================================================================
// Publish Release
// =============================================================================

// PublishRelease handles POST /api/releases/:releaseHash/publish
func (h *ReleaseHandler) PublishRelease(c *gin.Context) {
	releaseHash := c.Param("releaseHash")
	if releaseHash == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "releaseHash is required"})
		return
	}

	ctx := c.Request.Context()

	release, err := h.getReleaseByHash(ctx, releaseHash)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Release not found"})
		return
	}

	if !middleware.UserManagesArtist(c, h.db, release.ArtistKey) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to publish this release"})
		return
	}

	if !release.IsUploaded {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Release files have not been uploaded yet"})
		return
	}

	// Move files in R2
	if err := h.storageClient.PublishRelease(ctx, release.ArtistKey, releaseHash); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish release files"})
		return
	}

	// Update database
	if err := h.updateReleaseStatus(ctx, releaseHash, "published", true); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update release status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "published", "releaseHash": releaseHash})
}

// =============================================================================
// Unpublish Release
// =============================================================================

// UnpublishRelease handles POST /api/releases/:releaseHash/unpublish
func (h *ReleaseHandler) UnpublishRelease(c *gin.Context) {
	releaseHash := c.Param("releaseHash")
	if releaseHash == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "releaseHash is required"})
		return
	}

	ctx := c.Request.Context()

	release, err := h.getReleaseByHash(ctx, releaseHash)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Release not found"})
		return
	}

	if !middleware.UserManagesArtist(c, h.db, release.ArtistKey) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to unpublish this release"})
		return
	}

	// Move files back to draft in R2
	if err := h.storageClient.UnpublishRelease(ctx, release.ArtistKey, releaseHash); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unpublish release files"})
		return
	}

	// Update database
	if err := h.updateReleaseStatus(ctx, releaseHash, "draft", false); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update release status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "unpublished", "releaseHash": releaseHash})
}

// =============================================================================
// Delete Release
// =============================================================================

// DeleteRelease handles DELETE /api/releases/:releaseHash
func (h *ReleaseHandler) DeleteRelease(c *gin.Context) {
	releaseHash := c.Param("releaseHash")
	if releaseHash == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "releaseHash is required"})
		return
	}

	ctx := c.Request.Context()

	release, err := h.getReleaseByHash(ctx, releaseHash)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Release not found"})
		return
	}

	if !middleware.UserManagesArtist(c, h.db, release.ArtistKey) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this release"})
		return
	}

	// Delete files from R2
	if err := h.storageClient.DeleteRelease(ctx, release.ArtistKey, releaseHash); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete release files"})
		return
	}

	// Delete from database
	if err := h.deleteReleaseFromDatabase(ctx, release.Key, releaseHash); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted", "releaseHash": releaseHash})
}

// deleteReleaseFromDatabase removes release and tracks in a transaction
func (h *ReleaseHandler) deleteReleaseFromDatabase(ctx context.Context, releaseKey, releaseHash string) error {
	return database.WithTransaction(ctx, h.db,
		[]string{"ReleaseDrafts", "TrackDrafts"},
		[]string{"ReleaseDrafts", "TrackDrafts"},
		func(tx *database.Transaction) error {
			// Delete tracks
			err := database.TxExec(
				tx.TransactionContext(ctx), tx,
				`FOR t IN TrackDrafts FILTER t.releaseKey == @key REMOVE t IN TrackDrafts`,
				map[string]interface{}{"key": releaseKey},
			)
			if err != nil {
				return fmt.Errorf("failed to delete tracks: %w", err)
			}

			// Delete release
			err = database.TxExec(
				tx.TransactionContext(ctx), tx,
				`FOR a IN ReleaseDrafts FILTER a.hash == @hash REMOVE a IN ReleaseDrafts`,
				map[string]interface{}{"hash": releaseHash},
			)
			if err != nil {
				return fmt.Errorf("failed to delete release: %w", err)
			}

			return nil
		})
}

// =============================================================================
// Helper Functions
// =============================================================================

// getReleaseByHash retrieves a release by its hash
func (h *ReleaseHandler) getReleaseByHash(ctx context.Context, releaseHash string) (*models.Release, error) {
	return database.QueryOne[models.Release](ctx, h.db,
		`FOR a IN ReleaseDrafts FILTER a.hash == @hash RETURN a`,
		map[string]interface{}{"hash": releaseHash},
	)
}

// updateReleaseStatus updates the published status of a release
func (h *ReleaseHandler) updateReleaseStatus(ctx context.Context, releaseHash, status string, published bool) error {
	_, err := database.QueryOne[string](ctx, h.db,
		`FOR a IN ReleaseDrafts FILTER a.hash == @hash UPDATE a WITH { status: @status, published: @published, updatedAt: @now } IN ReleaseDrafts RETURN NEW._key`,
		map[string]interface{}{"hash": releaseHash, "status": status, "published": published, "now": time.Now()},
	)
	return err
}

func (h *ReleaseHandler) GetDraftByKey(c *gin.Context) {
	releaseKey := c.Param("releaseKey")
	if releaseKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "releaseKey is required"})
		return
	}

	draft, err := database.QueryOne[models.ReleaseWithTracks](c.Request.Context(), h.db,
		/*aql*/ `
		FOR draft IN ReleaseDrafts FILTER draft._key == @key
		LET tracks = FIRST(
			FOR track IN TrackDrafts
			FILTER track.releaseKey == draft._key
			RETURN track
		)
		RETURN MERGE(draft, { tracks: tracks })
		`,
		map[string]interface{}{"key": releaseKey},
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Draft not found"})
		return
	}

	fmt.Println("draft", draft)

	c.JSON(http.StatusOK, gin.H{"draft": draft})
}
