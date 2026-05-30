package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/brandon-v-seeters/go-silk-wave/internal/logger"
	"github.com/brandon-v-seeters/go-silk-wave/internal/middleware"
	"github.com/brandon-v-seeters/go-silk-wave/internal/models"
	"github.com/brandon-v-seeters/go-silk-wave/internal/repository"
	"github.com/brandon-v-seeters/go-silk-wave/internal/slug"
	"github.com/brandon-v-seeters/go-silk-wave/internal/storage"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (h *ReleaseHandler) GetDraftsByArtistKey(c *gin.Context) {
	artistKey := c.Query("artistKey")
	if artistKey == "" {
		respondError(c, http.StatusBadRequest, "invalid_request", "artistKey is required")
		return
	}

	if !middleware.UserManagesArtist(c, h.db, artistKey) {
		respondError(c, http.StatusForbidden, "forbidden", "user does not manage this artist")
		return
	}

	drafts, err := h.releases.ListDraftsByArtistKey(c.Request.Context(), artistKey)
	if err != nil {
		logger.Error("failed to get release drafts", err, zap.String("artistKey", artistKey))
		respondError(c, http.StatusInternalServerError, "internal_error", "failed to get drafts")
		return
	}

	if drafts == nil {
		drafts = []models.ReleaseWithArtist{}
	}

	respondOK(c, gin.H{"drafts": drafts}, "")
}

func (h *ReleaseHandler) SaveDraftRelease(c *gin.Context) {
	var req models.CreateDraftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	if !middleware.UserManagesArtist(c, h.db, req.ArtistKey) {
		respondError(c, http.StatusForbidden, "forbidden", "You are not authorized to manage this artist")
		return
	}

	ctx := c.Request.Context()
	releaseId := uuid.New().String()

	presignedUrls, trackRecords, coverPath, err := h.generatePresignedURLs(ctx, req.ArtistKey, releaseId, req)
	if err != nil {
		logger.Error("failed to generate release upload URLs", err, zap.String("artistKey", req.ArtistKey))
		respondError(c, http.StatusInternalServerError, "internal_error", "Failed to generate upload URLs")
		return
	}

	release := h.buildDraftRelease(req, releaseId, coverPath)
	releaseKey, err := h.releases.CreateDraftWithTracks(ctx, release, trackRecords)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateReleaseSlug) {
			respondError(c, http.StatusConflict, "duplicate_slug", "A release with this slug already exists for this artist")
			return
		}

		logger.Error("failed to save draft release", err, zap.String("artistKey", req.ArtistKey))
		respondError(c, http.StatusInternalServerError, "internal_error", "Failed to save draft release")
		return
	}

	respondOK(c, models.CreateDraftResponse{
		ReleaseId:     releaseId,
		ReleaseKey:    releaseKey,
		ArtistKey:     req.ArtistKey,
		PresignedUrls: presignedUrls,
	}, "")
}

func (h *ReleaseHandler) generatePresignedURLs(
	ctx context.Context,
	artistKey, releaseId string,
	req models.CreateDraftRequest,
) (models.PresignedUrlsDTO, []models.Track, string, error) {
	resolver := h.storageClient.Resolver()
	presignedUrls := models.PresignedUrlsDTO{
		Tracks: make([]models.TrackUrlDTO, 0, len(req.Tracks)),
	}
	var trackRecords []models.Track
	var coverPath string

	if req.CoverArt != nil {
		ext := storage.GetExtensionFromFileName(req.CoverArt.FileName)
		if ext == "" {
			ext = storage.GetExtensionFromContentType(req.CoverArt.FileType)
		}

		coverPath = resolver.ReleaseDraftCover(artistKey, releaseId, ext)
		coverPresignedURL, err := h.storageClient.GetPresignedUploadURL(ctx, coverPath, req.CoverArt.FileType, presignedURLExpireSeconds)
		if err != nil {
			return presignedUrls, nil, "", fmt.Errorf("generate cover upload URL: %w", err)
		}
		presignedUrls.CoverArt = &coverPresignedURL
	}

	for _, trackReq := range req.Tracks {
		fileId := uuid.New().String()[:8]

		url, storagePath, err := h.storageClient.GetUploadDraftTrackURL(
			ctx,
			artistKey,
			releaseId,
			fileId,
			trackReq.FileType,
			presignedURLExpireSeconds,
		)
		if err != nil {
			return presignedUrls, nil, "", fmt.Errorf("generate track upload URL: %w", err)
		}

		presignedUrls.Tracks = append(presignedUrls.Tracks, models.TrackUrlDTO{
			Id:           fileId,
			FileName:     trackReq.FileName,
			PresignedUrl: url,
			StoragePath:  storagePath,
		})

		trackRecords = append(trackRecords, models.Track{
			Id:    fileId,
			Title: trackReq.Title,
			Order: trackReq.Order,
			Files: models.TrackFiles{
				Original: models.OriginalFile{
					Path:     storagePath,
					Format:   trackReq.FileType,
					FileSize: trackReq.FileSize,
				},
			},
			Uploaded:         false,
			StreamingEnabled: true,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		})
	}

	return presignedUrls, trackRecords, coverPath, nil
}

