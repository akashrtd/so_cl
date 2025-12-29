package core

import (
	"testing"
)

func TestGeneratePFP(t *testing.T) {
	feedRef := "@alice.ed25519"

	pfp := GeneratePFP(feedRef)

	if len(pfp) == 0 {
		t.Fatal("PFP is empty")
	}

	pfp2 := GeneratePFP(feedRef)

	if pfp != pfp2 {
		t.Error("PFP is not deterministic for same feedRef")
	}

	pfp3 := GeneratePFP("@bob.ed25519")

	if pfp == pfp3 {
		t.Error("PFP should be different for different feedRef")
	}
}

func TestRenderColoredPFP(t *testing.T) {
	pattern := patterns[0]
	fg := ColorRed
	bg := ColorBlue

	result := renderColoredPFP(pattern, fg, bg)

	if len(result) == 0 {
		t.Fatal("renderColoredPFP returned empty string")
	}

	if !contains(result, ColorReset) {
		t.Error("Result should contain color reset code")
	}

	if !contains(result, fg) && !contains(result, bg) {
		t.Error("Result should contain at least one of the colors")
	}
}

func TestSplitLines(t *testing.T) {
	s := "line1\nline2\nline3"
	lines := splitLines(s)

	if len(lines) != 3 {
		t.Fatalf("Expected 3 lines, got %d", len(lines))
	}

	if lines[0] != "line1" || lines[1] != "line2" || lines[2] != "line3" {
		t.Error("Lines not split correctly")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && containsHelper(s, substr)))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
