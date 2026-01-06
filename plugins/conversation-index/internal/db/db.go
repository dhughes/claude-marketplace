package db

import (
	"database/sql"
	"fmt"

	"github.com/doughughes/claude-marketplace/plugins/conversation-index/internal/shared"
	_ "modernc.org/sqlite"
)

// DB interface defines all database operations
type DB interface {
	InitSchema() error
	TruncateAll() error
	SaveConversation(conv *Conversation) error
	SaveMessages(messages []Message) error
	GetIndexState(uuid string) (*IndexState, error)
	UpdateIndexState(state *IndexState) error
	DeleteConversation(uuid string) error
	DeleteIndexState(uuid string) error
	GetFirstUserMessage(uuid string) (string, error)
	Search(query, scope, projectPath string, limit int) ([]Match, error)
	Close() error
}

// sqliteDB implements the DB interface using SQLite
type sqliteDB struct {
	conn *sql.DB
}

// Open opens a SQLite database at the given path
func Open(path string) (DB, error) {
	// Add mode=rwc to ensure read-write-create access
	// Add busy_timeout to handle concurrent access (wait up to 5 seconds)
	connStr := path + "?mode=rwc&_busy_timeout=5000"
	conn, err := sql.Open("sqlite", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Note: WAL mode disabled due to sandbox restrictions
	// WAL creates additional files (-wal, -shm) which may be blocked by sandbox

	return &sqliteDB{conn: conn}, nil
}

// TruncateAll deletes all data from all tables
func (db *sqliteDB) TruncateAll() error {
	// Delete in order to respect foreign key constraints
	queries := []string{
		"DELETE FROM messages",
		"DELETE FROM conversations",
		"DELETE FROM index_state",
	}

	for _, query := range queries {
		if _, err := db.conn.Exec(query); err != nil {
			return fmt.Errorf("failed to truncate tables: %w", err)
		}
	}

	return nil
}

// InitSchema creates all tables, indexes, and triggers
func (db *sqliteDB) InitSchema() error {
	schema := `
		CREATE TABLE IF NOT EXISTS conversations (
			uuid TEXT PRIMARY KEY,
			project_path TEXT NOT NULL,
			encoded_path TEXT NOT NULL,
			created_at TEXT NOT NULL,
			last_updated TEXT NOT NULL,
			message_count INTEGER DEFAULT 0
		);

		CREATE TABLE IF NOT EXISTS messages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			conversation_uuid TEXT NOT NULL,
			timestamp TEXT NOT NULL,
			role TEXT NOT NULL,
			content TEXT,
			FOREIGN KEY (conversation_uuid) REFERENCES conversations(uuid)
		);

		CREATE VIRTUAL TABLE IF NOT EXISTS messages_fts USING fts5(
			conversation_uuid,
			content,
			content=messages,
			content_rowid=id
		);

		CREATE TABLE IF NOT EXISTS index_state (
			conversation_uuid TEXT PRIMARY KEY,
			last_indexed_line INTEGER DEFAULT 0,
			last_modified_time TEXT
		);

		CREATE TRIGGER IF NOT EXISTS messages_ai AFTER INSERT ON messages BEGIN
			INSERT INTO messages_fts(rowid, conversation_uuid, content)
			VALUES (new.id, new.conversation_uuid, new.content);
		END;

		CREATE TRIGGER IF NOT EXISTS messages_ad AFTER DELETE ON messages BEGIN
			DELETE FROM messages_fts WHERE rowid = old.id;
		END;

		CREATE TRIGGER IF NOT EXISTS messages_au AFTER UPDATE ON messages BEGIN
			DELETE FROM messages_fts WHERE rowid = old.id;
			INSERT INTO messages_fts(rowid, conversation_uuid, content)
			VALUES (new.id, new.conversation_uuid, new.content);
		END;
	`

	_, err := db.conn.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to initialize schema: %w", err)
	}

	return nil
}

// SaveConversation inserts or updates a conversation record
func (db *sqliteDB) SaveConversation(conv *Conversation) error {
	query := `
		INSERT INTO conversations (uuid, project_path, encoded_path, created_at, last_updated, message_count)
		VALUES (?, ?, ?, ?, ?, COALESCE((SELECT message_count FROM conversations WHERE uuid = ?), 0))
		ON CONFLICT(uuid) DO UPDATE SET
			project_path = excluded.project_path,
			encoded_path = excluded.encoded_path,
			last_updated = excluded.last_updated
	`

	_, err := db.conn.Exec(query,
		conv.UUID,
		conv.ProjectPath,
		conv.EncodedPath,
		shared.FormatTimestamp(conv.CreatedAt),
		shared.FormatTimestamp(conv.LastUpdated),
		conv.UUID,
	)

	if err != nil {
		return fmt.Errorf("failed to save conversation: %w", err)
	}

	return nil
}

