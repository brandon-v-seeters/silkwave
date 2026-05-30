package handlers

import (
	"errors"
	"net/http"

	"github.com/brandon-v-seeters/go-silk-wave/internal/middleware"
	"github.com/brandon-v-seeters/go-silk-wave/internal/models"
	"github.com/brandon-v-seeters/go-silk-wave/internal/repository"
	"github.com/gin-gonic/gin"
)

type ArtistHandler struct {
	artists *repository.ArtistRepository
}

func NewArtistHandler(artists *repository.ArtistRepository) *ArtistHandler {
	return &ArtistHandler{artists: artists}
}

// GetArtistBySlug handles GET /api/artists/:artistSlug
func (h *ArtistHandler) GetArtistBySlug(c *gin.Context) {
	artistSlug := c.Param("artistSlug")

	artist, err := h.artists.GetBySlug(c.Request.Context(), artistSlug)
	if err != nil {
		if errors.Is(err, repository.ErrArtistNotFound) {
			respondError(c, http.StatusNotFound, "not_found", "Artist not found.")
			return
		}
		respondError(c, http.StatusInternalServerError, "server_error", "Artist could not be loaded.")
		return
	}

	respondOK(c, gin.H{"artist": artist}, "")
}

// RegisterArtistName handles POST /api/register/artist-name
func (h *ArtistHandler) RegisterArtistName(c *gin.Context) {
	var req models.RegisterArtistNameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "invalid_request", "An artist name is required.")
		return
	}

	userKey, ok := middleware.GetUserKey(c)
	if !ok {
		respondError(c, http.StatusUnauthorized, "unauthorized", "Not authenticated.")
		return
	}

	if err := h.artists.RegisterName(c.Request.Context(), userKey, req.Name); err != nil {
		if errors.Is(err, repository.ErrInvalidArtistName) {
			respondError(c, http.StatusBadRequest, "invalid_request", "An artist name is required.")
			return
		}
		respondError(c, http.StatusBadRequest, "artist_name_unavailable", "Artist name could not be created.")
		return
	}

	respondOK(c, nil, "Artist name created successfully.")
}

func (h *ArtistHandler) FollowArtist(c *gin.Context) {
	userKey, ok := middleware.GetUserKey(c)
	if !ok {
		respondError(c, http.StatusUnauthorized, "unauthorized", "Not authenticated.")
		return
	}

	response, err := h.artists.Follow(c.Request.Context(), userKey, c.Param("artistId"))
	if err != nil {
		respondArtistFollowError(c, err)
		return
	}

	respondOK(c, response, "")
}

func (h *ArtistHandler) UnfollowArtist(c *gin.Context) {
	userKey, ok := middleware.GetUserKey(c)
	if !ok {
		respondError(c, http.StatusUnauthorized, "unauthorized", "Not authenticated.")
		return
	}

	response, err := h.artists.Unfollow(c.Request.Context(), userKey, c.Param("artistId"))
	if err != nil {
		respondArtistFollowError(c, err)
		return
	}

	respondOK(c, response, "")
}

func respondArtistFollowError(c *gin.Context, err error) {
	if errors.Is(err, repository.ErrArtistNotFound) {
		respondError(c, http.StatusNotFound, "not_found", "Artist not found.")
		return
	}

	respondError(c, http.StatusInternalServerError, "server_error", "Follow state could not be updated.")
}
