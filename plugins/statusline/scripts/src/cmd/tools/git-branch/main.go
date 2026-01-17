package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Output struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// mutedColor returns text in a muted gray color
func mutedColor(text string) string {
	return fmt.Sprintf("\033[90m%s\033[0m", text)
}

// isGitRepo checks if the current directory is in a git repository
func isGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	err := cmd.Run()
	return err == nil
}

// getCurrentBranch gets the current git branch name
func getCurrentBranch() (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func main() {
	// Read input from stdin (for future use)
	var input map[string]interface{}
	json.NewDecoder(os.Stdin).Decode(&input)

	var value string

	// Check if we're in a git repository
	if !isGitRepo() {
		// Not in a git repo, show N/A in muted color
		value = mutedColor("N/A")
	} else {
		// Get current branch
		branch, err := getCurrentBranch()
		if err != nil || branch == "" {
			// Failed to get branch or detached HEAD, show N/A
			value = mutedColor("N/A")
		} else {
			value = branch
		}
	}

	// Output JSON with label and value
	output := Output{
		Label: "Branch",
		Value: value,
	}

	json.NewEncoder(os.Stdout).Encode(output)
}
