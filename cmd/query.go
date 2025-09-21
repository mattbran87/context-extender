package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"context-extender/internal/database"
)

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query completed conversations",
	Long:  "Search and retrieve completed conversation data.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var queryListCmd = &cobra.Command{
	Use:   "list",
	Short: "List completed conversations",
	Long:  "Display a list of all completed conversations with basic metadata.",
	Run: func(cmd *cobra.Command, args []string) {
		handleListConversations(cmd, args)
	},
}

var showCmd = &cobra.Command{
	Use:   "show [session-id]",
	Short: "Show detailed conversation",
	Long:  "Display the full details of a specific conversation.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		handleShowConversation(cmd, args)
	},
}

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search conversations",
	Long:  "Search for conversations containing specific terms or matching criteria.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		handleSearchConversations(cmd, args)
	},
}

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show conversation statistics",
	Long:  "Display aggregate statistics across all completed conversations.",
	Run: func(cmd *cobra.Command, args []string) {
		handleConversationStats(cmd, args)
	},
}

// Query command flags
var (
	queryFormat     string
	queryLimit      int
	queryProject    string
	queryDateFrom   string
	queryDateTo     string
	queryShowEvents bool
	queryShowSummary bool
)

func init() {
	rootCmd.AddCommand(queryCmd)
	queryCmd.AddCommand(queryListCmd)
	queryCmd.AddCommand(showCmd)
	queryCmd.AddCommand(searchCmd)
	queryCmd.AddCommand(statsCmd)

	// Global query flags
	queryCmd.PersistentFlags().StringVar(&queryFormat, "format", "table", "Output format (table, json)")
	queryCmd.PersistentFlags().IntVar(&queryLimit, "limit", 10, "Limit number of results")
	queryCmd.PersistentFlags().StringVar(&queryProject, "project", "", "Filter by project name")
	queryCmd.PersistentFlags().StringVar(&queryDateFrom, "from", "", "Filter from date (YYYY-MM-DD)")
	queryCmd.PersistentFlags().StringVar(&queryDateTo, "to", "", "Filter to date (YYYY-MM-DD)")

	// Show command specific flags
	showCmd.Flags().BoolVar(&queryShowEvents, "events", false, "Show all conversation events")
	showCmd.Flags().BoolVar(&queryShowSummary, "summary", true, "Show conversation summary")
}

func handleListConversations(cmd *cobra.Command, args []string) {
	// Initialize database manager
	config := database.DefaultDatabaseConfig()
	manager := database.NewManager(config)

	ctx := context.Background()
	if err := manager.Initialize(ctx); err != nil {
		fmt.Printf("❌ Failed to initialize database: %v\n", err)
		os.Exit(1)
	}
	defer manager.Close()

	backend, err := manager.GetBackend()
	if err != nil {
		fmt.Printf("❌ Failed to get database backend: %v\n", err)
		os.Exit(1)
	}

	// Get all sessions
	sessions, err := backend.ListSessions(ctx, &database.SessionFilters{})
	if err != nil {
		fmt.Printf("❌ Failed to list sessions: %v\n", err)
		os.Exit(1)
	}

	// Convert to conversation metadata format
	conversations := convertSessionsToMetadata(ctx, backend, sessions)

	// Apply filters
	conversations = applyFilters(conversations)

	// Apply limit
	if queryLimit > 0 && len(conversations) > queryLimit {
		conversations = conversations[:queryLimit]
	}

	// Output results
	if queryFormat == "json" {
		outputJSON(conversations)
	} else {
		outputTable(conversations)
	}
}

