# Context Extender API Reference

## Encryption Package

### Core Interfaces

#### EncryptionProvider

The main interface for encryption operations:

```go
type EncryptionProvider interface {
    // Initialize prepares the provider for use
    Initialize() error

    // Close releases resources
    Close() error

    // EncryptField encrypts a single field
    EncryptField(data []byte) (*FieldEncryption, error)

    // DecryptField decrypts a single field
    DecryptField(encrypted *FieldEncryption) ([]byte, error)

    // BulkEncrypt encrypts multiple fields
    BulkEncrypt(data [][]byte) ([]*FieldEncryption, error)

    // BulkDecrypt decrypts multiple fields
    BulkDecrypt(encrypted []*FieldEncryption) ([][]byte, error)

    // RotateKey generates a new encryption key
    RotateKey() error

    // GetCurrentKeyID returns the current key identifier
    GetCurrentKeyID() string

    // GetMetrics returns performance metrics
    GetMetrics() *EncryptionMetrics

    // IsInitialized checks if provider is ready
    IsInitialized() bool
}
```

#### FieldMiddleware

Interface for field-level encryption middleware:

```go
type FieldMiddleware struct {
    provider EncryptionProvider
    rules    *FieldRules
}

// Methods
func (fm *FieldMiddleware) EncryptRow(tableName string, row map[string]interface{}) (map[string]interface{}, error)
func (fm *FieldMiddleware) DecryptRow(tableName string, row map[string]interface{}) (map[string]interface{}, error)
func (fm *FieldMiddleware) EncryptStruct(v interface{}) error
func (fm *FieldMiddleware) DecryptStruct(v interface{}) error
```

### Data Types

#### EncryptionConfig

```go
type EncryptionConfig struct {
    Method           EncryptionMethod `json:"method"`
    KeyPath          string           `json:"key_path"`
    KeyDerivation    string           `json:"key_derivation"`
    PBKDF2Iterations int              `json:"pbkdf2_iterations"`
    Salt             []byte           `json:"-"`
    Enabled          bool             `json:"enabled"`
}
```

#### FieldEncryption

```go
type FieldEncryption struct {
    EncryptedData []byte `json:"encrypted_data"`
    Nonce         []byte `json:"nonce"`
    KeyID         string `json:"key_id"`
    Method        string `json:"method"`
    Timestamp     int64  `json:"timestamp"`
}
```

#### EncryptionMetrics

```go
type EncryptionMetrics struct {
    EncryptionOps   uint64        `json:"encryption_ops"`
    DecryptionOps   uint64        `json:"decryption_ops"`
    Errors          uint64        `json:"errors"`
    TotalTime       time.Duration `json:"total_time"`
    AverageTime     time.Duration `json:"average_time"`
    KeyRotations    uint64        `json:"key_rotations"`
    LastRotation    time.Time     `json:"last_rotation"`
}
```

### Factory Functions

#### NewAESGCMProvider

```go
func NewAESGCMProvider(config *EncryptionConfig) (*AESGCMProvider, error)
```

Creates a new AES-GCM encryption provider.

**Parameters:**
- `config`: Encryption configuration

**Returns:**
- Provider instance or error

**Example:**
```go
config := &EncryptionConfig{
    Method:           EncryptionMethodAESGCM,
    KeyPath:          "/path/to/key.key",
    PBKDF2Iterations: 100000,
    Enabled:          true,
}
provider, err := NewAESGCMProvider(config)
```

#### CreateProvider

```go
func (f *ProviderFactory) CreateProvider(config *EncryptionConfig) (EncryptionProvider, error)
```

Factory method for creating encryption providers.

**Parameters:**
- `config`: Encryption configuration

**Returns:**
- Provider interface or error

**Supported Methods:**
- `aes-gcm`: AES-256-GCM encryption
- `plaintext`: No encryption (for testing)

## Database Package

### Interfaces

#### Database

