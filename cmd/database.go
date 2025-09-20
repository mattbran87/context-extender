package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"context-extender/internal/database"
	"github.com/spf13/cobra"
)

var databaseCmd = &cobra.Command{
	Use:   "database",
	Short: "Database operations for conversation storage",
	Long: `Database operations for managing SQLite storage of Claude Code conversations.

This command provides database initialization, migration, and maintenance operations.`,
}

var initDbCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the database",
	Long: `Initialize the SQLite database with the required schema for storing conversations.

This creates the database file and runs all necessary migrations.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		config := database.DefaultDatabaseConfig()

		// Create directory if it doesn't exist
		dbDir := filepath.Dir(config.DatabasePath)
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			return fmt.Errorf("failed to create database directory: %w", err)
		}

		// Create database manager
		manager := database.NewManager(config)

		// Initialize database
		ctx := cmd.Context()
		if err := manager.Initialize(ctx); err != nil {
			return fmt.Errorf("failed to initialize database: %w", err)
		}
		defer manager.Close()

		// Get backend and create schema
		backend, err := manager.GetBackend()
		if err != nil {
			return fmt.Errorf("failed to get backend: %w", err)
		}

		// Create schema
		if err := backend.CreateSchema(ctx); err != nil {
			return fmt.Errorf("failed to create schema: %w", err)
		}

		// Run migrations
		if err := backend.MigrateSchema(ctx, 1); err != nil {
			return fmt.Errorf("failed to run migrations: %w", err)
		}

		fmt.Printf("Database initialized successfully at: %s\n", config.DatabasePath)
		fmt.Printf("Backend: %s\n", backend.GetBackendInfo().Name)
		return nil
	},
}

var migrateDbCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Long: `Run all pending database migrations to update the schema to the latest version.

This is safe to run multiple times and will only apply new migrations.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		config := database.DefaultDatabaseConfig()
		manager := database.NewManager(config)

		ctx := cmd.Context()
		if err := manager.Initialize(ctx); err != nil {
			return fmt.Errorf("failed to initialize database: %w", err)
		}
		defer manager.Close()

		backend, err := manager.GetBackend()
		if err != nil {
			return fmt.Errorf("failed to get backend: %w", err)
		}

		if err := backend.MigrateSchema(ctx, 1); err != nil {
			return fmt.Errorf("failed to run migrations: %w", err)
		}

		fmt.Println("Database migrations completed successfully")
		return nil
	},
}

var captureCmd = &cobra.Command{
	Use:   "capture",
	Short: "Capture conversation events to database",
	Long: `Capture conversation events and store them in the SQLite database.

This command is primarily used by Claude Code hooks for automatic conversation capture.`,
}

