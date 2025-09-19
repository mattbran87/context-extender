# Process Governance SME

## Role Definition
The Process Governance SME maintains process integrity, ensures compliance with established workflows, and drives continuous process improvement throughout the context-extender project lifecycle.

## Consultation Protocol for Claude

### When to Consult
Consult this SME when:
- Deviating from standard processes
- Optimizing workflow efficiency
- Defining new processes
- Measuring process effectiveness
- Resolving process conflicts
- Implementing process improvements
- Conducting process audits

### How to Consult
```markdown
As the Process Governance SME, evaluate [specific process issue/improvement].

Analyze:
1. Current process state and gaps
2. Impact on project efficiency
3. Compliance requirements
4. Resource implications
5. Change management needs
6. Measurable benefits

Provide:
- Process assessment and recommendations
- Implementation approach and timeline
- Change impact analysis
- Success metrics
- Risk assessment
- Training/communication needs
```

## Process Framework

### Core Processes

#### Phase Transition Process
- **Purpose**: Ensure quality gates are met
- **Owner**: Claude (with User approval)
- **Frequency**: End of each phase
- **Compliance**: Mandatory

#### Risk Management Process
- **Purpose**: Identify and mitigate risks
- **Owner**: Claude
- **Frequency**: Continuous
- **Compliance**: Mandatory

#### Quality Assurance Process
- **Purpose**: Maintain quality standards
- **Owner**: Claude
- **Frequency**: Continuous
- **Compliance**: Mandatory

#### Stakeholder Communication Process
- **Purpose**: Ensure effective communication
- **Owner**: Claude
- **Frequency**: As defined in matrix
- **Compliance**: Mandatory

#### Continuous Improvement Process
- **Purpose**: Enhance processes
- **Owner**: Claude
- **Frequency**: Each retrospective
- **Compliance**: Recommended

### Process Maturity Model

| Level | Description | Characteristics | Target Timeline |
|-------|------------|-----------------|-----------------|
| **1 - Initial** | Ad-hoc processes | Inconsistent, reactive | Starting point |
| **2 - Defined** | Documented processes | Consistent, proactive | Current state |
| **3 - Managed** | Measured processes | Metrics-driven, optimized | 3-6 months |
| **4 - Optimized** | Continuous improvement | Self-improving, predictive | 12+ months |

## Process Compliance

### Compliance Framework

#### Mandatory Processes
Must be followed without exception:
- Phase transition gates
- User approval requirements
- Risk escalation protocol
- Security review process
- Documentation standards

#### Recommended Processes
Should be followed unless justified:
- Code review practices
- Testing strategies
- Meeting cadences
- Metrics collection
- Retrospective format

#### Optional Processes
Can be adapted as needed:
- Communication methods
- Tool selection
- Documentation formats
- Review techniques

### Deviation Management

#### Deviation Request Process
```markdown
## Process Deviation Request

**Process**: [Process name]
**Requested By**: Claude/User
**Date**: YYYY-MM-DD

### Deviation Details
**Standard Process**: [What should be done]
**Requested Deviation**: [What is proposed]
**Duration**: [Temporary/Permanent]

### Justification
[Why deviation is needed]

### Impact Analysis
- **Quality Impact**: [Assessment]
- **Timeline Impact**: [Assessment]
- **Risk Impact**: [Assessment]

### Mitigation
[How to minimize negative impacts]

### Approval
**Status**: Approved/Rejected
**Approver**: User
**Conditions**: [Any conditions]
```

## Process Optimization

### Optimization Methodology

#### PDCA Cycle (Plan-Do-Check-Act)
1. **Plan**: Identify improvement opportunity
2. **Do**: Implement change in pilot
3. **Check**: Measure effectiveness
4. **Act**: Standardize if successful

### Process Metrics

#### Efficiency Metrics
- **Cycle Time**: Time from start to finish
- **Process Time**: Actual work time
- **Wait Time**: Delays and handoffs
- **Rework Rate**: Work that needs redoing
- **First Pass Yield**: Success on first attempt

#### Effectiveness Metrics
- **Goal Achievement**: Process outcomes met
- **Compliance Rate**: Adherence to process
- **Stakeholder Satisfaction**: Process usability
- **Error Rate**: Process-related errors
- **Value Delivery**: Business value created

### Improvement Identification

#### Sources of Improvement
- **Retrospectives**: Team feedback
- **Metrics Analysis**: Data-driven insights
- **Stakeholder Feedback**: External perspective
- **Benchmarking**: Industry best practices
- **Incident Analysis**: Learning from failures

#### Improvement Prioritization Matrix

```
Impact ↑
High    | Quick Win | Major Project
Medium  | Fill In   | Good Option
Low     | Don't Do  | Low Priority
        +------------------------
          Low         High
             Effort →
```

## Process Documentation

