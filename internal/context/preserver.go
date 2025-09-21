package context

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// ContextPreserver manages context preservation during conversation compression
type ContextPreserver struct {
	CriticalContext map[string]interface{}
	Timestamp       time.Time
}

// CompressionSummary represents what should be preserved during compression
type CompressionSummary struct {
	// Core Project Context
	ProjectName        string                 `json:"project_name"`
	CurrentObjective   string                 `json:"current_objective"`
	CurrentPhase       string                 `json:"current_phase"`
	WorkingDirectory   string                 `json:"working_directory"`

	// Technical Decisions
	TechnicalStack     []string               `json:"technical_stack"`
	Constraints        []string               `json:"constraints"`
	KeyDecisions       map[string]string      `json:"key_decisions"`

	// User Preferences
	UserPreferences    map[string]string      `json:"user_preferences"`
	WorkflowStyle      string                 `json:"workflow_style"`
	CommunicationStyle string                 `json:"communication_style"`

	// Current State
	CompletedTasks     []string               `json:"completed_tasks"`
	PendingTasks       []string               `json:"pending_tasks"`
	ActiveProblems     []string               `json:"active_problems"`

	// Implementation Details
	KeyCodePatterns    map[string]string      `json:"key_code_patterns"`
	ErrorsToAvoid      []string               `json:"errors_to_avoid"`
	SuccessfulSolutions map[string]string     `json:"successful_solutions"`

	// Relationship Context
	TrustLevel         string                 `json:"trust_level"`
	ApprovalPatterns   []string               `json:"approval_patterns"`
	AvoidTopics        []string               `json:"avoid_topics"`
}

// ExtractCriticalContext analyzes conversation and extracts critical context
func ExtractCriticalContext(conversations []string) (*CompressionSummary, error) {
	summary := &CompressionSummary{
		KeyDecisions:        make(map[string]string),
		UserPreferences:     make(map[string]string),
		KeyCodePatterns:     make(map[string]string),
		SuccessfulSolutions: make(map[string]string),
	}

	// Analyze conversations for patterns
	for _, conv := range conversations {
		// Extract project context
		if strings.Contains(conv, "context-extender") || strings.Contains(conv, "Context-Extender") {
			summary.ProjectName = "Context-Extender"
		}

		// Extract technical decisions
		if strings.Contains(conv, "Pure Go") || strings.Contains(conv, "zero CGO") {
			summary.Constraints = append(summary.Constraints, "Zero CGO dependencies required")
			summary.KeyDecisions["database"] = "Pure Go SQLite using modernc.org/sqlite"
		}

		// Extract workflow preferences
		if strings.Contains(conv, "pragmatic") || strings.Contains(conv, "practical") {
			summary.WorkflowStyle = "Pragmatic and outcome-focused"
		}

		// Extract completed work
		if strings.Contains(conv, "v1.0.1") && strings.Contains(conv, "released") {
			summary.CompletedTasks = append(summary.CompletedTasks, "v1.0.1 release with critical fixes")
		}

		// Extract current objectives
		if strings.Contains(conv, "Cycle 5") {
			summary.CurrentObjective = "Planning Cycle 5"
			summary.CurrentPhase = "Between cycles - planning next iteration"
		}

		// Extract errors to avoid
		if strings.Contains(conv, "import cycle") {
			summary.ErrorsToAvoid = append(summary.ErrorsToAvoid, "Import cycles in Go packages")
		}

		// Extract user preferences from patterns
		if strings.Contains(conv, "simplified workflow") {
			summary.UserPreferences["workflow"] = "Simplified 5-day adaptive cycles"
			summary.UserPreferences["documentation"] = "Minimal, value-focused documentation"
		}
	}

	// Set defaults if not detected
	if summary.ProjectName == "" {
		summary.ProjectName = "Context-Extender CLI Tool"
	}

	if len(summary.TechnicalStack) == 0 {
		summary.TechnicalStack = []string{"Go", "SQLite", "Claude Code Hooks", "Cobra CLI"}
	}

	return summary, nil
}

// GenerateContextPrompt creates a prompt to reinject critical context
func GenerateContextPrompt(summary *CompressionSummary) string {
	var builder strings.Builder

	builder.WriteString("## Critical Context (Post-Compression)\n\n")

	// Project Context
	builder.WriteString(fmt.Sprintf("**Project**: %s\n", summary.ProjectName))
	builder.WriteString(fmt.Sprintf("**Current Objective**: %s\n", summary.CurrentObjective))
	builder.WriteString(fmt.Sprintf("**Phase**: %s\n\n", summary.CurrentPhase))

	// Technical Context
	builder.WriteString("### Technical Decisions\n")
	for key, value := range summary.KeyDecisions {
		builder.WriteString(fmt.Sprintf("- %s: %s\n", key, value))
	}
	builder.WriteString("\n")

	// Constraints
	if len(summary.Constraints) > 0 {
		builder.WriteString("### Constraints\n")
		for _, constraint := range summary.Constraints {
			builder.WriteString(fmt.Sprintf("- %s\n", constraint))
		}
		builder.WriteString("\n")
	}

	// User Preferences
	if len(summary.UserPreferences) > 0 {
		builder.WriteString("### User Preferences\n")
		for key, value := range summary.UserPreferences {
			builder.WriteString(fmt.Sprintf("- %s: %s\n", key, value))
		}
		builder.WriteString("\n")
	}

	// Current State
	if len(summary.CompletedTasks) > 0 {
		builder.WriteString("### Recently Completed\n")
		for _, task := range summary.CompletedTasks {
			builder.WriteString(fmt.Sprintf("- %s\n", task))
		}
		builder.WriteString("\n")
	}

	if len(summary.PendingTasks) > 0 {
		builder.WriteString("### Pending Tasks\n")
		for _, task := range summary.PendingTasks {
			builder.WriteString(fmt.Sprintf("- %s\n", task))
		}
		builder.WriteString("\n")
	}

	// Errors to Avoid
	if len(summary.ErrorsToAvoid) > 0 {
		builder.WriteString("### Known Issues to Avoid\n")
		for _, err := range summary.ErrorsToAvoid {
			builder.WriteString(fmt.Sprintf("- %s\n", err))
		}
	}

	return builder.String()
}

// SaveCompressionContext saves the compression context to JSON
func SaveCompressionContext(summary *CompressionSummary) (string, error) {
	data, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal compression summary: %w", err)
	}
	return string(data), nil
}

// LoadCompressionContext loads compression context from JSON
func LoadCompressionContext(data string) (*CompressionSummary, error) {
	var summary CompressionSummary
	if err := json.Unmarshal([]byte(data), &summary); err != nil {
		return nil, fmt.Errorf("failed to unmarshal compression summary: %w", err)
	}
	return &summary, nil
}