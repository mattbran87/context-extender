# Research Phase Subagent Framework

## Overview
The context-extender project benefits from specialized research phase subagents that automate architecture exploration, feasibility assessment, and knowledge discovery during the critical 2-day research phase. These subagents provide systematic investigation capabilities that complement SME strategic guidance.

## Research Phase Characteristics

### Why Research Subagents Are Valuable
- **Different Skills**: Research requires exploration and discovery vs execution
- **Systematic Analysis**: Automated comparison frameworks and evaluation matrices
- **Knowledge Synthesis**: Structured capture of discoveries and patterns
- **Decision Support**: Evidence-based recommendations with clear rationale
- **Handoff Optimization**: Structured outputs for planning phase consumption

## Recommended Research Subagents

### 🔴 Priority 1: Essential Research Automation (Recommended for Implementation)

#### 1. Architecture Discovery Specialist (Created)
**File**: `.claude/subagents/architecture_discovery_specialist.md`
- **Purpose**: Systematic architecture exploration and technology evaluation
- **Key Value**: Comprehensive option analysis with ADR generation
- **Integration**: Works with all Technical Specialists and feeds planning phase

#### 2. Feasibility Assessment Orchestrator (Future)
**File**: `.claude/subagents/feasibility_assessment_orchestrator.md`
- **Purpose**: Multi-dimensional feasibility analysis and synthesis
- **Key Value**: Holistic feasibility scoring beyond individual SME perspectives
- **Integration**: Orchestrates all SME inputs for comprehensive assessment

#### 3. Knowledge Discovery Agent (Future)
**File**: `.claude/subagents/knowledge_discovery_agent.md`
- **Purpose**: External solution research and best practice identification
- **Key Value**: Systematic knowledge gathering and anti-pattern documentation
- **Integration**: Feeds Knowledge Curator for cross-cycle learning

### 🟡 Priority 2: High-Value Research Enhancement (Future Consideration)

#### 4. Requirements Synthesis Specialist
- **Purpose**: Multi-stakeholder requirement consolidation and conflict resolution
- **Key Value**: Enhanced requirement quality before planning phase
- **Integration**: Direct handoff to Story Refinement Specialist

#### 5. Proof of Concept Automation Agent
- **Purpose**: POC planning, execution, and evaluation automation
- **Key Value**: Accelerated technical validation with standardized frameworks
- **Integration**: Works with Technical Specialists for rapid prototyping

### 🟢 Priority 3: Process Optimization (Long-term)

#### 6. Research Metrics Analyst
- **Purpose**: Research phase velocity and effectiveness tracking
- **Key Value**: Data-driven research process improvement
- **Integration**: Feeds Progress Tracker for comprehensive metrics

## Complete Subagent Lifecycle Coverage

```
RESEARCH PHASE (2 days) - ENHANCED
├── Architecture Discovery Specialist ✓
├── Feasibility Assessment Orchestrator (future)
├── Knowledge Discovery Agent (future)
└── SME Strategic Consultation

PLANNING PHASE (3 days) - COMPLETE
├── Story Refinement Specialist ✓
├── Implementation Planning Orchestrator ✓
└── Future: Architecture Decision Advisor, Acceptance Criteria Generator

IMPLEMENTATION PHASE (11 days) - COMPLETE
├── Test Automation Specialist ✓
├── Code Quality Enforcer ✓
├── Integration Orchestrator ✓
├── Progress Tracker and Reporter ✓
└── Knowledge Curator ✓

REVIEW PHASE (3 days) - SME DRIVEN
└── Retrospectives with subagent data inputs
```

## Context-Extender Specific Benefits

### Architecture Discovery Benefits
- **Go Architecture Patterns**: Systematic evaluation of Go-specific patterns
- **CLI Design Exploration**: Comprehensive CLI architecture analysis
- **Claude Code Integration**: Deep analysis of extension architectures
- **Performance Analysis**: Early performance architecture validation
- **Risk Identification**: Proactive architectural risk discovery

