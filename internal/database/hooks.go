package database

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type HookHandler struct {
	sessionSequence map[string]int
	mu              sync.RWMutex
}

type SessionStartData struct {
	SessionID string            `json:"session_id"`
	Timestamp time.Time         `json:"timestamp"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

type UserPromptData struct {
	SessionID string    `json:"session_id"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

type ClaudeResponseData struct {
	SessionID  string    `json:"session_id"`
	Response   string    `json:"response"`
	Timestamp  time.Time `json:"timestamp"`
	TokenCount *int      `json:"token_count,omitempty"`
	ModelInfo  *string   `json:"model_info,omitempty"`
}

type SessionEndData struct {
	SessionID string    `json:"session_id"`
	Timestamp time.Time `json:"timestamp"`
	Summary   string    `json:"summary,omitempty"`
}

var defaultHandler *HookHandler
var handlerOnce sync.Once

func GetHookHandler() *HookHandler {
	handlerOnce.Do(func() {
		defaultHandler = &HookHandler{
			sessionSequence: make(map[string]int),
		}
	})
	return defaultHandler
}

func (h *HookHandler) HandleSessionStart(data SessionStartData) error {
	start := time.Now()
	defer func() {
		RecordHookExecution("session_start", time.Since(start))
	}()

	h.mu.Lock()
	h.sessionSequence[data.SessionID] = 0
	h.mu.Unlock()

	metadataJSON, err := json.Marshal(data.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	session := &Session{
		ID:        data.SessionID,
		CreatedAt: data.Timestamp,
		UpdatedAt: data.Timestamp,
		Status:    "active",
		Metadata:  string(metadataJSON),
	}

	if err := CreateSession(session); err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	IncrementWriteCount()

	eventData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal session start data: %w", err)
	}

	event := &Event{
		ID:          uuid.New().String(),
		SessionID:   data.SessionID,
		EventType:   "session_start",
		Data:        string(eventData),
		Timestamp:   data.Timestamp,
		SequenceNum: h.getNextSequence(data.SessionID),
	}

	err = CreateEvent(event)
	if err == nil {
		IncrementWriteCount()
	}
	return err
}

func (h *HookHandler) HandleUserPrompt(data UserPromptData) error {
	start := time.Now()
	defer func() {
		RecordHookExecution("user_prompt", time.Since(start))
	}()

	eventData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal user prompt data: %w", err)
	}

	event := &Event{
		ID:          uuid.New().String(),
		SessionID:   data.SessionID,
		EventType:   "user_prompt",
		Data:        string(eventData),
		Timestamp:   data.Timestamp,
		SequenceNum: h.getNextSequence(data.SessionID),
	}

	if err := CreateEvent(event); err != nil {
		return fmt.Errorf("failed to create user prompt event: %w", err)
	}
	IncrementWriteCount()

	conversation := &Conversation{
		ID:          uuid.New().String(),
		SessionID:   data.SessionID,
		MessageType: "user",
		Content:     data.Message,
		Timestamp:   data.Timestamp,
	}

	if err := CreateConversation(conversation); err != nil {
		return fmt.Errorf("failed to create user conversation: %w", err)
	}
	IncrementWriteCount()

	return nil
}

func (h *HookHandler) HandleClaudeResponse(data ClaudeResponseData) error {
	eventData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal claude response data: %w", err)
	}

	event := &Event{
		ID:          uuid.New().String(),
		SessionID:   data.SessionID,
		EventType:   "claude_response",
		Data:        string(eventData),
		Timestamp:   data.Timestamp,
		SequenceNum: h.getNextSequence(data.SessionID),
	}

	if err := CreateEvent(event); err != nil {
		return fmt.Errorf("failed to create claude response event: %w", err)
	}

	tokenCount := 0
	if data.TokenCount != nil {
		tokenCount = *data.TokenCount
	}
	model := ""
	if data.ModelInfo != nil {
		model = *data.ModelInfo
	}

	conversation := &Conversation{
		ID:          uuid.New().String(),
		SessionID:   data.SessionID,
		MessageType: "assistant",
		Content:     data.Response,
		Timestamp:   data.Timestamp,
		TokenCount:  tokenCount,
		Model:       model,
	}

	if err := CreateConversation(conversation); err != nil {
		return fmt.Errorf("failed to create claude conversation: %w", err)
	}

	return nil
}

func (h *HookHandler) HandleSessionEnd(data SessionEndData) error {
	metadata := make(map[string]string)
	if data.Summary != "" {
		metadata["summary"] = data.Summary
	}

	if err := UpdateSession(data.SessionID, "completed", metadata); err != nil {
		return fmt.Errorf("failed to update session: %w", err)
	}

	eventData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal session end data: %w", err)
	}

	event := &Event{
		ID:          uuid.New().String(),
		SessionID:   data.SessionID,
		EventType:   "session_end",
		Data:        string(eventData),
		Timestamp:   data.Timestamp,
		SequenceNum: h.getNextSequence(data.SessionID),
	}

	if err := CreateEvent(event); err != nil {
		return fmt.Errorf("failed to create session end event: %w", err)
	}

	h.mu.Lock()
	delete(h.sessionSequence, data.SessionID)
	h.mu.Unlock()

	return nil
}

func (h *HookHandler) getNextSequence(sessionID string) int {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.sessionSequence[sessionID]++
	return h.sessionSequence[sessionID]
}

func HandleSessionStartHook(sessionID string) error {
	if sessionID == "" {
		sessionID = uuid.New().String()
	}

	data := SessionStartData{
		SessionID: sessionID,
		Timestamp: time.Now(),
		Metadata:  make(map[string]string),
	}

	return GetHookHandler().HandleSessionStart(data)
}

func HandleUserPromptHook(sessionID, message string) error {
	data := UserPromptData{
		SessionID: sessionID,
		Message:   message,
		Timestamp: time.Now(),
	}

	return GetHookHandler().HandleUserPrompt(data)
}

func HandleClaudeResponseHook(sessionID, response string, tokenCount *int, modelInfo *string) error {
	data := ClaudeResponseData{
		SessionID:  sessionID,
		Response:   response,
		Timestamp:  time.Now(),
		TokenCount: tokenCount,
		ModelInfo:  modelInfo,
	}

	return GetHookHandler().HandleClaudeResponse(data)
}

func HandleSessionEndHook(sessionID, summary string) error {
	data := SessionEndData{
		SessionID: sessionID,
		Timestamp: time.Now(),
		Summary:   summary,
	}

	return GetHookHandler().HandleSessionEnd(data)
}