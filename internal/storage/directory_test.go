package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// TestNewStorageManager tests storage manager creation
func TestNewStorageManager(t *testing.T) {
	tests := []struct {
		name     string
		config   *StorageConfig
		expected string
		wantErr  bool
	}{
		{
			name:     "default storage path",
			config:   nil,
			expected: "", // Will be set based on platform
			wantErr:  false,
		},
		{
			name:     "custom storage path",
			config:   &StorageConfig{CustomPath: "/tmp/test-context-extender"},
			expected: "/tmp/test-context-extender",
			wantErr:  false,
		},
		{
			name:     "empty custom path uses default",
			config:   &StorageConfig{CustomPath: ""},
			expected: "", // Will be set based on platform
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm, err := NewStorageManager(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewStorageManager() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if tt.expected != "" {
					// For custom paths, check exact match
					absExpected, _ := filepath.Abs(tt.expected)
					if sm.GetBaseDir() != absExpected {
						t.Errorf("NewStorageManager() baseDir = %v, expected %v", sm.GetBaseDir(), absExpected)
					}
				} else {
					// For default paths, check platform-specific expectations
					baseDir := sm.GetBaseDir()
					switch runtime.GOOS {
					case "windows":
						if !strings.Contains(baseDir, "context-extender") {
							t.Errorf("Windows path should contain 'context-extender', got: %v", baseDir)
						}
					case "linux":
						if !strings.Contains(baseDir, ".context-extender") && !strings.Contains(baseDir, "context-extender") {
							t.Errorf("Linux path should contain 'context-extender', got: %v", baseDir)
						}
					case "darwin":
						if !strings.Contains(baseDir, "context-extender") {
							t.Errorf("macOS path should contain 'context-extender', got: %v", baseDir)
						}
					}
				}
			}
		})
	}
}

// TestGetDefaultStorageDir tests platform-specific storage directory resolution
func TestGetDefaultStorageDir(t *testing.T) {
	dir, err := getDefaultStorageDir()
	if err != nil {
		t.Fatalf("getDefaultStorageDir() failed: %v", err)
	}

	if dir == "" {
		t.Error("getDefaultStorageDir() returned empty path")
	}

	// Verify platform-specific behavior
	switch runtime.GOOS {
	case "windows":
		if !strings.Contains(dir, "context-extender") {
			t.Errorf("Windows storage dir should contain 'context-extender', got: %v", dir)
		}
	case "linux":
		// Should be ~/.context-extender or $XDG_CONFIG_HOME/context-extender
		if !strings.Contains(dir, "context-extender") {
			t.Errorf("Linux storage dir should contain 'context-extender', got: %v", dir)
		}
	case "darwin":
		// Should be ~/Library/Application Support/context-extender
		if !strings.Contains(dir, "context-extender") {
			t.Errorf("macOS storage dir should contain 'context-extender', got: %v", dir)
		}
	}
}

// TestStorageManagerDirectories tests directory path generation
func TestStorageManagerDirectories(t *testing.T) {
	tempDir := t.TempDir()
	config := &StorageConfig{CustomPath: tempDir}

	sm, err := NewStorageManager(config)
	if err != nil {
		t.Fatalf("NewStorageManager() failed: %v", err)
	}

	// Test all directory getters
	tests := []struct {
		name     string
		getter   func() string
		expected string
	}{
		{
			name:     "base directory",
			getter:   sm.GetBaseDir,
			expected: tempDir,
		},
		{
			name:     "conversations directory",
			getter:   sm.GetConversationsDir,
			expected: filepath.Join(tempDir, "conversations"),
		},
		{
			name:     "config directory",
			getter:   sm.GetConfigDir,
			expected: filepath.Join(tempDir, "config"),
		},
		{
			name:     "logs directory",
			getter:   sm.GetLogsDir,
			expected: filepath.Join(tempDir, "logs"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			absExpected, _ := filepath.Abs(tt.expected)
			result := tt.getter()
			if result != absExpected {
				t.Errorf("%s = %v, expected %v", tt.name, result, absExpected)
			}
		})
	}
}

// TestEnsureStorageStructure tests directory creation
func TestEnsureStorageStructure(t *testing.T) {
	tempDir := t.TempDir()
	config := &StorageConfig{CustomPath: tempDir}

	sm, err := NewStorageManager(config)
	if err != nil {
		t.Fatalf("NewStorageManager() failed: %v", err)
	}

	// Ensure storage structure
	err = sm.EnsureStorageStructure()
	if err != nil {
		t.Fatalf("EnsureStorageStructure() failed: %v", err)
	}

	// Verify all directories were created
	dirs := []string{
		sm.GetBaseDir(),
		sm.GetConversationsDir(),
		sm.GetConfigDir(),
		sm.GetLogsDir(),
	}

	for _, dir := range dirs {
		if info, err := os.Stat(dir); err != nil {
			t.Errorf("Directory %s was not created: %v", dir, err)
		} else if !info.IsDir() {
			t.Errorf("Path %s exists but is not a directory", dir)
		}
	}
}

