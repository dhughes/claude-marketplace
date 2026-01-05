package db

import "fmt"

// MockDB is a simple in-memory mock implementation of the DB interface for testing
type MockDB struct {
	conversations map[string]*Conversation
	messages      map[string][]Message // keyed by conversation UUID
	indexStates   map[string]*IndexState
}

// NewMock creates a new mock database
func NewMock() *MockDB {
	return &MockDB{
		conversations: make(map[string]*Conversation),
		messages:      make(map[string][]Message),
		indexStates:   make(map[string]*IndexState),
	}
}

func (m *MockDB) InitSchema() error {
	// No-op for mock
	return nil
}

func (m *MockDB) TruncateAll() error {
	m.conversations = make(map[string]*Conversation)
	m.messages = make(map[string][]Message)
	m.indexStates = make(map[string]*IndexState)
	return nil
}

func (m *MockDB) SaveConversation(conv *Conversation) error {
	m.conversations[conv.UUID] = conv
	return nil
}

func (m *MockDB) SaveMessages(messages []Message) error {
	if len(messages) == 0 {
		return nil
	}

	uuid := messages[0].ConversationUUID
	m.messages[uuid] = append(m.messages[uuid], messages...)

	// Update message count
	if conv, exists := m.conversations[uuid]; exists {
		conv.MessageCount += len(messages)
	}

	return nil
}

func (m *MockDB) GetIndexState(uuid string) (*IndexState, error) {
	state, exists := m.indexStates[uuid]
	if !exists {
		return nil, nil
	}
	return state, nil
}

func (m *MockDB) UpdateIndexState(state *IndexState) error {
	m.indexStates[state.ConversationUUID] = state
	return nil
}

func (m *MockDB) DeleteConversation(uuid string) error {
	delete(m.messages, uuid)
	if conv, exists := m.conversations[uuid]; exists {
		conv.MessageCount = 0
	}
	return nil
}

func (m *MockDB) DeleteIndexState(uuid string) error {
	delete(m.indexStates, uuid)
	return nil
}

func (m *MockDB) GetFirstUserMessage(uuid string) (string, error) {
	messages, exists := m.messages[uuid]
	if !exists {
		return "", nil
	}

	for _, msg := range messages {
		if msg.Role == "user" {
			return msg.Content, nil
		}
	}

	return "", nil
}

func (m *MockDB) Search(query, scope, projectPath string, limit int) ([]Match, error) {
	// Simple mock: just return empty results
	// In real tests, you could populate this with test data
	return []Match{}, nil
}

func (m *MockDB) Close() error {
	return nil
}

// Helper methods for testing

// GetConversation retrieves a conversation for testing
func (m *MockDB) GetConversation(uuid string) (*Conversation, error) {
	conv, exists := m.conversations[uuid]
	if !exists {
		return nil, fmt.Errorf("conversation not found")
	}
	return conv, nil
}

// GetMessages retrieves all messages for a conversation
func (m *MockDB) GetMessages(uuid string) []Message {
	return m.messages[uuid]
}
