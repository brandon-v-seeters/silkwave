package handlers

import (
	"errors"
	"net/http"

	"github.com/brandon-v-seeters/go-silk-wave/internal/logger"
	"github.com/brandon-v-seeters/go-silk-wave/internal/middleware"
	"github.com/brandon-v-seeters/go-silk-wave/internal/models"
	"github.com/brandon-v-seeters/go-silk-wave/internal/pending"
	"github.com/brandon-v-seeters/go-silk-wave/internal/repository"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *ReleaseHandler) StagePendingEdit(c *gin.Context) {
	releaseId := c.Param("releaseId")
	if releaseId == "" {
		respondError(c, http.StatusBadRequest, "invalid_request", "releaseId is required")
		return
	}

	var req models.PendingReleaseEdit
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	release, ok := h.authorizedRelease(c, releaseId, "stage pending edit")
	if !ok {
		return
	}

	preview, err := h.releases.StagePendingEdit(c.Request.Context(), release.Id, req)
	if err != nil {
		h.respondPendingError(c, err, release.Id)
		return
	}

	respondOK(c, gin.H{"release": preview}, "")
}

func (h *ReleaseHandler) DiscardPendingEdit(c *gin.Context) {
	releaseId := c.Param("releaseId")
	if releaseId == "" {
		respondError(c, http.StatusBadRequest, "invalid_request", "releaseId is required")
		return
	}

	release, ok := h.authorizedRelease(c, releaseId, "discard pending edit")
	if !ok {
		return
	}

	if err := h.releases.DiscardPendingEdit(c.Request.Context(), release.Id); err != nil {
		logger.Error("failed to discard pending edit", err, zap.String("releaseId", release.Id))
		respondError(c, http.StatusInternalServerError, "internal_error", "Failed to discard pending edit")
		return
	}

	respondOK(c, gin.H{"releaseId": release.Id}, "")
}

func (h *ReleaseHandler) GetPendingPreview(c *gin.Context) {
	releaseId := c.Param("releaseId")
	if releaseId == "" {
		respondError(c, http.StatusBadRequest, "invalid_request", "releaseId is required")
		return
	}

	release, ok := h.authorizedRelease(c, releaseId, "preview pending edit")
	if !ok {
		return
	}

	preview, err := h.releases.PendingPreview(c.Request.Context(), release.Id)
	if err != nil {
		h.respondPendingError(c, err, release.Id)
		return
	}

	respondOK(c, gin.H{"release": preview}, "")
}

func (h *ReleaseHandler) PublishPendingEdit(c *gin.Context) {
	releaseId := c.Param("releaseId")
	if releaseId == "" {
		respondError(c, http.StatusBadRequest, "invalid_request", "releaseId is required")
		return
	}

	release, ok := h.authorizedRelease(c, releaseId, "publish pending edit")
	if !ok {
		return
	}

	if release.Pending != nil {
		if err := pending.Validate(*release.Pending, release.Tracks); err != nil {
			h.respondPendingError(c, err, release.Id)
			return
		}
	}

	if coverPath := pendingCoverPath(release.Pending); coverPath != "" {
		if err := h.storageClient.PublishPendingCover(c.Request.Context(), release.ArtistKey, release.Id, coverPath); err != nil {
			logger.Error("failed to publish pending cover", err, zap.String("releaseId", release.Id))
			respondError(c, http.StatusInternalServerError, "internal_error", "Failed to publish pending cover")
			return
		}
	}

	published, err := h.releases.PublishPendingEdit(c.Request.Context(), release.Id)
	if err != nil {
		h.respondPendingError(c, err, release.Id)
		return
	}

	respondOK(c, gin.H{"release": published}, "")
}

func (h *ReleaseHandler) authorizedRelease(c *gin.Context, releaseId, action string) (*models.ReleaseWithTracks, bool) {
	release, err := h.releases.GetReleaseWithTracksById(c.Request.Context(), releaseId)
	if err != nil {
		respondError(c, http.StatusNotFound, "not_found", "Release not found")
		return nil, false
	}

	if !middleware.UserManagesArtist(c, h.db, release.ArtistKey) {
		respondError(c, http.StatusForbidden, "forbidden", "You are not authorized to "+action)
		return nil, false
	}

	return release, true
}

func (h *ReleaseHandler) respondPendingError(c *gin.Context, err error, releaseId string) {
	switch {
	case errors.Is(err, repository.ErrReleaseMustBePublished):
		respondError(c, http.StatusBadRequest, "invalid_state", "Only published releases can have pending edits")
	case errors.Is(err, repository.ErrNoPendingEdit):
		respondError(c, http.StatusBadRequest, "invalid_state", "Release has no pending edit")
	case errors.Is(err, pending.ErrTrackMembershipChange):
		respondError(c, http.StatusBadRequest, "track_membership_change", "Pending edits cannot add or remove tracks")
	case errors.Is(err, pending.ErrAudioReplacement):
		respondError(c, http.StatusBadRequest, "audio_replacement", "Pending edits cannot replace audio files")
	case errors.Is(err, pending.ErrMalformedEdit):
		respondError(c, http.StatusBadRequest, "malformed_pending_edit", "Pending edit is malformed")
	default:
		logger.Error("pending edit failed", err, zap.String("releaseId", releaseId))
		respondError(c, http.StatusInternalServerError, "internal_error", "Pending edit failed")
	}
}

func pendingCoverPath(edit *models.PendingReleaseEdit) string {
	if edit == nil || edit.Assets == nil || edit.Assets.CoverArt == nil {
		return ""
	}

	return edit.Assets.CoverArt.Original
}
