# Context-Extender Complete Project Guide

---

## 📋 Section 1: Quick Start & Current Status

### **CRITICAL: Read This First**

**Current Project Status**: PLANNING PHASE - NO CODE CREATION
**Project Phase**: Structure planning and process definition
**Code Creation**: ❌ RESTRICTED until explicitly authorized

### **Immediate Actions for New Claude Instances:**
1. ✅ Read this complete guide
2. ✅ Check `.claude/current_status` file for phase information
3. ✅ Understand the 4-phase cyclical approach
4. ✅ Recognize user approval is required for ALL phase transitions
5. ✅ **Consult SME subagents** for specialized decisions (see Section 7)
6. ✅ **Follow established processes** for all activities (see Section 8)

### **Key Rules & Restrictions:**
- **NO CODE CREATION** until planning phase complete
- **USER APPROVAL REQUIRED** for every phase transition
- **16 USER INTERACTION POINTS** per development cycle
- **DOCUMENT EVERYTHING** according to established templates
- **CONSULT SME SUBAGENTS** for technical, quality, risk, and process decisions
- **FOLLOW ESTABLISHED PROCESSES** for risk, communication, governance, CI/CD, and feedback
- **MANDATORY GODOC COMMENTING** for all code - detailed documentation required for packages, types, functions, and methods

---

## 🎯 Section 2: Project Overview

### **Project Description**
Context-Extender is a Go module designed to enhance context handling capabilities. The project is currently in initial planning phase with minimal Go module configuration.

### **Core Philosophy**
The project uses a 4-phase cyclical development approach focusing on small incremental changes to deliver value continuously while maintaining quality and stakeholder engagement.

### **Development Approach: Cyclical Framework**
- **17-day cycles** (approximately 3 weeks)
- **4 phases per cycle**: Research → Planning → Implementation → Review
- **Scrumban methodology** integration
- **Go-specific optimizations** throughout

### **Team Structure**
- **Technical Lead**: Architecture decisions, technical design oversight
- **Senior Developers**: Feature implementation, code review, mentoring
- **DevOps Engineer**: CI/CD pipeline, infrastructure, deployment automation
- **Project Manager**: Cycle coordination, stakeholder communication, resource planning
- **Product Owner**: Customer advocacy, requirement definition, acceptance criteria validation

### **Success Metrics**
- **Cycle Completion Rate**: Percentage of cycles delivering planned value
- **Customer Satisfaction**: Post-cycle stakeholder feedback scores
- **Technical Quality**: Code coverage, performance benchmarks, security scores
- **Time to Market**: Cycle idea to production deployment time
- **Team Velocity**: Story points or features delivered per cycle

---

## 🔄 Section 3: Cyclical Framework Overview

### **Cycle Structure (17 days total)**

#### **Phase 1: Research (3 days / 18%)**
**Objective**: Investigate and define the increment scope

**Key Activities**:
- Problem definition and scope validation
- Technical feasibility assessment
- Customer feedback analysis and user research
- Risk analysis and mitigation strategies
- Competitive analysis and market research
- **User story creation**: Draft initial epics and user stories based on customer needs

**Key Deliverables**:
- Problem statement document
- Technical feasibility report
- Customer insights summary
- Risk assessment matrix
- Resource requirements estimate
- **Initial user story backlog**: High-level epics and stories with user value statements

#### **Phase 2: Planning (4 days / 24%)**
**Objective**: Design and prepare for implementation

**Key Activities**:
- Detailed technical design and architecture
- Customer validation of planned features
- **User story refinement**: Break down epics, add technical details, estimate complexity
- Acceptance criteria definition with stakeholders
- Resource allocation and timeline planning
- Test case planning and quality gate definition
- **Story point estimation**: Size stories for sprint planning

**Key Deliverables**:
- Technical design document
- **Sprint-ready user stories**: Refined stories with INVEST criteria met
- **Acceptance criteria**: Testable conditions for each story
- Implementation roadmap and timeline
- Test plan and quality checklist
- Resource allocation plan

#### **Phase 3: Implementation (7 days / 41%)**
**Objective**: Build and test the working increment

**Key Activities**:
- Test-driven development following Go best practices
- Continuous integration with automated testing
- Code review process with quality gates
- Performance monitoring and optimization
- Progressive feature rollout using feature flags

**Key Deliverables**:
- Working software increment (deployable)
- Comprehensive test suite (unit, integration, e2e)
- Code review reports and quality metrics
- Performance benchmarks and optimization reports
- Updated documentation

#### **Phase 4: Review (3 days / 17%)**
**Objective**: Evaluate, learn, and prepare for next cycle

**Key Activities**:
- Stakeholder demonstrations and feedback collection
- Metrics analysis and performance evaluation
- Retrospectives and lessons learned documentation
- Customer satisfaction measurement
- Next cycle preparation and backlog refinement

**Key Deliverables**:
- Stakeholder feedback summary
- Metrics and performance analysis report
- Retrospective action items
- Customer satisfaction scores
- Next cycle research preparation

### **Directory Structure**
```
context-extender/
├── cmd/                    # Application entry points
│   └── context-extender/   # Main application
├── internal/               # Private application code
│   ├── config/            # Configuration management
│   ├── handlers/          # HTTP/gRPC handlers
│   ├── models/            # Data models
│   ├── repository/        # Data access layer
│   └── services/          # Business logic
├── pkg/                   # Public library code (if needed)
├── api/                   # API definitions (OpenAPI/Swagger)
├── docs/                  # Project documentation
│   ├── adr/               # Architecture Decision Records
│   └── cycles/            # Cycle-specific documentation
├── scripts/               # Build and deployment scripts
├── deployments/           # Docker, k8s configurations
├── configs/               # Environment-specific configurations
├── test/                  # Test fixtures and data
├── cycles/                # Cyclical development artifacts
│   ├── templates/         # Phase templates and standards
│   ├── cycle-001/         # Individual cycle documentation
│   │   ├── research/
│   │   ├── planning/
│   │   ├── implementation/
│   │   └── review/
│   └── metrics/           # Cross-cycle analytics and reports
├── .claude/               # Claude configuration and guides
│   ├── project_guide.md   # This file
│   └── current_status     # Current phase tracking
├── .github/               # GitHub workflows and templates
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

### **User Story Management**

#### **Research Phase - Story Creation**
**When**: Days 1-3 of Research phase
**Purpose**: Capture user needs and create initial story backlog

**Activities**:
- Conduct customer interviews to identify pain points
- Create user personas and journey maps
- Draft initial epics and high-level user stories
- Focus on user value and business outcomes

**Output**: Initial story backlog with epics and broad user stories

#### **Planning Phase - Story Refinement**
**When**: Days 1-2 of Planning phase
**Purpose**: Refine stories for implementation readiness

**Activities**:
- Break down epics into implementable user stories
- Add technical constraints and Go-specific requirements
- Define testable acceptance criteria
- Estimate story points using team consensus
- Apply INVEST criteria validation

**Output**: Sprint-ready stories with acceptance criteria and estimates

#### **User Story Template**
```markdown
## User Story: [ID] [Title]

**As a** [user type]
**I want** [functionality]
**So that** [business value]

### Acceptance Criteria
- [ ] Given [context]
- [ ] When [action]
- [ ] Then [expected outcome]

### Technical Notes (Go-specific)
- Package/module boundaries
- Interface contracts required
- Error handling patterns
- Concurrency requirements
- Performance benchmarks

### Definition of Ready Checklist
- [ ] User value clearly defined
- [ ] Acceptance criteria testable
- [ ] Technical approach agreed
- [ ] Dependencies identified
- [ ] Story estimated (points)
- [ ] Fits within sprint capacity

