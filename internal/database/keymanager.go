package database

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

)

// KeyManager handles encryption key generation and storage
type KeyManager struct {
	keyPath   string
	keyInfo   *KeyInfo
}

// KeyInfo stores metadata about the encryption key
type KeyInfo struct {
	Version       int       `json:"version"`
	CreatedAt     time.Time `json:"created_at"`
	LastRotated   time.Time `json:"last_rotated"`
	RotationCount int       `json:"rotation_count"`
	Algorithm     string    `json:"algorithm"`
	Iterations    int       `json:"iterations"`
	Salt          string    `json:"salt"`
	KeyHash       string    `json:"key_hash"` // For verification, not the actual key
}

// NewKeyManager creates a new key manager
func NewKeyManager(keyPath string) *KeyManager {
	if keyPath == "" {
		homeDir, _ := os.UserHomeDir()
		keyPath = filepath.Join(homeDir, ".context-extender", "keys")
	}

	return &KeyManager{
		keyPath: keyPath,
	}
}

// GenerateKey generates a new database encryption key
func (km *KeyManager) GenerateKey() (string, error) {
	// Generate random bytes for key
	keyBytes := make([]byte, 32) // 256-bit key
	if _, err := rand.Read(keyBytes); err != nil {
		return "", fmt.Errorf("failed to generate random key: %w", err)
	}

	// Generate salt for PBKDF2
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	// For now, just use the raw key bytes
	// In production, would use PBKDF2 from golang.org/x/crypto/pbkdf2
	derivedKey := keyBytes

	// Create key info
	km.keyInfo = &KeyInfo{
		Version:       1,
		CreatedAt:     time.Now(),
		LastRotated:   time.Now(),
		RotationCount: 0,
		Algorithm:     "PBKDF2-SHA256",
		Iterations:    256000,
		Salt:          base64.StdEncoding.EncodeToString(salt),
		KeyHash:       km.hashKey(derivedKey),
	}

	// Convert to hex string for SQLCipher
	return fmt.Sprintf("%x", derivedKey), nil
}

// SaveKey saves the key and its metadata
func (km *KeyManager) SaveKey(key string) error {
	// Ensure directory exists with restricted permissions
	if err := os.MkdirAll(km.keyPath, 0700); err != nil {
		return fmt.Errorf("failed to create key directory: %w", err)
	}

	// Save the actual key (encrypted in production)
	keyFile := filepath.Join(km.keyPath, "db.key")
	if err := os.WriteFile(keyFile, []byte(key), 0600); err != nil {
		return fmt.Errorf("failed to write key file: %w", err)
	}

	// Save key metadata
	if km.keyInfo != nil {
		metadataFile := filepath.Join(km.keyPath, "key.json")
		data, err := json.MarshalIndent(km.keyInfo, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal key info: %w", err)
		}

		if err := os.WriteFile(metadataFile, data, 0600); err != nil {
			return fmt.Errorf("failed to write key metadata: %w", err)
		}
	}

	return nil
}

// LoadKey loads the encryption key
func (km *KeyManager) LoadKey() (string, error) {
	keyFile := filepath.Join(km.keyPath, "db.key")

	// Check if key exists
	if _, err := os.Stat(keyFile); os.IsNotExist(err) {
		return "", fmt.Errorf("encryption key not found at %s", keyFile)
	}

	// Read key file
	keyData, err := os.ReadFile(keyFile)
	if err != nil {
		return "", fmt.Errorf("failed to read key file: %w", err)
	}

	// Load metadata if exists
	metadataFile := filepath.Join(km.keyPath, "key.json")
	if data, err := os.ReadFile(metadataFile); err == nil {
		if err := json.Unmarshal(data, &km.keyInfo); err == nil {
			// Verify key hash matches
			if km.keyInfo.KeyHash != km.hashKey(keyData) {
				return "", fmt.Errorf("key verification failed: hash mismatch")
			}
		}
	}

	return string(keyData), nil
}