var sessionStartCmd = &cobra.Command{
	Use:   "session-start [session-id]",
	Short: "Capture session start event",
	Long: `Capture a session start event and store it in the database.

If no session ID is provided, a new UUID will be generated.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		config := database.DefaultDatabaseConfig()
		manager := database.NewManager(config)

		ctx := cmd.Context()
		if err := manager.Initialize(ctx); err != nil {
			return fmt.Errorf("failed to initialize database: %w", err)
		}
		defer manager.Close()

		sessionID := ""
		if len(args) > 0 {
			sessionID = args[0]
		}

		if err := database.HandleSessionStartHook(sessionID); err != nil {
			return fmt.Errorf("failed to handle session start: %w", err)
		}

		fmt.Printf("Session started successfully\n")
		return nil
	},
}

var userPromptCmd = &cobra.Command{
	Use:   "user-prompt [session-id] [message]",
	Short: "Capture user prompt event",
	Long: `Capture a user prompt event and store it in the database.

Both session ID and message are required parameters.`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		config := database.DefaultDatabaseConfig()
		manager := database.NewManager(config)

		ctx := cmd.Context()
		if err := manager.Initialize(ctx); err != nil {
			return fmt.Errorf("failed to initialize database: %w", err)
		}
		defer manager.Close()

		sessionID := args[0]
		message := args[1]

		if err := database.HandleUserPromptHook(sessionID, message); err != nil {
			return fmt.Errorf("failed to handle user prompt: %w", err)
		}

		fmt.Printf("User prompt captured successfully\n")
		return nil
	},
}

var claudeResponseCmd = &cobra.Command{
	Use:   "claude-response [session-id] [response]",
	Short: "Capture Claude response event",
	Long: `Capture a Claude response event and store it in the database.

Both session ID and response are required parameters.`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		config := database.DefaultDatabaseConfig()
		manager := database.NewManager(config)

		ctx := cmd.Context()
		if err := manager.Initialize(ctx); err != nil {
			return fmt.Errorf("failed to initialize database: %w", err)
		}
		defer manager.Close()

		sessionID := args[0]
		response := args[1]

		if err := database.HandleClaudeResponseHook(sessionID, response, nil, nil); err != nil {
			return fmt.Errorf("failed to handle claude response: %w", err)
		}

		fmt.Printf("Claude response captured successfully\n")
		return nil
	},
}

var sessionEndCmd = &cobra.Command{
	Use:   "session-end [session-id]",
	Short: "Capture session end event",
	Long: `Capture a session end event and store it in the database.

Session ID is required. An optional summary can be provided.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		config := database.DefaultDatabaseConfig()
		manager := database.NewManager(config)

		ctx := cmd.Context()
		if err := manager.Initialize(ctx); err != nil {
			return fmt.Errorf("failed to initialize database: %w", err)
		}
		defer manager.Close()

		sessionID := args[0]
		summary, _ := cmd.Flags().GetString("summary")

		if err := database.HandleSessionEndHook(sessionID, summary); err != nil {
			return fmt.Errorf("failed to handle session end: %w", err)
		}

		fmt.Printf("Session ended successfully\n")
		return nil
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show database status",
	Long: `Show the current status of the database including path, size, and basic statistics.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		config := database.DefaultDatabaseConfig()

		// Check if database file exists
		if _, err := os.Stat(config.DatabasePath); os.IsNotExist(err) {
			fmt.Printf("Database not initialized. Run 'context-extender database init' to create it.\n")
			fmt.Printf("Would be created at: %s\n", config.DatabasePath)
			return nil
		}

		// Create database manager
		manager := database.NewManager(config)
		ctx := cmd.Context()

		// Initialize database
		if err := manager.Initialize(ctx); err != nil {
			return fmt.Errorf("failed to connect to database: %w", err)
		}
		defer manager.Close()

		// Get backend
		backend, err := manager.GetBackend()
		if err != nil {
			return fmt.Errorf("failed to get backend: %w", err)
		}

		// Get backend info
		backendInfo := backend.GetBackendInfo()

		// Get file info
		fileInfo, err := os.Stat(config.DatabasePath)
		if err != nil {
			return fmt.Errorf("failed to get database file info: %w", err)
		}

		fmt.Printf("Database Status:\n")
		fmt.Printf("  Path: %s\n", config.DatabasePath)
		fmt.Printf("  Size: %d bytes\n", fileInfo.Size())
		fmt.Printf("  Modified: %s\n", fileInfo.ModTime().Format("2006-01-02 15:04:05"))
		fmt.Printf("  Backend: %s\n", backendInfo.Name)
		fmt.Printf("  Version: %s\n", backendInfo.Version)
		fmt.Printf("  CGO Required: %v\n", backendInfo.RequiresCGO)

		// Test connection
		if err := backend.Ping(ctx); err != nil {
			return fmt.Errorf("database connection failed: %w", err)
		}

		fmt.Printf("  Connection: ✓ Active\n")

		// Get database stats
		stats, err := backend.GetDatabaseStats(ctx)
		if err != nil {
			return fmt.Errorf("failed to get database stats: %w", err)
		}

		fmt.Printf("\nStatistics:\n")
		fmt.Printf("  Sessions: %d\n", stats.SessionCount)
		fmt.Printf("  Events: %d\n", stats.EventCount)
		fmt.Printf("  Conversations: %d\n", stats.ConversationCount)

		if stats.OldestRecord != nil {
			fmt.Printf("  Oldest Record: %s\n", stats.OldestRecord.Format("2006-01-02 15:04:05"))
		}
		if stats.NewestRecord != nil {
			fmt.Printf("  Newest Record: %s\n", stats.NewestRecord.Format("2006-01-02 15:04:05"))
		}

		return nil
	},
}

func init() {
	// Add session-end summary flag
	sessionEndCmd.Flags().StringP("summary", "s", "", "Optional session summary")

	// Build command hierarchy
	captureCmd.AddCommand(sessionStartCmd)
	captureCmd.AddCommand(userPromptCmd)
	captureCmd.AddCommand(claudeResponseCmd)
	captureCmd.AddCommand(sessionEndCmd)

	databaseCmd.AddCommand(initDbCmd)
	databaseCmd.AddCommand(migrateDbCmd)
	databaseCmd.AddCommand(captureCmd)
	databaseCmd.AddCommand(statusCmd)

	rootCmd.AddCommand(databaseCmd)
}