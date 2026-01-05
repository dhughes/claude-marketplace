package indexer

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/doughughes/claude-marketplace/plugins/conversation-index/internal/db"
	"github.com/doughughes/claude-marketplace/plugins/conversation-index/internal/shared"
)

// JSONLEntry represents a single line in a conversation JSONL file
type JSONLEntry struct {
	Type      string          `json:"type"`
	Timestamp string          `json:"timestamp"`
	CWD       string          `json:"cwd"`
	Message   *JSONLMessage   `json:"message"`
}

// JSONLMessage represents the message field in a JSONL entry
type JSONLMessage struct {
	Content interface{} `json:"content"`
}

// Parser handles parsing of JSONL conversation files
type Parser struct{}

// NewParser creates a new JSONL parser
func NewParser() *Parser {
	return &Parser{}
}

// ParseLine parses a single JSONL line and extracts searchable messages
func (p *Parser) ParseLine(line string) ([]db.Message, error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil, nil
	}

	var entry JSONLEntry
	if err := json.Unmarshal([]byte(line), &entry); err != nil {
		return nil, fmt.Errorf("failed to parse JSONL: %w", err)
	}

	return p.extractMessages(&entry)
}

// extractMessages extracts searchable content from a JSONL entry
func (p *Parser) extractMessages(entry *JSONLEntry) ([]db.Message, error) {
	var messages []db.Message

	timestamp, err := shared.ParseTimestamp(entry.Timestamp)
	if err != nil {
		timestamp = time.Now() // Fallback to current time
	}

	// Handle user messages
	if entry.Type == "user" && entry.Message != nil {
		if content, ok := entry.Message.Content.(string); ok && content != "" {
			messages = append(messages, db.Message{
				Timestamp: timestamp,
				Role:      "user",
				Content:   content,
			})
		}
	}

	// Handle assistant messages
	if entry.Type == "assistant" && entry.Message != nil {
		msgs := p.extractAssistantMessages(entry.Message.Content, timestamp)
		messages = append(messages, msgs...)
	}

	return messages, nil
}

// extractAssistantMessages extracts content from assistant message content array
func (p *Parser) extractAssistantMessages(content interface{}, timestamp time.Time) []db.Message {
	var messages []db.Message

	// Content can be a string or an array
	contentArray, ok := content.([]interface{})
	if !ok {
		return messages
	}

	for _, item := range contentArray {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		itemType, _ := itemMap["type"].(string)

		// Extract text responses
		if itemType == "text" {
			if text, ok := itemMap["text"].(string); ok && text != "" {
				messages = append(messages, db.Message{
					Timestamp: timestamp,
					Role:      "assistant",
					Content:   text,
				})
			}
		}

		// Extract tool use information
		if itemType == "tool_use" {
			toolMsg := p.extractToolUse(itemMap, timestamp)
			if toolMsg.Content != "" {
				messages = append(messages, toolMsg)
			}
		}
	}

	return messages
}

// extractToolUse extracts searchable content from tool use
func (p *Parser) extractToolUse(itemMap map[string]interface{}, timestamp time.Time) db.Message {
	var parts []string

	// Tool name
	if name, ok := itemMap["name"].(string); ok {
		parts = append(parts, "Tool: "+name)
	}

	// Tool input parameters
	if input, ok := itemMap["input"].(map[string]interface{}); ok {
		if filePath, ok := input["file_path"].(string); ok {
			parts = append(parts, "File: "+filePath)
		}
		if pattern, ok := input["pattern"].(string); ok {
			parts = append(parts, "Pattern: "+pattern)
		}
		if description, ok := input["description"].(string); ok {
			parts = append(parts, description)
		}
		if prompt, ok := input["prompt"].(string); ok {
			parts = append(parts, prompt)
		}
		if command, ok := input["command"].(string); ok {
			parts = append(parts, "Command: "+command)
		}
	}

	return db.Message{
		Timestamp: timestamp,
		Role:      "tool",
		Content:   strings.Join(parts, " "),
	}
}

// GetCWD extracts the working directory from the first JSONL line
func (p *Parser) GetCWD(firstLine string) (string, error) {
	var entry JSONLEntry
	if err := json.Unmarshal([]byte(firstLine), &entry); err != nil {
		return "", fmt.Errorf("failed to parse first line: %w", err)
	}

	return entry.CWD, nil
}

// GetTimestamp extracts the timestamp from the first JSONL line
func (p *Parser) GetTimestamp(firstLine string) (time.Time, error) {
	var entry JSONLEntry
	if err := json.Unmarshal([]byte(firstLine), &entry); err != nil {
		return time.Time{}, fmt.Errorf("failed to parse first line: %w", err)
	}

	return shared.ParseTimestamp(entry.Timestamp)
}
