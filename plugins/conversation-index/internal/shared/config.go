package shared

import (
	"os"
	"path/filepath"
)

// Constants for Claude Code directory structure
var (
	ClaudeDir   = filepath.Join(os.Getenv("HOME"), ".claude")
	ProjectsDir = filepath.Join(ClaudeDir, "projects")
	DBPath      = filepath.Join(ClaudeDir, "conversation-index.db")
)
