// +build sqlite3 sqlcipher

package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	_ "github.com/mutecomm/go-sqlcipher/v4"
)

var (
	encryptedDB   *sql.DB
	encryptedOnce sync.Once
	isEncrypted   bool
)

// EncryptedConfig extends Config with encryption settings
type EncryptedConfig struct {
	*Config
	EncryptionEnabled bool
	EncryptionKey     string
	KDFIterations     int
}

// DefaultEncryptedConfig returns default configuration with encryption
func DefaultEncryptedConfig() *EncryptedConfig {
	return &EncryptedConfig{
		Config:            DefaultConfig(),
		EncryptionEnabled: false,
		EncryptionKey:     "",
		KDFIterations:     256000, // SQLCipher 4 recommended
	}
}

// InitializeWithEncryption initializes database with optional encryption
func InitializeWithEncryption(config *EncryptedConfig) error {
	var err error

	encryptedOnce.Do(func() {
		if config == nil {
			config = DefaultEncryptedConfig()
		}

		if config.EncryptionEnabled {
			err = initializeEncryptedDatabase(config)
			isEncrypted = true
		} else {
			// Fall back to regular SQLite
			err = initializeDatabase(config.Config)
			isEncrypted = false
		}
	})

	return err
}

func initializeEncryptedDatabase(config *EncryptedConfig) error {
	var err error

	// Ensure database directory exists
	dbDir := filepath.Dir(config.DatabasePath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return fmt.Errorf("failed to create database directory: %w", err)
	}

	// Open with SQLCipher driver
	encryptedDB, err = sql.Open("sqlite3", config.DatabasePath)
	if err != nil {
		return fmt.Errorf("failed to open encrypted database: %w", err)
	}

	// Set connection pool settings
	encryptedDB.SetMaxOpenConns(config.MaxOpenConns)
	encryptedDB.SetMaxIdleConns(config.MaxIdleConns)

	// Set encryption key
	if config.EncryptionKey != "" {
		keyPragma := fmt.Sprintf("PRAGMA key = '%s'", config.EncryptionKey)
		if _, err := encryptedDB.Exec(keyPragma); err != nil {
			encryptedDB.Close()
			return fmt.Errorf("failed to set encryption key: %w", err)
		}
	}

	// Configure SQLCipher settings
	if err := configureEncryption(config); err != nil {
		encryptedDB.Close()
		return fmt.Errorf("failed to configure encryption: %w", err)
	}

	// Test the connection
	if err := encryptedDB.Ping(); err != nil {
		encryptedDB.Close()
		return fmt.Errorf("failed to ping encrypted database: %w", err)
	}

	// Set regular pragmas
	if err := enableEncryptedPragmas(); err != nil {
		encryptedDB.Close()
		return fmt.Errorf("failed to enable pragmas: %w", err)
	}

	// Update global db reference
	db = encryptedDB

	return nil
}

func configureEncryption(config *EncryptedConfig) error {
	encryptionPragmas := []string{
		fmt.Sprintf("PRAGMA kdf_iter = %d", config.KDFIterations),
		"PRAGMA cipher_page_size = 4096",
		"PRAGMA cipher_hmac_algorithm = HMAC_SHA512",
		"PRAGMA cipher_kdf_algorithm = PBKDF2_HMAC_SHA512",
		"PRAGMA cipher_plaintext_header_size = 0",
		"PRAGMA cipher_use_hmac = ON",
	}

	for _, pragma := range encryptionPragmas {
		if _, err := encryptedDB.Exec(pragma); err != nil {
			return fmt.Errorf("failed to execute encryption pragma %s: %w", pragma, err)
		}
	}

	return nil
}

func enableEncryptedPragmas() error {
	pragmas := []string{
		"PRAGMA foreign_keys = ON",
		"PRAGMA journal_mode = WAL",
		"PRAGMA synchronous = NORMAL",
		"PRAGMA cache_size = 10000",
		"PRAGMA temp_store = memory",
		"PRAGMA mmap_size = 268435456",
	}

	for _, pragma := range pragmas {
		if _, err := encryptedDB.Exec(pragma); err != nil {
			return fmt.Errorf("failed to execute pragma %s: %w", pragma, err)
		}
	}

	return nil
}

