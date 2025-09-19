package session

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"context-extender/internal/database"
	"context-extender/internal/encryption"
)

func TestAdvancedSessionManager_CreateSession(t *testing.T) {
	manager, cleanup := setupTestManager(t)
	defer cleanup()

	metadata := map[string]interface{}{
		"project": "test-project",
		"user":    "test-user",
	}

	session, err := manager.CreateSession(context.Background(), metadata)
	if err != nil {
		t.Fatalf("CreateSession failed: %v", err)
	}

	if session.ID == "" {
		t.Error("Session ID should not be empty")
	}

	if session.Status != StatusActive {
		t.Errorf("Expected status %s, got %s", StatusActive, session.Status)
	}

	if session.Version != 1 {
		t.Errorf("Expected version 1, got %d", session.Version)
	}

	if session.Analytics == nil {
		t.Error("Analytics should be initialized")
	}
}

func TestAdvancedSessionManager_CreateSessionWithEncryption(t *testing.T) {
	manager, cleanup := setupTestManagerWithEncryption(t)
	defer cleanup()

	metadata := map[string]interface{}{
		"secret_key": "sensitive-data",
		"api_token":  "very-secret-token",
	}

	session, err := manager.CreateSession(context.Background(), metadata)
	if err != nil {
		t.Fatalf("CreateSession with encryption failed: %v", err)
	}

	// Verify the data was encrypted (metadata should be different from input)
	if val, ok := session.Metadata["secret_key"].(string); ok && val == "sensitive-data" {
		t.Error("Metadata should be encrypted, but appears to be plaintext")
	}

	// Retrieve and verify decryption works
	retrieved, err := manager.GetSession(context.Background(), session.ID)
	if err != nil {
		t.Fatalf("GetSession failed: %v", err)
	}

	// After retrieval, data should be decrypted
	if retrieved.Metadata["secret_key"] != "sensitive-data" {
		t.Error("Decrypted metadata doesn't match original")
	}
}

func TestAdvancedSessionManager_UpdateSessionWithVersioning(t *testing.T) {
	manager, cleanup := setupTestManager(t)
	defer cleanup()

	session, err := manager.CreateSession(context.Background(), map[string]interface{}{})
	if err != nil {
		t.Fatalf("CreateSession failed: %v", err)
	}

	originalVersion := session.Version

	// Update session
	session.EventCount = 10
	err = manager.UpdateSession(context.Background(), session)
	if err != nil {
		t.Fatalf("UpdateSession failed: %v", err)
	}

	if session.Version != originalVersion+1 {
		t.Errorf("Expected version %d, got %d", originalVersion+1, session.Version)
	}

	// Test version conflict
	oldSession := &EnhancedSession{
		ID:      session.ID,
		Version: originalVersion, // Old version
	}

	err = manager.UpdateSession(context.Background(), oldSession)
	if err == nil {
		t.Error("Expected version conflict error")
	}
}

func TestAdvancedSessionManager_SessionLifecycle(t *testing.T) {
	manager, cleanup := setupTestManager(t)
	defer cleanup()

	ctx := context.Background()

	// Create session
	session, err := manager.CreateSession(ctx, map[string]interface{}{
		"project": "lifecycle-test",
	})
	if err != nil {
		t.Fatalf("CreateSession failed: %v", err)
	}

	if session.Status != StatusActive {
		t.Errorf("New session should be active, got %s", session.Status)
	}

	// Add some events
	session.EventCount = 5
	now := time.Now()
	session.LastEventAt = &now

	err = manager.UpdateSession(ctx, session)
	if err != nil {
		t.Fatalf("UpdateSession failed: %v", err)
	}

	// Complete session
	session.Status = StatusCompleted
	endTime := time.Now()
	session.EndTime = &endTime

	err = manager.UpdateSession(ctx, session)
	if err != nil {
		t.Fatalf("Complete session failed: %v", err)
	}

	// Verify completion
	retrieved, err := manager.GetSession(ctx, session.ID)
	if err != nil {
		t.Fatalf("GetSession failed: %v", err)
	}

	if retrieved.Status != StatusCompleted {
		t.Errorf("Expected status %s, got %s", StatusCompleted, retrieved.Status)
	}

	if retrieved.EndTime == nil {
		t.Error("EndTime should be set for completed session")
	}
}

