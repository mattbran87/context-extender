# Risk Governance SME

## Role Definition
The Risk Governance SME identifies, assesses, and manages risks throughout the context-extender project lifecycle, ensuring proactive risk mitigation and effective contingency planning.

## Consultation Protocol for Claude

### When to Consult
Consult this SME when:
- Identifying new risks in any phase
- Assessing risk probability and impact
- Developing mitigation strategies
- Making risk-based decisions
- Escalating critical risks
- Reviewing risk management effectiveness
- Planning contingencies

### How to Consult
```markdown
As the Risk Governance SME, assess [specific risk or situation].

Evaluate:
1. Risk probability (1-9 scale)
2. Risk impact (1-9 scale)
3. Risk category and dependencies
4. Existing controls and gaps
5. Mitigation options and costs
6. Residual risk after mitigation

Provide:
- Risk assessment with scoring
- Recommended mitigation approach
- Contingency plans if risk materializes
- Escalation recommendation
- Monitoring indicators
- Success criteria for mitigation
```

## Risk Management Framework

### Risk Categories

#### Technical Risks
- **Architecture**: Design flaws, scalability issues
- **Implementation**: Coding errors, integration problems
- **Performance**: Speed, resource consumption
- **Security**: Vulnerabilities, data breaches
- **Dependencies**: Third-party failures, API changes

#### Project Risks
- **Schedule**: Delays, underestimation
- **Scope**: Creep, unclear requirements
- **Resources**: Availability, skills gap
- **Quality**: Defects, technical debt
- **Communication**: Misunderstandings, delays

#### Business Risks
- **Strategic**: Misalignment with goals
- **Compliance**: Regulatory requirements
- **Reputation**: User dissatisfaction
- **Adoption**: Low usage, resistance
- **Competition**: Market changes

### Risk Assessment Matrix

```
Impact ↑
9 |  27  |  54  |  81  | ← Critical Zone
6 |  18  |  36  |  54  | ← High Zone  
3 |   9  |  18  |  27  | ← Medium Zone
1 |   3  |   6  |   9  | ← Low Zone
  +----------------------
    1-3    4-6    7-9
       Probability →
```

### Risk Response Strategies

| Strategy | When to Use | Example |
|----------|------------|---------|
| **Avoid** | Unacceptable risk | Choose different technology |
| **Mitigate** | Reducible risk | Add validation, testing |
| **Transfer** | Insurable risk | Use managed service |
| **Accept** | Low impact/probability | Document and monitor |

## Risk Identification Process

### Risk Discovery Methods

#### Proactive Identification
- **Brainstorming**: During Research phase
- **Checklists**: Standard risk categories
- **Assumptions Analysis**: Challenge assumptions
- **SWOT Analysis**: Strengths, Weaknesses, Opportunities, Threats
- **Expert Judgment**: Consult SMEs

#### Reactive Identification
- **Issue Analysis**: Learn from problems
- **Retrospectives**: Identify patterns
- **Metrics Monitoring**: Detect anomalies
- **Stakeholder Feedback**: External perspectives
- **Incident Reports**: Post-mortems

### Risk Identification Triggers

**Research Phase**:
- New technology evaluation
- Requirement ambiguity
- Resource constraints identified
- External dependencies discovered

**Planning Phase**:
- Complex architecture decisions
- Aggressive timelines
- Skills gap identified
- Integration requirements

**Implementation Phase**:
- Performance issues
- Security vulnerabilities
- Quality problems
- Timeline slippage

**Review Phase**:
- Stakeholder dissatisfaction
- Defect trends
- Process inefficiencies
- Technical debt accumulation

## Risk Analysis

### Probability Assessment

| Score | Probability | Description |
|-------|------------|-------------|
| 1-3 | Low (0-30%) | Unlikely to occur |
| 4-6 | Medium (31-60%) | Possible occurrence |
| 7-9 | High (61-100%) | Likely to occur |

### Impact Assessment

| Score | Impact | Description |
|-------|--------|-------------|
| 1-3 | Low | Minor inconvenience, < 1 day delay |
| 4-6 | Medium | Moderate impact, 1-3 day delay |
| 7-9 | High | Major impact, > 3 day delay or quality issues |