// IsEncrypted returns whether the database is encrypted
func IsEncrypted() bool {
	return isEncrypted
}

// VerifyEncryption checks if encryption is properly configured
func VerifyEncryption() error {
	if !isEncrypted {
		return fmt.Errorf("database is not encrypted")
	}

	// Try to query cipher version
	var version string
	err := encryptedDB.QueryRow("SELECT sqlite_version()").Scan(&version)
	if err != nil {
		return fmt.Errorf("failed to verify encryption: %w", err)
	}

	return nil
}

// ChangeEncryptionKey changes the database encryption key
func ChangeEncryptionKey(oldKey, newKey string) error {
	if !isEncrypted {
		return fmt.Errorf("cannot change key on unencrypted database")
	}

	// First verify old key works
	testPragma := fmt.Sprintf("PRAGMA key = '%s'", oldKey)
	if _, err := encryptedDB.Exec(testPragma); err != nil {
		return fmt.Errorf("old key verification failed: %w", err)
	}

	// Rekey the database
	rekeyPragma := fmt.Sprintf("PRAGMA rekey = '%s'", newKey)
	if _, err := encryptedDB.Exec(rekeyPragma); err != nil {
		return fmt.Errorf("failed to change encryption key: %w", err)
	}

	return nil
}

// ExportToUnencrypted exports encrypted database to unencrypted
func ExportToUnencrypted(encryptedPath, unencryptedPath, key string) error {
	// Open encrypted database
	encDB, err := sql.Open("sqlite3", encryptedPath)
	if err != nil {
		return fmt.Errorf("failed to open encrypted database: %w", err)
	}
	defer encDB.Close()

	// Set encryption key
	keyPragma := fmt.Sprintf("PRAGMA key = '%s'", key)
	if _, err := encDB.Exec(keyPragma); err != nil {
		return fmt.Errorf("failed to unlock encrypted database: %w", err)
	}

	// Attach unencrypted database
	attachSQL := fmt.Sprintf("ATTACH DATABASE '%s' AS plaintext KEY ''", unencryptedPath)
	if _, err := encDB.Exec(attachSQL); err != nil {
		return fmt.Errorf("failed to attach unencrypted database: %w", err)
	}

	// Export schema and data
	if _, err := encDB.Exec("SELECT sqlcipher_export('plaintext')"); err != nil {
		return fmt.Errorf("failed to export to unencrypted database: %w", err)
	}

	// Detach
	if _, err := encDB.Exec("DETACH DATABASE plaintext"); err != nil {
		return fmt.Errorf("failed to detach database: %w", err)
	}

	return nil
}

// ImportFromUnencrypted imports unencrypted database to encrypted
func ImportFromUnencrypted(unencryptedPath, encryptedPath, key string) error {
	// Open unencrypted database
	plainDB, err := sql.Open("sqlite3", unencryptedPath)
	if err != nil {
		return fmt.Errorf("failed to open unencrypted database: %w", err)
	}
	defer plainDB.Close()

	// Attach encrypted database with key
	attachSQL := fmt.Sprintf("ATTACH DATABASE '%s' AS encrypted KEY '%s'", encryptedPath, key)
	if _, err := plainDB.Exec(attachSQL); err != nil {
		return fmt.Errorf("failed to attach encrypted database: %w", err)
	}

	// Export to encrypted
	if _, err := plainDB.Exec("SELECT sqlcipher_export('encrypted')"); err != nil {
		return fmt.Errorf("failed to export to encrypted database: %w", err)
	}

	// Detach
	if _, err := plainDB.Exec("DETACH DATABASE encrypted"); err != nil {
		return fmt.Errorf("failed to detach database: %w", err)
	}

	return nil
}