# Quality Governance SME

## Role Definition
The Quality Governance SME ensures that all deliverables meet quality standards, establishes quality processes, and drives continuous quality improvement throughout the context-extender project.

## Consultation Protocol for Claude

### When to Consult
Consult this SME when dealing with:
- Quality standard definitions
- Testing strategy development
- Code review processes
- Documentation requirements
- Quality gate criteria
- Defect management
- Quality metrics and reporting
- Process improvement initiatives

### How to Consult
```markdown
As the Quality Governance SME, evaluate [specific quality concern/decision].

Consider:
1. Current quality standards and benchmarks
2. Impact on user satisfaction
3. Technical debt implications
4. Testing coverage and effectiveness
5. Process efficiency and overhead
6. Long-term maintainability

Provide:
- Quality assessment of current state
- Recommended quality improvements
- Quality risks and mitigation strategies
- Measurement criteria and metrics
- Implementation priority and timeline
```

## Quality Standards Framework

### Code Quality Standards

#### Mandatory Requirements
- **Test Coverage**: Minimum 80% for all packages
- **Documentation**: 100% GoDoc coverage for all public packages, types, functions, and methods
- **Documentation Detail**: Minimum 3 sentences per documentation block with examples
- **Linting**: Zero golangci-lint errors
- **Security**: Zero high/critical vulnerabilities
- **Performance**: Meet defined benchmarks

#### Quality Gates by Phase

**Planning Phase Gates**:
- [ ] User stories have clear acceptance criteria
- [ ] Test strategy defined for each story
- [ ] Quality metrics identified
- [ ] Definition of Done established

**Implementation Phase Gates**:
- [ ] All code has accompanying tests
- [ ] Code reviews completed
- [ ] Documentation updated
- [ ] No regression in quality metrics

**Review Phase Gates**:
- [ ] All acceptance criteria met
- [ ] Quality metrics within targets
- [ ] No critical defects open
- [ ] Stakeholder approval obtained

### Testing Strategy

#### Test Pyramid
```
         /\
        /  \    End-to-End Tests (10%)
       /____\   - User journey tests
      /      \  - System integration tests
     /________\ Integration Tests (20%)
    /          \ - API tests
   /____________\ - Component integration
  /              \ Unit Tests (70%)
 /________________\ - Function-level tests
                    - Edge cases
                    - Error conditions
```

#### Test Types and Coverage

| Test Type | Coverage Target | Responsibility | When |
|-----------|----------------|----------------|------|
| Unit Tests | 80% minimum | Claude | During development |
| Integration Tests | Critical paths | Claude | Before PR |
| E2E Tests | User journeys | Claude | Before deployment |
| Performance Tests | Key operations | Claude | Each cycle |
| Security Tests | All inputs | Claude | Each cycle |

### Code Review Standards

#### Review Checklist
- [ ] **Functionality**: Code does what it's supposed to do
- [ ] **Design**: Follows architectural patterns
- [ ] **Complexity**: Simple and readable
- [ ] **Tests**: Adequate test coverage
- [ ] **Naming**: Clear and consistent naming
- [ ] **Comments**: Necessary comments present
- [ ] **GoDoc Documentation**: Complete documentation for all public elements
- [ ] **Documentation Quality**: Detailed purpose, functionality, parameters, returns, errors, examples
- [ ] **Documentation Examples**: Working code examples for all public functions
- [ ] **Security**: No security vulnerabilities
- [ ] **Performance**: No obvious performance issues

#### Review Response Times
- Critical fixes: Within 2 hours
- Normal changes: Within 1 day
- Large changes: Within 2 days

## Defect Management

### Defect Classification

| Severity | Description | Response Time | Resolution Time |
|----------|------------|---------------|-----------------|
| **Critical** | System unusable, data loss | Immediate | 4 hours |
| **High** | Major feature broken | 2 hours | 1 day |
| **Medium** | Feature partially working | 1 day | 3 days |
| **Low** | Minor issue, cosmetic | Next cycle | As scheduled |

### Defect Lifecycle
```
New → Confirmed → In Progress → Fixed → Verified → Closed
                       ↓
                   Cannot Reproduce → Closed
```

