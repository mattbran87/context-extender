# Master Checklist Framework for Context-Extender

## Overview
This master checklist coordinates all phase-specific, process, and team checklists for the context-extender project's 4-phase development cycle.

## Checklist Categories

### 📋 Phase Checklists (Primary)
- **Research Phase Checklist** - Days 1-2 entry/exit criteria
- **Planning Phase Checklist** - Days 3-4 entry/exit criteria  
- **Implementation Phase Checklist** - Days 5-15 entry/exit criteria
- **Review Phase Checklist** - Days 16-17 entry/exit criteria

### 📑 Document Checklists
- **Required Artifacts Checklist** - Documents required per phase
- **Handoff Documents Checklist** - Inter-phase artifact transfer
- **Knowledge Capture Checklist** - Continuous documentation

### 🤖 Team Activation Checklists
- **SME Consultation Checklist** - When to consult which SME
- **Subagent Activation Checklist** - Phase-specific team activation
- **Integration Points Checklist** - Subagent coordination

### ✅ Quality & Process Checklists  
- **Quality Gates Checklist** - Pass/fail criteria
- **Risk Management Checklist** - Continuous risk tracking
- **User Interaction Checklist** - 16 mandatory touchpoints

## Quick Reference: Which Checklist When?

| Current Activity | Use This Checklist | Location |
|-----------------|-------------------|----------|
| Starting new cycle | Research Phase Checklist | `./research_phase_checklist.md` |
| Phase transition | Phase Exit/Entry Checklists | `./[phase]_checklist.md` |
| Daily standup | Daily Activities Checklist | Within phase checklist |
| Complex decision | SME Consultation Checklist | `./sme_consultation_checklist.md` |
| Activating subagents | Subagent Activation Checklist | `./subagent_activation_checklist.md` |
| Document handoff | Required Artifacts Checklist | `./required_artifacts_checklist.md` |
| Quality review | Quality Gates Checklist | `./quality_gates_checklist.md` |
| Risk assessment | Risk Management Checklist | `./risk_management_checklist.md` |

## Checklist Hierarchy

```
Master Checklist (this document)
├── Phase Checklists (4)
│   ├── Entry Requirements
│   ├── Daily Activities
│   └── Exit Criteria
├── Process Checklists (3)
│   ├── Quality Gates
│   ├── Risk Management
│   └── User Interactions
├── Team Checklists (2)
│   ├── SME Consultation
│   └── Subagent Activation
└── Document Checklists (3)
    ├── Required Artifacts
    ├── Handoff Documents
    └── Knowledge Capture
```

## Critical Checkpoints (Cannot Proceed Without)

### 🔴 Phase Gates (User Approval Required)
- [ ] Research → Planning transition approval
- [ ] Planning → Implementation transition approval
- [ ] Implementation → Review transition approval
- [ ] Cycle completion approval

### 🔴 Risk Escalation Points
- [ ] Critical risks identified and escalated
- [ ] High risks have mitigation plans
- [ ] Risk register updated

### 🔴 Quality Gates
- [ ] Test coverage > 80%
- [ ] Zero critical security vulnerabilities
- [ ] GoDoc 100% coverage for public APIs
- [ ] Code review approval obtained

## Usage Instructions

### At Cycle Start
1. Open `research_phase_checklist.md`
2. Verify all pre-phase requirements
3. Activate Architecture Discovery Specialist
4. Begin daily activities checklist

### At Phase Transitions
1. Complete current phase exit checklist
2. Get user approval for transition
3. Archive phase deliverables
4. Open next phase entry checklist
5. Activate appropriate subagent team

### Daily Usage
1. Review daily activities for current phase
2. Check SME consultation triggers
3. Update progress on deliverables
4. Note any blockers or issues

### At Cycle End
1. Complete Review phase checklist
2. Compile all cycle metrics
3. Archive all checklists with cycle number
4. Prepare checklists for next cycle

## Checklist Compliance Tracking

### Mandatory Compliance (100% Required)
- User approval checkpoints
- Security review completion
- Risk escalation protocols
- Phase gate criteria

### Strong Compliance (>95% Expected)
- Quality gate standards
- Documentation requirements
- Process adherence
- Subagent coordination

### Recommended Compliance (>80% Target)
- Template usage
- Communication protocols
- Metrics collection
- Knowledge capture

## Success Metrics

### Checklist Effectiveness
- **Completion Rate**: >95% of checklist items completed
- **Gate Success**: >90% first-time phase gate passage
- **Process Compliance**: >95% adherence to mandatory items
- **User Satisfaction**: >4.5/5 rating on checklist usefulness

### Process Improvement
- **Cycle Time**: Reduction in phase duration over time
- **Quality Metrics**: Improvement in defect rates
- **Knowledge Capture**: Increase in documented decisions
- **Team Efficiency**: Reduction in coordination overhead

## Continuous Improvement

### After Each Cycle
- Review checklist effectiveness
- Identify missing or unnecessary items
- Update based on retrospective feedback
- Version control checklist changes

### Quarterly Review
- Analyze checklist usage patterns
- Consolidate or expand as needed
- Align with process improvements
- Update success metrics

## Quick Links to All Checklists

### Phase Checklists
- [Research Phase Checklist](./research_phase_checklist.md)
- [Planning Phase Checklist](./planning_phase_checklist.md)
- [Implementation Phase Checklist](./implementation_phase_checklist.md)
- [Review Phase Checklist](./review_phase_checklist.md)

### Process Checklists
- [Quality Gates Checklist](./quality_gates_checklist.md)
- [Risk Management Checklist](./risk_management_checklist.md)
- [User Interaction Checklist](./user_interaction_checklist.md)

### Team Checklists
- [SME Consultation Checklist](./sme_consultation_checklist.md)
- [Subagent Activation Checklist](./subagent_activation_checklist.md)

### Document Checklists
- [Required Artifacts Checklist](./required_artifacts_checklist.md)
- [Handoff Documents Checklist](./handoff_documents_checklist.md)
- [Knowledge Capture Checklist](./knowledge_capture_checklist.md)

---
**Version**: 1.0
**Last Updated**: [Current Date]
**Next Review**: After Cycle 1 completion