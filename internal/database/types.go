package database

import (
	"context"
	"time"
)

// DatabaseBackend defines the interface all database implementations must satisfy
type DatabaseBackend interface {
	// Lifecycle Management
	Initialize(ctx context.Context, config *DatabaseConfig) error
	Close() error
	Ping(ctx context.Context) error

	// Schema Management
	CreateSchema(ctx context.Context) error
	GetSchemaVersion(ctx context.Context) (int, error)
	MigrateSchema(ctx context.Context, targetVersion int) error

	// Session Operations
	CreateSession(ctx context.Context, session *Session) error
	GetSession(ctx context.Context, id string) (*Session, error)
	UpdateSession(ctx context.Context, session *Session) error
	DeleteSession(ctx context.Context, id string) error
	ListSessions(ctx context.Context, filters *SessionFilters) ([]*Session, error)

	// Event Operations
	CreateEvent(ctx context.Context, event *Event) error
	GetEventsBySession(ctx context.Context, sessionID string) ([]*Event, error)
	CreateEventBatch(ctx context.Context, events []*Event) error

	// Conversation Operations
	CreateConversation(ctx context.Context, conv *Conversation) error
	GetConversationsBySession(ctx context.Context, sessionID string) ([]*Conversation, error)
	SearchConversations(ctx context.Context, query string, limit int) ([]*Conversation, error)

	// Statistics
	GetDatabaseStats(ctx context.Context) (*DatabaseStats, error)

	// Backend Information
	GetBackendInfo() *BackendInfo
}

// BackendFactory creates database backend instances
type BackendFactory interface {
	CreateBackend(config *DatabaseConfig) (DatabaseBackend, error)
	IsAvailable() bool
	GetCapabilities() *BackendCapabilities
}

// Session represents a conversation session
type Session struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Status    string    `json:"status"`
	Metadata  string    `json:"metadata,omitempty"`
}

// Event represents a session event
type Event struct {
	ID          string    `json:"id"`
	SessionID   string    `json:"session_id"`
	EventType   string    `json:"event_type"`
	Timestamp   time.Time `json:"timestamp"`
	SequenceNum int       `json:"sequence_num"`
	Data        string    `json:"data,omitempty"`
}

// Conversation represents a conversation message
type Conversation struct {
	ID          string    `json:"id"`
	SessionID   string    `json:"session_id"`
	MessageType string    `json:"message_type"`
	Content     string    `json:"content"`
	Timestamp   time.Time `json:"timestamp"`
	Metadata    string    `json:"metadata,omitempty"`
	TokenCount  int       `json:"token_count,omitempty"`
	Model       string    `json:"model,omitempty"`
}

// SessionFilters defines filters for session queries
type SessionFilters struct {
	Status        string     `json:"status,omitempty"`
	CreatedAfter  *time.Time `json:"created_after,omitempty"`
	CreatedBefore *time.Time `json:"created_before,omitempty"`
	Limit         int        `json:"limit,omitempty"`
	Offset        int        `json:"offset,omitempty"`
	SortBy        string     `json:"sort_by,omitempty"`
	SortOrder     string     `json:"sort_order,omitempty"`
}

// DatabaseStats provides database statistics
type DatabaseStats struct {
	SessionCount      int        `json:"session_count"`
	EventCount        int        `json:"event_count"`
	ConversationCount int        `json:"conversation_count"`
	ImportCount       int        `json:"import_count"`
	DatabaseSize      int64      `json:"database_size"`
	OldestRecord      *time.Time `json:"oldest_record,omitempty"`
	NewestRecord      *time.Time `json:"newest_record,omitempty"`
}

// BackendInfo provides information about a database backend
type BackendInfo struct {
	Name         string            `json:"name"`
	Version      string            `json:"version"`
	RequiresCGO  bool              `json:"requires_cgo"`
	Features     []string          `json:"features"`
	Capabilities map[string]bool   `json:"capabilities"`
}

// BackendCapabilities describes what a backend supports
type BackendCapabilities struct {
	SupportsEncryption   bool     `json:"supports_encryption"`
	SupportsFullText     bool     `json:"supports_full_text"`
	SupportsTransactions bool     `json:"supports_transactions"`
	RequiresCGO          bool     `json:"requires_cgo"`
	PlatformSupport      []string `json:"platform_support"`
}