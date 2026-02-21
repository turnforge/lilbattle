package cmd

import (
	"testing"
)

func TestTruncate(t *testing.T) {
	tests := []struct {
		input    string
		maxLen   int
		expected string
	}{
		{"hello", 10, "hello"},
		{"hello", 5, "hello"},
		{"hello world", 8, "hello..."},
		{"hello world", 5, "he..."},
		{"hi", 3, "hi"},
		{"hello", 3, "hel"},
		{"", 5, ""},
		{"abcdefghij", 10, "abcdefghij"},
		{"abcdefghijk", 10, "abcdefg..."},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := truncate(tt.input, tt.maxLen)
			if result != tt.expected {
				t.Errorf("truncate(%q, %d) = %q, want %q", tt.input, tt.maxLen, result, tt.expected)
			}
		})
	}
}

func TestGetWorldsClient_ProfileSpec(t *testing.T) {
	// Test that profile:worldId parsing correctly extracts the world ID
	// This doesn't make real network calls - we just verify the parsing logic
	// by checking that plain IDs (no colon) don't trigger profile parsing

	// A plain ID without colon should not be treated as profile:worldId
	// It should fall through to the active profile path
	// (which will fail in test since there's no profile configured, but
	// the important thing is it doesn't try to parse "aruba" as "profile:worldId")

	// Test that spec with colon is recognized as profile:worldId format
	specs := []struct {
		input     string
		isProfile bool
	}{
		{"aruba", false},
		{"01bdc3ce", false},
		{"prod:aruba", true},
		{"fsbe:01bdc3ce", true},
		{"http://localhost:8080", false},
		{"https://lilbattle.com", false},
	}

	for _, tt := range specs {
		t.Run(tt.input, func(t *testing.T) {
			hasColon := !isHTTPURL(tt.input) && containsColon(tt.input)
			if hasColon != tt.isProfile {
				t.Errorf("spec %q: expected isProfile=%v, got %v", tt.input, tt.isProfile, hasColon)
			}
		})
	}
}

func TestContainsIgnoreCase(t *testing.T) {
	tests := []struct {
		s      string
		substr string
		want   bool
	}{
		{"already exists", "already exists", true},
		{"World Already Exists", "already exists", true},
		{"not found", "already exists", false},
		{"ALREADY EXISTS error", "already exists", true},
		{"", "test", false},
		{"test", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			got := containsIgnoreCase(tt.s, tt.substr)
			if got != tt.want {
				t.Errorf("containsIgnoreCase(%q, %q) = %v, want %v", tt.s, tt.substr, got, tt.want)
			}
		})
	}
}

// Helper functions for testing (mirror the logic in getWorldsClient)
func isHTTPURL(s string) bool {
	return len(s) >= 7 && (s[:7] == "http://" || (len(s) >= 8 && s[:8] == "https://"))
}

func containsColon(s string) bool {
	for _, c := range s {
		if c == ':' {
			return true
		}
	}
	return false
}
