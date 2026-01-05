package indexer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/doughughes/claude-marketplace/plugins/conversation-index/internal/shared"
)

// ConversationFile represents metadata about a conversation file
type ConversationFile struct {
	UUID         string
	FilePath     string
	ProjectPath  string
	EncodedPath  string
	LastModified int64 // Unix timestamp in nanoseconds
}

// Scanner scans the filesystem for conversation files
type Scanner struct {
	projectsDir string
}

// NewScanner creates a new file scanner
func NewScanner(projectsDir string) *Scanner {
	return &Scanner{projectsDir: projectsDir}
}

// Scan walks the projects directory and finds all conversation JSONL files
func (s *Scanner) Scan() ([]ConversationFile, error) {
	if _, err := os.Stat(s.projectsDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("projects directory does not exist: %s", s.projectsDir)
	}

	var files []ConversationFile

	// Read all project directories
	entries, err := os.ReadDir(s.projectsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read projects directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		encodedPath := entry.Name()
		projectDir := filepath.Join(s.projectsDir, encodedPath)
		decodedPath := shared.DecodeProjectPath(encodedPath)

		// Find all .jsonl files in this project directory
		conversationFiles, err := s.scanProject(projectDir, decodedPath, encodedPath)
		if err != nil {
			// Log error but continue with other projects
			continue
		}

		files = append(files, conversationFiles...)
	}

	return files, nil
}

// scanProject scans a single project directory for conversation files
func (s *Scanner) scanProject(projectDir, decodedPath, encodedPath string) ([]ConversationFile, error) {
	entries, err := os.ReadDir(projectDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read project directory: %w", err)
	}

	var files []ConversationFile

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".jsonl") {
			continue
		}

		uuid := strings.TrimSuffix(entry.Name(), ".jsonl")
		filePath := filepath.Join(projectDir, entry.Name())

		// Get file modification time
		info, err := os.Stat(filePath)
		if err != nil {
			continue
		}

		files = append(files, ConversationFile{
			UUID:         uuid,
			FilePath:     filePath,
			ProjectPath:  decodedPath,
			EncodedPath:  encodedPath,
			LastModified: info.ModTime().UnixNano(),
		})
	}

	return files, nil
}
