package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// Types are now defined in types.go to avoid circular imports

func CreateSession(session *Session) error {
	db, err := GetConnection()
	if err != nil {
		return err
	}

	metadataJSON := session.Metadata

	query := `
		INSERT INTO sessions (id, created_at, updated_at, status, metadata)
		VALUES (?, ?, ?, ?, ?)
	`

	_, err = db.Exec(query,
		session.ID,
		session.CreatedAt.Format(time.RFC3339),
		session.UpdatedAt.Format(time.RFC3339),
		session.Status,
		metadataJSON,
	)

	return err
}

func UpdateSession(sessionID string, status string, metadata map[string]string) error {
	db, err := GetConnection()
	if err != nil {
		return err
	}

	metadataJSON := ""
	if metadata != nil {
		metadataBytes, err := json.Marshal(metadata)
		if err != nil {
			return fmt.Errorf("failed to marshal metadata: %w", err)
		}
		metadataJSON = string(metadataBytes)
	}

	query := `
		UPDATE sessions
		SET updated_at = ?, status = ?, metadata = ?
		WHERE id = ?
	`

	_, err = db.Exec(query,
		time.Now().Format(time.RFC3339),
		status,
		metadataJSON,
		sessionID,
	)

	return err
}

func GetSession(sessionID string) (*Session, error) {
	db, err := GetConnection()
	if err != nil {
		return nil, err
	}

	query := `
		SELECT id, created_at, updated_at, status, metadata
		FROM sessions
		WHERE id = ?
	`

	row := db.QueryRow(query, sessionID)

	var session Session
	var createdAt, updatedAt, metadataJSON string

	err = row.Scan(&session.ID, &createdAt, &updatedAt, &session.Status, &metadataJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	session.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	session.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

	if metadataJSON != "" {
		if err := json.Unmarshal([]byte(metadataJSON), &session.Metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}
	}

	return &session, nil
}

func CreateEvent(event *Event) error {
	db, err := GetConnection()
	if err != nil {
		return err
	}

	query := `
		INSERT INTO events (session_id, event_type, event_data, timestamp, sequence_number)
		VALUES (?, ?, ?, ?, ?)
	`

	result, err := db.Exec(query,
		event.SessionID,
		event.EventType,
		event.Data,
		event.Timestamp.Format(time.RFC3339),
		event.SequenceNum,
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	event.ID = fmt.Sprintf("%d", id)
	return nil
}

func GetEventsBySession(sessionID string) ([]*Event, error) {
	db, err := GetConnection()
	if err != nil {
		return nil, err
	}

	query := `
		SELECT id, session_id, event_type, event_data, timestamp, sequence_number
		FROM events
		WHERE session_id = ?
		ORDER BY sequence_number ASC
	`

	rows, err := db.Query(query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*Event
	for rows.Next() {
		var event Event
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
			return nil, err
		}

		event.Timestamp, _ = time.Parse(time.RFC3339, timestampStr)
		events = append(events, &event)
	}

	return events, rows.Err()
}

func CreateConversation(conversation *Conversation) error {
	db, err := GetConnection()
	if err != nil {
		return err
	}

	query := `
		INSERT INTO conversations (session_id, message_type, content, timestamp, token_count, model_info)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	result, err := db.Exec(query,
		conversation.SessionID,
		conversation.MessageType,
		conversation.Content,
		conversation.Timestamp.Format(time.RFC3339),
		conversation.TokenCount,
		conversation.Model,
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	conversation.ID = fmt.Sprintf("%d", id)
	return nil
}

func GetConversationsBySession(sessionID string) ([]*Conversation, error) {
	db, err := GetConnection()
	if err != nil {
		return nil, err
	}

	query := `
		SELECT id, session_id, message_type, content, timestamp, token_count, model_info
		FROM conversations
		WHERE session_id = ?
		ORDER BY timestamp ASC
	`

	rows, err := db.Query(query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversations []*Conversation
	for rows.Next() {
		var conversation Conversation
		var timestampStr string

		err := rows.Scan(
			&conversation.ID,
			&conversation.SessionID,
			&conversation.MessageType,
			&conversation.Content,
			&timestampStr,
			&conversation.TokenCount,
			&conversation.Model,
		)
		if err != nil {
			return nil, err
		}

		conversation.Timestamp, _ = time.Parse(time.RFC3339, timestampStr)
		conversations = append(conversations, &conversation)
	}

	return conversations, rows.Err()
}