func TestAdvancedSessionManager_Analytics(t *testing.T) {
	manager, cleanup := setupTestManager(t)
	defer cleanup()

	ctx := context.Background()

	// Create multiple sessions with different patterns
	for i := 0; i < 5; i++ {
		session, err := manager.CreateSession(ctx, map[string]interface{}{
			"test_id": i,
		})
		if err != nil {
			t.Fatalf("CreateSession %d failed: %v", i, err)
		}

		// Simulate some activity
		session.EventCount = 10 + i*5
		session.Analytics.UserInteractions = i * 2
		session.Analytics.CommandExecutions = i * 3
		session.Analytics.EventPatterns["test_pattern"] = i + 1

		err = manager.UpdateSession(ctx, session)
		if err != nil {
			t.Fatalf("UpdateSession %d failed: %v", i, err)
		}
	}

	// Get global analytics
	analytics, err := manager.GetAnalytics(ctx)
	if err != nil {
		t.Fatalf("GetAnalytics failed: %v", err)
	}

	if analytics.TotalSessions != 5 {
		t.Errorf("Expected 5 total sessions, got %d", analytics.TotalSessions)
	}

	if analytics.ActiveSessions != 5 {
		t.Errorf("Expected 5 active sessions, got %d", analytics.ActiveSessions)
	}

	if len(analytics.TopEventPatterns) == 0 {
		t.Error("Should have event patterns")
	}
}

func TestAdvancedSessionManager_ExportImport(t *testing.T) {
	manager, cleanup := setupTestManager(t)
	defer cleanup()

	ctx := context.Background()

	// Create a session with rich data
	originalMetadata := map[string]interface{}{
		"project":     "export-test",
		"description": "Test session for export/import",
		"tags":        []string{"test", "export"},
	}

	session, err := manager.CreateSession(ctx, originalMetadata)
	if err != nil {
		t.Fatalf("CreateSession failed: %v", err)
	}

	// Add some context
	session.Context = []byte(`{"events": [{"type": "test", "data": "sample"}]}`)
	session.EventCount = 3
	session.Analytics.UserInteractions = 5

	err = manager.UpdateSession(ctx, session)
	if err != nil {
		t.Fatalf("UpdateSession failed: %v", err)
	}

	// Test JSON export
	exported, err := manager.ExportSession(ctx, session.ID, FormatJSON)
	if err != nil {
		t.Fatalf("ExportSession failed: %v", err)
	}

	if len(exported) == 0 {
		t.Error("Exported data should not be empty")
	}

	// Verify JSON structure
	var exportedSession EnhancedSession
	err = json.Unmarshal(exported, &exportedSession)
	if err != nil {
		t.Fatalf("Exported JSON is invalid: %v", err)
	}

	// Test import
	importedSession, err := manager.ImportSession(ctx, exported, FormatJSON)
	if err != nil {
		t.Fatalf("ImportSession failed: %v", err)
	}

	// Verify imported data (should have new ID but same content)
	if importedSession.ID == session.ID {
		t.Error("Imported session should have new ID")
	}

	if importedSession.EventCount != session.EventCount {
		t.Errorf("Event count mismatch: expected %d, got %d",
			session.EventCount, importedSession.EventCount)
	}

	// Test Markdown export
	markdownExport, err := manager.ExportSession(ctx, session.ID, FormatMarkdown)
	if err != nil {
		t.Fatalf("Markdown export failed: %v", err)
	}

	markdownStr := string(markdownExport)
	if !containsString(markdownStr, session.ID) {
		t.Error("Markdown export should contain session ID")
	}

	if !containsString(markdownStr, "export-test") {
		t.Error("Markdown export should contain project name")
	}
}

func TestAdvancedSessionManager_MaxActiveSessions(t *testing.T) {
	// Create manager with low limit
	config := DefaultManagerConfig()
	config.MaxActiveSessions = 2

	manager, cleanup := setupTestManagerWithConfig(t, config)
	defer cleanup()

	ctx := context.Background()

	// Create maximum allowed sessions
	for i := 0; i < 2; i++ {
		_, err := manager.CreateSession(ctx, map[string]interface{}{
			"test_id": i,
		})
		if err != nil {
			t.Fatalf("CreateSession %d failed: %v", i, err)
		}
	}

	// Try to create one more (should fail)
	_, err := manager.CreateSession(ctx, map[string]interface{}{
		"test_id": 3,
	})
	if err == nil {
		t.Error("Should not be able to create session beyond limit")
	}

	// Complete one session
	sessions, err := manager.ListSessions(ctx, SessionFilter{Limit: 1})
	if err != nil {
		t.Fatalf("ListSessions failed: %v", err)
	}

	if len(sessions) > 0 {
		sessions[0].Status = StatusCompleted
		err = manager.UpdateSession(ctx, sessions[0])
		if err != nil {
			t.Fatalf("UpdateSession failed: %v", err)
		}

		// Now should be able to create another
		_, err = manager.CreateSession(ctx, map[string]interface{}{
			"test_id": 4,
		})
		if err != nil {
			t.Errorf("Should be able to create session after completing one: %v", err)
		}
	}
}

