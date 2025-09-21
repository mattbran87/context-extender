package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	golang_context "context"
	preserver "context-extender/internal/context"
	"context-extender/internal/database"
	"github.com/spf13/cobra"
)

var contextCmd = &cobra.Command{
	Use:   "context",
	Short: "Manage conversation context preservation",
	Long: `Manage and analyze conversation context for preservation during compression.

This command helps preserve critical context when Claude conversations are compressed,
ensuring important decisions, preferences, and technical details are not lost.`,
}

var contextAnalyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze current session for critical context",
	Long:  `Analyze the current session's conversations to extract critical context that should be preserved.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		sessionID, _ := cmd.Flags().GetString("session")
		// limit, _ := cmd.Flags().GetInt("limit") // TODO: Implement limiting

		// Initialize database
		config := database.DefaultDatabaseConfig()
		manager := database.NewManager(config)

		ctx := golang_context.Background()
		if err := manager.Initialize(ctx); err != nil {
			return fmt.Errorf("failed to initialize database: %w", err)
		}
		defer manager.Close()

		backend, err := manager.GetBackend()
		if err != nil {
			return fmt.Errorf("failed to get backend: %w", err)
		}

		// Get conversations for analysis
		conversations, err := backend.GetConversationsBySession(ctx, sessionID)
		if err != nil {
			return fmt.Errorf("failed to get conversations: %w", err)
		}

		// Extract content for analysis
		var contents []string
		for _, conv := range conversations {
			contents = append(contents, conv.Content)
		}

		// Analyze for critical context
		summary, err := preserver.ExtractCriticalContext(contents)
		if err != nil {
			return fmt.Errorf("failed to extract context: %w", err)
		}

		// Display analysis
		fmt.Println("🔍 Context Analysis Results")
		fmt.Println("=" + strings.Repeat("=", 50))

		jsonOutput, _ := cmd.Flags().GetBool("json")
		if jsonOutput {
			data, _ := json.MarshalIndent(summary, "", "  ")
			fmt.Println(string(data))
		} else {
			fmt.Printf("\n📦 Project: %s\n", summary.ProjectName)
			fmt.Printf("🎯 Current Objective: %s\n", summary.CurrentObjective)
			fmt.Printf("📍 Phase: %s\n", summary.CurrentPhase)

			if len(summary.TechnicalStack) > 0 {
				fmt.Println("\n🔧 Technical Stack:")
				for _, tech := range summary.TechnicalStack {
					fmt.Printf("  - %s\n", tech)
				}
			}

			if len(summary.Constraints) > 0 {
				fmt.Println("\n⚠️ Constraints:")
				for _, constraint := range summary.Constraints {
					fmt.Printf("  - %s\n", constraint)
				}
			}

			if len(summary.KeyDecisions) > 0 {
				fmt.Println("\n📌 Key Decisions:")
				for key, value := range summary.KeyDecisions {
					fmt.Printf("  - %s: %s\n", key, value)
				}
			}

			if len(summary.UserPreferences) > 0 {
				fmt.Println("\n👤 User Preferences:")
				for key, value := range summary.UserPreferences {
					fmt.Printf("  - %s: %s\n", key, value)
				}
			}

			if len(summary.CompletedTasks) > 0 {
				fmt.Println("\n✅ Completed Tasks:")
				for _, task := range summary.CompletedTasks {
					fmt.Printf("  - %s\n", task)
				}
			}

			if len(summary.ErrorsToAvoid) > 0 {
				fmt.Println("\n❌ Errors to Avoid:")
				for _, err := range summary.ErrorsToAvoid {
					fmt.Printf("  - %s\n", err)
				}
			}
		}

		return nil
	},
}

var contextPreserveCmd = &cobra.Command{
	Use:   "preserve",
	Short: "Preserve current context for later reinjection",
	Long:  `Analyze and preserve the current session's critical context for reinjection after compression.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		sessionID, _ := cmd.Flags().GetString("session")
		// limit, _ := cmd.Flags().GetInt("limit") // TODO: Implement limiting

		// Initialize database
		config := database.DefaultDatabaseConfig()
		manager := database.NewManager(config)

		ctx := golang_context.Background()
		if err := manager.Initialize(ctx); err != nil {
			return fmt.Errorf("failed to initialize database: %w", err)
		}
		defer manager.Close()

		backend, err := manager.GetBackend()
		if err != nil {
			return fmt.Errorf("failed to get backend: %w", err)
		}

		// Get conversations for analysis
		conversations, err := backend.GetConversationsBySession(ctx, sessionID)
		if err != nil {
			return fmt.Errorf("failed to get conversations: %w", err)
		}

		// Extract content for analysis
		var contents []string
		for _, conv := range conversations {
			contents = append(contents, conv.Content)
		}

		// Analyze for critical context
		summary, err := preserver.ExtractCriticalContext(contents)
		if err != nil {
			return fmt.Errorf("failed to extract context: %w", err)
		}

		// Save as compression event
		contextJSON, err := preserver.SaveCompressionContext(summary)
		if err != nil {
			return fmt.Errorf("failed to save context: %w", err)
		}

		// Store in database
		event := &database.Event{
			ID:        fmt.Sprintf("%s_preserve_%d", sessionID, 0),
			SessionID: sessionID,
			EventType: "context_preservation",
			Data:      contextJSON,
		}

		if err := backend.CreateEvent(ctx, event); err != nil {
			return fmt.Errorf("failed to store preservation event: %w", err)
		}

		fmt.Println("✅ Context preserved successfully!")
		fmt.Printf("Session: %s\n", sessionID)
		fmt.Printf("Preserved: %d key decisions, %d constraints, %d preferences\n",
			len(summary.KeyDecisions), len(summary.Constraints), len(summary.UserPreferences))

		return nil
	},
}

