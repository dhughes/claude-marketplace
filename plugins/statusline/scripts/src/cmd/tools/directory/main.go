package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Input struct {
	Workspaces []string `json:"workspaces"`
}

type Output struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// makeHyperlink creates a terminal hyperlink using OSC 8 escape sequences
// Format: ESC]8;;URL BEL TEXT ESC]8;; BEL
// Using BEL (\a) terminator for better terminal compatibility
func makeHyperlink(url, text string) string {
	return fmt.Sprintf("\033]8;;%s\a%s\033]8;;\a", url, text)
}

func main() {
	// Read input from stdin
	var input Input
	json.NewDecoder(os.Stdin).Decode(&input)

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		// Output empty result if we can't get the directory
		fmt.Println("{}")
		return
	}

	// Determine the primary directory to show (workspace or cwd)
	var primaryDir string
	var secondaryDir string

	if len(input.Workspaces) > 0 && input.Workspaces[0] != "" {
		// Use workspace directory as primary
		primaryDir = input.Workspaces[0]

		// If CWD is different from workspace, note it
		if cwd != primaryDir {
			secondaryDir = cwd
		}
	} else {
		// No workspace info, fall back to CWD
		primaryDir = cwd
	}

	// Get just the directory name (last component of path)
	primaryDirName := filepath.Base(primaryDir)

	// Create file:// URL for the primary directory
	fileURL := fmt.Sprintf("file://%s", primaryDir)

	// Create hyperlink with directory name
	linkedDirName := makeHyperlink(fileURL, primaryDirName)

	// If there's a secondary directory (CWD different from workspace), add it
	var value string
	if secondaryDir != "" {
		secondaryDirName := filepath.Base(secondaryDir)
		value = fmt.Sprintf("%s (in: %s)", linkedDirName, secondaryDirName)
	} else {
		value = linkedDirName
	}

	// Output JSON with label and value
	output := Output{
		Label: "Directory",
		Value: value,
	}

	json.NewEncoder(os.Stdout).Encode(output)
}
