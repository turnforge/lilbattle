package services

import (
	"testing"
)

func TestNormalizeWorldID(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Aruba", "aruba"},
		{"ARUBA", "aruba"},
		{"aruba", "aruba"},
		{"aRuBa", "aruba"},
		{"Desert-Map", "desert-map"},
		{"01bdc3ce", "01bdc3ce"},
		{"MyWorld123", "myworld123"},
		{"", ""},
		{"already-lower", "already-lower"},
		{"UPPER", "upper"},
		{"MixedCase-With-Dashes", "mixedcase-with-dashes"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := NormalizeWorldID(tt.input)
			if result != tt.expected {
				t.Errorf("NormalizeWorldID(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestNormalizeWorldID_Idempotent(t *testing.T) {
	ids := []string{"Aruba", "DESERT", "myWorld", "test-123"}
	for _, id := range ids {
		first := NormalizeWorldID(id)
		second := NormalizeWorldID(first)
		if first != second {
			t.Errorf("NormalizeWorldID is not idempotent: NormalizeWorldID(%q) = %q, NormalizeWorldID(%q) = %q",
				id, first, first, second)
		}
	}
}
