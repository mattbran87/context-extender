# Implementation Architecture - Cycle 3

**Cycle**: 3 (Production Deployment Readiness)
**Phase**: Planning Phase
**Date**: September 18, 2025
**Architecture Focus**: Pure Go SQLite + Application Encryption

---

## 🏗️ Architecture Overview

### Current Architecture (Cycle 2)
```
Context Extender CLI
├── cmd/ (Cobra CLI)
├── internal/database/ (CGO SQLite + SQLCipher)
├── internal/importer/ (Claude JSONL)
├── internal/graphql/ (Query API)
└── internal/hooks/ (Claude Code integration)
```

### Target Architecture (Cycle 3)
```
Context Extender CLI (Pure Go)
├── cmd/ (Cobra CLI - unchanged)
├── internal/database/ (Abstraction Layer)
│   ├── interface.go (DatabaseBackend interface)
│   ├── manager.go (Backend selection and management)
│   ├── purgo/ (Pure Go SQLite implementation)
│   └── legacy/ (CGO SQLite for compatibility)
├── internal/encryption/ (Application-level encryption)
│   ├── aes.go (AES-256-GCM implementation)
│   ├── keys.go (Key management)
│   └── fields.go (Field-level encryption)
├── internal/importer/ (Claude JSONL - unchanged)
├── internal/graphql/ (Query API - updated for new backend)
└── internal/hooks/ (Claude Code integration - unchanged)
```

---

## 📦 Package Structure Design

### New Package: `internal/database/`

#### Interface Definition (`interface.go`)
```go
package database

import (
    "context"
    "time"
)

// DatabaseBackend defines the interface all database implementations must satisfy
type DatabaseBackend interface {
    // Lifecycle Management
    Initialize(ctx context.Context, config *Config) error
    Close() error
    Ping(ctx context.Context) error

    // Schema Management
    CreateSchema(ctx context.Context) error
    GetSchemaVersion(ctx context.Context) (int, error)
    MigrateSchema(ctx context.Context, targetVersion int) error

    // Session Operations
    CreateSession(ctx context.Context, session *Session) error
    GetSession(ctx context.Context, id string) (*Session, error)
    UpdateSession(ctx context.Context, session *Session) error
    DeleteSession(ctx context.Context, id string) error
    ListSessions(ctx context.Context, filters *SessionFilters) ([]*Session, error)

    // Event Operations
    CreateEvent(ctx context.Context, event *Event) error
    GetEventsBySession(ctx context.Context, sessionID string) ([]*Event, error)
    CreateEventBatch(ctx context.Context, events []*Event) error

    // Conversation Operations
    CreateConversation(ctx context.Context, conv *Conversation) error
    GetConversationsBySession(ctx context.Context, sessionID string) ([]*Conversation, error)
    SearchConversations(ctx context.Context, query string, limit int) ([]*Conversation, error)

    // Statistics
    GetDatabaseStats(ctx context.Context) (*DatabaseStats, error)

    // Backend Information
    GetBackendInfo() *BackendInfo
}

// Config represents database configuration
type Config struct {
    Backend         BackendType        `json:"backend"`
    DatabasePath    string            `json:"database_path"`
    EncryptionEnabled bool            `json:"encryption_enabled"`
    EncryptionConfig  *EncryptionConfig `json:"encryption_config,omitempty"`
    ConnectionTimeout time.Duration   `json:"connection_timeout"`
    QueryTimeout     time.Duration    `json:"query_timeout"`
    BackendOptions   map[string]interface{} `json:"backend_options,omitempty"`
}

type BackendType string

const (
    BackendPureGoSQLite BackendType = "pure_go_sqlite"
    BackendCGOSQLite    BackendType = "cgo_sqlite"
    BackendAuto         BackendType = "auto"
)

type BackendInfo struct {
    Name         string   `json:"name"`
    Version      string   `json:"version"`
    RequiresCGO  bool     `json:"requires_cgo"`
    Features     []string `json:"features"`
    Capabilities map[string]bool `json:"capabilities"`
}
```

