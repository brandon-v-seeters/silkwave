package storage

import (
	"context"
	"strings"
	"testing"
)

func TestNewClientRejectsMissingR2Config(t *testing.T) {
	_, err := NewClient(context.Background(), R2Config{})
	if err == nil {
		t.Fatal("expected missing R2 config error")
	}

	message := err.Error()
	for _, want := range []string{
		"R2_ACCOUNT_ID",
		"R2_ACCESS_KEY_ID",
		"R2_SECRET_ACCESS_KEY",
		"R2_BUCKET_NAME",
	} {
		if !strings.Contains(message, want) {
			t.Fatalf("expected error %q to mention %s", message, want)
		}
	}
}

func TestNewClientPresignsWithCompleteR2Config(t *testing.T) {
	client, err := NewClient(context.Background(), R2Config{
		AccountID:       "test-account",
		AccessKeyID:     "test-access-key",
		SecretAccessKey: "test-secret-key",
		BucketName:      "test-bucket",
	})
	if err != nil {
		t.Fatalf("expected complete R2 config to create client: %v", err)
	}

	url, err := client.GetPresignedUploadURL(
		context.Background(),
		"artist_content/artist_framer/releases/release-id/draft/cover.jpg",
		"image/jpeg",
		3600,
	)
	if err != nil {
		t.Fatalf("expected complete R2 config to presign upload URL: %v", err)
	}

	if !strings.Contains(url, "test-account.r2.cloudflarestorage.com") {
		t.Fatalf("expected presigned URL to use R2 account endpoint, got %q", url)
	}
	if !strings.Contains(url, "X-Amz-Credential=test-access-key") {
		t.Fatalf("expected presigned URL to use configured access key, got %q", url)
	}
}
