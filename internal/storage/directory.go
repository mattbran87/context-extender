package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// StorageManager handles storage directory operations
type StorageManager struct {
	baseDir string
}

// StorageConfig holds configuration for storage operations
type StorageConfig struct {
	CustomPath string
}

// NewStorageManager creates a new storage manager with default or custom path
func NewStorageManager(config *StorageConfig) (*StorageManager, error) {
	var baseDir string
	var err error

	if config != nil && config.CustomPath != "" {
		// Use custom path if provided
		baseDir = config.CustomPath
	} else {
		// Use platform-specific default path
		baseDir, err = getDefaultStorageDir()
		if err != nil {
			return nil, fmt.Errorf("failed to determine default storage directory: %w", err)
		}
	}

	// Convert to absolute path
	baseDir, err = filepath.Abs(baseDir)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	return &StorageManager{
		baseDir: baseDir,
	}, nil
}

// getDefaultStorageDir returns the platform-specific default storage directory
func getDefaultStorageDir() (string, error) {
	switch runtime.GOOS {
	case "windows":
		// Use %APPDATA%\context-extender\
		appData := os.Getenv("APPDATA")
		if appData == "" {
			userConfigDir, err := os.UserConfigDir()
			if err != nil {
				return "", fmt.Errorf("failed to get user config directory: %w", err)
			}
			return filepath.Join(userConfigDir, "context-extender"), nil
		}
		return filepath.Join(appData, "context-extender"), nil

	case "darwin":
		// Use ~/Library/Application Support/context-extender/
		userConfigDir, err := os.UserConfigDir()
		if err != nil {
			return "", fmt.Errorf("failed to get user config directory: %w", err)
		}
		return filepath.Join(userConfigDir, "context-extender"), nil

	case "linux":
		// Use ~/.context-extender/ or $XDG_CONFIG_HOME/context-extender/
		configDir := os.Getenv("XDG_CONFIG_HOME")
		if configDir == "" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return "", fmt.Errorf("failed to get user home directory: %w", err)
			}
			return filepath.Join(homeDir, ".context-extender"), nil
		}
		return filepath.Join(configDir, "context-extender"), nil

	default:
		// Fallback for other Unix-like systems
		userConfigDir, err := os.UserConfigDir()
		if err != nil {
			// Final fallback to home directory
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return "", fmt.Errorf("failed to get user home directory: %w", err)
			}
			return filepath.Join(homeDir, ".context-extender"), nil
		}
		return filepath.Join(userConfigDir, "context-extender"), nil
	}
}

// GetBaseDir returns the base storage directory path
func (sm *StorageManager) GetBaseDir() string {
	return sm.baseDir
}

// GetConversationsDir returns the conversations storage directory path
func (sm *StorageManager) GetConversationsDir() string {
	return filepath.Join(sm.baseDir, "conversations")
}

// GetConfigDir returns the configuration directory path
func (sm *StorageManager) GetConfigDir() string {
	return filepath.Join(sm.baseDir, "config")
}

// GetLogsDir returns the logs directory path
func (sm *StorageManager) GetLogsDir() string {
	return filepath.Join(sm.baseDir, "logs")
}

