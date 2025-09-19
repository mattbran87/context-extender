# SME Subagent System Documentation

## Overview
The context-extender project uses an enhanced SME (Subject Matter Expert) subagent system that combines 4 primary SMEs with specialized technical competencies to provide comprehensive guidance throughout the development lifecycle.

## SME Structure

### Primary SMEs (Decision Authority)
1. **Technical Governance SME** - Technical decision authority with specialized competencies
2. **Quality Governance SME** - Quality standards and processes
3. **Risk Governance SME** - Risk identification and mitigation  
4. **Process Governance SME** - Process compliance and improvement

### Technical Specialized Competencies (Advisory)
Embedded within the Technical Governance SME:
1. **Go Language Specialist** - Advanced Go patterns, performance, concurrency
2. **CLI Development Specialist** - Command-line interface design and UX
3. **Claude Code Specialist** - Claude Code extension and integration patterns

## Documentation Structure

```
.claude/sme/
├── README.md                    # This overview document
├── technical_governance.md      # Primary technical SME with competencies
├── go_language_specialist.md    # Go patterns and optimization expertise  
├── cli_development_specialist.md # CLI design and UX expertise
├── claude_code_specialist.md    # Claude Code extension expertise
├── quality_governance.md        # Quality standards and processes
├── risk_governance.md          # Risk management and assessment
├── process_governance.md       # Process compliance and improvement
├── consultation_workflows.md   # Detailed consultation workflows
└── team_organization.md        # Team scaling and organization
```

## Quick Reference: When to Consult Which SME

### Technical Decisions
- **Go Language Questions** → Go Language Specialist → Technical SME
- **CLI Design Questions** → CLI Development Specialist → Technical SME  
- **Claude Code Integration** → Claude Code Specialist → Technical SME
- **Architecture Decisions** → Multiple Specialists → Technical SME → User (if major)
- **Cross-cutting Technical** → Relevant Specialists → Technical SME

### Quality & Process
- **Quality Standards** → Quality Governance SME
- **Testing Strategy** → Quality SME + relevant Technical Specialist
- **Risk Assessment** → Risk Governance SME
- **Process Questions** → Process Governance SME

## Consultation Workflows

### Simple Decisions (Claude Autonomous)
Standard patterns, basic implementations, common practices

### Specialized Decisions (Specialist → Technical SME)
Domain-specific expertise needed, optimization decisions, design patterns

### Complex Decisions (Multiple Specialists → Technical SME)  
Cross-cutting concerns, architectural implications, performance vs maintainability

### Strategic Decisions (All SMEs → User)
Major architecture changes, technology stack decisions, significant risk

## Key Benefits of Enhanced Structure

### Specialized Expertise
- **Deep Knowledge**: Go, CLI, and Claude Code specific expertise
- **Pattern Recognition**: Domain-specific best practices
- **Optimization**: Performance and design optimization in specialized areas

### Maintained Simplicity  
- **Single Technical Authority**: Technical SME maintains decision authority
- **Clear Escalation**: Straightforward escalation to User
- **Efficient Process**: Specialists advise, SME decides

### Scalability
- **Team Growth**: Scales from 2-person to large teams
- **Knowledge Transfer**: Specialists become centers of excellence
- **Consistent Standards**: Uniform practices across teams

## Context-Extender Specific Applications

### Core Project Focus
Since context-extender enhances Claude Code functionality:

#### Primary Consultation Areas
1. **Context Manipulation Patterns** (Go + Claude Code Specialists)
2. **CLI Interface Design** (CLI + Claude Code Specialists)  
3. **Extension Architecture** (Claude Code + Technical SME)
4. **Performance Optimization** (Go + Claude Code Specialists)

#### Common Decision Scenarios
- **Context Processing Implementation**: Go patterns for efficient context handling
- **Claude Code Hook Development**: Extension points and integration patterns
- **CLI Command Structure**: User experience and Claude Code workflow integration
- **MCP Server Development**: Model Context Protocol integration patterns

## Getting Started

### For New Claude Instances
1. **Read this overview** to understand the SME structure
2. **Review consultation_workflows.md** for detailed decision trees
3. **Consult specialists first** for domain-specific decisions
4. **Escalate through Technical SME** for final technical decisions
5. **Involve User** for strategic and architectural decisions

### For Specific Technical Areas

#### Go Development
- Start with Go Language Specialist for patterns and optimization  
- Escalate complex decisions to Technical SME
- Consider Claude Code integration implications

#### CLI Development
- Consult CLI Development Specialist for UX and design
- Integrate with Claude Code Specialist for workflow compatibility
- Technical SME for final architecture decisions

#### Claude Code Integration  
- Always consult Claude Code Specialist for extension decisions
- Consider security and performance implications
- Technical SME for architectural integration
- User approval for major integration changes

## Success Metrics

The SME system is effective when:
- **Decisions are well-informed** with appropriate specialist input
- **Implementation is consistent** with established patterns
- **Quality is maintained** across all technical areas
- **User involvement** is limited to strategic decisions
- **Development velocity** is maintained or improved

## Continuous Evolution

This SME structure is designed to:
- **Adapt to team growth** from 2-person to large teams
- **Incorporate new competencies** as project needs expand
- **Maintain decision efficiency** while providing specialized expertise
- **Support knowledge transfer** and organizational learning

The enhanced SME system ensures context-extender development leverages both broad governance expertise and deep technical specialization while maintaining clear decision authority and efficient workflows.