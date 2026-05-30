package storage

import (
	"context"
	"fmt"
)

func (c *Client) PublishPendingCover(ctx context.Context, artistKey, releaseId, stagedCoverKey string) error {
	if stagedCoverKey == "" {
		return fmt.Errorf("pending cover key is required")
	}

	publishedKey, err := c.resolver.PendingCoverPublishedKey(artistKey, releaseId, stagedCoverKey)
	if err != nil {
		return err
	}

	if err := c.Move(ctx, stagedCoverKey, publishedKey); err != nil {
		return fmt.Errorf("publish pending cover: %w", err)
	}

	return nil
}