### Definition of Done
- [ ] Code complete and follows Go conventions
- [ ] Unit tests written (>80% coverage)
- [ ] Integration tests pass
- [ ] Code review approved
- [ ] Documentation updated
- [ ] Performance benchmarks met
```

### **Quality Standards and Gates**

#### **Code Quality Requirements**
- Test coverage: Minimum 80%
- All code must pass `gofmt`, `go vet`, and `golangci-lint`
- Security scanning with `gosec` for vulnerability detection
- Performance benchmarks for critical code paths
- All public APIs must be documented

#### **Phase Completion Criteria**
- **Research → Planning**: Research findings documented, scope confirmed, stakeholder approval
- **Planning → Implementation**: Technical design approved, acceptance criteria defined, team capacity confirmed
- **Implementation → Review**: Code complete, tests passing, documentation updated, performance benchmarks met
- **Review → Next Research**: Retrospective completed, stakeholder feedback collected, lessons learned documented

### **Go-Specific Optimizations**
- Leverage Go's fast compilation for continuous testing during Implementation
- Use `go test -race` for concurrent safety validation
- Implement `golangci-lint` quality gates in CI/CD pipeline
- Include performance benchmarks using Go's built-in testing tools
- Maintain minimal external dependencies with regular `go mod tidy` reviews

### **Documentation Standards**

#### **Mandatory GoDoc Commenting**
All code must include comprehensive GoDoc documentation following Go conventions:

**Package Documentation**:
```go
// Package contextextender provides enhanced context handling capabilities
// for Go applications, enabling advanced context manipulation and extension
// patterns for complex request processing workflows.
//
// The package supports context enrichment, conditional context forwarding,
// and context-aware request routing with built-in performance monitoring
// and error handling.
//
// Example usage:
//
//	extender := contextextender.New()
//	ctx := extender.WithMetadata(context.Background(), "key", "value")
//	result, err := extender.Process(ctx, request)
package contextextender
```

**Type Documentation**:
```go
// ContextExtender manages context enhancement and processing operations.
// It provides thread-safe methods for context manipulation, metadata
// injection, and request routing with configurable middleware support.
//
// The extender maintains internal state for performance metrics and
// supports graceful shutdown with proper cleanup of resources.
type ContextExtender struct {
    // config holds the configuration settings for the extender
    config *Config
    // metrics tracks performance and usage statistics
    metrics *Metrics
    // middleware contains the processing pipeline components
    middleware []Middleware
}
```

**Function Documentation**:
```go
// WithMetadata creates a new context with the specified metadata key-value pair.
// The metadata is stored in the context and can be retrieved downstream
// using the GetMetadata function.
//
// Parameters:
//   - ctx: the parent context to extend
//   - key: the metadata key (must be non-empty string)
//   - value: the metadata value (any type supported)
//
// Returns:
//   - context.Context: new context with metadata attached
//
// Example:
//   ctx := extender.WithMetadata(context.Background(), "userID", "12345")
//   userID, exists := extender.GetMetadata(ctx, "userID")
//
// The function is thread-safe and does not modify the parent context.
// If the key is empty, the function panics with ErrInvalidKey.
func (e *ContextExtender) WithMetadata(ctx context.Context, key string, value interface{}) context.Context {
    // implementation...
}
```

**Method Documentation**:
```go
// Process executes the configured middleware pipeline on the given request
// within the provided context. Each middleware component is executed in
// sequence, with the ability to modify the context, request, or response.
//
// The processing includes automatic error handling, performance tracking,
// and graceful degradation in case of middleware failures. If any
// middleware returns an error, processing stops and the error is returned
// with appropriate context information.
//
// Parameters:
//   - ctx: request context with timeout and cancellation support
//   - req: the request object to be processed
//
// Returns:
//   - *Response: processed response object (nil if error occurred)
//   - error: processing error with detailed error context
//
// Errors:
//   - ErrInvalidRequest: when request is nil or invalid
//   - ErrContextCanceled: when context is canceled during processing
//   - ErrMiddlewareFailure: when middleware component fails
//
// The method is thread-safe and supports concurrent request processing.
func (e *ContextExtender) Process(ctx context.Context, req *Request) (*Response, error) {
    // implementation...
}
```

#### **Documentation Requirements**

**All documentation must include**:
1. **Purpose**: Clear description of what the code does
2. **Functionality**: Detailed explanation of how it works
3. **Parameters**: Complete parameter documentation with types and constraints
4. **Return Values**: All return values with types and possible states
5. **Errors**: All possible error conditions and their meanings
6. **Examples**: Working code examples showing typical usage
7. **Thread Safety**: Concurrent access behavior and safety guarantees
8. **Side Effects**: Any state changes or external interactions

**Documentation Quality Gates**:
- **Coverage**: 100% of public packages, types, functions, and methods
- **Detail Level**: Minimum 3 sentences per documentation block
- **Examples**: At least one example per public function
- **Error Cases**: Document all error conditions and handling
- **Performance Notes**: Include performance characteristics for critical paths

#### **Documentation Validation**
- Use `go doc` to verify documentation renders correctly
- Run `godoc -http=:6060` to preview documentation locally
- Include documentation review in code review process
- Automated checks in CI/CD pipeline for documentation completeness

---

## 📝 Section 4: Granular Step-by-Step Workflow

### **PHASE 1: RESEARCH (3 DAYS)**

#### **Day 1: Problem Discovery**

##### **Step 1.1: Phase Initialization** 🚀 START HERE
- Update `.claude/current_status` to Research phase
- Create `cycles/cycle-XXX/research/` directory structure
- Initialize research documentation templates

##### **Step 1.2: Problem Definition Session** 🔴 USER INTERACTION REQUIRED
**Claude Action:**
```
"Starting Research Phase for Cycle-XXX. 
Let's define the problem statement and scope for this cycle.
What specific problem or feature should we focus on?"
```
**Wait for:** User to discuss and define the problem
**Continue when:** User says "Continue the phase" or similar

##### **Step 1.3: Stakeholder Identification** 🔴 USER INTERACTION REQUIRED
**Claude Action:**
```
"Now let's identify the stakeholders and user groups.
Who will be affected by this change? Who are the primary users?"
```
**Wait for:** User input on stakeholders
**Continue when:** User approval received

##### **Step 1.4: Pain Points Documentation** 🔴 USER INTERACTION REQUIRED
**Claude Action:**
```
"What are the current pain points related to this problem?
What specific issues are users experiencing?"
```
**Document:** Create `research/pain-points.md`
**Continue when:** User confirms pain points are captured

##### **Step 1.5: Competitive Analysis**
- Research similar solutions in the ecosystem
- Document findings in `research/competitive-analysis.md`
- Identify unique value propositions

##### **Step 1.6: Risk Assessment**
- Create initial risk matrix
- Document in `research/risk-assessment.md`
- Identify mitigation strategies

#### **Day 2: User Analysis**

##### **Step 2.1: User Persona Creation** 🔴 USER INTERACTION REQUIRED
**Claude Action:**
```
"Let's create user personas for this feature.
Can you describe the typical users and their characteristics?"
```
**Document:** Create `research/user-personas.md`
**Continue when:** User approves personas

##### **Step 2.2: User Journey Mapping**
- Map out user workflows
- Identify touchpoints and interactions
- Document in `research/user-journeys.md`

##### **Step 2.3: Requirements Gathering**
- Compile functional requirements
- Identify non-functional requirements
- Create `research/requirements.md`

##### **Step 2.4: Epic Creation**
- Draft initial epics from user needs
- Structure in `research/epics.md`

##### **Step 2.5: High-Level User Story Creation** 🔴 USER INTERACTION REQUIRED
**Claude Action:**
```
"Based on our analysis, here are the initial user stories:
[List stories]
Do these capture the user needs correctly?"
```
**Document:** Create `research/initial-stories.md`
**Continue when:** User approves initial stories

##### **Step 2.6: Technical Constraints**
- Identify technical limitations
- Document Go-specific considerations
- Create `research/technical-constraints.md`

#### **Day 3: Feasibility & Documentation**

##### **Step 3.1: Technical Feasibility Assessment**
- Evaluate technical approach options
- Assess implementation complexity
- Document in `research/feasibility-report.md`

##### **Step 3.2: Problem Statement Finalization**
- Consolidate all research findings
- Create final `research/problem-statement.md`

##### **Step 3.3: Customer Insights Summary**
- Compile all user feedback
- Create `research/customer-insights.md`

##### **Step 3.4: Risk Mitigation Documentation**
- Finalize mitigation strategies
- Update `research/risk-mitigation.md`

##### **Step 3.5: User Story Backlog Preparation**
- Organize all stories by priority
- Create `research/story-backlog.md`

##### **Step 3.6: Research Summary Report**
- Compile comprehensive phase summary
- Create `research/phase-summary.md`

#### **Phase 1 Completion Gate** 🔴 USER APPROVAL REQUIRED
**Claude Action:**
```
"Research Phase Complete! 
Deliverables ready:
✓ Problem statement
✓ User personas and journeys
✓ Initial story backlog
✓ Feasibility report
✓ Risk assessment

