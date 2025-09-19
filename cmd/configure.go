package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"context-extender/internal/config"
	"context-extender/internal/hooks"
)

// configureCmd represents the configure command
var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure Claude Code hooks for conversation capture",
	Long: `Configure automatically installs the necessary hooks into Claude Code
to enable automatic conversation capture.

This command will:
- Validate Claude Code installation
- Create a backup of existing settings
- Install context-extender hooks for conversation capture
- Verify the installation was successful

The hooks will capture:
- Session start/end events
- User prompts submitted to Claude
- Claude responses and tool usage

Example:
  context-extender configure          # Install hooks
  context-extender configure --status # Check installation status
  context-extender configure --remove # Remove hooks`,

	RunE: func(cmd *cobra.Command, args []string) error {
		// Check flags
		status, _ := cmd.Flags().GetBool("status")
		remove, _ := cmd.Flags().GetBool("remove")

		if status {
			return showInstallationStatus()
		}

		if remove {
			return removeHooks()
		}

		return installHooks()
	},
}

func init() {
	rootCmd.AddCommand(configureCmd)

	// Add flags
	configureCmd.Flags().BoolP("status", "s", false, "Show current hook installation status")
	configureCmd.Flags().BoolP("remove", "r", false, "Remove context-extender hooks")
}

func installHooks() error {
	fmt.Println("Configuring Claude Code hooks for context-extender...")

	// Validate Claude Code installation
	fmt.Print("Validating Claude Code installation... ")
	if err := config.ValidateClaudeInstallation(); err != nil {
		fmt.Println("❌ FAILED")
		return fmt.Errorf("Claude Code validation failed: %w", err)
	}
	fmt.Println("✅ OK")

	// Check if already installed
	fmt.Print("Checking existing installation... ")
	installed, err := hooks.IsInstalled()
	if err != nil {
		fmt.Println("❌ ERROR")
		return fmt.Errorf("failed to check installation status: %w", err)
	}

	if installed {
		fmt.Println("⚠️  Already installed")
		fmt.Println("Context-extender hooks are already installed.")
		fmt.Println("Use --remove to uninstall or --status to see details.")
		return nil
	}
	fmt.Println("✅ Ready to install")

	// Install hooks
	fmt.Print("Installing conversation capture hooks... ")
	if err := hooks.InstallHooks(); err != nil {
		fmt.Println("❌ FAILED")
		return fmt.Errorf("failed to install hooks: %w", err)
	}
	fmt.Println("✅ SUCCESS")

	// Verify installation
	fmt.Print("Verifying installation... ")
	installed, err = hooks.IsInstalled()
	if err != nil {
		fmt.Println("❌ ERROR")
		return fmt.Errorf("failed to verify installation: %w", err)
	}

	if !installed {
		fmt.Println("❌ FAILED")
		return fmt.Errorf("installation verification failed")
	}
	fmt.Println("✅ VERIFIED")

	// Show success message
	fmt.Println("\n🎉 Configuration completed successfully!")
	fmt.Println("\nContext-extender is now configured to automatically capture")
	fmt.Println("Claude Code conversations. The following hooks are installed:")
	fmt.Println("  • SessionStart - Captures session initialization")
	fmt.Println("  • UserPromptSubmit - Captures your prompts to Claude")
	fmt.Println("  • Stop - Captures Claude's responses")
	fmt.Println("  • SessionEnd - Captures session completion")
	fmt.Println("\nYour conversations will be automatically captured starting")
	fmt.Println("with your next Claude Code session.")

	return nil
}

func removeHooks() error {
	fmt.Println("Removing context-extender hooks from Claude Code...")

	// Check if installed
	fmt.Print("Checking current installation... ")
	installed, err := hooks.IsInstalled()
	if err != nil {
		fmt.Println("❌ ERROR")
		return fmt.Errorf("failed to check installation status: %w", err)
	}

	if !installed {
		fmt.Println("⚠️  Not installed")
		fmt.Println("Context-extender hooks are not currently installed.")
		return nil
	}
	fmt.Println("✅ Found hooks")

	// Remove hooks
	fmt.Print("Removing hooks... ")
	if err := hooks.UninstallHooks(); err != nil {
		fmt.Println("❌ FAILED")
		return fmt.Errorf("failed to remove hooks: %w", err)
	}
	fmt.Println("✅ SUCCESS")

	// Verify removal
	fmt.Print("Verifying removal... ")
	installed, err = hooks.IsInstalled()
	if err != nil {
		fmt.Println("❌ ERROR")
		return fmt.Errorf("failed to verify removal: %w", err)
	}

	if installed {
		fmt.Println("❌ FAILED")
		return fmt.Errorf("removal verification failed")
	}
	fmt.Println("✅ VERIFIED")

	fmt.Println("\n✅ Context-extender hooks have been successfully removed.")
	fmt.Println("Claude Code will no longer capture conversations automatically.")

	return nil
}

func showInstallationStatus() error {
	fmt.Println("Context-extender hook installation status:")
	fmt.Println()

	// Check Claude Code installation
	fmt.Print("Claude Code installation: ")
	if err := config.ValidateClaudeInstallation(); err != nil {
		fmt.Printf("❌ Not found or inaccessible (%v)\n", err)
		return nil
	}
	fmt.Println("✅ Found and accessible")

	// Get detailed status
	status, err := hooks.GetInstallationStatus()
	if err != nil {
		return fmt.Errorf("failed to get installation status: %w", err)
	}

	// Check overall installation
	installed, err := hooks.IsInstalled()
	if err != nil {
		return fmt.Errorf("failed to check installation status: %w", err)
	}

	fmt.Printf("Overall status: ")
	if installed {
		fmt.Println("✅ Installed and ready")
	} else {
		fmt.Println("❌ Not installed")
	}

	fmt.Println("\nHook details:")
	hookTypes := []string{"SessionStart", "UserPromptSubmit", "Stop", "SessionEnd"}
	for _, hookType := range hookTypes {
		fmt.Printf("  %-18s ", hookType+":")
		if status[hookType] {
			fmt.Println("✅ Installed")
		} else {
			fmt.Println("❌ Missing")
		}
	}

	if !installed {
		fmt.Println("\nTo install hooks, run: context-extender configure")
	}

	// Show Claude settings path
	settingsPath, err := config.GetClaudeSettingsPath()
	if err == nil {
		fmt.Printf("\nClaude settings location: %s\n", settingsPath)
	}

	return nil
}