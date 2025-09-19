# Context Extender CLI - Testing Guide

**Version**: v2.0.0 (Database Integration)
**Date**: September 2025
**Target Audience**: Testers and Early Adopters

## 🎯 What is Context Extender?

Context Extender is a CLI tool that automatically captures Claude Code conversations and stores them in a secure SQLite database. It enables you to:
- **Preserve conversation history** beyond Claude's session limits
- **Search and analyze** your Claude interactions
- **Import existing conversations** from Claude's native files
- **Query data** using a powerful GraphQL interface

## 📋 Testing Checklist Overview

This guide covers testing of 4 major features:
- ✅ **Database Setup** (Core functionality)
- ✅ **Import System** (Claude conversation import)
- ✅ **Encryption** (Database security)
- ✅ **GraphQL API** (Query interface)

---

## 🚀 Getting Started

### Prerequisites
- Windows 10/11, macOS, or Linux
- Go 1.21+ (if building from source)
- C compiler (GCC on Linux/macOS, MinGW on Windows) for SQLite support
- Access to existing Claude conversations (optional, for import testing)

### Installation

#### Option 1: Pre-built Binary (Recommended for Testing)
```bash
# Download the appropriate binary for your platform
# Extract to a directory in your PATH
```

#### Option 2: Build from Source
```bash
git clone <repository-url>
cd context-extender

# Install dependencies (requires C compiler)
go mod tidy

# Build with SQLite support
CGO_ENABLED=1 go build -tags sqlite3 -o context-extender.exe .

# Verify installation
./context-extender.exe version
```

**Expected Output:**
```
Context Extender CLI v2.0.0
Build Date: 2025-09-18
Platform: windows/amd64
```

---

## 🧪 Test Scenario 1: Basic Database Setup

### Test 1.1: Initialize Database
```bash
# Initialize the database
./context-extender.exe database init
```

**Expected Results:**
- ✅ Database file created at `~/.context-extender/conversations.db`
- ✅ Success message: "Database initialized successfully"
- ✅ All tables created (sessions, events, conversations, etc.)

### Test 1.2: Check Database Status
```bash
# Check database status
./context-extender.exe database status
```

**Expected Output:**
```
Database Status:
  Path: C:\Users\[username]\.context-extender\conversations.db
  Size: [size] bytes
  Modified: [timestamp]
  Connection: ✓ Active

Table Statistics:
  Sessions: 0
  Events: 0
  Conversations: 0
  Import History: 0
```

### Test 1.3: Test Database Commands
```bash
# Test session creation
./context-extender.exe database capture session-start test-session-123

# Test user prompt capture
./context-extender.exe database capture user-prompt test-session-123 "Hello, this is a test message"

# Test session end
./context-extender.exe database capture session-end test-session-123

# Check status again - should show 1 session, 2 events
./context-extender.exe database status
```

**Expected Results:**
- ✅ Commands execute without errors
- ✅ Database status shows: Sessions: 1, Events: 2

---

## 📥 Test Scenario 2: Import System

### Test 2.1: Discover Claude Conversations
```bash
# Run the import wizard
./context-extender.exe import wizard
```

**Expected Behavior:**
1. ✅ Searches for Claude conversation files automatically
2. ✅ Lists found projects and file counts
3. ✅ Shows interactive menu with options:
   - Import all conversations
   - Import specific project
   - Skip duplicate imports
   - Cancel

### Test 2.2: Test Auto Import
```bash
# Auto-import all Claude conversations
./context-extender.exe import auto --dry-run --verbose
```

**Expected Results:**
- ✅ Finds Claude conversation files in standard locations:
  - `~/.claude/projects/`
  - `~/Library/Application Support/Claude/projects/` (macOS)
  - `~/AppData/Roaming/Claude/projects/` (Windows)
- ✅ Shows preview of files to be imported
- ✅ With `--dry-run`: Shows what would be imported without changes

### Test 2.3: Import Specific File
```bash
# Find a Claude JSONL file and import it
./context-extender.exe import file ~/.claude/projects/[project-name]/[uuid].jsonl --verbose
```

**Expected Results:**
- ✅ Successfully parses Claude JSONL format
- ✅ Imports messages, sessions, and metadata
- ✅ Records import in import_history table
- ✅ Shows progress: "✅ Successfully imported session [session-id]"

