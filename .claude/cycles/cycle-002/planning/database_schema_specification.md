# Database Schema Specification - Cycle 2
**Project**: Context Extender CLI Tool
**Phase**: Planning Phase
**Date**: 2025-09-16
**Component**: SQLite Database Design

## 📊 **Database Architecture Overview**

### **Technology Stack**
- **Database**: SQLite 3.45+ (latest stable)
- **Driver**: github.com/mattn/go-sqlite3
- **Encryption**: SQLCipher 4.x via mutecomm/go-sqlcipher
- **Query Builder**: sqlc for type-safe SQL generation
- **Migration Tool**: golang-migrate/migrate

### **Database Configuration**
```sql
-- Performance and reliability settings
PRAGMA journal_mode = WAL;          -- Write-Ahead Logging for concurrent reads
PRAGMA synchronous = NORMAL;        -- Balance between safety and performance
PRAGMA temp_store = MEMORY;         -- Use memory for temporary tables
PRAGMA mmap_size = 30000000000;     -- Memory-mapped I/O for performance
PRAGMA page_size = 4096;            -- Optimal page size for most systems
PRAGMA cache_size = 10000;          -- ~40MB cache (4KB pages)
PRAGMA foreign_keys = ON;           -- Enforce referential integrity
PRAGMA busy_timeout = 10000;        -- 10-second timeout for write locks

-- Security settings (when using SQLCipher)
PRAGMA cipher_page_size = 4096;     -- Match page_size for encrypted DB
PRAGMA kdf_iter = 256000;           -- Key derivation iterations
PRAGMA cipher_hmac = ON;            -- Enable HMAC for authentication
PRAGMA cipher_plaintext_header_size = 0;  -- Fully encrypted header
```

## 🗄️ **Core Database Schema**

### **1. Sessions Table**
Stores conversation session metadata and lifecycle information.

```sql
CREATE TABLE sessions (
    -- Primary identification
    id TEXT PRIMARY KEY NOT NULL,                    -- UUID from Claude or generated

    -- Session metadata
    start_time DATETIME NOT NULL,                    -- Session start timestamp
    end_time DATETIME,                              -- Session end (NULL if active)
    working_dir TEXT NOT NULL,                      -- Working directory path
    project_name TEXT,                              -- Extracted project name

    -- Session state
    status TEXT NOT NULL DEFAULT 'active'           -- active, completed, archived
        CHECK(status IN ('active', 'completed', 'archived')),
    event_count INTEGER DEFAULT 0,                  -- Total events in session

    -- Source information
    source TEXT NOT NULL DEFAULT 'hook'             -- hook, import, manual
        CHECK(source IN ('hook', 'import', 'manual')),
    import_path TEXT,                                -- Original file path if imported

    -- Timestamps
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    -- Constraints
    CHECK(end_time IS NULL OR end_time >= start_time),
    CHECK(event_count >= 0)
);

-- Indexes for common queries
CREATE INDEX idx_sessions_working_dir ON sessions(working_dir);
CREATE INDEX idx_sessions_start_time ON sessions(start_time DESC);
CREATE INDEX idx_sessions_status ON sessions(status) WHERE status = 'active';
CREATE INDEX idx_sessions_project ON sessions(project_name);
CREATE INDEX idx_sessions_source ON sessions(source);
```

### **2. Events Table**
Stores individual conversation events (prompts, responses, tool usage).

```sql
CREATE TABLE events (
    -- Primary identification
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id TEXT NOT NULL,

    -- Event metadata
    event_type TEXT NOT NULL                        -- Event type enumeration
        CHECK(event_type IN (
            'session_start',
            'user_prompt',
            'claude_response',
            'tool_use',
            'notification',
            'session_end'
        )),
    timestamp DATETIME NOT NULL,
    sequence_number INTEGER NOT NULL,               -- Order within session

    -- Event content
    data TEXT NOT NULL,                             -- JSON blob of event data
    data_size INTEGER GENERATED ALWAYS AS (LENGTH(data)) STORED,

    -- Additional metadata
    tool_name TEXT,                                 -- Tool name if tool_use event
    error_flag BOOLEAN DEFAULT FALSE,               -- Mark error events

    -- Timestamps
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    -- Foreign key
    FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE,

    -- Constraints
    UNIQUE(session_id, sequence_number),
    CHECK(sequence_number > 0),
    CHECK(LENGTH(data) > 0)
);

-- Indexes for performance
CREATE INDEX idx_events_session_id ON events(session_id);
CREATE INDEX idx_events_timestamp ON events(timestamp);
CREATE INDEX idx_events_type ON events(event_type);
CREATE INDEX idx_events_session_sequence ON events(session_id, sequence_number);
CREATE INDEX idx_events_tool ON events(tool_name) WHERE tool_name IS NOT NULL;
```