// EnsureStorageStructure creates the storage directory structure if it doesn't exist
func (sm *StorageManager) EnsureStorageStructure() error {
	dirs := []string{
		sm.GetBaseDir(),
		sm.GetConversationsDir(),
		sm.GetConfigDir(),
		sm.GetLogsDir(),
	}

	for _, dir := range dirs {
		if err := sm.ensureDirectory(dir); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// ensureDirectory creates a directory with proper permissions if it doesn't exist
func (sm *StorageManager) ensureDirectory(path string) error {
	// Check if directory already exists
	if info, err := os.Stat(path); err == nil {
		if !info.IsDir() {
			return fmt.Errorf("path exists but is not a directory: %s", path)
		}
		// Directory exists, check permissions
		return sm.validateDirectoryPermissions(path)
	}

	// Create directory with proper permissions
	// 0755 = rwxr-xr-x (owner: read/write/execute, group/others: read/execute)
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	return nil
}

// validateDirectoryPermissions checks if the directory has proper permissions
func (sm *StorageManager) validateDirectoryPermissions(path string) error {
	// Test write permissions by creating a temporary file
	testFile := filepath.Join(path, ".write_test")
	file, err := os.Create(testFile)
	if err != nil {
		return fmt.Errorf("directory is not writable: %w", err)
	}
	file.Close()

	// Clean up test file
	if err := os.Remove(testFile); err != nil {
		// Log warning but don't fail - the main operation succeeded
		fmt.Printf("Warning: failed to clean up test file %s: %v\n", testFile, err)
	}

	return nil
}

// ValidateStorageAccess checks if the storage directory is accessible and writable
func (sm *StorageManager) ValidateStorageAccess() error {
	// Ensure directory structure exists
	if err := sm.EnsureStorageStructure(); err != nil {
		return fmt.Errorf("failed to ensure storage structure: %w", err)
	}

	// Test access to each directory
	dirs := []string{
		sm.GetBaseDir(),
		sm.GetConversationsDir(),
		sm.GetConfigDir(),
		sm.GetLogsDir(),
	}

	for _, dir := range dirs {
		if err := sm.validateDirectoryPermissions(dir); err != nil {
			return fmt.Errorf("storage validation failed for %s: %w", dir, err)
		}
	}

	return nil
}

// GetStorageInfo returns information about the storage setup
func (sm *StorageManager) GetStorageInfo() *StorageInfo {
	return &StorageInfo{
		BaseDir:          sm.GetBaseDir(),
		ConversationsDir: sm.GetConversationsDir(),
		ConfigDir:        sm.GetConfigDir(),
		LogsDir:          sm.GetLogsDir(),
		Platform:         runtime.GOOS,
	}
}

// StorageInfo contains information about storage directories
type StorageInfo struct {
	BaseDir          string `json:"base_dir"`
	ConversationsDir string `json:"conversations_dir"`
	ConfigDir        string `json:"config_dir"`
	LogsDir          string `json:"logs_dir"`
	Platform         string `json:"platform"`
}

// GetStorageUsage returns storage usage statistics (basic implementation)
func (sm *StorageManager) GetStorageUsage() (*StorageUsage, error) {
	usage := &StorageUsage{}

	// Calculate directory sizes
	if size, err := sm.calculateDirectorySize(sm.GetConversationsDir()); err == nil {
		usage.ConversationsSize = size
	}

	if size, err := sm.calculateDirectorySize(sm.GetConfigDir()); err == nil {
		usage.ConfigSize = size
	}

	if size, err := sm.calculateDirectorySize(sm.GetLogsDir()); err == nil {
		usage.LogsSize = size
	}

	usage.TotalSize = usage.ConversationsSize + usage.ConfigSize + usage.LogsSize

	return usage, nil
}

// StorageUsage contains storage usage statistics
type StorageUsage struct {
	ConversationsSize int64 `json:"conversations_size"`
	ConfigSize        int64 `json:"config_size"`
	LogsSize          int64 `json:"logs_size"`
	TotalSize         int64 `json:"total_size"`
}

// calculateDirectorySize calculates the total size of a directory
func (sm *StorageManager) calculateDirectorySize(dirPath string) (int64, error) {
	var size int64

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// Skip files that can't be accessed
			return nil
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})

	if err != nil {
		return 0, err
	}

	return size, nil
}

// CleanupStorage removes temporary files and performs basic cleanup
func (sm *StorageManager) CleanupStorage() error {
	// Remove any temporary files (files starting with .)
	dirs := []string{
		sm.GetBaseDir(),
		sm.GetConversationsDir(),
		sm.GetConfigDir(),
		sm.GetLogsDir(),
	}

	for _, dir := range dirs {
		if err := sm.cleanupTempFiles(dir); err != nil {
			return fmt.Errorf("failed to cleanup temp files in %s: %w", dir, err)
		}
	}

	return nil
}

// cleanupTempFiles removes temporary files from a directory
func (sm *StorageManager) cleanupTempFiles(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		// Directory doesn't exist, nothing to clean
		return nil
	}

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// Remove temporary files (starting with . or ending with .tmp)
		name := entry.Name()
		if name[0] == '.' && name != ".gitkeep" ||
		   len(name) > 4 && name[len(name)-4:] == ".tmp" {
			filePath := filepath.Join(dirPath, name)
			if err := os.Remove(filePath); err != nil {
				// Log warning but don't fail
				fmt.Printf("Warning: failed to remove temp file %s: %v\n", filePath, err)
			}
		}
	}

	return nil
}