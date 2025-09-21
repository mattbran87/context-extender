#!/bin/bash
# Hook to preserve context when conversation compresses

# Capture current context
cat > .claude/last_compression.md << 'EOF'
# Context Preserved at Compression

## Quick Reinjection
Working on Context-Extender v1.0.1+. Completed Cycles 1-4 successfully.
- Zero CGO achieved with modernc.org/sqlite
- Prefer pragmatic 5-day cycles over rigid 17-day documentation
- Direct implementation approach, minimal docs
- Build: CGO_ENABLED=0 go build

## Current State
- Phase: Planning Cycle 5
- Last: Fixed critical issues in 1-day sprint
- Focus: Workflow improvements and compression handling

## Key Decisions
- Database: Pure Go SQLite (modernc.org/sqlite v1.39.0)
- Architecture: Clean backend abstraction
- Workflow: Pragmatic, outcome-focused

## Don't
- Follow rigid 17-day cycles
- Add encryption/auth
- Over-document
- Wait for formal approvals
EOF

echo "Context preserved to .claude/last_compression.md"