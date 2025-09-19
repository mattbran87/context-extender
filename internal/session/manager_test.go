package session

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"context-extender/internal/storage"
)

// TestNewSessionManager tests session manager creation
func TestNewSessionManager(t *testing.T) {
	// Test with nil storage manager
	sm, err := NewSessionManager(nil)
	if err == nil {
		t.Error("NewSessionManager should fail with nil storage manager")
	}
	if sm != nil {
		t.Error("NewSessionManager should return nil with error")
	}

	// Test with valid storage manager
	tempDir := t.TempDir()
	storageConfig := &storage.StorageConfig{CustomPath: tempDir}
	storageMgr, err := storage.NewStorageManager(storageConfig)
	if err != nil {
		t.Fatalf("Failed to create storage manager: %v", err)
	}

	sm, err = NewSessionManager(storageMgr)
	if err != nil {
		t.Fatalf("NewSessionManager failed: %v", err)
	}

	if sm == nil {
		t.Error("NewSessionManager should return non-nil manager")
	}

	if sm.sessionTimeout != 30*time.Minute {
		t.Errorf("Default timeout should be 30 minutes, got %v", sm.sessionTimeout)
	}
}

// TestStartSession tests session creation and start
func TestStartSession(t *testing.T) {
	tempDir := t.TempDir()
	storageConfig := &storage.StorageConfig{CustomPath: tempDir}
	storageMgr, err := storage.NewStorageManager(storageConfig)
	if err != nil {
		t.Fatalf("Failed to create storage manager: %v", err)
	}

	sm, err := NewSessionManager(storageMgr)
	if err != nil {
		t.Fatalf("NewSessionManager failed: %v", err)
	}

	// Start a new session
	workingDir := "/test/project"
	projectName := "test-project"

	metadata, err := sm.StartSession(workingDir, projectName)
	if err != nil {
		t.Fatalf("StartSession failed: %v", err)
	}

	// Verify metadata
	if metadata == nil {
		t.Fatal("StartSession should return metadata")
	}

	if metadata.ID == "" {
		t.Error("Session ID should not be empty")
	}

	if metadata.WorkingDir != workingDir {
		t.Errorf("WorkingDir = %v, expected %v", metadata.WorkingDir, workingDir)
	}

	if metadata.ProjectName != projectName {
		t.Errorf("ProjectName = %v, expected %v", metadata.ProjectName, projectName)
	}

	if metadata.Status != StatusActive {
		t.Errorf("Status = %v, expected %v", metadata.Status, StatusActive)
	}

	if metadata.EventCount != 0 {
		t.Errorf("EventCount = %v, expected 0", metadata.EventCount)
	}

	// Verify persistence
	metadataPath := filepath.Join(tempDir, "conversations", "active", metadata.ID, "metadata.json")
	if _, err := os.Stat(metadataPath); os.IsNotExist(err) {
		t.Error("Session metadata file should be created")
	}
}

// TestEndSession tests session completion
func TestEndSession(t *testing.T) {
	tempDir := t.TempDir()
	storageConfig := &storage.StorageConfig{CustomPath: tempDir}
	storageMgr, err := storage.NewStorageManager(storageConfig)
	if err != nil {
		t.Fatalf("Failed to create storage manager: %v", err)
	}

	sm, err := NewSessionManager(storageMgr)
	if err != nil {
		t.Fatalf("NewSessionManager failed: %v", err)
	}

	// Start a session
	metadata, err := sm.StartSession("/test", "test")
	if err != nil {
		t.Fatalf("StartSession failed: %v", err)
	}

	sessionID := metadata.ID

	// End the session
	err = sm.EndSession(sessionID)
	if err != nil {
		t.Fatalf("EndSession failed: %v", err)
	}

	// Verify session is moved to completed
	activePath := filepath.Join(tempDir, "conversations", "active", sessionID)
	if _, err := os.Stat(activePath); !os.IsNotExist(err) {
		t.Error("Active session directory should be removed")
	}

	completedPath := filepath.Join(tempDir, "conversations", "completed", sessionID, "metadata.json")
	if _, err := os.Stat(completedPath); os.IsNotExist(err) {
		t.Error("Completed session metadata should exist")
	}

	// Verify metadata is updated
	data, err := os.ReadFile(completedPath)
	if err != nil {
		t.Fatalf("Failed to read completed metadata: %v", err)
	}

	var completedMetadata SessionMetadata
	if err := json.Unmarshal(data, &completedMetadata); err != nil {
		t.Fatalf("Failed to unmarshal metadata: %v", err)
	}

	if completedMetadata.Status != StatusCompleted {
		t.Errorf("Status = %v, expected %v", completedMetadata.Status, StatusCompleted)
	}

	if completedMetadata.EndTime == nil {
		t.Error("EndTime should be set")
	}
}

