package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"context-extender/internal/database"
	"github.com/spf13/cobra"
)

// captureRootCmd represents the root capture command
var captureRootCmd = &cobra.Command{
	Use:   "capture",
	Short: "Capture conversation events for Claude Code integration",
	Long: `Capture conversation events from Claude Code hooks and store them in the database.

This command is primarily used by Claude Code hooks for automatic conversation capture.
It accepts event types via the --event flag and routes them to appropriate handlers.

Supported events:
  - session-start: Start of a new Claude Code session
  - user-prompt: User prompt submitted to Claude
  - claude-response: Claude's response (including tool usage)
  - session-end: End of Claude Code session

Example usage:
  context-extender capture --event=session-start
  context-extender capture --event=user-prompt --data='{"message":"Hello"}'
  context-extender capture --event=claude-response
  context-extender capture --event=session-end`,
	RunE: func(cmd *cobra.Command, args []string) error {
		event, _ := cmd.Flags().GetString("event")
		data, _ := cmd.Flags().GetString("data")

		if event == "" {
			return fmt.Errorf("--event flag is required")
		}

		// Initialize database with new manager approach
		config := database.DefaultDatabaseConfig()
		manager := database.NewManager(config)

		ctx := context.Background()
		if err := manager.Initialize(ctx); err != nil {
			return fmt.Errorf("failed to initialize database: %w", err)
		}
		defer manager.Close()

		// Route to appropriate handler based on event type
		switch event {
		case "session-start":
			return handleSessionStart(ctx, manager, data)
		case "user-prompt":
			return handleUserPrompt(ctx, manager, data)
		case "claude-response":
			return handleClaudeResponse(ctx, manager, data)
		case "session-end":
			return handleSessionEnd(ctx, manager, data)
		default:
			return fmt.Errorf("unknown event type: %s", event)
		}
	},
}

// handleSessionStart processes session start events
func handleSessionStart(ctx context.Context, manager *database.Manager, data string) error {
	// Get session ID from environment or generate new one
	sessionID := os.Getenv("CLAUDE_SESSION_ID")
	if sessionID == "" {
		// Generate a new session ID based on working directory and timestamp
		wd, _ := os.Getwd()
		sessionID = fmt.Sprintf("%s_%d", filepath.Base(wd), os.Getpid())
	}

	backend, err := manager.GetBackend()
	if err != nil {
		return fmt.Errorf("failed to get backend: %w", err)
	}

	// Create session record
	session := &database.Session{
		ID:        sessionID,
		Status:    "active",
		Metadata:  fmt.Sprintf(`{"working_directory":"%s","data":%q}`, os.Getenv("PWD"), data),
	}

	if err := backend.CreateSession(ctx, session); err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	fmt.Printf("Session %s started\n", sessionID)
	return nil
}

// handleUserPrompt processes user prompt events
func handleUserPrompt(ctx context.Context, manager *database.Manager, data string) error {
	sessionID := os.Getenv("CLAUDE_SESSION_ID")
	if sessionID == "" {
		return fmt.Errorf("CLAUDE_SESSION_ID environment variable not set")
	}

	backend, err := manager.GetBackend()
	if err != nil {
		return fmt.Errorf("failed to get backend: %w", err)
	}

	// Create conversation record for user prompt
	conversation := &database.Conversation{
		ID:          fmt.Sprintf("%s_user_%d", sessionID, os.Getpid()),
		SessionID:   sessionID,
		MessageType: "user",
		Content:     data,
		Metadata:    "{}",
	}

	if err := backend.CreateConversation(ctx, conversation); err != nil {
		return fmt.Errorf("failed to create conversation: %w", err)
	}

	fmt.Printf("User prompt captured for session %s\n", sessionID)
	return nil
}

// handleClaudeResponse processes Claude response events
func handleClaudeResponse(ctx context.Context, manager *database.Manager, data string) error {
	sessionID := os.Getenv("CLAUDE_SESSION_ID")
	if sessionID == "" {
		return fmt.Errorf("CLAUDE_SESSION_ID environment variable not set")
	}

	backend, err := manager.GetBackend()
	if err != nil {
		return fmt.Errorf("failed to get backend: %w", err)
	}

	// Create conversation record for Claude response
	conversation := &database.Conversation{
		ID:          fmt.Sprintf("%s_claude_%d", sessionID, os.Getpid()),
		SessionID:   sessionID,
		MessageType: "assistant",
		Content:     data,
		Metadata:    "{}",
	}

	if err := backend.CreateConversation(ctx, conversation); err != nil {
		return fmt.Errorf("failed to create conversation: %w", err)
	}

	fmt.Printf("Claude response captured for session %s\n", sessionID)
	return nil
}

// handleSessionEnd processes session end events
func handleSessionEnd(ctx context.Context, manager *database.Manager, data string) error {
	sessionID := os.Getenv("CLAUDE_SESSION_ID")
	if sessionID == "" {
		return fmt.Errorf("CLAUDE_SESSION_ID environment variable not set")
	}

	backend, err := manager.GetBackend()
	if err != nil {
		return fmt.Errorf("failed to get backend: %w", err)
	}

	// Update session status to completed
	session, err := backend.GetSession(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}

	session.Status = "completed"
	if err := backend.UpdateSession(ctx, session); err != nil {
		return fmt.Errorf("failed to update session: %w", err)
	}

	fmt.Printf("Session %s ended\n", sessionID)
	return nil
}

func init() {
	rootCmd.AddCommand(captureRootCmd)

	// Add flags
	captureRootCmd.Flags().StringP("event", "e", "", "Event type (session-start, user-prompt, claude-response, session-end)")
	captureRootCmd.Flags().StringP("data", "d", "", "Event data (JSON or plain text)")
}