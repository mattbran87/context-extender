#!/bin/bash
# Script to start a new development cycle with proper directory structure

# Get the next cycle number
CYCLE_NUM=$1

if [ -z "$CYCLE_NUM" ]; then
    # Auto-detect next cycle number
    LAST_CYCLE=$(ls -1 .claude/cycles/ | grep -E '^cycle-[0-9]+$' | sort -V | tail -1)
    if [ -z "$LAST_CYCLE" ]; then
        CYCLE_NUM="001"
    else
        LAST_NUM=$(echo $LAST_CYCLE | sed 's/cycle-//')
        NEXT_NUM=$((10#$LAST_NUM + 1))
        CYCLE_NUM=$(printf "%03d" $NEXT_NUM)
    fi
fi

CYCLE_DIR=".claude/cycles/cycle-$CYCLE_NUM"

echo "🚀 Starting Cycle $CYCLE_NUM"
echo "Creating directory structure: $CYCLE_DIR"

# Create the cycle directory structure
mkdir -p "$CYCLE_DIR/research"
mkdir -p "$CYCLE_DIR/planning"
mkdir -p "$CYCLE_DIR/implementation"
mkdir -p "$CYCLE_DIR/review"

# Create initial README for the cycle
cat > "$CYCLE_DIR/README.md" << EOF
# Cycle $CYCLE_NUM

**Start Date**: $(date +"%Y-%m-%d")
**Objective**: TBD
**Duration**: TBD
**Status**: Planning

## Directory Structure

- \`research/\` - Problem analysis, feasibility studies, technical research
- \`planning/\` - Design documents, user stories, implementation plans
- \`implementation/\` - Daily progress reports, code quality reports
- \`review/\` - Demo feedback, retrospectives, cycle summary

## Phase Checklist

### Research Phase
- [ ] Problem definition
- [ ] Technical feasibility analysis
- [ ] Risk assessment
- [ ] Requirements gathering

### Planning Phase
- [ ] User stories created
- [ ] Technical design completed
- [ ] Implementation plan ready
- [ ] Success criteria defined

### Implementation Phase
- [ ] Development completed
- [ ] Tests passing
- [ ] Code reviewed
- [ ] Documentation updated

### Review Phase
- [ ] Demo conducted
- [ ] Feedback collected
- [ ] Retrospective completed
- [ ] Cycle summary created

## Notes

- Follow pragmatic 5-day cycle approach
- Focus on working solutions over documentation
- Update this README as cycle progresses
EOF

# Create template files for each phase
cat > "$CYCLE_DIR/research/problem_definition.md" << EOF
# Problem Definition - Cycle $CYCLE_NUM

## Problem Statement
TBD

## Objectives
TBD

## Success Criteria
TBD

## Constraints
- Maintain zero CGO dependencies
- Preserve Pure Go SQLite implementation
- Follow pragmatic development approach

## Next Steps
TBD
EOF

cat > "$CYCLE_DIR/planning/implementation_plan.md" << EOF
# Implementation Plan - Cycle $CYCLE_NUM

## User Stories
TBD

## Technical Approach
TBD

## Daily Schedule
- Day 1: Research and planning
- Days 2-4: Implementation
- Day 5: Testing and review

## Risk Mitigation
TBD
EOF

cat > "$CYCLE_DIR/implementation/daily_progress.md" << EOF
# Daily Progress - Cycle $CYCLE_NUM

## Day 1 - $(date +"%Y-%m-%d")
**Focus**: Setup and initial research
**Progress**: Cycle directory structure created
**Next**: Define objectives and begin implementation

---

<!-- Add daily entries as cycle progresses -->
EOF

cat > "$CYCLE_DIR/review/cycle_summary.md" << EOF
# Cycle $CYCLE_NUM Summary

**Duration**: TBD
**Objective**: TBD
**Status**: TBD

## Results
TBD

## Lessons Learned
TBD

## Next Cycle Recommendations
TBD
EOF

echo "✅ Cycle $CYCLE_NUM structure created successfully!"
echo "📁 Directory: $CYCLE_DIR"
echo "📝 Next: Define cycle objectives in research/problem_definition.md"

# Update current_status file
echo "📊 Updating current status..."
sed -i "s/Cycle 5 Planning/Cycle $CYCLE_NUM Planning/" .claude/current_status
sed -i "s/Ready for Cycle 5/Ready for Cycle $CYCLE_NUM/" .claude/current_status

echo "🎯 Ready to begin Cycle $CYCLE_NUM development!"