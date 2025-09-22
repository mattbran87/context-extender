package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"context-extender/internal/config"
	"context-extender/internal/database"
	"context-extender/internal/hooks"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Interactive installation wizard for first-time setup",
	Long: `Interactive installation wizard that guides you through the complete
setup process for Context-Extender.

This wizard will:
- Check system requirements
- Explain what Context-Extender does
- Guide you through Claude Code integration
- Test the installation
- Show you how to get started

Perfect for first-time users who want step-by-step guidance.

For quick installation, use: context-extender configure`,

	RunE: func(cmd *cobra.Command, args []string) error {
		return runInstallWizard()
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}

func runInstallWizard() error {
	fmt.Println("🎉 Welcome to Context-Extender Installation Wizard!")
	fmt.Println("==================================================")
	fmt.Println()

	// Introduction
	showIntroduction()
	if !promptContinue("Ready to start the installation?") {
		fmt.Println("Installation cancelled. Run 'context-extender install' anytime to try again.")
		return nil
	}

	fmt.Println()
	fmt.Println("📋 Step 1: System Requirements Check")
	fmt.Println("=====================================")

	// Check Claude Code
	if err := checkClaudeCode(); err != nil {
		return handleClaudeCodeError(err)
	}

	fmt.Println()
	fmt.Println("⚙️  Step 2: Installation Process")
	fmt.Println("=================================")

	// Explain what will happen
	explainInstallation()
	if !promptContinue("Proceed with hook installation?") {
		fmt.Println("Installation cancelled.")
		return nil
	}

	// Perform installation
	if err := performInstallation(); err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("🧪 Step 3: Testing Installation")
	fmt.Println("================================")

	// Test installation
	if err := testInstallation(); err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("🚀 Step 4: Getting Started")
	fmt.Println("===========================")

	// Show getting started guide
	showGettingStarted()

	fmt.Println()
	fmt.Println("🎉 Installation Complete!")
	fmt.Println("=========================")
	fmt.Println("Context-Extender is now ready to capture your Claude Code conversations!")
	fmt.Println()
	fmt.Println("💡 Tip: Start a new Claude Code session to begin capturing conversations.")
	fmt.Println("    Then run 'context-extender query list' to see your captured data.")

	return nil
}

func showIntroduction() {
	fmt.Println("📚 What is Context-Extender?")
	fmt.Println("Context-Extender automatically captures your Claude Code conversations")
	fmt.Println("and stores them in a local database. This enables:")
	fmt.Println()
	fmt.Println("  ✅ Conversation history across sessions")
	fmt.Println("  ✅ Search through previous conversations")
	fmt.Println("  ✅ Export data for analysis (CSV, JSON, Excel)")
	fmt.Println("  ✅ Context sharing between different Claude sessions")
	fmt.Println("  ✅ Analytics and usage insights")
	fmt.Println()
	fmt.Println("🔒 Privacy: All data is stored locally on your computer.")
	fmt.Println("🚀 Performance: Zero impact on Claude Code performance.")
	fmt.Println()
}

func checkClaudeCode() error {
	fmt.Print("Checking for Claude Code installation... ")

	if err := config.ValidateClaudeInstallation(); err != nil {
		fmt.Println("❌ NOT FOUND")
		return err
	}

	fmt.Println("✅ FOUND")

	// Get settings path to show user
	settingsPath, err := config.GetClaudeSettingsPath()
	if err == nil {
		fmt.Printf("   Claude settings: %s\n", settingsPath)
	}

	fmt.Println("✅ System requirements met!")
	return nil
}

func handleClaudeCodeError(err error) error {
	fmt.Println()
	fmt.Println("❌ Claude Code Not Found")
	fmt.Println("========================")
	fmt.Println()
	fmt.Println("Context-Extender requires Claude Code to be installed on your system.")
	fmt.Println()
	fmt.Println("📥 To install Claude Code:")
	fmt.Println("   1. Visit: https://claude.ai/code")
	fmt.Println("   2. Download and install Claude Code for your platform")
	fmt.Println("   3. Run Claude Code at least once to complete setup")
	fmt.Println("   4. Come back and run 'context-extender install' again")
	fmt.Println()
	fmt.Printf("Technical details: %v\n", err)
	fmt.Println()
	fmt.Println("Need help? Check the Context-Extender documentation or file an issue.")

	return fmt.Errorf("Claude Code installation required")
}

