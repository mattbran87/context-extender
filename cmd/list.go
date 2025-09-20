package cmd

import (
	"github.com/spf13/cobra"
)

// listCmd provides a convenient alias for 'query list'
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List captured conversations (alias for 'query list')",
	Long: `List all captured Claude Code conversations stored in the database.
This is a convenience alias for 'query list'.

The output shows:
- Session ID
- Project/directory
- Status
- Number of events
- Duration
- Start time

Example:
  context-extender list                    # List all conversations
  context-extender list --limit 10         # List last 10 conversations
  context-extender list --project myapp    # Filter by project`,
	Run: func(cmd *cobra.Command, args []string) {
		// Call the same handler function as query list
		handleListConversations(cmd, args)
	},
}

func init() {
	// Copy flags from queryListCmd to maintain consistency
	listCmd.Flags().IntP("limit", "l", 10, "Maximum number of sessions to display")
	listCmd.Flags().StringP("project", "p", "", "Filter by project name")
	listCmd.Flags().StringP("status", "s", "", "Filter by status (active, completed)")
	listCmd.Flags().BoolP("verbose", "v", false, "Show detailed session information")

	rootCmd.AddCommand(listCmd)
}