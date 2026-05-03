package models

import "time"

// Artist represents an artist in the system
type Artist struct {
	DocumentMeta `tstype:",extends"`

	Name        string    `json:"name" binding:"required"`
	Slug        string    `json:"slug" binding:"required"`
	Bio         string    `json:"bio"`
	CreatedAt   time.Time `json:"createdAt"`
	PublishedAt time.Time `json:"publishedAt,omitempty"`
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
