package cmd

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"context-extender/internal/database"
	"context-extender/internal/export"
)

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export conversation data to various formats",
	Long: `Export conversation data to CSV, JSON, and other formats for analysis.

Supported formats:
  - csv: Comma-separated values for Excel/spreadsheet analysis
  - json: Structured JSON for programmatic access

Examples:
  # Export all conversations to CSV
  context-extender export --format csv --output conversations.csv

  # Export with date range filter
  context-extender export --format csv --from 2024-01-01 --to 2024-01-31 --output monthly.csv

  # Export specific sessions to JSON
  context-extender export --format json --sessions session1,session2 --output specific.json

  # Export with custom CSV columns
  context-extender export --format csv --columns session_id,start_time,duration --output summary.csv`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return handleExport(cmd, args)
	},
}

// Export command flags
var (
	exportFormat      string
	exportOutput      string
	exportColumns     []string
	exportFrom        string
	exportTo          string
	exportProject     string
	exportSessions    []string
	exportStatus      string
	exportMinDuration string
	exportMaxDuration string
	exportPretty      bool
	exportCompress    bool
	exportProgress    bool
	exportPreview     bool
	exportStats       bool
	exportValidate    bool
)

func init() {
	rootCmd.AddCommand(exportCmd)

	// Core export flags
	exportCmd.Flags().StringVar(&exportFormat, "format", "csv", "Export format (csv, json)")
	exportCmd.Flags().StringVar(&exportOutput, "output", "", "Output file path (required)")

	// Filtering flags
	exportCmd.Flags().StringVar(&exportFrom, "from", "", "Filter from date (YYYY-MM-DD)")
	exportCmd.Flags().StringVar(&exportTo, "to", "", "Filter to date (YYYY-MM-DD)")
	exportCmd.Flags().StringVar(&exportProject, "project", "", "Filter by project name")
	exportCmd.Flags().StringSliceVar(&exportSessions, "sessions", []string{}, "Filter by specific session IDs (comma-separated)")
	exportCmd.Flags().StringVar(&exportStatus, "status", "", "Filter by session status (active, completed, error)")
	exportCmd.Flags().StringVar(&exportMinDuration, "min-duration", "", "Filter by minimum session duration (e.g. 5m, 1h, 30s)")
	exportCmd.Flags().StringVar(&exportMaxDuration, "max-duration", "", "Filter by maximum session duration (e.g. 2h, 45m)")

	// CSV-specific flags
	exportCmd.Flags().StringSliceVar(&exportColumns, "columns", []string{}, "Custom CSV columns (comma-separated)")

	// JSON-specific flags
	exportCmd.Flags().BoolVar(&exportPretty, "pretty", false, "Pretty-print JSON output")
	exportCmd.Flags().BoolVar(&exportCompress, "compress", false, "Compress output file")

	// General flags
	exportCmd.Flags().BoolVar(&exportProgress, "progress", false, "Show progress indicators for large exports")
	exportCmd.Flags().BoolVar(&exportPreview, "preview", false, "Show preview of export without creating file")
	exportCmd.Flags().BoolVar(&exportStats, "stats", false, "Show detailed export statistics")
	exportCmd.Flags().BoolVar(&exportValidate, "validate", false, "Validate export completeness and integrity")

	// Mark output as required only when not in preview mode
	// This will be validated in handleExport function
}

