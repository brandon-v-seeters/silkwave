package access

import (
	"context"
	"errors"
	"testing"
)

func TestCanStream(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name     string
		access   ReleaseAccess
		expected bool
	}{
		{name: "subscribed user", access: ReleaseAccess{ActiveSubscription: true}, expected: true},
		{name: "subscribed but inactive", access: ReleaseAccess{}, expected: false},
		{name: "non-subscriber", access: ReleaseAccess{}, expected: false},
		{name: "purchaser without subscription", access: ReleaseAccess{Purchased: true}, expected: true},
		{name: "purchaser with subscription", access: ReleaseAccess{Purchased: true, ActiveSubscription: true}, expected: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			allowed, err := CanStream(ctx, fakeSource{releaseAccess: tt.access}, "user-key", "release-id")
			if err != nil {
				t.Fatalf("CanStream returned error: %v", err)
			}
			if allowed != tt.expected {
				t.Fatalf("expected %t, got %t", tt.expected, allowed)
			}
		})
	}
}

func TestCanDownload(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name     string
		access   ReleaseAccess
		expected bool
	}{
		{name: "purchaser", access: ReleaseAccess{Purchased: true}, expected: true},
		{name: "active subscriber without purchase", access: ReleaseAccess{ActiveSubscription: true}, expected: false},
		{name: "inactive subscriber without purchase", access: ReleaseAccess{}, expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			allowed, err := CanDownload(ctx, fakeSource{releaseAccess: tt.access}, "user-key", "release-id")
			if err != nil {
				t.Fatalf("CanDownload returned error: %v", err)
			}
			if allowed != tt.expected {
				t.Fatalf("expected %t, got %t", tt.expected, allowed)
			}
		})
	}
}

func TestSubscriberDiscount(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name     string
		discount int
		expected int
	}{
		{name: "active subscriber discount", discount: 20, expected: 20},
		{name: "inactive or no subscription", discount: 0, expected: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			discount, err := SubscriberDiscount(ctx, fakeSource{discount: tt.discount}, "user-key", "artist-id")
			if err != nil {
				t.Fatalf("SubscriberDiscount returned error: %v", err)
			}
			if discount != tt.expected {
				t.Fatalf("expected %d, got %d", tt.expected, discount)
			}
		})
	}
}

func TestAccessPropagatesSourceErrors(t *testing.T) {
	expected := errors.New("source failed")
	source := fakeSource{err: expected}

	if _, err := CanStream(context.Background(), source, "user-key", "release-id"); !errors.Is(err, expected) {
		t.Fatalf("expected source error, got %v", err)
	}
}

type fakeSource struct {
	releaseAccess ReleaseAccess
	discount      int
	err           error
}

func (f fakeSource) ReleaseAccess(context.Context, string, string) (ReleaseAccess, error) {
	if f.err != nil {
		return ReleaseAccess{}, f.err
	}

	return f.releaseAccess, nil
}

func (f fakeSource) SubscriberDiscount(context.Context, string, string) (int, error) {
	if f.err != nil {
		return 0, f.err
	}

	return f.discount, nil
}
