package main

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
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

func TestParseInput(t *testing.T) {
	tests := []struct {
		name          string
		inputJSON     string
		wantID        string
		wantDisplay   string
	}{
		{
			name:          "full model info",
			inputJSON:     `{"model":{"id":"claude-opus-4-1","display_name":"Opus"}}`,
			wantID:        "claude-opus-4-1",
			wantDisplay:   "Opus",
		},
		{
			name:          "only id",
			inputJSON:     `{"model":{"id":"claude-sonnet-4-20250514"}}`,
			wantID:        "claude-sonnet-4-20250514",
			wantDisplay:   "",
		},
		{
			name:          "empty model",
			inputJSON:     `{"model":{}}`,
			wantID:        "",
			wantDisplay:   "",
		},
		{
			name:          "no model field",
			inputJSON:     `{}`,
			wantID:        "",
			wantDisplay:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var input Input
			err := json.NewDecoder(strings.NewReader(tt.inputJSON)).Decode(&input)
			if err != nil {
				t.Fatalf("failed to decode input: %v", err)
			}

			if input.Model.ID != tt.wantID {
				t.Errorf("Model.ID = %q, want %q", input.Model.ID, tt.wantID)
			}
			if input.Model.DisplayName != tt.wantDisplay {
				t.Errorf("Model.DisplayName = %q, want %q", input.Model.DisplayName, tt.wantDisplay)
			}
		})
	}
}

func TestMainOutput(t *testing.T) {
	tests := []struct {
		name        string
		inputJSON   string
		wantLabel   string
		wantValue   string
		wantMuted   bool
	}{
		{
			name:        "displays display_name when available",
			inputJSON:   `{"model":{"id":"claude-opus-4-1","display_name":"Opus"}}`,
			wantLabel:   "Model",
			wantValue:   "Opus",
			wantMuted:   false,
		},
		{
			name:        "falls back to id when no display_name",
			inputJSON:   `{"model":{"id":"claude-sonnet-4-20250514"}}`,
			wantLabel:   "Model",
			wantValue:   "claude-sonnet-4-20250514",
			wantMuted:   false,
		},
		{
			name:        "shows N/A when no model info",
			inputJSON:   `{}`,
			wantLabel:   "Model",
			wantValue:   "N/A",
			wantMuted:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldStdin := os.Stdin
			oldStdout := os.Stdout
			defer func() {
				os.Stdin = oldStdin
				os.Stdout = oldStdout
			}()

			r, w, _ := os.Pipe()
			os.Stdin = r
			go func() {
				w.Write([]byte(tt.inputJSON))
				w.Close()
			}()

			outR, outW, _ := os.Pipe()
			os.Stdout = outW

			main()

			outW.Close()
			var buf bytes.Buffer
			buf.ReadFrom(outR)

			var output Output
			if err := json.Unmarshal(buf.Bytes(), &output); err != nil {
				t.Fatalf("failed to unmarshal output: %v", err)
			}

			if output.Label != tt.wantLabel {
				t.Errorf("Label = %q, want %q", output.Label, tt.wantLabel)
			}

			if tt.wantMuted {
				expected := mutedColor(tt.wantValue)
				if output.Value != expected {
					t.Errorf("Value = %q, want %q (muted)", output.Value, expected)
				}
			} else {
				if output.Value != tt.wantValue {
					t.Errorf("Value = %q, want %q", output.Value, tt.wantValue)
				}
			}
		})
	}
}
