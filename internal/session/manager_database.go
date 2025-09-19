package session

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"context-extender/internal/database"
)

// SessionStatus represents the status of a session
type SessionStatus string

const (
	StatusActive    SessionStatus = "active"
	StatusCompleted SessionStatus = "completed"
	StatusTimeout   SessionStatus = "timeout"
	StatusError     SessionStatus = "error"
)

// EventType represents different hook event types
type EventType string

const (
	EventSessionStart   EventType = "session-start"
	EventUserPrompt     EventType = "user-prompt"
	EventClaudeResponse EventType = "claude-response"
	EventSessionEnd     EventType = "session-end"
)

// SessionMetadata contains metadata about a session
type SessionMetadata struct {
	ID            string        `json:"id"`
	StartTime     time.Time     `json:"start_time"`
	EndTime       *time.Time    `json:"end_time,omitempty"`
	Status        SessionStatus `json:"status"`
	WorkingDir    string        `json:"working_dir"`
	ProjectName   string        `json:"project_name,omitempty"`
	EventCount    int           `json:"event_count"`
	LastEventTime time.Time     `json:"last_event_time"`
}

// Event represents a captured event from Claude Code hooks
type Event struct {
	SessionID   string          `json:"session_id"`
	EventType   EventType       `json:"event_type"`
	Timestamp   time.Time       `json:"timestamp"`
	Data        json.RawMessage `json:"data"`
	SequenceNum int             `json:"sequence_num"`
}

// DatabaseSessionManager manages session lifecycle using SQLite database
type DatabaseSessionManager struct {
	activeSessions map[string]*SessionMetadata
	sessionLock    sync.RWMutex
	eventSequence  map[string]int
	sessionTimeout time.Duration
	hookHandler    *database.HookHandler
}

// NewDatabaseSessionManager creates a new database-based session manager
func NewDatabaseSessionManager() (*DatabaseSessionManager, error) {
	// Initialize database
	config := database.DefaultConfig()
	if err := database.Initialize(config); err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Run migrations
	if err := database.RunMigrations(); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return &DatabaseSessionManager{
		activeSessions: make(map[string]*SessionMetadata),
		eventSequence:  make(map[string]int),
		sessionTimeout: 30 * time.Minute, // Default 30 minute timeout
		hookHandler:    database.GetHookHandler(),
	}, nil
}

// StartSession creates and starts a new session
func (dsm *DatabaseSessionManager) StartSession(workingDir, projectName string) (*SessionMetadata, error) {
	dsm.sessionLock.Lock()
	defer dsm.sessionLock.Unlock()

	// Generate new session ID
	sessionID := uuid.New().String()

	// Create session metadata
	metadata := &SessionMetadata{
		ID:            sessionID,
		StartTime:     time.Now(),
		Status:        StatusActive,
		WorkingDir:    workingDir,
		ProjectName:   projectName,
		EventCount:    0,
		LastEventTime: time.Now(),
	}

	// Store in active sessions
	dsm.activeSessions[sessionID] = metadata
	dsm.eventSequence[sessionID] = 0

	// Handle session start via database hook
	sessionStartData := database.SessionStartData{
		SessionID: sessionID,
		Timestamp: metadata.StartTime,
		Metadata: map[string]string{
			"working_dir":  workingDir,
			"project_name": projectName,
		},
	}

	if err := dsm.hookHandler.HandleSessionStart(sessionStartData); err != nil {
		return nil, fmt.Errorf("failed to handle session start: %w", err)
	}

	return metadata, nil
}

// EndSession marks a session as completed
func (dsm *DatabaseSessionManager) EndSession(sessionID string) error {
	dsm.sessionLock.Lock()
	defer dsm.sessionLock.Unlock()

	metadata, exists := dsm.activeSessions[sessionID]
	if !exists {
		// Try to load from database
		var err error
		metadata, err = dsm.loadSessionFromDatabase(sessionID)
		if err != nil {
			return fmt.Errorf("session not found: %s", sessionID)
		}
	}

	// Update metadata
	endTime := time.Now()
	metadata.EndTime = &endTime
	metadata.Status = StatusCompleted

	// Handle session end via database hook
	sessionEndData := database.SessionEndData{
		SessionID: sessionID,
		Timestamp: endTime,
		Summary:   fmt.Sprintf("Session completed with %d events", metadata.EventCount),
	}

	if err := dsm.hookHandler.HandleSessionEnd(sessionEndData); err != nil {
		return fmt.Errorf("failed to handle session end: %w", err)
	}

	// Remove from active sessions
	delete(dsm.activeSessions, sessionID)
	delete(dsm.eventSequence, sessionID)

	return nil
}

