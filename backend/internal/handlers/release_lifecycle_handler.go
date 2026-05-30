package handlers

import (
	"net/http"

	"github.com/brandon-v-seeters/go-silk-wave/internal/logger"
	"github.com/brandon-v-seeters/go-silk-wave/internal/middleware"
	"github.com/brandon-v-seeters/go-silk-wave/internal/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *ReleaseHandler) PublishRelease(c *gin.Context) {
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
		respondError(c, http.StatusForbidden, "forbidden", "You are not authorized to publish this release")
		return
	}
	if release.Status != models.ReleaseStatusDraft && release.Status != models.ReleaseStatusArchived {
		respondError(c, http.StatusBadRequest, "invalid_state", "Only draft or archived releases can be published")
		return
	}

	if release.Status == models.ReleaseStatusDraft {
		uploadsComplete, err := h.releases.UploadsComplete(ctx, release.Key)
		if err != nil {
			logger.Error("failed to check release upload readiness", err, zap.String("releaseId", releaseId))
			respondError(c, http.StatusInternalServerError, "internal_error", "Failed to check release upload readiness")
			return
		}
		if !uploadsComplete {
			respondError(c, http.StatusBadRequest, "invalid_state", "Release files have not been uploaded yet")
			return
		}

		if err := h.storageClient.PublishRelease(ctx, release.ArtistKey, releaseId); err != nil {
			logger.Error("failed to publish release files", err, zap.String("releaseId", releaseId))
			respondError(c, http.StatusInternalServerError, "internal_error", "Failed to publish release files")
			return
		}
	}

	if err := h.releases.UpdateReleaseStatus(ctx, releaseId, models.ReleaseStatusPublished); err != nil {
		logger.Error("failed to update release status", err, zap.String("releaseId", releaseId))
		respondError(c, http.StatusInternalServerError, "internal_error", "Failed to update release status")
		return
	}

	respondOK(c, gin.H{"status": string(models.ReleaseStatusPublished), "releaseId": releaseId}, "")
}

func (h *ReleaseHandler) ArchiveRelease(c *gin.Context) {
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
		respondError(c, http.StatusForbidden, "forbidden", "You are not authorized to archive this release")
		return
	}

	if release.Status != models.ReleaseStatusPublished {
		respondError(c, http.StatusBadRequest, "invalid_state", "Only published releases can be archived")
		return
	}

	if err := h.releases.UpdateReleaseStatus(ctx, releaseId, models.ReleaseStatusArchived); err != nil {
		logger.Error("failed to archive release", err, zap.String("releaseId", releaseId))
		respondError(c, http.StatusInternalServerError, "internal_error", "Failed to archive release")
		return
	}

	respondOK(c, gin.H{"status": string(models.ReleaseStatusArchived), "releaseId": releaseId}, "")
}

func (h *ReleaseHandler) DeleteRelease(c *gin.Context) {
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
		respondError(c, http.StatusForbidden, "forbidden", "You are not authorized to delete this release")
		return
	}

	if err := h.storageClient.DeleteRelease(ctx, release.ArtistKey, releaseId); err != nil {
		logger.Error("failed to delete release files", err, zap.String("releaseId", releaseId))
		respondError(c, http.StatusInternalServerError, "internal_error", "Failed to delete release files")
		return
	}

	if err := h.releases.DeleteReleaseWithTracks(ctx, release.Key, releaseId); err != nil {
		logger.Error("failed to delete release", err, zap.String("releaseId", releaseId))
		respondError(c, http.StatusInternalServerError, "internal_error", "Failed to delete release")
		return
	}

	respondOK(c, gin.H{"releaseId": releaseId}, "")
}
