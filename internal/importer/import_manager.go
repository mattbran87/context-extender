package importer

import (
	"context-extender/internal/database"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// ImportManager manages the import of Claude conversations
type ImportManager struct {
	parser       *ClaudeParser
	verbose      bool
	dryRun       bool
	skipExisting bool
}

// ImportOptions configures import behavior
type ImportOptions struct {
	Verbose      bool
	DryRun       bool
	SkipExisting bool
	MaxFiles     int
}

// ImportResult contains import statistics
type ImportResult struct {
	TotalFiles       int
	SuccessfulFiles  int
	FailedFiles      int
	SkippedFiles     int
	TotalMessages    int
	TotalSessions    int
	ImportDuration   time.Duration
	Errors           []string
}

// NewImportManager creates a new import manager
func NewImportManager(options ImportOptions) *ImportManager {
	return &ImportManager{
		parser:       NewClaudeParser(options.Verbose),
		verbose:      options.Verbose,
		dryRun:       options.DryRun,
		skipExisting: options.SkipExisting,
	}
}

// ImportFile imports a single Claude JSONL file
func (im *ImportManager) ImportFile(filePath string) error {
	if im.verbose {
		fmt.Printf("Importing file: %s\n", filePath)
	}

	// Check if already imported
	if im.skipExisting {
		if exists, err := im.isFileImported(filePath); err == nil && exists {
			if im.verbose {
				fmt.Println("  File already imported, skipping...")
			}
			return nil
		}
	}

	// Parse the file
	conversation, err := im.parser.ParseFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to parse file: %w", err)
	}

	if im.verbose {
		fmt.Printf("  Found %d messages, %d summaries\n",
			len(conversation.Messages), len(conversation.Summaries))
	}

	// If dry run, don't actually import
	if im.dryRun {
		fmt.Printf("  [DRY RUN] Would import %d messages from session %s\n",
			len(conversation.Messages), conversation.SessionID)
		return nil
	}

	// Import to database
	if err := im.importConversation(conversation, filePath); err != nil {
		return fmt.Errorf("failed to import conversation: %w", err)
	}

	// Record import
	if err := im.recordImport(filePath, conversation); err != nil {
		return fmt.Errorf("failed to record import: %w", err)
	}

	if im.verbose {
		fmt.Printf("  ✅ Successfully imported session %s\n", conversation.SessionID)
	}

	return nil
}

// ImportDirectory imports all JSONL files from a directory
func (im *ImportManager) ImportDirectory(dirPath string) (*ImportResult, error) {
	result := &ImportResult{
		Errors: []string{},
	}
	startTime := time.Now()

	// Find all JSONL files
	var files []string
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors
		}

		if !info.IsDir() && filepath.Ext(path) == ".jsonl" {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		return result, fmt.Errorf("failed to walk directory: %w", err)
	}

	result.TotalFiles = len(files)

	// Import each file
	for _, file := range files {
		if err := im.ImportFile(file); err != nil {
			result.FailedFiles++
			result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", file, err))
			if im.verbose {
				fmt.Printf("  ❌ Failed to import %s: %v\n", file, err)
			}
		} else {
			result.SuccessfulFiles++
		}
	}

	result.ImportDuration = time.Since(startTime)
	return result, nil
}

// ImportAllClaude imports all Claude conversations found on the system
func (im *ImportManager) ImportAllClaude() (*ImportResult, error) {
	result := &ImportResult{
		Errors: []string{},
	}
	startTime := time.Now()

	// Find all Claude conversation files
	files, err := FindClaudeConversations()
	if err != nil {
		return result, fmt.Errorf("failed to find Claude conversations: %w", err)
	}

	result.TotalFiles = len(files)

	if im.verbose {
		fmt.Printf("Found %d Claude conversation files\n", len(files))
	}

	// Import each file
	for i, file := range files {
		if im.verbose {
			fmt.Printf("\n[%d/%d] Processing %s\n", i+1, len(files), filepath.Base(file))
		}

		if err := im.ImportFile(file); err != nil {
			result.FailedFiles++
			result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", file, err))
		} else {
			result.SuccessfulFiles++
		}
	}

	result.ImportDuration = time.Since(startTime)
	return result, nil
}