Please review the research outputs.
Approve transition to Planning phase?"
```
**Wait for:** "Approved" or phase transition confirmation

---

### **PHASE 2: PLANNING (4 DAYS)**

#### **Day 1: Story Refinement**

##### **Step 1.1: Phase Initialization** 🚀 START HERE
- Review all research phase outputs
- Create `cycles/cycle-XXX/planning/` directory
- Update `.claude/current_status` to Planning phase

##### **Step 1.2: Epic Breakdown** 🔴 USER INTERACTION REQUIRED
**Claude Action:**
```
"Breaking down epics into implementable user stories:
[List detailed stories]
Please review and provide feedback."
```
**Document:** Create `planning/user-stories.md`
**Wait for:** User approval to proceed

##### **Step 1.3: Technical Details Addition** 🔴 USER INTERACTION REQUIRED
**Claude Action:**
```
"Adding technical details to stories:
[Show stories with technical notes]
Do these technical approaches align with your expectations?"
```
**Update:** `planning/user-stories.md` with technical details
**Wait for:** User approval

##### **Step 1.4: Acceptance Criteria Definition**
- Define testable criteria for each story
- Apply Given/When/Then format
- Update stories with criteria

##### **Step 1.5: INVEST Validation**
- Validate each story against INVEST criteria
- Document validation in `planning/invest-validation.md`

##### **Step 1.6: Dependency Mapping**
- Create story dependency graph
- Document in `planning/dependencies.md`

#### **Day 2: Technical Design**

##### **Step 2.1: Architecture Design** 🔴 USER INTERACTION REQUIRED
**Claude Action:**
```
"Here's the proposed high-level architecture:
[Present architecture diagram/description]
Please review and provide feedback."
```
**Document:** Create `planning/architecture.md`
**Wait for:** User approval

##### **Step 2.2: API Contract Definition**
- Define all API endpoints
- Specify request/response formats
- Document in `planning/api-contracts.md`

##### **Step 2.3: Go Package Structure**
- Design module organization
- Define package boundaries
- Create `planning/package-structure.md`

##### **Step 2.4: Data Model Design**
- Define data structures
- Create schema definitions
- Document in `planning/data-models.md`

##### **Step 2.5: Database Planning** 🔴 USER INTERACTION REQUIRED
**Claude Action:**
```
"Database structure plan:
[Present database design]
Is this database design appropriate?"
```
**Document:** Create `planning/database-design.md`
**Wait for:** User approval

##### **Step 2.6: Security Requirements** 🔴 USER INTERACTION REQUIRED
**Claude Action:**
```
"Let's define security requirements:
What security measures are critical for this feature?"
```
**Document:** Create `planning/security-requirements.md`
**Continue when:** User confirms requirements

#### **Day 3: Estimation & Resource Planning**

##### **Step 3.1: Story Point Estimation**
- Estimate each story's complexity
- Use Fibonacci sequence (1,2,3,5,8,13)
- Document in `planning/estimations.md`

##### **Step 3.2: Story Prioritization**
- Rank stories by value and dependencies
- Create implementation order
- Update `planning/sprint-backlog.md`

##### **Step 3.3: Resource Allocation**
- Assign stories to team members (if applicable)
- Balance workload across implementation days
- Document in `planning/resource-plan.md`

##### **Step 3.4: Timeline Creation**
- Map stories to implementation days
- Account for dependencies
- Create `planning/timeline.md`

##### **Step 3.5: Tool & Dependency Identification**
- List required libraries and tools
- Document Go module dependencies
- Create `planning/tool-dependencies.md`

##### **Step 3.6: Test Strategy Planning**
- Define testing approach
- Plan test coverage targets
- Document in `planning/test-strategy.md`

#### **Day 4: Finalization & Validation**

##### **Step 4.1: Technical Design Compilation**
- Consolidate all technical documents
- Create comprehensive `planning/technical-design.md`

##### **Step 4.2: Stakeholder Validation**
- Review plans with stakeholders
- Incorporate feedback
- Update relevant documents

##### **Step 4.3: Sprint Backlog Finalization**
- Lock in story selection for implementation
- Finalize `planning/sprint-backlog.md`

##### **Step 4.4: Test Plan Creation**
- Detail test cases for each story
- Create `planning/test-plan.md`

##### **Step 4.5: Environment Setup**
- Prepare development environment
- Configure tools and dependencies
- Document in `planning/dev-setup.md`

##### **Step 4.6: Implementation Kickoff Preparation**
- Create implementation checklist
- Prepare kickoff materials

#### **Phase 2 Completion Gate** 🔴 USER APPROVAL REQUIRED
**Claude Action:**
```
"Planning Phase Complete!
Ready for Implementation:
✓ User stories refined and estimated
✓ Technical design documented
✓ Sprint backlog prioritized
✓ Test strategy defined
✓ Development environment ready

