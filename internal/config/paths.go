package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// GetClaudeSettingsPath returns the path to Claude Code's settings.json file
func GetClaudeSettingsPath() (string, error) {
	var baseDir string

	switch runtime.GOOS {
	case "windows":
		baseDir = os.Getenv("APPDATA")
		if baseDir == "" {
			return "", fmt.Errorf("APPDATA environment variable not set")
		}
	case "darwin":
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get user home directory: %w", err)
		}
		baseDir = filepath.Join(homeDir, ".config")
	default: // linux and others
		baseDir = os.Getenv("XDG_CONFIG_HOME")
		if baseDir == "" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return "", fmt.Errorf("failed to get user home directory: %w", err)
			}
			baseDir = filepath.Join(homeDir, ".config")
		}
	}

	return filepath.Join(baseDir, "claude", "settings.json"), nil
}

// GetContextExtenderConfigPath returns the path to context-extender's config directory
func GetContextExtenderConfigPath() (string, error) {
	var baseDir string

	switch runtime.GOOS {
	case "windows":
		baseDir = os.Getenv("APPDATA")
		if baseDir == "" {
			return "", fmt.Errorf("APPDATA environment variable not set")
		}
	case "darwin":
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get user home directory: %w", err)
		}
		baseDir = filepath.Join(homeDir, ".context-extender")
	default: // linux and others
		configDir := os.Getenv("XDG_CONFIG_HOME")
		if configDir == "" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return "", fmt.Errorf("failed to get user home directory: %w", err)
			}
			configDir = filepath.Join(homeDir, ".config")
		}
		baseDir = filepath.Join(configDir, "context-extender")
	}

	return baseDir, nil
}

// GetContextExtenderStoragePath returns the path to conversation storage directory
func GetContextExtenderStoragePath() (string, error) {
	configPath, err := GetContextExtenderConfigPath()
	if err != nil {
		return "", err
	}

	return filepath.Join(configPath, "conversations"), nil
}

// EnsureDirectoryExists creates the directory if it doesn't exist
func EnsureDirectoryExists(path string) error {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return os.MkdirAll(path, 0755)
	}
	if err != nil {
		return fmt.Errorf("failed to check directory: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("path exists but is not a directory: %s", path)
	}
	return nil
}

// ValidatePermissions checks if we can read and write to the given path
func ValidatePermissions(path string) error {
	// Check if path exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Try to create parent directory to test write permissions
		parentDir := filepath.Dir(path)
		if err := EnsureDirectoryExists(parentDir); err != nil {
			return fmt.Errorf("cannot create parent directory %s: %w", parentDir, err)
		}
	}

	// Try to create a temporary file to test permissions
	tempFile := path + ".tmp"
	file, err := os.Create(tempFile)
	if err != nil {
		return fmt.Errorf("insufficient permissions to write to %s: %w", path, err)
	}
	file.Close()

	// Clean up temp file
	if err := os.Remove(tempFile); err != nil {
		// Log but don't fail - the main operation succeeded
		fmt.Printf("Warning: failed to clean up temp file %s: %v\n", tempFile, err)
	}

	return nil
}