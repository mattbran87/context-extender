# Subagent & SME Activation Checklist

## Quick Phase Detection
**Current Cycle Day**: ___ of 17
**Current Phase**: [ ] Research (1-2) [ ] Planning (3-4) [ ] Implementation (5-15) [ ] Review (16-17)

## 🔬 Research Phase Team (Days 1-2)

### Subagent Activation
- [ ] **Architecture Discovery Specialist** - ACTIVATED
  - Input: Problem statement, constraints
  - Output: ADRs, architecture options
  - Mode: Lead explorer

### SME Availability
- [ ] **Technical Governance SME** - ON STANDBY
  - Triggers: Architecture decisions, technology evaluation
- [ ] **Go Language Specialist** - ON STANDBY
  - Triggers: Go-specific patterns, performance questions
- [ ] **CLI Development Specialist** - ON STANDBY
  - Triggers: CLI architecture, UX patterns
- [ ] **Claude Code Specialist** - ON STANDBY
  - Triggers: Extension patterns, integration approaches
- [ ] **Risk Governance SME** - ON STANDBY
  - Triggers: Risk identification, mitigation strategies

### Coordination
- [ ] Architecture Discovery leads exploration
- [ ] SMEs consulted as needed
- [ ] Knowledge Curator captures decisions (passive)

## 📋 Planning Phase Team (Days 3-4)

### Subagent Activation
- [ ] **Story Refinement Specialist** - ACTIVATED
  - Input: ADRs, constraints, user needs
  - Output: INVEST stories, technical details
  - Mode: Lead planner
  
- [ ] **Implementation Planning Orchestrator** - ACTIVATED
  - Input: Refined stories
  - Output: Estimates, schedule, resource plan
  - Mode: Support planner

### Subagent Deactivation
- [ ] **Architecture Discovery Specialist** - DEACTIVATED
  - Handoff: ADRs transferred to planning team

### SME Availability
- [ ] **Process Governance SME** - ACTIVE
  - Role: Planning process oversight
- [ ] **All Technical Specialists** - ON STANDBY
  - Role: Estimation support, technical details

### Coordination
- [ ] Story Refinement leads epic breakdown
- [ ] Implementation Planning handles estimation
- [ ] Sequential handoff between specialists

## 🔨 Implementation Phase Team (Days 5-15)

### Full Team Activation (5 Subagents)
- [ ] **Test Automation Specialist** - ACTIVATED
  - Input: Stories with acceptance criteria
  - Output: Automated tests, coverage reports
  - Mode: Continuous operation

- [ ] **Code Quality Enforcer** - ACTIVATED
  - Input: Code changes
  - Output: Quality validation, standards enforcement
  - Mode: Real-time monitoring

- [ ] **Integration Orchestrator** - ACTIVATED
  - Input: Components, dependencies
  - Output: Integration validation, compatibility checks
  - Mode: Continuous validation

- [ ] **Progress Tracker and Reporter** - ACTIVATED
  - Input: Development activities
  - Output: Metrics, reports, trends
  - Mode: Continuous monitoring

- [ ] **Knowledge Curator** - ACTIVATED
  - Input: Decisions, patterns, lessons
  - Output: Documentation, knowledge base
  - Mode: Continuous capture

### Planning Team Deactivation
- [ ] **Story Refinement Specialist** - DEACTIVATED
- [ ] **Implementation Planning Orchestrator** - DEACTIVATED
  - Handoff: Stories and schedule to implementation team

### SME Governance
- [ ] **Technical Governance SME** - ACTIVE
  - Role: Technical decision authority
- [ ] **Quality Governance SME** - ACTIVE
  - Role: Quality standards enforcement
- [ ] **Risk Governance SME** - MONITORING
  - Role: Risk tracking and escalation

### Coordination Matrix
```
Test Automation ↔ Code Quality (validation loop)
Code Quality ↔ Integration (standards alignment)
Integration ↔ Progress Tracker (status updates)
All → Knowledge Curator (documentation flow)
```

## 📊 Review Phase Team (Days 16-17)

### Implementation Team Deactivation
- [ ] **Test Automation Specialist** - DEACTIVATED
- [ ] **Code Quality Enforcer** - DEACTIVATED
- [ ] **Integration Orchestrator** - DEACTIVATED
- [ ] **Progress Tracker** - READ-ONLY MODE
- [ ] **Knowledge Curator** - READ-ONLY MODE

### SME-Led Review
- [ ] **Process Governance SME** - LEAD
  - Role: Retrospective facilitation
- [ ] **All SMEs** - ACTIVE
  - Role: Domain-specific review
- [ ] **User** - PRIMARY
  - Role: Stakeholder demos, feedback, approval

### Data Access
- [ ] Progress Tracker data available for metrics
- [ ] Knowledge Curator data available for lessons
- [ ] No new subagent activation

## 🔄 Phase Transition Protocols

### Research → Planning (Day 2 → 3)
- [ ] Complete Architecture Discovery outputs
- [ ] Archive research artifacts
- [ ] Activate Story Refinement Specialist
- [ ] Activate Implementation Planning Orchestrator
- [ ] Transfer ADRs and constraints

### Planning → Implementation (Day 4 → 5)
- [ ] Complete story refinement
- [ ] Finalize implementation schedule
- [ ] Activate all 5 implementation subagents
- [ ] Deactivate planning specialists
- [ ] Transfer stories and schedule

### Implementation → Review (Day 15 → 16)
- [ ] Complete all implementation tasks
- [ ] Compile implementation metrics
- [ ] Deactivate implementation subagents
- [ ] Switch to SME-led review mode
- [ ] Prepare demo materials

## 🚨 Activation Errors to Avoid

### ❌ Never Do This
- [ ] Run planning subagents during implementation
- [ ] Activate implementation team before planning complete
- [ ] Mix Research and Implementation subagents
- [ ] Skip deactivation protocols
- [ ] Activate Review subagents (it's SME-led)

### ✅ Always Do This
- [ ] Check current phase before activation
- [ ] Complete handoff before switching teams
- [ ] Maintain Knowledge Curator throughout
- [ ] Keep SMEs available for consultation
- [ ] Follow phase boundaries strictly

## 📝 Activation Log

### Research Phase
- **Activated**: Date: ___ Time: ___
- **Subagents**: Architecture Discovery
- **Issues**: ___

### Planning Phase
- **Activated**: Date: ___ Time: ___
- **Subagents**: Story Refinement, Implementation Planning
- **Issues**: ___

### Implementation Phase
- **Activated**: Date: ___ Time: ___
- **Subagents**: All 5 implementation specialists
- **Issues**: ___

### Review Phase
- **Activated**: Date: ___ Time: ___
- **Mode**: SME-led (no new subagents)
- **Issues**: ___

## Quick Reference Commands

### Research Phase Start
```
"Activating Research Team: Architecture Discovery Specialist as lead"
```

### Planning Phase Start
```
"Activating Planning Team: Story Refinement + Implementation Planning"
```

### Implementation Phase Start
```
"Activating Full Implementation Team: 5 specialists operational"
```

### Review Phase Start
```
"Switching to Review Mode: SME-led with Progress/Knowledge data support"
```

---
**Checklist Version**: 1.0
**Team Status**: [ ] Research [ ] Planning [ ] Implementation [ ] Review
**Active Subagents**: ___________________