func TestAdvancedSessionManager_Metrics(t *testing.T) {
	manager, cleanup := setupTestManager(t)
	defer cleanup()

	ctx := context.Background()

	initialMetrics := manager.GetMetrics()
	if initialMetrics.TotalSessions != 0 {
		t.Errorf("Initial total sessions should be 0, got %d", initialMetrics.TotalSessions)
	}

	// Create and complete sessions
	for i := 0; i < 3; i++ {
		session, err := manager.CreateSession(ctx, map[string]interface{}{})
		if err != nil {
			t.Fatalf("CreateSession %d failed: %v", i, err)
		}

		if i < 2 { // Complete first two sessions
			session.Status = StatusCompleted
			endTime := time.Now()
			session.EndTime = &endTime
			err = manager.UpdateSession(ctx, session)
			if err != nil {
				t.Fatalf("Complete session %d failed: %v", i, err)
			}
		}
	}

	metrics := manager.GetMetrics()

	if metrics.TotalSessions != 3 {
		t.Errorf("Expected 3 total sessions, got %d", metrics.TotalSessions)
	}

	if metrics.ActiveSessions != 1 {
		t.Errorf("Expected 1 active session, got %d", metrics.ActiveSessions)
	}

	if metrics.CompletedSessions != 2 {
		t.Errorf("Expected 2 completed sessions, got %d", metrics.CompletedSessions)
	}
}

func TestAdvancedSessionManager_SessionFilter(t *testing.T) {
	manager, cleanup := setupTestManager(t)
	defer cleanup()

	ctx := context.Background()

	// Create sessions with different statuses
	activeSession, _ := manager.CreateSession(ctx, map[string]interface{}{"type": "active"})

	completedSession, _ := manager.CreateSession(ctx, map[string]interface{}{"type": "completed"})
	completedSession.Status = StatusCompleted
	manager.UpdateSession(ctx, completedSession)

	timeoutSession, _ := manager.CreateSession(ctx, map[string]interface{}{"type": "timeout"})
	timeoutSession.Status = StatusTimeout
	manager.UpdateSession(ctx, timeoutSession)

	// Test status filter
	activeSessions, err := manager.ListSessions(ctx, SessionFilter{
		Status: StatusActive,
	})
	if err != nil {
		t.Fatalf("ListSessions failed: %v", err)
	}

	if len(activeSessions) != 1 {
		t.Errorf("Expected 1 active session, got %d", len(activeSessions))
	}

	// Test limit
	limitedSessions, err := manager.ListSessions(ctx, SessionFilter{
		Limit: 2,
	})
	if err != nil {
		t.Fatalf("ListSessions with limit failed: %v", err)
	}

	if len(limitedSessions) != 2 {
		t.Errorf("Expected 2 sessions with limit, got %d", len(limitedSessions))
	}
}

// Benchmark tests

func BenchmarkAdvancedSessionManager_CreateSession(b *testing.B) {
	manager, cleanup := setupTestManager(&testing.T{})
	defer cleanup()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := manager.CreateSession(context.Background(), map[string]interface{}{
			"benchmark_id": i,
		})
		if err != nil {
			b.Fatalf("CreateSession failed: %v", err)
		}
	}
}

