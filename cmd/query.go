package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"context-extender/internal/converter"
	"context-extender/internal/storage"
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
	// Initialize storage manager
	storageManager, err := storage.NewStorageManager(nil)
	if err != nil {
		fmt.Printf("❌ Failed to initialize storage: %v\n", err)
		os.Exit(1)
	}

	// Initialize converter
	sessionConverter := converter.NewSessionConverter(storageManager.GetConversationsDir())

	// Get all completed conversations
	conversations, err := sessionConverter.ListCompletedConversations()
	if err != nil {
		fmt.Printf("❌ Failed to list conversations: %v\n", err)
		os.Exit(1)
	}

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

	// Initialize storage manager
	storageManager, err := storage.NewStorageManager(nil)
	if err != nil {
		fmt.Printf("❌ Failed to initialize storage: %v\n", err)
		os.Exit(1)
	}

	// Initialize converter
	sessionConverter := converter.NewSessionConverter(storageManager.GetConversationsDir())

	// Load the conversation
	conversation, err := sessionConverter.LoadCompletedConversation(sessionID)
	if err != nil {
		fmt.Printf("❌ Failed to load conversation %s: %v\n", sessionID, err)
		os.Exit(1)
	}

	// Output conversation details
	if queryFormat == "json" {
		data, err := json.MarshalIndent(conversation, "", "  ")
		if err != nil {
			fmt.Printf("❌ Failed to marshal conversation: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(string(data))
	} else {
		outputConversationDetails(conversation)
	}
}

func handleSearchConversations(cmd *cobra.Command, args []string) {
	searchQuery := strings.ToLower(args[0])

	// Initialize storage manager
	storageManager, err := storage.NewStorageManager(nil)
	if err != nil {
		fmt.Printf("❌ Failed to initialize storage: %v\n", err)
		os.Exit(1)
	}

	// Initialize converter
	sessionConverter := converter.NewSessionConverter(storageManager.GetConversationsDir())

	// Get all completed conversations
	conversations, err := sessionConverter.ListCompletedConversations()
	if err != nil {
		fmt.Printf("❌ Failed to list conversations: %v\n", err)
		os.Exit(1)
	}

	// Filter conversations that match search query
	var matchingConversations []converter.ConversationMetadata
	for _, conv := range conversations {
		if matchesSearch(conv, searchQuery, sessionConverter) {
			matchingConversations = append(matchingConversations, conv)
		}
	}

	// Apply other filters
	matchingConversations = applyFilters(matchingConversations)

	// Apply limit
	if queryLimit > 0 && len(matchingConversations) > queryLimit {
		matchingConversations = matchingConversations[:queryLimit]
	}

	fmt.Printf("🔍 Found %d conversations matching '%s'\n\n", len(matchingConversations), args[0])

	// Output results
	if queryFormat == "json" {
		outputJSON(matchingConversations)
	} else {
		outputTable(matchingConversations)
	}
}

func handleConversationStats(cmd *cobra.Command, args []string) {
	// Initialize storage manager
	storageManager, err := storage.NewStorageManager(nil)
	if err != nil {
		fmt.Printf("❌ Failed to initialize storage: %v\n", err)
		os.Exit(1)
	}

	// Initialize converter
	sessionConverter := converter.NewSessionConverter(storageManager.GetConversationsDir())

	// Get all completed conversations
	conversations, err := sessionConverter.ListCompletedConversations()
	if err != nil {
		fmt.Printf("❌ Failed to list conversations: %v\n", err)
		os.Exit(1)
	}

	// Apply filters
	conversations = applyFilters(conversations)

	// Calculate aggregate statistics
	stats := calculateAggregateStats(conversations)

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

func applyFilters(conversations []converter.ConversationMetadata) []converter.ConversationMetadata {
	var filtered []converter.ConversationMetadata

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

func matchesSearch(conv converter.ConversationMetadata, query string, sessionConverter *converter.SessionConverter) bool {
	// Check project name
	if strings.Contains(strings.ToLower(conv.ProjectName), query) {
		return true
	}

	// Check working directory
	if strings.Contains(strings.ToLower(conv.WorkingDir), query) {
		return true
	}

	// Check session ID
	if strings.Contains(strings.ToLower(conv.SessionID), query) {
		return true
	}

	// TODO: Search within conversation content (requires loading full conversation)
	return false
}

func outputJSON(data interface{}) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("❌ Failed to marshal JSON: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(jsonData))
}

func outputTable(conversations []converter.ConversationMetadata) {
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

func outputConversationDetails(conv *converter.CompletedConversation) {
	fmt.Printf("🗣️  Conversation Details\n")
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")

	// Metadata
	fmt.Printf("📋 Session Information\n")
	fmt.Printf("   Session ID:     %s\n", conv.Metadata.SessionID)
	fmt.Printf("   Project:        %s\n", conv.Metadata.ProjectName)
	fmt.Printf("   Working Dir:    %s\n", conv.Metadata.WorkingDir)
	fmt.Printf("   Status:         %s\n", conv.Metadata.Status)
	fmt.Printf("   Duration:       %s\n", conv.Metadata.Duration)
	fmt.Printf("   Start Time:     %s\n", conv.Metadata.StartTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("   End Time:       %s\n", conv.Metadata.EndTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("\n")

	// Summary
	if queryShowSummary {
		fmt.Printf("📊 Summary\n")
		fmt.Printf("   Total Events:    %d\n", conv.Summary.Statistics.TotalEvents)
		fmt.Printf("   User Prompts:    %d\n", conv.Summary.PromptCount)
		fmt.Printf("   Claude Replies:  %d\n", conv.Summary.ResponseCount)
		fmt.Printf("   Avg Gap Time:    %s\n", conv.Summary.AverageGapTime)
		fmt.Printf("   User Words:      %d\n", conv.Summary.Statistics.UserPromptWords)
		fmt.Printf("   Claude Words:    %d\n", conv.Summary.Statistics.ClaudeResponseWords)
		fmt.Printf("\n")
	}

	// Events
	if queryShowEvents {
		fmt.Printf("💬 Conversation Events\n")
		for i, event := range conv.Conversation {
			fmt.Printf("   %d. [%s] %s\n", i+1, event.Timestamp.Format("15:04:05"), event.EventType)

			switch event.EventType {
			case "session-start":
				if event.Content.SessionInfo != nil {
					fmt.Printf("      → Started session for project: %s\n", event.Content.SessionInfo.Project)
				}
			case "user-prompt":
				if event.Content.UserPrompt != nil {
					message := event.Content.UserPrompt.Message
					if len(message) > 100 {
						message = message[:97] + "..."
					}
					fmt.Printf("      → User: %s\n", message)
				}
			case "claude-response":
				if event.Content.ClaudeResponse != nil {
					response := event.Content.ClaudeResponse.Response
					if len(response) > 100 {
						response = response[:97] + "..."
					}
					fmt.Printf("      → Claude: %s\n", response)
				}
			case "session-end":
				fmt.Printf("      → Session ended\n")
			}
		}
		fmt.Printf("\n")
	}

	// Export info
	fmt.Printf("🔄 Export Information\n")
	fmt.Printf("   Exported:       %s\n", conv.Export.ExportTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("   Version:        %s\n", conv.Export.ExportVersion)
	fmt.Printf("   Source:         %s\n", conv.Export.SourceFormat)
	fmt.Printf("   Processed by:   %s\n", conv.Export.ProcessedBy)
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

func calculateAggregateStats(conversations []converter.ConversationMetadata) *AggregateStats {
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