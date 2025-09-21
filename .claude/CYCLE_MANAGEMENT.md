# Cycle Management Guidelines

## Overview

This document defines the standard procedures for managing development cycles in the Context-Extender project. These rules ensure consistent documentation and project organization.

## **MANDATORY RULE: Cycle Directory Creation**

### When Starting Any New Cycle

The **FIRST ACTION** when beginning a new cycle must be creating the proper directory structure:

```bash
# Manual creation (replace X with cycle number)
mkdir -p .claude/cycles/cycle-00X/{research,planning,implementation,review}

# OR use the automated script
.claude/scripts/start_new_cycle.sh [cycle_number]
```

### Why This Rule Exists

1. **Consistency**: Ensures all cycles follow the same documentation pattern
2. **Completeness**: Prevents missing documentation phases
3. **Organization**: Keeps project history properly structured
4. **Efficiency**: Templates and structure ready from day one

## Directory Structure Standard

```
.claude/cycles/cycle-00X/
├── README.md              # Cycle overview and checklist
├── research/               # Problem analysis phase
│   └── problem_definition.md
├── planning/               # Design and planning phase
│   └── implementation_plan.md
├── implementation/         # Development phase
│   └── daily_progress.md
└── review/                # Review and retrospective phase
    └── cycle_summary.md
```

## Automated Setup Script

### Usage

```bash
# Auto-detect next cycle number
.claude/scripts/start_new_cycle.sh

# Specify cycle number manually
.claude/scripts/start_new_cycle.sh 005
```

### What the Script Does

1. **Creates directory structure** with all required phases
2. **Generates template files** for each phase
3. **Creates cycle README** with checklist and guidelines
4. **Updates current_status** to reflect new cycle
5. **Provides next steps** for cycle planning

## Phase Documentation Requirements

### Research Phase
- **problem_definition.md**: Core problem and objectives
- **technical_feasibility.md**: Technical analysis (if needed)
- **risk_assessment.md**: Risk identification (if needed)

### Planning Phase
- **implementation_plan.md**: Development approach and schedule
- **user_stories.md**: Feature requirements (if applicable)
- **technical_design.md**: Architecture decisions (if needed)

### Implementation Phase
- **daily_progress.md**: Daily updates and progress tracking
- **code_quality_report.md**: Testing and quality metrics (if applicable)
- **weekly_retrospective.md**: Mid-cycle learning (for longer cycles)

### Review Phase
- **demo_feedback.md**: User feedback and validation
- **retrospective.md**: What worked, what didn't
- **cycle_summary.md**: Overall results and lessons learned

## Integration with Current Workflow

### Quick Fix Cycle (1-2 days)
1. ✅ Create cycle directory structure
2. Document problem in research/
3. Implement fix
4. Document results in review/

### Standard Cycle (5 days)
1. ✅ Create cycle directory structure
2. Day 1: Research and planning documentation
3. Days 2-4: Implementation with daily progress updates
4. Day 5: Review documentation and retrospective

## Cycle Numbering Convention

- **Format**: cycle-XXX (three digits, zero-padded)
- **Examples**: cycle-001, cycle-002, cycle-005, cycle-015
- **Auto-increment**: Script automatically detects next number

## Template Customization

Templates in the script can be customized based on:
- Cycle type (quick fix vs standard)
- Project phase (research vs feature development)
- Specific requirements (performance, security, etc.)

## Compliance and Quality

### Mandatory Elements
- ✅ Directory structure created before starting work
- ✅ README.md with cycle overview
- ✅ problem_definition.md with clear objectives
- ✅ cycle_summary.md with results and lessons

### Optional Elements
- Technical design documents (for complex changes)
- Detailed user stories (for feature work)
- Performance benchmarks (for optimization cycles)

## Integration with Existing Tools

### Version Control
- All cycle directories committed to git
- Document major decisions and changes
- Tag releases with cycle references

### Current Status Tracking
- Script automatically updates `.claude/current_status`
- Reflects current cycle and phase
- Maintains project history

### Compression Context
- Cycle documentation preserved in compression hooks
- Directory structure helps context restoration
- Consistent location for finding project state

## Benefits of This Approach

1. **Predictable Structure**: Always know where to find cycle information
2. **Complete Documentation**: No missing phases or deliverables
3. **Historical Record**: Clear project progression over time
4. **Efficient Setup**: Script eliminates manual setup errors
5. **Context Preservation**: Structured information survives compressions

## Enforcement

This rule should be followed for **ALL** new cycles, regardless of:
- Cycle duration (1 day or 5 days)
- Cycle type (bug fix or feature development)
- Urgency level (critical fix or enhancement)

**No exceptions** - the directory creation takes less than 30 seconds and provides immense organizational value.

## Example Usage

Starting Cycle 5:
```bash
# Create structure
.claude/scripts/start_new_cycle.sh

# Output:
# 🚀 Starting Cycle 005
# ✅ Cycle 005 structure created successfully!
# 📁 Directory: .claude/cycles/cycle-005
# 📝 Next: Define cycle objectives in research/problem_definition.md

# Begin work
cd .claude/cycles/cycle-005/research
# Edit problem_definition.md with cycle objectives
```

This systematic approach ensures Context-Extender maintains its excellent organizational standards while supporting the pragmatic development workflow that has proven 800% more efficient than rigid processes.