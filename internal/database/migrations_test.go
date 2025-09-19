package database

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunMigrations(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "context-extender-migrations-test-*")
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
	defer Close()

	err = RunMigrations()
	if err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	db, err := GetConnection()
	if err != nil {
		t.Fatalf("Failed to get connection: %v", err)
	}

	tables := []string{"sessions", "events", "conversations", "import_history", "settings", "schema_migrations"}

	for _, table := range tables {
		query := "SELECT name FROM sqlite_master WHERE type='table' AND name=?"
		var name string
		err = db.QueryRow(query, table).Scan(&name)
		if err != nil {
			t.Errorf("Table %s was not created: %v", table, err)
		}
	}

	var migrationCount int
	err = db.QueryRow("SELECT COUNT(*) FROM schema_migrations").Scan(&migrationCount)
	if err != nil {
		t.Fatalf("Failed to count migrations: %v", err)
	}

	if migrationCount != len(migrations) {
		t.Errorf("Expected %d migrations, got %d", len(migrations), migrationCount)
	}
}

func TestRunMigrationsIdempotent(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "context-extender-migrations-idempotent-test-*")
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
	defer Close()

	err = RunMigrations()
	if err != nil {
		t.Fatalf("Failed to run migrations first time: %v", err)
	}

	err = RunMigrations()
	if err != nil {
		t.Fatalf("Failed to run migrations second time: %v", err)
	}

	db, err := GetConnection()
	if err != nil {
		t.Fatalf("Failed to get connection: %v", err)
	}

	var migrationCount int
	err = db.QueryRow("SELECT COUNT(*) FROM schema_migrations").Scan(&migrationCount)
	if err != nil {
		t.Fatalf("Failed to count migrations: %v", err)
	}

	if migrationCount != len(migrations) {
		t.Errorf("Expected %d migrations, got %d after running twice", len(migrations), migrationCount)
	}
}