// RecordEvent records an event for a session
func (dsm *DatabaseSessionManager) RecordEvent(sessionID string, eventType EventType, data json.RawMessage) (*Event, error) {
	dsm.sessionLock.Lock()
	defer dsm.sessionLock.Unlock()

	// Check if session exists
	metadata, exists := dsm.activeSessions[sessionID]
	if !exists {
		// Try to load from database
		var err error
		metadata, err = dsm.loadSessionFromDatabase(sessionID)
		if err != nil {
			return nil, fmt.Errorf("session not found: %s", sessionID)
		}
		dsm.activeSessions[sessionID] = metadata
	}

	// Check if session is active
	if metadata.Status != StatusActive {
		return nil, fmt.Errorf("session %s is not active (status: %s)", sessionID, metadata.Status)
	}

	// Get next sequence number
	dsm.eventSequence[sessionID]++
	sequenceNum := dsm.eventSequence[sessionID]

	// Create event
	event := &Event{
		SessionID:   sessionID,
		EventType:   eventType,
		Timestamp:   time.Now(),
		Data:        data,
		SequenceNum: sequenceNum,
	}

	// Update session metadata
	metadata.EventCount++
	metadata.LastEventTime = time.Now()

	// Handle specific event types via database hooks
	switch eventType {
	case EventUserPrompt:
		var promptData struct {
			Message string `json:"message"`
		}
		if err := json.Unmarshal(data, &promptData); err == nil {
			userPromptData := database.UserPromptData{
				SessionID: sessionID,
				Message:   promptData.Message,
				Timestamp: event.Timestamp,
			}
			if err := dsm.hookHandler.HandleUserPrompt(userPromptData); err != nil {
				return nil, fmt.Errorf("failed to handle user prompt: %w", err)
			}
		}

	case EventClaudeResponse:
		var responseData struct {
			Response   string `json:"response"`
			TokenCount *int   `json:"token_count,omitempty"`
			ModelInfo  *string `json:"model_info,omitempty"`
		}
		if err := json.Unmarshal(data, &responseData); err == nil {
			claudeResponseData := database.ClaudeResponseData{
				SessionID:  sessionID,
				Response:   responseData.Response,
				Timestamp:  event.Timestamp,
				TokenCount: responseData.TokenCount,
				ModelInfo:  responseData.ModelInfo,
			}
			if err := dsm.hookHandler.HandleClaudeResponse(claudeResponseData); err != nil {
				return nil, fmt.Errorf("failed to handle claude response: %w", err)
			}
		}
	}

	return event, nil
}

// GetActiveSession retrieves an active session by ID
func (dsm *DatabaseSessionManager) GetActiveSession(sessionID string) (*SessionMetadata, error) {
	dsm.sessionLock.RLock()
	defer dsm.sessionLock.RUnlock()

	metadata, exists := dsm.activeSessions[sessionID]
	if !exists {
		// Try to load from database
		var err error
		metadata, err = dsm.loadSessionFromDatabase(sessionID)
		if err != nil {
			return nil, fmt.Errorf("session not found: %s", sessionID)
		}
	}

	return metadata, nil
}

// GetAllActiveSessions returns all active sessions
func (dsm *DatabaseSessionManager) GetAllActiveSessions() ([]*SessionMetadata, error) {
	dsm.sessionLock.RLock()
	defer dsm.sessionLock.RUnlock()

	// Load all active sessions from database
	if err := dsm.loadAllActiveSessionsFromDatabase(); err != nil {
		return nil, fmt.Errorf("failed to load active sessions: %w", err)
	}

	sessions := make([]*SessionMetadata, 0, len(dsm.activeSessions))
	for _, metadata := range dsm.activeSessions {
		sessions = append(sessions, metadata)
	}

	return sessions, nil
}

// FindSessionByWorkingDir finds an active session by working directory
func (dsm *DatabaseSessionManager) FindSessionByWorkingDir(workingDir string) (*SessionMetadata, error) {
	dsm.sessionLock.RLock()
	defer dsm.sessionLock.RUnlock()

	// Load all sessions from database to ensure we have the latest
	if err := dsm.loadAllActiveSessionsFromDatabase(); err != nil {
		return nil, fmt.Errorf("failed to load active sessions: %w", err)
	}

	for _, metadata := range dsm.activeSessions {
		if metadata.WorkingDir == workingDir && metadata.Status == StatusActive {
			// Check for timeout
			if time.Since(metadata.LastEventTime) > dsm.sessionTimeout {
				continue
			}
			return metadata, nil
		}
	}

	return nil, fmt.Errorf("no active session found for working directory: %s", workingDir)
}

