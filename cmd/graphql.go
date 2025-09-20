package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"context-extender/internal/database"
	"context-extender/internal/graphql"
	"github.com/spf13/cobra"
)

var (
	graphqlPort      int
	graphqlPretty    bool
	graphqlVariables string
)

var graphqlCmd = &cobra.Command{
	Use:   "graphql",
	Short: "GraphQL query interface for conversation data",
	Long: `GraphQL query interface providing powerful search and analysis capabilities for conversation data.

This command provides access to sessions, conversations, events, and search functionality
through a GraphQL interface. You can run queries directly or start an interactive server.`,
}

var graphqlServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Start GraphQL server",
	Long: `Start a GraphQL server with an interactive playground.

The server provides:
- GraphQL endpoint at /graphql
- Interactive playground at /
- CORS support for web applications
- Real-time query execution

Example:
  context-extender graphql server --port 8080`,
	RunE: func(cmd *cobra.Command, args []string) error {
		server, err := graphql.NewGraphQLServer(graphqlPort)
		if err != nil {
			return fmt.Errorf("failed to create GraphQL server: %w", err)
		}

		// This will block
		return server.Start()
	},
}

var graphqlExecCmd = &cobra.Command{
	Use:   "exec [query]",
	Short: "Execute a GraphQL query",
	Long: `Execute a GraphQL query directly and return the results.

Example queries:
  # Get database statistics
  context-extender graphql exec "{ stats { totalSessions totalConversations } }"

  # List recent sessions
  context-extender graphql exec "{ sessions(limit: 5) { id createdAt status } }"

  # Search conversations
  context-extender graphql exec "{ search(query: \"hello\") { totalCount conversations { content } } }"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := args[0]

		// Initialize database and schema
		if err := initializeGraphQL(); err != nil {
			return fmt.Errorf("failed to initialize GraphQL: %w", err)
		}

		// Parse variables if provided
		var variables map[string]interface{}
		if graphqlVariables != "" {
			if err := json.Unmarshal([]byte(graphqlVariables), &variables); err != nil {
				return fmt.Errorf("invalid variables JSON: %w", err)
			}
		}

		// Execute query
		result := graphql.ExecuteQuery(query, variables)

		// Format output
		var output []byte
		var err error

		if graphqlPretty {
			output, err = json.MarshalIndent(result, "", "  ")
		} else {
			output, err = json.Marshal(result)
		}

		if err != nil {
			return fmt.Errorf("failed to marshal result: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

var graphqlExamplesCmd = &cobra.Command{
	Use:   "examples",
	Short: "Show example GraphQL queries",
	Long:  `Display a collection of useful GraphQL query examples.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		examples := []struct {
			Title       string
			Description string
			Query       string
		}{
			{
				Title:       "Database Statistics",
				Description: "Get overall database statistics",
				Query: `{
  stats {
    totalSessions
    totalConversations
    totalEvents
    totalImports
    oldestSession
    newestSession
  }
}`,
			},
			{
				Title:       "Recent Sessions",
				Description: "List the 10 most recent sessions",
				Query: `{
  sessions(limit: 10, sortBy: "created_at", sortOrder: "DESC") {
    id
    createdAt
    status
    metadata
  }
}`,
			},
			{
				Title:       "Session with Conversations",
				Description: "Get a specific session with its conversations",
				Query: `{
  session(id: "your-session-id") {
    id
    createdAt
    status
    conversations(limit: 50) {
      messageType
      content
      timestamp
    }
  }
}`,
			},
			{
				Title:       "Search Conversations",
				Description: "Search for conversations containing specific text",
				Query: `{
  search(query: "database", limit: 20) {
    totalCount
    conversations {
      sessionId
      messageType
      content
      timestamp
    }
  }
}`,
			},
			{
				Title:       "Recent Events",
				Description: "Get recent events across all sessions",
				Query: `{
  events(limit: 20, eventType: "user_prompt") {
    id
    sessionId
    eventType
    timestamp
    sequenceNumber
  }
}`,
			},
			{
				Title:       "User Messages Only",
				Description: "Get only user messages from conversations",
				Query: `{
  conversations(messageType: "user", limit: 10) {
    sessionId
    content
    timestamp
  }
}`,
			},
			{
				Title:       "Active Sessions",
				Description: "List all active sessions",
				Query: `{
  sessions(status: "active") {
    id
    createdAt
    events(limit: 5) {
      eventType
      timestamp
    }
  }
}`,
			},
		}

		fmt.Println("📚 GraphQL Query Examples")
		fmt.Println("=" + strings.Repeat("=", 50))

		for i, example := range examples {
			fmt.Printf("\n%d. %s\n", i+1, example.Title)
			fmt.Printf("   %s\n\n", example.Description)
			fmt.Printf("   Query:\n")

			// Indent the query
			lines := strings.Split(example.Query, "\n")
			for _, line := range lines {
				fmt.Printf("   %s\n", line)
			}

			fmt.Printf("\n   Execute with:\n")
			fmt.Printf("   context-extender graphql exec '%s'\n", strings.ReplaceAll(example.Query, "\n", ""))

			if i < len(examples)-1 {
				fmt.Println("\n" + strings.Repeat("-", 60))
			}
		}

		fmt.Printf("\n\n💡 Tips:\n")
		fmt.Printf("• Use --pretty flag for formatted output\n")
		fmt.Printf("• Start the server for interactive queries: context-extender graphql server\n")
		fmt.Printf("• Variables can be passed with --variables '{\"key\": \"value\"}'\n")

		return nil
	},
}

var graphqlStatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Quick database statistics",
	Long:  `Display quick database statistics using GraphQL.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Initialize GraphQL
		if err := initializeGraphQL(); err != nil {
			return fmt.Errorf("failed to initialize GraphQL: %w", err)
		}

		// Execute stats query
		query := `{
			stats {
				totalSessions
				totalConversations
				totalEvents
				totalImports
				oldestSession
				newestSession
			}
		}`

		result := graphql.ExecuteQuery(query, nil)

		if len(result.Errors) > 0 {
			return fmt.Errorf("GraphQL errors: %v", result.Errors)
		}

		// Extract stats from result
		stats, ok := result.Data.(map[string]interface{})["stats"].(map[string]interface{})
		if !ok {
			return fmt.Errorf("unexpected result format")
		}

		fmt.Println("📊 Database Statistics")
		fmt.Println("=====================")
		fmt.Printf("Sessions:      %v\n", stats["totalSessions"])
		fmt.Printf("Conversations: %v\n", stats["totalConversations"])
		fmt.Printf("Events:        %v\n", stats["totalEvents"])
		fmt.Printf("Imports:       %v\n", stats["totalImports"])

		if oldest, ok := stats["oldestSession"].(string); ok && oldest != "" {
			fmt.Printf("Oldest:        %s\n", oldest)
		}
		if newest, ok := stats["newestSession"].(string); ok && newest != "" {
			fmt.Printf("Newest:        %s\n", newest)
		}

		return nil
	},
}

var graphqlSearchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search conversations using GraphQL",
	Long: `Search across all conversations for specific text.

Example:
  context-extender graphql search "database implementation"
  context-extender graphql search "GraphQL" --limit 10`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		searchQuery := args[0]
		limit, _ := cmd.Flags().GetInt("limit")

		// Initialize GraphQL
		if err := initializeGraphQL(); err != nil {
			return fmt.Errorf("failed to initialize GraphQL: %w", err)
		}

		// Build GraphQL query
		query := fmt.Sprintf(`{
			search(query: "%s", limit: %d) {
				totalCount
				conversations {
					sessionId
					messageType
					content
					timestamp
				}
			}
		}`, searchQuery, limit)

		result := graphql.ExecuteQuery(query, nil)

		if len(result.Errors) > 0 {
			return fmt.Errorf("GraphQL errors: %v", result.Errors)
		}

		// Format and display results
		if graphqlPretty {
			output, _ := json.MarshalIndent(result.Data, "", "  ")
			fmt.Println(string(output))
		} else {
			// Extract and format search results
			searchData := result.Data.(map[string]interface{})["search"].(map[string]interface{})
			totalCount := searchData["totalCount"]

			fmt.Printf("🔍 Search Results for: %s\n", searchQuery)
			fmt.Printf("Total matches: %v\n\n", totalCount)

			conversations := searchData["conversations"].([]interface{})
			for i, conv := range conversations {
				convData := conv.(map[string]interface{})
				content := convData["content"].(string)

				// Truncate long content
				if len(content) > 200 {
					content = content[:200] + "..."
				}

				fmt.Printf("%d. [%s] %s\n", i+1, convData["messageType"], content)
				fmt.Printf("   Session: %s | Time: %s\n\n", convData["sessionId"], convData["timestamp"])
			}
		}

		return nil
	},
}

// initializeGraphQL initializes the GraphQL system
func initializeGraphQL() error {
	// Initialize database with new backend manager
	config := database.DefaultDatabaseConfig()
	manager := database.NewManager(config)

	ctx := context.Background()
	if err := manager.Initialize(ctx); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	// TEMPORARY FIX: Initialize old database system for GraphQL compatibility
	// TODO: Update GraphQL resolvers to use new backend directly
	oldConfig := &database.Config{
		DriverName:   "sqlite",
		DatabasePath: config.DatabasePath,
	}
	if err := database.Initialize(oldConfig); err != nil {
		return fmt.Errorf("failed to initialize legacy database for GraphQL: %w", err)
	}

	// Initialize GraphQL schema
	if err := graphql.InitializeSchema(); err != nil {
		return fmt.Errorf("failed to initialize GraphQL schema: %w", err)
	}

	// Setup resolvers
	graphql.SetupResolvers()

	return nil
}

func init() {
	// Server command flags
	graphqlServerCmd.Flags().IntVarP(&graphqlPort, "port", "p", 8080, "Port to run GraphQL server on")

	// Exec command flags
	graphqlExecCmd.Flags().BoolVar(&graphqlPretty, "pretty", false, "Pretty print JSON output")
	graphqlExecCmd.Flags().StringVar(&graphqlVariables, "variables", "", "JSON variables for the query")

	// Search command flags
	graphqlSearchCmd.Flags().Int("limit", 10, "Maximum number of results")
	graphqlSearchCmd.Flags().BoolVar(&graphqlPretty, "pretty", false, "Pretty print JSON output")

	// Stats command flags
	graphqlStatsCmd.Flags().BoolVar(&graphqlPretty, "pretty", false, "Pretty print JSON output")

	// Add subcommands
	graphqlCmd.AddCommand(graphqlServerCmd)
	graphqlCmd.AddCommand(graphqlExecCmd)
	graphqlCmd.AddCommand(graphqlExamplesCmd)
	graphqlCmd.AddCommand(graphqlStatsCmd)
	graphqlCmd.AddCommand(graphqlSearchCmd)

	rootCmd.AddCommand(graphqlCmd)
}