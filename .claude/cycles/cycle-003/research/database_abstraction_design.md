# Database Abstraction Layer Design - Cycle 3 Research

**Research Phase**: Cycle 3 Day 2-3
**Date**: September 18, 2025
**Focus**: Flexible database backend architecture

---

## 🎯 Design Objectives

### Primary Goals
1. **Support multiple database backends** (Pure Go SQLite, CGO SQLite, future alternatives)
2. **Maintain current API compatibility** for minimal migration effort
3. **Enable runtime backend selection** based on user configuration
4. **Provide consistent encryption** across all backends
5. **Simplify testing** with backend switching capabilities

### Success Criteria
- ✅ Drop-in replacement for current database package
- ✅ <5% performance overhead from abstraction
- ✅ Consistent behavior across all backends
- ✅ Simple configuration-based backend selection
- ✅ Comprehensive error handling and fallback mechanisms

---

## 🏗️ Architecture Overview

### Current Architecture (Cycle 2)
```
CLI Commands → Database Package → SQLite+SQLCipher (CGO)
```

### Proposed Architecture (Cycle 3)
```
CLI Commands → Database Interface → Backend Manager → Selected Backend
                                                   ├── Pure Go SQLite
                                                   ├── CGO SQLite
                                                   └── Future Backends
```

---

## 📋 Interface Design

### Core Database Interface

```go
// DatabaseBackend defines the interface all database backends must implement
type DatabaseBackend interface {
    // Connection Management
    Initialize(config *Config) error
    Close() error
    Ping() error

    // Schema Management
    CreateSchema() error
    MigrateSchema(fromVersion, toVersion int) error
    GetSchemaVersion() (int, error)

    // Session Operations
    CreateSession(session *Session) error
    GetSession(id string) (*Session, error)
    UpdateSession(session *Session) error
    DeleteSession(id string) error
    ListSessions(filters *SessionFilters) ([]*Session, error)

    // Event Operations
    CreateEvent(event *Event) error
    GetEventsBySession(sessionID string) ([]*Event, error)
    CreateEventBatch(events []*Event) error

    // Conversation Operations
    CreateConversation(conv *Conversation) error
    GetConversationsBySession(sessionID string) ([]*Conversation, error)
    SearchConversations(query string, limit int) ([]*Conversation, error)

    // Statistics and Analytics
    GetDatabaseStats() (*DatabaseStats, error)
    GetSessionCount() (int, error)
    GetConversationCount() (int, error)

    // Backup and Recovery
    BackupDatabase(path string) error
    RestoreDatabase(path string) error

    // Backend Information
    GetBackendInfo() *BackendInfo
    SupportsFeature(feature string) bool
}
```

### Configuration Interface

```go
type Config struct {
    // Backend Selection
    Backend         BackendType    `json:"backend"`
    DatabasePath    string         `json:"database_path"`

    // Encryption Configuration
    EncryptionEnabled bool         `json:"encryption_enabled"`
    EncryptionConfig  *EncryptionConfig `json:"encryption_config,omitempty"`

    // Performance Tuning
    ConnectionPool    *PoolConfig   `json:"connection_pool,omitempty"`
    QueryTimeout      time.Duration `json:"query_timeout"`

    // Backend-Specific Settings
    BackendOptions    map[string]interface{} `json:"backend_options,omitempty"`
}

type BackendType string

const (
    BackendPureGoSQLite BackendType = "pure_go_sqlite"
    BackendCGOSQLite    BackendType = "cgo_sqlite"
    BackendBadgerDB     BackendType = "badgerdb"
    BackendAuto         BackendType = "auto"
)

type EncryptionConfig struct {
    Method          EncryptionMethod `json:"method"`
    KeyPath         string           `json:"key_path"`
    KeyDerivation   *KDFConfig       `json:"key_derivation,omitempty"`
    FieldEncryption []string         `json:"field_encryption,omitempty"`
}
```

### Backend Manager

```go
type BackendManager struct {
    config          *Config
    currentBackend  DatabaseBackend
    encryptionLayer *EncryptionLayer
    registry        map[BackendType]BackendFactory
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
    SupportsConcurrency  bool
    RequiresCGO          bool
    PlatformSupport      []Platform
}
```

---

## 🔧 Backend Implementations

### Pure Go SQLite Backend

