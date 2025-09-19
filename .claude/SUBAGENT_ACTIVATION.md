# 🚀 SUBAGENT ACTIVATION QUICK REFERENCE

## Current Phase Detection
**Check cycle day to determine phase:**
- Days 1-2: RESEARCH
- Days 3-4: PLANNING  
- Days 5-15: IMPLEMENTATION
- Days 16-17: REVIEW

## Phase-Specific Subagent Teams

### 🔬 RESEARCH (Days 1-2)
```bash
ACTIVATE:
✓ Architecture Discovery Specialist (lead)
✓ Technical SMEs (advisory)

OBJECTIVES:
- Explore architecture options
- Generate ADRs
- Assess feasibility
```

### 📋 PLANNING (Days 3-4)
```bash
ACTIVATE:
✓ Story Refinement Specialist (lead)
✓ Implementation Planning Orchestrator
✓ Knowledge Curator (passive)

OBJECTIVES:
- Break down epics
- Estimate stories
- Create schedule
```

### 🔨 IMPLEMENTATION (Days 5-15)
```bash
ACTIVATE ALL:
✓ Test Automation Specialist
✓ Code Quality Enforcer
✓ Integration Orchestrator
✓ Progress Tracker and Reporter
✓ Knowledge Curator

OBJECTIVES:
- Execute stories
- Maintain quality
- Track progress
```

### 📊 REVIEW (Days 16-17)
```bash
SME-LED (No new subagent activation):
✓ Progress Tracker (data only)
✓ Knowledge Curator (data only)
✓ Process Governance SME (lead)

OBJECTIVES:
- Stakeholder demos
- Retrospective
- Next cycle prep
```

## Decision Flow
```
1. What day of cycle? → Determines phase
2. What phase? → Determines team
3. What team? → Determines subagents
4. Activate subagents → Execute phase
```

## Quick Commands

### Start of Cycle (Day 1)
```
"Starting Research Phase - Activating Architecture Discovery Specialist"
```

### Planning Transition (Day 3)
```
"Starting Planning Phase - Activating Story Refinement and Implementation Planning"
```

### Implementation Start (Day 5)
```
"Starting Implementation Phase - Activating full implementation team (5 subagents)"
```

### Review Start (Day 16)
```
"Starting Review Phase - SME-led with Progress/Knowledge data support"
```

## Remember
- **Clear boundaries**: Don't mix phase teams
- **Clean handoffs**: Transfer artifacts between phases
- **SME oversight**: Always maintain governance
- **Knowledge capture**: Continuous throughout all phases

**For detailed protocols**: See `.claude/subagents/subagent_teams.md`