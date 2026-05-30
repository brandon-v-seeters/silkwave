package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/avelino/slugify"
	"github.com/brandon-v-seeters/go-silk-wave/internal/database"
	"github.com/brandon-v-seeters/go-silk-wave/internal/models"
	"github.com/google/uuid"
)

const (
	artistsCollection      = "Artists"
	followsCollection      = "Follows"
	usersArtistsCollection = "UsersArtists"
)

var (
	ErrArtistNotFound           = errors.New("artist not found")
	ErrInvalidArtistName        = errors.New("invalid artist name")
	ErrArtistRegistrationFailed = errors.New("artist registration failed")
)

type ArtistRepository struct {
	db *database.ArangoDB
}

func NewArtistRepository(db *database.ArangoDB) *ArtistRepository {
	return &ArtistRepository{db: db}
}

func (r *ArtistRepository) GetBySlug(ctx context.Context, slug string) (*models.ArtistProfile, error) {
	q := /*aql*/ `
		FOR artist IN Artists
			FILTER artist.slug == @slug
			LIMIT 1
			LET followers = (
				FOR follow IN Follows
					FILTER follow._to == CONCAT("Artists/", artist._key)
					FOR user IN Users
						FILTER follow._from == CONCAT("Users/", user._key)
						SORT follow.createdAt DESC
						RETURN { username: user.username, followedAt: follow.createdAt }
			)
			LET activeSubscriberKeys = UNIQUE(
				FOR subscriber IN Subscribers
					FILTER subscriber.artistKey == artist._key
					FILTER subscriber.status == @active
					RETURN subscriber.subscriberKey
			)
			RETURN MERGE(artist, {
				followerCount: LENGTH(followers),
				subscriberCount: LENGTH(activeSubscriberKeys),
				followers: followers
			})
	`

	artists, err := database.Query[models.ArtistProfile](ctx, r.db, q, map[string]interface{}{
		"active": string(models.SubscriptionStatusActive),
		"slug":   slug,
	})
	if err != nil {
		return nil, fmt.Errorf("get artist by slug: %w", err)
	}
	if len(artists) == 0 {
		return nil, ErrArtistNotFound
	}

	return &artists[0], nil
}

func (r *ArtistRepository) RegisterName(ctx context.Context, userKey, name string) error {
	name = strings.TrimSpace(name)
	baseSlug := slugify.Slugify(name)
	if name == "" || baseSlug == "" {
		return ErrInvalidArtistName
	}

	const maxAttempts = 5
	var result *registerArtistResult
	var err error

	for attempt := 0; attempt < maxAttempts; attempt++ {
		result, err = r.registerNameAttempt(ctx, userKey, name, baseSlug)
		if err == nil || !database.IsUniqueConstraintError(err) {
			break
		}
	}

	if err != nil {
		return fmt.Errorf("register artist name: %w", err)
	}
	if result == nil || !result.UserUpdated || !result.ArtistCreated || !result.UserArtistLinked {
		return ErrArtistRegistrationFailed
	}

	return nil
}

func (r *ArtistRepository) Follow(ctx context.Context, userKey, artistID string) (*models.FollowArtistResponse, error) {
	q := /*aql*/ `
		FOR artist IN Artists
			FILTER artist.id == @artistId
			LIMIT 1
			LET edge = FIRST(
				UPSERT { _from: CONCAT("Users/", @userKey), _to: CONCAT("Artists/", artist._key) }
					INSERT {
						_from: CONCAT("Users/", @userKey),
						_to: CONCAT("Artists/", artist._key),
						createdAt: @now
					}
					UPDATE {}
				IN Follows
				RETURN NEW
			)
			LET followerCount = LENGTH(
				FOR follow IN Follows
					FILTER follow._to == CONCAT("Artists/", artist._key)
					RETURN 1
			)
			RETURN { artistId: artist.id, following: edge != null, followerCount: followerCount }
	`

	return r.followState(ctx, q, userKey, artistID)
}

func (r *ArtistRepository) Unfollow(ctx context.Context, userKey, artistID string) (*models.FollowArtistResponse, error) {
	q := /*aql*/ `
		FOR artist IN Artists
			FILTER artist.id == @artistId
			LIMIT 1
			LET removed = (
				FOR follow IN Follows
					FILTER follow._from == CONCAT("Users/", @userKey)
					FILTER follow._to == CONCAT("Artists/", artist._key)
					REMOVE follow IN Follows
					RETURN OLD
			)
			LET followerCount = LENGTH(
				FOR follow IN Follows
					FILTER follow._to == CONCAT("Artists/", artist._key)
					RETURN 1
			)
			RETURN { artistId: artist.id, following: false, followerCount: followerCount }
	`

	return r.followState(ctx, q, userKey, artistID)
}

type registerArtistResult struct {
	UserUpdated      bool `json:"userUpdated"`
	ArtistCreated    bool `json:"artistCreated"`
	UserArtistLinked bool `json:"userArtistLinked"`
}

func (r *ArtistRepository) registerNameAttempt(ctx context.Context, userKey, name, baseSlug string) (*registerArtistResult, error) {
	q := /*aql*/ `
		LET base = @base
		LET taken = (
			FOR artist IN Artists
				FILTER artist.slug == base OR STARTS_WITH(artist.slug, CONCAT(base, "-"))
				LET rest = SUBSTRING(artist.slug, LENGTH(base))
				FILTER rest == "" OR REGEX_TEST(rest, "^-[0-9]+$")
				RETURN rest == "" ? 1 : TO_NUMBER(SUBSTRING(rest, 1))
		)
		LET nextSlug = LENGTH(taken) == 0 ? base : CONCAT(base, "-", MAX(taken) + 1)
		LET artist = FIRST(
			INSERT {
				bio: "",
				createdAt: @createdAt,
				id: @id,
				name: @name,
				slug: nextSlug
			} INTO Artists
			RETURN NEW
		)
		LET user = FIRST(
			UPDATE @userKey WITH { invalidArtistName: false } IN Users
			RETURN NEW
		)
		LET userArtist = FIRST(
			UPSERT { _from: CONCAT("Users/", @userKey), _to: CONCAT("Artists/", artist._key) }
				INSERT {
					_from: CONCAT("Users/", @userKey),
					_to: CONCAT("Artists/", artist._key),
					createdAt: @createdAt
				}
				UPDATE {}
			IN UsersArtists
			RETURN NEW
		)
		RETURN {
			userUpdated: user != null,
			artistCreated: artist != null,
			userArtistLinked: userArtist != null
		}
	`

	return database.QueryOne[registerArtistResult](ctx, r.db, q, map[string]interface{}{
		"base":      baseSlug,
		"createdAt": time.Now(),
		"id":        uuid.New().String(),
		"name":      name,
		"userKey":   userKey,
	})
}

func (r *ArtistRepository) followState(ctx context.Context, q, userKey, artistID string) (*models.FollowArtistResponse, error) {
	results, err := database.Query[models.FollowArtistResponse](ctx, r.db, q, map[string]interface{}{
		"artistId": artistID,
		"now":      time.Now(),
		"userKey":  userKey,
	})
	if err != nil {
		return nil, fmt.Errorf("follow state: %w", err)
	}
	if len(results) == 0 {
		return nil, ErrArtistNotFound
	}

	return &results[0], nil
}
