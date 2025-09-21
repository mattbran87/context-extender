# Cycle 5 Implementation Plan (Revised)

## Major Decision: Remove GraphQL, Focus on Export

**Rationale**: GraphQL adds complexity without value, causes bugs, and no users depend on it.

---

## Priority 1: Remove GraphQL & Fix Query System (Day 1)

### 🚨 **CRITICAL: Remove GraphQL Dependencies**

**Actions**:
1. Remove GraphQL dependency from query system
2. Delete dual database initialization
3. Ensure query commands use unified database backend
4. Remove GraphQL commands from CLI
5. Clean up codebase

**Files to Modify**:
- `cmd/query.go` - Update to use database.Manager directly
- `cmd/graphql.go` - DELETE entire file
- `cmd/root.go` - Remove GraphQL command registration
- `internal/graphql/` - DELETE entire package
- `internal/converter/converter.go` - Update to use new backend

**Files to Delete**:
```
cmd/graphql.go
internal/graphql/
internal/api/  (if GraphQL-specific)
```

**Success Criteria**:
- ✅ `context-extender query list` shows captured conversations
- ✅ `context-extender query show [session-id]` displays details
- ✅ No GraphQL commands in help text
- ✅ Single database initialization path
- ✅ All tests pass

---

## Priority 2: CSV & JSON Export (Days 2-3)

### 📊 **Core Export Functionality**

**New Command Structure**:
```bash
context-extender export --format csv --output conversations.csv
context-extender export --format json --output backup.json
context-extender export --format csv --from 2024-01-01 --to 2024-01-31
context-extender export --format json --session session-123
```

**CSV Export Features**:
- Configurable columns
- Date range filtering
- Project filtering
- Session filtering
- Automatic Excel-friendly formatting

**JSON Export Features**:
- Full conversation data
- Metadata preservation
- Structured output for programmatic access
- Optional pretty-printing
- Compression support for large exports

**Implementation**:
```go
// cmd/export.go (NEW)
type ExportCommand struct {
    Format   string   // csv, json
    Output   string   // output file path
    Columns  []string // for CSV
    From     string   // date filter
    To       string   // date filter
    Project  string   // project filter
    Sessions []string // specific sessions
}

// internal/export/exporter.go (NEW)
type Exporter interface {
    Export(data []Conversation, options ExportOptions) error
}

type CSVExporter struct{}
type JSONExporter struct{}
```

**Default CSV Columns**:
```
session_id, project, start_time, end_time, duration, user_prompts, claude_responses, total_words, status
```

---

## Priority 3: Slash Commands for Export (Day 4)

### 💬 **Claude Slash Commands**

**Commands to Create**:
```markdown
# .claude/commands/cx-export.md
Export conversations to CSV or JSON format

# .claude/commands/cx-search.md
Search conversation history

# .claude/commands/cx-stats.md
Show conversation statistics

# .claude/commands/cx-help.md
Show Context-Extender commands
```

**Installation Command**:
```bash
context-extender slash-commands install
# Creates .claude/commands/ directory
# Generates command files
# Shows success message
```

**Implementation**:
```go
// cmd/slash_commands.go (NEW)
type SlashCommandGenerator struct {
    commands []SlashCommand
}

func (g *SlashCommandGenerator) Install() error {
    // Create .claude/commands/ directory
    // Generate .md files for each command
    // Each file calls context-extender CLI
}
```

---

## Priority 4: Enhanced Query Commands (Day 5)

### 🔍 **Better Query Interface**

**Improvements**:
1. Fix query to work with new backend (after GraphQL removal)
2. Add export flags to query commands
3. Improve output formatting
4. Add summary statistics

**Enhanced Commands**:
```bash
# Direct export from query
context-extender query list --export csv --output sessions.csv
context-extender query search "database" --export json

# Better formatting
context-extender query list --format detailed
context-extender query stats --format json
```

---

## Technical Implementation Plan

### Day 1: GraphQL Removal
- [ ] Morning: Remove GraphQL from query system
- [ ] Afternoon: Delete GraphQL packages
- [ ] Evening: Test all query commands work
- [ ] Release: v1.0.2 with fix

### Day 2: Export Foundation
- [ ] Morning: Create export command structure
- [ ] Afternoon: Implement CSV exporter
- [ ] Evening: Add filtering options

### Day 3: Export Enhancement
- [ ] Morning: Implement JSON exporter
- [ ] Afternoon: Add column configuration
- [ ] Evening: Test with real data

### Day 4: Slash Commands
- [ ] Morning: Create command generator
- [ ] Afternoon: Generate export commands
- [ ] Evening: Test in Claude Code

### Day 5: Polish & Release
- [ ] Morning: Enhanced query commands
- [ ] Afternoon: Documentation
- [ ] Evening: Release v1.1.0

---

## Migration Guide for GraphQL Users

```markdown
## GraphQL Removal Notice

GraphQL has been removed in v1.1.0 to simplify the codebase and fix bugs.

### Migration Options:

**For programmatic access:**
- Use: `context-extender export --format json`
- Provides: Complete structured data
- Better: No server required, direct file output

**For queries:**
- Use: `context-extender query` commands
- Provides: Same data access via CLI
- Better: Simpler, faster, no initialization issues

**For analysis:**
- Use: `context-extender export --format csv`
- Provides: Excel-compatible data
- Better: Direct analysis in spreadsheet tools

If you need GraphQL, please open an issue explaining your use case.
```

---

## Success Metrics

### Critical (Day 1)
- ✅ Query commands show captured data
- ✅ No GraphQL code remains
- ✅ Single database initialization

### High Priority (Days 2-3)
- ✅ CSV export works with Excel
- ✅ JSON export preserves all data
- ✅ Filtering options functional

### Important (Days 4-5)
- ✅ Slash commands installable
- ✅ Export accessible from Claude
- ✅ Documentation complete

---

## Risk Mitigation

### Risk: Breaking existing workflows
**Mitigation**:
- Clear migration guide
- JSON export replaces GraphQL data access
- Keep v1.0.1 available for transition

### Risk: Large dataset exports
**Mitigation**:
- Implement streaming for CSV
- Add pagination for JSON
- Progress indicators for long operations

### Risk: Complex filtering needs
**Mitigation**:
- Start with basic filters
- Add advanced filters based on feedback
- Document filter combinations

---

## Deliverables

### v1.0.2 (Day 1)
- GraphQL removed
- Query system fixed
- Bug resolved

### v1.1.0-beta (Day 3)
- CSV export working
- JSON export working
- Basic filtering

### v1.1.0 (Day 5)
- Full export capabilities
- Slash commands
- Complete documentation
- Migration guide

---

## Communication Plan

### v1.0.2 Release Notes
```
FIXED: Query commands now show captured conversations
REMOVED: GraphQL support (deprecated)
IMPROVED: Simplified database initialization

Breaking change: GraphQL commands removed.
Use export --format json for programmatic access.
```

### v1.1.0 Release Notes
```
NEW: Export to CSV and JSON formats
NEW: Slash commands for Claude integration
NEW: Advanced filtering for exports
REMOVED: GraphQL completely removed
IMPROVED: Query performance and reliability

Migration: See docs/MIGRATION.md for GraphQL alternatives
```

---

## Benefits of This Approach

1. **Immediate Bug Fix**: Query works on Day 1
2. **Simpler Codebase**: ~30% less code to maintain
3. **User-Focused Features**: Export is what users actually need
4. **Better Performance**: No dual initialization overhead
5. **Clearer Purpose**: CLI tool, not an API server

---

**This plan delivers MORE VALUE with LESS COMPLEXITY**