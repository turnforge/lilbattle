package r2

import (
	"testing"
)

func TestValidatePath(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{"valid simple path", "screenshots/worlds/abc/default.png", false},
		{"valid nested path", "a/b/c/d.txt", false},
		{"valid single file", "file.txt", false},
		{"empty path", "", true},
		{"absolute path", "/etc/passwd", true},
		{"directory traversal", "../secrets/key.json", true},
		{"nested traversal", "a/../../etc/passwd", true},
		{"dot-dot only", "..", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePath(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("validatePath(%q) error = %v, wantErr %v", tt.path, err, tt.wantErr)
			}
		})
	}
}
