# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 🚨 CRITICAL: READ FIRST

**SINGLE SOURCE OF TRUTH**: `.claude/project_guide.md`

Before doing ANYTHING, read the complete project guide at `.claude/project_guide.md`. This contains:
- Current project status and restrictions
- Complete cyclical development framework
- Step-by-step workflow execution guide
- All rules, templates, and reference information

## Quick Reference

**Current Status**: v1.0.1 released, planning Cycle 5
**Project Phase**: Feature planning and workflow improvements
**Development Style**: Pragmatic, outcome-focused, direct implementation
**Workflow**: 5-day adaptive cycles (or faster for urgent fixes)

## Key Principles

1. ✅ **PRAGMATIC APPROACH** - Focus on delivering working solutions
2. 🚀 **RAPID ITERATION** - Test and ship quickly
3. 📦 **ZERO CGO** - Maintain Pure Go implementation
4. 🎯 **OUTCOME-FOCUSED** - Minimal process, maximum value

## Project Overview

**Context-Extender** is a CLI tool that captures Claude Code conversations automatically via hooks.

**Current State**: v1.0.1 production ready
**Core Achievement**: Zero CGO dependencies with Pure Go SQLite
**Key Features**: Conversation capture, GraphQL queries, context preservation

## Development Setup

**Go Version**: go1.25.1
**Build Command**: `CGO_ENABLED=0 "C:/Users/marko/sdk/go1.25.1/bin/go.exe" build -o context-extender.exe .`
**Database**: Pure Go SQLite using modernc.org/sqlite v1.39.0
**No CGO Required**: Builds work on all platforms without C compiler

## Actual Workflow (What Works)

### Quick Fix Cycle (1-2 days)
- Identify issue → Implement fix → Test → Release
- Used for critical bugs and small enhancements

### Standard Cycle (5 days)
- **Day 1**:
  1. ✅ **FIRST ACTION**: Create cycle directory structure
     ```bash
     mkdir -p .claude/cycles/cycle-00X/{research,planning,implementation,review}
     ```
  2. Identify objectives and plan
- Days 2-4: Implementation and testing
- Day 5: Validation and release

### Key Success Factors
- **Always start with directory creation** for proper documentation
- Direct implementation without over-planning
- Continuous testing during development
- Minimal but sufficient documentation
- Fast feedback and iteration

## Compression Context

If this conversation compresses, check:
- `.claude/COMPRESSION_CONTEXT.md` - Full context to restore
- `.claude/hooks/reinject_context.md` - Quick reinjection template
- `.claude/current_status` - Current project phase

## Cycle Management Rules

### **MANDATORY: New Cycle Setup**
When starting any new cycle, the **FIRST ACTION** must be:
```bash
# Replace X with cycle number (005, 006, etc.)
mkdir -p .claude/cycles/cycle-00X/{research,planning,implementation,review}
```

### **Directory Structure Standard**
```
.claude/cycles/cycle-00X/
├── research/          # Problem analysis, feasibility
├── planning/          # Design, user stories, estimates
├── implementation/    # Daily progress, code reports
└── review/           # Demo feedback, retrospectives
```

## Important Notes

**Don't follow** the rigid 17-day cycle in `.claude/project_guide.md` - it's outdated.
**Do follow** pragmatic approach that has delivered 800% efficiency.
**Always create** cycle directories first - this ensures proper documentation.
**Remember**: We've already completed 4 successful cycles and have a working v1.0.1.