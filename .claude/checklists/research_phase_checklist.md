# Research Phase Checklist - Cycle [XXX]

## Phase Information
- **Phase**: Research
- **Duration**: 2 days
- **Start Date**: ____________________
- **Expected End Date**: ____________________
- **Current Day**: ___ of 2

## ✅ Pre-Phase Requirements (Entry Criteria)

### Mandatory (Must Complete)
- [ ] Previous cycle Review phase completed (if applicable)
- [ ] User approval to begin new cycle obtained
- [ ] Cycle objectives discussed with user
- [ ] Research directory created: `cycles/cycle-XXX/research/`
- [ ] Status file updated to Research phase

### Team Activation
- [ ] Architecture Discovery Specialist activated
- [ ] Technical Governance SME available for consultation
- [ ] Risk Governance SME available for consultation
- [ ] Knowledge Curator activated (passive mode)

## 📅 Daily Activities Checklist

### Day 1: Problem Definition & Exploration
**Morning (User Interaction Required)**
- [ ] Conduct problem definition session with user
- [ ] Document business objectives and success criteria
- [ ] Identify key stakeholders and their needs
- [ ] Define scope boundaries and constraints

**Architecture Exploration**
- [ ] Activate Architecture Discovery Specialist for initial exploration
- [ ] Research existing Claude Code extension patterns
- [ ] Investigate Go context manipulation approaches
- [ ] Document technical constraints and requirements

**Risk Identification**
- [ ] Create initial risk register
- [ ] Identify technical risks (Go, CLI, Claude Code integration)
- [ ] Identify project risks (timeline, resources, scope)
- [ ] Assess risk probability and impact

**📋 MANDATORY DOCUMENTATION**
- [ ] **IMMEDIATELY after problem session**: Create `cycles/cycle-XXX/research/problem_definition.md`
  - Include: Problem statement, use cases, success criteria, constraints
  - Template: Standardized problem definition format
- [ ] **IMMEDIATELY after SME consultation**: Create `cycles/cycle-XXX/research/architecture_decision_records.md`
  - Include: All ADRs from Architecture Discovery Specialist
  - Format: Standard ADR format (Status, Context, Decision, Consequences)

### Day 2: Feasibility & Technical Analysis
**Technical Feasibility (User Interaction Required)**
- [ ] Review technical exploration findings with user
- [ ] Validate architectural direction
- [ ] Confirm performance requirements
- [ ] Agree on technology constraints

**Competitive Analysis**
- [ ] Research similar solutions and patterns
- [ ] Document best practices from Go community
- [ ] Identify Claude Code extension precedents
- [ ] Compile anti-patterns to avoid

**User Research**
- [ ] Define user personas for context-extender
- [ ] Map user journeys for key scenarios
- [ ] Document pain points and opportunities
- [ ] Prioritize user needs with user input

**Feasibility Assessment**
- [ ] Technical feasibility confirmation
- [ ] Resource feasibility validation
- [ ] Timeline feasibility assessment
- [ ] Risk-adjusted feasibility score

**📋 MANDATORY DOCUMENTATION**
- [ ] **IMMEDIATELY after feasibility analysis**: Create `cycles/cycle-XXX/research/technical_feasibility_analysis.md`
  - Include: Challenge analysis, implementation patterns, risk assessment
  - Include: Performance targets, dependencies, testing strategy
- [ ] **IMMEDIATELY after risk analysis**: Create `cycles/cycle-XXX/research/risk_register.md`
  - Include: Risk matrix, mitigation strategies, monitoring plans
  - Format: Standardized risk register with scoring

## 🤖 SME Consultation Triggers

### Technical Governance SME
- [ ] Consult if: Major architecture patterns being evaluated
- [ ] Consult if: Technology stack decisions needed
- [ ] Consult if: Performance requirements unclear
- [ ] Consult if: Security concerns identified

### Go Language Specialist
- [ ] Consult if: Go-specific patterns needed for context manipulation
- [ ] Consult if: Concurrency design questions arise
- [ ] Consult if: Performance optimization strategies needed

