package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long: `Display the current version of context-extender along with build information.

This includes the version number, build date, and git commit hash if available.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("context-extender version %s\n", Version)
		fmt.Printf("Build date: %s\n", BuildDate)
		fmt.Printf("Git commit: %s\n", GitCommit)
		fmt.Printf("Platform: %s\n", getPlatform())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}