Approve transition to Implementation phase?"
```
**Wait for:** User approval

---

### **PHASE 3: IMPLEMENTATION (7 DAYS)**

#### **Day 1-2: Setup & Core Development**

##### **Step 1.1: Phase Initialization** 🚀 START HERE
- Create `cycles/cycle-XXX/implementation/` directory
- Update `.claude/current_status` to Implementation
- Set up version control branches

##### **Step 1.2: Project Structure Setup**
- Initialize Go module structure
- Create package directories
- Set up initial files

##### **Step 1.3: CI/CD Configuration**
- Set up GitHub Actions workflows
- Configure testing pipeline
- Add linting and security scanning

##### **Step 1.4: Core Functionality Implementation**
- Implement first high-priority stories
- Follow TDD approach
- Commit changes regularly

##### **Step 1.5: Initial Unit Tests**
- Write tests for core functions
- Achieve >80% coverage for core
- Run tests continuously

##### **Step 1.6: Logging & Monitoring Setup**
- Implement structured logging
- Add metrics collection
- Configure monitoring hooks

##### **Step 1.7: Initial API Endpoints**
- Implement basic API structure
- Add request/response handling
- Include error handling

#### **Day 3-4: Feature Development**

##### **Step 2.1: Story Implementation Continuation**
- Continue implementing user stories
- Follow priority order
- Update story status regularly

##### **Step 2.2: Comprehensive Testing**
- Write unit tests for all new code
- Maintain >80% coverage
- Add edge case tests

##### **Step 2.3: Code Review Process**
- Submit PRs for review
- Address review feedback
- Ensure code quality standards

##### **Step 2.4: Error Handling Enhancement**
- Implement comprehensive error handling
- Add error recovery mechanisms
- Include proper error logging

##### **Step 2.5: Integration Test Development**
- Create integration test suite
- Test component interactions
- Validate API contracts

##### **Step 2.6: Documentation Updates**
- Update code documentation
- Add inline comments where needed
- Update API documentation

#### **Day 5-6: Integration & Testing**

##### **Step 3.1: Story Completion Push**
- Complete remaining stories
- Focus on MVP requirements
- Defer nice-to-haves if needed

##### **Step 3.2: Full Test Suite Execution**
- Run all unit tests
- Execute integration tests
- Perform end-to-end testing

##### **Step 3.3: Bug Fixing**
- Address failing tests
- Fix identified issues
- Retest after fixes

##### **Step 3.4: Security Scanning**
- Run gosec security scanner
- Address vulnerabilities
- Document security measures

##### **Step 3.5: Performance Benchmarking**
- Run Go benchmarks
- Profile code performance
- Optimize bottlenecks

##### **Step 3.6: Feature Flag Implementation**
- Add feature flags if needed
- Configure toggle mechanisms
- Test flag behaviors

#### **Day 7: Polish & Preparation**

##### **Step 4.1: Final Code Reviews**
- Complete all pending reviews
- Ensure approval on all PRs
- Merge approved changes

##### **Step 4.2: Documentation Finalization**
- Complete all documentation
- Review for accuracy
- Generate godoc documentation

##### **Step 4.3: Deployment Package Creation**
- Build release artifacts
- Create Docker images if applicable
- Prepare deployment scripts

##### **Step 4.4: Demo Material Preparation**
- Create demo scenarios
- Prepare presentation materials
- Test demo flow

##### **Step 4.5: Acceptance Test Execution**
- Run acceptance test suite
- Validate all acceptance criteria
- Document test results

##### **Step 4.6: Release Notes Creation**
- Document implemented features
- List known issues
- Include upgrade instructions

#### **Daily Recurring Activities**
- Morning standup (status update)
- Update story status on board
- Address and document blockers
- Commit and push code changes
- Run test suite before EOD

#### **Phase 3 Completion Gate** 🔴 USER APPROVAL REQUIRED
**Claude Action:**
```
"Implementation Phase Complete!
Results:
✓ X stories completed
✓ Test coverage: XX%
✓ All tests passing
✓ Code reviews approved
✓ Documentation updated
✓ Demo ready

Approve transition to Review phase?"
```
**Wait for:** User approval

---

### **PHASE 4: REVIEW (3 DAYS)**

#### **Day 1: Stakeholder Review**

##### **Step 1.1: Phase Initialization** 🚀 START HERE
- Create `cycles/cycle-XXX/review/` directory
- Update `.claude/current_status` to Review
- Prepare review materials

##### **Step 1.2: Stakeholder Demo** 🔴 USER INTERACTION REQUIRED
**Claude Action:**
```
"Ready to demo the implemented features.
Would you like to see the working implementation?"
[Conduct demo]
```
**Document:** Demo feedback in `review/demo-feedback.md`

##### **Step 1.3: Feedback Collection** 🔴 USER INTERACTION REQUIRED
**Claude Action:**
```
"Please provide feedback on the completed stories:
- What works well?
- What needs improvement?
- Any missing functionality?"
```
**Document:** Create `review/stakeholder-feedback.md`
**Continue when:** Feedback discussion complete

##### **Step 1.4: Enhancement Documentation**
- Document requested enhancements
- Prioritize for future cycles
- Create `review/enhancement-requests.md`

##### **Step 1.5: Acceptance Criteria Review**
- Verify all criteria met
- Document any gaps
- Update `review/acceptance-review.md`

##### **Step 1.6: Satisfaction Score Collection**
- Gather user satisfaction ratings
- Document scores and comments
- Create `review/satisfaction-scores.md`

##### **Step 1.7: Critical Issues Identification** 🔴 USER INTERACTION REQUIRED
**Claude Action:**
```
"Are there any critical issues that need immediate attention?
Should we address these before closing the cycle?"
```
**Document:** Create `review/critical-issues.md`
**Continue when:** Issues discussed and documented

#### **Day 2: Metrics & Analysis**

##### **Step 2.1: Velocity Analysis**
- Calculate story points completed
- Compare to estimates
- Document in `review/velocity-metrics.md`

##### **Step 2.2: Test Coverage Review**
- Analyze coverage reports
- Identify gaps
- Create `review/test-coverage.md`

##### **Step 2.3: Performance Assessment**
- Review benchmark results
- Compare to targets
- Document in `review/performance-metrics.md`

##### **Step 2.4: Completion Rate Calculation**
- Calculate story completion percentage
- Analyze incomplete items
- Create `review/completion-metrics.md`

##### **Step 2.5: Technical Debt Documentation**
- Identify accumulated debt
- Prioritize for addressing
- Create `review/technical-debt.md`

##### **Step 2.6: Cycle Metrics Report**
- Compile all metrics
- Create comprehensive `review/cycle-metrics-report.md`

#### **Day 3: Retrospective & Planning**

##### **Step 3.1: Team Retrospective**
- Conduct retrospective session
- Identify what worked well
- Discuss improvements needed
- Document in `review/retrospective.md`

##### **Step 3.2: Lessons Learned Documentation**
- Compile key learnings
- Create actionable insights
- Document in `review/lessons-learned.md`

##### **Step 3.3: Process Improvement Identification**
- Identify process bottlenecks
- Propose improvements
- Create `review/process-improvements.md`

##### **Step 3.4: Next Cycle Action Items**
- Create action item list
- Assign priorities
- Document in `review/action-items.md`

##### **Step 3.5: Deliverable Archival**
- Archive all cycle artifacts
- Organize documentation
- Update cycle index

##### **Step 3.6: Next Cycle Topic Preparation**
- Identify potential topics
- Prioritize based on feedback
- Create `review/next-cycle-topics.md`

#### **Phase 4 & Cycle Completion Gate** 🔴 USER APPROVAL REQUIRED
**Claude Action:**
```
"Review Phase and Cycle-XXX Complete!
Summary:
✓ Stakeholder feedback collected
✓ Metrics analyzed
✓ Retrospective completed
✓ Next cycle topics identified
✓ All deliverables archived

Approve cycle completion and start new cycle?"
```
**Wait for:** User approval for cycle completion

---

## 🔄 Section 5: Phase Tracking & Status Management

### **Phase Transition Workflow**

#### **IMPORTANT: User Approval Required**
**Phase transitions MUST be explicitly approved by the user. Claude cannot autonomously move between phases or cycles without user confirmation.**

#### **Phase Completion Protocol:**
1. **Claude Notifies User**: When phase deliverables are complete, Claude must inform the user and request approval to proceed
2. **User Reviews Deliverables**: User evaluates whether phase objectives have been met
3. **User Approves Transition**: User explicitly states "Approve phase transition" or similar confirmation
4. **Claude Executes Transition**: Only after approval, Claude updates status and begins next phase

#### **Example Interaction:**
```
Claude: "Research phase deliverables are complete:
- Problem statement ✓
- Feasibility report ✓  
- Risk assessment ✓
Ready to transition to Planning phase. Please review and approve."

User: "Approved. Begin Planning phase."