### Test 2.4: Check Import History
```bash
# View import history
./context-extender.exe import history
```

**Expected Output:**
```
📜 Import History:
────────────────────────────────────────────────────────────────────────────────

File: [filename].jsonl
  Path:     [full-path]
  Imported: 2025-09-18 14:30:00
  Messages: [number]
  Checksum: [hash]...

Total imported files: [count]
```

---

## 🔐 Test Scenario 3: Database Encryption

### Test 3.1: Initialize Encrypted Database
```bash
# Initialize with encryption (will generate key)
./context-extender.exe encrypt init
```

**Expected Results:**
- ✅ Generates new encryption key
- ✅ Saves key to `~/.context-extender/keys/db.key`
- ✅ Creates encrypted database
- ✅ Shows warning: "⚠️ IMPORTANT: Back up your encryption key"

### Test 3.2: Verify Encryption
```bash
# Check encryption status
./context-extender.exe encrypt verify
```

**Expected Output:**
```
✅ Database encryption verified successfully
   SQLite version: [version]
   Encryption: ENABLED
   Database path: [path]

Database Statistics:
   Sessions: [count]
   Events: [count]
```

### Test 3.3: Key Management
```bash
# View key information
./context-extender.exe encrypt key-info
```

**Expected Output:**
```
Encryption Key Information:
  Version: 1
  Created: 2025-09-18 14:30:00
  Last Rotated: 2025-09-18 14:30:00
  Rotation Count: 0
  Algorithm: PBKDF2-SHA256
  Iterations: 256000

Key Location: [path]
```

### Test 3.4: Convert Existing Database
```bash
# Convert unencrypted to encrypted (if you have unencrypted)
./context-extender.exe encrypt convert
```

**Expected Behavior:**
- ✅ Detects existing unencrypted database
- ✅ Creates encrypted copy
- ✅ Asks for confirmation before replacing original
- ✅ Backs up original database

---

## 🔍 Test Scenario 4: GraphQL Query Interface

### Test 4.1: Quick Statistics
```bash
# Get database statistics
./context-extender.exe graphql stats
```

**Expected Output:**
```
📊 Database Statistics
=====================
Sessions:      [count]
Conversations: [count]
Events:        [count]
Imports:       [count]
Oldest:        [timestamp]
Newest:        [timestamp]
```

### Test 4.2: Search Functionality
```bash
# Search conversations
./context-extender.exe graphql search "database" --limit 5
```

**Expected Results:**
- ✅ Returns matching conversations
- ✅ Shows session IDs and timestamps
- ✅ Truncates long content appropriately
- ✅ Shows total match count

### Test 4.3: Direct Query Execution
```bash
# Execute GraphQL query directly
./context-extender.exe graphql exec "{ sessions(limit: 3) { id createdAt status } }" --pretty
```

**Expected Output:**
```json
{
  "data": {
    "sessions": [
      {
        "id": "[session-id]",
        "createdAt": "[timestamp]",
        "status": "[status]"
      }
    ]
  }
}
```

### Test 4.4: Interactive GraphQL Server
```bash
# Start GraphQL server
./context-extender.exe graphql server --port 8080
```

**Expected Behavior:**
1. ✅ Server starts on http://localhost:8080
2. ✅ GraphQL endpoint available at `/graphql`
3. ✅ Interactive playground at `/`
4. ✅ Can execute queries through web interface

**Test in browser**: Navigate to `http://localhost:8080`
- ✅ Shows GraphQL playground interface
- ✅ Example queries execute successfully
- ✅ Auto-completion and syntax highlighting work

### Test 4.5: View Query Examples
```bash
# Show example queries
./context-extender.exe graphql examples
```

**Expected Output:**
- ✅ Shows 7+ example queries with descriptions
- ✅ Includes copy-paste ready queries
- ✅ Covers stats, sessions, search, events, conversations

---

## 🔧 Test Scenario 5: Advanced Features

### Test 5.1: Hook Integration (If Available)
```bash
# Configure Claude Code hooks
./context-extender.exe configure

# Check hook status
./context-extender.exe configure --status
```

**Expected Results:**
- ✅ Detects Claude Code installation
- ✅ Configures hooks in settings.json
- ✅ Shows hook status and session correlation

### Test 5.2: Data Export/Backup
```bash
# Create backup of database
cp ~/.context-extender/conversations.db ./backup-conversations.db

# For encrypted databases, test decryption
./context-extender.exe encrypt decrypt
```