func (h *ReleaseHandler) buildDraftRelease(req models.CreateDraftRequest, releaseId, coverPath string) models.Release {
	resolver := h.storageClient.Resolver()

	return models.Release{
		Id:        releaseId,
		ArtistKey: req.ArtistKey,
		Title:     req.Title,
		Slug:      slug.Slugify(req.Title),
		Status:    models.ReleaseStatusDraft,
		Metadata:  models.ReleaseMetadata{Genres: req.Genres},
		Assets: models.ReleaseAssets{
			BasePath: resolver.ReleaseDraftPrefix(req.ArtistKey, releaseId),
			CoverArt: &models.CoverArtAsset{Original: coverPath},
		},
		TrackCount:       len(req.Tracks),
		StreamingEnabled: true,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
}

func (h *ReleaseHandler) ConfirmDraftRelease(c *gin.Context) {
	var req models.ConfirmUploadsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	releaseId := c.Param("releaseId")
	if releaseId == "" {
		respondError(c, http.StatusBadRequest, "invalid_request", "releaseId is required")
		return
	}

	ctx := c.Request.Context()
	release, err := h.releases.GetReleaseById(ctx, releaseId)
	if err != nil {
		respondError(c, http.StatusNotFound, "not_found", "Release not found")
		return
	}

	if !middleware.UserManagesArtist(c, h.db, release.ArtistKey) {
		respondError(c, http.StatusForbidden, "forbidden", "You are not authorized to confirm this release")
		return
	}
	if release.Status != models.ReleaseStatusDraft {
		respondError(c, http.StatusBadRequest, "invalid_state", "Only draft releases can confirm uploads")
		return
	}

	expectedPaths := make([]string, 0, len(req.Tracks)+1)
	if req.CoverArtPath != nil && *req.CoverArtPath != "" {
		expectedPaths = append(expectedPaths, *req.CoverArtPath)
	}
	for _, track := range req.Tracks {
		expectedPaths = append(expectedPaths, track.StoragePath)
	}

	missingFiles, err := h.storageClient.VerifyDraftUploads(ctx, release.ArtistKey, releaseId, expectedPaths)
	if err != nil {
		logger.Error("failed to verify release uploads", err, zap.String("releaseId", releaseId))
		respondError(c, http.StatusInternalServerError, "internal_error", "Failed to verify uploads")
		return
	}
	if len(missingFiles) > 0 {
		respondError(c, http.StatusBadRequest, "missing_uploads", "Some files were not uploaded")
		return
	}

	trackIds := make([]string, len(req.Tracks))
	for i, track := range req.Tracks {
		trackIds[i] = track.TrackID
	}

	if err := h.releases.MarkUploadsComplete(ctx, releaseId, trackIds); err != nil {
		logger.Error("failed to mark release uploads complete", err, zap.String("releaseId", releaseId))
		respondError(c, http.StatusInternalServerError, "internal_error", "Failed to update release uploads")
		return
	}

	respondOK(c, gin.H{"releaseId": releaseId}, "Release uploads confirmed")
}

func (h *ReleaseHandler) GetDraftByKey(c *gin.Context) {
	releaseKey := c.Param("releaseKey")
	if releaseKey == "" {
		respondError(c, http.StatusBadRequest, "invalid_request", "releaseKey is required")
		return
	}

	draft, err := h.releases.GetDraftWithTracksByKey(c.Request.Context(), releaseKey)
	if err != nil {
		respondError(c, http.StatusNotFound, "not_found", "Draft not found")
		return
	}

	respondOK(c, gin.H{"draft": draft}, "")
}
