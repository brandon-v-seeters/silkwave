package config

import "testing"

func TestLoadUsesDocumentedR2AccountID(t *testing.T) {
	t.Setenv("R2_ACCOUNT_ID", "documented-account")
	t.Setenv("CF_ACCOUNT_ID", "legacy-account")

	cfg := Load()

	if cfg.R2AccountID != "documented-account" {
		t.Fatalf("expected documented R2 account ID, got %q", cfg.R2AccountID)
	}
}

func TestLoadFallsBackToLegacyCFAccountID(t *testing.T) {
	t.Setenv("R2_ACCOUNT_ID", "")
	t.Setenv("CF_ACCOUNT_ID", "legacy-account")

	cfg := Load()

	if cfg.R2AccountID != "legacy-account" {
		t.Fatalf("expected legacy CF account ID fallback, got %q", cfg.R2AccountID)
	}
}
