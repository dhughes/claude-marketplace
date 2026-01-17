package main

import (
	"encoding/json"
	"testing"
)

func TestToolOutputJSON(t *testing.T) {
	tests := []struct {
		name     string
		output   ToolOutput
		expected string
	}{
		{
			name: "simple output",
			output: ToolOutput{
				Label: "Directory",
				Value: "test-project",
			},
			expected: `{"label":"Directory","value":"test-project"}`,
		},
		{
			name: "output with special characters",
			output: ToolOutput{
				Label: "Branch",
				Value: "feature/test-123",
			},
			expected: `{"label":"Branch","value":"feature/test-123"}`,
		},
		{
			name: "output with escape sequences",
			output: ToolOutput{
				Label: "Directory",
				Value: "\033]8;;file:///tmp\atmp\033]8;;\a",
			},
			expected: `{"label":"Directory","value":"\u001b]8;;file:///tmp\u0007tmp\u001b]8;;\u0007"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBytes, err := json.Marshal(tt.output)
			if err != nil {
				t.Fatalf("json.Marshal() error = %v", err)
			}

			if string(jsonBytes) != tt.expected {
				t.Errorf("json.Marshal() = %q, want %q", string(jsonBytes), tt.expected)
			}
		})
	}
}

func TestToolOutputUnmarshal(t *testing.T) {
	jsonInput := `{"label":"Test","value":"value123"}`

	var output ToolOutput
	err := json.Unmarshal([]byte(jsonInput), &output)
	if err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if output.Label != "Test" {
		t.Errorf("Label = %q, want %q", output.Label, "Test")
	}

	if output.Value != "value123" {
		t.Errorf("Value = %q, want %q", output.Value, "value123")
	}
}
