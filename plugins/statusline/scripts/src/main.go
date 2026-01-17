package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type ToolOutput struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

func main() {
	// Read JSON input from stdin (Claude Code provides this)
	var input map[string]interface{}
	inputBytes := []byte("{}")
	if err := json.NewDecoder(os.Stdin).Decode(&input); err == nil {
		// Re-encode input to pass to tools
		inputBytes, _ = json.Marshal(input)
	}

	// Get the directory where this binary is located
	execPath, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding executable path: %v\n", err)
		os.Exit(1)
	}
	binDir := filepath.Dir(execPath)

	// Determine platform for tool binary names
	platform := fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH)

	// List of tools to run (hardcoded for now, will be config-driven later)
	tools := []string{"status-directory", "status-git-branch"}

	var outputs []ToolOutput

	// Execute each tool
	for _, toolName := range tools {
		toolBinary := filepath.Join(binDir, fmt.Sprintf("%s-%s", toolName, platform))

		// Check if binary exists
		if _, err := os.Stat(toolBinary); os.IsNotExist(err) {
			// Tool doesn't exist, skip silently
			continue
		}

		// Execute the tool
		cmd := exec.Command(toolBinary)
		cmd.Stdin = bytes.NewReader(inputBytes)

		output, err := cmd.Output()
		if err != nil {
			// Tool failed, skip it
			continue
		}

		// Parse tool output
		var toolOutput ToolOutput
		if err := json.Unmarshal(output, &toolOutput); err != nil {
			// Invalid JSON, skip
			continue
		}

		// Skip empty outputs
		if toolOutput.Label == "" {
			continue
		}

		outputs = append(outputs, toolOutput)
	}

	// If no outputs, exit silently
	if len(outputs) == 0 {
		return
	}

	// Calculate maximum label length for alignment
	maxLabelLen := 0
	for _, out := range outputs {
		if len(out.Label) > maxLabelLen {
			maxLabelLen = len(out.Label)
		}
	}

	// Format and print output
	var result strings.Builder
	for _, out := range outputs {
		// Pad label to max length, add colon, then value
		result.WriteString(fmt.Sprintf("%-*s:   %s\n", maxLabelLen, out.Label, out.Value))
	}

	// Print without trailing newline (statusline adds its own)
	fmt.Print(strings.TrimSuffix(result.String(), "\n"))
}