// CheckAndTimeoutSessions checks for and handles timed out sessions
func (dsm *DatabaseSessionManager) CheckAndTimeoutSessions() error {
	dsm.sessionLock.Lock()
	defer dsm.sessionLock.Unlock()

	// Load all sessions from database
	if err := dsm.loadAllActiveSessionsFromDatabase(); err != nil {
		return fmt.Errorf("failed to load active sessions: %w", err)
	}

	timedOutSessions := []string{}

	for sessionID, metadata := range dsm.activeSessions {
		if metadata.Status == StatusActive && time.Since(metadata.LastEventTime) > dsm.sessionTimeout {
			timedOutSessions = append(timedOutSessions, sessionID)
		}
	}

	// Process timed out sessions
	for _, sessionID := range timedOutSessions {
		metadata := dsm.activeSessions[sessionID]
		endTime := metadata.LastEventTime.Add(dsm.sessionTimeout)

		// Handle session end via database hook with timeout status
		sessionEndData := database.SessionEndData{
			SessionID: sessionID,
			Timestamp: endTime,
			Summary:   fmt.Sprintf("Session timed out after %v", dsm.sessionTimeout),
		}

		if err := dsm.hookHandler.HandleSessionEnd(sessionEndData); err != nil {
			return fmt.Errorf("failed to handle timeout for session %s: %w", sessionID, err)
		}

		// Remove from active sessions
		delete(dsm.activeSessions, sessionID)
		delete(dsm.eventSequence, sessionID)
	}

	return nil
}

// SetSessionTimeout sets the session timeout duration
func (dsm *DatabaseSessionManager) SetSessionTimeout(timeout time.Duration) {
	dsm.sessionTimeout = timeout
}

// loadSessionFromDatabase loads session metadata from database
func (dsm *DatabaseSessionManager) loadSessionFromDatabase(sessionID string) (*SessionMetadata, error) {
	session, err := database.GetSession(sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session from database: %w", err)
	}

	if session == nil {
		return nil, fmt.Errorf("session not found: %s", sessionID)
	}

	// Convert database session to session metadata
	metadata := &SessionMetadata{
		ID:        session.ID,
		StartTime: session.CreatedAt,
		Status:    SessionStatus(session.Status),
	}

	if session.Status == "completed" {
		metadata.EndTime = &session.UpdatedAt
	}

	// Extract metadata from session metadata JSON
	if session.Metadata != "" {
		var metadataMap map[string]string
		if err := json.Unmarshal([]byte(session.Metadata), &metadataMap); err == nil {
			if workingDir, ok := metadataMap["working_dir"]; ok {
				metadata.WorkingDir = workingDir
			}
			if projectName, ok := metadataMap["project_name"]; ok {
				metadata.ProjectName = projectName
			}
		}
	}

	// Get event count from database
	events, err := database.GetEventsBySession(sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get events for session: %w", err)
	}

	metadata.EventCount = len(events)
	if len(events) > 0 {
		metadata.LastEventTime = events[len(events)-1].Timestamp
	} else {
		metadata.LastEventTime = metadata.StartTime
	}

	return metadata, nil
}

// loadAllActiveSessionsFromDatabase loads all active sessions from database
func (dsm *DatabaseSessionManager) loadAllActiveSessionsFromDatabase() error {
	// For now, we'll maintain active sessions in memory
	// In a future enhancement, we could query the database for active sessions
	// This method is kept for compatibility with the existing interface
	return nil
}

// GetSessionEvents retrieves all events for a session
func (dsm *DatabaseSessionManager) GetSessionEvents(sessionID string) ([]*Event, error) {
	events, err := database.GetEventsBySession(sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get events from database: %w", err)
	}

	// Convert database events to session events
	var sessionEvents []*Event
	for _, dbEvent := range events {
		sessionEvent := &Event{
			SessionID:   dbEvent.SessionID,
			EventType:   EventType(dbEvent.EventType),
			Timestamp:   dbEvent.Timestamp,
			Data:        json.RawMessage(dbEvent.Data),
			SequenceNum: dbEvent.SequenceNum,
		}
		sessionEvents = append(sessionEvents, sessionEvent)
	}

	return sessionEvents, nil
}