### Research-to-Planning Handoff
```
Architecture Discovery Specialist
├── ADRs → Architecture Decision Advisor (planning)
├── Technical Constraints → Story Refinement Specialist
├── Complexity Assessment → Implementation Planning Orchestrator
└── Discovered Patterns → Knowledge Curator

Feasibility Assessment Orchestrator
├── Feasibility Scores → Story prioritization
├── Risk Assessment → Risk mitigation planning
├── Resource Requirements → Resource allocation
└── POC Results → Technical validation

Knowledge Discovery Agent
├── Best Practices → Implementation patterns
├── Anti-Patterns → Quality enforcement rules
├── External Solutions → Architecture decisions
└── Lessons Learned → Process improvements
```

## Implementation Strategy

### Phase 1: Foundation (Next Cycle)
**Deploy Architecture Discovery Specialist**
- Integrate with Technical SME and Specialists
- Establish ADR generation templates
- Connect to planning phase subagents
- **Expected Impact**: 30% improvement in architecture decision quality

### Phase 2: Enhancement (Future Cycles)
**Evaluate Additional Research Subagents**
- Monitor research phase effectiveness
- Identify remaining gaps and pain points
- Deploy Feasibility and Knowledge agents if warranted
- **Expected Impact**: 25% reduction in planning rework

### Phase 3: Optimization (Long-term)
**Research Process Intelligence**
- Deploy Research Metrics Analyst
- Implement continuous improvement cycle
- Cross-project knowledge transfer
- **Expected Impact**: Research phase optimization

## Integration with SME Framework

### Research Subagent Collaboration Model
```
Research Question
    ↓
Architecture Discovery Specialist (explores options)
    ↓
Technical Specialists (provide domain expertise)
    ↓
Technical Governance SME (validates approach)
    ↓
Risk Governance SME (assesses risks)
    ↓
Planning Subagents (receive structured outputs)
```

### Decision Authority
- **Research Subagents**: Explore, analyze, and recommend
- **SMEs**: Validate, govern, and decide
- **User**: Approve strategic decisions
- **Clear Boundaries**: Research subagents don't make decisions, they inform them

## Success Metrics

### Research Phase Effectiveness
- **Discovery Coverage**: > 90% of viable options identified
- **Decision Quality**: > 85% of research decisions validated in implementation
- **Research Velocity**: Complete research in 2-day window
- **Knowledge Capture**: 100% of discoveries documented
- **Planning Readiness**: 95% of research outputs directly usable in planning

### Long-term Value
- **Pattern Library Growth**: Continuous architectural pattern accumulation
- **Decision Accuracy**: Improving ADR quality over cycles
- **Risk Prevention**: Early identification preventing implementation issues
- **Knowledge Reuse**: 60% of research patterns reused in future cycles
- **Process Maturity**: Evolution from ad-hoc to systematic research

## Key Differentiators: Research vs Other Phases

### Research Phase Uniqueness
- **Exploration vs Execution**: Discovery and analysis vs implementation
- **Uncertainty Management**: High uncertainty requiring systematic exploration
- **Knowledge Synthesis**: Combining external and internal knowledge
- **Decision Support**: Evidence gathering for strategic decisions
- **Foundation Setting**: Establishing technical direction for entire cycle

### Why Specialized Subagents Help
- **Different Skillset**: Research requires different capabilities than coding
- **Systematic Approach**: Automated frameworks ensure comprehensive analysis
- **Consistency**: Standardized evaluation across cycles
- **Knowledge Accumulation**: Structured capture for future reference
- **Time Efficiency**: Parallel exploration of multiple options

## Recommendation Summary

### Immediate Action
1. **Deploy Architecture Discovery Specialist** for systematic architecture exploration
2. **Monitor Effectiveness** during next 2-3 cycles
3. **Evaluate Gaps** to determine need for additional research subagents

### Future Consideration
- **Feasibility Assessment Orchestrator** if feasibility analysis becomes bottleneck
- **Knowledge Discovery Agent** if external research needs increase
- **Requirements Synthesis** if requirement quality issues persist

### Success Indicators
- Improved architecture decision quality
- Reduced planning phase rework
- Accelerated research phase completion
- Enhanced knowledge accumulation
- Better research-to-planning handoff

The research phase subagent framework provides targeted automation for the unique challenges of technical exploration and discovery, complementing the comprehensive planning and implementation subagent coverage while maintaining the strategic guidance of the SME framework.