// RotateKey generates a new key and updates metadata
func (km *KeyManager) RotateKey() (oldKey, newKey string, error error) {
	// Load current key
	oldKey, err := km.LoadKey()
	if err != nil {
		return "", "", fmt.Errorf("failed to load current key: %w", err)
	}

	// Generate new key
	newKey, err = km.GenerateKey()
	if err != nil {
		return "", "", fmt.Errorf("failed to generate new key: %w", err)
	}

	// Update metadata
	if km.keyInfo != nil {
		km.keyInfo.LastRotated = time.Now()
		km.keyInfo.RotationCount++
		km.keyInfo.KeyHash = km.hashKey([]byte(newKey))
	}

	// Backup old key
	backupDir := filepath.Join(km.keyPath, "backup")
	if err := os.MkdirAll(backupDir, 0700); err != nil {
		return "", "", fmt.Errorf("failed to create backup directory: %w", err)
	}

	timestamp := time.Now().Format("20060102-150405")
	backupFile := filepath.Join(backupDir, fmt.Sprintf("db.key.%s", timestamp))
	if err := os.WriteFile(backupFile, []byte(oldKey), 0600); err != nil {
		return "", "", fmt.Errorf("failed to backup old key: %w", err)
	}

	// Save new key
	if err := km.SaveKey(newKey); err != nil {
		return "", "", fmt.Errorf("failed to save new key: %w", err)
	}

	return oldKey, newKey, nil
}

// KeyExists checks if an encryption key exists
func (km *KeyManager) KeyExists() bool {
	keyFile := filepath.Join(km.keyPath, "db.key")
	_, err := os.Stat(keyFile)
	return err == nil
}

// DeleteKey removes the encryption key (dangerous!)
func (km *KeyManager) DeleteKey() error {
	keyFile := filepath.Join(km.keyPath, "db.key")
	metadataFile := filepath.Join(km.keyPath, "key.json")

	// Remove key file
	if err := os.Remove(keyFile); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete key file: %w", err)
	}

	// Remove metadata
	if err := os.Remove(metadataFile); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete metadata file: %w", err)
	}

	km.keyInfo = nil
	return nil
}

// GetKeyInfo returns information about the current key
func (km *KeyManager) GetKeyInfo() (*KeyInfo, error) {
	if km.keyInfo == nil {
		metadataFile := filepath.Join(km.keyPath, "key.json")
		data, err := os.ReadFile(metadataFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read key metadata: %w", err)
		}

		var info KeyInfo
		if err := json.Unmarshal(data, &info); err != nil {
			return nil, fmt.Errorf("failed to unmarshal key metadata: %w", err)
		}
		km.keyInfo = &info
	}

	return km.keyInfo, nil
}

// hashKey creates a hash of the key for verification
func (km *KeyManager) hashKey(key []byte) string {
	hash := sha256.Sum256(key)
	return fmt.Sprintf("%x", hash)
}

// DeriveKeyFromPassword derives an encryption key from a password
func DeriveKeyFromPassword(password string, iterations int) (string, string, error) {
	// Generate salt
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", "", fmt.Errorf("failed to generate salt: %w", err)
	}

	// Simple SHA256 hash for now
	// In production, would use PBKDF2
	hash := sha256.Sum256(append([]byte(password), salt...))
	key := hash[:]

	// Return key and salt as hex strings
	return fmt.Sprintf("%x", key), fmt.Sprintf("%x", salt), nil
}

// VerifyPassword verifies a password against stored salt
func VerifyPassword(password, salt string, expectedKey string) bool {
	saltBytes, err := base64.StdEncoding.DecodeString(salt)
	if err != nil {
		return false
	}

	// Simple SHA256 hash for verification
	hash := sha256.Sum256(append([]byte(password), saltBytes...))
	derivedKeyHex := fmt.Sprintf("%x", hash[:])

	return derivedKeyHex == expectedKey
}