package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

// Version information (set during build)
var (
	Version   = "0.1.0"
	BuildDate = "unknown"
	GitCommit = "unknown"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "context-extender",
	Short: "Automatically capture and manage Claude Code conversations",
	Long: `Context-Extender is a CLI tool that automatically captures Claude Code
conversations and enables context sharing between different Claude sessions.

Features:
- Automatic conversation capture via Claude Code hooks
- Session correlation and storage management
- Context sharing between Claude Code sessions
- Cross-platform support (Windows, Mac, Linux)

Example usage:
  context-extender configure               # Set up Claude Code hooks
  context-extender list                    # List captured conversations
  context-extender share session-123      # Share conversation context
  context-extender config show            # Show current configuration`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// getPlatform returns a human-readable platform string
func getPlatform() string {
	return fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
}

func init() {
	// Global flags can be added here
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.context-extender.yaml)")
}