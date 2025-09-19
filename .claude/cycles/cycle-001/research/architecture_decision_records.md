# Architecture Decision Records - Context-Extender

**Date**: 2024-09-16
**Phase**: Research - Day 1-2
**Cycle**: 001

---

## ADR-001: CLI Framework Selection

**Status**: Decided
**Date**: 2024-09-16
**Context**: Need to select a Go CLI framework for context-extender tool

### Decision
Use **spf13/cobra** for CLI framework

### Rationale
- **Mature Ecosystem**: 35k+ GitHub stars, extensive community support
- **Feature Rich**: Excellent documentation, command hierarchies, flag management
- **Industry Adoption**: Used by major Go projects (Kubernetes, Docker, Hugo)
- **Viper Integration**: Strong configuration management through Viper
- **Long-term Maintainability**: Better support for complex CLI evolution

### Alternatives Considered
- **urfave/cli**: Simpler but less feature-rich, limited command hierarchy support
- **kingpin**: Less active development, superseded by kong
- **Standard library**: Too low-level for complex CLI requirements

### Consequences
- **Positive**: Strong foundation for feature expansion, excellent documentation
- **Negative**: Slightly larger binary size than minimal alternatives
- **Neutral**: Learning curve for team members unfamiliar with Cobra

### Implementation Notes
```go
import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
    Use:   "context-extender",
    Short: "Manage Claude Code conversation context",
    Long:  "Automatically capture and manage Claude Code conversations...",
}
```

---

## ADR-002: Storage Strategy (File vs Database Progression)

**Status**: Decided
**Date**: 2024-09-16
**Context**: Need to choose data storage approach for conversation data

### Decision
Start with **JSON file storage** for Cycle 1 MVP, design for database migration

### Rationale

#### Cycle 1 Benefits (File Storage)
- **Faster Development**: No external dependencies (SQLite, PostgreSQL)
- **Cross-Platform Simplicity**: No installation complexity
- **Easy Backup/Sharing**: Files are portable and human-readable
- **Sufficient Scale**: Adequate for MVP user base and data volumes
- **Debugging Friendly**: Direct file inspection during development

#### Future Migration Strategy
- **Repository Pattern**: Abstract storage implementation for easy swapping
- **Database Migration Tools**: Graceful transition from files to database
- **Planned Timeline**: SQLite (Cycle 2-3), PostgreSQL (Cycle 4+)

### Implementation Architecture
```go
// Repository pattern enables future database migration
type ConversationRepository interface {
    Save(conversation *Conversation) error
    FindByID(id string) (*Conversation, error)
    Search(query SearchQuery) ([]*Conversation, error)
    Delete(id string) error
}

// File implementation for Cycle 1
type FileRepository struct {
    basePath string
    // ... implementation details
}

// Future database implementations
type SQLiteRepository struct { /* ... */ }
type PostgreSQLRepository struct { /* ... */ }
```

### Storage Structure
```
~/.context-extender/
├── conversations/
│   ├── active/           # JSONL files for ongoing sessions
│   │   └── session-abc123.jsonl
│   └── completed/        # JSON files for finished conversations
│       └── 2024-09-15/
│           └── session-abc123.json
├── config.json           # User configuration
└── index.json           # Search index for performance
```

### Migration Timeline
- **Cycle 1-2**: JSON file storage with optimization
- **Cycle 3-4**: SQLite for local performance improvements
- **Cycle 5+**: PostgreSQL for team sharing features

---

## ADR-003: Claude Code Integration Approach

**Status**: Decided
**Date**: 2024-09-16
**Context**: Need to integrate with Claude Code to capture user requests and responses

### Decision
Use **Claude Code hooks system** with asynchronous processing

### Integration Points
1. **SessionStart Hook**: Initialize conversation tracking
2. **UserPromptSubmit Hook**: Capture user input in real-time
3. **Stop Hook**: Capture Claude responses and session state
4. **SessionEnd Hook**: Finalize conversation and trigger archival

### Hook Configuration Strategy
```json
{
  "hooks": {
    "SessionStart": [{
      "matcher": "",
      "hooks": [{
        "type": "command",
        "command": "context-extender capture --event=session-start",
        "timeout": 30
      }]
    }],
    "UserPromptSubmit": [{
      "matcher": "",
      "hooks": [{
        "type": "command",
        "command": "context-extender capture --event=user-prompt",
        "timeout": 30
      }]
    }],
    "Stop": [{
      "matcher": "",
      "hooks": [{
        "type": "command",
        "command": "context-extender capture --event=claude-response",
        "timeout": 30
      }]
    }],
    "SessionEnd": [{
      "matcher": "",
      "hooks": [{
        "type": "command",
        "command": "context-extender capture --event=session-end",
        "timeout": 30
      }]
    }]
  }
}
```