### CLI Development Specialist
- [ ] Consult if: CLI command structure being designed
- [ ] Consult if: Cross-platform considerations identified

### Claude Code Specialist
- [ ] Consult if: Extension mechanism selection needed
- [ ] Consult if: SDK integration patterns being evaluated
- [ ] Consult if: Compatibility requirements unclear

### Risk Governance SME
- [ ] Consult if: Critical risks identified
- [ ] Consult if: Risk mitigation strategies needed
- [ ] Consult if: Risk escalation required

## 📋 Deliverable Requirements

### Required Documents
- [ ] Problem statement document (`problem_statement.md`)
- [ ] Technical feasibility report (`feasibility_report.md`)
- [ ] Risk assessment matrix (`risk_matrix.md`)
- [ ] User personas and journeys (`user_research.md`)
- [ ] Architecture exploration findings (`architecture_notes.md`)
- [ ] Initial ADRs for key decisions (minimum 2)

### Architecture Discovery Outputs
- [ ] Architecture option comparison matrix
- [ ] Technology evaluation results
- [ ] Performance analysis findings
- [ ] Integration approach recommendations

### Knowledge Capture
- [ ] All decisions documented with rationale
- [ ] Assumptions clearly stated
- [ ] Constraints identified and documented
- [ ] Lessons from research captured

## 🚪 Post-Phase Requirements (Exit Criteria)

### 📋 MANDATORY DOCUMENTATION CHECK
**🔴 NO PHASE TRANSITION WITHOUT ALL DOCUMENTS CREATED:**
- [ ] `cycles/cycle-XXX/research/problem_definition.md` - Complete and validated
- [ ] `cycles/cycle-XXX/research/architecture_decision_records.md` - All ADRs documented
- [ ] `cycles/cycle-XXX/research/technical_feasibility_analysis.md` - Confidence rating included
- [ ] `cycles/cycle-XXX/research/risk_register.md` - Mitigation plans complete
- [ ] All documents cross-referenced and internally consistent

### Mandatory (Must Complete for Phase Transition)
- [ ] Problem statement finalized and documented
- [ ] Technical feasibility confirmed
- [ ] At least 2 ADRs completed and reviewed
- [ ] Risk assessment with mitigation strategies
- [ ] User personas and journeys approved
- [ ] **USER APPROVAL: Ready to proceed to Planning phase**

### Quality Checks
- [ ] All research documented in standard format
- [ ] Key decisions have clear rationale
- [ ] Risks are quantified (probability × impact)
- [ ] Feasibility has evidence backing

### Handoff to Planning Phase
- [ ] ADRs ready for story refinement input
- [ ] Constraints documented for estimation
- [ ] Risk register available for planning
- [ ] User research ready for acceptance criteria
- [ ] Technical requirements clear for design

## 🔄 Phase Transition Protocol

### Deactivation
- [ ] Architecture Discovery Specialist findings archived
- [ ] Research artifacts organized in cycle folder
- [ ] Knowledge Curator updated with research insights

### Activation for Next Phase
- [ ] Story Refinement Specialist ready to activate
- [ ] Implementation Planning Orchestrator ready to activate
- [ ] Research outputs formatted for planning input

## 📝 Notes and Issues

### Blockers
- Issue: ____________________
  - Resolution: ____________________
- Issue: ____________________
  - Resolution: ____________________

### Key Decisions
- Decision: ____________________
  - Rationale: ____________________
- Decision: ____________________
  - Rationale: ____________________

### Parking Lot (Future Consideration)
- [ ] ____________________
- [ ] ____________________

## ✍️ Sign-off

### Phase Completion
- **Claude Confirmation**: [ ] All research objectives met
- **Deliverables Complete**: [ ] All required documents ready
- **Quality Validated**: [ ] Research meets quality standards

### User Approval
- **User Review**: [ ] Research findings reviewed
- **User Approval**: [ ] Approved to proceed to Planning
- **Date**: ____________________
- **Notes**: ____________________

---
**Checklist Version**: 1.0
**Phase Status**: [ ] Not Started [ ] In Progress [ ] Complete
**Next Phase**: Planning (Days 3-4)