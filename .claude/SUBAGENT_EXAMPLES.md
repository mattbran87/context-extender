# Subagent Usage Examples for Context-Extender

## Real-World Subagent Invocation Examples

### 🔬 Research Phase Examples

#### Example 1: Architecture Exploration
```markdown
PROMPT TO USE:
"As the Architecture Discovery Specialist, explore architecture options for 
implementing context manipulation in the context-extender Go project.

Requirements:
- Must safely extend Claude Code context without breaking existing functionality
- Need high-performance context operations (< 1ms overhead)
- Support concurrent context access from multiple goroutines
- Enable context persistence for session recovery

Consider:
1. In-memory vs persistent storage architectures
2. Synchronous vs asynchronous context updates
3. Event-driven vs polling patterns for Claude Code integration
4. Centralized vs distributed context management

Generate:
- ADR for context storage architecture
- Performance comparison matrix
- Risk assessment for each option
- Recommended approach with rationale"
```

### 📋 Planning Phase Examples

#### Example 2: Epic Breakdown
```markdown
PROMPT TO USE:
"As the Story Refinement Specialist, break down this epic into user stories:

Epic: 'Implement Context Extension Features for Claude Code'
- Users need to extend context with custom metadata
- Context extensions must persist across Claude Code sessions
- Extensions should be searchable and filterable
- Performance must not degrade Claude Code operations

Requirements:
- Each story should be 1-2 days of work
- Include Go technical details
- Define CLI commands for each feature
- Specify Claude Code integration points

Generate INVEST-compliant stories with:
- Clear acceptance criteria
- Technical implementation notes
- Dependencies identified
- Risk assessment"
```

#### Example 3: Implementation Planning
```markdown
PROMPT TO USE:
"As the Implementation Planning Orchestrator, create an implementation schedule 
for these refined stories:

Stories:
1. Core context data structures (CE-001)
2. Context persistence layer (CE-002)
3. CLI context commands (CE-003)
4. Claude Code hook integration (CE-004)
5. Context search functionality (CE-005)

Available capacity: 11 days (Days 5-15)
Team: 1 developer (Claude)

Consider:
- Go package dependencies
- Integration testing windows
- Claude Code compatibility testing
- Performance optimization time

Generate:
- Day-by-day schedule with specific tasks
- Risk-adjusted timeline with buffers
- Critical path identification
- Resource allocation matrix"
```

### 🔨 Implementation Phase Examples

#### Example 4: Test Generation
```markdown
PROMPT TO USE:
"As the Test Automation Specialist, generate comprehensive tests for this 
context manipulation function:

Package: context_handler
Function: ExtendContext(ctx context.Context, key string, value interface{}) (context.Context, error)

Requirements:
- Should safely add metadata to context
- Must validate key format (alphanumeric + underscore only)
- Should handle nil contexts gracefully
- Must be thread-safe
- Should maintain Claude Code context integrity

Generate:
- Unit tests with table-driven approach
- Edge case tests (nil, empty, invalid inputs)
- Concurrency tests with race detection
- Performance benchmarks
- Integration tests with Claude Code context"
```

#### Example 5: Code Quality Check
```markdown
PROMPT TO USE:
"As the Code Quality Enforcer, review this Go code for quality issues:

[paste code here]

Check for:
- Go idioms and best practices
- GoDoc documentation (100% coverage, 3+ sentences)
- Error handling patterns
- Concurrency safety
- Memory efficiency
- Security vulnerabilities
- Claude Code SDK usage patterns

Provide:
- Quality score (1-10)
- Specific issues found with severity
- Fix recommendations with code examples
- GoDoc templates for undocumented functions
- Performance optimization suggestions"
```

#### Example 6: Integration Testing
```markdown
PROMPT TO USE:
"As the Integration Orchestrator, design integration tests for:

Component: Context-extender CLI
Integration Points:
- Go context package integration
- Claude Code SDK integration
- File system persistence
- Cross-platform CLI behavior

Test Scenarios:
1. CLI → Go context → Claude Code flow
2. Context persistence and recovery
3. Concurrent CLI command execution
4. Cross-platform path handling

Generate:
- Integration test plan with scenarios
- Test environment setup requirements
- Expected results and validation criteria
- Performance benchmarks for integration
- Compatibility matrix for Claude Code versions"
```

### 📊 Review Phase Examples

#### Example 7: Progress Analysis (Data Request)
```markdown
PROMPT TO USE:
"As the Progress Tracker and Reporter, provide comprehensive metrics for 
Cycle 1 implementation phase:

Analyze:
- Story completion rate vs planned
- Velocity trends across days 5-15
- Blocker frequency and resolution time
- Test coverage progression
- Quality gate pass rates

Generate:
- Cycle summary dashboard
- Velocity chart with trend line
- Blocker analysis report
- Recommendations for Cycle 2
- Risk indicators for next cycle"
```

### Cross-Phase Examples

#### Example 8: Knowledge Documentation
```markdown
PROMPT TO USE:
"As the Knowledge Curator, document this architectural decision:

Decision: Using in-memory context store with optional persistence
Context: Need fast context operations with session recovery capability
Rationale: [provided rationale]
Trade-offs: Performance vs durability
Implementation: sync.Map with file-based snapshots

Create:
- ADR document with standard format
- Implementation pattern documentation
- Lessons learned entry
- Reusable code pattern
- Cross-reference to related decisions"
```

## SME Consultation Examples

#### Example 9: Technical SME Consultation
```markdown
PROMPT TO USE:
"As the Technical Governance SME, evaluate this technical decision:

We need to choose between:
1. Using sync.Map for concurrent context access
2. Using regular map with RWMutex
3. Using channel-based context updates

Consider:
- Read-heavy workload (90% reads, 10% writes)
- Need for iteration over all context entries
- Memory efficiency requirements
- Go 1.21 compatibility

Provide:
- Recommended approach with rationale
- Performance implications
- Code examples for recommended pattern
- Risk assessment
- Alternative solutions if requirements change"
```

#### Example 10: Quality SME Consultation
```markdown
PROMPT TO USE:
"As the Quality Governance SME, define quality standards for the 
context-extender CLI tool:

Project characteristics:
- Go CLI application
- Extends Claude Code functionality
- Critical for user workflows
- Must maintain high reliability

Define:
- Test coverage requirements
- Documentation standards
- Code review criteria
- Performance benchmarks
- Security requirements
- Release quality gates"
```

## Tips for Effective Subagent Usage

### 1. Provide Context
Always include:
- Project context (context-extender specifics)
- Technical constraints (Go, CLI, Claude Code)
- Current phase and objectives
- Dependencies and prerequisites

### 2. Be Specific
Instead of: "Generate tests"
Use: "Generate Go unit tests for ExtendContext function with concurrency testing"

### 3. Request Structured Output
Ask for:
- Specific formats (ADR, user story template)
- Measurable outcomes (metrics, scores)
- Actionable recommendations
- Clear next steps

### 4. Chain Subagents Appropriately
```
Architecture Discovery → generates constraints
    ↓
Story Refinement → uses constraints for stories
    ↓
Test Automation → uses stories for test generation
```

### 5. Validate with SMEs
For critical decisions:
1. Get subagent recommendation
2. Validate with relevant SME
3. Get User approval if strategic

These examples demonstrate practical, context-specific usage of subagents throughout the development cycle.