### **3. Conversations Table**
Stores processed conversation data with analysis and metadata.

```sql
CREATE TABLE conversations (
    -- Primary identification (same as session_id)
    id TEXT PRIMARY KEY NOT NULL,

    -- Conversation metadata
    title TEXT,                                     -- Generated or extracted title
    summary TEXT,                                   -- AI-generated summary
    duration_seconds INTEGER,                       -- Total conversation duration

    -- Analysis data (JSON arrays)
    topics TEXT,                                    -- JSON array of extracted topics
    keywords TEXT,                                  -- JSON array of keywords
    tools_used TEXT,                               -- JSON array of tools used

    -- Statistics
    total_events INTEGER NOT NULL DEFAULT 0,
    user_prompts INTEGER DEFAULT 0,
    claude_responses INTEGER DEFAULT 0,
    total_tokens INTEGER,                          -- Token count if available

    -- Content
    conversation_data TEXT NOT NULL,               -- Complete conversation JSON
    conversation_size INTEGER GENERATED ALWAYS AS (LENGTH(conversation_data)) STORED,

    -- Timestamps
    processed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    -- Foreign key
    FOREIGN KEY (id) REFERENCES sessions(id) ON DELETE CASCADE,

    -- Constraints
    CHECK(total_events > 0),
    CHECK(LENGTH(conversation_data) > 0)
);

-- Indexes for search and analysis
CREATE INDEX idx_conversations_title ON conversations(title);
CREATE INDEX idx_conversations_processed ON conversations(processed_at DESC);
-- Full-text search on summary
CREATE VIRTUAL TABLE conversations_fts USING fts5(
    id UNINDEXED,
    title,
    summary,
    topics,
    keywords,
    content=conversations,
    content_rowid=rowid
);

-- Triggers to keep FTS index updated
CREATE TRIGGER conversations_ai AFTER INSERT ON conversations BEGIN
    INSERT INTO conversations_fts(rowid, id, title, summary, topics, keywords)
    VALUES (new.rowid, new.id, new.title, new.summary, new.topics, new.keywords);
END;

CREATE TRIGGER conversations_ad AFTER DELETE ON conversations BEGIN
    DELETE FROM conversations_fts WHERE rowid = old.rowid;
END;

CREATE TRIGGER conversations_au AFTER UPDATE ON conversations BEGIN
    UPDATE conversations_fts
    SET title = new.title,
        summary = new.summary,
        topics = new.topics,
        keywords = new.keywords
    WHERE rowid = new.rowid;
END;
```

### **4. Import History Table**
Tracks Claude JSONL imports for deduplication and auditing.

```sql
CREATE TABLE import_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    -- Import metadata
    file_path TEXT NOT NULL,                       -- Source file path
    file_hash TEXT NOT NULL,                       -- SHA-256 of file content
    file_size INTEGER NOT NULL,

    -- Import results
    sessions_imported INTEGER DEFAULT 0,
    events_imported INTEGER DEFAULT 0,
    import_status TEXT NOT NULL DEFAULT 'pending'
        CHECK(import_status IN ('pending', 'success', 'partial', 'failed')),
    error_message TEXT,

    -- Timestamps
    imported_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    -- Constraints
    UNIQUE(file_hash),                            -- Prevent duplicate imports
    CHECK(file_size > 0)
);

-- Index for lookup
CREATE INDEX idx_import_history_hash ON import_history(file_hash);
CREATE INDEX idx_import_history_status ON import_history(import_status);
```

### **5. Settings Table**
Stores application configuration and preferences.

```sql
CREATE TABLE settings (
    key TEXT PRIMARY KEY NOT NULL,
    value TEXT NOT NULL,
    description TEXT,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    CHECK(LENGTH(key) > 0),
    CHECK(LENGTH(value) > 0)
);

-- Default settings
INSERT INTO settings (key, value, description) VALUES
    ('schema_version', '1', 'Database schema version'),
    ('encryption_enabled', 'true', 'Whether database encryption is active'),
    ('auto_import', 'true', 'Automatically import new Claude conversations'),
    ('import_directory', '', 'Claude conversation directory path'),
    ('graphql_enabled', 'false', 'Whether GraphQL endpoint is active'),
    ('graphql_port', '8080', 'GraphQL server port');
```

## 🔄 **Migration Scripts**

