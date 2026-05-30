package storage

import "testing"

func TestPendingCoverPublishedKey(t *testing.T) {
	t.Parallel()

	resolver := NewKeyResolver()
	key, err := resolver.PendingCoverPublishedKey(
		"artist-key",
		"release-id",
		"artist_content/artist-key/releases/release-id/draft/cover.jpg",
	)
	if err != nil {
		t.Fatalf("PendingCoverPublishedKey returned error: %v", err)
	}
	if key != "artist_content/artist-key/releases/release-id/cover.jpg" {
		t.Fatalf("expected published cover key, got %q", key)
	}
}

func TestPendingCoverPublishedKeyRejectsNonCoverPath(t *testing.T) {
	t.Parallel()

	resolver := NewKeyResolver()
	_, err := resolver.PendingCoverPublishedKey(
		"artist-key",
		"release-id",
		"artist_content/artist-key/releases/release-id/draft/wavs/track.wav",
	)
	if err == nil {
		t.Fatal("expected non-cover path to fail")
	}
}
