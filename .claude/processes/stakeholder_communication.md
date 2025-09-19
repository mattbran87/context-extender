# Stakeholder Communication Matrix

## Overview
This document defines the stakeholder communication strategy for the context-extender project, adaptable from a single stakeholder (current state) to multiple stakeholders (future state).

## Stakeholder Categories

### Current State (Single Stakeholder)
- **Primary Stakeholder**: User (Product Owner, Developer, Customer)
  - All roles consolidated in single person
  - Direct communication through Claude interface
  - Immediate feedback and decision-making

### Future State (Multiple Stakeholders)
- **Product Owner**: Vision, priorities, acceptance
- **Development Team**: Technical implementation
- **End Users**: Feature consumers
- **Technical Reviewers**: Code quality, architecture
- **Operations Team**: Deployment, monitoring

## Communication Matrix

### Phase-Specific Communication Points

| Phase | Stakeholder | Communication Type | Frequency | Method | Content |
|-------|------------|-------------------|-----------|---------|---------|
| **Research** | User/Product Owner | Problem Definition | Day 1 | Interactive Discussion | Problem scope, objectives, constraints |
| **Research** | User/Product Owner | Requirements Gathering | Day 2 | Interactive Discussion | User needs, acceptance criteria |
| **Research** | User/Product Owner | Feasibility Review | Day 3 | Document Review | Technical feasibility, risks |
| **Planning** | User/Technical Lead | Design Review | Day 2 | Document + Discussion | Architecture, technical approach |
| **Planning** | User/Product Owner | Story Refinement | Day 1 | Interactive Review | User stories, priorities |
| **Planning** | User/All Roles | Sprint Planning | Day 4 | Approval Gate | Sprint commitment, resource allocation |
| **Implementation** | User/Developer | Progress Update | Daily (async) | Status File | Story completion, blockers |
| **Implementation** | User/All Roles | Mid-Sprint Check | Day 4 | Optional Sync | Progress review, scope adjustment |
| **Review** | User/Product Owner | Demo | Day 1 | Interactive Demo | Feature demonstration |
| **Review** | User/All Roles | Retrospective | Day 3 | Interactive Discussion | Process improvements |

## Communication Protocols

### Synchronous Communication (Current State)
**When**: User interaction points as defined in workflow
**Format**: 
```
Claude: "[Phase] [Topic]: [Specific Question/Update]"
User: "[Decision/Feedback]"
Claude: "[Confirmation and Next Steps]"
```

### Asynchronous Communication (Future Enhancement)
**Status Updates**: Via `.claude/current_status` file
**Daily Summaries**: In `cycles/cycle-XXX/daily-updates/day-X.md`
**Blockers**: Immediate notification via agreed channel

## Stakeholder Engagement Levels

### RACI Matrix Template

| Activity | User (Current) | Product Owner (Future) | Dev Team (Future) | Reviewers (Future) |
|----------|---------------|----------------------|------------------|-------------------|
| Problem Definition | R,A,C,I | A,C | C,I | I |
| Technical Design | R,A,C,I | C | R,A | C,I |
| Story Creation | R,A,C,I | R,A | C | I |
| Implementation | R,A,C,I | I | R,A | C |
| Code Review | R,A,C,I | I | R | A |
| Acceptance Testing | R,A,C,I | A | R | C |
| Deployment Decision | R,A,C,I | A | R,C | C |

**Legend**:
- R = Responsible (does the work)
- A = Accountable (final approval)
- C = Consulted (provides input)
- I = Informed (kept updated)

## Communication Templates

### Phase Transition Request
```markdown
## Phase Transition: [Current] → [Next]

**Cycle**: cycle-XXX
**Date**: YYYY-MM-DD

### Completed Deliverables
- ✓ [Deliverable 1]
- ✓ [Deliverable 2]
- ✓ [Deliverable 3]

### Key Outcomes
[Summary of phase achievements]

### Risks/Issues for Next Phase
[Any concerns to address]

**Request**: Approve transition to [Next] phase?
```

### Daily Status Update
```markdown
## Daily Status - [Phase] Day [X]

**Date**: YYYY-MM-DD
**Progress**: X of Y stories complete

### Completed Today
- [Achievement 1]
- [Achievement 2]

### In Progress
- [Current work]

### Blockers
- [Any blockers]

### Tomorrow's Plan
- [Next priorities]
```

### Risk Escalation
```markdown
## 🔴 Risk Escalation

**Risk**: [Title]
**Severity**: Critical | High
**Impact**: [Description]

### Mitigation Options
1. [Option 1]
2. [Option 2]

**Recommendation**: [Claude's recommendation]
**Decision Required**: [What user needs to decide]
```

## Feedback Collection Methods

### Research Phase Feedback
- **Method**: Interactive discussion
- **Topics**: Problem validation, scope agreement
- **Documentation**: `research/stakeholder-feedback.md`

### Planning Phase Feedback
- **Method**: Document review + discussion
- **Topics**: Technical approach, story priorities
- **Documentation**: `planning/stakeholder-input.md`

### Implementation Phase Feedback
- **Method**: Async updates + optional check-ins
- **Topics**: Progress, blockers, scope changes
- **Documentation**: `implementation/daily-updates/`

### Review Phase Feedback
- **Method**: Demo + structured feedback session
- **Topics**: Feature acceptance, improvements
- **Documentation**: `review/stakeholder-feedback.md`

## Stakeholder Satisfaction Metrics

### Measurement Points
- End of each phase (quick pulse check)
- End of each cycle (comprehensive review)
- Quarterly trend analysis

### Satisfaction Categories
1. **Communication Effectiveness** (1-5 scale)
   - Clarity of updates
   - Timeliness of information
   - Relevance of content

2. **Engagement Quality** (1-5 scale)
   - Opportunity for input
   - Response to feedback
   - Decision involvement

3. **Delivery Satisfaction** (1-5 scale)
   - Met expectations
   - Quality of deliverables
   - Timeline adherence

### Feedback Questions Template
```markdown
## Cycle XXX Stakeholder Feedback

**Communication** (1-5): [ ]
- Were updates clear and timely?
- Did you have enough information for decisions?

**Engagement** (1-5): [ ]
- Were you involved at the right times?
- Was your feedback incorporated?

**Delivery** (1-5): [ ]
- Did deliverables meet expectations?
- Are you satisfied with the cycle outcome?

**Comments**: 
[Open feedback]

**Improvements for Next Cycle**:
[Suggestions]
```

## Scaling Communication Strategy

### Current State (1 Stakeholder)
- All communication through Claude interface
- Immediate feedback loops
- Single approval authority
- Unified vision and priorities

### Transition State (2-3 Stakeholders)
- Add email/Slack notifications for key events
- Introduce stakeholder-specific views
- Maintain single approval authority
- Begin role separation

### Future State (4+ Stakeholders)
- Implement communication platform (Slack/Teams)
- Create stakeholder dashboard
- Establish approval hierarchy
- Specialized communication by role

## Communication Effectiveness Monitoring

### Key Performance Indicators
- Response time to stakeholder queries
- Stakeholder satisfaction scores
- Number of miscommunications per cycle
- Feedback incorporation rate

### Continuous Improvement
- Review communication effectiveness each retrospective
- Adjust frequency based on stakeholder needs
- Optimize documentation for clarity
- Streamline approval processes

## Integration with Existing Workflow

This communication matrix integrates with the existing 4-phase cycle:
- Leverages existing interaction points
- Adds structure to current ad-hoc communication
- Scales from single to multiple stakeholders
- Maintains user approval gates