```go
type Database interface {
    // Connection management
    Open() error
    Close() error

    // Transaction support
    BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)

    // Query operations
    Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
    QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row
    Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error)

    // Prepared statements
    Prepare(ctx context.Context, query string) (*sql.Stmt, error)

    // Health checks
    Ping(ctx context.Context) error
    Stats() sql.DBStats
}
```

### Configuration

#### DatabaseConfig

```go
type DatabaseConfig struct {
    // Connection settings
    Path           string `json:"path"`
    MaxConnections int    `json:"max_connections"`
    MaxIdleConns   int    `json:"max_idle_conns"`
    ConnMaxLifetime time.Duration `json:"conn_max_lifetime"`

    // Performance settings
    CacheSize    int    `json:"cache_size"`
    WALMode      bool   `json:"wal_mode"`
    Synchronous  string `json:"synchronous"`
    BusyTimeout  int    `json:"busy_timeout"`

    // Encryption
    EncryptionConfig *encryption.EncryptionConfig `json:"encryption,omitempty"`
}
```

### Operations

#### Session Management

```go
type SessionOps struct {
    db Database
    encryption *FieldMiddleware
}

// Methods
func (s *SessionOps) CreateSession(ctx context.Context, session *Session) error
func (s *SessionOps) GetSession(ctx context.Context, id string) (*Session, error)
func (s *SessionOps) UpdateSession(ctx context.Context, session *Session) error
func (s *SessionOps) DeleteSession(ctx context.Context, id string) error
func (s *SessionOps) ListSessions(ctx context.Context, opts ListOptions) ([]*Session, error)
```

#### Event Management

```go
type EventOps struct {
    db Database
    encryption *FieldMiddleware
}

// Methods
func (e *EventOps) CreateEvent(ctx context.Context, event *Event) error
func (e *EventOps) GetEvent(ctx context.Context, id string) (*Event, error)
func (e *EventOps) ListEvents(ctx context.Context, sessionID string) ([]*Event, error)
func (e *EventOps) CountEvents(ctx context.Context, sessionID string) (int64, error)
```

## CLI Commands

### Command Structure

```go
type Command struct {
    Use   string
    Short string
    Long  string
    Run   func(cmd *cobra.Command, args []string) error
}
```

### Available Commands

#### Root Command

```go
var rootCmd = &cobra.Command{
    Use:   "context-extender",
    Short: "A CLI tool for managing context extensions",
}
```

#### Encrypt Commands

```go
var encryptCmd = &cobra.Command{
    Use:   "encrypt",
    Short: "Manage encryption settings",
}

// Subcommands
encrypt init     // Initialize encryption
encrypt status   // Show encryption status
encrypt rotate   // Rotate encryption keys
encrypt migrate  // Migrate unencrypted data
encrypt verify   // Verify encrypted data
```

#### Config Commands

```go
var configCmd = &cobra.Command{
    Use:   "config",
    Short: "Manage configuration",
}

// Subcommands
config show      // Display configuration
config init      // Initialize configuration
config set       // Set configuration value
config get       // Get configuration value
config validate  // Validate configuration
```

#### Database Commands

```go
var dbCmd = &cobra.Command{
    Use:   "db",
    Short: "Database operations",
}

// Subcommands
db migrate       // Run migrations
db backup        // Create backup
db restore       // Restore from backup
db vacuum        // Optimize database
db stats         // Show statistics
```

## HTTP API (Future)

### Endpoints

#### Health Check

```
GET /health
```

**Response:**
```json
{
  "status": "healthy",
  "version": "1.0.0",
  "uptime": 3600,
  "checks": {
    "database": "ok",
    "encryption": "ok"
  }
}
```

#### Metrics

```
GET /metrics
```

**Response:** Prometheus-formatted metrics

#### Sessions API

```
POST   /api/v1/sessions       Create session
GET    /api/v1/sessions/{id}  Get session
PUT    /api/v1/sessions/{id}  Update session
DELETE /api/v1/sessions/{id}  Delete session
GET    /api/v1/sessions       List sessions
```