#### Backend Manager (`manager.go`)
```go
package database

import (
    "context"
    "fmt"
    "sync"
)

// Manager handles backend selection and lifecycle
type Manager struct {
    config         *Config
    backend        DatabaseBackend
    encryptionLayer *encryption.Layer
    mu             sync.RWMutex
    registry       map[BackendType]BackendFactory
}

type BackendFactory interface {
    CreateBackend(config *Config) (DatabaseBackend, error)
    IsAvailable() bool
    GetCapabilities() *BackendCapabilities
}

type BackendCapabilities struct {
    SupportsEncryption   bool
    SupportsFullText     bool
    SupportsTransactions bool
    RequiresCGO          bool
    PlatformSupport      []string
}

func NewManager(config *Config) *Manager {
    m := &Manager{
        config:   config,
        registry: make(map[BackendType]BackendFactory),
    }

    // Register available backends
    m.registerBackends()

    return m
}

func (m *Manager) registerBackends() {
    // Register Pure Go SQLite backend
    m.registry[BackendPureGoSQLite] = &purgo.Factory{}

    // Register CGO SQLite backend if available
    if cgoAvailable() {
        m.registry[BackendCGOSQLite] = &legacy.Factory{}
    }
}

func (m *Manager) Initialize(ctx context.Context) error {
    m.mu.Lock()
    defer m.mu.Unlock()

    backendType := m.config.Backend
    if backendType == BackendAuto {
        backendType = m.autoSelectBackend()
    }

    factory, exists := m.registry[backendType]
    if !exists {
        return fmt.Errorf("backend %s not available", backendType)
    }

    backend, err := factory.CreateBackend(m.config)
    if err != nil {
        return fmt.Errorf("failed to create backend: %w", err)
    }

    if err := backend.Initialize(ctx, m.config); err != nil {
        return fmt.Errorf("failed to initialize backend: %w", err)
    }

    m.backend = backend

    // Initialize encryption layer if enabled
    if m.config.EncryptionEnabled {
        m.encryptionLayer = encryption.NewLayer(m.config.EncryptionConfig)
    }

    return nil
}

func (m *Manager) autoSelectBackend() BackendType {
    // Prefer Pure Go SQLite for CGO-free operation
    if factory := m.registry[BackendPureGoSQLite]; factory.IsAvailable() {
        return BackendPureGoSQLite
    }

    // Fallback to CGO SQLite if available
    if factory := m.registry[BackendCGOSQLite]; factory.IsAvailable() {
        return BackendCGOSQLite
    }

    // Default to Pure Go SQLite
    return BackendPureGoSQLite
}
```

### New Package: `internal/database/purgo/`

#### Pure Go SQLite Implementation (`backend.go`)
```go
package purgo

import (
    "context"
    "database/sql"
    "fmt"

    _ "modernc.org/sqlite"
    "context-extender/internal/encryption"
)

type Backend struct {
    db              *sql.DB
    config          *database.Config
    encryptionLayer *encryption.Layer
    statements      map[string]*sql.Stmt
}

func (b *Backend) Initialize(ctx context.Context, config *database.Config) error {
    // Open pure Go SQLite database
    db, err := sql.Open("sqlite", config.DatabasePath)
    if err != nil {
        return fmt.Errorf("failed to open pure Go SQLite database: %w", err)
    }

    // Configure SQLite for optimal performance
    if err := b.configureSQLite(ctx, db); err != nil {
        return fmt.Errorf("failed to configure SQLite: %w", err)
    }

    b.db = db
    b.config = config

    // Initialize encryption if enabled
    if config.EncryptionEnabled {
        b.encryptionLayer = encryption.NewLayer(config.EncryptionConfig)
    }

    // Prepare common statements
    if err := b.prepareStatements(ctx); err != nil {
        return fmt.Errorf("failed to prepare statements: %w", err)
    }

    return nil
}

func (b *Backend) configureSQLite(ctx context.Context, db *sql.DB) error {
    // Enable WAL mode for better concurrency
    if _, err := db.ExecContext(ctx, "PRAGMA journal_mode=WAL"); err != nil {
        return err
    }

    // Set connection limits
    db.SetMaxOpenConns(10)
    db.SetMaxIdleConns(5)
    db.SetConnMaxLifetime(time.Hour)

    // Enable foreign keys
    if _, err := db.ExecContext(ctx, "PRAGMA foreign_keys=ON"); err != nil {
        return err
    }

    // Optimize for performance
    pragmas := []string{
        "PRAGMA synchronous=NORMAL",
        "PRAGMA cache_size=10000",
        "PRAGMA temp_store=memory",
        "PRAGMA mmap_size=268435456", // 256MB
    }

    for _, pragma := range pragmas {
        if _, err := db.ExecContext(ctx, pragma); err != nil {
            return fmt.Errorf("failed to set pragma %s: %w", pragma, err)
        }
    }

    return nil
}

func (b *Backend) CreateSession(ctx context.Context, session *database.Session) error {
    // Apply encryption if enabled
    if b.encryptionLayer != nil {
        if err := b.encryptionLayer.EncryptSession(session); err != nil {
            return fmt.Errorf("failed to encrypt session: %w", err)
        }
    }

    stmt := b.statements["insert_session"]
    _, err := stmt.ExecContext(ctx,
        session.ID,
        session.CreatedAt,
        session.UpdatedAt,
        session.Status,
        session.Metadata)

    return err
}

func (b *Backend) GetSession(ctx context.Context, id string) (*database.Session, error) {
    stmt := b.statements["get_session"]
    row := stmt.QueryRowContext(ctx, id)

    session := &database.Session{}
    err := row.Scan(
        &session.ID,
        &session.CreatedAt,
        &session.UpdatedAt,
        &session.Status,
        &session.Metadata,
    )

    if err != nil {
        return nil, err
    }

    // Apply decryption if enabled
    if b.encryptionLayer != nil {
        if err := b.encryptionLayer.DecryptSession(session); err != nil {
            return nil, fmt.Errorf("failed to decrypt session: %w", err)
        }
    }

    return session, nil
}

func (b *Backend) GetBackendInfo() *database.BackendInfo {
    return &database.BackendInfo{
        Name:        "Pure Go SQLite",
        Version:     "modernc.org/sqlite v1.39.0",
        RequiresCGO: false,
        Features: []string{
            "application_encryption",
            "full_text_search",
            "cross_platform",
            "no_cgo_dependency",
        },
        Capabilities: map[string]bool{
            "encryption":   true,
            "full_text":    true,
            "transactions": true,
            "concurrent":   true,
        },
    }
}
```

