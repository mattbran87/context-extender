# SME Subagent Team Organization Framework

## Overview
This document defines how our 4 SME subagents (Technical Governance, Quality Governance, Risk Governance, Process Governance) integrate with team structures and scale as the project grows from a 2-person team to larger organizations.

## Current Team Structure (2-Person Team)

### Team Composition
- **User**: Product Owner, Stakeholder, Decision Authority
- **Claude**: Developer, Process Owner, SME Orchestrator

### SME Subagent Roles in 2-Person Team

#### Technical Governance SME (Enhanced with Specialized Competencies)
- **Primary Role**: Technical decision authority with specialized advisory competencies
- **Specialized Competencies**: 
  - Go Language Specialist (advanced patterns, performance, concurrency)
  - CLI Development Specialist (UX design, command structure, cross-platform)
  - Claude Code Specialist (extension patterns, SDK integration, hooks/MCP)
- **Consultation Frequency**: As-needed for technical decisions, specialists consulted first
- **Team Integration**: Direct consultation during Implementation phase, specialists advise on domain-specific decisions
- **Decision Flow**: Specialists → Technical SME → Implementation (User for major architecture)
- **Escalation**: To User for architectural decisions, complex Claude Code integrations

#### Quality Governance SME
- **Primary Role**: Quality standard enforcer
- **Consultation Frequency**: Continuous during development
- **Team Integration**: Integrated into CI/CD pipeline
- **Escalation**: To User for quality standard changes

#### Risk Governance SME
- **Primary Role**: Risk identification and mitigation
- **Consultation Frequency**: Daily risk assessment
- **Team Integration**: Risk alerts during all phases
- **Escalation**: To User for high/critical risks

#### Process Governance SME
- **Primary Role**: Process compliance and improvement
- **Consultation Frequency**: Phase transitions and retrospectives
- **Team Integration**: Embedded in workflow validation
- **Escalation**: To User for major process changes

## Team Scaling Strategy

### 3-5 Person Team (Next Growth Phase)

#### New Roles
- **Technical Lead**: Senior developer
- **QA Engineer**: Dedicated quality assurance
- **DevOps Engineer**: Infrastructure and deployment

#### SME Subagent Adaptation

**Technical Governance SME (with Specialized Competencies)**:
- **Expanded Role**: Senior technical advisor to Technical Lead and team
- **New Responsibilities**: Architecture review, technology evaluation, specialist competency coordination
- **Specialized Competencies**: Go Language, CLI Development, Claude Code integration expertise
- **Team Integration**: Technical design reviews, specialist consultation, code review guidance
- **Consultation Protocol**: 
  - Team members → Relevant Specialist → Technical SME → Technical Lead
  - Complex decisions → Technical SME → Technical Lead → User
- **Competency Distribution**: Specialists can advise team directly on domain-specific questions

**Quality Governance SME**:
- **Expanded Role**: Partner with QA Engineer
- **New Responsibilities**: Test strategy validation, quality metrics analysis
- **Team Integration**: QA process design, automated quality checks
- **Consultation Protocol**: QA Engineer implements SME recommendations

**Risk Governance SME**:
- **Expanded Role**: Team-wide risk advisor
- **New Responsibilities**: Cross-team risk coordination
- **Team Integration**: Risk assessment in sprint planning
- **Consultation Protocol**: All team members can request risk assessment

**Process Governance SME**:
- **Expanded Role**: Process coach for entire team
- **New Responsibilities**: Team process training, compliance monitoring
- **Team Integration**: Scrum Master equivalent for process adherence
- **Consultation Protocol**: Process questions go to SME first

### 6-10 Person Team (Medium Team)

#### New Structure
- **Development Teams**: 2 feature teams of 3-4 developers each
- **Specialized Roles**: UX Designer, Product Manager, Security Engineer
- **Support Roles**: Scrum Master, Technical Writer

#### SME Subagent Evolution

**Technical Governance SME (with Specialized Competencies)**:
- **Cross-Team Role**: Architecture consistency and specialized expertise across teams
- **New Focus**: Technology standardization, cross-team integration, competency center
- **Specialized Competencies**: Become centers of excellence for Go, CLI, and Claude Code
- **Consultation Model**: 
  - Weekly technical reviews with leads
  - Specialist competencies available to all teams
  - Cross-team knowledge sharing sessions
- **Competency Distribution**: Each team may have specialist liaisons
- **Escalation Path**: Team → Specialist → Technical SME → Technical Lead → User

**Quality Governance SME**:
- **Quality Standards**: Uniform quality across all teams
- **New Focus**: Quality metrics aggregation, cross-team quality trends
- **Consultation Model**: Quality representatives from each team
- **Integration**: Quality gates in cross-team dependencies

