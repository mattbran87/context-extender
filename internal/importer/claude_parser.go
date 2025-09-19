package importer

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// ClaudeEntry represents a single line in Claude's JSONL file
type ClaudeEntry struct {
	Type        string          `json:"type"`
	ParentUUID  string          `json:"parentUuid,omitempty"`
	UserType    string          `json:"userType,omitempty"`
	CWD         string          `json:"cwd,omitempty"`
	SessionID   string          `json:"sessionId,omitempty"`
	Version     string          `json:"version,omitempty"`
	Message     *ClaudeMessage  `json:"message,omitempty"`
	Summary     string          `json:"summary,omitempty"`
	LeafUUID    string          `json:"leafUuid,omitempty"`
	Timestamp   string          `json:"timestamp,omitempty"`
	IsSidechain bool            `json:"isSidechain,omitempty"`
	GitBranch   string          `json:"gitBranch,omitempty"`
}

// ClaudeMessage represents a message in the conversation
type ClaudeMessage struct {
	Role    string                 `json:"role"`
	Content []ClaudeMessageContent `json:"content"`
}

// ClaudeMessageContent represents content within a message
type ClaudeMessageContent struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

// ClaudeConversation represents a parsed Claude conversation
type ClaudeConversation struct {
	SessionID   string
	ProjectPath string
	StartTime   time.Time
	EndTime     time.Time
	Messages    []ParsedMessage
	Summaries   []string
	Metadata    map[string]interface{}
}

// ParsedMessage represents a normalized message
type ParsedMessage struct {
	ID        string
	ParentID  string
	Role      string // "user" or "assistant"
	Content   string
	Timestamp time.Time
	Metadata  map[string]interface{}
}

// ClaudeParser parses Claude JSONL files
type ClaudeParser struct {
	verbose bool
}

// NewClaudeParser creates a new Claude parser
func NewClaudeParser(verbose bool) *ClaudeParser {
	return &ClaudeParser{
		verbose: verbose,
	}
}

// ParseFile parses a single Claude JSONL file
func (cp *ClaudeParser) ParseFile(filePath string) (*ClaudeConversation, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	conversation := &ClaudeConversation{
		Messages:  []ParsedMessage{},
		Summaries: []string{},
		Metadata:  make(map[string]interface{}),
	}

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		if line == "" {
			continue
		}

		var entry ClaudeEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			if cp.verbose {
				fmt.Printf("Warning: Failed to parse line %d: %v\n", lineNum, err)
			}
			continue
		}

		// Process entry based on type
		switch entry.Type {
		case "user", "assistant":
			if err := cp.processMessage(&entry, conversation); err != nil {
				if cp.verbose {
					fmt.Printf("Warning: Failed to process message at line %d: %v\n", lineNum, err)
				}
			}

		case "summary":
			if entry.Summary != "" {
				conversation.Summaries = append(conversation.Summaries, entry.Summary)
			}

		case "session-start", "session_start":
			cp.processSessionStart(&entry, conversation)

		case "session-end", "session_end":
			cp.processSessionEnd(&entry, conversation)

		default:
			if cp.verbose {
				fmt.Printf("Info: Unknown entry type '%s' at line %d\n", entry.Type, lineNum)
			}
		}

		// Update session info
		if entry.SessionID != "" && conversation.SessionID == "" {
			conversation.SessionID = entry.SessionID
		}
		if entry.CWD != "" && conversation.ProjectPath == "" {
			conversation.ProjectPath = entry.CWD
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	// Post-process conversation
	cp.postProcess(conversation)

	return conversation, nil
}

