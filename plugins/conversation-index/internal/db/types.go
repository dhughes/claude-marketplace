package db

import "time"

// Conversation represents a conversation metadata record
type Conversation struct {
	UUID         string
	ProjectPath  string
	EncodedPath  string
	CreatedAt    time.Time
	LastUpdated  time.Time
	MessageCount int
}

// Message represents a single message in a conversation
type Message struct {
	ID              int64
	ConversationUUID string
	Timestamp       time.Time
	Role            string // user, assistant, tool
	Content         string
}

// IndexState tracks the indexing progress for a conversation
type IndexState struct {
	ConversationUUID  string
	LastIndexedLine   int
	LastModifiedTime  time.Time
}

// Match represents a search result
type Match struct {
	UUID           string  `json:"uuid"`
	ProjectPath    string  `json:"project_path"`
	EncodedPath    string  `json:"encoded_path"`
	CreatedAt      string  `json:"created_at"`
	LastUpdated    string  `json:"last_updated"`
	MessageCount   int     `json:"message_count"`
	Summary        string  `json:"summary"`
	RelevanceScore float64 `json:"relevance_score"`
}

// SearchResult represents the full search response
type SearchResult struct {
	Query          string  `json:"query"`
	Scope          string  `json:"scope"`
	CurrentProject string  `json:"current_project,omitempty"`
	TotalMatches   int     `json:"total_matches"`
	Matches        []Match `json:"matches"`
}