**Risk Governance SME**:
- **Enterprise Risk**: Portfolio-level risk management
- **New Focus**: Inter-team risks, dependency risks
- **Consultation Model**: Risk champions in each team
- **Integration**: Risk dashboard for all stakeholders

**Process Governance SME**:
- **Process Standardization**: Consistent processes across teams
- **New Focus**: Process efficiency optimization, team coordination
- **Consultation Model**: Process owners in each team
- **Integration**: Cross-team retrospectives and improvement

### 10+ Person Team (Large Team/Organization)

#### Organizational Structure
- **Multiple Product Teams**: Independent feature teams
- **Platform Teams**: Infrastructure and shared services
- **Center of Excellence**: Architecture, Security, Quality
- **Management Layer**: Engineering Manager, Product Directors

#### SME Subagent Transformation

**Technical Governance SME (with Specialized Competencies)**:
- **Strategic Role**: Technology strategy and standards with specialized expertise
- **Focus**: Platform decisions, technology roadmap, competency development
- **Specialized Competencies**: Organizational centers of excellence
  - Go Language Center of Excellence
  - CLI Development Center of Excellence  
  - Claude Code Integration Center of Excellence
- **Consultation Model**: Architecture Review Board with specialist panels
- **Integration**: Technology governance committee with competency councils
- **Knowledge Management**: Organizational learning and best practice development

**Quality Governance SME**:
- **Quality Program**: Organization-wide quality initiatives
- **Focus**: Quality engineering practices, tooling strategy
- **Consultation Model**: Quality Council with team representatives
- **Integration**: Quality metrics program, quality engineering team

**Risk Governance SME**:
- **Enterprise Risk Management**: Systematic risk program
- **Focus**: Business continuity, compliance, operational risk
- **Consultation Model**: Risk Management Committee
- **Integration**: Risk register, risk reporting dashboards

**Process Governance SME**:
- **Process Excellence**: Continuous improvement program
- **Focus**: Process maturity, operational excellence
- **Consultation Model**: Process Improvement Office
- **Integration**: Process metrics, improvement initiatives

## SME Collaboration Patterns

### Intra-SME Collaboration

#### Daily Coordination
```
Technical ←→ Quality: Architecture impacts on testability
Quality ←→ Risk: Quality trends as risk indicators  
Risk ←→ Process: Risk mitigation through process
Process ←→ Technical: Process efficiency vs technical debt
```

#### Weekly Integration
- **Technical + Quality**: Code review and testing strategy alignment
- **Risk + Process**: Risk assessment methodology and process gaps
- **All SMEs**: Cross-cutting issue identification

#### Cycle-End Synthesis
- Joint recommendation for cycle improvements
- Integrated risk and quality assessment
- Technical debt vs process debt analysis
- Collective input for next cycle planning

### SME-Team Integration Patterns

#### Embedded Consultation (Current 2-Person)
- SMEs directly accessible to Claude
- Real-time decision support
- Immediate feedback loops

#### Representative Model (3-10 Person)
- Team members designated as SME liaisons
- Regular consultation schedules
- SME guidance implementation by representatives

#### Center of Excellence (10+ Person)
- SMEs become organizational capabilities
- Formal governance structures
- Strategic guidance and standards setting

## Decision-Making Framework

### SME Authority Levels by Team Size

#### 2-Person Team
| SME | Autonomous Decisions | Collaborative Decisions | Escalation Required |
|-----|---------------------|------------------------|---------------------|
| Technical | Tool selection, patterns | Architecture, tech stack | Major platform changes |
| Quality | Test strategies, metrics | Quality standards | Quality policy changes |
| Risk | Risk assessments | Mitigation strategies | High/Critical risks |
| Process | Process adjustments | New processes | Major process changes |

#### 3-10 Person Team
| SME | Advisory Role | Approval Role | Governance Role |
|-----|---------------|---------------|-----------------|
| Technical | Code review guidance | Architecture decisions | Technology standards |
| Quality | Quality coaching | Test approaches | Quality gates |
| Risk | Risk identification | Mitigation plans | Risk thresholds |
| Process | Process support | Process changes | Process compliance |

#### 10+ Person Team
| SME | Strategic Input | Policy Development | Standards Enforcement |
|-----|-----------------|-------------------|---------------------|
| Technical | Technology roadmap | Architecture principles | Development standards |
| Quality | Quality strategy | Quality policies | Quality metrics |
| Risk | Risk strategy | Risk policies | Risk monitoring |
| Process | Process strategy | Process frameworks | Process auditing |

## Communication and Reporting

### SME Reporting Structure

#### Current (2-Person)
```
SME Subagents → Claude → User
```

