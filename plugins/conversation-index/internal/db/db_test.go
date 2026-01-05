package db

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSQLiteDB_Integration(t *testing.T) {
	// Create a temporary database file
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := Open(dbPath)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	// Test schema initialization
	if err := db.InitSchema(); err != nil {
		t.Fatalf("failed to initialize schema: %v", err)
	}

	// Test saving a conversation
	conv := &Conversation{
		UUID:         "test-uuid-1",
		ProjectPath:  "/Users/test/project",
		EncodedPath:  "-Users-test-project",
		CreatedAt:    time.Now(),
		LastUpdated:  time.Now(),
		MessageCount: 0,
	}

	if err := db.SaveConversation(conv); err != nil {
		t.Fatalf("failed to save conversation: %v", err)
	}

	// Test saving messages
	messages := []Message{
		{
			ConversationUUID: "test-uuid-1",
			Timestamp:        time.Now(),
			Role:             "user",
			Content:          "Hello, this is a test message",
		},
		{
			ConversationUUID: "test-uuid-1",
			Timestamp:        time.Now(),
			Role:             "assistant",
			Content:          "This is a response message",
		},
	}

	if err := db.SaveMessages(messages); err != nil {
		t.Fatalf("failed to save messages: %v", err)
	}

	// Test getting first user message
	firstMsg, err := db.GetFirstUserMessage("test-uuid-1")
	if err != nil {
		t.Fatalf("failed to get first user message: %v", err)
	}

	if firstMsg != "Hello, this is a test message" {
		t.Errorf("expected first message to be 'Hello, this is a test message', got %q", firstMsg)
	}

	// Test index state operations
	state := &IndexState{
		ConversationUUID:  "test-uuid-1",
		LastIndexedLine:   10,
		LastModifiedTime:  time.Now(),
	}

	if err := db.UpdateIndexState(state); err != nil {
		t.Fatalf("failed to update index state: %v", err)
	}

	retrievedState, err := db.GetIndexState("test-uuid-1")
	if err != nil {
		t.Fatalf("failed to get index state: %v", err)
	}

	if retrievedState == nil {
		t.Fatal("expected index state, got nil")
	}

	if retrievedState.LastIndexedLine != 10 {
		t.Errorf("expected last indexed line to be 10, got %d", retrievedState.LastIndexedLine)
	}

	// Test search (FTS5)
	matches, err := db.Search("test message", "all_projects", "", 10)
	if err != nil {
		t.Fatalf("failed to search: %v", err)
	}

	if len(matches) != 1 {
		t.Errorf("expected 1 match, got %d", len(matches))
	}

	if len(matches) > 0 {
		if matches[0].UUID != "test-uuid-1" {
			t.Errorf("expected UUID 'test-uuid-1', got %q", matches[0].UUID)
		}

		if matches[0].Summary != "Hello, this is a test message" {
			t.Errorf("expected summary 'Hello, this is a test message', got %q", matches[0].Summary)
		}
	}

	// Test delete conversation
	if err := db.DeleteConversation("test-uuid-1"); err != nil {
		t.Fatalf("failed to delete conversation: %v", err)
	}

	// Verify messages were deleted
	firstMsg, err = db.GetFirstUserMessage("test-uuid-1")
	if err != nil {
		t.Fatalf("failed to get first user message after delete: %v", err)
	}

	if firstMsg != "" {
		t.Errorf("expected no messages after delete, got %q", firstMsg)
	}
}

func TestSQLiteDB_SearchScoping(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := Open(dbPath)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	if err := db.InitSchema(); err != nil {
		t.Fatalf("failed to initialize schema: %v", err)
	}

	// Create two conversations in different projects
	conv1 := &Conversation{
		UUID:         "project1-conv",
		ProjectPath:  "/Users/test/project1",
		EncodedPath:  "-Users-test-project1",
		CreatedAt:    time.Now(),
		LastUpdated:  time.Now(),
	}

	conv2 := &Conversation{
		UUID:         "project2-conv",
		ProjectPath:  "/Users/test/project2",
		EncodedPath:  "-Users-test-project2",
		CreatedAt:    time.Now(),
		LastUpdated:  time.Now(),
	}

	if err := db.SaveConversation(conv1); err != nil {
		t.Fatalf("failed to save conversation 1: %v", err)
	}

	if err := db.SaveConversation(conv2); err != nil {
		t.Fatalf("failed to save conversation 2: %v", err)
	}

	// Add messages to both
	msg1 := []Message{{
		ConversationUUID: "project1-conv",
		Timestamp:        time.Now(),
		Role:             "user",
		Content:          "database query in project 1",
	}}

	msg2 := []Message{{
		ConversationUUID: "project2-conv",
		Timestamp:        time.Now(),
		Role:             "user",
		Content:          "database query in project 2",
	}}

	if err := db.SaveMessages(msg1); err != nil {
		t.Fatalf("failed to save messages 1: %v", err)
	}

	if err := db.SaveMessages(msg2); err != nil {
		t.Fatalf("failed to save messages 2: %v", err)
	}

	// Search all projects
	allMatches, err := db.Search("database query", "all_projects", "", 10)
	if err != nil {
		t.Fatalf("failed to search all projects: %v", err)
	}

	if len(allMatches) != 2 {
		t.Errorf("expected 2 matches in all projects, got %d", len(allMatches))
	}

	// Search current project only
	currentMatches, err := db.Search("database query", "current_project", "-Users-test-project1", 10)
	if err != nil {
		t.Fatalf("failed to search current project: %v", err)
	}

	if len(currentMatches) != 1 {
		t.Errorf("expected 1 match in current project, got %d", len(currentMatches))
	}

	if len(currentMatches) > 0 && currentMatches[0].UUID != "project1-conv" {
		t.Errorf("expected match from project1-conv, got %q", currentMatches[0].UUID)
	}
}

func TestSQLiteDB_FileCreation(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "new-db.db")

	// Verify file doesn't exist
	if _, err := os.Stat(dbPath); !os.IsNotExist(err) {
		t.Fatal("database file should not exist yet")
	}

	// Open database
	db, err := Open(dbPath)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	// Initialize schema (triggers file creation)
	if err := db.InitSchema(); err != nil {
		t.Fatalf("failed to initialize schema: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(dbPath); err != nil {
		t.Fatalf("database file was not created: %v", err)
	}
}
