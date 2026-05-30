package pending

import (
	"errors"
	"fmt"

	"github.com/brandon-v-seeters/go-silk-wave/internal/models"
)

var (
	ErrTrackMembershipChange = errors.New("pending edit cannot change track membership")
	ErrAudioReplacement      = errors.New("pending edit cannot replace audio files")
	ErrMalformedEdit         = errors.New("malformed pending edit")
)

func Stage(release models.Release, tracks []models.Track, edit models.PendingReleaseEdit) (models.Release, error) {
	if err := Validate(edit, tracks); err != nil {
		return models.Release{}, err
	}

	release.Pending = &edit
	return release, nil
}

func Discard(release models.Release) models.Release {
	release.Pending = nil
	return release
}

func Apply(release models.ReleaseWithTracks, edit models.PendingReleaseEdit) (models.ReleaseWithTracks, error) {
	if err := Validate(edit, release.Tracks); err != nil {
		return models.ReleaseWithTracks{}, err
	}

	merged := cloneReleaseWithTracks(release)
	applyReleaseFields(&merged.Release, edit)
	applyTrackFields(merged.Tracks, edit.Tracks)

	return merged, nil
}

func PublishOnto(release models.ReleaseWithTracks, edit models.PendingReleaseEdit) (models.ReleaseWithTracks, error) {
	merged, err := Apply(release, edit)
	if err != nil {
		return models.ReleaseWithTracks{}, err
	}

	merged.Pending = nil
	return merged, nil
}

func Validate(edit models.PendingReleaseEdit, tracks []models.Track) error {
	if edit.TrackCount != nil || edit.TrackList != nil {
		return ErrTrackMembershipChange
	}

	trackKeys := make(map[string]struct{}, len(tracks))
	for _, track := range tracks {
		trackKeys[track.Key] = struct{}{}
	}

	seen := make(map[string]struct{}, len(edit.Tracks))
	for _, trackEdit := range edit.Tracks {
		if trackEdit.TrackKey == "" {
			return fmt.Errorf("%w: track edit is missing _key", ErrMalformedEdit)
		}
		if _, ok := seen[trackEdit.TrackKey]; ok {
			return fmt.Errorf("%w: duplicate track edit for %s", ErrMalformedEdit, trackEdit.TrackKey)
		}
		seen[trackEdit.TrackKey] = struct{}{}

		if _, ok := trackKeys[trackEdit.TrackKey]; !ok {
			return ErrTrackMembershipChange
		}
		if trackEdit.Remove != nil && *trackEdit.Remove {
			return ErrTrackMembershipChange
		}
		if trackEdit.Files != nil {
			return ErrAudioReplacement
		}
	}

	return nil
}

func cloneReleaseWithTracks(release models.ReleaseWithTracks) models.ReleaseWithTracks {
	clone := release
	if release.Tracks != nil {
		clone.Tracks = append([]models.Track(nil), release.Tracks...)
	}

	return clone
}