### Documentation Standards

#### Process Document Template
```markdown
# [Process Name]

## Purpose
[Why this process exists]

## Scope
[What this process covers]

## Process Owner
[Who is responsible]

## Process Steps
1. [Step 1]
   - Input: [What's needed]
   - Activity: [What to do]
   - Output: [What's produced]
2. [Step 2]
   ...

## Roles and Responsibilities
- **Role 1**: [Responsibilities]
- **Role 2**: [Responsibilities]

## Process Metrics
- [Metric 1]: [Target]
- [Metric 2]: [Target]

## Related Documents
- [Document 1]
- [Document 2]

## Revision History
| Version | Date | Changes | Author |
|---------|------|---------|--------|
| 1.0 | YYYY-MM-DD | Initial | Claude |
```

### Process Repository Structure
```
.claude/
├── processes/
│   ├── core/
│   │   ├── phase_transition.md
│   │   ├── risk_management.md
│   │   ├── quality_assurance.md
│   │   └── stakeholder_communication.md
│   ├── supporting/
│   │   ├── code_review.md
│   │   ├── testing_strategy.md
│   │   └── documentation.md
│   └── improvement/
│       ├── retrospective.md
│       └── metrics_collection.md
```

## Process Governance

### Governance Structure

#### Process Governance Board (Future)
- Reviews process performance
- Approves major changes
- Resolves process conflicts
- Sets process standards

#### Process Owners
- **Phase Processes**: Claude
- **Technical Processes**: Claude
- **Communication**: User/Claude
- **Governance**: User

### Change Control

#### Process Change Levels

| Level | Impact | Approval | Implementation |
|-------|--------|----------|----------------|
| **Minor** | Single step, no risk | Claude | Immediate |
| **Moderate** | Multiple steps, low risk | User review | Next cycle |
| **Major** | Cross-process, medium risk | User approval | Pilot first |
| **Critical** | High risk, major change | User decision | Phased rollout |

### Process Audit

#### Audit Schedule
- **Self-audit**: Each cycle during Review phase
- **Deep audit**: Quarterly
- **External audit**: Annually (if applicable)

#### Audit Checklist
- [ ] Process documentation current
- [ ] Process being followed
- [ ] Metrics being collected
- [ ] Improvements identified
- [ ] Training adequate
- [ ] Tools appropriate

## Process Training

### Training Requirements

#### New Process Introduction
1. Document process clearly
2. Communicate changes
3. Provide examples
4. Pilot with support
5. Gather feedback
6. Refine and standardize

#### Ongoing Training
- Process refreshers in retrospectives
- Updates on process changes
- Lessons learned sharing
- Best practice documentation

## Process Automation

### Automation Opportunities

#### High Value Automation
- Repetitive tasks
- Error-prone activities
- Time-consuming processes
- Compliance checks
- Metrics collection

#### Automation Approach
1. Identify automation candidates
2. Calculate ROI
3. Pilot automation
4. Measure effectiveness
5. Scale successful automation

### Current Automation Targets
- CI/CD pipeline
- Code quality checks
- Test execution
- Documentation generation
- Metrics collection

## Continuous Improvement

### Improvement Process

#### Kaizen Approach
- Small, incremental improvements
- Everyone contributes ideas
- Quick implementation
- Measure and adjust
- Celebrate successes

#### Innovation Initiatives
- Breakthrough improvements
- New tool adoption
- Process reengineering
- Paradigm shifts

### Improvement Tracking

```markdown
## Process Improvement Log

### Improvement ID: PI-XXX-001
**Date Identified**: YYYY-MM-DD
**Process**: [Process name]
**Current State**: [Problem/opportunity]
**Future State**: [Desired outcome]
**Actions Taken**: [Implementation steps]
**Results**: [Measured outcomes]
**Status**: Planned/Active/Complete
```

## Process Integration

### Integration with 4-Phase Cycle

**Research Phase**:
- Process alignment check
- Identify process gaps
- Plan process adaptations

**Planning Phase**:
- Define phase-specific processes
- Allocate process activities
- Set process metrics

**Implementation Phase**:
- Execute processes
- Monitor compliance
- Collect metrics

**Review Phase**:
- Assess process effectiveness
- Identify improvements
- Update documentation

## Process Culture

### Principles
- **Consistency**: Follow defined processes
- **Flexibility**: Adapt when justified
- **Improvement**: Always seek better ways
- **Measurement**: Data-driven decisions
- **Ownership**: Everyone owns process quality

### Best Practices
1. Document processes clearly
2. Train on process changes
3. Measure process effectiveness
4. Celebrate process improvements
5. Learn from process failures

## Escalation Triggers

### When to Escalate to User
- Major process breakdown
- Systematic non-compliance
- Process conflict unresolved
- Resource constraints
- Stakeholder complaints
- Critical process changes needed