Claude: [Updates status file and creates planning directories]
```

### **Status File System**
Claude uses a combination of status file and directory structure to understand the current project phase:

#### **Status File Location**
```
C:\Users\marko\IdeaProjects\context-extender\.claude\current_status
```

#### **Status File Format:**
```
CURRENT_CYCLE: cycle-001
CURRENT_PHASE: research
PHASE_START_DATE: 2024-01-15
PHASE_END_DATE: 2024-01-18
CYCLE_GOAL: Define context parsing requirements and architecture
NEXT_MILESTONE: Technical design completion
PHASE_PROGRESS: Day 2 of 3
BLOCKERS: None
NOTES: Initial feasibility assessment complete
```

#### **Directory Structure Indicators**
The presence and content of cycle directories indicate phase status:

**Completed Phase**: Contains deliverable files and artifacts
```
cycles/cycle-001/research/
├── problem-statement.md      # Present = Research completed
├── feasibility-report.md
└── risk-assessment.md
```

**Current Phase**: Contains in-progress work files
```
cycles/cycle-001/planning/
├── technical-design.md       # In progress
├── user-stories.md          # Draft
└── wip/                     # Work in progress folder
```

**Future Phase**: Empty directory or non-existent
```
cycles/cycle-001/implementation/   # Empty = Not started
cycles/cycle-001/review/           # Empty = Not started
```

#### **When Starting a New Phase (After User Approval):**
1. **Confirm User Approval**: Ensure explicit approval has been given
2. **Update Status File**: Change CURRENT_PHASE, dates, and goals
3. **Create Phase Directory**: If it doesn't exist
4. **Archive Previous Phase**: Ensure previous phase deliverables are complete
5. **Initialize Phase Templates**: Copy from `cycles/templates/[phase]/`

#### **When Starting a New Cycle (After User Approval):**
1. **Confirm User Approval**: User must explicitly approve cycle completion and new cycle start
2. **Complete Review Phase**: Archive all cycle deliverables
3. **Update Cycle Counter**: Increment to cycle-002, cycle-003, etc.
4. **Reset to Research Phase**: New cycle always starts with Research
5. **Create New Cycle Directory Structure**:
```bash
mkdir -p cycles/cycle-002/{research,planning,implementation,review}
```

### **Status Checking for Claude**
Claude will check status using this priority order:

1. **Read Status File**: Primary source of truth for current phase
2. **Check Directory Structure**: Validate phase progress through file presence
3. **Ask User if Unclear**: Request clarification if status is ambiguous

### **Phase Status Indicators**

#### **Research Phase Status:**
- **Not Started**: No research directory or empty
- **In Progress**: Research directory exists with WIP files
- **Complete**: All research deliverables present and finalized

#### **Planning Phase Status:**
- **Not Started**: No planning directory
- **In Progress**: Planning directory with draft documents
- **Complete**: Technical design approved, user stories finalized

#### **Implementation Phase Status:**
- **Not Started**: No implementation artifacts
- **In Progress**: Code commits, active branches, failing/passing tests
- **Complete**: All acceptance criteria met, code reviewed, tests passing

#### **Review Phase Status:**
- **Not Started**: No review artifacts
- **In Progress**: Partial feedback, incomplete retrospective
- **Complete**: All deliverables reviewed, retrospective complete, next cycle prepared

---

## 📚 Section 6: Templates & Reference Information

### **User Interaction Points Summary**
- **Research Phase**: 5 interaction points
- **Planning Phase**: 6 interaction points  
- **Implementation Phase**: 1 interaction point (approval)
- **Review Phase**: 4 interaction points
- **Total**: 16 user interaction points per cycle

### **Cross-Phase Workflows**

#### **Status Management Workflow**
1. At each phase start:
   - Update `.claude/current_status`
   - Record phase start date
   - Set phase end date target
   - Update cycle documentation

2. At each phase end:
   - Archive phase deliverables
   - Update completion status
   - Request user approval
   - Prepare next phase

#### **Communication Workflow**
1. Phase kickoff notification
2. Daily progress updates during Implementation
3. Blocker escalation as needed
4. Phase completion notification
5. Stakeholder updates at milestones

#### **Documentation Workflow**
1. Create phase directory structure
2. Use templates from `cycles/templates/`
3. Maintain consistent naming conventions
4. Update documents progressively
5. Archive completed phase docs

### **Integration with Existing Tools**
- **Version Control**: Git with feature branch workflow
- **CI/CD**: GitHub Actions with Go-specific pipelines
- **Project Management**: Compatible with Jira, Azure DevOps, or similar tools
- **Communication**: Slack/Teams integration for notifications
- **Documentation**: Markdown-based with automated generation support

### **Adaptation Guidelines**
- Cycle length can vary from 15-21 days based on increment complexity
- Phase overlap is permitted when appropriate (e.g., next cycle Research during Implementation)
- Process should be continuously refined based on retrospective feedback
- Scale approach as team size and project complexity grow

### **INVEST Criteria for Go Projects**
- **Independent**: Stories can be developed without blocking dependencies
- **Negotiable**: Implementation details flexible while maintaining user value
- **Valuable**: Each story delivers measurable user or business value
- **Estimable**: Team can estimate effort with Go development context
- **Small**: Fits within 7-day Implementation phase
- **Testable**: Clear pass/fail criteria with Go testing framework

### **Anti-patterns to Avoid**
- **Technical Stories**: "Refactor auth module" → Instead: "As a user, I want faster login"
- **Vague Criteria**: "Should work well" → Instead: "Response time < 200ms"
- **Solution-First**: "Add Redis cache" → Instead: "As a user, I want fast data retrieval"
- **Monster Stories**: Stories too large for single sprint
- **Orphan Stories**: Stories without clear user value

---

## 🚀 Getting Started

### **For New Claude Instances:**
1. **Initialize Status System**: Create first status file and cycle directory
2. Read and understand this complete guide
3. Check current project status in `.claude/current_status`
4. Follow the granular workflow steps
5. Always request user approval for phase transitions

### **For Continuing Development:**
1. Check `.claude/current_status` for current phase
2. Review cycle directory structure for progress
3. Follow workflow steps from current position
4. Maintain documentation standards
5. Scale and optimize for sustained development velocity

---

## 🤖 Section 7: SME Subagent System

### **Overview**
Claude has access to 4 specialized SME (Subject Matter Expert) subagents that provide expert guidance on complex decisions throughout the project lifecycle, plus 8 execution subagents organized into phase-specific teams.

### **When to Consult SME Subagents**
- Complex technical decisions
- Quality standard questions
- Risk assessment needs
- Process optimization opportunities
- Architecture design choices
- Performance optimization
- Security implementation

### **Available SME Subagents**

#### **Technical Governance SME (Enhanced with Specialized Competencies)**
**Primary File**: `.claude/sme/technical_governance.md`
**Specialized Competency Files**:
- `.claude/sme/go_language_specialist.md` - Advanced Go patterns and optimization
- `.claude/sme/cli_development_specialist.md` - CLI design and user experience
- `.claude/sme/claude_code_specialist.md` - Claude Code extension patterns
**Consult for**:
- Architecture and design patterns
- Technology stack selection  
- Performance optimization strategies
- Security implementation approaches
- Technical debt management
- Go development best practices (Go Language Specialist competency)
- CLI tool design and user experience (CLI Development Specialist competency)
- Claude Code integration and extension patterns (Claude Code Specialist competency)

**Consultation Protocol**:
1. For Go-specific decisions: Consult Go Language Specialist → Technical SME validates
2. For CLI-specific decisions: Consult CLI Development Specialist → Technical SME validates  
3. For Claude Code integration: Consult Claude Code Specialist → Technical SME validates
4. For cross-cutting decisions: Consult relevant specialists → Technical SME synthesizes

#### **Quality Governance SME**  
**File**: `.claude/sme/quality_governance.md`
**Consult for**:
- Quality standard definitions
- Testing strategy development
- Code review processes
- Documentation requirements
- Quality gate criteria
- Defect management

#### **Risk Governance SME**
**File**: `.claude/sme/risk_governance.md`
**Consult for**:
- Risk identification and assessment
- Mitigation strategy development
- Risk escalation decisions
- Contingency planning
- Risk monitoring approaches

#### **Process Governance SME**
**File**: `.claude/sme/process_governance.md`
**Consult for**:
- Process deviation requests
- Workflow optimizations
- Phase gate criteria
- Retrospective action items
- Process metrics analysis

### **How to Use Subagents - Practical Guide**

#### **Step 1: Determine Current Phase**
```markdown
Check cycle day → Identify phase → Activate appropriate team
- Days 1-2: Research Phase → Architecture Discovery Specialist
- Days 3-4: Planning Phase → Story Refinement + Implementation Planning
- Days 5-15: Implementation Phase → Full 5-subagent team
- Days 16-17: Review Phase → SME-led (no new subagents)
```

#### **Step 2: Activate Phase-Appropriate Subagents**

**For SME Consultation (Any Phase):**
```markdown
1. Use Task tool with subagent_type: "general-purpose"
2. Copy consultation prompt from SME file
3. Example:
   "As the Technical Governance SME, evaluate [specific decision].
   Consider: [factors]
   Provide: [recommendations]"
