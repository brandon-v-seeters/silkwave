package models

import "time"

// Edge represents a base edge document in ArangoDB
type Edge struct {
	Key       string    `json:"_key,omitempty"`
	ID        string    `json:"_id,omitempty"`
	Rev       string    `json:"_rev,omitempty"`
	From      string    `json:"_from" binding:"required"`
	To        string    `json:"_to" binding:"required"`
	CreatedAt time.Time `json:"createdAt"`
}

// UserArtist represents the edge between a User and an Artist they manage
type UserArtist struct {
	Edge
}

// Follow represents the free User-to-Artist follow relationship.
type Follow struct {
	Edge
}
