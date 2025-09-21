package export

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	"context-extender/internal/database"
)

// CSVExporter implements the Exporter interface for CSV format
type CSVExporter struct {
	supportedColumns []string
}

// NewCSVExporter creates a new CSV exporter
func NewCSVExporter() *CSVExporter {
	return &CSVExporter{
		supportedColumns: AllAvailableColumns,
	}
}

// Export performs the CSV export operation
func (e *CSVExporter) Export(ctx context.Context, backend database.DatabaseBackend, sessions []*database.Session, options *ExportOptions) error {
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

	// Determine columns to export
	columns := options.Columns
	if len(columns) == 0 {
		columns = DefaultCSVColumns
	}

	// Create output file
	file, err := os.Create(options.Output)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	// Create CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header row
	if err := writer.Write(columns); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write data rows
	for _, session := range exportData {
		record := make([]string, len(columns))

		for i, column := range columns {
			record[i] = e.getColumnValue(session, column)
		}

		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write CSV record: %w", err)
		}
	}

	return nil
}

// GetColumnValue extracts the value for a specific column from session data (public method)
func (e *CSVExporter) GetColumnValue(session *SessionExportData, column string) string {
	return e.getColumnValue(session, column)
}

// getColumnValue extracts the value for a specific column from session data
func (e *CSVExporter) getColumnValue(session *SessionExportData, column string) string {
	switch column {
	case "session_id":
		return session.SessionID
	case "project_name":
		return session.ProjectName
	case "working_dir":
		return session.WorkingDir
	case "start_time":
		return session.StartTime.Format("2006-01-02 15:04:05")
	case "end_time":
		return session.EndTime.Format("2006-01-02 15:04:05")
	case "duration":
		return session.Duration
	case "status":
		return session.Status
	case "event_count":
		return strconv.Itoa(session.EventCount)
	case "user_prompts":
		return strconv.Itoa(session.UserPrompts)
	case "claude_replies":
		return strconv.Itoa(session.ClaudeReplies)
	case "total_words":
		return strconv.Itoa(session.TotalWords)
	case "user_words":
		return strconv.Itoa(session.UserWords)
	case "claude_words":
		return strconv.Itoa(session.ClaudeWords)
	case "avg_response_time":
		return session.AvgResponseTime
	case "compression_events":
		return strconv.Itoa(session.CompressionEvents)
	case "tool_usage_count":
		return strconv.Itoa(session.ToolUsageCount)
	case "session_tags":
		return strings.Join(session.SessionTags, ", ")
	case "first_prompt":
		return session.FirstPrompt
	case "last_activity":
		return session.LastActivity
	case "working_dir_name":
		return session.WorkingDirName
	case "raw_metadata":
		// Clean up metadata for CSV (remove newlines, escape quotes)
		metadata := strings.ReplaceAll(session.RawMetadata, "\n", " ")
		metadata = strings.ReplaceAll(metadata, "\"", "\"\"")
		return metadata
	default:
		return ""
	}
}

// GetSupportedColumns returns the list of columns this exporter supports
func (e *CSVExporter) GetSupportedColumns() []string {
	return e.supportedColumns
}

// ValidateOptions validates the export options for CSV export
func (e *CSVExporter) ValidateOptions(options *ExportOptions) error {
	// Validate format
	if options.Format != "csv" {
		return fmt.Errorf("CSV exporter only supports 'csv' format, got: %s", options.Format)
	}

	// Validate output file path
	if options.Output == "" {
		return fmt.Errorf("output file path is required")
	}

	// Validate columns if specified
	if len(options.Columns) > 0 {
		if err := ValidateColumns(options.Columns); err != nil {
			return fmt.Errorf("invalid columns: %w", err)
		}
	}

	// CSV-specific validations
	if options.Pretty {
		return fmt.Errorf("pretty-print option is not applicable for CSV format")
	}

	if options.Compress {
		return fmt.Errorf("compression is not yet supported for CSV export")
	}

	return nil
}

// ExportWithProgress exports CSV data with progress reporting
func (e *CSVExporter) ExportWithProgress(ctx context.Context, backend database.DatabaseBackend, sessions []*database.Session, options *ExportOptions) error {
	if !options.ShowProgress {
		return e.Export(ctx, backend, sessions, options)
	}

	fmt.Printf("🔄 Preparing CSV export for %d sessions...\n", len(sessions))

	// Validate options
	if err := e.ValidateOptions(options); err != nil {
		return fmt.Errorf("invalid export options: %w", err)
	}

	// Prepare session data with progress
	fmt.Print("📊 Processing session data... ")
	exportData, err := PrepareSessionData(ctx, backend, sessions)
	if err != nil {
		return fmt.Errorf("failed to prepare session data: %w", err)
	}
	fmt.Println("✅")

	if len(exportData) == 0 {
		return fmt.Errorf("no data to export")
	}

	// Determine columns
	columns := options.Columns
	if len(columns) == 0 {
		columns = DefaultCSVColumns
	}

	fmt.Printf("📝 Writing CSV with %d columns: %s\n", len(columns), strings.Join(columns, ", "))

	// Create output file
	file, err := os.Create(options.Output)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	// Create CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	if err := writer.Write(columns); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write data with progress
	total := len(exportData)
	for i, session := range exportData {
		record := make([]string, len(columns))

		for j, column := range columns {
			record[j] = e.getColumnValue(session, column)
		}

		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write CSV record: %w", err)
		}

		// Show progress every 10% or for small datasets, every record
		if total <= 10 || (i+1)%(total/10) == 0 || i+1 == total {
			percent := float64(i+1) / float64(total) * 100
			fmt.Printf("📈 Progress: %.0f%% (%d/%d sessions)\n", percent, i+1, total)
		}
	}

	return nil
}

// GetCSVPreview returns a preview of what the CSV export would look like
func (e *CSVExporter) GetCSVPreview(ctx context.Context, backend database.DatabaseBackend, sessions []*database.Session, options *ExportOptions, maxRows int) ([][]string, error) {
	// Prepare limited session data
	limitedSessions := sessions
	if len(sessions) > maxRows {
		limitedSessions = sessions[:maxRows]
	}

	exportData, err := PrepareSessionData(ctx, backend, limitedSessions)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare session data: %w", err)
	}

	// Determine columns
	columns := options.Columns
	if len(columns) == 0 {
		columns = DefaultCSVColumns
	}

	// Build preview data
	var preview [][]string

	// Add header
	preview = append(preview, columns)

	// Add data rows
	for _, session := range exportData {
		record := make([]string, len(columns))
		for i, column := range columns {
			record[i] = e.getColumnValue(session, column)
		}
		preview = append(preview, record)
	}

	return preview, nil
}