// SaveMessages inserts multiple messages in a transaction
func (db *sqliteDB) SaveMessages(messages []Message) error {
	if len(messages) == 0 {
		return nil
	}

	tx, err := db.conn.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO messages (conversation_uuid, timestamp, role, content)
		VALUES (?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, msg := range messages {
		_, err := stmt.Exec(
			msg.ConversationUUID,
			shared.FormatTimestamp(msg.Timestamp),
			msg.Role,
			msg.Content,
		)
		if err != nil {
			return fmt.Errorf("failed to insert message: %w", err)
		}
	}

	// Update message count
	updateQuery := `
		UPDATE conversations
		SET message_count = message_count + ?
		WHERE uuid = ?
	`
	_, err = tx.Exec(updateQuery, len(messages), messages[0].ConversationUUID)
	if err != nil {
		return fmt.Errorf("failed to update message count: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetIndexState retrieves the index state for a conversation
func (db *sqliteDB) GetIndexState(uuid string) (*IndexState, error) {
	query := `
		SELECT conversation_uuid, last_indexed_line, last_modified_time
		FROM index_state
		WHERE conversation_uuid = ?
	`

	var state IndexState
	var modifiedStr string

	err := db.conn.QueryRow(query, uuid).Scan(
		&state.ConversationUUID,
		&state.LastIndexedLine,
		&modifiedStr,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get index state: %w", err)
	}

	// Parse timestamp
	modifiedTime, err := shared.ParseTimestamp(modifiedStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse timestamp: %w", err)
	}
	state.LastModifiedTime = modifiedTime

	return &state, nil
}

// UpdateIndexState inserts or updates the index state for a conversation
func (db *sqliteDB) UpdateIndexState(state *IndexState) error {
	query := `
		INSERT OR REPLACE INTO index_state (conversation_uuid, last_indexed_line, last_modified_time)
		VALUES (?, ?, ?)
	`

	_, err := db.conn.Exec(query,
		state.ConversationUUID,
		state.LastIndexedLine,
		shared.FormatTimestamp(state.LastModifiedTime),
	)

	if err != nil {
		return fmt.Errorf("failed to update index state: %w", err)
	}

	return nil
}

// DeleteConversation deletes all messages for a conversation
func (db *sqliteDB) DeleteConversation(uuid string) error {
	query := `DELETE FROM messages WHERE conversation_uuid = ?`

	_, err := db.conn.Exec(query, uuid)
	if err != nil {
		return fmt.Errorf("failed to delete conversation messages: %w", err)
	}

	// Reset message count
	updateQuery := `UPDATE conversations SET message_count = 0 WHERE uuid = ?`
	_, err = db.conn.Exec(updateQuery, uuid)
	if err != nil {
		return fmt.Errorf("failed to reset message count: %w", err)
	}

	return nil
}

// DeleteIndexState deletes the index state for a conversation
func (db *sqliteDB) DeleteIndexState(uuid string) error {
	query := `DELETE FROM index_state WHERE conversation_uuid = ?`

	_, err := db.conn.Exec(query, uuid)
	if err != nil {
		return fmt.Errorf("failed to delete index state: %w", err)
	}

	return nil
}

// GetFirstUserMessage retrieves the first user message from a conversation
func (db *sqliteDB) GetFirstUserMessage(uuid string) (string, error) {
	query := `
		SELECT content
		FROM messages
		WHERE conversation_uuid = ? AND role = 'user'
		ORDER BY timestamp ASC
		LIMIT 1
	`

	var content string
	err := db.conn.QueryRow(query, uuid).Scan(&content)

	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("failed to get first user message: %w", err)
	}

	return content, nil
}

// Search performs an FTS5 search across conversations
func (db *sqliteDB) Search(query, scope, projectPath string, limit int) ([]Match, error) {
	sqlQuery := `
		SELECT
			c.uuid,
			c.project_path,
			c.encoded_path,
			c.created_at,
			c.last_updated,
			c.message_count,
			messages_fts.rank as relevance_score
		FROM messages_fts
		JOIN messages m ON messages_fts.rowid = m.id
		JOIN conversations c ON m.conversation_uuid = c.uuid
		WHERE messages_fts MATCH ?
	`

	args := []interface{}{query}

	// Add project scope filtering
	if scope == "current_project" && projectPath != "" {
		sqlQuery += ` AND c.encoded_path = ?`
		args = append(args, projectPath)
	}

	sqlQuery += `
		GROUP BY c.uuid
		ORDER BY relevance_score DESC
		LIMIT ?
	`
	args = append(args, limit)

	rows, err := db.conn.Query(sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute search: %w", err)
	}
	defer rows.Close()

	var matches []Match
	for rows.Next() {
		var match Match
		err := rows.Scan(
			&match.UUID,
			&match.ProjectPath,
			&match.EncodedPath,
			&match.CreatedAt,
			&match.LastUpdated,
			&match.MessageCount,
			&match.RelevanceScore,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Get summary (first user message)
		summary, err := db.GetFirstUserMessage(match.UUID)
		if err != nil {
			return nil, fmt.Errorf("failed to get summary: %w", err)
		}

		// Truncate summary to 150 characters (UTF-8 safe)
		if summary == "" {
			summary = "No summary available"
		} else {
			summary = shared.TruncateString(summary, 150)
		}
		match.Summary = summary

		// Convert relevance score to absolute value
		if match.RelevanceScore < 0 {
			match.RelevanceScore = -match.RelevanceScore
		}

		matches = append(matches, match)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return matches, nil
}

// Close closes the database connection
func (db *sqliteDB) Close() error {
	return db.conn.Close()
}