// TestRecordEvent tests event recording
func TestRecordEvent(t *testing.T) {
	tempDir := t.TempDir()
	storageConfig := &storage.StorageConfig{CustomPath: tempDir}
	storageMgr, err := storage.NewStorageManager(storageConfig)
	if err != nil {
		t.Fatalf("Failed to create storage manager: %v", err)
	}

	sm, err := NewSessionManager(storageMgr)
	if err != nil {
		t.Fatalf("NewSessionManager failed: %v", err)
	}

	// Start a session
	metadata, err := sm.StartSession("/test", "test")
	if err != nil {
		t.Fatalf("StartSession failed: %v", err)
	}

	sessionID := metadata.ID

	// Record events
	eventData := json.RawMessage(`{"message": "test event"}`)

	event1, err := sm.RecordEvent(sessionID, EventUserPrompt, eventData)
	if err != nil {
		t.Fatalf("RecordEvent failed: %v", err)
	}

	if event1.SessionID != sessionID {
		t.Errorf("Event SessionID = %v, expected %v", event1.SessionID, sessionID)
	}

	if event1.EventType != EventUserPrompt {
		t.Errorf("EventType = %v, expected %v", event1.EventType, EventUserPrompt)
	}

	if event1.SequenceNum != 1 {
		t.Errorf("SequenceNum = %v, expected 1", event1.SequenceNum)
	}

	// Record another event
	event2, err := sm.RecordEvent(sessionID, EventClaudeResponse, eventData)
	if err != nil {
		t.Fatalf("Second RecordEvent failed: %v", err)
	}

	if event2.SequenceNum != 2 {
		t.Errorf("Second event SequenceNum = %v, expected 2", event2.SequenceNum)
	}

	// Verify events are persisted
	eventsPath := filepath.Join(tempDir, "conversations", "active", sessionID, "events.jsonl")
	if _, err := os.Stat(eventsPath); os.IsNotExist(err) {
		t.Error("Events file should be created")
	}

	// Verify metadata is updated
	updatedMetadata, err := sm.GetActiveSession(sessionID)
	if err != nil {
		t.Fatalf("GetActiveSession failed: %v", err)
	}

	if updatedMetadata.EventCount != 2 {
		t.Errorf("EventCount = %v, expected 2", updatedMetadata.EventCount)
	}
}

// TestGetActiveSession tests retrieving active sessions
func TestGetActiveSession(t *testing.T) {
	tempDir := t.TempDir()
	storageConfig := &storage.StorageConfig{CustomPath: tempDir}
	storageMgr, err := storage.NewStorageManager(storageConfig)
	if err != nil {
		t.Fatalf("Failed to create storage manager: %v", err)
	}

	sm, err := NewSessionManager(storageMgr)
	if err != nil {
		t.Fatalf("NewSessionManager failed: %v", err)
	}

	// Test non-existent session
	_, err = sm.GetActiveSession("non-existent")
	if err == nil {
		t.Error("GetActiveSession should fail for non-existent session")
	}

	// Start a session
	metadata, err := sm.StartSession("/test", "test")
	if err != nil {
		t.Fatalf("StartSession failed: %v", err)
	}

	// Get active session
	retrieved, err := sm.GetActiveSession(metadata.ID)
	if err != nil {
		t.Fatalf("GetActiveSession failed: %v", err)
	}

	if retrieved.ID != metadata.ID {
		t.Errorf("Retrieved ID = %v, expected %v", retrieved.ID, metadata.ID)
	}

	if retrieved.WorkingDir != metadata.WorkingDir {
		t.Errorf("Retrieved WorkingDir = %v, expected %v", retrieved.WorkingDir, metadata.WorkingDir)
	}
}