func handleShowConversation(cmd *cobra.Command, args []string) {
	sessionID := args[0]

	// Initialize database manager
	config := database.DefaultDatabaseConfig()
	manager := database.NewManager(config)

	ctx := context.Background()
	if err := manager.Initialize(ctx); err != nil {
		fmt.Printf("❌ Failed to initialize database: %v\n", err)
		os.Exit(1)
	}
	defer manager.Close()

	backend, err := manager.GetBackend()
	if err != nil {
		fmt.Printf("❌ Failed to get database backend: %v\n", err)
		os.Exit(1)
	}

	// Get session details
	session, err := backend.GetSession(ctx, sessionID)
	if err != nil {
		fmt.Printf("❌ Failed to get session %s: %v\n", sessionID, err)
		os.Exit(1)
	}

	// Get conversations for this session
	conversations, err := backend.GetConversationsBySession(ctx, sessionID)
	if err != nil {
		fmt.Printf("❌ Failed to get conversations for session %s: %v\n", sessionID, err)
		os.Exit(1)
	}

	// Get events for this session
	events, err := backend.GetEventsBySession(ctx, sessionID)
	if err != nil {
		fmt.Printf("❌ Failed to get events for session %s: %v\n", sessionID, err)
		os.Exit(1)
	}

	// Output conversation details
	if queryFormat == "json" {
		sessionData := map[string]interface{}{
			"session":       session,
			"conversations": conversations,
			"events":        events,
		}
		data, err := json.MarshalIndent(sessionData, "", "  ")
		if err != nil {
			fmt.Printf("❌ Failed to marshal conversation: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(string(data))
	} else {
		outputSessionDetails(session, conversations, events)
	}
}

func handleSearchConversations(cmd *cobra.Command, args []string) {
	searchQuery := strings.ToLower(args[0])

	// Initialize database manager
	config := database.DefaultDatabaseConfig()
	manager := database.NewManager(config)

	ctx := context.Background()
	if err := manager.Initialize(ctx); err != nil {
		fmt.Printf("❌ Failed to initialize database: %v\n", err)
		os.Exit(1)
	}
	defer manager.Close()

	backend, err := manager.GetBackend()
	if err != nil {
		fmt.Printf("❌ Failed to get database backend: %v\n", err)
		os.Exit(1)
	}

	// Search conversations using database search
	conversations, err := backend.SearchConversations(ctx, searchQuery, queryLimit)
	if err != nil {
		fmt.Printf("❌ Failed to search conversations: %v\n", err)
		os.Exit(1)
	}

	// Get unique sessions from search results
	sessionIDs := make(map[string]bool)
	for _, conv := range conversations {
		sessionIDs[conv.SessionID] = true
	}

	// Get session details for matching conversations
	var matchingSessions []*database.Session
	for sessionID := range sessionIDs {
		session, err := backend.GetSession(ctx, sessionID)
		if err == nil {
			matchingSessions = append(matchingSessions, session)
		}
	}

	// Convert to conversation metadata
	matchingConversations := convertSessionsToMetadata(ctx, backend, matchingSessions)

	// Apply filters
	matchingConversations = applyFilters(matchingConversations)

	fmt.Printf("🔍 Found %d conversations matching '%s'\n\n", len(matchingConversations), args[0])

	// Output results
	if queryFormat == "json" {
		outputJSON(matchingConversations)
	} else {
		outputTable(matchingConversations)
	}
}

func handleConversationStats(cmd *cobra.Command, args []string) {
	// Initialize database manager
	config := database.DefaultDatabaseConfig()
	manager := database.NewManager(config)

	ctx := context.Background()
	if err := manager.Initialize(ctx); err != nil {
		fmt.Printf("❌ Failed to initialize database: %v\n", err)
		os.Exit(1)
	}
	defer manager.Close()

	backend, err := manager.GetBackend()
	if err != nil {
		fmt.Printf("❌ Failed to get database backend: %v\n", err)
		os.Exit(1)
	}

	// Get database statistics
	dbStats, err := backend.GetDatabaseStats(ctx)
	if err != nil {
		fmt.Printf("❌ Failed to get database statistics: %v\n", err)
		os.Exit(1)
	}

	// Get all sessions for detailed stats
	sessions, err := backend.ListSessions(ctx, &database.SessionFilters{})
	if err != nil {
		fmt.Printf("❌ Failed to list sessions: %v\n", err)
		os.Exit(1)
	}

	// Convert to conversation metadata and apply filters
	conversations := convertSessionsToMetadata(ctx, backend, sessions)
	conversations = applyFilters(conversations)

	// Calculate aggregate statistics
	stats := calculateAggregateStats(conversations, dbStats)

	// Output statistics
	if queryFormat == "json" {
		data, err := json.MarshalIndent(stats, "", "  ")
		if err != nil {
			fmt.Printf("❌ Failed to marshal stats: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(string(data))
	} else {
		outputStatsTable(stats)
	}
}

// Helper functions

func applyFilters(conversations []ConversationMetadata) []ConversationMetadata {
	var filtered []ConversationMetadata

	for _, conv := range conversations {
		// Project filter
		if queryProject != "" && !strings.Contains(strings.ToLower(conv.ProjectName), strings.ToLower(queryProject)) {
			continue
		}

		// Date filters
		if queryDateFrom != "" {
			fromDate, err := time.Parse("2006-01-02", queryDateFrom)
			if err == nil && conv.StartTime.Before(fromDate) {
				continue
			}
		}

		if queryDateTo != "" {
			toDate, err := time.Parse("2006-01-02", queryDateTo)
			if err == nil && conv.StartTime.After(toDate.Add(24*time.Hour)) {
				continue
			}
		}

		filtered = append(filtered, conv)
	}

	return filtered
}

// matchesSearch function removed - now using database backend SearchConversations

func outputJSON(data interface{}) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("❌ Failed to marshal JSON: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(jsonData))
}

func outputTable(conversations []ConversationMetadata) {
	if len(conversations) == 0 {
		fmt.Println("📭 No conversations found")
		return
	}

	fmt.Printf("📊 Found %d conversations:\n\n", len(conversations))
	fmt.Printf("%-12s %-20s %-15s %-8s %-10s %-25s\n",
		"SESSION", "PROJECT", "STATUS", "EVENTS", "DURATION", "START TIME")
	fmt.Printf("%-12s %-20s %-15s %-8s %-10s %-25s\n",
		"━━━━━━━", "━━━━━━━━━━━━━━━━━━━", "━━━━━━━━━━━━━━", "━━━━━━━", "━━━━━━━━━", "━━━━━━━━━━━━━━━━━━━━━━━━")

	for _, conv := range conversations {
		sessionShort := conv.SessionID[:8] + "..."
		projectName := conv.ProjectName
		if len(projectName) > 18 {
			projectName = projectName[:15] + "..."
		}

		startTime := conv.StartTime.Format("2006-01-02 15:04:05")

		fmt.Printf("%-12s %-20s %-15s %-8d %-10s %-25s\n",
			sessionShort, projectName, conv.Status, conv.EventCount, conv.Duration, startTime)
	}
}

func outputSessionDetails(session *database.Session, conversations []*database.Conversation, events []*database.Event) {
	fmt.Printf("🗣️  Session Details\n")
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")

	// Parse metadata for additional info
	var metadata map[string]interface{}
	projectName := ""
	workingDir := ""
	if session.Metadata != "" {
		if json.Unmarshal([]byte(session.Metadata), &metadata) == nil {
			if wd, ok := metadata["working_directory"].(string); ok {
				workingDir = wd
			}
			if proj, ok := metadata["project"].(string); ok {
				projectName = proj
			}
		}
	}

	// Session Information
	fmt.Printf("📋 Session Information\n")
	fmt.Printf("   Session ID:     %s\n", session.ID)
	fmt.Printf("   Project:        %s\n", projectName)
	fmt.Printf("   Working Dir:    %s\n", workingDir)
	fmt.Printf("   Status:         %s\n", session.Status)
	fmt.Printf("   Duration:       %s\n", session.UpdatedAt.Sub(session.CreatedAt).String())
	fmt.Printf("   Start Time:     %s\n", session.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("   End Time:       %s\n", session.UpdatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("\n")

	// Summary
	if queryShowSummary {
		userPrompts := 0
		claudeReplies := 0
		for _, conv := range conversations {
			switch conv.MessageType {
			case "user":
				userPrompts++
			case "assistant":
				claudeReplies++
			}
		}

		fmt.Printf("📊 Summary\n")
		fmt.Printf("   Total Events:    %d\n", len(events))
		fmt.Printf("   User Prompts:    %d\n", userPrompts)
		fmt.Printf("   Claude Replies:  %d\n", claudeReplies)
		fmt.Printf("   Total Conversations: %d\n", len(conversations))
		fmt.Printf("\n")
	}

	// Events
	if queryShowEvents {
		fmt.Printf("💬 Session Events\n")
		for i, event := range events {
			fmt.Printf("   %d. [%s] %s\n", i+1, event.Timestamp.Format("15:04:05"), event.EventType)
			if event.Data != "" && len(event.Data) > 0 {
				preview := event.Data
				if len(preview) > 100 {
					preview = preview[:97] + "..."
				}
				fmt.Printf("      → %s\n", preview)
			}
		}
		fmt.Printf("\n")
	}

	// Conversations
	if len(conversations) > 0 {
		fmt.Printf("💬 Conversations\n")
		for i, conv := range conversations {
			fmt.Printf("   %d. [%s] %s\n", i+1, conv.Timestamp.Format("15:04:05"), conv.MessageType)
			content := conv.Content
			if len(content) > 100 {
				content = content[:97] + "..."
			}
			fmt.Printf("      → %s\n", content)
		}
		fmt.Printf("\n")
	}
}

// ConversationMetadata represents metadata about a conversation (compatible with converter package)
type ConversationMetadata struct {
	SessionID     string    `json:"session_id"`
	ProjectName   string    `json:"project_name,omitempty"`
	WorkingDir    string    `json:"working_dir"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	Duration      string    `json:"duration"`
	Status        string    `json:"status"`
	EventCount    int       `json:"event_count"`
	UserPrompts   int       `json:"user_prompts"`
	ClaudeReplies int       `json:"claude_replies"`
}

// convertSessionsToMetadata converts database sessions to conversation metadata
func convertSessionsToMetadata(ctx context.Context, backend database.DatabaseBackend, sessions []*database.Session) []ConversationMetadata {
	var conversations []ConversationMetadata

	for _, session := range sessions {
		// Get conversations for this session to count prompts/replies
		convs, err := backend.GetConversationsBySession(ctx, session.ID)
		if err != nil {
			continue // Skip this session if we can't get conversations
		}

		// Count user prompts and Claude replies
		userPrompts := 0
		claudeReplies := 0
		for _, conv := range convs {
			switch conv.MessageType {
			case "user":
				userPrompts++
			case "assistant":
				claudeReplies++
			}
		}

		// Get events to count them
		events, err := backend.GetEventsBySession(ctx, session.ID)
		eventCount := 0
		if err == nil {
			eventCount = len(events)
		}

		// Parse metadata to get project and working directory
		projectName := ""
		workingDir := ""
		if session.Metadata != "" {
			// Try to parse JSON metadata
			var metadata map[string]interface{}
			if json.Unmarshal([]byte(session.Metadata), &metadata) == nil {
				if wd, ok := metadata["working_directory"].(string); ok {
					workingDir = wd
				}
				if proj, ok := metadata["project"].(string); ok {
					projectName = proj
				}
			}
		}

		// Calculate duration
		duration := session.UpdatedAt.Sub(session.CreatedAt).String()

		conversations = append(conversations, ConversationMetadata{
			SessionID:     session.ID,
			ProjectName:   projectName,
			WorkingDir:    workingDir,
			StartTime:     session.CreatedAt,
			EndTime:       session.UpdatedAt,
			Duration:      duration,
			Status:        session.Status,
			EventCount:    eventCount,
			UserPrompts:   userPrompts,
			ClaudeReplies: claudeReplies,
		})
	}

	// Sort by start time (newest first)
	sort.Slice(conversations, func(i, j int) bool {
		return conversations[i].StartTime.After(conversations[j].StartTime)
	})

	return conversations
}

// AggregateStats represents statistics across all conversations
type AggregateStats struct {
	TotalConversations int                                    `json:"total_conversations"`
	TotalEvents        int                                    `json:"total_events"`
	TotalPrompts       int                                    `json:"total_prompts"`
	TotalReplies       int                                    `json:"total_replies"`
	TotalDuration      time.Duration                          `json:"total_duration"`
	AverageDuration    time.Duration                          `json:"average_duration"`
	ProjectBreakdown   map[string]int                         `json:"project_breakdown"`
	StatusBreakdown    map[string]int                         `json:"status_breakdown"`
	DateRange          struct {
		Earliest time.Time `json:"earliest"`
		Latest   time.Time `json:"latest"`
	} `json:"date_range"`
}

func calculateAggregateStats(conversations []ConversationMetadata, dbStats *database.DatabaseStats) *AggregateStats {
	stats := &AggregateStats{
		ProjectBreakdown: make(map[string]int),
		StatusBreakdown:  make(map[string]int),
	}

	if len(conversations) == 0 {
		return stats
	}

	stats.TotalConversations = len(conversations)

	var totalDurationNanos int64
	earliest := conversations[0].StartTime
	latest := conversations[0].StartTime

	for _, conv := range conversations {
		// Count events and prompts
		stats.TotalEvents += conv.EventCount
		stats.TotalPrompts += conv.UserPrompts
		stats.TotalReplies += conv.ClaudeReplies

		// Duration calculation (parse duration string)
		if duration, err := time.ParseDuration(conv.Duration); err == nil {
			totalDurationNanos += duration.Nanoseconds()
		}

		// Project breakdown
		project := conv.ProjectName
		if project == "" {
			project = "(unnamed)"
		}
		stats.ProjectBreakdown[project]++

		// Status breakdown
		stats.StatusBreakdown[conv.Status]++

		// Date range
		if conv.StartTime.Before(earliest) {
			earliest = conv.StartTime
		}
		if conv.StartTime.After(latest) {
			latest = conv.StartTime
		}
	}

	stats.TotalDuration = time.Duration(totalDurationNanos)
	stats.AverageDuration = time.Duration(totalDurationNanos / int64(len(conversations)))
	stats.DateRange.Earliest = earliest
	stats.DateRange.Latest = latest

	return stats
}

func outputStatsTable(stats *AggregateStats) {
	fmt.Printf("📈 Conversation Statistics\n")
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")

	fmt.Printf("📊 Overall Statistics\n")
	fmt.Printf("   Total Conversations:  %d\n", stats.TotalConversations)
	fmt.Printf("   Total Events:         %d\n", stats.TotalEvents)
	fmt.Printf("   Total Prompts:        %d\n", stats.TotalPrompts)
	fmt.Printf("   Total Replies:        %d\n", stats.TotalReplies)
	fmt.Printf("   Total Duration:       %s\n", stats.TotalDuration.String())
	fmt.Printf("   Average Duration:     %s\n", stats.AverageDuration.String())
	fmt.Printf("\n")

	fmt.Printf("📅 Date Range\n")
	fmt.Printf("   Earliest:  %s\n", stats.DateRange.Earliest.Format("2006-01-02 15:04:05"))
	fmt.Printf("   Latest:    %s\n", stats.DateRange.Latest.Format("2006-01-02 15:04:05"))
	fmt.Printf("\n")

	if len(stats.ProjectBreakdown) > 0 {
		fmt.Printf("🚀 Projects\n")

		// Sort projects by count
		type projectCount struct {
			name  string
			count int
		}
		var projects []projectCount
		for name, count := range stats.ProjectBreakdown {
			projects = append(projects, projectCount{name, count})
		}
		sort.Slice(projects, func(i, j int) bool {
			return projects[i].count > projects[j].count
		})

		for _, p := range projects {
			fmt.Printf("   %-20s  %d conversations\n", p.name, p.count)
		}
		fmt.Printf("\n")
	}

	if len(stats.StatusBreakdown) > 0 {
		fmt.Printf("✅ Status Breakdown\n")
		for status, count := range stats.StatusBreakdown {
			fmt.Printf("   %-15s  %d conversations\n", status, count)
		}
	}
}