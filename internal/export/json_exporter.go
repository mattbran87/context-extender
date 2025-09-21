package export

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
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

	// Create export structure
	exportStructure := map[string]interface{}{
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

	// Create output file
	file, err := os.Create(options.Output)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	// Create JSON encoder
	encoder := json.NewEncoder(file)
	if options.Pretty {
		encoder.SetIndent("", "  ")
	}

	// Write JSON data
	if err := encoder.Encode(exportStructure); err != nil {
		return fmt.Errorf("failed to write JSON data: %w", err)
	}

	return nil
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

	// Compression not yet supported
	if options.Compress {
		return fmt.Errorf("compression is not yet supported for JSON export")
	}

	return nil
}