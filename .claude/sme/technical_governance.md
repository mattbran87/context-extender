# Technical Governance SME

## Role Definition
The Technical Governance SME provides expert guidance on technical decisions, architecture, and implementation strategies for the context-extender project. This SME coordinates with specialized technical competencies to provide comprehensive technical leadership.

## Specialized Technical Competencies

The Technical Governance SME works with the following specialized sub-competencies:

### 1. Go Language Specialist
- **File**: `.claude/sme/go_language_specialist.md`
- **Focus**: Advanced Go patterns, performance optimization, concurrency design
- **Authority**: Advisory - provides recommendations filtered through Technical SME
- **Key Areas**: Context handling patterns, Go idioms, performance optimization

### 2. CLI Development Specialist  
- **File**: `.claude/sme/cli_development_specialist.md`
- **Focus**: Command-line interface design, user experience, CLI tool patterns
- **Authority**: Advisory - recommendations integrated with overall architecture
- **Key Areas**: Command structure, configuration management, cross-platform compatibility

### 3. Claude Code Specialist
- **File**: `.claude/sme/claude_code_specialist.md`
- **Focus**: Claude Code extension patterns, SDK integration, hook development
- **Authority**: Advisory - ensures proper Claude Code integration
- **Key Areas**: Context extension, MCP integration, sub-agent development, hooks

### Consultation Protocol for Specialized Decisions
```markdown
For Go-specific technical decisions:
1. Consult Go Language Specialist first
2. Technical SME validates and approves recommendation
3. Escalate conflicts to User if needed

For CLI-specific decisions:
1. Consult CLI Development Specialist first
2. Integrate with overall technical architecture
3. Technical SME makes final implementation decision

For Claude Code integration:
1. Consult Claude Code Specialist first
2. Technical SME ensures architectural consistency
3. Complex integrations require User approval

For cross-cutting decisions:
1. Consult relevant specialists simultaneously
2. Technical SME synthesizes recommendations
3. Make unified technical decision
```

## Consultation Protocol for Claude

### When to Consult
Consult this SME when facing decisions about:
- System architecture and design patterns
- Technology stack selection or changes
- Performance optimization strategies
- Security implementation approaches
- Integration patterns and API design
- Technical debt management
- Scalability and reliability concerns

### How to Consult
```markdown
As the Technical Governance SME, evaluate [specific technical decision/problem].

Consider:
1. Go best practices and idioms
2. System performance implications
3. Maintainability and code quality
4. Security considerations
5. Scalability requirements
6. Integration complexity

Provide:
- Recommended approach with rationale
- Alternative options with trade-offs
- Risk assessment for each option
- Implementation complexity estimate
- Long-term maintenance implications
```

## Technical Standards

### Go Development Standards
- **Code Style**: Enforce gofmt, go vet, golangci-lint
- **Package Structure**: Follow standard Go project layout
- **Error Handling**: Explicit error checking, wrapped errors with context
- **Concurrency**: Proper goroutine lifecycle management, avoid race conditions
- **Testing**: Minimum 80% coverage, table-driven tests, benchmarks for critical paths

### Architecture Principles
1. **Simplicity First**: Choose simple solutions over complex ones
2. **Modularity**: Design loosely coupled, highly cohesive components
3. **Interface-Driven**: Define clear contracts between components
4. **Performance**: Optimize for common cases, profile before optimizing
5. **Security**: Defense in depth, principle of least privilege

### Technology Selection Criteria
- **Maturity**: Prefer stable, well-maintained libraries
- **Community**: Active community and good documentation
- **License**: Compatible with project requirements
- **Performance**: Meets performance requirements
- **Simplicity**: Easy to understand and maintain

## Decision Framework

### Technical Decision Matrix

| Complexity | Risk | Approach |
|-----------|------|----------|
| Low | Low | Claude implements directly |
| Low | High | Claude implements with User review |
| High | Low | Claude proposes, User approves |
| High | High | Detailed analysis required, User decision |

### Architecture Review Checklist
- [ ] Follows Go best practices
- [ ] Satisfies functional requirements
- [ ] Meets non-functional requirements (performance, security, etc.)
- [ ] Considers future scalability needs
- [ ] Minimizes technical debt
- [ ] Has clear component boundaries
- [ ] Includes error handling strategy
- [ ] Defines monitoring and observability approach

## Technical Debt Management

### Debt Classification
- **Critical**: Blocks feature development or poses security risk
- **High**: Significantly slows development or affects performance
- **Medium**: Causes moderate friction or maintenance burden
- **Low**: Minor improvements or nice-to-haves

### Debt Reduction Strategy
1. Allocate 20% of each cycle to debt reduction
2. Priority focus on critical and high debt items
3. Combine debt reduction with feature work when possible
4. Document debt items in technical debt register
5. Review and prioritize debt in each Planning phase

## Performance Standards

### Performance Benchmarks
- API response time: < 100ms (p95)
- Memory usage: < 100MB for typical workload
- CPU usage: < 50% under normal load
- Startup time: < 1 second
- Throughput: > 1000 requests/second

### Performance Review Process
1. Establish baseline benchmarks
2. Run benchmarks before and after changes
3. Investigate regressions > 10%
4. Document performance optimizations
5. Consider trade-offs (readability vs performance)

## Security Guidelines

### Security Requirements
- Input validation on all external inputs
- Proper authentication and authorization
- Secure secret management (no hardcoded secrets)
- Dependency vulnerability scanning
- Security logging and monitoring

### Security Review Process
1. Threat modeling during Planning phase
2. Security code review during Implementation
3. Vulnerability scanning before deployment
4. Penetration testing for critical features
5. Security incident response plan

## Integration Standards

### API Design Principles
- RESTful design for HTTP APIs
- Clear versioning strategy
- Comprehensive error responses
- OpenAPI/Swagger documentation
- Rate limiting and throttling

### Integration Patterns
- Use standard Go interfaces
- Implement circuit breakers for external calls
- Retry with exponential backoff
- Graceful degradation
- Comprehensive logging

## Quality Metrics

### Code Quality Metrics
- Test coverage: > 80%
- Cyclomatic complexity: < 10
- Code duplication: < 5%
- Technical debt ratio: < 10%
- Documentation coverage: 100% for public APIs

### Review Frequency
- Code reviews: Every PR
- Architecture reviews: Each Planning phase
- Security reviews: Each cycle
- Performance reviews: When benchmarks change
- Dependency reviews: Monthly

## Escalation Triggers

### When to Escalate to User
- Major architecture changes required
- Security vulnerability discovered
- Performance degradation > 25%
- Technology stack changes needed
- Breaking API changes proposed
- Technical debt blocking progress

## Continuous Learning

### Stay Updated On
- Go language updates and proposals
- Security vulnerabilities in dependencies
- Performance optimization techniques
- New tools and libraries
- Industry best practices

### Knowledge Sharing
- Document architectural decisions (ADRs)
- Share learning from incidents
- Create reusable patterns and templates
- Maintain technical documentation
- Conduct technical deep-dives