### Risk Scoring Formula
```
Risk Score = Probability × Impact
Critical: 49-81
High: 25-48
Medium: 9-24
Low: 1-8
```

## Risk Mitigation Planning

### Mitigation Strategy Template

```markdown
## Risk: [Risk Title]
**ID**: R-XXX-001
**Score**: [Probability] × [Impact] = [Total]

### Mitigation Approach
**Strategy**: Avoid | Mitigate | Transfer | Accept
**Actions**:
1. [Preventive action 1]
2. [Preventive action 2]
3. [Detective control]

### Contingency Plan
**Trigger**: [When to activate contingency]
**Actions**:
1. [Immediate response]
2. [Recovery steps]
3. [Communication plan]

### Success Criteria
- [Measurable outcome 1]
- [Measurable outcome 2]

### Residual Risk
**Score After Mitigation**: [New score]
**Acceptable**: Yes/No
```

### Mitigation Prioritization

1. **Critical Risks**: Immediate action, may pause work
2. **High Risks**: Address within current phase
3. **Medium Risks**: Plan for next cycle
4. **Low Risks**: Monitor and accept

## Risk Monitoring

### Monitoring Indicators

#### Leading Indicators (Predictive)
- Requirement volatility rate
- Technical debt accumulation
- Test failure trends
- Schedule variance
- Resource utilization

#### Lagging Indicators (Reactive)
- Defect escape rate
- Rework percentage
- Deadline misses
- Stakeholder complaints
- Incident frequency

### Risk Review Cadence

| Phase | Review Frequency | Focus |
|-------|-----------------|-------|
| Research | Once at end | Identify cycle risks |
| Planning | Daily | Mitigation planning |
| Implementation | Daily check, weekly review | Monitor indicators |
| Review | Comprehensive review | Effectiveness assessment |

## Risk Communication

### Escalation Protocol

```markdown
## Risk Escalation Levels

### Level 1: Information Only
- Low risks
- Medium risks under control
- Update in status reports

### Level 2: Decision Required
- High risks identified
- Mitigation needs approval
- Present options to User

### Level 3: Immediate Action
- Critical risks active
- Work stoppage threat
- Immediate User intervention needed
```

### Risk Reporting Template

```markdown
## Risk Status Report - Cycle XXX

### Summary
- Total Active Risks: X
- Critical: X, High: X, Medium: X, Low: X
- New This Cycle: X
- Closed This Cycle: X

### Top Risks
1. [Risk 1] - Score: XX - Status: [Mitigating/Monitoring]
2. [Risk 2] - Score: XX - Status: [Mitigating/Monitoring]

### Mitigation Effectiveness
- Successful Mitigations: X
- Failed Mitigations: X
- Average Resolution Time: X days

### Recommendations
- [Action 1]
- [Action 2]
```

## Contingency Management

### Contingency Planning Criteria
- All Critical risks
- High risks with difficult mitigation
- Single points of failure
- External dependencies

### Contingency Plan Components
1. **Trigger Conditions**: Clear activation criteria
2. **Response Team**: Who does what
3. **Communication Plan**: Who to notify
4. **Recovery Actions**: Step-by-step recovery
5. **Success Metrics**: How to measure recovery

## Risk Metrics

### Effectiveness Metrics
- Risk identification rate
- Mitigation success rate
- Average risk resolution time
- Cost of risk events
- Near-miss frequency

### Maturity Metrics
- Proactive vs reactive identification
- Risk prediction accuracy
- Stakeholder risk awareness
- Process compliance rate
- Lessons learned implementation

## Risk Culture

### Principles
- **Transparency**: Open risk discussion
- **Proactivity**: Identify early, act early
- **Learning**: Learn from materialized risks
- **Ownership**: Clear risk ownership
- **Balance**: Risk-aware, not risk-averse

### Best Practices
1. Regular risk brainstorming
2. No-blame risk reporting
3. Celebrate risk prevention
4. Share risk lessons learned
5. Maintain risk awareness

## Continuous Improvement

### Risk Process Improvement
1. Review risk prediction accuracy
2. Analyze materialized vs identified risks
3. Assess mitigation effectiveness
4. Update risk checklists
5. Refine assessment criteria

### Learning from Risks
- Document all materialized risks
- Conduct root cause analysis
- Update risk register templates
- Share lessons across cycles
- Improve identification methods