#### 3-10 Person Team
```
SME Subagents → Team Leads → Claude → User
```

#### 10+ Person Team
```
SME Subagents → Governance Committees → Management → User
```

### Communication Cadence

#### Daily (All Team Sizes)
- Automated risk alerts
- Quality gate notifications
- Process compliance checks
- Technical decision logging

#### Weekly
- SME consultation summaries
- Cross-SME collaboration notes
- Team feedback on SME guidance
- Metrics and trends analysis

#### Cycle-End
- Comprehensive SME assessment
- Recommendations for next cycle
- Process improvement suggestions
- Risk and quality trends

## Success Metrics

### SME Effectiveness by Team Size

#### 2-Person Team Metrics
- **Consultation Response Time**: < 5 minutes
- **Decision Quality**: Measured by outcomes
- **User Satisfaction**: Feedback on SME guidance
- **Process Adherence**: Compliance rates

#### 3-10 Person Team Metrics
- **Team Adoption Rate**: % of SME recommendations implemented
- **Cross-Team Consistency**: Standardization measures
- **Efficiency Gains**: Process improvement metrics
- **Quality Improvements**: Trend analysis

#### 10+ Person Team Metrics
- **Strategic Alignment**: Technology/process roadmap adherence
- **Risk Reduction**: Enterprise risk mitigation
- **Quality Excellence**: Organization-wide quality metrics
- **Process Maturity**: Process optimization indicators

## Implementation Roadmap

### Phase 1: Current State Optimization (Next 2 Cycles)
- Refine SME consultation protocols
- Improve SME collaboration patterns
- Document decision-making effectiveness
- Prepare for first team expansion

### Phase 2: Small Team Scaling (Cycles 3-6)
- Implement representative model
- Train team members on SME consultation
- Establish cross-team coordination
- Develop governance rhythms

### Phase 3: Medium Team Evolution (Future)
- Transform to advisory committees
- Implement governance structures
- Establish Centers of Excellence
- Create strategic guidance framework

### Phase 4: Large Team Transformation (Long-term)
- Full organizational integration
- Strategic governance roles
- Policy and standards development
- Enterprise risk and quality management

## Continuous Evolution

### SME Capability Maturity

**Level 1: Reactive Consultation** (Current)
- On-demand advice
- Tactical guidance
- Individual decision support

**Level 2: Proactive Advisory** (3-10 person target)
- Preventive guidance
- Strategic input
- Team capability building

**Level 3: Systematic Governance** (10+ person target)
- Policy development
- Standards enforcement
- Organizational capability

**Level 4: Predictive Excellence** (Future)
- Anticipatory guidance
- Self-optimizing systems
- Autonomous decision support

## Summary: Enhanced SME Structure for Context-Extender

### What We've Achieved
We've successfully enhanced our SME structure to provide specialized technical expertise while maintaining decision efficiency:

#### Enhanced Technical Governance SME
- **Go Language Specialist**: Advanced Go patterns, performance optimization, concurrency design
- **CLI Development Specialist**: Command-line interface design, user experience, cross-platform compatibility  
- **Claude Code Specialist**: Claude Code extension patterns, SDK integration, hooks and MCP development

#### Key Benefits
1. **Specialized Expertise**: Deep domain knowledge for Go, CLI, and Claude Code development
2. **Maintained Authority**: Single Technical SME maintains decision authority
3. **Efficient Consultation**: Specialists advise, Technical SME decides
4. **Clear Escalation**: Straightforward path to User for strategic decisions
5. **Scalable Framework**: Grows from 2-person team to large organization

#### Project-Specific Value
For context-extender development:
- **Context Manipulation**: Go + Claude Code specialist collaboration
- **CLI Interface**: CLI + Claude Code integration expertise
- **Extension Architecture**: Claude Code specialist + Technical SME oversight
- **Performance Optimization**: Go specialist + Claude Code compatibility

### Implementation Status
✅ **Completed**:
- Enhanced Technical Governance SME with specialized competencies
- Created Claude Code Specialist SME for extension patterns
- Developed comprehensive consultation workflows
- Updated team organization framework for all team sizes
- Integrated with existing Quality, Risk, and Process SMEs
- Created detailed documentation and reference guides

### Next Steps for Context-Extender Development
With the enhanced SME structure in place, we're ready to:
1. **Begin Research Phase** with proper technical guidance structure
2. **Leverage Specialists** for Go, CLI, and Claude Code decisions
3. **Maintain Quality** through coordinated SME consultation
4. **Scale Effectively** as the project and team grow

This framework ensures our SME subagents remain valuable and properly integrated regardless of team size, providing the right level of guidance and governance for each stage of organizational growth.