func BenchmarkAdvancedSessionManager_GetSession(b *testing.B) {
	manager, cleanup := setupTestManager(&testing.T{})
	defer cleanup()

	// Create a session first
	session, err := manager.CreateSession(context.Background(), map[string]interface{}{})
	if err != nil {
		b.Fatalf("Setup CreateSession failed: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := manager.GetSession(context.Background(), session.ID)
		if err != nil {
			b.Fatalf("GetSession failed: %v", err)
		}
	}
}

func BenchmarkAdvancedSessionManager_UpdateSession(b *testing.B) {
	manager, cleanup := setupTestManager(&testing.T{})
	defer cleanup()

	// Create a session first
	session, err := manager.CreateSession(context.Background(), map[string]interface{}{})
	if err != nil {
		b.Fatalf("Setup CreateSession failed: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		session.EventCount = i
		err := manager.UpdateSession(context.Background(), session)
		if err != nil {
			b.Fatalf("UpdateSession failed: %v", err)
		}
	}
}

func BenchmarkAdvancedSessionManager_ExportSession(b *testing.B) {
	manager, cleanup := setupTestManager(&testing.T{})
	defer cleanup()

	// Create a session with some data
	session, err := manager.CreateSession(context.Background(), map[string]interface{}{
		"large_data": generateLargeMetadata(),
	})
	if err != nil {
		b.Fatalf("Setup CreateSession failed: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := manager.ExportSession(context.Background(), session.ID, FormatJSON)
		if err != nil {
			b.Fatalf("ExportSession failed: %v", err)
		}
	}
}

// Helper functions

func setupTestManager(t *testing.T) (*AdvancedSessionManager, func()) {
	return setupTestManagerWithConfig(t, nil)
}

func setupTestManagerWithEncryption(t *testing.T) (*AdvancedSessionManager, func()) {
	config := DefaultManagerConfig()
	config.EnableEncryption = true

	db, dbCleanup := setupTestDatabase(t)

	// Setup encryption
	encConfig := encryption.DefaultEncryptionConfig()
	encConfig.KeyPath = "./test_session_encryption.key"
	encConfig.Enabled = true

	encProvider, err := encryption.NewAESGCMProvider(encConfig)
	if err != nil {
		t.Fatalf("Failed to create encryption provider: %v", err)
	}

	if err := encProvider.Initialize(); err != nil {
		t.Fatalf("Failed to initialize encryption: %v", err)
	}

	manager, err := NewAdvancedSessionManager(db, encProvider, config)
	if err != nil {
		t.Fatalf("Failed to create session manager: %v", err)
	}

	return manager, func() {
		encProvider.Close()
		dbCleanup()
	}
}

func setupTestManagerWithConfig(t *testing.T, config *ManagerConfig) (*AdvancedSessionManager, func()) {
	if config == nil {
		config = DefaultManagerConfig()
		config.EnableAutoSave = false // Disable for tests
	}

	db, cleanup := setupTestDatabase(t)

	manager, err := NewAdvancedSessionManager(db, nil, config)
	if err != nil {
		t.Fatalf("Failed to create session manager: %v", err)
	}

	return manager, cleanup
}

func setupTestDatabase(t *testing.T) (database.Database, func()) {
	// Use in-memory database for tests
	config := &database.DatabaseConfig{
		Path: ":memory:",
	}

	db, err := database.New(config)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// Create sessions table
	_, err = db.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS sessions (
			id TEXT PRIMARY KEY,
			status TEXT NOT NULL,
			start_time DATETIME NOT NULL,
			end_time DATETIME,
			metadata TEXT,
			context BLOB,
			event_count INTEGER DEFAULT 0,
			last_event_at DATETIME,
			tags TEXT,
			version INTEGER DEFAULT 1,
			analytics TEXT
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create sessions table: %v", err)
	}

	return db, func() {
		db.Close()
	}
}

func generateLargeMetadata() map[string]interface{} {
	metadata := make(map[string]interface{})
	for i := 0; i < 100; i++ {
		metadata[fmt.Sprintf("key_%d", i)] = fmt.Sprintf("value_%d_with_some_longer_content", i)
	}
	return metadata
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && contains(s, substr)))
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Add missing types and functions

type SessionFilter struct {
	Status    SessionStatus
	StartTime *time.Time
	EndTime   *time.Time
	Tags      []string
	Limit     int
	Offset    int
}

func (m *AdvancedSessionManager) ListSessions(ctx context.Context, filter SessionFilter) ([]*EnhancedSession, error) {
	// Implementation for listing sessions with filters
	// This would build appropriate SQL query based on filter parameters
	query := "SELECT id, status, start_time, end_time, metadata, context, event_count, last_event_at, tags, version, analytics FROM sessions WHERE 1=1"
	args := []interface{}{}

	if filter.Status != "" {
		query += " AND status = ?"
		args = append(args, filter.Status)
	}

	if filter.StartTime != nil {
		query += " AND start_time >= ?"
		args = append(args, filter.StartTime)
	}

	if filter.EndTime != nil {
		query += " AND end_time <= ?"
		args = append(args, filter.EndTime)
	}

	query += " ORDER BY start_time DESC"

	if filter.Limit > 0 {
		query += " LIMIT ?"
		args = append(args, filter.Limit)
	}

	rows, err := m.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*EnhancedSession
	for rows.Next() {
		session := &EnhancedSession{}
		var metadataJSON, tagsJSON, analyticsJSON []byte

		err := rows.Scan(
			&session.ID,
			&session.Status,
			&session.StartTime,
			&session.EndTime,
			&metadataJSON,
			&session.Context,
			&session.EventCount,
			&session.LastEventAt,
			&tagsJSON,
			&session.Version,
			&analyticsJSON,
		)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(metadataJSON, &session.Metadata); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(tagsJSON, &session.Tags); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(analyticsJSON, &session.Analytics); err != nil {
			session.Analytics = &SessionAnalytics{
				EventPatterns: make(map[string]int),
				TokenUsage:    make(map[string]int),
			}
		}

		// Decrypt if needed
		if m.encryption != nil {
			decryptedSession, err := m.decryptSessionIfNeeded(session)
			if err != nil {
				return nil, err
			}
			session = decryptedSession
		}

		sessions = append(sessions, session)
	}

	return sessions, nil
}