package export

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
	"context-extender/internal/database"
)

// ExcelExporter implements the Exporter interface for Excel (.xlsx) format
type ExcelExporter struct {
	supportedColumns []string
}

// NewExcelExporter creates a new Excel exporter
func NewExcelExporter() *ExcelExporter {
	return &ExcelExporter{
		supportedColumns: AllAvailableColumns,
	}
}

// Export performs the Excel export operation
func (e *ExcelExporter) Export(ctx context.Context, backend database.DatabaseBackend, sessions []*database.Session, options *ExportOptions) error {
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

	// Create a new Excel file
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("Warning: failed to close Excel file: %v\n", err)
		}
	}()

	// Set up the worksheet
	sheetName := "Sessions"
	f.SetSheetName("Sheet1", sheetName)

	// Set headers with styling
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 12,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#E7E6E6"},
			Pattern: 1,
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create header style: %w", err)
	}

	// Write headers
	for i, col := range columns {
		cell := fmt.Sprintf("%s1", getColumnLetter(i))
		f.SetCellValue(sheetName, cell, col)
		f.SetCellStyle(sheetName, cell, cell, headerStyle)
	}

	// Create data row style
	dataStyle, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "CCCCCC", Style: 1},
			{Type: "top", Color: "CCCCCC", Style: 1},
			{Type: "bottom", Color: "CCCCCC", Style: 1},
			{Type: "right", Color: "CCCCCC", Style: 1},
		},
		Alignment: &excelize.Alignment{
			Vertical: "top",
			WrapText: true,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create data style: %w", err)
	}

	// Write data rows with progress reporting
	total := len(exportData)
	for rowIdx, session := range exportData {
		excelRow := rowIdx + 2 // Excel rows are 1-indexed, plus header row

		for colIdx, column := range columns {
			cell := fmt.Sprintf("%s%d", getColumnLetter(colIdx), excelRow)
			value := e.getColumnValue(session, column)

			// Set appropriate cell type based on column
			if e.isNumericColumn(column) {
				if numValue, err := strconv.Atoi(value); err == nil {
					f.SetCellValue(sheetName, cell, numValue)
				} else {
					f.SetCellValue(sheetName, cell, value)
				}
			} else {
				f.SetCellValue(sheetName, cell, value)
			}

			f.SetCellStyle(sheetName, cell, cell, dataStyle)
		}

		// Show progress for large exports
		if options.ShowProgress && total > 10 && (rowIdx+1)%(total/10) == 0 {
			percent := float64(rowIdx+1) / float64(total) * 100
			fmt.Printf("📊 Excel progress: %.0f%% (%d/%d sessions)\n", percent, rowIdx+1, total)
		}
	}

	// Auto-size columns for better readability
	for i := range columns {
		col := getColumnLetter(i)
		f.SetColWidth(sheetName, col, col, 15)
	}

	// Add a summary sheet with statistics
	if err := e.addSummarySheet(f, exportData, options); err != nil {
		fmt.Printf("Warning: failed to create summary sheet: %v\n", err)
	}

	// Save the file
	if err := f.SaveAs(options.Output); err != nil {
		return fmt.Errorf("failed to save Excel file: %w", err)
	}

	return nil
}

// addSummarySheet creates a summary statistics sheet
func (e *ExcelExporter) addSummarySheet(f *excelize.File, exportData []*SessionExportData, options *ExportOptions) error {
	summarySheet := "Summary"
	index, err := f.NewSheet(summarySheet)
	if err != nil {
		return err
	}

	// Set the summary sheet as active
	f.SetActiveSheet(index)

	// Calculate summary statistics
	totalSessions := len(exportData)
	totalUserPrompts := 0
	totalClaudeReplies := 0
	totalWords := 0
	statusCount := make(map[string]int)
	tagCount := make(map[string]int)

	for _, session := range exportData {
		totalUserPrompts += session.UserPrompts
		totalClaudeReplies += session.ClaudeReplies
		totalWords += session.TotalWords
		statusCount[session.Status]++

		for _, tag := range session.SessionTags {
			tagCount[tag]++
		}
	}

	// Write summary data
	f.SetCellValue(summarySheet, "A1", "Export Summary")
	f.SetCellValue(summarySheet, "A3", "Total Sessions:")
	f.SetCellValue(summarySheet, "B3", totalSessions)
	f.SetCellValue(summarySheet, "A4", "Total User Prompts:")
	f.SetCellValue(summarySheet, "B4", totalUserPrompts)
	f.SetCellValue(summarySheet, "A5", "Total Claude Replies:")
	f.SetCellValue(summarySheet, "B5", totalClaudeReplies)
	f.SetCellValue(summarySheet, "A6", "Total Words:")
	f.SetCellValue(summarySheet, "B6", totalWords)

	// Status breakdown
	row := 8
	f.SetCellValue(summarySheet, "A"+strconv.Itoa(row), "Status Breakdown:")
	row++
	for status, count := range statusCount {
		f.SetCellValue(summarySheet, "A"+strconv.Itoa(row), status+":")
		f.SetCellValue(summarySheet, "B"+strconv.Itoa(row), count)
		row++
	}

	// Top session tags
	row += 2
	f.SetCellValue(summarySheet, "A"+strconv.Itoa(row), "Top Session Tags:")
	row++

	type tagInfo struct {
		tag   string
		count int
	}
	var tags []tagInfo
	for tag, count := range tagCount {
		tags = append(tags, tagInfo{tag, count})
	}

	// Show top 10 tags
	maxTags := 10
	if len(tags) > maxTags {
		tags = tags[:maxTags]
	}

	for _, tagInfo := range tags {
		f.SetCellValue(summarySheet, "A"+strconv.Itoa(row), tagInfo.tag+":")
		f.SetCellValue(summarySheet, "B"+strconv.Itoa(row), tagInfo.count)
		row++
	}

	return nil
}

// getColumnValue extracts the value for a specific column from session data
func (e *ExcelExporter) getColumnValue(session *SessionExportData, column string) string {
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
		return session.RawMetadata
	default:
		return ""
	}
}

// isNumericColumn determines if a column contains numeric data
func (e *ExcelExporter) isNumericColumn(column string) bool {
	numericColumns := map[string]bool{
		"event_count":        true,
		"user_prompts":       true,
		"claude_replies":     true,
		"total_words":        true,
		"user_words":         true,
		"claude_words":       true,
		"compression_events": true,
		"tool_usage_count":   true,
	}
	return numericColumns[column]
}

// getColumnLetter converts column index to Excel column letter (A, B, C, ..., Z, AA, AB, ...)
func getColumnLetter(colIndex int) string {
	result := ""
	for colIndex >= 0 {
		result = string(rune('A'+colIndex%26)) + result
		colIndex = colIndex/26 - 1
	}
	return result
}

// GetSupportedColumns returns the list of columns this exporter supports
func (e *ExcelExporter) GetSupportedColumns() []string {
	return e.supportedColumns
}

// ValidateOptions validates the export options for Excel export
func (e *ExcelExporter) ValidateOptions(options *ExportOptions) error {
	// Validate format
	if options.Format != "xlsx" {
		return fmt.Errorf("Excel exporter only supports 'xlsx' format, got: %s", options.Format)
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

	// Excel-specific validations
	if options.Pretty {
		fmt.Println("💡 Note: Pretty-print option is not applicable for Excel format (formatting is automatic)")
	}

	if options.Compress {
		return fmt.Errorf("compression is not supported for Excel export (Excel files are already compressed)")
	}

	return nil
}