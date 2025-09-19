package database

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
)

// EncryptionConfig holds encryption-related configuration
type EncryptionConfig struct {
	Enabled    bool
	KeyPath    string
	KeySize    int
	Iterations int
}

// DefaultEncryptionConfig returns default encryption settings
func DefaultEncryptionConfig() *EncryptionConfig {
	homeDir, _ := os.UserHomeDir()
	keyPath := filepath.Join(homeDir, ".context-extender", "db.key")

	return &EncryptionConfig{
		Enabled:    false, // Disabled by default for compatibility
		KeyPath:    keyPath,
		KeySize:    32, // 256-bit key
		Iterations: 64000, // SQLCipher 4 default
	}
}

// GenerateEncryptionKey generates a new encryption key
func GenerateEncryptionKey(size int) (string, error) {
	key := make([]byte, size)
	if _, err := rand.Read(key); err != nil {
		return "", fmt.Errorf("failed to generate random key: %w", err)
	}
	return hex.EncodeToString(key), nil
}

// SaveEncryptionKey saves the encryption key to a file
func SaveEncryptionKey(key, path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create key directory: %w", err)
	}

	// Write key with restricted permissions
	if err := os.WriteFile(path, []byte(key), 0600); err != nil {
		return fmt.Errorf("failed to write key file: %w", err)
	}

	return nil
}

// LoadEncryptionKey loads the encryption key from a file
func LoadEncryptionKey(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read key file: %w", err)
	}
	return string(data), nil
}

// InitializeEncryptedDatabase initializes an encrypted database
func InitializeEncryptedDatabase(config *Config, encConfig *EncryptionConfig) error {
	if !encConfig.Enabled {
		// Use regular initialization if encryption is disabled
		return Initialize(config)
	}

	// Check if key exists, generate if not
	var key string
	if _, err := os.Stat(encConfig.KeyPath); os.IsNotExist(err) {
		// Generate new key
		newKey, err := GenerateEncryptionKey(encConfig.KeySize)
		if err != nil {
			return fmt.Errorf("failed to generate encryption key: %w", err)
		}

		// Save key
		if err := SaveEncryptionKey(newKey, encConfig.KeyPath); err != nil {
			return fmt.Errorf("failed to save encryption key: %w", err)
		}
		key = newKey
	} else {
		// Load existing key
		existingKey, err := LoadEncryptionKey(encConfig.KeyPath)
		if err != nil {
			return fmt.Errorf("failed to load encryption key: %w", err)
		}
		key = existingKey
	}

	// Modify connection string for SQLCipher
	// Format: file:path?_pragma_key=x'hexkey'&_pragma_cipher_page_size=4096
	encryptedPath := fmt.Sprintf(
		"file:%s?_pragma_key=x'%s'&_pragma_cipher_page_size=4096&_pragma_kdf_iter=%d",
		config.DatabasePath,
		key,
		encConfig.Iterations,
	)

	// Update config with encrypted path
	config.DatabasePath = encryptedPath

	// Initialize with encryption
	return Initialize(config)
}

// SetupEncryptionPragmas configures SQLCipher-specific pragmas
func SetupEncryptionPragmas(key string, iterations int) error {
	db, err := GetConnection()
	if err != nil {
		return err
	}

	// SQLCipher pragmas for encryption
	pragmas := []string{
		fmt.Sprintf("PRAGMA key = \"x'%s'\"", key),
		fmt.Sprintf("PRAGMA kdf_iter = %d", iterations),
		"PRAGMA cipher_page_size = 4096",
		"PRAGMA cipher_use_hmac = ON",
		"PRAGMA cipher_plaintext_header_size = 0",
		"PRAGMA cipher_default_plaintext_header_size = 0",
	}

	for _, pragma := range pragmas {
		if _, err := db.Exec(pragma); err != nil {
			return fmt.Errorf("failed to execute pragma: %w", err)
		}
	}

	// Verify encryption is working
	var result string
	err = db.QueryRow("PRAGMA cipher_version").Scan(&result)
	if err != nil {
		return fmt.Errorf("SQLCipher not available or not working: %w", err)
	}

	return nil
}

// MigrateToEncrypted migrates an unencrypted database to encrypted
func MigrateToEncrypted(sourcePath, destPath, key string) error {
	// This would use SQLCipher's ATTACH and sqlcipher_export()
	// For now, return a placeholder
	return fmt.Errorf("migration to encrypted database not yet implemented")
}

// RotateEncryptionKey rotates the database encryption key
func RotateEncryptionKey(oldKey, newKey string) error {
	db, err := GetConnection()
	if err != nil {
		return err
	}

	// SQLCipher rekey operation
	rekeyPragma := fmt.Sprintf("PRAGMA rekey = \"x'%s'\"", newKey)
	if _, err := db.Exec(rekeyPragma); err != nil {
		return fmt.Errorf("failed to rotate encryption key: %w", err)
	}

	return nil
}

// VerifyDatabaseIntegrity checks database integrity
func VerifyDatabaseIntegrity() error {
	db, err := GetConnection()
	if err != nil {
		return err
	}

	var result string
	err = db.QueryRow("PRAGMA integrity_check").Scan(&result)
	if err != nil {
		return fmt.Errorf("integrity check failed: %w", err)
	}

	if result != "ok" {
		return fmt.Errorf("database integrity check failed: %s", result)
	}

	return nil
}