func explainInstallation() {
	fmt.Println("🔧 What happens during installation:")
	fmt.Println()
	fmt.Println("  1. ✅ Create backup of your current Claude settings")
	fmt.Println("  2. ✅ Install conversation capture hooks")
	fmt.Println("  3. ✅ Test that hooks are working properly")
	fmt.Println("  4. ✅ Initialize local database for storage")
	fmt.Println()
	fmt.Println("⚠️  This modifies your Claude Code settings to add conversation capture.")
	fmt.Println("   Your existing Claude settings will be backed up first.")
	fmt.Println("   You can always uninstall with 'context-extender configure --remove'")
	fmt.Println()
	fmt.Println("🔒 Security: Hooks only capture conversation data, no passwords or keys.")
	fmt.Println()
}

func performInstallation() error {
	fmt.Println("🔄 Installing Context-Extender...")
	fmt.Println()

	// Check if already installed
	fmt.Print("⚡ Checking existing installation... ")
	installed, err := hooks.IsInstalled()
	if err != nil {
		fmt.Println("❌ ERROR")
		return fmt.Errorf("failed to check installation status: %w", err)
	}

	if installed {
		fmt.Println("⚠️  ALREADY INSTALLED")
		fmt.Println()
		fmt.Println("Context-Extender hooks are already installed!")
		fmt.Println("Installation completed successfully (no changes needed).")
		return nil
	}
	fmt.Println("✅ READY")

	// Install hooks
	fmt.Print("⚡ Installing conversation capture hooks... ")
	if err := hooks.InstallHooks(); err != nil {
		fmt.Println("❌ FAILED")
		fmt.Println()
		fmt.Println("Installation failed. Your Claude settings have not been modified.")
		return fmt.Errorf("failed to install hooks: %w", err)
	}
	fmt.Println("✅ SUCCESS")

	// Verify installation
	fmt.Print("⚡ Verifying installation... ")
	installed, err = hooks.IsInstalled()
	if err != nil || !installed {
		fmt.Println("❌ VERIFICATION FAILED")
		return fmt.Errorf("installation verification failed")
	}
	fmt.Println("✅ VERIFIED")

	fmt.Println()
	fmt.Println("✅ Installation completed successfully!")

	return nil
}

func testInstallation() error {
	fmt.Println("🧪 Testing your installation...")
	fmt.Println()

	// Test database initialization
	fmt.Print("⚡ Testing database system... ")
	// Test database connection by trying to initialize
	if err := testDatabase(); err != nil {
		fmt.Println("❌ FAILED")
		return fmt.Errorf("database test failed: %w", err)
	}
	fmt.Println("✅ OK")

	// Test hook status
	fmt.Print("⚡ Testing hook integration... ")
	status, err := hooks.GetInstallationStatus()
	if err != nil {
		fmt.Println("❌ FAILED")
		return fmt.Errorf("hook status test failed: %w", err)
	}

	allInstalled := true
	for _, installed := range status {
		if !installed {
			allInstalled = false
			break
		}
	}

	if !allInstalled {
		fmt.Println("❌ FAILED")
		return fmt.Errorf("not all hooks are installed properly")
	}
	fmt.Println("✅ OK")

	fmt.Println()
	fmt.Println("✅ All tests passed! Installation is working correctly.")

	return nil
}

func showGettingStarted() {
	fmt.Println("🎓 How to use Context-Extender:")
	fmt.Println()
	fmt.Println("📝 1. CAPTURE CONVERSATIONS")
	fmt.Println("   • Start a new Claude Code session")
	fmt.Println("   • Have conversations with Claude as normal")
	fmt.Println("   • Context-Extender automatically captures everything!")
	fmt.Println()
	fmt.Println("🔍 2. VIEW YOUR DATA")
	fmt.Println("   • context-extender query list           # List all conversations")
	fmt.Println("   • context-extender query show <id>      # View specific conversation")
	fmt.Println("   • context-extender database status      # Check database stats")
	fmt.Println()
	fmt.Println("📊 3. EXPORT & ANALYZE")
	fmt.Println("   • context-extender export --format csv  # Export to Excel")
	fmt.Println("   • context-extender export --format json # Export to JSON")
	fmt.Println("   • context-extender export --format xlsx # Export to Excel with charts")
	fmt.Println()
	fmt.Println("⚙️  4. MANAGE INSTALLATION")
	fmt.Println("   • context-extender configure --status   # Check hook status")
	fmt.Println("   • context-extender configure --remove   # Uninstall hooks")
	fmt.Println("   • context-extender --help               # See all commands")
	fmt.Println()
}

func testDatabase() error {
	config := database.DefaultDatabaseConfig()
	manager := database.NewManager(config)

	ctx := context.Background()
	if err := manager.Initialize(ctx); err != nil {
		return err
	}
	defer manager.Close()

	// Test getting backend and ensure schema exists
	backend, err := manager.GetBackend()
	if err != nil {
		return err
	}

	// Test basic database operation to ensure schema is created
	_, err = backend.ListSessions(ctx, nil)
	return err
}

func promptContinue(message string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s (y/N): ", message)

	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes"
}