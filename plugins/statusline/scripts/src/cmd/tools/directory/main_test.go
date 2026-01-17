package main

import (
	"encoding/json"
	"testing"
)

func TestMakeHyperlink(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		text     string
		expected string
	}{
		{
			name:     "simple file URL",
			url:      "file:///Users/test/project",
			text:     "project",
			expected: "\033]8;;file:///Users/test/project\aproject\033]8;;\a",
		},
		{
			name:     "URL with spaces",
			url:      "file:///Users/test/my project",
			text:     "my project",
			expected: "\033]8;;file:///Users/test/my project\amy project\033]8;;\a",
		},
		{
			name:     "empty text",
			url:      "file:///tmp",
			text:     "",
			expected: "\033]8;;file:///tmp\a\033]8;;\a",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := makeHyperlink(tt.url, tt.text)
			if result != tt.expected {
				t.Errorf("makeHyperlink(%q, %q) = %q, want %q", tt.url, tt.text, result, tt.expected)
			}
		})
	}
}

func TestInputParsing(t *testing.T) {
	tests := []struct {
		name     string
		jsonInput string
		wantWorkspaces []string
	}{
		{
			name:     "single workspace",
			jsonInput: `{"workspaces": ["/Users/test/project"]}`,
			wantWorkspaces: []string{"/Users/test/project"},
		},
		{
			name:     "multiple workspaces",
			jsonInput: `{"workspaces": ["/Users/test/project1", "/Users/test/project2"]}`,
			wantWorkspaces: []string{"/Users/test/project1", "/Users/test/project2"},
		},
		{
			name:     "empty workspaces",
			jsonInput: `{"workspaces": []}`,
			wantWorkspaces: []string{},
		},
		{
			name:     "no workspaces field",
			jsonInput: `{}`,
			wantWorkspaces: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var input Input
			err := json.Unmarshal([]byte(tt.jsonInput), &input)
			if err != nil {
				t.Fatalf("json.Unmarshal() error = %v", err)
			}

			if len(input.Workspaces) != len(tt.wantWorkspaces) {
				t.Errorf("len(Workspaces) = %d, want %d", len(input.Workspaces), len(tt.wantWorkspaces))
			}

			for i, workspace := range input.Workspaces {
				if workspace != tt.wantWorkspaces[i] {
					t.Errorf("Workspaces[%d] = %q, want %q", i, workspace, tt.wantWorkspaces[i])
				}
			}
		})
	}
}
