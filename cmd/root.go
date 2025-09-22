package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

// Version information (set during build)
var (
	Version   = "1.2.0"
	BuildDate = "unknown"
	GitCommit = "unknown"
)

// SetBuildInfo sets the build information from main
func SetBuildInfo(version, buildDate, gitCommit string) {
	if version != "" && version != "dev" {
		Version = version
	}
	if buildDate != "" && buildDate != "unknown" {
		BuildDate = buildDate
	}
	if gitCommit != "" && gitCommit != "unknown" {
		GitCommit = gitCommit
	}
}

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
- Export to CSV, JSON, Excel formats
- Cross-platform support (Windows, Mac, Linux)

🚀 First-time setup:
  context-extender install                 # Interactive installation wizard

🔧 Quick setup:
  context-extender configure               # Set up Claude Code hooks
  context-extender configure --status      # Check installation status

📊 Using the tool:
  context-extender query list              # List captured conversations
  context-extender export --format xlsx    # Export to Excel
  context-extender database status         # Check database status

🗑️  Uninstall:
  context-extender uninstall               # Remove completely from system
  context-extender uninstall --keep-data   # Remove but keep conversations`,
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