// processMessage processes a message entry
func (cp *ClaudeParser) processMessage(entry *ClaudeEntry, conv *ClaudeConversation) error {
	if entry.Message == nil {
		return fmt.Errorf("message entry has no message content")
	}

	// Extract text content
	var textContent strings.Builder
	for _, content := range entry.Message.Content {
		if content.Type == "text" && content.Text != "" {
			textContent.WriteString(content.Text)
			textContent.WriteString("\n")
		}
	}

	message := ParsedMessage{
		ID:       entry.LeafUUID,
		ParentID: entry.ParentUUID,
		Role:     entry.Message.Role,
		Content:  strings.TrimSpace(textContent.String()),
		Metadata: map[string]interface{}{
			"user_type":   entry.UserType,
			"version":     entry.Version,
			"is_sidechain": entry.IsSidechain,
			"git_branch":  entry.GitBranch,
		},
	}

	// Parse timestamp if available
	if entry.Timestamp != "" {
		if t, err := time.Parse(time.RFC3339, entry.Timestamp); err == nil {
			message.Timestamp = t
		}
	}

	conv.Messages = append(conv.Messages, message)
	return nil
}

// processSessionStart processes session start events
func (cp *ClaudeParser) processSessionStart(entry *ClaudeEntry, conv *ClaudeConversation) {
	if entry.Timestamp != "" {
		if t, err := time.Parse(time.RFC3339, entry.Timestamp); err == nil {
			conv.StartTime = t
		}
	}
}

// processSessionEnd processes session end events
func (cp *ClaudeParser) processSessionEnd(entry *ClaudeEntry, conv *ClaudeConversation) {
	if entry.Timestamp != "" {
		if t, err := time.Parse(time.RFC3339, entry.Timestamp); err == nil {
			conv.EndTime = t
		}
	}
}

// postProcess performs post-processing on the conversation
func (cp *ClaudeParser) postProcess(conv *ClaudeConversation) {
	// Set session ID if not found
	if conv.SessionID == "" {
		conv.SessionID = uuid.New().String()
	}

	// Calculate time bounds if not set
	if len(conv.Messages) > 0 {
		if conv.StartTime.IsZero() && !conv.Messages[0].Timestamp.IsZero() {
			conv.StartTime = conv.Messages[0].Timestamp
		}
		if conv.EndTime.IsZero() && !conv.Messages[len(conv.Messages)-1].Timestamp.IsZero() {
			conv.EndTime = conv.Messages[len(conv.Messages)-1].Timestamp
		}
	}

	// If still no timestamps, use file modification time as approximation
	if conv.StartTime.IsZero() {
		conv.StartTime = time.Now().Add(-24 * time.Hour) // Default to 24 hours ago
	}
	if conv.EndTime.IsZero() {
		conv.EndTime = time.Now()
	}

	// Store additional metadata
	conv.Metadata["message_count"] = len(conv.Messages)
	conv.Metadata["summary_count"] = len(conv.Summaries)
	conv.Metadata["duration"] = conv.EndTime.Sub(conv.StartTime).String()
}

// FindClaudeConversations finds all Claude conversation files
func FindClaudeConversations() ([]string, error) {
	var conversationFiles []string

	// Check standard Claude directories
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	claudeDirs := []string{
		filepath.Join(homeDir, ".claude", "projects"),
		filepath.Join(homeDir, "Library", "Application Support", "Claude", "projects"),
		filepath.Join(homeDir, "AppData", "Roaming", "Claude", "projects"),
	}

	for _, dir := range claudeDirs {
		if _, err := os.Stat(dir); err != nil {
			continue // Directory doesn't exist
		}

		// Walk through all subdirectories
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil // Skip errors
			}

			// Look for JSONL files
			if !info.IsDir() && strings.HasSuffix(info.Name(), ".jsonl") {
				conversationFiles = append(conversationFiles, path)
			}

			return nil
		})

		if err != nil {
			return nil, fmt.Errorf("error walking directory %s: %w", dir, err)
		}
	}

	return conversationFiles, nil
}

// GetProjectName extracts project name from the directory path
func GetProjectName(dirPath string) string {
	// Claude uses format like "C--Users-marko-IdeaProjects-context-extender"
	// Convert back to readable format
	projectName := filepath.Base(dirPath)

	// Replace common separators
	projectName = strings.ReplaceAll(projectName, "--", "/")
	projectName = strings.ReplaceAll(projectName, "-", " ")

	// Extract just the last part for brevity
	parts := strings.Split(projectName, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}

	return projectName
}