// TestEnsureStorageStructureExistingFile tests error handling when file exists at directory path
func TestEnsureStorageStructureExistingFile(t *testing.T) {
	tempDir := t.TempDir()

	// Create a file where we want to create a directory
	filePath := filepath.Join(tempDir, "conversations")
	file, err := os.Create(filePath)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	file.Close()

	config := &StorageConfig{CustomPath: tempDir}
	sm, err := NewStorageManager(config)
	if err != nil {
		t.Fatalf("NewStorageManager() failed: %v", err)
	}

	// This should fail because a file exists where we want to create a directory
	err = sm.EnsureStorageStructure()
	if err == nil {
		t.Error("EnsureStorageStructure() should have failed when file exists at directory path")
	}

	if !strings.Contains(err.Error(), "not a directory") {
		t.Errorf("Expected 'not a directory' error, got: %v", err)
	}
}

// TestValidateStorageAccess tests storage access validation
func TestValidateStorageAccess(t *testing.T) {
	tempDir := t.TempDir()
	config := &StorageConfig{CustomPath: tempDir}

	sm, err := NewStorageManager(config)
	if err != nil {
		t.Fatalf("NewStorageManager() failed: %v", err)
	}

	// Validate storage access (should create directories and test write access)
	err = sm.ValidateStorageAccess()
	if err != nil {
		t.Fatalf("ValidateStorageAccess() failed: %v", err)
	}

	// Verify directories exist and are writable
	dirs := []string{
		sm.GetBaseDir(),
		sm.GetConversationsDir(),
		sm.GetConfigDir(),
		sm.GetLogsDir(),
	}

	for _, dir := range dirs {
		// Test write access by creating a test file
		testFile := filepath.Join(dir, "write_test.txt")
		file, err := os.Create(testFile)
		if err != nil {
			t.Errorf("Directory %s is not writable: %v", dir, err)
			continue
		}
		file.Close()
		os.Remove(testFile)
	}
}

// TestGetStorageInfo tests storage information retrieval
func TestGetStorageInfo(t *testing.T) {
	tempDir := t.TempDir()
	config := &StorageConfig{CustomPath: tempDir}

	sm, err := NewStorageManager(config)
	if err != nil {
		t.Fatalf("NewStorageManager() failed: %v", err)
	}

	info := sm.GetStorageInfo()

	// Verify storage info
	if info == nil {
		t.Fatal("GetStorageInfo() returned nil")
	}

	absExpected, _ := filepath.Abs(tempDir)
	if info.BaseDir != absExpected {
		t.Errorf("StorageInfo.BaseDir = %v, expected %v", info.BaseDir, absExpected)
	}

	if info.Platform != runtime.GOOS {
		t.Errorf("StorageInfo.Platform = %v, expected %v", info.Platform, runtime.GOOS)
	}

	// Verify all directory paths are set
	if info.ConversationsDir == "" {
		t.Error("StorageInfo.ConversationsDir is empty")
	}
	if info.ConfigDir == "" {
		t.Error("StorageInfo.ConfigDir is empty")
	}
	if info.LogsDir == "" {
		t.Error("StorageInfo.LogsDir is empty")
	}
}

