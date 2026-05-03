package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/avelino/slugify"
	"github.com/brandon-v-seeters/go-silk-wave/internal/database"
	"github.com/brandon-v-seeters/go-silk-wave/internal/middleware"
	"github.com/brandon-v-seeters/go-silk-wave/internal/models"
	"github.com/gin-gonic/gin"
)

type ArtistHandler struct {
	db *database.ArangoDB
}

func NewArtistHandler(db *database.ArangoDB) *ArtistHandler {
	return &ArtistHandler{db: db}
}

// GetArtistBySlug handles GET /api/artists/:artistSlug
func (h *ArtistHandler) GetArtistBySlug(c *gin.Context) {
	artistSlug := c.Param("artistSlug")

	artist, err := database.QueryOne[models.Artist](c.Request.Context(), h.db,
		"FOR artist IN Artists FILTER artist.slug == @artistSlug RETURN artist",
		map[string]interface{}{"artistSlug": artistSlug})

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Artist not found."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"artist": artist})
}

// RegisterArtistName handles POST /api/register/artist-name
func (h *ArtistHandler) RegisterArtistName(c *gin.Context) {
	var req models.RegisterArtistNameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "An artist name is required."})
		return
	}


	if req.Name == "" || strings.TrimSpace(req.Name) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "An artist name is required."})
		return
	}

	baseSlug := slugify.Slugify(req.Name)

	// Get user key from context
	userKey, ok := middleware.GetUserKey(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Not authenticated."})
		return
	}

	type RegisterResult struct {
		UserCreated   bool `json:"userCreated"`
		ArtistCreated bool `json:"artistCreated"`
	}

	// Compute next available slug and insert atomically. Retry on unique-constraint
	// conflict in case a concurrent request claimed the same computed slug.
	const maxAttempts = 5
	var result *RegisterResult
	var err error
	for attempt := 0; attempt < maxAttempts; attempt++ {
		result, err = database.QueryOne[RegisterResult](c.Request.Context(), h.db, `
			LET base = @base
			LET taken = (
				FOR a IN Artists
					FILTER a.slug == base OR STARTS_WITH(a.slug, CONCAT(base, "-"))
					LET rest = SUBSTRING(a.slug, LENGTH(base))
					FILTER rest == "" OR REGEX_TEST(rest, "^-[0-9]+$")
					RETURN rest == "" ? 1 : TO_NUMBER(SUBSTRING(rest, 1))
			)
			LET nextSlug = LENGTH(taken) == 0 ? base : CONCAT(base, "-", MAX(taken) + 1)

			LET artist = FIRST(
				INSERT {
					bio: "",
					createdAt: @createdAt,
					name: @name,
					slug: nextSlug,
					userKey: @userKey
				} INTO Artists
				RETURN NEW
			)

			LET user = FIRST(
				UPDATE @userKey WITH {
					invalidArtistName: false,
					artistKeys: [artist._key]
				} IN Users
				RETURN NEW
			)

			RETURN {
				userCreated: user != null,
				artistCreated: artist != null
			}
		`, map[string]interface{}{
			"userKey":   userKey,
			"name":      req.Name,
			"base":      baseSlug,
			"createdAt": time.Now(),
		})

		if err == nil || !isUniqueConstraintError(err) {
			break
		}
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Artist name could not be created."})
		return
	}

	if result == nil || !result.UserCreated || !result.ArtistCreated {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Artist name could not be created."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Artist name created successfully."})
}