// importConversation imports a parsed conversation to the database
func (im *ImportManager) importConversation(conv *ClaudeConversation, filePath string) error {
	// Create session
	metadataMap := map[string]string{
		"source":       "claude",
		"project_path": conv.ProjectPath,
		"file_path":    filePath,
		"import_date":  time.Now().Format(time.RFC3339),
	}
	metadataJSON, _ := json.Marshal(metadataMap)

	session := &database.Session{
		ID:        conv.SessionID,
		CreatedAt: conv.StartTime,
		UpdatedAt: conv.EndTime,
		Status:    "imported",
		Metadata:  string(metadataJSON),
	}

	if err := database.CreateSession(session); err != nil {
		// Session might already exist, try to update
		metadataMap := map[string]string{
			"source":       "claude",
			"project_path": conv.ProjectPath,
			"file_path":    filePath,
			"import_date":  time.Now().Format(time.RFC3339),
		}
		if err := database.UpdateSession(session.ID, "imported", metadataMap); err != nil {
			return fmt.Errorf("failed to create/update session: %w", err)
		}
	}

	// Import messages as conversations
	for i, msg := range conv.Messages {
		conversation := &database.Conversation{
			SessionID:   conv.SessionID,
			MessageType: msg.Role,
			Content:     msg.Content,
			Timestamp:   msg.Timestamp,
		}

		// Add metadata
		if metadata, err := json.Marshal(msg.Metadata); err == nil {
			conversation.Metadata = string(metadata)
		}

		if err := database.CreateConversation(conversation); err != nil {
			if im.verbose {
				fmt.Printf("  Warning: Failed to import message %d: %v\n", i+1, err)
			}
		}
	}

	// Create events for session lifecycle
	startEvent := &database.Event{
		SessionID: conv.SessionID,
		EventType: "import_start",
		Data: fmt.Sprintf(`{"file":"%s","time":"%s"}`,
			filepath.Base(filePath), conv.StartTime.Format(time.RFC3339)),
		Timestamp:   conv.StartTime,
		SequenceNum: 1,
	}
	database.CreateEvent(startEvent)

	endEvent := &database.Event{
		SessionID: conv.SessionID,
		EventType: "import_end",
		Data: fmt.Sprintf(`{"file":"%s","time":"%s","messages":%d}`,
			filepath.Base(filePath), conv.EndTime.Format(time.RFC3339), len(conv.Messages)),
		Timestamp:   conv.EndTime,
		SequenceNum: 2,
	}
	database.CreateEvent(endEvent)

	return nil
}

// recordImport records the import in the database
func (im *ImportManager) recordImport(filePath string, conv *ClaudeConversation) error {
	// Calculate file checksum
	checksum, err := im.calculateChecksum(filePath)
	if err != nil {
		checksum = "" // Non-fatal error
	}

	// Get database connection
	db, err := database.GetConnection()
	if err != nil {
		return err
	}

	// Insert or update import record
	query := `
		INSERT OR REPLACE INTO import_history
		(file_path, imported_at, session_count, event_count, checksum)
		VALUES (?, ?, ?, ?, ?)
	`

	_, err = db.Exec(query,
		filePath,
		time.Now().Format(time.RFC3339),
		1, // One session per file
		len(conv.Messages),
		checksum,
	)

	return err
}

// isFileImported checks if a file has already been imported
func (im *ImportManager) isFileImported(filePath string) (bool, error) {
	db, err := database.GetConnection()
	if err != nil {
		return false, err
	}

	var count int
	query := "SELECT COUNT(*) FROM import_history WHERE file_path = ?"
	err = db.QueryRow(query, filePath).Scan(&count)

	return count > 0, err
}

// calculateChecksum calculates MD5 checksum of a file
func (im *ImportManager) calculateChecksum(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// GetImportHistory returns the import history
func GetImportHistory() ([]ImportRecord, error) {
	db, err := database.GetConnection()
	if err != nil {
		return nil, err
	}

	query := `
		SELECT file_path, imported_at, session_count, event_count, checksum
		FROM import_history
		ORDER BY imported_at DESC
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []ImportRecord
	for rows.Next() {
		var record ImportRecord
		err := rows.Scan(
			&record.FilePath,
			&record.ImportedAt,
			&record.SessionCount,
			&record.EventCount,
			&record.Checksum,
		)
		if err != nil {
			continue
		}
		records = append(records, record)
	}

	return records, rows.Err()
}

// ImportRecord represents an import history record
type ImportRecord struct {
	FilePath     string
	ImportedAt   string
	SessionCount int
	EventCount   int
	Checksum     string
}