// TestGetStorageUsage tests storage usage calculation
func TestGetStorageUsage(t *testing.T) {
	tempDir := t.TempDir()
	config := &StorageConfig{CustomPath: tempDir}

	sm, err := NewStorageManager(config)
	if err != nil {
		t.Fatalf("NewStorageManager() failed: %v", err)
	}

	// Create storage structure
	err = sm.EnsureStorageStructure()
	if err != nil {
		t.Fatalf("EnsureStorageStructure() failed: %v", err)
	}

	// Create some test files
	testFiles := []struct {
		path    string
		content string
	}{
		{filepath.Join(sm.GetConversationsDir(), "test1.json"), "test content 1"},
		{filepath.Join(sm.GetConfigDir(), "test2.json"), "test content 2"},
		{filepath.Join(sm.GetLogsDir(), "test3.log"), "test content 3"},
	}

	for _, tf := range testFiles {
		err = os.WriteFile(tf.path, []byte(tf.content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", tf.path, err)
		}
	}

	// Get storage usage
	usage, err := sm.GetStorageUsage()
	if err != nil {
		t.Fatalf("GetStorageUsage() failed: %v", err)
	}

	// Verify usage calculation
	if usage.ConversationsSize == 0 {
		t.Error("ConversationsSize should be > 0")
	}
	if usage.ConfigSize == 0 {
		t.Error("ConfigSize should be > 0")
	}
	if usage.LogsSize == 0 {
		t.Error("LogsSize should be > 0")
	}

	expectedTotal := usage.ConversationsSize + usage.ConfigSize + usage.LogsSize
	if usage.TotalSize != expectedTotal {
		t.Errorf("TotalSize = %v, expected %v", usage.TotalSize, expectedTotal)
	}
}

// TestCleanupStorage tests temporary file cleanup
func TestCleanupStorage(t *testing.T) {
	tempDir := t.TempDir()
	config := &StorageConfig{CustomPath: tempDir}

	sm, err := NewStorageManager(config)
	if err != nil {
		t.Fatalf("NewStorageManager() failed: %v", err)
	}

	// Create storage structure
	err = sm.EnsureStorageStructure()
	if err != nil {
		t.Fatalf("EnsureStorageStructure() failed: %v", err)
	}

	// Create test files including temporary ones
	testFiles := []struct {
		path      string
		content   string
		shouldRemove bool
	}{
		{filepath.Join(sm.GetConversationsDir(), "normal.json"), "normal file", false},
		{filepath.Join(sm.GetConversationsDir(), ".temp_file"), "temp file", true},
		{filepath.Join(sm.GetConversationsDir(), "backup.tmp"), "temp backup", true},
		{filepath.Join(sm.GetConfigDir(), ".gitkeep"), "gitkeep file", false}, // Should not remove .gitkeep
		{filepath.Join(sm.GetLogsDir(), "app.log"), "normal log", false},
	}

	for _, tf := range testFiles {
		err = os.WriteFile(tf.path, []byte(tf.content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", tf.path, err)
		}
	}

	// Cleanup storage
	err = sm.CleanupStorage()
	if err != nil {
		t.Fatalf("CleanupStorage() failed: %v", err)
	}

	// Verify cleanup results
	for _, tf := range testFiles {
		_, err := os.Stat(tf.path)
		if tf.shouldRemove {
			if !os.IsNotExist(err) {
				t.Errorf("Temporary file %s should have been removed", tf.path)
			}
		} else {
			if os.IsNotExist(err) {
				t.Errorf("Normal file %s should not have been removed", tf.path)
			}
		}
	}
}

// TestValidateDirectoryPermissions tests directory permission validation
func TestValidateDirectoryPermissions(t *testing.T) {
	tempDir := t.TempDir()
	config := &StorageConfig{CustomPath: tempDir}

	sm, err := NewStorageManager(config)
	if err != nil {
		t.Fatalf("NewStorageManager() failed: %v", err)
	}

	// Create a test directory
	testDir := filepath.Join(tempDir, "test_permissions")
	err = os.MkdirAll(testDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Test permission validation
	err = sm.validateDirectoryPermissions(testDir)
	if err != nil {
		t.Errorf("validateDirectoryPermissions() failed for writable directory: %v", err)
	}
}

// TestCalculateDirectorySize tests directory size calculation
func TestCalculateDirectorySize(t *testing.T) {
	tempDir := t.TempDir()
	config := &StorageConfig{CustomPath: tempDir}

	sm, err := NewStorageManager(config)
	if err != nil {
		t.Fatalf("NewStorageManager() failed: %v", err)
	}

	// Create test directory with files
	testDir := filepath.Join(tempDir, "size_test")
	err = os.MkdirAll(testDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Create files with known sizes
	testContent1 := "Hello, World!" // 13 bytes
	testContent2 := "Test content for size calculation" // 33 bytes

	err = os.WriteFile(filepath.Join(testDir, "file1.txt"), []byte(testContent1), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file 1: %v", err)
	}

	err = os.WriteFile(filepath.Join(testDir, "file2.txt"), []byte(testContent2), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file 2: %v", err)
	}

	// Calculate directory size
	size, err := sm.calculateDirectorySize(testDir)
	if err != nil {
		t.Fatalf("calculateDirectorySize() failed: %v", err)
	}

	expectedSize := int64(len(testContent1) + len(testContent2))
	if size != expectedSize {
		t.Errorf("calculateDirectorySize() = %v, expected %v", size, expectedSize)
	}
}

// Benchmark tests for performance validation

func BenchmarkNewStorageManager(b *testing.B) {
	config := &StorageConfig{CustomPath: "/tmp/benchmark-context-extender"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := NewStorageManager(config)
		if err != nil {
			b.Fatalf("NewStorageManager() failed: %v", err)
		}
	}
}

func BenchmarkEnsureStorageStructure(b *testing.B) {
	tempDir := b.TempDir()
	config := &StorageConfig{CustomPath: tempDir}

	sm, err := NewStorageManager(config)
	if err != nil {
		b.Fatalf("NewStorageManager() failed: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := sm.EnsureStorageStructure()
		if err != nil {
			b.Fatalf("EnsureStorageStructure() failed: %v", err)
		}
	}
}

func BenchmarkGetStorageUsage(b *testing.B) {
	tempDir := b.TempDir()
	config := &StorageConfig{CustomPath: tempDir}

	sm, err := NewStorageManager(config)
	if err != nil {
		b.Fatalf("NewStorageManager() failed: %v", err)
	}

	// Create storage structure and some test files
	err = sm.EnsureStorageStructure()
	if err != nil {
		b.Fatalf("EnsureStorageStructure() failed: %v", err)
	}

	// Create some test content
	for i := 0; i < 10; i++ {
		content := strings.Repeat("test content ", 100)
		err = os.WriteFile(filepath.Join(sm.GetConversationsDir(), fmt.Sprintf("test%d.json", i)), []byte(content), 0644)
		if err != nil {
			b.Fatalf("Failed to create test file: %v", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := sm.GetStorageUsage()
		if err != nil {
			b.Fatalf("GetStorageUsage() failed: %v", err)
		}
	}
}