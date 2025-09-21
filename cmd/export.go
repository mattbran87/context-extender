package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	exportFormat   string
	exportOutput   string
	exportColumns  []string
	exportFrom     string
	exportTo       string
	exportProject  string
	exportSessions []string
	exportPretty   bool
	exportCompress bool
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

	// CSV-specific flags
	exportCmd.Flags().StringSliceVar(&exportColumns, "columns", []string{}, "Custom CSV columns (comma-separated)")

	// JSON-specific flags
	exportCmd.Flags().BoolVar(&exportPretty, "pretty", false, "Pretty-print JSON output")
	exportCmd.Flags().BoolVar(&exportCompress, "compress", false, "Compress output file")

	// Mark required flags
	exportCmd.MarkFlagRequired("output")
}

func handleExport(cmd *cobra.Command, args []string) error {
	// Validate export format
	if exportFormat != "csv" && exportFormat != "json" {
		return fmt.Errorf("unsupported export format: %s (supported: csv, json)", exportFormat)
	}

	// Validate output path
	if exportOutput == "" {
		return fmt.Errorf("output file path is required")
	}

	// Auto-detect format from file extension if not explicitly set
	if !cmd.Flags().Changed("format") {
		ext := strings.ToLower(filepath.Ext(exportOutput))
		switch ext {
		case ".csv":
			exportFormat = "csv"
		case ".json":
			exportFormat = "json"
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
		Format:    exportFormat,
		Output:    exportOutput,
		Columns:   exportColumns,
		From:      exportFrom,
		To:        exportTo,
		Project:   exportProject,
		Sessions:  exportSessions,
		Pretty:    exportPretty,
		Compress:  exportCompress,
	}

	// Create exporter based on format
	var exporter export.Exporter
	switch exportFormat {
	case "csv":
		exporter = export.NewCSVExporter()
	case "json":
		exporter = export.NewJSONExporter()
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

	fmt.Printf("🔄 Exporting %d sessions to %s format...\n", len(sessions), exportFormat)

	// Perform export
	if err := exporter.Export(ctx, backend, sessions, options); err != nil {
		return fmt.Errorf("export failed: %w", err)
	}

	fmt.Printf("✅ Export completed successfully: %s\n", exportOutput)

	// Show file size
	if stat, err := os.Stat(exportOutput); err == nil {
		fmt.Printf("📊 File size: %d bytes\n", stat.Size())
	}

	return nil
}

// getFilteredSessions retrieves sessions based on export filters
func getFilteredSessions(ctx context.Context, backend database.DatabaseBackend, options *export.ExportOptions) ([]*database.Session, error) {
	// Build session filters
	filters := &database.SessionFilters{}

	// Date filters
	if options.From != "" {
		// TODO: Parse date and set CreatedAfter filter
	}
	if options.To != "" {
		// TODO: Parse date and set CreatedBefore filter
	}

	// Get all sessions (filtering will be applied later)
	allSessions, err := backend.ListSessions(ctx, filters)
	if err != nil {
		return nil, err
	}

	var filteredSessions []*database.Session

	for _, session := range allSessions {
		// Project filter
		if options.Project != "" {
			// Parse metadata to check project
			if session.Metadata != "" {
				// TODO: Implement project filtering based on metadata
				// For now, skip project filtering - will be implemented in next iteration
			}
		}

		// Session ID filter
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

		filteredSessions = append(filteredSessions, session)
	}

	return filteredSessions, nil
}