### Root Cause Analysis
For all High and Critical defects:
1. Identify root cause
2. Document in defect report
3. Identify prevention measures
4. Update processes/tests to prevent recurrence
5. Share learnings in retrospective

## Quality Metrics

### Product Quality Metrics

#### Code Metrics
- **Coverage**: Test coverage percentage
- **Complexity**: Cyclomatic complexity
- **Duplication**: Code duplication percentage
- **Debt**: Technical debt ratio
- **Vulnerabilities**: Security vulnerability count

#### Defect Metrics
- **Defect Density**: Defects per story point
- **Escape Rate**: Defects found after deployment
- **Resolution Time**: Average time to fix
- **Regression Rate**: Reintroduced defects
- **Root Cause Distribution**: By category

### Process Quality Metrics

#### Efficiency Metrics
- **First Time Pass Rate**: Stories passing review first time
- **Rework Percentage**: Time spent on rework
- **Review Turnaround**: Code review completion time
- **Test Execution Time**: Total test suite runtime
- **Build Success Rate**: CI/CD pipeline success

#### Effectiveness Metrics
- **Customer Satisfaction**: Stakeholder feedback scores
- **Cycle Time**: Idea to production time
- **Predictability**: Actual vs estimated effort
- **Quality Trends**: Quality metrics over time

## Documentation Standards

### Documentation Requirements

| Document Type | Required | Format | Owner | Review |
|--------------|----------|--------|-------|--------|
| API Documentation | Yes | OpenAPI/Swagger | Claude | Each change |
| Code Comments | Yes | GoDoc format | Claude | Code review |
| Architecture Docs | Yes | Markdown + diagrams | Claude | Planning phase |
| User Guides | Future | Markdown | Claude | Each release |
| Process Docs | Yes | Markdown | Claude | Each cycle |

### Documentation Quality Criteria
- **Accuracy**: Matches implementation
- **Completeness**: Covers all features
- **Clarity**: Easy to understand
- **Currency**: Up to date
- **Accessibility**: Easy to find and navigate

## Continuous Improvement

### Quality Improvement Process

1. **Measure**: Collect quality metrics
2. **Analyze**: Identify trends and patterns
3. **Improve**: Implement improvements
4. **Control**: Monitor effectiveness
5. **Standardize**: Update processes

### Improvement Opportunities

#### Short-term (Each Cycle)
- Review and update quality gates
- Optimize test execution time
- Improve code review turnaround
- Reduce defect escape rate

#### Medium-term (Quarterly)
- Enhance automation coverage
- Implement new quality tools
- Refine quality metrics
- Update quality standards

#### Long-term (Annually)
- Achieve quality certifications
- Implement advanced testing techniques
- Establish quality culture
- Benchmark against industry standards

## Quality Risk Management

### Quality Risk Categories
- **Technical Debt**: Accumulation affecting quality
- **Skills Gap**: Lack of quality expertise
- **Process Gaps**: Missing quality processes
- **Tool Limitations**: Inadequate quality tools
- **Time Pressure**: Quality compromised for speed

### Risk Mitigation Strategies
- Regular quality audits
- Continuous training and learning
- Process automation
- Tool evaluation and upgrades
- Quality-first mindset

## Escalation Triggers

### When to Escalate to User
- Quality gates consistently failing
- Critical defect trends increasing
- Customer satisfaction declining
- Quality debt blocking progress
- Process compliance issues
- Resource constraints affecting quality

## Quality Tools and Automation

### Recommended Tools
- **Testing**: Go test, testify, gomock
- **Coverage**: go test -cover, coveralls
- **Linting**: golangci-lint
- **Security**: gosec, snyk
- **Performance**: pprof, go test -bench
- **Documentation**: godoc, swagger

### Automation Strategy
1. Automate repetitive quality checks
2. Integrate quality gates in CI/CD
3. Automated test execution on commit
4. Automated quality reporting
5. Automated dependency updates

## Quality Culture

### Principles
- **Quality is everyone's responsibility**
- **Prevention over detection**
- **Continuous improvement mindset**
- **Data-driven quality decisions**
- **Customer-focused quality**

### Practices
- Regular quality reviews
- Quality metrics visibility
- Celebrate quality achievements
- Learn from quality failures
- Share quality best practices