### Test 5.3: Performance Testing
```bash
# Import large conversation file
./context-extender.exe import file [large-file.jsonl] --verbose

# Test search performance
time ./context-extender.exe graphql search "common-term"

# Test stats performance
time ./context-extender.exe graphql stats
```

---

## 🐛 Expected Issues and Troubleshooting

### Common Issues

#### Issue 1: "cgo: C compiler not found"
**Symptoms**: Build or CGO-related commands fail
**Solution**: Install C compiler:
- **Windows**: Install MinGW-w64 or Visual Studio Build Tools
- **macOS**: Install Xcode Command Line Tools: `xcode-select --install`
- **Linux**: Install gcc: `sudo apt-get install gcc` (Ubuntu) or `sudo yum install gcc` (CentOS)

#### Issue 2: "Database not initialized"
**Symptoms**: Commands fail with database connection errors
**Solution**: Run `./context-extender.exe database init`

#### Issue 3: "Permission denied" on key files
**Symptoms**: Encryption operations fail
**Solution**: Check file permissions on `~/.context-extender/keys/` directory

#### Issue 4: No Claude conversations found
**Symptoms**: Import wizard finds no files
**Solution**:
- Check Claude installation location
- Look for `.claude` directory in home folder
- Try manual path: `./context-extender.exe import dir /path/to/claude/projects`

### Performance Expectations

#### Acceptable Performance
- **Database init**: < 5 seconds
- **Small import** (< 1MB): < 10 seconds
- **GraphQL query**: < 2 seconds
- **Search**: < 5 seconds for reasonable datasets

#### Memory Usage
- **Normal operation**: < 50MB
- **Large imports**: < 200MB
- **GraphQL server**: < 100MB

---

## 📊 Test Results Template

### Tester Information
- **Name**: _______________
- **Platform**: _______________
- **Go Version**: _______________
- **Test Date**: _______________

### Test Results

| Test Scenario | Status | Notes |
|---------------|--------|-------|
| Database Setup | ✅ ❌ | |
| Import Auto-Discovery | ✅ ❌ | |
| Import Single File | ✅ ❌ | |
| Encryption Init | ✅ ❌ | |
| Encryption Verify | ✅ ❌ | |
| GraphQL Stats | ✅ ❌ | |
| GraphQL Search | ✅ ❌ | |
| GraphQL Server | ✅ ❌ | |
| GraphQL Examples | ✅ ❌ | |

### Overall Assessment

**Functionality**: ⭐⭐⭐⭐⭐ (1-5 stars)
**Performance**: ⭐⭐⭐⭐⭐ (1-5 stars)
**Documentation**: ⭐⭐⭐⭐⭐ (1-5 stars)
**Ease of Use**: ⭐⭐⭐⭐⭐ (1-5 stars)

**Additional Comments:**
```
[Your feedback here]
```

---

## 🆘 Getting Help

### Command Help
```bash
# General help
./context-extender.exe --help

# Command-specific help
./context-extender.exe database --help
./context-extender.exe import --help
./context-extender.exe encrypt --help
./context-extender.exe graphql --help
```

### Debug Mode
Add `--verbose` to most commands for detailed output:
```bash
./context-extender.exe import auto --verbose
./context-extender.exe database status --verbose
```

### Log Files
Check for log files in:
- `~/.context-extender/logs/`
- Current directory (if running from development)

### Support Information
For bug reports, please include:
1. Platform and Go version
2. Exact command that failed
3. Full error message
4. Output of `./context-extender.exe database status`

---

## ✅ Success Criteria

The Context Extender CLI tool is considered **ready for production** if:

### Core Functionality (Must Pass)
- ✅ Database initializes successfully
- ✅ Can capture conversation events
- ✅ Import discovers and processes Claude files
- ✅ GraphQL queries return expected results
- ✅ Encryption setup works correctly

### Performance (Should Pass)
- ✅ Commands respond within expected timeframes
- ✅ Memory usage stays within reasonable limits
- ✅ Large imports complete successfully

### User Experience (Should Pass)
- ✅ Help text is clear and accurate
- ✅ Error messages are informative
- ✅ Interactive features work as expected

**Thank you for testing Context Extender! Your feedback helps make the tool better for everyone.** 🚀