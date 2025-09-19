package graphql

import (
	"context-extender/internal/database"
	"fmt"
	"strings"
	"time"

	"github.com/graphql-go/graphql"
)

// SetupResolvers configures all GraphQL resolvers
func SetupResolvers() {
	// Session resolver
	SessionType.AddFieldConfig("events", &graphql.Field{
		Type: graphql.NewList(EventType),
		Args: graphql.FieldConfigArgument{
			"limit": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"offset": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			session, ok := p.Source.(*database.Session)
			if !ok {
				return nil, fmt.Errorf("invalid session type")
			}

			events, err := database.GetEventsBySession(session.ID)
			if err != nil {
				return nil, err
			}

			// Apply pagination if specified
			limit := getIntArg(p.Args, "limit", len(events))
			offset := getIntArg(p.Args, "offset", 0)

			end := offset + limit
			if end > len(events) {
				end = len(events)
			}
			if offset > len(events) {
				return []*database.Event{}, nil
			}

			return events[offset:end], nil
		},
	})

	SessionType.AddFieldConfig("conversations", &graphql.Field{
		Type: graphql.NewList(ConversationType),
		Args: graphql.FieldConfigArgument{
			"limit": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"offset": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			session, ok := p.Source.(*database.Session)
			if !ok {
				return nil, fmt.Errorf("invalid session type")
			}

			conversations, err := database.GetConversationsBySession(session.ID)
			if err != nil {
				return nil, err
			}

			// Apply pagination if specified
			limit := getIntArg(p.Args, "limit", len(conversations))
			offset := getIntArg(p.Args, "offset", 0)

			end := offset + limit
			if end > len(conversations) {
				end = len(conversations)
			}
			if offset > len(conversations) {
				return []*database.Conversation{}, nil
			}

			return conversations[offset:end], nil
		},
	})

	// Query resolvers
	QueryType.AddFieldConfig("session", &graphql.Field{
		Type: SessionType,
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			id, ok := p.Args["id"].(string)
			if !ok {
				return nil, fmt.Errorf("id argument is required")
			}

			session, err := database.GetSession(id)
			if err != nil {
				return nil, err
			}

			return session, nil
		},
	})

	QueryType.AddFieldConfig("sessions", &graphql.Field{
		Type: graphql.NewList(SessionType),
		Args: graphql.FieldConfigArgument{
			"limit": &graphql.ArgumentConfig{
				Type:         graphql.Int,
				DefaultValue: 50,
			},
			"offset": &graphql.ArgumentConfig{
				Type:         graphql.Int,
				DefaultValue: 0,
			},
			"status": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"sortBy": &graphql.ArgumentConfig{
				Type:         graphql.String,
				DefaultValue: "created_at",
			},
			"sortOrder": &graphql.ArgumentConfig{
				Type:         graphql.String,
				DefaultValue: "DESC",
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			limit := getIntArg(p.Args, "limit", 50)
			offset := getIntArg(p.Args, "offset", 0)
			status := getStringArg(p.Args, "status", "")
			sortBy := getStringArg(p.Args, "sortBy", "created_at")
			sortOrder := getStringArg(p.Args, "sortOrder", "DESC")

			return querySessionsWithFilters(limit, offset, status, sortBy, sortOrder)
		},
	})

	QueryType.AddFieldConfig("events", &graphql.Field{
		Type: graphql.NewList(EventType),
		Args: graphql.FieldConfigArgument{
			"sessionId": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"eventType": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"limit": &graphql.ArgumentConfig{
				Type:         graphql.Int,
				DefaultValue: 100,
			},
			"offset": &graphql.ArgumentConfig{
				Type:         graphql.Int,
				DefaultValue: 0,
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			sessionId := getStringArg(p.Args, "sessionId", "")
			eventType := getStringArg(p.Args, "eventType", "")
			limit := getIntArg(p.Args, "limit", 100)
			offset := getIntArg(p.Args, "offset", 0)

			return queryEventsWithFilters(sessionId, eventType, limit, offset)
		},
	})

	QueryType.AddFieldConfig("conversations", &graphql.Field{
		Type: graphql.NewList(ConversationType),
		Args: graphql.FieldConfigArgument{
			"sessionId": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"messageType": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"limit": &graphql.ArgumentConfig{
				Type:         graphql.Int,
				DefaultValue: 100,
			},
			"offset": &graphql.ArgumentConfig{
				Type:         graphql.Int,
				DefaultValue: 0,
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			sessionId := getStringArg(p.Args, "sessionId", "")
			messageType := getStringArg(p.Args, "messageType", "")
			limit := getIntArg(p.Args, "limit", 100)
			offset := getIntArg(p.Args, "offset", 0)

			return queryConversationsWithFilters(sessionId, messageType, limit, offset)
		},
	})

	QueryType.AddFieldConfig("search", &graphql.Field{
		Type: SearchResultType,
		Args: graphql.FieldConfigArgument{
			"query": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"limit": &graphql.ArgumentConfig{
				Type:         graphql.Int,
				DefaultValue: 50,
			},
			"searchSessions": &graphql.ArgumentConfig{
				Type:         graphql.Boolean,
				DefaultValue: true,
			},
			"searchConversations": &graphql.ArgumentConfig{
				Type:         graphql.Boolean,
				DefaultValue: true,
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			query, ok := p.Args["query"].(string)
			if !ok {
				return nil, fmt.Errorf("query argument is required")
			}

			limit := getIntArg(p.Args, "limit", 50)
			searchSessions := getBoolArg(p.Args, "searchSessions", true)
			searchConversations := getBoolArg(p.Args, "searchConversations", true)

			return performSearch(query, limit, searchSessions, searchConversations)
		},
	})

	QueryType.AddFieldConfig("stats", &graphql.Field{
		Type: StatsType,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return getDatabaseStats()
		},
	})
}

// Helper functions for argument extraction
func getIntArg(args map[string]interface{}, key string, defaultValue int) int {
	if val, ok := args[key]; ok {
		if intVal, ok := val.(int); ok {
			return intVal
		}
	}
	return defaultValue
}

func getStringArg(args map[string]interface{}, key string, defaultValue string) string {
	if val, ok := args[key]; ok {
		if strVal, ok := val.(string); ok {
			return strVal
		}
	}
	return defaultValue
}

func getBoolArg(args map[string]interface{}, key string, defaultValue bool) bool {
	if val, ok := args[key]; ok {
		if boolVal, ok := val.(bool); ok {
			return boolVal
		}
	}
	return defaultValue
}

// Query implementations
func querySessionsWithFilters(limit, offset int, status, sortBy, sortOrder string) ([]*database.Session, error) {
	db, err := database.GetConnection()
	if err != nil {
		return nil, err
	}

	query := "SELECT id, created_at, updated_at, status, metadata FROM sessions"
	args := []interface{}{}

	// Add WHERE clause if status filter is provided
	if status != "" {
		query += " WHERE status = ?"
		args = append(args, status)
	}

	// Add ORDER BY clause
	validSortFields := map[string]bool{
		"created_at": true,
		"updated_at": true,
		"status":     true,
		"id":         true,
	}
	if !validSortFields[sortBy] {
		sortBy = "created_at"
	}
	if sortOrder != "ASC" && sortOrder != "DESC" {
		sortOrder = "DESC"
	}
	query += fmt.Sprintf(" ORDER BY %s %s", sortBy, sortOrder)

	// Add LIMIT and OFFSET
	query += " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*database.Session
	for rows.Next() {
		var session database.Session
		var createdAt, updatedAt, metadataJSON string

		err := rows.Scan(&session.ID, &createdAt, &updatedAt, &session.Status, &metadataJSON)
		if err != nil {
			continue
		}

		// Parse timestamps
		session.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		session.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

		sessions = append(sessions, &session)
	}

	return sessions, rows.Err()
}

func queryEventsWithFilters(sessionId, eventType string, limit, offset int) ([]*database.Event, error) {
	db, err := database.GetConnection()
	if err != nil {
		return nil, err
	}

	query := "SELECT id, session_id, event_type, event_data, timestamp, sequence_number FROM events"
	args := []interface{}{}
	conditions := []string{}

	if sessionId != "" {
		conditions = append(conditions, "session_id = ?")
		args = append(args, sessionId)
	}

	if eventType != "" {
		conditions = append(conditions, "event_type = ?")
		args = append(args, eventType)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY timestamp DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*database.Event
	for rows.Next() {
		var event database.Event
		var timestampStr string

		err := rows.Scan(
			&event.ID,
			&event.SessionID,
			&event.EventType,
			&event.Data,
			&timestampStr,
			&event.SequenceNum,
		)
		if err != nil {
			continue
		}

		event.Timestamp, _ = time.Parse(time.RFC3339, timestampStr)
		events = append(events, &event)
	}

	return events, rows.Err()
}

func queryConversationsWithFilters(sessionId, messageType string, limit, offset int) ([]*database.Conversation, error) {
	db, err := database.GetConnection()
	if err != nil {
		return nil, err
	}

	query := "SELECT id, session_id, message_type, content, timestamp, token_count, model_info FROM conversations"
	args := []interface{}{}
	conditions := []string{}

	if sessionId != "" {
		conditions = append(conditions, "session_id = ?")
		args = append(args, sessionId)
	}

	if messageType != "" {
		conditions = append(conditions, "message_type = ?")
		args = append(args, messageType)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY timestamp ASC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversations []*database.Conversation
	for rows.Next() {
		var conv database.Conversation
		var timestampStr string

		err := rows.Scan(
			&conv.ID,
			&conv.SessionID,
			&conv.MessageType,
			&conv.Content,
			&timestampStr,
			&conv.TokenCount,
			&conv.Model,
		)
		if err != nil {
			continue
		}

		conv.Timestamp, _ = time.Parse(time.RFC3339, timestampStr)
		conversations = append(conversations, &conv)
	}

	return conversations, rows.Err()
}

type SearchResult struct {
	Sessions      []*database.Session      `json:"sessions"`
	Conversations []*database.Conversation `json:"conversations"`
	TotalCount    int                      `json:"totalCount"`
}

func performSearch(query string, limit int, searchSessions, searchConversations bool) (*SearchResult, error) {
	result := &SearchResult{
		Sessions:      []*database.Session{},
		Conversations: []*database.Conversation{},
		TotalCount:    0,
	}

	db, err := database.GetConnection()
	if err != nil {
		return result, err
	}

	searchTerm := "%" + query + "%"

	// Search conversations
	if searchConversations {
		convQuery := `
			SELECT id, session_id, message_type, content, timestamp, token_count, model_info
			FROM conversations
			WHERE content LIKE ?
			ORDER BY timestamp DESC
			LIMIT ?
		`
		rows, err := db.Query(convQuery, searchTerm, limit/2)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var conv database.Conversation
				var timestampStr string

				err := rows.Scan(
					&conv.ID,
					&conv.SessionID,
					&conv.MessageType,
					&conv.Content,
					&timestampStr,
					&conv.TokenCount,
					&conv.Model,
				)
				if err == nil {
					conv.Timestamp, _ = time.Parse(time.RFC3339, timestampStr)
					result.Conversations = append(result.Conversations, &conv)
				}
			}
		}
	}

	// Search sessions (by metadata)
	if searchSessions {
		sessQuery := `
			SELECT id, created_at, updated_at, status, metadata
			FROM sessions
			WHERE metadata LIKE ? OR id LIKE ?
			ORDER BY created_at DESC
			LIMIT ?
		`
		rows, err := db.Query(sessQuery, searchTerm, searchTerm, limit/2)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var session database.Session
				var createdAt, updatedAt, metadataJSON string

				err := rows.Scan(&session.ID, &createdAt, &updatedAt, &session.Status, &metadataJSON)
				if err == nil {
					session.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
					session.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
					result.Sessions = append(result.Sessions, &session)
				}
			}
		}
	}

	result.TotalCount = len(result.Sessions) + len(result.Conversations)
	return result, nil
}

type DatabaseStats struct {
	TotalSessions      int    `json:"totalSessions"`
	TotalEvents        int    `json:"totalEvents"`
	TotalConversations int    `json:"totalConversations"`
	TotalImports       int    `json:"totalImports"`
	OldestSession      string `json:"oldestSession"`
	NewestSession      string `json:"newestSession"`
}

func getDatabaseStats() (*DatabaseStats, error) {
	stats := &DatabaseStats{}

	db, err := database.GetConnection()
	if err != nil {
		return stats, err
	}

	// Get counts
	db.QueryRow("SELECT COUNT(*) FROM sessions").Scan(&stats.TotalSessions)
	db.QueryRow("SELECT COUNT(*) FROM events").Scan(&stats.TotalEvents)
	db.QueryRow("SELECT COUNT(*) FROM conversations").Scan(&stats.TotalConversations)
	db.QueryRow("SELECT COUNT(*) FROM import_history").Scan(&stats.TotalImports)

	// Get oldest and newest sessions
	db.QueryRow("SELECT MIN(created_at) FROM sessions").Scan(&stats.OldestSession)
	db.QueryRow("SELECT MAX(created_at) FROM sessions").Scan(&stats.NewestSession)

	return stats, nil
}