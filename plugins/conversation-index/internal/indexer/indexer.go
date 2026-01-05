package indexer

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/doughughes/claude-marketplace/plugins/conversation-index/internal/db"
)

// Indexer coordinates the indexing of conversation files
type Indexer struct {
	db      db.DB
	parser  *Parser
	scanner *Scanner
}

// NewIndexer creates a new indexer
func NewIndexer(database db.DB, projectsDir string) *Indexer {
	return &Indexer{
		db:      database,
		parser:  NewParser(),
		scanner: NewScanner(projectsDir),
	}
}

// IndexStats holds statistics about the indexing run
type IndexStats struct {
	TotalIndexed       int
	TotalSkipped       int
	TotalConversations int
}

// IndexAll indexes all conversations, optionally doing a full reindex
func (idx *Indexer) IndexAll(fullReindex bool) error {
	// Initialize schema
	if err := idx.db.InitSchema(); err != nil {
		return fmt.Errorf("failed to initialize schema: %w", err)
	}

	if fullReindex {
		log.Println("Performing full reindex...")
		if err := idx.db.TruncateAll(); err != nil {
			return fmt.Errorf("failed to truncate database: %w", err)
		}
	}

	startTime := time.Now()
	stats := &IndexStats{}

	// Scan for conversation files
	files, err := idx.scanner.Scan()
	if err != nil {
		return fmt.Errorf("failed to scan conversations: %w", err)
	}

	// Index each conversation
	for _, file := range files {
		indexed, skipped, err := idx.indexConversation(file)
		if err != nil {
			log.Printf("Failed to index conversation %s: %v", file.UUID, err)
			continue
		}

		stats.TotalConversations++
		if skipped {
			stats.TotalSkipped++
		} else {
			stats.TotalIndexed += indexed
		}
	}

	elapsed := time.Since(startTime)
	log.Printf("Indexed %d messages from %d conversations (%d skipped) in %dms",
		stats.TotalIndexed, stats.TotalConversations, stats.TotalSkipped, elapsed.Milliseconds())

	return nil
}

// indexConversation indexes a single conversation file
func (idx *Indexer) indexConversation(file ConversationFile) (indexed int, skipped bool, err error) {
	// Get index state
	state, err := idx.db.GetIndexState(file.UUID)
	if err != nil {
		return 0, false, fmt.Errorf("failed to get index state: %w", err)
	}

	// Check if file has been modified
	lastModified := time.Unix(0, file.LastModified)
	if state != nil && state.LastModifiedTime.Equal(lastModified) {
		return 0, true, nil
	}

	// Read file lines
	lines, err := idx.readLines(file.FilePath)
	if err != nil {
		return 0, false, fmt.Errorf("failed to read file: %w", err)
	}

	if len(lines) == 0 {
		return 0, true, nil
	}

	// Detect rollback: file has fewer lines than last indexed
	if state != nil && len(lines) < state.LastIndexedLine {
		log.Printf("Rollback detected for conversation %s (was %d lines, now %d lines)",
			file.UUID, state.LastIndexedLine, len(lines))

		// Clear this conversation's data
		if err := idx.db.DeleteConversation(file.UUID); err != nil {
			return 0, false, fmt.Errorf("failed to delete conversation: %w", err)
		}

		if err := idx.db.DeleteIndexState(file.UUID); err != nil {
			return 0, false, fmt.Errorf("failed to delete index state: %w", err)
		}

		// Re-index from scratch
		state = nil
	}

	// Determine starting line
	startLine := 0
	if state != nil && len(lines) >= state.LastIndexedLine {
		startLine = state.LastIndexedLine
	}

	if len(lines) <= startLine {
		return 0, true, nil
	}

	// Get or create conversation record
	var createdAt time.Time
	var actualProjectPath string

	if state == nil {
		// Extract metadata from first line
		firstTimestamp, err := idx.parser.GetTimestamp(lines[0])
		if err == nil {
			createdAt = firstTimestamp
		} else {
			createdAt = time.Now()
		}

		actualCWD, err := idx.parser.GetCWD(lines[0])
		if err == nil && actualCWD != "" {
			actualProjectPath = actualCWD
		} else {
			actualProjectPath = file.ProjectPath
		}
	} else {
		// Use existing values (would need to query from DB in full implementation)
		createdAt = time.Now()
		actualProjectPath = file.ProjectPath
	}

	// Save conversation record
	conv := &db.Conversation{
		UUID:         file.UUID,
		ProjectPath:  actualProjectPath,
		EncodedPath:  file.EncodedPath,
		CreatedAt:    createdAt,
		LastUpdated:  lastModified,
		MessageCount: 0, // Will be updated by SaveMessages
	}

	if err := idx.db.SaveConversation(conv); err != nil {
		return 0, false, fmt.Errorf("failed to save conversation: %w", err)
	}

	// Index new lines
	var allMessages []db.Message
	for i := startLine; i < len(lines); i++ {
		messages, err := idx.parser.ParseLine(lines[i])
		if err != nil {
			// Skip invalid lines
			continue
		}

		for _, msg := range messages {
			msg.ConversationUUID = file.UUID
			allMessages = append(allMessages, msg)
		}
	}

	// Save messages in a batch
	if len(allMessages) > 0 {
		if err := idx.db.SaveMessages(allMessages); err != nil {
			return 0, false, fmt.Errorf("failed to save messages: %w", err)
		}
	}

	// Update index state
	newState := &db.IndexState{
		ConversationUUID:  file.UUID,
		LastIndexedLine:   len(lines),
		LastModifiedTime:  lastModified,
	}

	if err := idx.db.UpdateIndexState(newState); err != nil {
		return 0, false, fmt.Errorf("failed to update index state: %w", err)
	}

	return len(allMessages), false, nil
}

// readLines reads all non-empty lines from a file
func (idx *Indexer) readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)

	// Increase buffer size to handle large JSONL lines (default is 64KB)
	// Some assistant responses with code can exceed this
	const maxCapacity = 10 * 1024 * 1024 // 10MB per line
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
