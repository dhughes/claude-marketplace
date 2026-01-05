package shared

import "testing"

func TestEncodeProjectPath(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "absolute path with leading slash",
			input: "/Users/foo/bar",
			want:  "-Users-foo-bar",
		},
		{
			name:  "relative path",
			input: "foo/bar",
			want:  "foo-bar",
		},
		{
			name:  "single directory",
			input: "/Users",
			want:  "-Users",
		},
		{
			name:  "empty path",
			input: "",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EncodeProjectPath(tt.input)
			if got != tt.want {
				t.Errorf("EncodeProjectPath(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestDecodeProjectPath(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "encoded path",
			input: "-Users-foo-bar",
			want:  "/Users/foo/bar",
		},
		{
			name:  "no leading dash",
			input: "foo-bar",
			want:  "foo-bar",
		},
		{
			name:  "single encoded directory",
			input: "-Users",
			want:  "/Users",
		},
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DecodeProjectPath(tt.input)
			if got != tt.want {
				t.Errorf("DecodeProjectPath(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestRoundTrip(t *testing.T) {
	// Note: Round trip is lossy for paths with dashes in directory names
	// This matches the Node.js implementation behavior (see indexer.sh line 86)
	tests := []struct {
		name     string
		input    string
		expected string // What we expect after round trip (may differ if path has dashes)
	}{
		{
			name:     "simple path without dashes",
			input:    "/Users/foo/bar",
			expected: "/Users/foo/bar",
		},
		{
			name:     "path with dashes is lossy",
			input:    "/Users/doughughes/code/doug/claude-marketplace",
			expected: "/Users/doughughes/code/doug/claude/marketplace", // Dash becomes slash
		},
		{
			name:     "tmp test path",
			input:    "/tmp/test",
			expected: "/tmp/test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded := EncodeProjectPath(tt.input)
			decoded := DecodeProjectPath(encoded)
			if decoded != tt.expected {
				t.Errorf("Round trip: %q -> %q -> %q, expected %q", tt.input, encoded, decoded, tt.expected)
			}
		})
	}
}