var contextReinjectCmd = &cobra.Command{
	Use:   "reinject",
	Short: "Generate context reinjection prompt",
	Long:  `Generate a prompt to reinject preserved context after a conversation compression.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		sessionID, _ := cmd.Flags().GetString("session")

		// Initialize database
		config := database.DefaultDatabaseConfig()
		manager := database.NewManager(config)

		ctx := golang_context.Background()
		if err := manager.Initialize(ctx); err != nil {
			return fmt.Errorf("failed to initialize database: %w", err)
		}
		defer manager.Close()

		backend, err := manager.GetBackend()
		if err != nil {
			return fmt.Errorf("failed to get backend: %w", err)
		}

		// Get all events for session
		events, err := backend.GetEventsBySession(ctx, sessionID)
		if err != nil {
			return fmt.Errorf("failed to get events: %w", err)
		}

		// Find latest preservation event
		var latestPreservation *database.Event
		for _, event := range events {
			if event.EventType == "context_preservation" {
				latestPreservation = event
			}
		}

		if latestPreservation == nil {
			return fmt.Errorf("no preserved context found for session %s", sessionID)
		}

		// Load preserved context
		summary, err := preserver.LoadCompressionContext(latestPreservation.Data)
		if err != nil {
			return fmt.Errorf("failed to load context: %w", err)
		}

		// Generate reinjection prompt
		prompt := preserver.GenerateContextPrompt(summary)
		fmt.Println(prompt)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(contextCmd)
	contextCmd.AddCommand(contextAnalyzeCmd)
	contextCmd.AddCommand(contextPreserveCmd)
	contextCmd.AddCommand(contextReinjectCmd)

	// Common flags
	contextAnalyzeCmd.Flags().StringP("session", "s", "", "Session ID to analyze")
	contextAnalyzeCmd.Flags().IntP("limit", "l", 100, "Maximum conversations to analyze")
	contextAnalyzeCmd.Flags().Bool("json", false, "Output in JSON format")

	contextPreserveCmd.Flags().StringP("session", "s", "", "Session ID to preserve")
	contextPreserveCmd.Flags().IntP("limit", "l", 100, "Maximum conversations to analyze")

	contextReinjectCmd.Flags().StringP("session", "s", "", "Session ID to reinject context for")
}