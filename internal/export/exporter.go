package export

import (
	"context"
	"fmt"
	"time"

	"context-extender/internal/database"
)

// Exporter defines the interface for all export implementations
type Exporter interface {
	// Export performs the export operation
	Export(ctx context.Context, backend database.DatabaseBackend, sessions []*database.Session, options *ExportOptions) error

	// GetSupportedColumns returns the list of columns this exporter supports
	GetSupportedColumns() []string

	// ValidateOptions validates the export options for this exporter
	ValidateOptions(options *ExportOptions) error
}

// ExportOptions contains all export configuration options
type ExportOptions struct {
	// Core options
	Format   string `json:"format"`   // csv, json
	Output   string `json:"output"`   // output file path

	// Column customization (CSV only)
	Columns  []string `json:"columns,omitempty"`  // custom column selection

	// Filtering options
	From     string   `json:"from,omitempty"`     // date filter (YYYY-MM-DD)
	To       string   `json:"to,omitempty"`       // date filter (YYYY-MM-DD)
	Project  string   `json:"project,omitempty"`  // project name filter
	Sessions []string `json:"sessions,omitempty"` // specific session IDs
	Status   string   `json:"status,omitempty"`   // session status filter

	// Output formatting options
	Pretty   bool `json:"pretty,omitempty"`   // pretty-print JSON
	Compress bool `json:"compress,omitempty"` // compress output file

	// Progress reporting
	ShowProgress bool `json:"show_progress,omitempty"` // show progress indicators
}

// ExportMetadata contains information about the export operation
type ExportMetadata struct {
	ExportTime    time.Time `json:"export_time"`
	ExportVersion string    `json:"export_version"`
	Format        string    `json:"format"`
	SessionCount  int       `json:"session_count"`
	TotalRecords  int       `json:"total_records"`
	FilePath      string    `json:"file_path"`
	FileSize      int64     `json:"file_size"`
	Options       *ExportOptions `json:"options"`
}

// SessionExportData represents a session's data prepared for export
type SessionExportData struct {
	// Core session information
	SessionID     string    `json:"session_id"`
	ProjectName   string    `json:"project_name"`
	WorkingDir    string    `json:"working_dir"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	Duration      string    `json:"duration"`
	Status        string    `json:"status"`

	// Conversation statistics
	EventCount    int `json:"event_count"`
	UserPrompts   int `json:"user_prompts"`
	ClaudeReplies int `json:"claude_replies"`
	TotalWords    int `json:"total_words"`
	UserWords     int `json:"user_words"`
	ClaudeWords   int `json:"claude_words"`

	// Additional metadata
	Conversations []*database.Conversation `json:"conversations,omitempty"` // for JSON export
	Events        []*database.Event        `json:"events,omitempty"`        // for JSON export
	RawMetadata   string                   `json:"raw_metadata,omitempty"`  // original session metadata
}

// DefaultCSVColumns defines the standard columns for CSV export
var DefaultCSVColumns = []string{
	"session_id",
	"project_name",
	"start_time",
	"end_time",
	"duration",
	"status",
	"user_prompts",
	"claude_replies",
	"total_words",
	"working_dir",
}

// AllAvailableColumns lists all possible columns for export
var AllAvailableColumns = []string{
	"session_id",
	"project_name",
	"working_dir",
	"start_time",
	"end_time",
	"duration",
	"status",
	"event_count",
	"user_prompts",
	"claude_replies",
	"total_words",
	"user_words",
	"claude_words",
	"raw_metadata",
}

// PrepareSessionData converts database sessions to export format
func PrepareSessionData(ctx context.Context, backend database.DatabaseBackend, sessions []*database.Session) ([]*SessionExportData, error) {
	var exportData []*SessionExportData

	for _, session := range sessions {
		// Get conversations for this session
		conversations, err := backend.GetConversationsBySession(ctx, session.ID)
		if err != nil {
			// Continue with what we have, don't fail entire export
			conversations = []*database.Conversation{}
		}

		// Get events for this session
		events, err := backend.GetEventsBySession(ctx, session.ID)
		if err != nil {
			events = []*database.Event{}
		}

		// Calculate statistics
		userPrompts := 0
		claudeReplies := 0
		userWords := 0
		claudeWords := 0

		for _, conv := range conversations {
			switch conv.MessageType {
			case "user":
				userPrompts++
				userWords += countWords(conv.Content)
			case "assistant":
				claudeReplies++
				claudeWords += countWords(conv.Content)
			}
		}

		// Parse metadata for project and working directory
		projectName := ""
		workingDir := ""
		if session.Metadata != "" {
			// TODO: Parse JSON metadata safely
			// For now, extract basic info
		}

		// Calculate duration
		duration := session.UpdatedAt.Sub(session.CreatedAt).String()

		exportData = append(exportData, &SessionExportData{
			SessionID:     session.ID,
			ProjectName:   projectName,
			WorkingDir:    workingDir,
			StartTime:     session.CreatedAt,
			EndTime:       session.UpdatedAt,
			Duration:      duration,
			Status:        session.Status,
			EventCount:    len(events),
			UserPrompts:   userPrompts,
			ClaudeReplies: claudeReplies,
			TotalWords:    userWords + claudeWords,
			UserWords:     userWords,
			ClaudeWords:   claudeWords,
			Conversations: conversations,
			Events:        events,
			RawMetadata:   session.Metadata,
		})
	}

	return exportData, nil
}

// countWords provides a simple word count for content
func countWords(content string) int {
	if content == "" {
		return 0
	}
	// Simple word count by splitting on whitespace
	words := 0
	inWord := false
	for _, char := range content {
		if char == ' ' || char == '\t' || char == '\n' || char == '\r' {
			inWord = false
		} else if !inWord {
			words++
			inWord = true
		}
	}
	return words
}

// ValidateColumns checks if the specified columns are valid
func ValidateColumns(columns []string) error {
	validColumns := make(map[string]bool)
	for _, col := range AllAvailableColumns {
		validColumns[col] = true
	}

	for _, col := range columns {
		if !validColumns[col] {
			return fmt.Errorf("invalid column: %s", col)
		}
	}

	return nil
}