package indexer

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/doughughes/claude-marketplace/plugins/conversation-index/internal/db"
)

func TestIndexer_IncrementalIndexing(t *testing.T) {
	mockDB := db.NewMock()
	tmpDir := t.TempDir()

	// Create a test JSONL file
	conversationPath := filepath.Join(tmpDir, "test-uuid.jsonl")
	content := `{"type":"user","timestamp":"2026-01-05T10:00:00Z","message":{"content":"First message"},"cwd":"/test/project"}
{"type":"assistant","timestamp":"2026-01-05T10:00:01Z","message":{"content":[{"type":"text","text":"First response"}]}}
{"type":"user","timestamp":"2026-01-05T10:00:02Z","message":{"content":"Second message"}}
`

	if err := os.WriteFile(conversationPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// First index
	indexer := NewIndexer(mockDB, tmpDir)

	file := ConversationFile{
		UUID:         "test-uuid",
		FilePath:     conversationPath,
		ProjectPath:  "/test/project",
		EncodedPath:  "-test-project",
		LastModified: time.Now().UnixNano(),
	}

	indexed, skipped, err := indexer.indexConversation(file)
	if err != nil {
		t.Fatalf("failed to index conversation: %v", err)
	}

	if skipped {
		t.Error("expected conversation to be indexed, not skipped")
	}

	if indexed != 3 {
		t.Errorf("expected 3 messages indexed, got %d", indexed)
	}

	// Verify messages were saved
	messages := mockDB.GetMessages("test-uuid")
	if len(messages) != 3 {
		t.Errorf("expected 3 messages in database, got %d", len(messages))
	}

	// Test incremental update: add another line
	content += `{"type":"user","timestamp":"2026-01-05T10:00:03Z","message":{"content":"Third message"}}
`
	if err := os.WriteFile(conversationPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to update test file: %v", err)
	}

	file.LastModified = time.Now().UnixNano()

	indexed, skipped, err = indexer.indexConversation(file)
	if err != nil {
		t.Fatalf("failed to index updated conversation: %v", err)
	}

	if skipped {
		t.Error("expected conversation to be indexed, not skipped")
	}

	if indexed != 1 {
		t.Errorf("expected 1 new message indexed, got %d", indexed)
	}

	// Verify total messages
	messages = mockDB.GetMessages("test-uuid")
	if len(messages) != 4 {
		t.Errorf("expected 4 total messages, got %d", len(messages))
	}
}

func TestIndexer_RollbackDetection(t *testing.T) {
	mockDB := db.NewMock()
	tmpDir := t.TempDir()

	conversationPath := filepath.Join(tmpDir, "test-uuid.jsonl")
	content := `{"type":"user","timestamp":"2026-01-05T10:00:00Z","message":{"content":"Message 1"},"cwd":"/test"}
{"type":"user","timestamp":"2026-01-05T10:00:01Z","message":{"content":"Message 2"}}
{"type":"user","timestamp":"2026-01-05T10:00:02Z","message":{"content":"Message 3"}}
`

	if err := os.WriteFile(conversationPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	indexer := NewIndexer(mockDB, tmpDir)

	file := ConversationFile{
		UUID:         "test-uuid",
		FilePath:     conversationPath,
		ProjectPath:  "/test",
		EncodedPath:  "-test",
		LastModified: time.Now().UnixNano(),
	}

	// Index all 3 messages
	indexed, _, err := indexer.indexConversation(file)
	if err != nil {
		t.Fatalf("failed to index conversation: %v", err)
	}

	if indexed != 3 {
		t.Errorf("expected 3 messages indexed, got %d", indexed)
	}

	// Simulate rollback: remove last line
	shortContent := `{"type":"user","timestamp":"2026-01-05T10:00:00Z","message":{"content":"Message 1"},"cwd":"/test"}
`
	if err := os.WriteFile(conversationPath, []byte(shortContent), 0644); err != nil {
		t.Fatalf("failed to update test file: %v", err)
	}

	file.LastModified = time.Now().UnixNano()

	// Re-index should detect rollback and reindex from scratch
	indexed, _, err = indexer.indexConversation(file)
	if err != nil {
		t.Fatalf("failed to index after rollback: %v", err)
	}

	if indexed != 1 {
		t.Errorf("expected 1 message after rollback reindex, got %d", indexed)
	}

	// Verify only 1 message remains
	messages := mockDB.GetMessages("test-uuid")
	if len(messages) != 1 {
		t.Errorf("expected 1 message after rollback, got %d", len(messages))
	}
}

func TestIndexer_SkipUnchanged(t *testing.T) {
	mockDB := db.NewMock()
	tmpDir := t.TempDir()

	conversationPath := filepath.Join(tmpDir, "test-uuid.jsonl")
	content := `{"type":"user","timestamp":"2026-01-05T10:00:00Z","message":{"content":"Test message"},"cwd":"/test"}
`

	if err := os.WriteFile(conversationPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	fileInfo, _ := os.Stat(conversationPath)
	lastModified := fileInfo.ModTime().UnixNano()

	indexer := NewIndexer(mockDB, tmpDir)

	file := ConversationFile{
		UUID:         "test-uuid",
		FilePath:     conversationPath,
		ProjectPath:  "/test",
		EncodedPath:  "-test",
		LastModified: lastModified,
	}

	// First index
	indexed, skipped, err := indexer.indexConversation(file)
	if err != nil {
		t.Fatalf("failed to index conversation: %v", err)
	}

	if indexed != 1 {
		t.Errorf("expected 1 message indexed, got %d", indexed)
	}

	// Second index with same modification time should skip
	indexed, skipped, err = indexer.indexConversation(file)
	if err != nil {
		t.Fatalf("failed to index conversation second time: %v", err)
	}

	if !skipped {
		t.Error("expected conversation to be skipped on second index")
	}

	if indexed != 0 {
		t.Errorf("expected 0 messages indexed on skip, got %d", indexed)
	}
}