### **Initial Migration (001_initial_schema.up.sql)**
```sql
-- Enable required extensions
PRAGMA foreign_keys = ON;
PRAGMA journal_mode = WAL;

-- Create all tables
-- [Include all CREATE TABLE statements from above]

-- Create initial indexes
-- [Include all CREATE INDEX statements from above]

-- Insert default settings
-- [Include INSERT statements for settings]
```

### **Migration Rollback (001_initial_schema.down.sql)**
```sql
-- Disable foreign keys during drop
PRAGMA foreign_keys = OFF;

-- Drop tables in reverse order
DROP TABLE IF EXISTS conversations_fts;
DROP TABLE IF EXISTS conversations;
DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS import_history;
DROP TABLE IF EXISTS settings;

-- Re-enable foreign keys
PRAGMA foreign_keys = ON;
```

## 📈 **Performance Optimization**

### **Index Strategy**
```sql
-- Composite indexes for common queries
CREATE INDEX idx_events_session_type_time
    ON events(session_id, event_type, timestamp);

CREATE INDEX idx_sessions_dir_time
    ON sessions(working_dir, start_time DESC);

-- Partial indexes for active sessions
CREATE INDEX idx_active_sessions
    ON sessions(status, start_time DESC)
    WHERE status = 'active';

-- Covering index for session queries
CREATE INDEX idx_sessions_covering
    ON sessions(id, working_dir, project_name, status, start_time);
```

### **Query Optimization Examples**
```sql
-- Efficient active session lookup
SELECT id, working_dir, event_count,
       (julianday('now') - julianday(start_time)) * 86400 as duration_seconds
FROM sessions
WHERE status = 'active' AND working_dir = ?
ORDER BY start_time DESC
LIMIT 1;

-- Fast event retrieval with pagination
SELECT e.*, s.working_dir, s.project_name
FROM events e
INNER JOIN sessions s ON e.session_id = s.id
WHERE e.session_id = ?
ORDER BY e.sequence_number
LIMIT ? OFFSET ?;

-- Full-text search across conversations
SELECT c.id, c.title, c.summary,
       snippet(conversations_fts, -1, '<mark>', '</mark>', '...', 32) as highlight
FROM conversations c
INNER JOIN conversations_fts ON c.id = conversations_fts.id
WHERE conversations_fts MATCH ?
ORDER BY rank
LIMIT 20;
```

## 🔐 **Security Considerations**

### **Encryption Configuration**
```go
// SQLCipher initialization
func InitializeEncryptedDB(dbPath string, key []byte) (*sql.DB, error) {
    // Format key as hex string
    keyHex := hex.EncodeToString(key)

    // Connection string with encryption parameters
    dsn := fmt.Sprintf("%s?_pragma_key=x'%s'&_pragma_cipher_page_size=4096",
        dbPath, keyHex)

    db, err := sql.Open("sqlite3", dsn)
    if err != nil {
        return nil, err
    }

    // Verify encryption is working
    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("encryption verification failed: %w", err)
    }

    // Apply security settings
    pragmas := []string{
        "PRAGMA cipher_hmac = ON",
        "PRAGMA kdf_iter = 256000",
        "PRAGMA cipher_plaintext_header_size = 0",
    }

    for _, pragma := range pragmas {
        if _, err := db.Exec(pragma); err != nil {
            return nil, fmt.Errorf("failed to set %s: %w", pragma, err)
        }
    }

    return db, nil
}
```

### **SQL Injection Prevention**
- All queries use parameterized statements
- Input validation on all user-provided data
- Prepared statement caching for performance
- Query builder (sqlc) generates type-safe code

## 📊 **Capacity Planning**

### **Storage Estimates**
```
Per Session:
- Session record: ~200 bytes
- Events (avg 50/session): ~25KB
- Conversation data: ~50KB
Total per session: ~75KB

Storage projections:
- 100 sessions: ~7.5MB
- 1,000 sessions: ~75MB
- 10,000 sessions: ~750MB

With encryption overhead (+15%):
- 10,000 sessions: ~860MB
```

### **Performance Targets**
```
Operation targets:
- Session creation: <2ms
- Event insertion: <1ms
- Batch import: 1000 events/second
- Query response: <10ms for indexed queries
- Full-text search: <50ms for 10,000 documents
```

## ✅ **Schema Validation Checklist**

- [x] Foreign key relationships properly defined
- [x] Indexes on all foreign keys
- [x] Unique constraints where appropriate
- [x] Check constraints for data validation
- [x] Generated columns for computed values
- [x] Timestamps on all tables
- [x] UTF-8 support for all text fields
- [x] Full-text search capability
- [x] Audit trail via import_history
- [x] Settings storage for configuration

---

**Status**: ✅ **SPECIFICATION COMPLETE**
**Next**: Hook integration architecture design