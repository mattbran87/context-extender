# Risk Management Process

## Overview
This document defines the risk management process for the context-extender project, integrated into the 4-phase cyclical development framework.

## Risk Categories

### Technical Risks
- Architecture and design flaws
- Technology stack limitations
- Integration challenges
- Performance bottlenecks
- Security vulnerabilities

### Project Risks
- Scope creep
- Timeline delays
- Resource unavailability
- Requirement changes
- Quality issues

### External Risks
- Dependency failures
- Third-party service outages
- Go ecosystem changes
- Breaking changes in libraries

## Risk Management Integration by Phase

### Research Phase (Days 1-3)
**Step 1.6: Risk Assessment** (Enhanced)
- Identify new risks related to cycle objectives
- Review and update existing risk register
- Assess risk probability and impact
- Document in `research/risk-assessment.md`

### Planning Phase (Days 1-4)
**Day 2: Risk Mitigation Planning**
- Develop mitigation strategies for high-priority risks
- Assign risk owners (User or Claude)
- Define risk triggers and indicators
- Update risk register with mitigation plans

### Implementation Phase (Days 1-7)
**Daily Risk Monitoring**
- Check risk indicators during daily activities
- Log any triggered risks in implementation notes
- Escalate critical risks immediately to user
- Update risk status in daily commits

### Review Phase (Days 1-3)
**Day 2: Risk Review**
- Analyze risks that materialized during cycle
- Evaluate effectiveness of mitigation strategies
- Update risk register with lessons learned
- Archive cycle-specific risks

## Risk Register Template

```markdown
# Risk Register - Cycle XXX

## Active Risks

### Risk ID: R-XXX-001
**Risk Title**: [Short descriptive title]
**Category**: Technical | Project | External
**Status**: Open | Mitigating | Closed
**Probability**: Low (1-3) | Medium (4-6) | High (7-9)
**Impact**: Low (1-3) | Medium (4-6) | High (7-9)
**Risk Score**: [Probability × Impact]
**Description**: [Detailed description of the risk]
**Trigger Indicators**: 
- [What signs indicate this risk is occurring]
**Mitigation Strategy**: 
- [Preventive actions to reduce probability]
- [Contingency actions if risk occurs]
**Owner**: User | Claude
**Date Identified**: YYYY-MM-DD
**Last Updated**: YYYY-MM-DD
**Notes**: [Any additional context or updates]

---

### Risk ID: R-XXX-002
[Repeat template for each risk]
```

## Risk Assessment Matrix

```
Impact ↑
High    | Medium Risk | High Risk   | Critical Risk
Medium  | Low Risk    | Medium Risk | High Risk
Low     | Low Risk    | Low Risk    | Medium Risk
        +---------------------------------
          Low          Medium        High
                    Probability →
```

### Risk Priority Levels
- **Critical Risk** (Score 49-81): Immediate mitigation required, may block phase completion
- **High Risk** (Score 25-48): Active mitigation required, monitor daily
- **Medium Risk** (Score 9-24): Monitor weekly, have contingency plan ready
- **Low Risk** (Score 1-8): Accept and monitor, document if occurs

## Risk Escalation Protocol

### Critical Risk Escalation
```
1. Claude identifies critical risk
2. Immediately notify user: "🔴 CRITICAL RISK: [Title]"
3. Present mitigation options
4. Wait for user decision before proceeding
5. Document decision and outcome
```

### High Risk Escalation
```
1. Claude identifies high risk during daily activities
2. Complete current task
3. Notify user at next interaction point
4. Propose mitigation approach
5. Implement approved mitigation
```

## Risk Documentation Structure

```
cycles/
└── cycle-XXX/
    └── risk-management/
        ├── risk-register.md         # Current cycle risk register
        ├── risk-assessments/         # Detailed risk analyses
        │   ├── technical-risks.md
        │   ├── project-risks.md
        │   └── external-risks.md
        └── mitigation-tracking.md   # Mitigation effectiveness tracking
```

## Risk Metrics

### Cycle-Level Metrics
- Number of risks identified
- Number of risks materialized
- Mitigation effectiveness rate
- Average risk resolution time

### Project-Level Metrics
- Risk trends across cycles
- Most common risk categories
- Mitigation strategy success rate
- Cost of risk impacts (time/rework)

## Risk Review Questions

### During Research Phase
1. What risks could prevent achieving cycle objectives?
2. Have any new technical risks emerged?
3. Are there external dependencies that could fail?

### During Planning Phase
1. Do our plans adequately address identified risks?
2. Are mitigation strategies realistic and achievable?
3. Have we allocated enough buffer for risk impacts?

### During Implementation Phase
1. Are any risk indicators showing warning signs?
2. Have new risks emerged during development?
3. Are mitigation strategies working as expected?

### During Review Phase
1. Which risks actually occurred?
2. How effective were our mitigation strategies?
3. What new risks should we track for next cycle?

## Integration with Existing Workflow

This risk management process integrates seamlessly with the existing 4-phase cycle:
- Uses existing phase structure and timing
- Adds minimal overhead (15-30 minutes per phase)
- Leverages existing documentation structure
- Maintains user approval gates for critical decisions