package models

import "time"

// Artist represents an artist in the system
type Artist struct {
	DocumentMeta `tstype:",extends"`

	Id          string    `json:"id"`
	Name        string    `json:"name" binding:"required"`
	Slug        string    `json:"slug" binding:"required"`
	Bio         string    `json:"bio"`
	CreatedAt   time.Time `json:"createdAt"`
	PublishedAt time.Time `json:"publishedAt,omitempty"`
}

// ArtistFollower is the public-safe user shape shown on an Artist profile.
type ArtistFollower struct {
	Username   string    `json:"username"`
	FollowedAt time.Time `json:"followedAt"`
}

// ArtistProfile represents an Artist with public relationship counts.
type ArtistProfile struct {
	Artist `tstype:",extends"`

	FollowerCount   int              `json:"followerCount"`
	SubscriberCount int              `json:"subscriberCount"`
	Followers       []ArtistFollower `json:"followers"`
}

// FollowArtistResponse represents the current Follow state after mutation.
type FollowArtistResponse struct {
	ArtistID      string `json:"artistId"`
	Following     bool   `json:"following"`
	FollowerCount int    `json:"followerCount"`
}

// RegisterArtistNameRequest represents the request to create an artist name
type RegisterArtistNameRequest struct {
	Name string `json:"name" binding:"required"`
}

// CreateArtistRequest represents the request body for creating an artist
type CreateArtistRequest struct {
	Name string `json:"name" binding:"required"`
	Slug string `json:"slug" binding:"required"`
	Bio  string `json:"bio,omitempty"`
}

// UpdateArtistRequest represents the request body for updating an artist
type UpdateArtistRequest struct {
	Name string `json:"name,omitempty"`
	Bio  string `json:"bio,omitempty"`
}