```

**For Execution Subagents (Phase-Specific):**
```markdown
Research Phase Example:
"As the Architecture Discovery Specialist, explore architecture options for 
context manipulation in our Go-based CLI tool that extends Claude Code.
Consider: Performance, scalability, Claude Code integration patterns
Generate: ADRs with trade-offs and recommendations"

Planning Phase Example:
"As the Story Refinement Specialist, break down the epic 'Context Enhancement 
Features' into INVEST-compliant user stories with Go technical details"

Implementation Phase Example:
"As the Test Automation Specialist, generate comprehensive unit tests for 
the context manipulation functions in package context_handler"
```

#### **Step 3: Subagent Interaction Patterns**

**Sequential Consultation:**
```markdown
1. Architecture Discovery → generates ADRs
2. Story Refinement → uses ADRs to create stories
3. Implementation Planning → estimates stories
4. Test Automation → generates tests from stories
```

**Parallel Consultation (Same Phase):**
```markdown
During Implementation:
- Test Automation + Code Quality (run simultaneously)
- Integration Orchestrator + Progress Tracker (continuous)
- All feed → Knowledge Curator
```

**Cross-Phase Handoffs:**
```markdown
Research outputs → Planning inputs
Planning outputs → Implementation inputs
Implementation data → Review analysis
```

#### **Step 4: Common Usage Scenarios**

**Scenario 1: Starting a New Cycle**
```markdown
Day 1 Action:
1. Activate Architecture Discovery Specialist
2. Consult Technical SME for constraints
3. Generate ADRs for key decisions
4. Document in Knowledge Curator
```

**Scenario 2: Complex Technical Decision**
```markdown
Any Phase Action:
1. Consult Go Language Specialist (specific expertise)
2. Consult Claude Code Specialist (integration)
3. Technical SME synthesizes recommendations
4. User approval if architectural impact
```

**Scenario 3: Story Implementation**
```markdown
Day 5-15 Action:
1. Test Automation generates tests
2. Code Quality enforces standards real-time
3. Integration Orchestrator validates integration
4. Progress Tracker monitors velocity
5. Knowledge Curator captures patterns
```

#### **Decision Tree for Subagent Usage**
```
Is it a strategic decision? → Consult SME → User approval
Is it phase-specific work? → Use phase team subagents
Is it technical expertise? → Consult specialist → SME validates
Is it execution task? → Use implementation subagents
Is it review/retrospective? → SME-led with data support
```

### **Phase-Based Subagent Teams**

#### **Research Team (Days 1-2)**
- **Lead**: Architecture Discovery Specialist
- **Support**: Technical SMEs and Specialists
- **See**: `.claude/subagents/subagent_teams.md` for activation protocols

#### **Planning Team (Days 3-4)**
- **Lead**: Story Refinement Specialist
- **Support**: Implementation Planning Orchestrator
- **See**: `.claude/subagents/planning_framework.md` for workflows

#### **Implementation Team (Days 5-15)**
- **Active**: Test Automation, Code Quality, Integration, Progress, Knowledge
- **See**: `.claude/subagents/README.md` for coordination patterns

#### **Review Team (Days 16-17)**
- **Lead**: Process Governance SME (human-driven)
- **Data**: Progress Tracker and Knowledge Curator (read-only)
- **See**: `.claude/subagents/review_phase_strategy.md` for approach

**IMPORTANT**: Use `.claude/subagents/subagent_teams.md` to determine which subagents to activate for the current phase.

### **Common Subagent Usage Mistakes to Avoid**

#### **❌ Don't Do This:**
- **Mix phase teams**: Don't use planning subagents during implementation
- **Skip handoffs**: Always transfer artifacts between phases properly
- **Bypass SMEs**: Don't skip SME validation for major decisions
- **Parallel incompatible subagents**: Don't run conflicting subagents simultaneously
- **Ignore dependencies**: Respect input/output dependencies between subagents

#### **✅ Do This Instead:**
- **Respect phase boundaries**: Use only phase-appropriate subagents
- **Clean handoffs**: Ensure artifacts transfer between phase teams
- **SME oversight**: Maintain governance throughout all phases
- **Smart parallelization**: Only run compatible subagents in parallel
- **Follow dependencies**: Execute subagents in proper sequence

### **Quick Subagent Reference**

| If You Need To... | Use This Subagent | In Phase |
|-------------------|-------------------|----------|
| Explore architecture | Architecture Discovery Specialist | Research |
| Break down epics | Story Refinement Specialist | Planning |
| Estimate stories | Implementation Planning Orchestrator | Planning |
| Generate tests | Test Automation Specialist | Implementation |
| Check code quality | Code Quality Enforcer | Implementation |
| Test integrations | Integration Orchestrator | Implementation |
| Track progress | Progress Tracker and Reporter | Implementation |
| Capture knowledge | Knowledge Curator | All phases |
| Facilitate retrospective | Process Governance SME | Review |

---

## ✅ Section 8: Checklists and Quality Gates

### **Master Checklist Framework**
**File**: `.claude/checklists/MASTER_CHECKLIST.md`

The project uses comprehensive checklists to ensure consistency, quality, and completeness across all phases. Checklists are mandatory for phase transitions and quality validation.

### **How to Use Checklists - Step by Step**

#### **1. Identify Current Phase**
```markdown
Check cycle day (1-17) → Determine phase → Open appropriate checklist
Example: Day 7 → Implementation Phase → Open implementation_phase_checklist.md
```

#### **2. Complete Entry Requirements**
Before starting any phase:
- [ ] Verify previous phase deliverables
- [ ] Get user approval for phase transition
- [ ] Activate appropriate subagent team
- [ ] Set up phase-specific environment

#### **3. Follow Daily Activities**
Each phase has daily tasks:
- Use checklist as daily guide
- Mark items complete as you go
- Note blockers immediately
- Update progress metrics

#### **4. Validate Quality Gates**
Before phase exit:
- Run quality gate checks
- Ensure all mandatory items pass
- Document any exceptions
- Get remediation plan for failures

### **Available Checklists**

#### **Phase Checklists** (Entry → Daily → Exit)
| Phase | Checklist | Key Focus | User Touchpoints |
|-------|-----------|-----------|------------------|
| Research (Days 1-2) | `research_phase_checklist.md` | Problem definition, architecture exploration | 5 interactions |
| Planning (Days 3-4) | `planning_phase_checklist.md` | Story refinement, estimation | 6 interactions |
| Implementation (Days 5-15) | `implementation_phase_checklist.md` | Development, testing, quality | 1 interaction |
| Review (Days 16-17) | `review_phase_checklist.md` | Demo, retrospective, metrics | 4 interactions |

#### **Process Checklists**
| Checklist | Purpose | When to Use |
|-----------|---------|-------------|
| `quality_gates_checklist.md` | Pass/fail criteria | Before phase transitions |
| `subagent_activation_checklist.md` | Team coordination | Phase starts/transitions |
| `MASTER_CHECKLIST.md` | Checklist coordination | Daily reference |

### **Critical Quality Gates** 🔴

#### **Cannot Proceed Without**
```markdown
MANDATORY (Block phase transition if failed):
✓ Test coverage >80%
✓ Zero critical security vulnerabilities  
✓ 100% GoDoc coverage for public APIs
✓ User approval at phase transitions
✓ All CI/CD builds passing

