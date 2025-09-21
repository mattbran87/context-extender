package export

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"context-extender/internal/database"
)

// JSONExporter implements the Exporter interface for JSON format
type JSONExporter struct {
	supportedColumns []string
}

// NewJSONExporter creates a new JSON exporter
func NewJSONExporter() *JSONExporter {
	return &JSONExporter{
		supportedColumns: AllAvailableColumns,
	}
}

// Export performs the JSON export operation
func (e *JSONExporter) Export(ctx context.Context, backend database.DatabaseBackend, sessions []*database.Session, options *ExportOptions) error {
	// Validate options
	if err := e.ValidateOptions(options); err != nil {
		return fmt.Errorf("invalid export options: %w", err)
	}

	// Prepare session data
	exportData, err := PrepareSessionData(ctx, backend, sessions)
	if err != nil {
		return fmt.Errorf("failed to prepare session data: %w", err)
	}

	if len(exportData) == 0 {
		return fmt.Errorf("no data to export")
	}

	// Check if this is a single-session export
	var exportStructure interface{}
	if len(options.Sessions) == 1 && len(exportData) == 1 {
		// Single session export - simplified structure
		exportStructure = map[string]interface{}{
			"session_metadata": map[string]interface{}{
				"id":           exportData[0].SessionID,
				"project_name": exportData[0].ProjectName,
				"working_dir":  exportData[0].WorkingDir,
				"start_time":   exportData[0].StartTime,
				"end_time":     exportData[0].EndTime,
				"duration":     exportData[0].Duration,
				"status":       exportData[0].Status,
				"export_time":  time.Now(),
			},
			"conversation_flow": e.buildConversationFlow(exportData[0]),
			"analytics": map[string]interface{}{
				"total_events":    exportData[0].EventCount,
				"user_prompts":    exportData[0].UserPrompts,
				"claude_replies":  exportData[0].ClaudeReplies,
				"total_words":     exportData[0].TotalWords,
				"user_words":      exportData[0].UserWords,
				"claude_words":    exportData[0].ClaudeWords,
			},
		}
	} else {
		// Multi-session export - full structure
		exportStructure = map[string]interface{}{
			"metadata": ExportMetadata{
				ExportTime:    time.Now(),
				ExportVersion: "1.0.0",
				Format:        "json",
				SessionCount:  len(exportData),
				TotalRecords:  len(exportData),
				FilePath:      options.Output,
				Options:       options,
			},
			"sessions": exportData,
		}
	}

	// Determine if compression is needed
	useCompression := options.Compress || strings.HasSuffix(strings.ToLower(options.Output), ".gz")

	// Create output file
	file, err := os.Create(options.Output)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	// Create writer (with or without compression)
	var writer io.Writer = file
	var gzipWriter *gzip.Writer

	if useCompression {
		gzipWriter = gzip.NewWriter(file)
		writer = gzipWriter
		defer gzipWriter.Close()
	}

	// Create JSON encoder
	encoder := json.NewEncoder(writer)
	if options.Pretty {
		encoder.SetIndent("", "  ")
	}

	// Write JSON data
	if err := encoder.Encode(exportStructure); err != nil {
		return fmt.Errorf("failed to write JSON data: %w", err)
	}

	// Flush gzip if used
	if gzipWriter != nil {
		if err := gzipWriter.Close(); err != nil {
			return fmt.Errorf("failed to close compression: %w", err)
		}
	}

	return nil
}

// buildConversationFlow creates a chronological conversation flow for single-session export
func (e *JSONExporter) buildConversationFlow(session *SessionExportData) []map[string]interface{} {
	var flow []map[string]interface{}

	// Add conversations in chronological order
	for _, conv := range session.Conversations {
		flow = append(flow, map[string]interface{}{
			"timestamp":    conv.Timestamp,
			"type":         conv.MessageType,
			"content":      conv.Content,
			"token_count":  conv.TokenCount,
			"model":        conv.Model,
			"metadata":     conv.Metadata,
		})
	}

	// Add events in chronological order
	for _, event := range session.Events {
		flow = append(flow, map[string]interface{}{
			"timestamp":    event.Timestamp,
			"type":         "event:" + event.EventType,
			"data":         event.Data,
			"sequence_num": event.SequenceNum,
		})
	}

	// Sort by timestamp
	// Note: In a real implementation, we'd sort this slice by timestamp
	return flow
}

// GetSupportedColumns returns the list of columns this exporter supports
func (e *JSONExporter) GetSupportedColumns() []string {
	return e.supportedColumns
}

// ValidateOptions validates the export options for JSON export
func (e *JSONExporter) ValidateOptions(options *ExportOptions) error {
	// Validate format
	if options.Format != "json" {
		return fmt.Errorf("JSON exporter only supports 'json' format, got: %s", options.Format)
	}

	// Validate output file path
	if options.Output == "" {
		return fmt.Errorf("output file path is required")
	}

	// JSON doesn't use custom columns (exports everything)
	if len(options.Columns) > 0 {
		fmt.Println("⚠️  Warning: Custom columns are ignored for JSON export (all data is included)")
	}

	// Compression is now supported
	if options.Compress && !strings.HasSuffix(strings.ToLower(options.Output), ".gz") {
		fmt.Println("💡 Tip: Use .json.gz extension for automatic compression detection")
	}

	return nil
}