package repository

import (
	"context"
	"testing"
	"time"

	"github.com/brandon-v-seeters/go-silk-wave/internal/models"
)

func TestArtistRepositoryFollowAndUnfollowAreIdempotent(t *testing.T) {
	ctx := context.Background()
	db := setupUserRepositoryTestDB(t, ctx)
	repo := NewArtistRepository(db)

	userKey := seedUserWithCredentials(t, ctx, db, models.UserWithCredentials{
		User: models.User{
			Username:       "follow-user",
			Email:          "follow@example.com",
			ValidatedEmail: true,
			Role:           models.UserRoleUser,
			Settings:       models.UserSettings{ReceiveEmails: true},
			CreatedAt:      time.Now(),
		},
		Password:              "hashed-password",
		EmailUnsubscribeToken: "unsubscribe-token",
		Socials:               map[string]models.SocialWithCredentials{},
	})
	seedArtist(t, ctx, db, models.Artist{
		DocumentMeta: models.DocumentMeta{Key: "artist_followed"},
		Id:           "artist-followed-id",
		Name:         "Artist Followed",
		Slug:         "artist-followed",
		CreatedAt:    time.Now(),
	})

	firstFollow, err := repo.Follow(ctx, userKey, "artist-followed-id")
	if err != nil {
		t.Fatalf("Follow returned error: %v", err)
	}
	secondFollow, err := repo.Follow(ctx, userKey, "artist-followed-id")
	if err != nil {
		t.Fatalf("second Follow returned error: %v", err)
	}
	if firstFollow.FollowerCount != 1 || secondFollow.FollowerCount != 1 {
		t.Fatalf("expected idempotent follower count 1, got %d then %d", firstFollow.FollowerCount, secondFollow.FollowerCount)
	}

	profile, err := repo.GetBySlug(ctx, "artist-followed")
	if err != nil {
		t.Fatalf("GetBySlug returned error: %v", err)
	}
	if profile.FollowerCount != 1 {
		t.Fatalf("expected follower count 1, got %d", profile.FollowerCount)
	}
	if len(profile.Followers) != 1 || profile.Followers[0].Username != "follow-user" {
		t.Fatalf("expected follow-user in followers, got %#v", profile.Followers)
	}

	firstUnfollow, err := repo.Unfollow(ctx, userKey, "artist-followed-id")
	if err != nil {
		t.Fatalf("Unfollow returned error: %v", err)
	}
	secondUnfollow, err := repo.Unfollow(ctx, userKey, "artist-followed-id")
	if err != nil {
		t.Fatalf("second Unfollow returned error: %v", err)
	}
	if firstUnfollow.FollowerCount != 0 || secondUnfollow.FollowerCount != 0 {
		t.Fatalf("expected idempotent follower count 0, got %d then %d", firstUnfollow.FollowerCount, secondUnfollow.FollowerCount)
	}
}

func TestArtistRepositoryFollowMissingArtist(t *testing.T) {
	ctx := context.Background()
	db := setupUserRepositoryTestDB(t, ctx)
	repo := NewArtistRepository(db)

	_, err := repo.Follow(ctx, "user-key", "missing-artist")
	if err == nil {
		t.Fatal("expected missing artist follow to fail")
	}
}
