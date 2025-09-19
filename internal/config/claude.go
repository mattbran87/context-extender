package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// HookConfig represents a single hook configuration
type HookConfig struct {
	Type    string `json:"type"`
	Command string `json:"command"`
	Timeout int    `json:"timeout,omitempty"`
}

// HookEntry represents a complete hook entry with matcher
type HookEntry struct {
	Matcher string       `json:"matcher"`
	Hooks   []HookConfig `json:"hooks"`
}

// ClaudeSettings represents the structure of Claude Code's settings.json
type ClaudeSettings struct {
	Hooks map[string][]HookEntry `json:"hooks,omitempty"`
	// Add other fields as needed, using interface{} to preserve unknown fields
	Other map[string]interface{} `json:"-"`
}

// MarshalJSON custom marshaling to preserve unknown fields
func (cs *ClaudeSettings) MarshalJSON() ([]byte, error) {
	// Create a temporary map to combine known and unknown fields
	result := make(map[string]interface{})

	// Add known fields
	if cs.Hooks != nil {
		result["hooks"] = cs.Hooks
	}

	// Add unknown fields
	for k, v := range cs.Other {
		result[k] = v
	}

	return json.Marshal(result)
}

// UnmarshalJSON custom unmarshaling to preserve unknown fields
func (cs *ClaudeSettings) UnmarshalJSON(data []byte) error {
	// First unmarshal into a generic map
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Initialize the Other map
	cs.Other = make(map[string]interface{})

	// Extract known fields
	if hooks, exists := raw["hooks"]; exists {
		hooksBytes, err := json.Marshal(hooks)
		if err != nil {
			return fmt.Errorf("failed to marshal hooks: %w", err)
		}
		if err := json.Unmarshal(hooksBytes, &cs.Hooks); err != nil {
			return fmt.Errorf("failed to unmarshal hooks: %w", err)
		}
		delete(raw, "hooks")
	}

	// Store unknown fields
	for k, v := range raw {
		cs.Other[k] = v
	}

	return nil
}

// ReadClaudeSettings reads and parses Claude Code's settings.json
func ReadClaudeSettings() (*ClaudeSettings, error) {
	settingsPath, err := GetClaudeSettingsPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get Claude settings path: %w", err)
	}

	// Check if file exists
	if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
		// Return empty settings if file doesn't exist
		return &ClaudeSettings{
			Other: make(map[string]interface{}),
		}, nil
	}

	file, err := os.Open(settingsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open settings file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read settings file: %w", err)
	}

	// Handle empty file
	if len(data) == 0 {
		return &ClaudeSettings{
			Other: make(map[string]interface{}),
		}, nil
	}

	var settings ClaudeSettings
	if err := json.Unmarshal(data, &settings); err != nil {
		return nil, fmt.Errorf("failed to parse settings JSON: %w", err)
	}

	// Ensure Other map is initialized
	if settings.Other == nil {
		settings.Other = make(map[string]interface{})
	}

	return &settings, nil
}

// WriteClaudeSettings writes Claude Code settings with backup
func WriteClaudeSettings(settings *ClaudeSettings) error {
	settingsPath, err := GetClaudeSettingsPath()
	if err != nil {
		return fmt.Errorf("failed to get Claude settings path: %w", err)
	}

	// Create backup before writing
	if err := CreateBackup(settingsPath); err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	// Ensure parent directory exists
	if err := EnsureDirectoryExists(filepath.Dir(settingsPath)); err != nil {
		return fmt.Errorf("failed to create settings directory: %w", err)
	}

	// Marshal settings to JSON with pretty printing
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}

	// Write atomically using a temporary file
	tempPath := settingsPath + ".tmp"
	if err := os.WriteFile(tempPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write temporary file: %w", err)
	}

	// Atomic rename
	if err := os.Rename(tempPath, settingsPath); err != nil {
		// Clean up temp file on failure
		os.Remove(tempPath)
		return fmt.Errorf("failed to rename temporary file: %w", err)
	}

	return nil
}

// CreateBackup creates a timestamped backup of the settings file
func CreateBackup(settingsPath string) error {
	// Check if original file exists
	if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
		// No file to backup
		return nil
	}

	// Create backup filename with timestamp
	timestamp := time.Now().Format("20060102-150405")
	backupPath := fmt.Sprintf("%s.backup.%s", settingsPath, timestamp)

	// Copy file
	return copyFile(settingsPath, backupPath)
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

// ValidateClaudeInstallation checks if Claude Code appears to be installed
func ValidateClaudeInstallation() error {
	settingsPath, err := GetClaudeSettingsPath()
	if err != nil {
		return fmt.Errorf("failed to get Claude settings path: %w", err)
	}

	// Check if Claude directory exists
	claudeDir := filepath.Dir(settingsPath)
	if _, err := os.Stat(claudeDir); os.IsNotExist(err) {
		return fmt.Errorf("Claude Code configuration directory not found: %s", claudeDir)
	}

	// Validate permissions
	if err := ValidatePermissions(settingsPath); err != nil {
		return fmt.Errorf("insufficient permissions for Claude settings: %w", err)
	}

	return nil
}