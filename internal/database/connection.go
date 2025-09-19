package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	_ "modernc.org/sqlite"
)

var (
	db   *sql.DB
	once sync.Once
)

type Config struct {
	DatabasePath string
	DriverName   string
	MaxOpenConns int
	MaxIdleConns int
}

func DefaultConfig() *Config {
	return &Config{
		DatabasePath: getDefaultDatabasePath(),
		DriverName:   "sqlite3",
		MaxOpenConns: 25,
		MaxIdleConns: 5,
	}
}

func getDefaultDatabasePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "./context-extender.db"
	}

	dbDir := filepath.Join(homeDir, ".context-extender")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return "./context-extender.db"
	}

	return filepath.Join(dbDir, "conversations.db")
}

func Initialize(config *Config) error {
	var err error

	once.Do(func() {
		if config == nil {
			config = DefaultConfig()
		}

		err = initializeDatabase(config)
	})

	return err
}

func initializeDatabase(config *Config) error {
	var err error

	db, err = sql.Open(config.DriverName, config.DatabasePath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)

	if err = db.Ping(); err != nil {
		db.Close()
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return enablePragmas()
}

func enablePragmas() error {
	pragmas := []string{
		"PRAGMA foreign_keys = ON",
		"PRAGMA journal_mode = WAL",
		"PRAGMA synchronous = NORMAL",
		"PRAGMA cache_size = 10000",
		"PRAGMA temp_store = memory",
		"PRAGMA mmap_size = 268435456",
	}

	for _, pragma := range pragmas {
		if _, err := db.Exec(pragma); err != nil {
			return fmt.Errorf("failed to execute pragma %s: %w", pragma, err)
		}
	}

	return nil
}

func GetConnection() (*sql.DB, error) {
	if db == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	return db, nil
}

func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}