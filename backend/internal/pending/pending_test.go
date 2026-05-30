package pending

import (
	"errors"
	"testing"

	"github.com/brandon-v-seeters/go-silk-wave/internal/models"
)

func TestApplyMergeCases(t *testing.T) {
	t.Parallel()

	title := "Corrected Title"
	description := "new description"
	trackTitle := "Corrected Track"
	downloadEnabled := true
	metadata := models.ReleaseMetadata{Genres: []string{"ambient"}}
	trackMetadata := models.TrackMetadata{Tags: []string{"focus"}}

	release := releaseWithTracks()
	merged, err := Apply(release, models.PendingReleaseEdit{
		Title:       &title,
		Description: &description,
		Metadata:    &metadata,
		Tracks: []models.PendingTrackEdit{{
			TrackKey:        "track-key-1",
			Title:           &trackTitle,
			DownloadEnabled: &downloadEnabled,
			Metadata:        &trackMetadata,
		}},
	})
	if err != nil {
		t.Fatalf("Apply returned error: %v", err)
	}

	if merged.Title != title {
		t.Fatalf("expected title override %q, got %q", title, merged.Title)
	}
	if merged.Description != description {
		t.Fatalf("expected description override %q, got %q", description, merged.Description)
	}
	if merged.Slug != release.Slug {
		t.Fatalf("expected slug preservation %q, got %q", release.Slug, merged.Slug)
	}
	if merged.Tracks[0].Title != trackTitle {
		t.Fatalf("expected track title override %q, got %q", trackTitle, merged.Tracks[0].Title)
	}
	if !merged.Tracks[0].DownloadEnabled {
		t.Fatal("expected track download flag override")
	}
	if merged.Tracks[0].Files.Original.Path != release.Tracks[0].Files.Original.Path {
		t.Fatal("expected audio file path preservation")
	}
	if merged.Tracks[1].Title != release.Tracks[1].Title {
		t.Fatal("expected non-overridden track preservation")
	}
}

func TestApplyNoOpOnEmptyPartial(t *testing.T) {
	t.Parallel()

	release := releaseWithTracks()
	merged, err := Apply(release, models.PendingReleaseEdit{})
	if err != nil {
		t.Fatalf("Apply returned error: %v", err)
	}

	if merged.Title != release.Title {
		t.Fatalf("expected title %q, got %q", release.Title, merged.Title)
	}
	if merged.Tracks[0].Title != release.Tracks[0].Title {
		t.Fatalf("expected track title %q, got %q", release.Tracks[0].Title, merged.Tracks[0].Title)
	}
}

func TestValidateRejectsMembershipAndAudioChanges(t *testing.T) {
	t.Parallel()

	trackCount := 3
	remove := true

	tests := []struct {
		name      string
		edit      models.PendingReleaseEdit
		expectErr error
	}{
		{
			name: "track add by unknown key",
			edit: models.PendingReleaseEdit{Tracks: []models.PendingTrackEdit{{
				TrackKey: "missing-track",
			}}},
			expectErr: ErrTrackMembershipChange,
		},
		{
			name: "track remove marker",
			edit: models.PendingReleaseEdit{Tracks: []models.PendingTrackEdit{{
				TrackKey: "track-key-1",
				Remove:   &remove,
			}}},
			expectErr: ErrTrackMembershipChange,
		},
		{
			name:      "track count override",
			edit:      models.PendingReleaseEdit{TrackCount: &trackCount},
			expectErr: ErrTrackMembershipChange,
		},
		{
			name: "audio replacement",
			edit: models.PendingReleaseEdit{Tracks: []models.PendingTrackEdit{{
				TrackKey: "track-key-1",
				Files:    &models.TrackFiles{Original: models.OriginalFile{Path: "new.wav"}},
			}}},
			expectErr: ErrAudioReplacement,
		},
		{
			name: "missing track key",
			edit: models.PendingReleaseEdit{Tracks: []models.PendingTrackEdit{{
				Title: stringPtr("No Key"),
			}}},
			expectErr: ErrMalformedEdit,
		},
		{
			name: "duplicate track key",
			edit: models.PendingReleaseEdit{Tracks: []models.PendingTrackEdit{
				{TrackKey: "track-key-1"},
				{TrackKey: "track-key-1"},
			}},
			expectErr: ErrMalformedEdit,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.edit, releaseWithTracks().Tracks)
			if !errors.Is(err, tt.expectErr) {
				t.Fatalf("expected %v, got %v", tt.expectErr, err)
			}
		})
	}
}

func TestStageDiscardAndPublishOnto(t *testing.T) {
	t.Parallel()

	title := "Pending Title"
	release := releaseWithTracks()
	stagedRelease, err := Stage(release.Release, release.Tracks, models.PendingReleaseEdit{Title: &title})
	if err != nil {
		t.Fatalf("Stage returned error: %v", err)
	}
	if stagedRelease.Pending == nil || stagedRelease.Pending.Title == nil {
		t.Fatal("expected pending edit on staged release")
	}

	discarded := Discard(stagedRelease)
	if discarded.Pending != nil {
		t.Fatal("expected pending edit to be discarded")
	}

	published, err := PublishOnto(release, models.PendingReleaseEdit{Title: &title})
	if err != nil {
		t.Fatalf("PublishOnto returned error: %v", err)
	}
	if published.Title != title {
		t.Fatalf("expected published title %q, got %q", title, published.Title)
	}
	if published.Pending != nil {
		t.Fatal("expected published release pending edit to be cleared")
	}
}

func releaseWithTracks() models.ReleaseWithTracks {
	return models.ReleaseWithTracks{
		Release: models.Release{
			Title:       "Live at the Roxy",
			Slug:        "live-at-the-roxy",
			Description: "original description",
		},
		Tracks: []models.Track{
			{
				DocumentMeta: models.DocumentMeta{Key: "track-key-1"},
				Title:        "Track One",
				Files: models.TrackFiles{
					Original: models.OriginalFile{Path: "artist_content/a/releases/r/track-one.wav"},
				},
			},
			{
				DocumentMeta: models.DocumentMeta{Key: "track-key-2"},
				Title:        "Track Two",
				Files: models.TrackFiles{
					Original: models.OriginalFile{Path: "artist_content/a/releases/r/track-two.wav"},
				},
			},
		},
	}
}

func stringPtr(value string) *string {
	return &value
}
