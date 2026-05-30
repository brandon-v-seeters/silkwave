package repository

import (
	"context"
	"fmt"

	"github.com/brandon-v-seeters/go-silk-wave/internal/access"
	"github.com/brandon-v-seeters/go-silk-wave/internal/database"
	"github.com/brandon-v-seeters/go-silk-wave/internal/models"
)

type AccessRepository struct {
	db *database.ArangoDB
}

func NewAccessRepository(db *database.ArangoDB) *AccessRepository {
	return &AccessRepository{db: db}
}

func (r *AccessRepository) ReleaseAccess(ctx context.Context, userKey, releaseID string) (access.ReleaseAccess, error) {
	q := /*aql*/ `
		FOR release IN Releases
			FILTER release.id == @releaseId
			LIMIT 1
			LET activeSubscription = LENGTH(
				FOR subscriber IN Subscribers
					FILTER subscriber.artistKey == release.artistKey
					FILTER subscriber.subscriberKey == @userKey
					FILTER subscriber.status == @active
					LIMIT 1
					RETURN 1
			) > 0
			RETURN { purchased: false, activeSubscription: activeSubscription }
	`

	result, err := database.QueryOne[access.ReleaseAccess](ctx, r.db, q, map[string]interface{}{
		"active":    string(models.SubscriptionStatusActive),
		"releaseId": releaseID,
		"userKey":   userKey,
	})
	if err != nil {
		return access.ReleaseAccess{}, fmt.Errorf("get release access: %w", err)
	}

	return *result, nil
}

func (r *AccessRepository) SubscriberDiscount(ctx context.Context, userKey, artistID string) (int, error) {
	q := /*aql*/ `
		FOR artist IN Artists
			FILTER artist.id == @artistId
			LIMIT 1
			LET discounts = (
				FOR subscriber IN Subscribers
					FILTER subscriber.artistKey == artist._key
					FILTER subscriber.subscriberKey == @userKey
					FILTER subscriber.status == @active
					FOR subscription IN Subscriptions
						FILTER subscription._key == subscriber.subscriptionKey
						FILTER subscription.subscriberDiscountPercent != null
						RETURN subscription.subscriberDiscountPercent
			)
			RETURN LENGTH(discounts) == 0 ? 0 : FLOOR(MAX(discounts))
	`

	discount, err := database.QueryOne[int](ctx, r.db, q, map[string]interface{}{
		"active":   string(models.SubscriptionStatusActive),
		"artistId": artistID,
		"userKey":  userKey,
	})
	if err != nil {
		return 0, fmt.Errorf("get subscriber discount: %w", err)
	}

	return *discount, nil
}