func handleExport(cmd *cobra.Command, args []string) error {
	// Validate export format
	if exportFormat != "csv" && exportFormat != "json" && exportFormat != "xlsx" {
		return fmt.Errorf("unsupported export format: %s (supported: csv, json, xlsx)", exportFormat)
	}

	// Validate output path (not required for preview mode)
	if exportOutput == "" && !exportPreview {
		return fmt.Errorf("output file path is required (unless using --preview mode)")
	}

	// Auto-detect format from file extension if not explicitly set
	if !cmd.Flags().Changed("format") && exportOutput != "" {
		ext := strings.ToLower(filepath.Ext(exportOutput))
		switch ext {
		case ".csv":
			exportFormat = "csv"
		case ".json":
			exportFormat = "json"
		case ".xlsx":
			exportFormat = "xlsx"
		default:
			if ext != "" {
				fmt.Printf("💡 Unknown file extension '%s', using default format: %s\n", ext, exportFormat)
			}
		}
	}

	// Validate format and file extension consistency
	if exportOutput != "" && !exportPreview {
		expectedExt := ""
		switch exportFormat {
		case "csv":
			expectedExt = ".csv"
		case "json":
			expectedExt = ".json"
		case "xlsx":
			expectedExt = ".xlsx"
		}

		actualExt := strings.ToLower(filepath.Ext(exportOutput))
		if actualExt != expectedExt {
			fmt.Printf("⚠️  Warning: File extension '%s' doesn't match format '%s' (expected '%s')\n",
				actualExt, exportFormat, expectedExt)
		}
	}

	// Ensure output directory exists
	outputDir := filepath.Dir(exportOutput)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Initialize database manager
	config := database.DefaultDatabaseConfig()
	manager := database.NewManager(config)

	ctx := context.Background()
	if err := manager.Initialize(ctx); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer manager.Close()

	backend, err := manager.GetBackend()
	if err != nil {
		return fmt.Errorf("failed to get database backend: %w", err)
	}

	// Build export options
	options := &export.ExportOptions{
		Format:       exportFormat,
		Output:       exportOutput,
		Columns:      exportColumns,
		From:         exportFrom,
		To:           exportTo,
		Project:      exportProject,
		Sessions:     exportSessions,
		Status:       exportStatus,
		MinDuration:  exportMinDuration,
		MaxDuration:  exportMaxDuration,
		Pretty:       exportPretty,
		Compress:     exportCompress,
		ShowProgress: exportProgress,
		Preview:      exportPreview,
		ShowStats:    exportStats,
		Validate:     exportValidate,
	}

	// Create exporter based on format
	var exporter export.Exporter
	switch exportFormat {
	case "csv":
		exporter = export.NewCSVExporter()
	case "json":
		exporter = export.NewJSONExporter()
	case "xlsx":
		exporter = export.NewExcelExporter()
	default:
		return fmt.Errorf("unsupported export format: %s", exportFormat)
	}

	// Get filtered sessions
	sessions, err := getFilteredSessions(ctx, backend, options)
	if err != nil {
		return fmt.Errorf("failed to get sessions: %w", err)
	}

	if len(sessions) == 0 {
		fmt.Println("📭 No sessions found matching the specified criteria")
		return nil
	}

	// Handle preview mode
	if options.Preview {
		return handlePreviewMode(ctx, backend, sessions, options)
	}

	fmt.Printf("🔄 Exporting %d sessions to %s format...\n", len(sessions), exportFormat)

	// Perform export with progress if enabled
	if options.ShowProgress {
		if csvExporter, ok := exporter.(*export.CSVExporter); ok {
			if err := csvExporter.ExportWithProgress(ctx, backend, sessions, options); err != nil {
				return fmt.Errorf("export failed: %w", err)
			}
		} else {
			// For other exporters, fall back to regular export
			if err := exporter.Export(ctx, backend, sessions, options); err != nil {
				return fmt.Errorf("export failed: %w", err)
			}
		}
	} else {
		if err := exporter.Export(ctx, backend, sessions, options); err != nil {
			return fmt.Errorf("export failed: %w", err)
		}
	}

	fmt.Printf("✅ Export completed successfully: %s\n", exportOutput)

	// Show file size
	if stat, err := os.Stat(exportOutput); err == nil {
		fmt.Printf("📊 File size: %d bytes\n", stat.Size())
	}

	// Perform export validation if requested
	if options.Validate {
		if err := validateExport(exportOutput, options, len(sessions)); err != nil {
			fmt.Printf("⚠️  Export validation warning: %v\n", err)
		} else {
			fmt.Printf("✅ Export validation passed\n")
		}
	}

	return nil
}

