package export

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
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
	Format   string `json:"format"`   // csv, json, xlsx
	Output   string `json:"output"`   // output file path

	// Column customization (CSV only)
	Columns  []string `json:"columns,omitempty"`  // custom column selection

	// Filtering options
	From        string   `json:"from,omitempty"`         // date filter (YYYY-MM-DD)
	To          string   `json:"to,omitempty"`           // date filter (YYYY-MM-DD)
	Project     string   `json:"project,omitempty"`      // project name filter
	Sessions    []string `json:"sessions,omitempty"`     // specific session IDs
	Status      string   `json:"status,omitempty"`       // session status filter
	MinDuration string   `json:"min_duration,omitempty"` // minimum session duration
	MaxDuration string   `json:"max_duration,omitempty"` // maximum session duration

	// Output formatting options
	Pretty   bool `json:"pretty,omitempty"`   // pretty-print JSON
	Compress bool `json:"compress,omitempty"` // compress output file

	// Display options
	ShowProgress bool `json:"show_progress,omitempty"` // show progress indicators
	Preview      bool `json:"preview,omitempty"`       // preview mode
	ShowStats    bool `json:"show_stats,omitempty"`    // show export statistics
	Validate     bool `json:"validate,omitempty"`      // validate export
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

	// Extended analytics
	AvgResponseTime   string   `json:"avg_response_time"`
	CompressionEvents int      `json:"compression_events"`
	ToolUsageCount    int      `json:"tool_usage_count"`
	SessionTags       []string `json:"session_tags"`
	FirstPrompt       string   `json:"first_prompt"`
	LastActivity      string   `json:"last_activity"`
	WorkingDirName    string   `json:"working_dir_name"`

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
	"avg_response_time",
	"compression_events",
	"tool_usage_count",
	"session_tags",
	"first_prompt",
	"last_activity",
	"working_dir_name",
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
			// Simple metadata extraction
			if strings.Contains(session.Metadata, "working_directory") {
				// Extract working directory from JSON-like metadata
				if start := strings.Index(session.Metadata, "working_directory\":\""); start != -1 {
					start += len("working_directory\":\"")
					if end := strings.Index(session.Metadata[start:], "\""); end != -1 {
						workingDir = session.Metadata[start : start+end]
					}
				}
			}
		}

		// Calculate analytics
		avgResponseTime := calculateAverageResponseTime(conversations)
		compressionEvents := countCompressionEvents(events)
		toolUsageCount := countToolUsage(conversations)
		sessionTags := generateSessionTags(session, conversations, events)
		firstPrompt := getFirstPrompt(conversations)
		lastActivity := getLastActivity(conversations, events)
		workingDirName := ""
		if workingDir != "" {
			workingDirName = filepath.Base(workingDir)
		}

		// Calculate duration
		duration := session.UpdatedAt.Sub(session.CreatedAt).String()

		exportData = append(exportData, &SessionExportData{
			SessionID:         session.ID,
			ProjectName:       projectName,
			WorkingDir:        workingDir,
			StartTime:         session.CreatedAt,
			EndTime:           session.UpdatedAt,
			Duration:          duration,
			Status:            session.Status,
			EventCount:        len(events),
			UserPrompts:       userPrompts,
			ClaudeReplies:     claudeReplies,
			TotalWords:        userWords + claudeWords,
			UserWords:         userWords,
			ClaudeWords:       claudeWords,
			AvgResponseTime:   avgResponseTime,
			CompressionEvents: compressionEvents,
			ToolUsageCount:    toolUsageCount,
			SessionTags:       sessionTags,
			FirstPrompt:       firstPrompt,
			LastActivity:      lastActivity,
			WorkingDirName:    workingDirName,
			Conversations:     conversations,
			Events:            events,
			RawMetadata:       session.Metadata,
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

// Helper functions for advanced analytics

// calculateAverageResponseTime calculates average time between user prompts and Claude responses
func calculateAverageResponseTime(conversations []*database.Conversation) string {
	if len(conversations) < 2 {
		return "0s"
	}

	var totalGap time.Duration
	var gapCount int

	for i := 0; i < len(conversations)-1; i++ {
		if conversations[i].MessageType == "user" && conversations[i+1].MessageType == "assistant" {
			gap := conversations[i+1].Timestamp.Sub(conversations[i].Timestamp)
			totalGap += gap
			gapCount++
		}
	}

	if gapCount == 0 {
		return "0s"
	}

	avgGap := totalGap / time.Duration(gapCount)
	return avgGap.String()
}

// countCompressionEvents counts compression-related events
func countCompressionEvents(events []*database.Event) int {
	count := 0
	for _, event := range events {
		if event.EventType == "compression" || strings.Contains(event.EventType, "compress") {
			count++
		}
	}
	return count
}

// countToolUsage counts tool usage in conversations
func countToolUsage(conversations []*database.Conversation) int {
	count := 0
	for _, conv := range conversations {
		// Look for tool usage indicators in content
		content := strings.ToLower(conv.Content)
		if strings.Contains(content, "tool_calls") ||
		   strings.Contains(content, "function_call") ||
		   strings.Contains(content, "<function_calls>") {
			count++
		}
	}
	return count
}

// generateSessionTags creates tags based on session characteristics
func generateSessionTags(session *database.Session, conversations []*database.Conversation, events []*database.Event) []string {
	var tags []string

	// Duration-based tags
	duration := session.UpdatedAt.Sub(session.CreatedAt)
	if duration < 2*time.Minute {
		tags = append(tags, "quick-session")
	} else if duration > 30*time.Minute {
		tags = append(tags, "long-session")
	}

	// Activity-based tags
	if len(conversations) >= 10 {
		tags = append(tags, "active-conversation")
	} else if len(conversations) <= 3 {
		tags = append(tags, "brief-interaction")
	}

	// Status-based tags
	switch session.Status {
	case "completed":
		tags = append(tags, "completed-naturally")
	case "error":
		tags = append(tags, "error-terminated")
	case "active":
		tags = append(tags, "still-active")
	}

	// Time-based tags
	hour := session.CreatedAt.Hour()
	if hour >= 6 && hour < 12 {
		tags = append(tags, "morning")
	} else if hour >= 12 && hour < 18 {
		tags = append(tags, "afternoon")
	} else if hour >= 18 && hour < 22 {
		tags = append(tags, "evening")
	} else {
		tags = append(tags, "night")
	}

	// Content-based tags
	hasCodeContent := false
	for _, conv := range conversations {
		if strings.Contains(conv.Content, "```") || strings.Contains(conv.Content, "function") {
			hasCodeContent = true
			break
		}
	}
	if hasCodeContent {
		tags = append(tags, "coding-session")
	}

	return tags
}

// getFirstPrompt extracts the first user prompt
func getFirstPrompt(conversations []*database.Conversation) string {
	for _, conv := range conversations {
		if conv.MessageType == "user" {
			preview := conv.Content
			if len(preview) > 100 {
				preview = preview[:97] + "..."
			}
			return preview
		}
	}
	return ""
}

// getLastActivity finds the timestamp of the last activity
func getLastActivity(conversations []*database.Conversation, events []*database.Event) string {
	var lastTime time.Time

	// Check conversations
	for _, conv := range conversations {
		if conv.Timestamp.After(lastTime) {
			lastTime = conv.Timestamp
		}
	}

	// Check events
	for _, event := range events {
		if event.Timestamp.After(lastTime) {
			lastTime = event.Timestamp
		}
	}

	if lastTime.IsZero() {
		return ""
	}

	return lastTime.Format("2006-01-02 15:04:05")
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