### New Package: `internal/encryption/`

#### AES Encryption Implementation (`aes.go`)
```go
package encryption

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
    "fmt"
)

type Layer struct {
    gcm            cipher.AEAD
    encryptedFields map[string]bool
}

func NewLayer(config *EncryptionConfig) *Layer {
    key := deriveKeyFromConfig(config)

    block, err := aes.NewCipher(key)
    if err != nil {
        panic(fmt.Sprintf("failed to create cipher: %v", err))
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        panic(fmt.Sprintf("failed to create GCM: %v", err))
    }

    return &Layer{
        gcm: gcm,
        encryptedFields: map[string]bool{
            "conversation_content": true,
            "session_metadata":     true,
            "event_data":          true,
        },
    }
}

func (l *Layer) EncryptSession(session *database.Session) error {
    if l.encryptedFields["session_metadata"] && session.Metadata != "" {
        encrypted, err := l.encryptString(session.Metadata)
        if err != nil {
            return err
        }
        session.Metadata = encrypted
    }
    return nil
}

func (l *Layer) DecryptSession(session *database.Session) error {
    if l.encryptedFields["session_metadata"] && session.Metadata != "" {
        decrypted, err := l.decryptString(session.Metadata)
        if err != nil {
            return err
        }
        session.Metadata = decrypted
    }
    return nil
}

func (l *Layer) encryptString(plaintext string) (string, error) {
    if plaintext == "" {
        return "", nil
    }

    // Generate random nonce
    nonce := make([]byte, l.gcm.NonceSize())
    if _, err := rand.Read(nonce); err != nil {
        return "", err
    }

    // Encrypt the data
    ciphertext := l.gcm.Seal(nonce, nonce, []byte(plaintext), nil)

    // Return base64 encoded result
    return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (l *Layer) decryptString(ciphertext string) (string, error) {
    if ciphertext == "" {
        return "", nil
    }

    // Decode from base64
    data, err := base64.StdEncoding.DecodeString(ciphertext)
    if err != nil {
        return "", err
    }

    // Extract nonce and ciphertext
    nonceSize := l.gcm.NonceSize()
    if len(data) < nonceSize {
        return "", fmt.Errorf("ciphertext too short")
    }

    nonce, ciphertext := data[:nonceSize], data[nonceSize:]

    // Decrypt the data
    plaintext, err := l.gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return "", err
    }

    return string(plaintext), nil
}
```

---

## 🔄 Migration Strategy

### Database Migration Architecture

#### Migration Manager (`internal/migration/`)
```go
package migration

type Manager struct {
    sourceBackend DatabaseBackend
    targetBackend DatabaseBackend
    encryptionLayer *encryption.Layer
}

func (m *Manager) MigrateDatabase(ctx context.Context) error {
    // 1. Backup current database
    if err := m.createBackup(ctx); err != nil {
        return err
    }

    // 2. Initialize target backend
    if err := m.targetBackend.CreateSchema(ctx); err != nil {
        return err
    }

    // 3. Migrate data with progress tracking
    return m.migrateAllData(ctx)
}

func (m *Manager) migrateAllData(ctx context.Context) error {
    // Migrate sessions
    sessions, err := m.sourceBackend.ListSessions(ctx, nil)
    if err != nil {
        return err
    }

    for i, session := range sessions {
        if err := m.migrateSession(ctx, session); err != nil {
            return err
        }

        // Progress reporting every 100 sessions
        if i%100 == 0 {
            fmt.Printf("Migrated %d/%d sessions\n", i, len(sessions))
        }
    }

    return nil
}
```

