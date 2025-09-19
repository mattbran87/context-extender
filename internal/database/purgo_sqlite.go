package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

// PureGoSQLiteBackend implements DatabaseBackend using modernc.org/sqlite
type PureGoSQLiteBackend struct {
	db     *sql.DB
	config *DatabaseConfig
}

// NewPureGoSQLiteBackend creates a new pure Go SQLite backend
func NewPureGoSQLiteBackend() *PureGoSQLiteBackend {
	return &PureGoSQLiteBackend{}
}

// Initialize initializes the pure Go SQLite backend
func (b *PureGoSQLiteBackend) Initialize(ctx context.Context, config *DatabaseConfig) error {
	b.config = config

	// Open database connection using modernc.org/sqlite driver
	db, err := sql.Open("sqlite", config.DatabasePath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test connection
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Enable WAL mode for better concurrency
	if _, err := db.ExecContext(ctx, "PRAGMA journal_mode=WAL"); err != nil {
		log.Printf("Warning: failed to enable WAL mode: %v", err)
	}

	// Enable foreign keys
	if _, err := db.ExecContext(ctx, "PRAGMA foreign_keys=ON"); err != nil {
		log.Printf("Warning: failed to enable foreign keys: %v", err)
	}

	b.db = db
	log.Printf("Pure Go SQLite backend initialized at %s", config.DatabasePath)

	return nil
}

// Close closes the database connection
func (b *PureGoSQLiteBackend) Close() error {
	if b.db != nil {
		return b.db.Close()
	}
	return nil
}

// GetConnection returns the database connection
func (b *PureGoSQLiteBackend) GetConnection() (*sql.DB, error) {
	if b.db == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	return b.db, nil
}

// CreateSchema creates the database schema
func (b *PureGoSQLiteBackend) CreateSchema(ctx context.Context) error {
	if b.db == nil {
		return fmt.Errorf("database not initialized")
	}

	schema := `
	CREATE TABLE IF NOT EXISTS sessions (
		id TEXT PRIMARY KEY,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		status TEXT,
		metadata TEXT
	);

	CREATE TABLE IF NOT EXISTS events (
		id TEXT PRIMARY KEY,
		session_id TEXT,
		event_type TEXT,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
		sequence_num INTEGER,
		data TEXT,
		FOREIGN KEY (session_id) REFERENCES sessions(id)
	);

	CREATE TABLE IF NOT EXISTS conversations (
		id TEXT PRIMARY KEY,
		session_id TEXT,
		message_type TEXT,
		content TEXT,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
		metadata TEXT,
		token_count INTEGER,
		model TEXT,
		FOREIGN KEY (session_id) REFERENCES sessions(id)
	);

	CREATE INDEX IF NOT EXISTS idx_events_session_id ON events(session_id);
	CREATE INDEX IF NOT EXISTS idx_events_timestamp ON events(timestamp);
	CREATE INDEX IF NOT EXISTS idx_conversations_session_id ON conversations(session_id);
	CREATE INDEX IF NOT EXISTS idx_conversations_timestamp ON conversations(timestamp);
	`

	if _, err := b.db.ExecContext(ctx, schema); err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	return nil
}

// RunMigrations runs database migrations
func (b *PureGoSQLiteBackend) RunMigrations(ctx context.Context) error {
	// For now, just ensure schema exists
	return b.CreateSchema(ctx)
}

// GetBackendInfo returns backend information
func (b *PureGoSQLiteBackend) GetBackendInfo() *BackendInfo {
	capabilities := map[string]bool{
		"encryption":   false, // Application-level encryption can be added
		"full_text":    true,  // SQLite FTS available
		"transactions": true,
		"cgo":          false, // Pure Go!
	}

	return &BackendInfo{
		Name:         "Pure Go SQLite",
		Version:      "modernc.org/sqlite v1.39.0",
		RequiresCGO:  false,
		Features:     []string{"WAL mode", "Foreign keys", "FTS5", "JSON support"},
		Capabilities: capabilities,
	}
}

// CreateSession creates a new session
func (b *PureGoSQLiteBackend) CreateSession(ctx context.Context, session *Session) error {
	query := `
		INSERT INTO sessions (id, created_at, updated_at, status, metadata)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := b.db.ExecContext(ctx, query,
		session.ID,
		session.CreatedAt,
		session.UpdatedAt,
		session.Status,
		session.Metadata,
	)
	return err
}

// GetSession retrieves a session by ID
func (b *PureGoSQLiteBackend) GetSession(ctx context.Context, sessionID string) (*Session, error) {
	query := `
		SELECT id, created_at, updated_at, status, metadata
		FROM sessions WHERE id = ?
	`

	session := &Session{}
	err := b.db.QueryRowContext(ctx, query, sessionID).Scan(
		&session.ID,
		&session.CreatedAt,
		&session.UpdatedAt,
		&session.Status,
		&session.Metadata,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("session not found")
	}

	return session, err
}

// CreateEvent creates a new event
func (b *PureGoSQLiteBackend) CreateEvent(ctx context.Context, event *Event) error {
	query := `
		INSERT INTO events (id, session_id, event_type, timestamp, sequence_num, data)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err := b.db.ExecContext(ctx, query,
		event.ID,
		event.SessionID,
		event.EventType,
		event.Timestamp,
		event.SequenceNum,
		event.Data,
	)
	return err
}

// CreateConversation creates a new conversation entry
func (b *PureGoSQLiteBackend) CreateConversation(ctx context.Context, conv *Conversation) error {
	query := `
		INSERT INTO conversations (id, session_id, message_type, content, timestamp, metadata, token_count, model)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := b.db.ExecContext(ctx, query,
		conv.ID,
		conv.SessionID,
		conv.MessageType,
		conv.Content,
		conv.Timestamp,
		conv.Metadata,
		conv.TokenCount,
		conv.Model,
	)
	return err
}

// ExecuteQuery executes a raw SQL query
func (b *PureGoSQLiteBackend) ExecuteQuery(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return b.db.QueryContext(ctx, query, args...)
}

// ExecuteExec executes a raw SQL exec command
func (b *PureGoSQLiteBackend) ExecuteExec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return b.db.ExecContext(ctx, query, args...)
}

// BeginTransaction starts a new transaction
func (b *PureGoSQLiteBackend) BeginTransaction(ctx context.Context) (*sql.Tx, error) {
	return b.db.BeginTx(ctx, nil)
}

// Missing interface methods

// Ping tests database connectivity
func (b *PureGoSQLiteBackend) Ping(ctx context.Context) error {
	if b.db == nil {
		return fmt.Errorf("database not initialized")
	}
	return b.db.PingContext(ctx)
}

// GetSchemaVersion returns the current schema version
func (b *PureGoSQLiteBackend) GetSchemaVersion(ctx context.Context) (int, error) {
	// For now, return version 1 (initial schema)
	// In the future, this could be stored in a schema_info table
	return 1, nil
}

// MigrateSchema runs migrations to reach target version
func (b *PureGoSQLiteBackend) MigrateSchema(ctx context.Context, targetVersion int) error {
	currentVersion, err := b.GetSchemaVersion(ctx)
	if err != nil {
		return err
	}

	if currentVersion == targetVersion {
		return nil // Already at target version
	}

	// For now, just ensure schema exists
	return b.CreateSchema(ctx)
}

// UpdateSession updates an existing session
func (b *PureGoSQLiteBackend) UpdateSession(ctx context.Context, session *Session) error {
	query := `
		UPDATE sessions
		SET updated_at = ?, status = ?, metadata = ?
		WHERE id = ?
	`
	_, err := b.db.ExecContext(ctx, query,
		session.UpdatedAt,
		session.Status,
		session.Metadata,
		session.ID,
	)
	return err
}

// DeleteSession deletes a session by ID
func (b *PureGoSQLiteBackend) DeleteSession(ctx context.Context, id string) error {
	query := `DELETE FROM sessions WHERE id = ?`
	_, err := b.db.ExecContext(ctx, query, id)
	return err
}

// ListSessions returns sessions based on filters
func (b *PureGoSQLiteBackend) ListSessions(ctx context.Context, filters *SessionFilters) ([]*Session, error) {
	query := `SELECT id, created_at, updated_at, status, metadata FROM sessions WHERE 1=1`
	args := []interface{}{}

	if filters != nil {
		if filters.Status != "" {
			query += " AND status = ?"
			args = append(args, filters.Status)
		}
		if filters.CreatedAfter != nil {
			query += " AND created_at > ?"
			args = append(args, filters.CreatedAfter)
		}
		if filters.CreatedBefore != nil {
			query += " AND created_at < ?"
			args = append(args, filters.CreatedBefore)
		}
		if filters.SortBy != "" {
			query += " ORDER BY " + filters.SortBy
			if filters.SortOrder == "DESC" {
				query += " DESC"
			}
		}
		if filters.Limit > 0 {
			query += " LIMIT ?"
			args = append(args, filters.Limit)
		}
		if filters.Offset > 0 {
			query += " OFFSET ?"
			args = append(args, filters.Offset)
		}
	}

	rows, err := b.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*Session
	for rows.Next() {
		session := &Session{}
		err := rows.Scan(
			&session.ID,
			&session.CreatedAt,
			&session.UpdatedAt,
			&session.Status,
			&session.Metadata,
		)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}

	return sessions, rows.Err()
}

// GetEventsBySession returns all events for a session
func (b *PureGoSQLiteBackend) GetEventsBySession(ctx context.Context, sessionID string) ([]*Event, error) {
	query := `
		SELECT id, session_id, event_type, timestamp, sequence_num, data
		FROM events WHERE session_id = ? ORDER BY sequence_num
	`

	rows, err := b.db.QueryContext(ctx, query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*Event
	for rows.Next() {
		event := &Event{}
		err := rows.Scan(
			&event.ID,
			&event.SessionID,
			&event.EventType,
			&event.Timestamp,
			&event.SequenceNum,
			&event.Data,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, rows.Err()
}

// CreateEventBatch creates multiple events in a single transaction
func (b *PureGoSQLiteBackend) CreateEventBatch(ctx context.Context, events []*Event) error {
	tx, err := b.BeginTransaction(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO events (id, session_id, event_type, timestamp, sequence_num, data)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, event := range events {
		_, err := stmt.ExecContext(ctx,
			event.ID,
			event.SessionID,
			event.EventType,
			event.Timestamp,
			event.SequenceNum,
			event.Data,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// GetConversationsBySession returns all conversations for a session
func (b *PureGoSQLiteBackend) GetConversationsBySession(ctx context.Context, sessionID string) ([]*Conversation, error) {
	query := `
		SELECT id, session_id, message_type, content, timestamp, metadata, token_count, model
		FROM conversations WHERE session_id = ? ORDER BY timestamp
	`

	rows, err := b.db.QueryContext(ctx, query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversations []*Conversation
	for rows.Next() {
		conv := &Conversation{}
		err := rows.Scan(
			&conv.ID,
			&conv.SessionID,
			&conv.MessageType,
			&conv.Content,
			&conv.Timestamp,
			&conv.Metadata,
			&conv.TokenCount,
			&conv.Model,
		)
		if err != nil {
			return nil, err
		}
		conversations = append(conversations, conv)
	}

	return conversations, rows.Err()
}

// SearchConversations searches conversations by content
func (b *PureGoSQLiteBackend) SearchConversations(ctx context.Context, query string, limit int) ([]*Conversation, error) {
	// Simple text search - can be enhanced with FTS later
	searchQuery := `
		SELECT id, session_id, message_type, content, timestamp, metadata, token_count, model
		FROM conversations
		WHERE content LIKE ?
		ORDER BY timestamp DESC
		LIMIT ?
	`

	rows, err := b.db.QueryContext(ctx, searchQuery, "%"+query+"%", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversations []*Conversation
	for rows.Next() {
		conv := &Conversation{}
		err := rows.Scan(
			&conv.ID,
			&conv.SessionID,
			&conv.MessageType,
			&conv.Content,
			&conv.Timestamp,
			&conv.Metadata,
			&conv.TokenCount,
			&conv.Model,
		)
		if err != nil {
			return nil, err
		}
		conversations = append(conversations, conv)
	}

	return conversations, rows.Err()
}

// GetDatabaseStats returns database statistics
func (b *PureGoSQLiteBackend) GetDatabaseStats(ctx context.Context) (*DatabaseStats, error) {
	stats := &DatabaseStats{}

	// Count sessions
	err := b.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM sessions").Scan(&stats.SessionCount)
	if err != nil {
		return nil, err
	}

	// Count events
	err = b.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM events").Scan(&stats.EventCount)
	if err != nil {
		return nil, err
	}

	// Count conversations
	err = b.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM conversations").Scan(&stats.ConversationCount)
	if err != nil {
		return nil, err
	}

	// Get oldest record
	var oldestTime sql.NullTime
	err = b.db.QueryRowContext(ctx, `
		SELECT MIN(created_at) FROM (
			SELECT created_at FROM sessions
			UNION ALL
			SELECT timestamp FROM events
			UNION ALL
			SELECT timestamp FROM conversations
		)
	`).Scan(&oldestTime)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if oldestTime.Valid {
		stats.OldestRecord = &oldestTime.Time
	}

	// Get newest record
	var newestTime sql.NullTime
	err = b.db.QueryRowContext(ctx, `
		SELECT MAX(created_at) FROM (
			SELECT created_at FROM sessions
			UNION ALL
			SELECT timestamp FROM events
			UNION ALL
			SELECT timestamp FROM conversations
		)
	`).Scan(&newestTime)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if newestTime.Valid {
		stats.NewestRecord = &newestTime.Time
	}

	return stats, nil
}