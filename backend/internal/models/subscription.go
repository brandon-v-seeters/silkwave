package models

import "time"

// SubscriptionStatus represents the status of a subscriber
type SubscriptionStatus string

const (
	SubscriptionStatusActive   SubscriptionStatus = "active"
	SubscriptionStatusInactive SubscriptionStatus = "inactive"
)

// Subscription represents a paid support level offered by an artist
type Subscription struct {
	DocumentMeta `tstype:",extends"`

	ArtistKey                 string    `json:"artistKey" binding:"required"`
	Name                      string    `json:"name" binding:"required"`
	Description               string    `json:"description"`
	Price                     float64   `json:"price" binding:"required"`
	Currency                  string    `json:"currency" binding:"required"`
	SubscriberDiscountPercent *float64  `json:"subscriberDiscountPercent,omitempty"`
	CreatedAt                 time.Time `json:"createdAt"`
}

// CreateSubscriptionRequest represents the request body for creating a subscription
type CreateSubscriptionRequest struct {
	ArtistKey                 string   `json:"artistKey" binding:"required"`
	Name                      string   `json:"name" binding:"required"`
	Description               string   `json:"description"`
	Price                     float64  `json:"price" binding:"required"`
	Currency                  string   `json:"currency" binding:"required"`
	SubscriberDiscountPercent *float64 `json:"subscriberDiscountPercent,omitempty"`
}

// UpdateSubscriptionRequest represents the request body for updating a subscription
type UpdateSubscriptionRequest struct {
	Name                      string   `json:"name,omitempty"`
	Description               string   `json:"description,omitempty"`
	Price                     *float64 `json:"price,omitempty"`
	Currency                  string   `json:"currency,omitempty"`
	SubscriberDiscountPercent *float64 `json:"subscriberDiscountPercent,omitempty"`
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

	ArtistKey                 string             `json:"artistKey"`
	ArtistName                string             `json:"artistName"`
	SubscriberKey             string             `json:"subscriberKey"`
	SubscriptionKey           string             `json:"subscriptionKey"`
	SubscriptionName          string             `json:"subscriptionName"`
	SubscriptionPrice         float64            `json:"subscriptionPrice"`
	SubscriptionCurrency      string             `json:"subscriptionCurrency"`
	SubscriberDiscountPercent *float64           `json:"subscriberDiscountPercent,omitempty"`
	Status                    SubscriptionStatus `json:"status"`
}
