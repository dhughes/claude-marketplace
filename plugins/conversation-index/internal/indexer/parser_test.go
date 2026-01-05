package indexer

import (
	"testing"
)

func TestParser_ParseLine_UserMessage(t *testing.T) {
	parser := NewParser()

	input := `{"type":"user","timestamp":"2026-01-05T15:32:32.836Z","message":{"content":"Hello, this is a test"},"cwd":"/Users/test"}`

	messages, err := parser.ParseLine(input)
	if err != nil {
		t.Fatalf("failed to parse line: %v", err)
	}

	if len(messages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(messages))
	}

	msg := messages[0]
	if msg.Role != "user" {
		t.Errorf("expected role 'user', got %q", msg.Role)
	}

	if msg.Content != "Hello, this is a test" {
		t.Errorf("expected content 'Hello, this is a test', got %q", msg.Content)
	}
}

func TestParser_ParseLine_AssistantText(t *testing.T) {
	parser := NewParser()

	input := `{"type":"assistant","timestamp":"2026-01-05T15:32:33.000Z","message":{"content":[{"type":"text","text":"Here is my response"}]}}`

	messages, err := parser.ParseLine(input)
	if err != nil {
		t.Fatalf("failed to parse line: %v", err)
	}

	if len(messages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(messages))
	}

	msg := messages[0]
	if msg.Role != "assistant" {
		t.Errorf("expected role 'assistant', got %q", msg.Role)
	}

	if msg.Content != "Here is my response" {
		t.Errorf("expected content 'Here is my response', got %q", msg.Content)
	}
}

func TestParser_ParseLine_ToolUse(t *testing.T) {
	parser := NewParser()

	input := `{"type":"assistant","timestamp":"2026-01-05T15:32:34.000Z","message":{"content":[{"type":"tool_use","name":"Bash","input":{"command":"ls -la","description":"List files"}}]}}`

	messages, err := parser.ParseLine(input)
	if err != nil {
		t.Fatalf("failed to parse line: %v", err)
	}

	if len(messages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(messages))
	}

	msg := messages[0]
	if msg.Role != "tool" {
		t.Errorf("expected role 'tool', got %q", msg.Role)
	}

	expectedContent := "Tool: Bash List files Command: ls -la"
	if msg.Content != expectedContent {
		t.Errorf("expected content %q, got %q", expectedContent, msg.Content)
	}
}

func TestParser_ParseLine_MultipleContent(t *testing.T) {
	parser := NewParser()

	input := `{"type":"assistant","timestamp":"2026-01-05T15:32:35.000Z","message":{"content":[{"type":"text","text":"Let me check that"},{"type":"tool_use","name":"Read","input":{"file_path":"/path/to/file.go"}}]}}`

	messages, err := parser.ParseLine(input)
	if err != nil {
		t.Fatalf("failed to parse line: %v", err)
	}

	if len(messages) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(messages))
	}

	// First should be assistant text
	if messages[0].Role != "assistant" {
		t.Errorf("expected first message role 'assistant', got %q", messages[0].Role)
	}

	// Second should be tool
	if messages[1].Role != "tool" {
		t.Errorf("expected second message role 'tool', got %q", messages[1].Role)
	}
}

func TestParser_ParseLine_EmptyLine(t *testing.T) {
	parser := NewParser()

	messages, err := parser.ParseLine("")
	if err != nil {
		t.Fatalf("failed to parse empty line: %v", err)
	}

	if len(messages) != 0 {
		t.Errorf("expected 0 messages for empty line, got %d", len(messages))
	}
}

func TestParser_ParseLine_InvalidJSON(t *testing.T) {
	parser := NewParser()

	_, err := parser.ParseLine("{invalid json")
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestParser_GetCWD(t *testing.T) {
	parser := NewParser()

	input := `{"type":"user","timestamp":"2026-01-05T15:32:32.836Z","cwd":"/Users/doughughes/code/project","message":{"content":"test"}}`

	cwd, err := parser.GetCWD(input)
	if err != nil {
		t.Fatalf("failed to get CWD: %v", err)
	}

	expected := "/Users/doughughes/code/project"
	if cwd != expected {
		t.Errorf("expected CWD %q, got %q", expected, cwd)
	}
}

func TestParser_GetTimestamp(t *testing.T) {
	parser := NewParser()

	input := `{"type":"user","timestamp":"2026-01-05T15:32:32.836Z","message":{"content":"test"}}`

	timestamp, err := parser.GetTimestamp(input)
	if err != nil {
		t.Fatalf("failed to get timestamp: %v", err)
	}

	if timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}

	if timestamp.Year() != 2026 {
		t.Errorf("expected year 2026, got %d", timestamp.Year())
	}
}