// getFilteredSessions retrieves sessions based on export filters
func getFilteredSessions(ctx context.Context, backend database.DatabaseBackend, options *export.ExportOptions) ([]*database.Session, error) {
	// Build session filters
	filters := &database.SessionFilters{}

	// Parse and apply date filters
	if options.From != "" {
		fromDate, err := parseDate(options.From)
		if err != nil {
			return nil, fmt.Errorf("invalid 'from' date format: %w", err)
		}
		filters.CreatedAfter = &fromDate
	}

	if options.To != "" {
		toDate, err := parseDate(options.To)
		if err != nil {
			return nil, fmt.Errorf("invalid 'to' date format: %w", err)
		}
		// Add 24 hours to include the entire "to" date
		endOfDay := toDate.Add(24 * time.Hour)
		filters.CreatedBefore = &endOfDay
	}

	// Status filter
	if options.Status != "" {
		filters.Status = options.Status
	}

	// Get sessions with database-level filtering
	allSessions, err := backend.ListSessions(ctx, filters)
	if err != nil {
		return nil, err
	}

	var filteredSessions []*database.Session

	for _, session := range allSessions {
		// Project filter (application-level since it requires metadata parsing)
		if options.Project != "" {
			if !matchesProject(session, options.Project) {
				continue
			}
		}

		// Session ID filter (application-level for multiple IDs)
		if len(options.Sessions) > 0 {
			found := false
			for _, sessionID := range options.Sessions {
				if session.ID == sessionID {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// Duration filters
		sessionDuration := session.UpdatedAt.Sub(session.CreatedAt)

		if options.MinDuration != "" {
			minDur, err := parseDuration(options.MinDuration)
			if err != nil {
				return nil, fmt.Errorf("invalid min-duration format: %w", err)
			}
			if sessionDuration < minDur {
				continue
			}
		}

		if options.MaxDuration != "" {
			maxDur, err := parseDuration(options.MaxDuration)
			if err != nil {
				return nil, fmt.Errorf("invalid max-duration format: %w", err)
			}
			if sessionDuration > maxDur {
				continue
			}
		}

		filteredSessions = append(filteredSessions, session)
	}

	return filteredSessions, nil
}

// parseDate parses a date string in YYYY-MM-DD format
func parseDate(dateStr string) (time.Time, error) {
	// Try common date formats
	formats := []string{
		"2006-01-02",           // YYYY-MM-DD
		"2006-01-02 15:04:05",  // YYYY-MM-DD HH:MM:SS
		"01/02/2006",           // MM/DD/YYYY
		"2006/01/02",           // YYYY/MM/DD
	}

	for _, format := range formats {
		if parsed, err := time.Parse(format, dateStr); err == nil {
			return parsed, nil
		}
	}

	return time.Time{}, fmt.Errorf("unsupported date format: %s (supported: YYYY-MM-DD, YYYY-MM-DD HH:MM:SS, MM/DD/YYYY, YYYY/MM/DD)", dateStr)
}

// parseDuration parses a duration string like "5m", "1h30m", "2h", "45s"
func parseDuration(durationStr string) (time.Duration, error) {
	return time.ParseDuration(durationStr)
}

// matchesProject checks if a session matches the project filter
func matchesProject(session *database.Session, projectFilter string) bool {
	if session.Metadata == "" {
		return false
	}

	// Simple string matching in metadata (case-insensitive)
	metadataLower := strings.ToLower(session.Metadata)
	projectLower := strings.ToLower(projectFilter)

	// Check if project name appears in metadata
	if strings.Contains(metadataLower, fmt.Sprintf("\"project\":\"%s\"", projectLower)) {
		return true
	}

	// Check for partial matches in working directory or data fields
	if strings.Contains(metadataLower, projectLower) {
		return true
	}

	return false
}

// handlePreviewMode shows a preview of the export without creating a file
func handlePreviewMode(ctx context.Context, backend database.DatabaseBackend, sessions []*database.Session, options *export.ExportOptions) error {
	// Limit to first 5 sessions for preview
	previewSessions := sessions
	if len(sessions) > 5 {
		previewSessions = sessions[:5]
	}

	// Prepare export data
	exportData, err := export.PrepareSessionData(ctx, backend, previewSessions)
	if err != nil {
		return fmt.Errorf("failed to prepare preview data: %w", err)
	}

	fmt.Printf("📋 Export Preview (showing first %d of %d sessions):\n\n", len(exportData), len(sessions))

	if options.Format == "csv" {
		// Show CSV preview as a table
		columns := options.Columns
		if len(columns) == 0 {
			columns = export.DefaultCSVColumns
		}

		// Print headers
		fmt.Print("| ")
		for _, col := range columns {
			fmt.Printf("%-15s | ", col)
		}
		fmt.Println()

		// Print separator
		fmt.Print("| ")
		for range columns {
			fmt.Print("--------------- | ")
		}
		fmt.Println()

		// Print data rows
		csvExporter := export.NewCSVExporter()
		for _, session := range exportData {
			fmt.Print("| ")
			for _, col := range columns {
				value := csvExporter.GetColumnValue(session, col)
				if len(value) > 13 {
					value = value[:10] + "..."
				}
				fmt.Printf("%-15s | ", value)
			}
			fmt.Println()
		}
	} else {
		// Show JSON preview (first session structure)
		if len(exportData) > 0 {
			fmt.Println("JSON Structure Preview:")
			previewData := map[string]interface{}{
				"session_metadata": map[string]interface{}{
					"id":           exportData[0].SessionID,
					"project_name": exportData[0].ProjectName,
					"start_time":   exportData[0].StartTime,
					"duration":     exportData[0].Duration,
					"status":       exportData[0].Status,
				},
				"analytics": map[string]interface{}{
					"user_prompts":    exportData[0].UserPrompts,
					"claude_replies":  exportData[0].ClaudeReplies,
					"total_words":     exportData[0].TotalWords,
					"session_tags":    exportData[0].SessionTags,
				},
			}

			jsonData, _ := json.MarshalIndent(previewData, "", "  ")
			fmt.Println(string(jsonData))
		}
	}

	// Show summary statistics
	fmt.Printf("\n💡 Full export would include:\n")
	fmt.Printf("   - %d sessions total\n", len(sessions))

	// Calculate estimated file size
	estimatedSize := len(sessions) * 100 // rough estimate
	if options.Format == "json" {
		estimatedSize *= 10 // JSON is larger
	}
	fmt.Printf("   - Estimated file size: ~%d KB\n", estimatedSize/1024)

	if options.Format == "csv" {
		columns := options.Columns
		if len(columns) == 0 {
			columns = export.DefaultCSVColumns
		}
		fmt.Printf("   - CSV columns: %s\n", strings.Join(columns, ", "))
	}

	fmt.Printf("\n🔄 Run without --preview to perform actual export\n")
	return nil
}

// validateExport performs basic validation on the exported file
func validateExport(filePath string, options *export.ExportOptions, expectedSessions int) error {
	// Check if file exists and is not empty
	stat, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("exported file not found: %w", err)
	}

	if stat.Size() == 0 {
		return fmt.Errorf("exported file is empty")
	}

	// Basic format-specific validation
	switch options.Format {
	case "csv":
		return validateCSVExport(filePath, expectedSessions)
	case "json":
		return validateJSONExport(filePath)
	case "xlsx":
		return validateExcelExport(filePath)
	}

	return nil
}

// validateCSVExport validates CSV file structure
func validateCSVExport(filePath string, expectedSessions int) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Count lines (should be expectedSessions + 1 for header)
	lineCount := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineCount++
	}

	if lineCount < 1 {
		return fmt.Errorf("CSV file has no header")
	}

	actualSessions := lineCount - 1
	if actualSessions != expectedSessions {
		return fmt.Errorf("session count mismatch: expected %d, found %d", expectedSessions, actualSessions)
	}

	return nil
}

// validateJSONExport validates JSON file structure
func validateJSONExport(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Try to parse as JSON
	var data interface{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return fmt.Errorf("invalid JSON format: %w", err)
	}

	return nil
}

// validateExcelExport validates Excel file structure
func validateExcelExport(filePath string) error {
	// For Excel files, just check that the file can be opened
	// More detailed validation would require parsing the Excel file
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Check if it starts with the Excel file signature
	header := make([]byte, 4)
	if _, err := file.Read(header); err != nil {
		return fmt.Errorf("failed to read file header: %w", err)
	}

	// Excel files start with PK (ZIP signature)
	if header[0] != 'P' || header[1] != 'K' {
		return fmt.Errorf("file does not appear to be a valid Excel file")
	}

	return nil
}