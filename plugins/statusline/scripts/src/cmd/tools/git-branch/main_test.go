package main

import (
	"testing"
)

func TestMutedColor(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected string
	}{
		{
			name:     "N/A text",
			text:     "N/A",
			expected: "\033[90mN/A\033[0m",
		},
		{
			name:     "empty text",
			text:     "",
			expected: "\033[90m\033[0m",
		},
		{
			name:     "with spaces",
			text:     "not available",
			expected: "\033[90mnot available\033[0m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mutedColor(tt.text)
			if result != tt.expected {
				t.Errorf("mutedColor(%q) = %q, want %q", tt.text, result, tt.expected)
			}
		})
	}
}

// Note: isGitRepo and getCurrentBranch are not easily unit testable
// without mocking the exec.Command calls. These should be tested
// via integration tests instead.
