package cli

import "testing"

func TestResolveServerURLUsesFlagWhenChanged(t *testing.T) {
	got := resolveServerURL(
		"http://localhost:3000/v1/api",
		"http://localhost:3100/v1/api",
		true,
	)

	if got != "http://localhost:3100/v1/api" {
		t.Fatalf("expected explicit server flag to win, got %q", got)
	}
}

func TestResolveServerURLUsesConfigWhenFlagNotChanged(t *testing.T) {
	got := resolveServerURL(
		"http://localhost:3200/v1/api",
		defaultServerURL,
		false,
	)

	if got != "http://localhost:3200/v1/api" {
		t.Fatalf("expected config server to win, got %q", got)
	}
}

func TestResolveServerURLFallsBackToFlagValue(t *testing.T) {
	got := resolveServerURL("", defaultServerURL, false)

	if got != defaultServerURL {
		t.Fatalf("expected fallback server %q, got %q", defaultServerURL, got)
	}
}