### Processing Architecture
```
Hook Trigger → Lightweight Capture Script → Background Processor → File Storage
     ↓              ↓                          ↓                    ↓
 <50ms         Queue Event               Process with             Atomic
Response      (Non-blocking)            Context (<5s)           Write Ops
```

### Benefits
- **Minimal Performance Impact**: Hooks exit quickly, processing happens async
- **Resilient Integration**: Uses stable Claude Code hook API
- **Cross-Platform**: Works consistently across Windows/Mac/Linux
- **Graceful Degradation**: Continues working if context-extender unavailable

### Risk Mitigations
- **Hook Failures**: Retry logic and comprehensive error logging
- **Format Changes**: Flexible parsing with version detection
- **Large Files**: Streaming processing for memory efficiency
- **Privacy**: Configuration for sensitive data filtering

### Future Enhancements
- Real-time session sharing between multiple Claude Code instances
- Cloud storage integration for team collaboration
- Smart context suggestions based on conversation patterns

---

## ADR-004: Data Model and File Formats

**Status**: Decided
**Date**: 2024-09-16
**Context**: Need to define data structures for conversation storage

### Decision
Use **structured JSON/JSONL** with defined schemas for conversation data

### Data Model
```go
// Core conversation structures
type Conversation struct {
    ID           string              `json:"id"`
    StartTime    time.Time           `json:"startTime"`
    EndTime      *time.Time          `json:"endTime,omitempty"`
    ProjectPath  string              `json:"projectPath"`
    Events       []ConversationEvent `json:"events"`
    Summary      string              `json:"summary"`
    Statistics   ConversationStats   `json:"statistics"`
    Metadata     ConversationMeta    `json:"metadata"`
}

type ConversationEvent struct {
    ID          string                 `json:"id"`
    Timestamp   time.Time              `json:"timestamp"`
    Type        EventType              `json:"type"`
    Content     string                 `json:"content"`
    Metadata    map[string]interface{} `json:"metadata"`
}

type EventType string
const (
    EventSessionStart   EventType = "session_start"
    EventUserPrompt     EventType = "user_prompt"
    EventClaudeResponse EventType = "claude_response"
    EventSessionEnd     EventType = "session_end"
)
```

### File Format Strategy
- **Active Sessions**: JSONL format for append-only performance
- **Completed Sessions**: Structured JSON for query efficiency
- **Metadata**: Separate index files for fast searching

### Benefits
- **Schema Evolution**: JSON allows for backward-compatible changes
- **Human Readable**: Easy debugging and manual inspection
- **Tool Compatibility**: Standard formats work with existing tools
- **Performance**: JSONL optimized for real-time append operations

---

## ADR-005: Cross-Platform Path Management

**Status**: Decided
**Date**: 2024-09-16
**Context**: Need consistent file path handling across Windows, Mac, and Linux

### Decision
Use **Go standard library** with platform-specific configuration directories

### Implementation
```go
func getStoragePath() (string, error) {
    switch runtime.GOOS {
    case "windows":
        return filepath.Join(os.Getenv("APPDATA"), "context-extender"), nil
    case "darwin":
        home, err := os.UserHomeDir()
        if err != nil {
            return "", err
        }
        return filepath.Join(home, ".context-extender"), nil
    default: // linux and others
        configDir := os.Getenv("XDG_CONFIG_HOME")
        if configDir == "" {
            home, err := os.UserHomeDir()
            if err != nil {
                return "", err
            }
            configDir = filepath.Join(home, ".config")
        }
        return filepath.Join(configDir, "context-extender"), nil
    }
}
```

### Path Standards
- **Windows**: `%APPDATA%\context-extender\`
- **macOS**: `~/.context-extender/`
- **Linux**: `~/.config/context-extender/` (XDG Base Directory compliant)

### Benefits
- **Native Integration**: Follows OS conventions
- **User Familiarity**: Predictable locations
- **Permission Compliance**: Standard user directories

---

## Summary

These ADRs establish the foundational technical decisions for the Context-Extender CLI tool:

1. **Cobra CLI Framework**: Robust, scalable CLI foundation
2. **File-First Storage**: Simple MVP with clear database migration path
3. **Claude Code Hooks**: Reliable, performant integration strategy
4. **Structured Data Model**: Flexible, evolvable conversation schema
5. **Cross-Platform Paths**: Standard, predictable file locations

All decisions prioritize **MVP delivery speed** while maintaining **future scalability** and **cross-platform compatibility**.