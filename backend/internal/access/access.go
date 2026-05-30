package access

import (
	"context"
)

type ReleaseAccess struct {
	Purchased          bool `json:"purchased"`
	ActiveSubscription bool `json:"activeSubscription"`
}

type Source interface {
	ReleaseAccess(ctx context.Context, userKey, releaseID string) (ReleaseAccess, error)
	SubscriberDiscount(ctx context.Context, userKey, artistID string) (int, error)
}

func CanStream(ctx context.Context, source Source, userKey, releaseID string) (bool, error) {
	if userKey == "" {
		return false, nil
	}

	access, err := source.ReleaseAccess(ctx, userKey, releaseID)
	if err != nil {
		return false, err
	}

	return access.Purchased || access.ActiveSubscription, nil
}

func CanDownload(ctx context.Context, source Source, userKey, releaseID string) (bool, error) {
	if userKey == "" {
		return false, nil
	}

	access, err := source.ReleaseAccess(ctx, userKey, releaseID)
	if err != nil {
		return false, err
	}

	return access.Purchased, nil
}

func SubscriberDiscount(ctx context.Context, source Source, userKey, artistID string) (int, error) {
	if userKey == "" {
		return 0, nil
	}

	return source.SubscriberDiscount(ctx, userKey, artistID)
}