### Data Format Migration

#### Encryption Format Migration
```go
// Handle migration from SQLCipher to application-level encryption
func (m *Manager) migrateEncryptionFormat(session *database.Session) error {
    // If source is encrypted with SQLCipher, data is already decrypted
    // Apply new application-level encryption
    if m.encryptionLayer != nil {
        return m.encryptionLayer.EncryptSession(session)
    }
    return nil
}
```

---

## 🧪 Testing Architecture

### Testing Framework Structure

#### Backend Testing (`internal/database/testing/`)
```go
package testing

type BackendTestSuite struct {
    backend DatabaseBackend
    config  *database.Config
    ctx     context.Context
}

func (s *BackendTestSuite) TestSessionCRUD() error {
    // Create test session
    session := &database.Session{
        ID:        "test-session-1",
        CreatedAt: time.Now(),
        Status:    "active",
        Metadata:  "test metadata",
    }

    // Test Create
    if err := s.backend.CreateSession(s.ctx, session); err != nil {
        return fmt.Errorf("create failed: %w", err)
    }

    // Test Read
    retrieved, err := s.backend.GetSession(s.ctx, session.ID)
    if err != nil {
        return fmt.Errorf("get failed: %w", err)
    }

    // Verify data integrity
    if retrieved.Metadata != session.Metadata {
        return fmt.Errorf("metadata mismatch: got %s, want %s",
            retrieved.Metadata, session.Metadata)
    }

    return nil
}

func TestAllBackends(t *testing.T) {
    backends := []database.BackendType{
        database.BackendPureGoSQLite,
    }

    for _, backendType := range backends {
        t.Run(string(backendType), func(t *testing.T) {
            suite := NewBackendTestSuite(backendType)
            if err := suite.RunComplianceTests(); err != nil {
                t.Fatalf("Backend %s failed: %v", backendType, err)
            }
        })
    }
}
```

### Performance Testing

#### Benchmark Framework
```go
func BenchmarkPureGoSQLite(b *testing.B) {
    backend := setupPureGoBackend()
    ctx := context.Background()

    b.Run("SessionCreate", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            session := generateTestSession(i)
            backend.CreateSession(ctx, session)
        }
    })

    b.Run("SessionGet", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            backend.GetSession(ctx, fmt.Sprintf("session-%d", i%1000))
        }
    })

    b.Run("ConversationSearch", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            backend.SearchConversations(ctx, "test query", 10)
        }
    })
}
```

---

## 📦 Build and Distribution Architecture

### GitHub Actions Workflow (`.github/workflows/release.yml`)
```yaml
name: Release Binaries

on:
  push:
    tags:
      - 'v*'

jobs:
  build:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        goos: [linux, windows, darwin]
        goarch: [amd64, arm64]
        exclude:
          - goos: windows
            goarch: arm64

    steps:
    - uses: actions/checkout@v4

    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'

    - name: Build binary
      env:
        GOOS: ${{ matrix.goos }}
        GOARCH: ${{ matrix.goarch }}
        CGO_ENABLED: 0
      run: |
        go build -ldflags="-w -s" -o context-extender-${{ matrix.goos }}-${{ matrix.goarch }} .

    - name: Generate checksums
      run: |
        sha256sum context-extender-* > checksums.txt

    - name: Upload to release
      uses: softprops/action-gh-release@v1
      with:
        files: |
          context-extender-*
          checksums.txt
```

### Binary Distribution Strategy
1. **Automated builds** for Windows, macOS, Linux
2. **Checksum verification** for security
3. **GitHub Releases** for distribution
4. **Installation scripts** for different platforms

---

## 🎯 Performance Optimization Strategy

### Database Performance
1. **Connection pooling** with optimal settings
2. **Prepared statements** for common queries
3. **WAL mode** for concurrent access
4. **Index optimization** for search operations

### Encryption Performance
1. **AES-256-GCM** for authenticated encryption
2. **Field-level encryption** only for sensitive data
3. **Batch operations** for multiple encryptions
4. **Memory-efficient** key handling

### Build Performance
1. **CGO-free builds** for fast compilation
2. **Static linking** for portable binaries
3. **Size optimization** with build flags
4. **Parallel builds** in CI/CD

---

**Architecture Status**: ✅ **READY FOR IMPLEMENTATION**
**Complexity**: Well-defined interfaces and clear separation of concerns
**Risk Level**: 🟢 **LOW** - Proven patterns and technologies

---

*Implementation Architecture by: Development Team*
*Date: September 18, 2025*
*Planning Phase: Day 2 of 4*