ESCALATION REQUIRED:
- Coverage <80% → Quality SME consultation
- Security issues → Risk SME + User notification
- Failed builds → Technical SME investigation
```

### **Checklist Usage by Role**

#### **Claude's Daily Checklist Routine**
```markdown
Morning:
1. Open current phase checklist
2. Review today's activities
3. Check subagent activation status
4. Verify quality gate progress

During Day:
- Mark activities complete
- Note decisions in Knowledge Curator
- Update Progress Tracker
- Check SME consultation triggers

End of Day:
- Update checklist progress
- Document blockers
- Prepare for next day
- Check if phase transition approaching
```

#### **User Interaction Checkpoints**
The checklists mark 16 mandatory user interactions:
- Research: 5 touchpoints
- Planning: 6 touchpoints  
- Implementation: 1 touchpoint
- Review: 4 touchpoints

### **Common Checklist Patterns**

#### **Phase Transition Checklist**
```markdown
Current Phase Exit:
□ All deliverables complete
□ Quality gates passed
□ Documentation updated
□ User approval obtained

Next Phase Entry:
□ Handoff artifacts received
□ Subagent team activated
□ Environment prepared
□ Objectives understood
```

#### **Daily Stand-up Checklist**
```markdown
□ Yesterday's progress reviewed
□ Today's objectives identified
□ Blockers documented
□ Metrics updated
□ Team coordination confirmed
```

### **Checklist Compliance Tracking**

| Compliance Level | Required Rate | Items |
|-----------------|---------------|-------|
| **Mandatory** | 100% | User approvals, security, phase gates |
| **Essential** | >95% | Quality standards, documentation |
| **Recommended** | >80% | Process items, templates |

### **Quick Reference: Which Checklist When?**

```markdown
Starting new cycle? → research_phase_checklist.md
Changing phases? → Current phase exit + Next phase entry checklists
Daily work? → Current phase daily activities section
Quality check? → quality_gates_checklist.md
Activating subagents? → subagent_activation_checklist.md
Lost? → MASTER_CHECKLIST.md
```

---

## 📄 Section 9: Documentation Requirements

### **Core Principle: Real-Time Documentation**
**Documentation is created as work happens, not after work completes.**

All phase deliverables must be documented to files immediately when created. This ensures knowledge capture, enables collaboration, and prevents loss of context between sessions.

### **Documentation Standards**

#### **File Naming Convention**
```
.claude/cycles/cycle-XXX/[phase]/[document_name].md
```

#### **Required Templates**
- Use consistent markdown formatting
- Include date, phase, and cycle information in headers
- Add cross-references to related documents
- Maintain version history for significant changes

#### **Quality Requirements**
- Documents must be complete and self-contained
- Include context for future readers (assume no prior knowledge)
- Link to relevant code, decisions, or external resources
- Use clear, professional language suitable for stakeholders

### **Phase-Specific Documentation Requirements**

#### **Research Phase Documentation (Days 1-2)**

**Day 1 - Problem Definition Session**:
```markdown
IMMEDIATELY after user problem session:
→ Create: cycles/cycle-XXX/research/problem_definition.md
→ Include: Problem statement, use cases, success criteria, constraints
→ Template: Use standardized problem definition template
```

**Day 1 - Architecture Exploration**:
```markdown
IMMEDIATELY after SME consultation completes:
→ Create: cycles/cycle-XXX/research/architecture_decision_records.md
→ Include: All ADRs from Architecture Discovery Specialist
→ Format: Standard ADR format (Status, Context, Decision, Consequences)
```

**Day 2 - Technical Feasibility**:
```markdown
IMMEDIATELY after feasibility analysis completes:
→ Create: cycles/cycle-XXX/research/technical_feasibility_analysis.md
→ Include: Challenge analysis, implementation patterns, risk assessment
→ Include: Performance targets, dependencies, testing strategy
```

**Day 2 - Risk Assessment**:
```markdown
IMMEDIATELY after risk analysis completes:
→ Create: cycles/cycle-XXX/research/risk_register.md
→ Include: Risk matrix, mitigation strategies, monitoring plans
→ Format: Standardized risk register with scoring
```

**Research Phase Exit Requirements**:
- [ ] problem_definition.md complete and validated
- [ ] architecture_decision_records.md with all ADRs documented
- [ ] technical_feasibility_analysis.md with confidence rating
- [ ] risk_register.md with mitigation plans
- [ ] All documents cross-referenced and internally consistent

#### **Planning Phase Documentation (Days 3-4)**

**Day 3 - Story Creation**:
```markdown
IMMEDIATELY after Story Refinement Specialist completes:
→ Create: cycles/cycle-XXX/planning/user_stories.md
→ Include: All user stories with acceptance criteria
→ Include: INVEST validation, dependencies, story points
```

**Day 3 - Epic Breakdown**:
```markdown
SIMULTANEOUSLY with story creation:
→ Update: cycles/cycle-XXX/planning/epic_breakdown.md
→ Include: Epic decomposition rationale, scope decisions
→ Include: Story prioritization and cycle allocation
```

**Day 4 - Implementation Planning**:
```markdown
IMMEDIATELY after Implementation Planning Orchestrator completes:
→ Create: cycles/cycle-XXX/planning/implementation_plan.md
→ Include: Daily schedule, risk mitigation, quality gates
→ Include: Resource allocation, dependencies, success metrics
```

**Day 4 - Sprint Backlog**:
```markdown
BEFORE Planning Phase exit:
→ Create: cycles/cycle-XXX/planning/sprint_backlog.md
→ Include: Prioritized stories, estimation rationale
→ Include: Definition of Ready, Definition of Done
```

**Planning Phase Exit Requirements**:
- [ ] user_stories.md with all stories and acceptance criteria
- [ ] epic_breakdown.md with scope decisions documented
- [ ] implementation_plan.md with daily schedule
- [ ] sprint_backlog.md ready for implementation
- [ ] All estimates validated and realistic

#### **Implementation Phase Documentation (Days 5-15)**

**Daily Requirements**:
```markdown
END OF EACH DAY:
→ Update: cycles/cycle-XXX/implementation/daily_progress.md
→ Include: Completed tasks, blockers, tomorrow's plan
→ Include: Risk updates, quality metrics, time tracking
```

**Weekly Requirements**:
```markdown
END OF EACH WEEK:
→ Create: cycles/cycle-XXX/implementation/weekly_retrospective_[week].md
→ Include: Lessons learned, process improvements
→ Include: Velocity tracking, quality metrics
```

**Story Completion**:
```markdown
WHEN EACH STORY COMPLETES:
→ Update: cycles/cycle-XXX/implementation/story_completion_log.md
→ Include: Actual vs estimated effort, lessons learned
→ Include: Quality metrics, testing results
```

**Implementation Phase Exit Requirements**:
- [ ] daily_progress.md updated through final day
- [ ] weekly_retrospective_[X].md for each week
- [ ] story_completion_log.md with all completed stories
- [ ] implementation_summary.md with final status
- [ ] code_quality_report.md with metrics and coverage

#### **Review Phase Documentation (Days 16-17)**

**Day 16 - Demonstration**:
```markdown
AFTER demo session:
→ Create: cycles/cycle-XXX/review/demo_feedback.md
→ Include: User feedback, feature validation results
→ Include: Performance metrics, usability observations
```

**Day 17 - Retrospective**:
```markdown
AFTER retrospective session:
→ Create: cycles/cycle-XXX/review/retrospective_analysis.md
→ Include: What worked well, what didn't, improvements
→ Include: Process effectiveness, team satisfaction
```

**Day 17 - Cycle Completion**:
```markdown
BEFORE cycle closure:
→ Create: cycles/cycle-XXX/review/cycle_summary.md
→ Include: Overall success metrics, lessons learned
→ Include: Recommendations for next cycle
```

**Review Phase Exit Requirements**:
- [ ] demo_feedback.md with user validation results
- [ ] retrospective_analysis.md with process insights
- [ ] cycle_summary.md with complete cycle assessment
- [ ] next_cycle_recommendations.md for planning input

### **Documentation Automation Rules**

#### **SME/Subagent Output Capture**
```markdown
MANDATORY: When using Task tool for SME consultation:
1. Immediately save SME output to appropriate document
2. Include SME name, consultation date, and context
3. Format output for readability and future reference
4. Cross-reference with related decisions or documents
```

#### **Phase Transition Gates**
```markdown
NO PHASE TRANSITION without:
1. All required documents created and complete
2. Document quality review completed
3. Cross-references validated
4. User approval on documented deliverables
```

#### **Document Maintenance**
```markdown
ONGOING responsibilities:
1. Update documents when decisions change
2. Maintain version history for significant changes
3. Ensure documents remain current and accurate
4. Regular document quality audits
```

### **Documentation Quality Gates**

#### **Completeness Checklist**
- [ ] All required sections present
- [ ] Context provided for future readers
- [ ] Decisions explained with rationale
- [ ] Cross-references working and relevant
- [ ] Templates followed consistently

#### **Clarity Standards**
- [ ] Professional language appropriate for stakeholders
- [ ] Technical details balanced with accessibility
- [ ] Clear action items and next steps
- [ ] Consistent terminology throughout

#### **Accuracy Validation**
- [ ] Information current and up-to-date
- [ ] Technical details verified
- [ ] Dependencies and relationships correct
- [ ] Quality metrics accurate

### **Document Templates and Examples**

Templates are available in:
- `.claude/cycles/templates/` - Standard document templates
- `.claude/cycles/cycle-001/` - Reference examples from Cycle 1
- `.claude/processes/documentation_standards.md` - Detailed formatting guidelines

### **Compliance and Monitoring**

#### **Daily Checks**
- Document creation requirements met
- Template usage consistent
- Quality standards maintained

#### **Phase Exit Reviews**
- All required documents present
- Quality gates passed
- Cross-references validated
- User approval obtained

#### **Cycle Retrospectives**
- Documentation effectiveness assessment
- Template improvements identified
- Process refinements recommended

**Remember**: Documentation is not overhead—it's the foundation that enables effective collaboration, knowledge transfer, and continuous improvement.

---

## 📊 Section 10: Established Processes

### **Process Overview**
The project uses 5 core processes integrated into the 4-phase cycle:

### **1. Risk Management Process**
**File**: `.claude/processes/risk_management.md`

**Integration Points**:
- **Research Phase**: Risk identification and assessment
- **Planning Phase**: Mitigation strategy development
- **Implementation Phase**: Risk monitoring and issue tracking
- **Review Phase**: Risk effectiveness evaluation

**Key Components**:
- Risk register template with scoring matrix
- Escalation protocols (Critical/High/Medium/Low)
- Risk categories (Technical/Project/External)
- Mitigation tracking and effectiveness measurement

### **2. Stakeholder Communication Process**
**File**: `.claude/processes/stakeholder_communication.md`

**Integration Points**:
- **Research Phase**: Problem validation and requirements gathering
- **Planning Phase**: Technical approach validation and story refinement
- **Implementation Phase**: Progress updates and scope adjustments
- **Review Phase**: Demo, feedback collection, and satisfaction measurement

**Key Components**:
- Communication matrix by phase and stakeholder type
- RACI framework for decision involvement
- Feedback collection templates and satisfaction metrics
- Scalable from single stakeholder to multiple stakeholders

### **3. Project Governance Process**
**File**: `.claude/processes/project_governance.md`

**Integration Points**:
- **All Phases**: Decision-making authority and approval gates
- **Phase Transitions**: Governance review and approval
- **Cross-Phase**: Escalation handling and conflict resolution

**Key Components**:
- Decision authority matrix (Strategic/Technical/Operational/Administrative)
- Governance gates and approval criteria
- Escalation matrix and response protocols
- Decision log and accountability tracking

### **4. CI/CD Pipeline Process**
**File**: `.claude/processes/cicd_pipeline.md`

**Integration Points**:
- **Planning Phase**: Pipeline configuration and quality gates
- **Implementation Phase**: Continuous integration and deployment
- **Review Phase**: Pipeline metrics and deployment validation

**Key Components**:
- Complete GitHub Actions workflows
- Quality gates (Build/Test/Security/Lint)
- Branch protection and deployment strategies
- Feature flag management system
- Pipeline monitoring and notifications

### **5. Enhanced Feedback Loops Process**
**File**: `.claude/processes/feedback_loops.md`

**Integration Points**:
- **All Phases**: Multi-level feedback collection and response
- **Continuous**: Immediate, daily, weekly, and cycle feedback
- **Strategic**: Quarterly business reviews and trend analysis

**Key Components**:
- 6-level feedback architecture (Immediate → Daily → Weekly → Phase → Cycle → Quarterly)
- Automated metrics collection and trend analysis
- Weekly micro-retrospectives and improvement experiments
- Feedback response protocols and escalation triggers

### **Process Compliance Requirements**

#### **Mandatory Processes**
Must be followed without exception:
- Phase transition gates with user approval
- Risk escalation for Critical/High risks
- Security review process
- Documentation standards

#### **Recommended Processes**
Should be followed unless justified deviation:
- Code review practices
- Testing strategies
- Communication protocols
- Metrics collection

### **Process Integration Matrix**

| Process | Research | Planning | Implementation | Review |
|---------|----------|-----------|----------------|---------|
| **Risk Management** | Identification | Mitigation Planning | Monitoring | Effectiveness Review |
| **Communication** | Problem Validation | Design Review | Progress Updates | Demo & Feedback |
| **Governance** | Alignment Check | Technical Approval | Operational Oversight | Gate Review |
| **CI/CD** | - | Configuration | Continuous Integration | Metrics Analysis |
| **Feedback** | Initial Setup | Enhancement Planning | Multi-level Collection | Comprehensive Review |

---

## 🔄 Section 9: Process Evolution and Improvement

### **Continuous Process Improvement**
- Processes are reviewed and improved each cycle
- SME subagents provide ongoing optimization recommendations
- User feedback drives process refinements
- Metrics and effectiveness data inform improvements

### **Process Maturity Roadmap**
- **Current State**: Defined processes with basic automation
- **6-Month Goal**: Managed processes with comprehensive metrics
- **12-Month Goal**: Optimized processes with predictive capabilities

This complete guide with integrated SME subagents and established processes ensures Claude can make informed decisions, follow proven workflows, and continuously improve project delivery while maintaining quality and stakeholder satisfaction.