// TestFindSessionByWorkingDir tests finding sessions by working directory
func TestFindSessionByWorkingDir(t *testing.T) {
	tempDir := t.TempDir()
	storageConfig := &storage.StorageConfig{CustomPath: tempDir}
	storageMgr, err := storage.NewStorageManager(storageConfig)
	if err != nil {
		t.Fatalf("Failed to create storage manager: %v", err)
	}

	sm, err := NewSessionManager(storageMgr)
	if err != nil {
		t.Fatalf("NewSessionManager failed: %v", err)
	}

	// Test non-existent working directory
	_, err = sm.FindSessionByWorkingDir("/non-existent")
	if err == nil {
		t.Error("FindSessionByWorkingDir should fail for non-existent directory")
	}

	// Start sessions with different working directories
	workingDir1 := "/test/project1"
	workingDir2 := "/test/project2"

	metadata1, err := sm.StartSession(workingDir1, "project1")
	if err != nil {
		t.Fatalf("StartSession failed: %v", err)
	}

	metadata2, err := sm.StartSession(workingDir2, "project2")
	if err != nil {
		t.Fatalf("StartSession failed: %v", err)
	}

	// Find by working directory
	found1, err := sm.FindSessionByWorkingDir(workingDir1)
	if err != nil {
		t.Fatalf("FindSessionByWorkingDir failed: %v", err)
	}

	if found1.ID != metadata1.ID {
		t.Errorf("Found ID = %v, expected %v", found1.ID, metadata1.ID)
	}

	found2, err := sm.FindSessionByWorkingDir(workingDir2)
	if err != nil {
		t.Fatalf("FindSessionByWorkingDir failed: %v", err)
	}

	if found2.ID != metadata2.ID {
		t.Errorf("Found ID = %v, expected %v", found2.ID, metadata2.ID)
	}
}

// TestGetAllActiveSessions tests retrieving all active sessions
func TestGetAllActiveSessions(t *testing.T) {
	tempDir := t.TempDir()
	storageConfig := &storage.StorageConfig{CustomPath: tempDir}
	storageMgr, err := storage.NewStorageManager(storageConfig)
	if err != nil {
		t.Fatalf("Failed to create storage manager: %v", err)
	}

	sm, err := NewSessionManager(storageMgr)
	if err != nil {
		t.Fatalf("NewSessionManager failed: %v", err)
	}

	// Initially should have no sessions
	sessions, err := sm.GetAllActiveSessions()
	if err != nil {
		t.Fatalf("GetAllActiveSessions failed: %v", err)
	}

	if len(sessions) != 0 {
		t.Errorf("Initial sessions count = %v, expected 0", len(sessions))
	}

	// Start multiple sessions
	for i := 0; i < 3; i++ {
		_, err := sm.StartSession(
			filepath.Join("/test", string(rune('a'+i))),
			string(rune('a'+i)),
		)
		if err != nil {
			t.Fatalf("StartSession %d failed: %v", i, err)
		}
	}

	// Get all active sessions
	sessions, err = sm.GetAllActiveSessions()
	if err != nil {
		t.Fatalf("GetAllActiveSessions failed: %v", err)
	}

	if len(sessions) != 3 {
		t.Errorf("Sessions count = %v, expected 3", len(sessions))
	}

	// End one session
	if len(sessions) > 0 {
		err = sm.EndSession(sessions[0].ID)
		if err != nil {
			t.Fatalf("EndSession failed: %v", err)
		}
	}

	// Should have 2 active sessions
	sessions, err = sm.GetAllActiveSessions()
	if err != nil {
		t.Fatalf("GetAllActiveSessions after end failed: %v", err)
	}

	if len(sessions) != 2 {
		t.Errorf("Sessions count after end = %v, expected 2", len(sessions))
	}
}

