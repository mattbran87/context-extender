# Context-Extender Compression Hooks

## Overview

Context-Extender now includes sophisticated compression detection and context preservation features to handle Claude conversation compressions gracefully. When a conversation gets compressed, critical context is preserved and can be reinjected to maintain continuity.

## What Gets Lost During Compression

During conversation compression, typically the following is lost:

1. **Implementation Details**
   - Specific code snippets and solutions
   - Technical decisions and their rationale
   - Error messages and debugging paths
   - Configuration details

2. **User Preferences & Context**
   - Working style preferences
   - Project-specific conventions
   - Implicit agreements
   - Communication patterns

3. **Progressive Refinements**
   - Iterative improvements made
   - Failed approaches to avoid
   - Optimization decisions
   - Trade-offs accepted

4. **Relationship Context**
   - Trust level established
   - Areas of user autonomy
   - Topics to avoid or emphasize
   - Informal decisions

## Setting Up Compression Hooks

### 1. Install the Hooks

```bash
# Configure Context-Extender with compression support
context-extender configure

# The following hooks will be installed:
# - conversation-compress: Captures compression events
# - context-request: Handles context reinjection requests
```

### 2. Hook Configuration

Add to your Claude Code settings:

```json
{
  "hooks": {
    "onConversationCompress": "context-extender capture --event=conversation-compress --data='${context}'",
    "onContextRequest": "context-extender capture --event=context-request"
  }
}
```

## Using Context Preservation

### Analyze Current Context

Analyze your current session to see what critical context would be preserved:

```bash
# Analyze current session
context-extender context analyze --session=current

# Output in JSON format
context-extender context analyze --session=current --json
```

### Preserve Context Manually

Manually trigger context preservation before an expected compression:

```bash
# Preserve current session context
context-extender context preserve --session=current

# This creates a preservation event that can be retrieved later
```

### Reinject Context After Compression

After a compression occurs, reinject the preserved context:

```bash
# Generate reinjection prompt
context-extender context reinject --session=current

# This outputs a formatted prompt with all critical context
```

## Example Workflow

### Before Compression

```bash
# Working on Context-Extender project...
# Conversation getting long, compression may occur soon

# Manually preserve context
$ context-extender context preserve --session=current
✅ Context preserved successfully!
Session: context-extender-dev
Preserved: 5 key decisions, 3 constraints, 4 preferences
```

### After Compression

```bash
# Conversation was compressed, need to restore context
$ context-extender context reinject --session=current

## Critical Context (Post-Compression)

**Project**: Context-Extender CLI Tool
**Current Objective**: Implementing compression hooks
**Phase**: Feature Development

### Technical Decisions
- database: Pure Go SQLite using modernc.org/sqlite
- architecture: Zero CGO dependencies required

### Constraints
- Zero CGO dependencies required
- Cross-platform compatibility mandatory

### User Preferences
- workflow: Simplified 5-day adaptive cycles
- documentation: Minimal, value-focused documentation

### Recently Completed
- v1.0.1 release with critical fixes
- Compression hook implementation

### Known Issues to Avoid
- Import cycles in Go packages
```

## Context Analysis Features

The context analyzer automatically identifies:

- **Project Information**: Name, current objectives, phase
- **Technical Stack**: Languages, frameworks, tools
- **Constraints**: Technical limitations, requirements
- **Key Decisions**: Architectural choices, design patterns
- **User Preferences**: Workflow style, communication patterns
- **Completed Work**: Recent achievements, milestones
- **Pending Tasks**: Current objectives, next steps
- **Errors to Avoid**: Known issues, antipatterns

## Integration with Claude Code

### Automatic Capture

When properly configured, Context-Extender automatically:
1. Detects conversation compression events
2. Extracts critical context from the conversation
3. Stores it in the database
4. Makes it available for reinjection

### Manual Intervention

You can also manually:
1. Analyze context at any time
2. Preserve context preemptively
3. Reinject context when needed
4. Review preservation history

## Advanced Usage

### Custom Context Patterns

You can customize what gets preserved by modifying the extraction patterns in:
```go
// internal/context/preserver.go
func ExtractCriticalContext(conversations []string) (*CompressionSummary, error)
```

### Context Storage

Preserved context is stored as events in the database:
- Event Type: `context_preservation` or `compression`
- Data: JSON-encoded CompressionSummary
- Retrieval: Latest preservation per session

### GraphQL Queries

Query preserved context via GraphQL:
```graphql
{
  events(eventType: "context_preservation") {
    id
    sessionId
    timestamp
    data
  }
}
```

## Benefits

1. **Continuity**: Maintain project context across compressions
2. **Efficiency**: Quickly restore working state
3. **Consistency**: Preserve technical decisions and constraints
4. **Productivity**: Reduce time explaining context again
5. **Quality**: Avoid repeating resolved issues

## Troubleshooting

### Context Not Preserving

```bash
# Check if hooks are installed
context-extender configure --status

# Verify database is working
context-extender database status

# Check for preservation events
context-extender query list --type=events
```

### Reinject Not Working

```bash
# Ensure session has preserved context
context-extender context analyze --session=current

# Check for preservation events
context-extender graphql exec "{ events(eventType: \"context_preservation\") { id sessionId timestamp } }"
```

## Future Enhancements

Planned improvements:
- Automatic compression detection
- Smart context extraction using AI
- Context priority levels
- Cross-session context sharing
- Context templates for common scenarios
- Integration with Claude's native context system

## Summary

The compression hook feature ensures that critical project context survives conversation compressions, enabling seamless continuity in long-running development sessions. By automatically preserving and allowing reinjection of key information, Context-Extender helps maintain productivity and reduces friction in the development workflow.