package database

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config == nil {
		t.Fatal("DefaultConfig() returned nil")
	}

	if config.DriverName != "sqlite3" {
		t.Errorf("Expected driver name 'sqlite3', got '%s'", config.DriverName)
	}

	if config.MaxOpenConns != 25 {
		t.Errorf("Expected MaxOpenConns 25, got %d", config.MaxOpenConns)
	}

	if config.MaxIdleConns != 5 {
		t.Errorf("Expected MaxIdleConns 5, got %d", config.MaxIdleConns)
	}

	if config.DatabasePath == "" {
		t.Error("Expected non-empty DatabasePath")
	}
}

func TestInitializeDatabase(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "context-extender-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	config := &Config{
		DatabasePath: filepath.Join(tempDir, "test.db"),
		DriverName:   "sqlite3",
		MaxOpenConns: 10,
		MaxIdleConns: 2,
	}

	err = Initialize(config)
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	conn, err := GetConnection()
	if err != nil {
		t.Fatalf("Failed to get connection: %v", err)
	}

	if conn == nil {
		t.Fatal("Connection is nil")
	}

	err = conn.Ping()
	if err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}

	err = Close()
	if err != nil {
		t.Fatalf("Failed to close database: %v", err)
	}
}

func TestGetConnectionWithoutInitialize(t *testing.T) {
	db = nil

	_, err := GetConnection()
	if err == nil {
		t.Error("Expected error when getting connection without initialization")
	}
}