// TestCheckAndTimeoutSessions tests session timeout handling
func TestCheckAndTimeoutSessions(t *testing.T) {
	tempDir := t.TempDir()
	storageConfig := &storage.StorageConfig{CustomPath: tempDir}
	storageMgr, err := storage.NewStorageManager(storageConfig)
	if err != nil {
		t.Fatalf("Failed to create storage manager: %v", err)
	}

	sm, err := NewSessionManager(storageMgr)
	if err != nil {
		t.Fatalf("NewSessionManager failed: %v", err)
	}

	// Set very short timeout for testing
	sm.SetSessionTimeout(1 * time.Millisecond)

	// Start a session
	metadata, err := sm.StartSession("/test", "test")
	if err != nil {
		t.Fatalf("StartSession failed: %v", err)
	}

	sessionID := metadata.ID

	// Wait for timeout
	time.Sleep(10 * time.Millisecond)

	// Check and timeout sessions
	err = sm.CheckAndTimeoutSessions()
	if err != nil {
		t.Fatalf("CheckAndTimeoutSessions failed: %v", err)
	}

	// Verify session is timed out (should not be in active sessions)
	activeSessions, err := sm.GetAllActiveSessions()
	if err != nil {
		t.Fatalf("GetAllActiveSessions failed: %v", err)
	}

	for _, session := range activeSessions {
		if session.ID == sessionID {
			t.Error("Session should be timed out and not in active sessions")
		}
	}

	// Verify session is in completed storage with timeout status
	completedPath := filepath.Join(tempDir, "conversations", "completed", sessionID, "metadata.json")
	data, err := os.ReadFile(completedPath)
	if err != nil {
		t.Fatalf("Failed to read completed metadata: %v", err)
	}

	var timedOutMetadata SessionMetadata
	if err := json.Unmarshal(data, &timedOutMetadata); err != nil {
		t.Fatalf("Failed to unmarshal metadata: %v", err)
	}

	if timedOutMetadata.Status != StatusTimeout {
		t.Errorf("Status = %v, expected %v", timedOutMetadata.Status, StatusTimeout)
	}
}

// TestGetSessionEvents tests retrieving session events
func TestGetSessionEvents(t *testing.T) {
	tempDir := t.TempDir()
	storageConfig := &storage.StorageConfig{CustomPath: tempDir}
	storageMgr, err := storage.NewStorageManager(storageConfig)
	if err != nil {
		t.Fatalf("Failed to create storage manager: %v", err)
	}

	sm, err := NewSessionManager(storageMgr)
	if err != nil {
		t.Fatalf("NewSessionManager failed: %v", err)
	}

	// Start a session
	metadata, err := sm.StartSession("/test", "test")
	if err != nil {
		t.Fatalf("StartSession failed: %v", err)
	}

	sessionID := metadata.ID

	// Record multiple events
	eventTypes := []EventType{EventSessionStart, EventUserPrompt, EventClaudeResponse, EventSessionEnd}
	for i, eventType := range eventTypes {
		data := json.RawMessage([]byte(`{"index": ` + string(rune('0'+i)) + `}`))
		_, err := sm.RecordEvent(sessionID, eventType, data)
		if err != nil {
			t.Fatalf("RecordEvent %d failed: %v", i, err)
		}
	}

	// Get session events
	events, err := sm.GetSessionEvents(sessionID)
	if err != nil {
		t.Fatalf("GetSessionEvents failed: %v", err)
	}

	if len(events) != len(eventTypes) {
		t.Errorf("Events count = %v, expected %v", len(events), len(eventTypes))
	}

	// Verify event order and types
	for i, event := range events {
		if event.EventType != eventTypes[i] {
			t.Errorf("Event[%d] type = %v, expected %v", i, event.EventType, eventTypes[i])
		}

		if event.SequenceNum != i+1 {
			t.Errorf("Event[%d] sequence = %v, expected %v", i, event.SequenceNum, i+1)
		}
	}

	// End session and verify events are still retrievable
	err = sm.EndSession(sessionID)
	if err != nil {
		t.Fatalf("EndSession failed: %v", err)
	}

	eventsAfterEnd, err := sm.GetSessionEvents(sessionID)
	if err != nil {
		t.Fatalf("GetSessionEvents after end failed: %v", err)
	}

	if len(eventsAfterEnd) != len(events) {
		t.Errorf("Events count after end = %v, expected %v", len(eventsAfterEnd), len(events))
	}
}