```go
type PureGoSQLiteBackend struct {
    db              *sql.DB
    config          *Config
    encryptionLayer *EncryptionLayer
    queryBuilder    *QueryBuilder
}

func (p *PureGoSQLiteBackend) Initialize(config *Config) error {
    // Open modernc.org/sqlite database
    db, err := sql.Open("sqlite", config.DatabasePath)
    if err != nil {
        return fmt.Errorf("failed to open pure Go SQLite: %w", err)
    }

    // Configure WAL mode and performance settings
    p.configureDatabase(db)

    // Initialize encryption layer if enabled
    if config.EncryptionEnabled {
        p.encryptionLayer = NewApplicationEncryption(config.EncryptionConfig)
    }

    p.db = db
    return nil
}

func (p *PureGoSQLiteBackend) CreateSession(session *Session) error {
    // Encrypt sensitive fields if encryption enabled
    if p.encryptionLayer != nil {
        if err := p.encryptionLayer.EncryptSession(session); err != nil {
            return err
        }
    }

    // Standard database insertion
    query := `INSERT INTO sessions (id, created_at, status, metadata) VALUES (?, ?, ?, ?)`
    _, err := p.db.Exec(query, session.ID, session.CreatedAt, session.Status, session.Metadata)
    return err
}

func (p *PureGoSQLiteBackend) GetBackendInfo() *BackendInfo {
    return &BackendInfo{
        Name:        "Pure Go SQLite",
        Version:     "modernc.org/sqlite v1.39.0",
        RequiresCGO: false,
        Features: []string{
            "application_encryption",
            "full_text_search",
            "cross_platform",
        },
    }
}
```

### CGO SQLite Backend

```go
type CGOSQLiteBackend struct {
    db              *sql.DB
    config          *Config
    queryBuilder    *QueryBuilder
}

func (c *CGOSQLiteBackend) Initialize(config *Config) error {
    var dsnURI string

    if config.EncryptionEnabled {
        // Use SQLCipher with encryption key
        key := config.EncryptionConfig.Key
        dsnURI = fmt.Sprintf("%s?_pragma_key=x'%s'", config.DatabasePath, key)
    } else {
        dsnURI = config.DatabasePath
    }

    db, err := sql.Open("sqlite3", dsnURI)
    if err != nil {
        return fmt.Errorf("failed to open CGO SQLite: %w", err)
    }

    c.db = db
    return nil
}

func (c *CGOSQLiteBackend) GetBackendInfo() *BackendInfo {
    return &BackendInfo{
        Name:        "CGO SQLite",
        Version:     "go-sqlite3 + SQLCipher",
        RequiresCGO: true,
        Features: []string{
            "native_encryption",
            "full_sqlite_features",
            "maximum_performance",
        },
    }
}
```

---

## 🔐 Encryption Layer Design

### Application-Level Encryption

```go
type EncryptionLayer struct {
    cipher          cipher.AEAD
    key            []byte
    encryptedFields map[string]bool
}

func NewApplicationEncryption(config *EncryptionConfig) *EncryptionLayer {
    key := deriveKey(config.KeyPath)

    // Use AES-256-GCM for authenticated encryption
    block, _ := aes.NewCipher(key)
    gcm, _ := cipher.NewGCM(block)

    return &EncryptionLayer{
        cipher: gcm,
        key:    key,
        encryptedFields: map[string]bool{
            "conversation_content": true,
            "session_metadata":     true,
            "event_data":          true,
        },
    }
}

func (e *EncryptionLayer) EncryptSession(session *Session) error {
    if e.encryptedFields["session_metadata"] && session.Metadata != "" {
        encrypted, err := e.encrypt([]byte(session.Metadata))
        if err != nil {
            return err
        }
        session.Metadata = base64.StdEncoding.EncodeToString(encrypted)
    }
    return nil
}

func (e *EncryptionLayer) encrypt(data []byte) ([]byte, error) {
    nonce := make([]byte, e.cipher.NonceSize())
    if _, err := rand.Read(nonce); err != nil {
        return nil, err
    }

    ciphertext := e.cipher.Seal(nonce, nonce, data, nil)
    return ciphertext, nil
}
```

---

## 📊 Migration Strategy

### Backend Detection and Auto-Selection

```go
func (bm *BackendManager) AutoSelectBackend() (BackendType, error) {
    // Priority order for auto-selection
    priorities := []BackendType{
        BackendPureGoSQLite,  // Preferred: No CGO dependency
        BackendCGOSQLite,     // Fallback: If CGO available
    }

    for _, backend := range priorities {
        factory := bm.registry[backend]
        if factory.IsAvailable() {
            log.Printf("Auto-selected backend: %s", backend)
            return backend, nil
        }
    }

    return "", errors.New("no available database backend found")
}
```

### Data Migration Between Backends

```go
type MigrationManager struct {
    sourceBackend DatabaseBackend
    targetBackend DatabaseBackend
}

func (mm *MigrationManager) MigrateData() error {
    // 1. Export all data from source backend
    sessions, err := mm.sourceBackend.ListSessions(nil)
    if err != nil {
        return err
    }

    // 2. Initialize target backend
    if err := mm.targetBackend.CreateSchema(); err != nil {
        return err
    }

    // 3. Migrate data with progress tracking
    for i, session := range sessions {
        if err := mm.migrateSession(session); err != nil {
            return fmt.Errorf("failed to migrate session %s: %w", session.ID, err)
        }

        // Progress reporting
        if i%100 == 0 {
            log.Printf("Migrated %d/%d sessions", i, len(sessions))
        }
    }

    return nil
}
```

---

## 🧪 Testing Strategy