#### Events API

```
POST   /api/v1/events         Create event
GET    /api/v1/events/{id}    Get event
GET    /api/v1/events         List events
```

## Error Codes

### Encryption Errors

| Code | Error | Description |
|------|-------|-------------|
| E001 | ErrNotInitialized | Provider not initialized |
| E002 | ErrKeyNotFound | Encryption key not found |
| E003 | ErrInvalidKey | Invalid encryption key |
| E004 | ErrDecryptionFailed | Failed to decrypt data |
| E005 | ErrKeyIDMismatch | Key ID doesn't match |

### Database Errors

| Code | Error | Description |
|------|-------|-------------|
| D001 | ErrDatabaseClosed | Database connection closed |
| D002 | ErrTransactionFailed | Transaction failed |
| D003 | ErrRecordNotFound | Record not found |
| D004 | ErrDuplicateKey | Duplicate key violation |
| D005 | ErrLockTimeout | Database lock timeout |

## Examples

### Complete Encryption Workflow

```go
package main

import (
    "context"
    "log"

    "context-extender/internal/encryption"
    "context-extender/internal/database"
)

func main() {
    // Configure encryption
    encConfig := &encryption.EncryptionConfig{
        Method:           encryption.EncryptionMethodAESGCM,
        KeyPath:          "./encryption.key",
        PBKDF2Iterations: 100000,
        Enabled:          true,
    }

    // Create provider
    provider, err := encryption.NewAESGCMProvider(encConfig)
    if err != nil {
        log.Fatal(err)
    }
    defer provider.Close()

    // Initialize
    if err := provider.Initialize(); err != nil {
        log.Fatal(err)
    }

    // Configure database with encryption
    dbConfig := &database.DatabaseConfig{
        Path:             "./data.db",
        MaxConnections:   10,
        WALMode:          true,
        EncryptionConfig: encConfig,
    }

    // Create database
    db, err := database.New(dbConfig)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Use with encryption middleware
    middleware := encryption.NewFieldMiddleware(provider, nil)

    // Encrypt data
    data := map[string]interface{}{
        "user_id": "123",
        "token":   "secret-token",
        "metadata": map[string]string{
            "key": "value",
        },
    }

    encrypted, err := middleware.EncryptRow("sessions", data)
    if err != nil {
        log.Fatal(err)
    }

    // Store in database...
}
```

### Custom Provider Implementation

```go
type CustomProvider struct {
    key []byte
}

func (p *CustomProvider) Initialize() error {
    // Initialize your provider
    return nil
}

func (p *CustomProvider) EncryptField(data []byte) (*FieldEncryption, error) {
    // Implement encryption logic
    return &FieldEncryption{
        EncryptedData: encrypted,
        KeyID:         "custom-key-1",
        Method:        "custom",
    }, nil
}

// Implement other interface methods...
```

## Migration Guide

### From v1 to v2

```go
// v1 (CGO-based)
import "github.com/mattn/go-sqlite3"

// v2 (Pure Go)
import "modernc.org/sqlite"

// Update connection string
// v1: "file:data.db?_foreign_keys=on"
// v2: "file:data.db?_pragma=foreign_keys(1)"
```

## Performance Tips

1. **Batch Operations**: Use bulk encryption/decryption for multiple fields
2. **Connection Pooling**: Configure appropriate pool size
3. **Selective Encryption**: Only encrypt sensitive fields
4. **Caching**: Cache frequently accessed decrypted data
5. **Async Processing**: Use goroutines for parallel operations

## Security Best Practices

1. **Key Storage**: Never store keys in code or version control
2. **Key Rotation**: Rotate keys periodically
3. **Audit Logging**: Log all encryption operations
4. **Access Control**: Restrict key file permissions (600)
5. **Backup Keys**: Maintain secure key backups