package converter

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Event represents a captured event from Claude Code hooks
type Event struct {
	SessionID   string          `json:"session_id"`
	EventType   string          `json:"event_type"`
	Timestamp   time.Time       `json:"timestamp"`
	Data        json.RawMessage `json:"data"`
	SequenceNum int             `json:"sequence_num"`
}

// SessionMetadata contains metadata about a session
type SessionMetadata struct {
	ID            string     `json:"id"`
	StartTime     time.Time  `json:"start_time"`
	EndTime       *time.Time `json:"end_time,omitempty"`
	Status        string     `json:"status"`
	WorkingDir    string     `json:"working_dir"`
	ProjectName   string     `json:"project_name,omitempty"`
	EventCount    int        `json:"event_count"`
	LastEventTime time.Time  `json:"last_event_time"`
}

// CompletedConversation represents a fully processed conversation in structured JSON format
type CompletedConversation struct {
	// Metadata section
	Metadata ConversationMetadata `json:"metadata"`

	// Conversation flow with chronological events
	Conversation []ConversationEvent `json:"conversation"`

	// Summary information
	Summary ConversationSummary `json:"summary"`

	// Export metadata
	Export ExportMetadata `json:"export"`
}

// ConversationMetadata contains high-level conversation information
type ConversationMetadata struct {
	SessionID     string    `json:"session_id"`
	ProjectName   string    `json:"project_name,omitempty"`
	WorkingDir    string    `json:"working_dir"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	Duration      string    `json:"duration"` // Human readable duration
	Status        string    `json:"status"`   // completed, timeout, error
	EventCount    int       `json:"event_count"`
	UserPrompts   int       `json:"user_prompts"`
	ClaudeReplies int       `json:"claude_replies"`
}

// ConversationEvent represents a single event in the conversation timeline
type ConversationEvent struct {
	Timestamp   time.Time       `json:"timestamp"`
	EventType   string          `json:"event_type"`
	SequenceNum int             `json:"sequence_num"`
	Content     EventContent    `json:"content"`
	Metadata    json.RawMessage `json:"metadata,omitempty"` // Original event metadata
}

// EventContent contains the processed content for each event type
type EventContent struct {
	// For session-start events
	SessionInfo *SessionStartInfo `json:"session_info,omitempty"`

	// For user-prompt events
	UserPrompt *UserPromptInfo `json:"user_prompt,omitempty"`

	// For claude-response events (future extension)
	ClaudeResponse *ClaudeResponseInfo `json:"claude_response,omitempty"`

	// For session-end events
	SessionEnd *SessionEndInfo `json:"session_end,omitempty"`
}

// SessionStartInfo contains session start event details
type SessionStartInfo struct {
	Event      string `json:"event"`
	Project    string `json:"project"`
	WorkingDir string `json:"working_dir"`
}

// UserPromptInfo contains user prompt event details
type UserPromptInfo struct {
	Message   string `json:"message"`
	Tokens    int    `json:"tokens,omitempty"`    // Future: token count
	Wordcount int    `json:"wordcount,omitempty"` // Word count for analysis
}

// ClaudeResponseInfo contains Claude response event details (future extension)
type ClaudeResponseInfo struct {
	Response  string `json:"response"`
	Tokens    int    `json:"tokens,omitempty"`
	Wordcount int    `json:"wordcount,omitempty"`
}

// SessionEndInfo contains session end event details
type SessionEndInfo struct {
	Event  string `json:"event"`
	Reason string `json:"reason,omitempty"` // normal, timeout, error
}

// ConversationSummary contains aggregated conversation analysis
type ConversationSummary struct {
	TotalDuration    string                  `json:"total_duration"`
	PromptCount      int                     `json:"prompt_count"`
	ResponseCount    int                     `json:"response_count"`
	AverageGapTime   string                  `json:"average_gap_time"` // Time between prompts
	PeakActivity     *ActivityPeak           `json:"peak_activity,omitempty"`
	TopicKeywords    []string                `json:"topic_keywords,omitempty"`    // Future: extracted keywords
	ConversationTags []string                `json:"conversation_tags,omitempty"` // Future: auto-generated tags
	Statistics       ConversationStatistics  `json:"statistics"`
}

// ActivityPeak represents the period of highest conversation activity
type ActivityPeak struct {
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	EventCount int       `json:"event_count"`
	Duration   string    `json:"duration"`
}

// ConversationStatistics contains numerical analysis of the conversation
type ConversationStatistics struct {
	TotalEvents         int     `json:"total_events"`
	UserPromptWords     int     `json:"user_prompt_words"`
	ClaudeResponseWords int     `json:"claude_response_words"`
	AveragePromptLength float64 `json:"average_prompt_length"`
	LongestGapMinutes   float64 `json:"longest_gap_minutes"`
	ShortestGapSeconds  float64 `json:"shortest_gap_seconds"`
}

// ExportMetadata contains information about the export process
type ExportMetadata struct {
	ExportTime    time.Time `json:"export_time"`
	ExportVersion string    `json:"export_version"`
	SourceFormat  string    `json:"source_format"`
	ProcessedBy   string    `json:"processed_by"`
	OriginalFile  string    `json:"original_file,omitempty"`
}

// SessionConverter handles conversion from JSONL to structured JSON
type SessionConverter struct {
	conversationsDir string
	exportVersion    string
}

// NewSessionConverter creates a new session converter
func NewSessionConverter(conversationsDir string) *SessionConverter {
	return &SessionConverter{
		conversationsDir: conversationsDir,
		exportVersion:    "1.0.0",
	}
}

// ConvertSession converts a completed session from JSONL to structured JSON
func (sc *SessionConverter) ConvertSession(sessionID string, metadata *SessionMetadata) (*CompletedConversation, error) {
	// Load events from JSONL file
	events, err := sc.loadSessionEvents(sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to load session events: %w", err)
	}

	// Sort events by sequence number to ensure proper order
	sort.Slice(events, func(i, j int) bool {
		return events[i].SequenceNum < events[j].SequenceNum
	})

	// Build completed conversation structure
	conversation := &CompletedConversation{
		Metadata:     sc.buildMetadata(metadata, events),
		Conversation: sc.buildConversationEvents(events),
		Summary:      sc.buildSummary(metadata, events),
		Export:       sc.buildExportMetadata(sessionID),
	}

	return conversation, nil
}

// SaveCompletedConversation saves the structured conversation to JSON file
func (sc *SessionConverter) SaveCompletedConversation(sessionID string, conversation *CompletedConversation) error {
	// Create completed conversation file path
	completedPath := filepath.Join(sc.conversationsDir, "completed", sessionID)
	conversationFile := filepath.Join(completedPath, "conversation.json")

	// Ensure directory exists
	if err := os.MkdirAll(completedPath, 0755); err != nil {
		return fmt.Errorf("failed to create completed directory: %w", err)
	}

	// Marshal to pretty JSON
	data, err := json.MarshalIndent(conversation, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal conversation: %w", err)
	}

	// Write to file
	if err := os.WriteFile(conversationFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write conversation file: %w", err)
	}

	return nil
}

// ArchiveOriginalJSONL moves the original JSONL file to archive location
func (sc *SessionConverter) ArchiveOriginalJSONL(sessionID string) error {
	completedPath := filepath.Join(sc.conversationsDir, "completed", sessionID)
	eventsFile := filepath.Join(completedPath, "events.jsonl")
	archiveFile := filepath.Join(completedPath, "events.jsonl.archive")

	// Check if JSONL file exists
	if _, err := os.Stat(eventsFile); os.IsNotExist(err) {
		return nil // No JSONL file to archive
	}

	// Rename to archive
	if err := os.Rename(eventsFile, archiveFile); err != nil {
		return fmt.Errorf("failed to archive JSONL file: %w", err)
	}

	return nil
}

// LoadCompletedConversation loads a completed conversation from JSON file
func (sc *SessionConverter) LoadCompletedConversation(sessionID string) (*CompletedConversation, error) {
	conversationFile := filepath.Join(sc.conversationsDir, "completed", sessionID, "conversation.json")

	data, err := os.ReadFile(conversationFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read conversation file: %w", err)
	}

	var conversation CompletedConversation
	if err := json.Unmarshal(data, &conversation); err != nil {
		return nil, fmt.Errorf("failed to unmarshal conversation: %w", err)
	}

	return &conversation, nil
}

// ListCompletedConversations returns a list of all completed conversations
func (sc *SessionConverter) ListCompletedConversations() ([]ConversationMetadata, error) {
	completedDir := filepath.Join(sc.conversationsDir, "completed")

	entries, err := os.ReadDir(completedDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []ConversationMetadata{}, nil
		}
		return nil, fmt.Errorf("failed to read completed directory: %w", err)
	}

	var conversations []ConversationMetadata

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		sessionID := entry.Name()
		conversation, err := sc.LoadCompletedConversation(sessionID)
		if err != nil {
			// Skip invalid conversations
			continue
		}

		conversations = append(conversations, conversation.Metadata)
	}

	// Sort by start time (newest first)
	sort.Slice(conversations, func(i, j int) bool {
		return conversations[i].StartTime.After(conversations[j].StartTime)
	})

	return conversations, nil
}

// loadSessionEvents loads events from the JSONL file
func (sc *SessionConverter) loadSessionEvents(sessionID string) ([]*Event, error) {
	// Check both active and completed directories for events file
	var eventsPath string

	// Try completed directory first
	completedPath := filepath.Join(sc.conversationsDir, "completed", sessionID, "events.jsonl")
	if _, err := os.Stat(completedPath); err == nil {
		eventsPath = completedPath
	} else {
		// Try active directory
		activePath := filepath.Join(sc.conversationsDir, "active", sessionID, "events.jsonl")
		if _, err := os.Stat(activePath); err == nil {
			eventsPath = activePath
		} else {
			return nil, fmt.Errorf("events file not found for session %s", sessionID)
		}
	}

	// Read JSONL file directly
	file, err := os.Open(eventsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open events file: %w", err)
	}
	defer file.Close()

	var events []*Event
	decoder := json.NewDecoder(file)

	for decoder.More() {
		var event Event
		if err := decoder.Decode(&event); err != nil {
			// Skip invalid lines but continue processing
			continue
		}
		events = append(events, &event)
	}

	return events, nil
}

// Helper methods for building conversation components

func (sc *SessionConverter) buildMetadata(metadata *SessionMetadata, events []*Event) ConversationMetadata {
	duration := ""
	endTime := time.Now() // Default to now if EndTime is nil

	if metadata.EndTime != nil {
		duration = metadata.EndTime.Sub(metadata.StartTime).String()
		endTime = *metadata.EndTime
	} else {
		duration = time.Since(metadata.StartTime).String()
	}

	userPrompts := 0
	claudeReplies := 0

	for _, event := range events {
		switch event.EventType {
		case "user-prompt":
			userPrompts++
		case "claude-response":
			claudeReplies++
		}
	}

	return ConversationMetadata{
		SessionID:     metadata.ID,
		ProjectName:   metadata.ProjectName,
		WorkingDir:    metadata.WorkingDir,
		StartTime:     metadata.StartTime,
		EndTime:       endTime,
		Duration:      duration,
		Status:        metadata.Status,
		EventCount:    metadata.EventCount,
		UserPrompts:   userPrompts,
		ClaudeReplies: claudeReplies,
	}
}

func (sc *SessionConverter) buildConversationEvents(events []*Event) []ConversationEvent {
	var conversationEvents []ConversationEvent

	for _, event := range events {
		convEvent := ConversationEvent{
			Timestamp:   event.Timestamp,
			EventType:   event.EventType,
			SequenceNum: event.SequenceNum,
			Content:     sc.buildEventContent(event),
			Metadata:    event.Data,
		}

		conversationEvents = append(conversationEvents, convEvent)
	}

	return conversationEvents
}

func (sc *SessionConverter) buildEventContent(event *Event) EventContent {
	var content EventContent

	switch event.EventType {
	case "session-start":
		content.SessionInfo = sc.parseSessionStartData(event.Data)
	case "user-prompt":
		content.UserPrompt = sc.parseUserPromptData(event.Data)
	case "claude-response":
		content.ClaudeResponse = sc.parseClaudeResponseData(event.Data)
	case "session-end":
		content.SessionEnd = sc.parseSessionEndData(event.Data)
	}

	return content
}

func (sc *SessionConverter) parseSessionStartData(data json.RawMessage) *SessionStartInfo {
	var sessionData map[string]interface{}
	if err := json.Unmarshal(data, &sessionData); err != nil {
		return &SessionStartInfo{Event: "session-start"}
	}

	info := &SessionStartInfo{}
	if event, ok := sessionData["event"].(string); ok {
		info.Event = event
	}
	if project, ok := sessionData["project"].(string); ok {
		info.Project = project
	}
	if workingDir, ok := sessionData["working_dir"].(string); ok {
		info.WorkingDir = workingDir
	}

	return info
}

func (sc *SessionConverter) parseUserPromptData(data json.RawMessage) *UserPromptInfo {
	var promptData map[string]interface{}
	if err := json.Unmarshal(data, &promptData); err != nil {
		return &UserPromptInfo{}
	}

	info := &UserPromptInfo{}
	if message, ok := promptData["message"].(string); ok {
		info.Message = message
		// Basic word count calculation
		info.Wordcount = len(strings.Fields(message))
	}

	return info
}

func (sc *SessionConverter) parseClaudeResponseData(data json.RawMessage) *ClaudeResponseInfo {
	var responseData map[string]interface{}
	if err := json.Unmarshal(data, &responseData); err != nil {
		return &ClaudeResponseInfo{}
	}

	info := &ClaudeResponseInfo{}
	if response, ok := responseData["response"].(string); ok {
		info.Response = response
		info.Wordcount = len(strings.Fields(response))
	}

	return info
}

func (sc *SessionConverter) parseSessionEndData(data json.RawMessage) *SessionEndInfo {
	var endData map[string]interface{}
	if err := json.Unmarshal(data, &endData); err != nil {
		return &SessionEndInfo{Event: "session-end"}
	}

	info := &SessionEndInfo{}
	if event, ok := endData["event"].(string); ok {
		info.Event = event
	}
	if reason, ok := endData["reason"].(string); ok {
		info.Reason = reason
	}

	return info
}

func (sc *SessionConverter) buildSummary(metadata *SessionMetadata, events []*Event) ConversationSummary {
	summary := ConversationSummary{
		Statistics: ConversationStatistics{
			TotalEvents: len(events),
		},
	}

	if metadata.EndTime != nil {
		summary.TotalDuration = metadata.EndTime.Sub(metadata.StartTime).String()
	}

	// Count prompts and responses
	var userPromptWords, claudeResponseWords int
	var promptLengths []int
	var gaps []time.Duration
	var lastEventTime *time.Time
	var allUserText []string
	var eventTimes []time.Time

	for _, event := range events {
		eventTimes = append(eventTimes, event.Timestamp)

		if lastEventTime != nil {
			gap := event.Timestamp.Sub(*lastEventTime)
			gaps = append(gaps, gap)
		}
		lastEventTime = &event.Timestamp

		switch event.EventType {
		case "user-prompt":
			summary.PromptCount++
			if prompt := sc.parseUserPromptData(event.Data); prompt != nil {
				userPromptWords += prompt.Wordcount
				promptLengths = append(promptLengths, prompt.Wordcount)
				if prompt.Message != "" {
					allUserText = append(allUserText, prompt.Message)
				}
			}
		case "claude-response":
			summary.ResponseCount++
			if response := sc.parseClaudeResponseData(event.Data); response != nil {
				claudeResponseWords += response.Wordcount
			}
		}
	}

	summary.Statistics.UserPromptWords = userPromptWords
	summary.Statistics.ClaudeResponseWords = claudeResponseWords

	// Calculate averages
	if len(promptLengths) > 0 {
		total := 0
		for _, length := range promptLengths {
			total += length
		}
		summary.Statistics.AveragePromptLength = float64(total) / float64(len(promptLengths))
	}

	if len(gaps) > 0 {
		totalGap := time.Duration(0)
		var longestGap, shortestGap time.Duration

		for i, gap := range gaps {
			totalGap += gap
			if i == 0 || gap > longestGap {
				longestGap = gap
			}
			if i == 0 || gap < shortestGap {
				shortestGap = gap
			}
		}

		avgGap := totalGap / time.Duration(len(gaps))
		summary.AverageGapTime = avgGap.String()
		summary.Statistics.LongestGapMinutes = longestGap.Minutes()
		summary.Statistics.ShortestGapSeconds = shortestGap.Seconds()
	}

	// Enhanced summarization features
	summary.TopicKeywords = sc.extractTopicKeywords(allUserText)
	summary.PeakActivity = sc.detectActivityPeak(eventTimes)
	summary.ConversationTags = sc.generateConversationTags(metadata, events)

	return summary
}

func (sc *SessionConverter) buildExportMetadata(sessionID string) ExportMetadata {
	return ExportMetadata{
		ExportTime:    time.Now(),
		ExportVersion: sc.exportVersion,
		SourceFormat:  "jsonl",
		ProcessedBy:   "context-extender",
		OriginalFile:  "events.jsonl",
	}
}

// extractTopicKeywords extracts relevant keywords from user text using basic frequency analysis
func (sc *SessionConverter) extractTopicKeywords(userTexts []string) []string {
	if len(userTexts) == 0 {
		return []string{}
	}

	// Common stop words to filter out
	stopWords := map[string]bool{
		"the": true, "a": true, "an": true, "and": true, "or": true, "but": true, "in": true, "on": true, "at": true, "to": true,
		"for": true, "of": true, "with": true, "by": true, "is": true, "are": true, "was": true, "were": true, "be": true,
		"been": true, "have": true, "has": true, "had": true, "do": true, "does": true, "did": true, "will": true, "would": true,
		"could": true, "should": true, "may": true, "might": true, "can": true, "this": true, "that": true, "these": true,
		"those": true, "i": true, "you": true, "he": true, "she": true, "it": true, "we": true, "they": true, "me": true,
		"him": true, "her": true, "us": true, "them": true, "my": true, "your": true, "his": true, "its": true, "our": true,
		"their": true, "what": true, "where": true, "when": true, "why": true, "how": true, "if": true, "then": true,
		"please": true, "help": true, "thanks": true, "thank": true,
	}

	// Count word frequencies
	wordFreq := make(map[string]int)
	for _, text := range userTexts {
		words := strings.Fields(strings.ToLower(text))
		for _, word := range words {
			// Clean word (remove punctuation)
			cleaned := strings.Trim(word, ".,!?;:()[]{}\"'")
			if len(cleaned) > 2 && !stopWords[cleaned] {
				wordFreq[cleaned]++
			}
		}
	}

	// Find top keywords
	type wordCount struct {
		word  string
		count int
	}

	var wordCounts []wordCount
	for word, count := range wordFreq {
		if count > 1 || len(userTexts) == 1 { // Include words that appear multiple times, or all words if only one text
			wordCounts = append(wordCounts, wordCount{word, count})
		}
	}

	// Sort by frequency
	sort.Slice(wordCounts, func(i, j int) bool {
		return wordCounts[i].count > wordCounts[j].count
	})

	// Return top 5 keywords
	var keywords []string
	for i, wc := range wordCounts {
		if i >= 5 {
			break
		}
		keywords = append(keywords, wc.word)
	}

	return keywords
}

// detectActivityPeak finds periods of high conversation activity
func (sc *SessionConverter) detectActivityPeak(eventTimes []time.Time) *ActivityPeak {
	if len(eventTimes) < 3 {
		return nil // Need at least 3 events to detect a peak
	}

	// Sort times to ensure chronological order
	sort.Slice(eventTimes, func(i, j int) bool {
		return eventTimes[i].Before(eventTimes[j])
	})

	// Use a sliding window approach to find the period with most events
	windowSize := 5 * time.Minute // 5-minute window
	maxEvents := 0
	var peakStart, peakEnd time.Time

	for i := 0; i < len(eventTimes); i++ {
		windowStart := eventTimes[i]
		windowEnd := windowStart.Add(windowSize)
		eventsInWindow := 0

		for j := i; j < len(eventTimes) && eventTimes[j].Before(windowEnd); j++ {
			eventsInWindow++
		}

		if eventsInWindow > maxEvents {
			maxEvents = eventsInWindow
			peakStart = windowStart
			peakEnd = windowEnd
		}
	}

	// Only return peak if it contains more than half the events
	if maxEvents < len(eventTimes)/2 {
		return nil
	}

	return &ActivityPeak{
		StartTime:  peakStart,
		EndTime:    peakEnd,
		EventCount: maxEvents,
		Duration:   peakEnd.Sub(peakStart).String(),
	}
}

// generateConversationTags creates automatic tags based on conversation characteristics
func (sc *SessionConverter) generateConversationTags(metadata *SessionMetadata, events []*Event) []string {
	var tags []string

	// Duration-based tags
	if metadata.EndTime != nil {
		duration := metadata.EndTime.Sub(metadata.StartTime)
		if duration < 2*time.Minute {
			tags = append(tags, "quick-session")
		} else if duration > 30*time.Minute {
			tags = append(tags, "long-session")
		}
	}

	// Event count based tags
	if len(events) >= 10 {
		tags = append(tags, "active-conversation")
	} else if len(events) <= 3 {
		tags = append(tags, "brief-interaction")
	}

	// Status-based tags
	switch metadata.Status {
	case "completed":
		tags = append(tags, "completed-naturally")
	case "timeout":
		tags = append(tags, "auto-timeout")
	case "error":
		tags = append(tags, "error-terminated")
	}

	// Time-based tags
	hour := metadata.StartTime.Hour()
	if hour >= 6 && hour < 12 {
		tags = append(tags, "morning")
	} else if hour >= 12 && hour < 18 {
		tags = append(tags, "afternoon")
	} else if hour >= 18 && hour < 22 {
		tags = append(tags, "evening")
	} else {
		tags = append(tags, "night")
	}

	// Project-based tags
	if metadata.ProjectName != "" {
		projectLower := strings.ToLower(metadata.ProjectName)
		if strings.Contains(projectLower, "test") {
			tags = append(tags, "testing")
		}
		if strings.Contains(projectLower, "web") {
			tags = append(tags, "web-dev")
		}
		if strings.Contains(projectLower, "api") {
			tags = append(tags, "api-dev")
		}
		if strings.Contains(projectLower, "extender") {
			tags = append(tags, "context-extender")
		}
	}

	// Interaction pattern tags
	promptCount := 0
	responseCount := 0
	for _, event := range events {
		switch event.EventType {
		case "user-prompt":
			promptCount++
		case "claude-response":
			responseCount++
		}
	}

	if promptCount > 0 && responseCount == 0 {
		tags = append(tags, "prompts-only")
	} else if promptCount > 0 && responseCount > 0 {
		ratio := float64(responseCount) / float64(promptCount)
		if ratio > 0.8 {
			tags = append(tags, "interactive")
		}
	}

	return tags
}