### Backend Testing Framework

```go
type BackendTestSuite struct {
    backend DatabaseBackend
    config  *Config
}

func (ts *BackendTestSuite) RunComplianceTests() error {
    tests := []struct {
        name string
        test func() error
    }{
        {"Session CRUD", ts.testSessionCRUD},
        {"Event Operations", ts.testEventOperations},
        {"Conversation Search", ts.testConversationSearch},
        {"Encryption Consistency", ts.testEncryption},
        {"Performance Benchmarks", ts.testPerformance},
    }

    for _, test := range tests {
        if err := test.test(); err != nil {
            return fmt.Errorf("%s failed: %w", test.name, err)
        }
    }

    return nil
}

func TestAllBackends(t *testing.T) {
    backends := []BackendType{
        BackendPureGoSQLite,
        BackendCGOSQLite,
    }

    for _, backendType := range backends {
        t.Run(string(backendType), func(t *testing.T) {
            suite := NewBackendTestSuite(backendType)
            if err := suite.RunComplianceTests(); err != nil {
                t.Fatalf("Backend %s failed compliance tests: %v", backendType, err)
            }
        })
    }
}
```

---

## 📈 Performance Considerations

### Abstraction Overhead Mitigation

```go
type QueryBuilder struct {
    dialect SQLDialect
    cache   map[string]*sql.Stmt
}

// Pre-compile common queries for each backend
func (qb *QueryBuilder) PrepareCommonQueries(db *sql.DB) error {
    queries := map[string]string{
        "insert_session":      qb.dialect.InsertSessionQuery(),
        "get_session":         qb.dialect.GetSessionQuery(),
        "list_sessions":       qb.dialect.ListSessionsQuery(),
        "search_conversations": qb.dialect.SearchConversationsQuery(),
    }

    for name, query := range queries {
        stmt, err := db.Prepare(query)
        if err != nil {
            return err
        }
        qb.cache[name] = stmt
    }

    return nil
}
```

### Benchmark Framework

```go
func BenchmarkBackendOperations(b *testing.B) {
    backends := []BackendType{BackendPureGoSQLite, BackendCGOSQLite}

    for _, backend := range backends {
        b.Run(string(backend), func(b *testing.B) {
            db := setupBackend(backend)

            b.Run("SessionCreate", func(b *testing.B) {
                for i := 0; i < b.N; i++ {
                    session := generateTestSession()
                    db.CreateSession(session)
                }
            })

            b.Run("ConversationSearch", func(b *testing.B) {
                for i := 0; i < b.N; i++ {
                    db.SearchConversations("test query", 10)
                }
            })
        })
    }
}
```

---

## 🔄 Configuration Management

### Runtime Configuration

```go
// Default configuration with auto-detection
func DefaultConfig() *Config {
    return &Config{
        Backend:      BackendAuto,
        DatabasePath: getDefaultDatabasePath(),
        EncryptionEnabled: true,
        EncryptionConfig: &EncryptionConfig{
            Method:          EncryptionAES256GCM,
            KeyPath:         getDefaultKeyPath(),
            FieldEncryption: []string{"conversation_content", "session_metadata"},
        },
        QueryTimeout: 30 * time.Second,
    }
}

// Environment-based configuration
func ConfigFromEnvironment() *Config {
    config := DefaultConfig()

    if backend := os.Getenv("CONTEXT_EXTENDER_BACKEND"); backend != "" {
        config.Backend = BackendType(backend)
    }

    if dbPath := os.Getenv("CONTEXT_EXTENDER_DB_PATH"); dbPath != "" {
        config.DatabasePath = dbPath
    }

    return config
}
```

---

## 🎯 Implementation Roadmap

### Phase 1: Core Interface (Days 3-4)
1. **Define database interface** with comprehensive method signatures
2. **Implement backend manager** with registration and selection logic
3. **Create configuration system** with validation and defaults
4. **Build query abstraction** for SQL dialect differences

### Phase 2: Pure Go Backend (Days 4-5)
1. **Implement Pure Go SQLite backend** using modernc.org/sqlite
2. **Create application-level encryption** with AES-256-GCM
3. **Build migration tools** from current CGO implementation
4. **Comprehensive testing** with compliance test suite

### Phase 3: CGO Backend Compatibility (Day 5-6)
1. **Adapt current CGO implementation** to new interface
2. **Ensure feature parity** between backends
3. **Performance optimization** and benchmarking
4. **Documentation and examples** for backend selection

### Success Metrics:
- ✅ Interface covers 100% of current functionality
- ✅ Pure Go backend passes all compliance tests
- ✅ Performance within 10% of direct implementation
- ✅ Seamless backend switching at runtime
- ✅ Comprehensive documentation and examples

---

**Design Status**: ✅ **COMPLETE**
**Next Phase**: Proof-of-concept implementation
**Ready for**: Planning phase transition

---

*Database Abstraction Layer Design by: Development Team*
*Date: September 18, 2025*
*Research Phase: Day 3 of 6*