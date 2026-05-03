package models

import "time"

// SubscriptionTier represents the tier of a subscription
type SubscriptionTier string

const (
	SubscriptionTierFree SubscriptionTier = "free"
	SubscriptionTierPaid SubscriptionTier = "paid"
)

// SubscriptionStatus represents the status of a subscriber
type SubscriptionStatus string

const (
	SubscriptionStatusActive            SubscriptionStatus = "active"
	SubscriptionStatusCancelled         SubscriptionStatus = "cancelled"
	SubscriptionStatusPaused            SubscriptionStatus = "paused"
	SubscriptionStatusTrialing          SubscriptionStatus = "trialing"
	SubscriptionStatusIncomplete        SubscriptionStatus = "incomplete"
	SubscriptionStatusIncompleteExpired SubscriptionStatus = "incomplete_expired"
	SubscriptionStatusPastDue           SubscriptionStatus = "past_due"
	SubscriptionStatusUnpaid            SubscriptionStatus = "unpaid"
)

// Subscription represents a subscription plan offered by an artist
type Subscription struct {
	DocumentMeta `tstype:",extends"`

	ArtistKey   string           `json:"artistKey" binding:"required"`
	Name        string           `json:"name" binding:"required"`
	Description string           `json:"description"`
	Price       float64          `json:"price" binding:"required"`
	Currency    string           `json:"currency" binding:"required"`
	Tier        SubscriptionTier `json:"tier" binding:"required"`
	CreatedAt   time.Time        `json:"createdAt"`
}

// CreateSubscriptionRequest represents the request body for creating a subscription
type CreateSubscriptionRequest struct {
	ArtistKey   string           `json:"artistKey" binding:"required"`
	Name        string           `json:"name" binding:"required"`
	Description string           `json:"description"`
	Price       float64          `json:"price" binding:"required"`
	Currency    string           `json:"currency" binding:"required"`
	Tier        SubscriptionTier `json:"tier" binding:"required"`
}

// UpdateSubscriptionRequest represents the request body for updating a subscription
type UpdateSubscriptionRequest struct {
	Name        string           `json:"name,omitempty"`
	Description string           `json:"description,omitempty"`
	Price       *float64         `json:"price,omitempty"`
	Currency    string           `json:"currency,omitempty"`
	Tier        SubscriptionTier `json:"tier,omitempty"`
}

// Subscriber represents a user's subscription to an artist's plan
type Subscriber struct {
	DocumentMeta `tstype:",extends"`

	ArtistKey       string             `json:"artistKey" binding:"required"`
	SubscriberKey   string             `json:"subscriberKey" binding:"required"`
	SubscriptionKey string             `json:"subscriptionKey" binding:"required"`
	Status          SubscriptionStatus `json:"status" binding:"required"`
	CreatedAt       time.Time          `json:"createdAt"`
}

// CreateSubscriberRequest represents the request body for creating a subscriber
type CreateSubscriberRequest struct {
	ArtistKey       string             `json:"artistKey" binding:"required"`
	SubscriberKey   string             `json:"subscriberKey" binding:"required"`
	SubscriptionKey string             `json:"subscriptionKey" binding:"required"`
	Status          SubscriptionStatus `json:"status" binding:"required"`
}

// UpdateSubscriberRequest represents the request body for updating a subscriber
type UpdateSubscriberRequest struct {
	Status SubscriptionStatus `json:"status,omitempty"`
}

// ClientSubscriber represents an enriched subscriber with joined data for client responses
type ClientSubscriber struct {
	DocumentMeta `tstype:",extends"`

	ArtistKey        string             `json:"artistKey"`
	ArtistName       string             `json:"artistName"`
	SubscriberKey    string             `json:"subscriberKey"`
	SubscriptionKey  string             `json:"subscriptionKey"`
	SubscriptionName string             `json:"subscriptionName"`
	Status           SubscriptionStatus `json:"status"`
	Tier             SubscriptionTier   `json:"tier"`
}