// TestConcurrentSessions tests handling multiple concurrent sessions
func TestConcurrentSessions(t *testing.T) {
	tempDir := t.TempDir()
	storageConfig := &storage.StorageConfig{CustomPath: tempDir}
	storageMgr, err := storage.NewStorageManager(storageConfig)
	if err != nil {
		t.Fatalf("Failed to create storage manager: %v", err)
	}

	sm, err := NewSessionManager(storageMgr)
	if err != nil {
		t.Fatalf("NewSessionManager failed: %v", err)
	}

	// Start multiple concurrent sessions
	sessions := make([]*SessionMetadata, 5)
	for i := 0; i < 5; i++ {
		metadata, err := sm.StartSession(
			filepath.Join("/test", string(rune('a'+i))),
			string(rune('a'+i)),
		)
		if err != nil {
			t.Fatalf("StartSession %d failed: %v", i, err)
		}
		sessions[i] = metadata
	}

	// Record events for each session
	for i, session := range sessions {
		for j := 0; j < 3; j++ {
			data := json.RawMessage([]byte(`{"session": ` + string(rune('0'+i)) + `, "event": ` + string(rune('0'+j)) + `}`))
			_, err := sm.RecordEvent(session.ID, EventUserPrompt, data)
			if err != nil {
				t.Fatalf("RecordEvent for session %d event %d failed: %v", i, j, err)
			}
		}
	}

	// Verify all sessions have correct event counts
	for i, session := range sessions {
		metadata, err := sm.GetActiveSession(session.ID)
		if err != nil {
			t.Fatalf("GetActiveSession %d failed: %v", i, err)
		}

		if metadata.EventCount != 3 {
			t.Errorf("Session %d EventCount = %v, expected 3", i, metadata.EventCount)
		}
	}

	// End some sessions
	for i := 0; i < 2; i++ {
		err := sm.EndSession(sessions[i].ID)
		if err != nil {
			t.Fatalf("EndSession %d failed: %v", i, err)
		}
	}

	// Verify active session count
	activeSessions, err := sm.GetAllActiveSessions()
	if err != nil {
		t.Fatalf("GetAllActiveSessions failed: %v", err)
	}

	if len(activeSessions) != 3 {
		t.Errorf("Active sessions count = %v, expected 3", len(activeSessions))
	}
}


// Benchmark tests for performance validation

func BenchmarkStartSession(b *testing.B) {
	tempDir := b.TempDir()
	storageConfig := &storage.StorageConfig{CustomPath: tempDir}
	storageMgr, err := storage.NewStorageManager(storageConfig)
	if err != nil {
		b.Fatalf("Failed to create storage manager: %v", err)
	}

	sm, err := NewSessionManager(storageMgr)
	if err != nil {
		b.Fatalf("NewSessionManager failed: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := sm.StartSession("/test", "test")
		if err != nil {
			b.Fatalf("StartSession failed: %v", err)
		}
	}
}

func BenchmarkRecordEvent(b *testing.B) {
	tempDir := b.TempDir()
	storageConfig := &storage.StorageConfig{CustomPath: tempDir}
	storageMgr, err := storage.NewStorageManager(storageConfig)
	if err != nil {
		b.Fatalf("Failed to create storage manager: %v", err)
	}

	sm, err := NewSessionManager(storageMgr)
	if err != nil {
		b.Fatalf("NewSessionManager failed: %v", err)
	}

	metadata, err := sm.StartSession("/test", "test")
	if err != nil {
		b.Fatalf("StartSession failed: %v", err)
	}

	eventData := json.RawMessage(`{"message": "test event"}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := sm.RecordEvent(metadata.ID, EventUserPrompt, eventData)
		if err != nil {
			b.Fatalf("RecordEvent failed: %v", err)
		}
	}
}