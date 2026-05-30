package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/brandon-v-seeters/go-silk-wave/internal/access"
	"github.com/brandon-v-seeters/go-silk-wave/internal/logger"
	"github.com/brandon-v-seeters/go-silk-wave/internal/middleware"
	"github.com/brandon-v-seeters/go-silk-wave/internal/models"
	"github.com/brandon-v-seeters/go-silk-wave/internal/repository"
	"github.com/brandon-v-seeters/go-silk-wave/internal/slug"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *ReleaseHandler) GetReleases(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	releases, err := h.releases.ListPublished(c.Request.Context(), repository.ListPublishedReleasesParams{
		Limit:       limit,
		Offset:      offset,
		SearchQuery: c.Query("q"),
	})
	if err != nil {
		logger.Error("failed to fetch releases", err)
		respondError(c, http.StatusInternalServerError, "internal_error", "Failed to fetch releases")
		return
	}

	if releases == nil {
		releases = []models.ReleaseWithArtist{}
	}

	respondOK(c, models.ReleasesResponse{
		Releases: releases,
		Limit:    limit,
		Offset:   offset,
	}, "")
}

func (h *ReleaseHandler) GetReleaseByArtistAndSlug(c *gin.Context) {
	artistSlug := c.Param("artistSlug")
	releaseSlug := c.Param("releaseSlug")

	release, err := slug.ResolveRelease(c.Request.Context(), h.releases, artistSlug, releaseSlug)
	if err != nil {
		respondError(c, http.StatusNotFound, "not_found", "Release not found")
		return
	}
	userKey, _ := middleware.GetUserKey(c)
	if err := h.attachPublicReleaseURLs(c.Request.Context(), userKey, release); err != nil {
		logger.Error("failed to attach public release URLs", err, zap.String("releaseId", release.Id))
		respondError(c, http.StatusInternalServerError, "internal_error", "Failed to load release media")
		return
	}

	respondOK(c, release, "")
}

func (h *ReleaseHandler) attachPublicReleaseURLs(ctx context.Context, userKey string, release *models.PublicRelease) error {
	if release == nil {
		return nil
	}

	coverPath := release.Cover
	if release.Assets.CoverArt != nil && release.Assets.CoverArt.Original != "" {
		coverPath = release.Assets.CoverArt.Original
	}
	if coverPath != "" {
		coverURL, err := h.storageClient.GetPublishedReleaseObjectURL(
			ctx,
			release.ArtistKey,
			release.Id,
			coverPath,
			presignedURLExpireSeconds,
		)
		if err != nil {
			return fmt.Errorf("generate cover URL: %w", err)
		}
		if coverURL != "" {
			release.Cover = coverURL
			if release.Assets.CoverArt != nil {
				release.Assets.CoverArt.Original = coverURL
			}
		}
	}

	if !release.StreamingEnabled {
		return nil
	}
	canStream, err := access.CanStream(ctx, h.accessSource, userKey, release.Id)
	if err != nil {
		return fmt.Errorf("check stream access: %w", err)
	}
	if !canStream {
		return nil
	}

	for i := range release.Tracks {
		track := &release.Tracks[i]
		if !track.Uploaded || !track.StreamingEnabled {
			continue
		}

		streamURL, err := h.storageClient.GetPublishedTrackURL(
			ctx,
			release.ArtistKey,
			release.Id,
			track.Files.Original.Path,
			presignedURLExpireSeconds,
		)
		if err != nil {
			return fmt.Errorf("generate stream URL for track %q: %w", track.Id, err)
		}
		if streamURL != "" {
			track.StreamURL = &streamURL
		}
	}

	return nil
}
