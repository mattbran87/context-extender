package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"

	"context-extender/internal/hooks"
)

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Completely remove Context-Extender from your system",
	Long: `Uninstall removes all Context-Extender components from your system:

This command will:
- Remove Claude Code hooks
- Delete database and conversation data (with confirmation)
- Remove system PATH installation (if present)
- Clean up all configuration files

Options:
  --keep-data    Remove hooks but keep database and conversations
  --force        Skip confirmation prompts (use with caution)

Examples:
  context-extender uninstall              # Interactive uninstall with prompts
  context-extender uninstall --keep-data  # Remove hooks but preserve data
  context-extender uninstall --force      # Force uninstall without prompts`,

	RunE: func(cmd *cobra.Command, args []string) error {
		keepData, _ := cmd.Flags().GetBool("keep-data")
		force, _ := cmd.Flags().GetBool("force")

		return runUninstall(keepData, force)
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)

	// Add flags
	uninstallCmd.Flags().BoolP("keep-data", "k", false, "Keep database and conversation data")
	uninstallCmd.Flags().BoolP("force", "f", false, "Skip confirmation prompts")
}

func runUninstall(keepData, force bool) error {
	fmt.Println("🗑️  Context-Extender Uninstall")
	fmt.Println("==============================")
	fmt.Println()

	// Show what will be removed
	showUninstallPlan(keepData)

	// Confirmation prompt (unless forced)
	if !force {
		if !confirmUninstall(keepData) {
			fmt.Println("Uninstall cancelled.")
			return nil
		}
	}

	fmt.Println()
	fmt.Println("🔄 Starting uninstall process...")
	fmt.Println()

	// Step 1: Remove Claude Code hooks
	if err := uninstallHooks(); err != nil {
		return fmt.Errorf("failed to remove hooks: %w", err)
	}

	// Step 2: Remove system PATH installation (if present)
	if err := removeSystemInstallation(); err != nil {
		// Non-fatal - just warn user
		fmt.Printf("⚠️  Warning: Could not remove system PATH installation: %v\n", err)
		fmt.Println("   You may need to manually remove context-extender from your PATH")
	}

	// Step 3: Remove data (if requested)
	if !keepData {
		if err := removeUserData(force); err != nil {
			return fmt.Errorf("failed to remove user data: %w", err)
		}
	} else {
		fmt.Println("📦 Keeping database and conversation data as requested")
	}

	// Show completion message
	showUninstallComplete(keepData)

	return nil
}

func showUninstallPlan(keepData bool) {
	fmt.Println("📋 Uninstall Plan:")
	fmt.Println()
	fmt.Println("✅ Remove Claude Code hooks (stop automatic capture)")
	fmt.Println("✅ Remove system PATH installation (if present)")

	if keepData {
		fmt.Println("📦 Keep database and conversation data")
	} else {
		fmt.Println("🗑️  Delete database and conversation data")
		fmt.Println("    ⚠️  This will permanently delete all captured conversations!")
	}

	fmt.Println()
}

func confirmUninstall(keepData bool) bool {
	reader := bufio.NewReader(os.Stdin)

	if keepData {
		fmt.Print("Remove Context-Extender but keep your data? (y/N): ")
	} else {
		fmt.Print("⚠️  This will DELETE ALL conversation data. Continue? (y/N): ")
	}

	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes"
}

func uninstallHooks() error {
	fmt.Print("🔗 Removing Claude Code hooks... ")

	// Check if hooks are installed
	installed, err := hooks.IsInstalled()
	if err != nil {
		fmt.Println("❌ ERROR")
		return fmt.Errorf("failed to check hook status: %w", err)
	}

	if !installed {
		fmt.Println("⚠️  Not installed")
		fmt.Println("   Claude Code hooks were not found (already removed)")
		return nil
	}

	// Remove hooks
	if err := hooks.UninstallHooks(); err != nil {
		fmt.Println("❌ FAILED")
		return err
	}

	// Verify removal
	installed, err = hooks.IsInstalled()
	if err != nil || installed {
		fmt.Println("❌ VERIFICATION FAILED")
		return fmt.Errorf("hook removal verification failed")
	}

	fmt.Println("✅ SUCCESS")
	return nil
}

func removeSystemInstallation() error {
	fmt.Print("🛠️  Checking system PATH installation... ")

	var pathFile string
	switch runtime.GOOS {
	case "windows":
		userApps := filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Local", "Microsoft", "WindowsApps")
		pathFile = filepath.Join(userApps, "context-extender.bat")
	case "darwin", "linux":
		// Check common locations for symlinks or scripts
		possiblePaths := []string{
			"/usr/local/bin/context-extender",
			filepath.Join(os.Getenv("HOME"), ".local", "bin", "context-extender"),
			filepath.Join(os.Getenv("HOME"), "bin", "context-extender"),
		}
		for _, path := range possiblePaths {
			if _, err := os.Stat(path); err == nil {
				pathFile = path
				break
			}
		}
	}

	// Check if PATH installation exists
	if pathFile == "" {
		fmt.Println("⚠️  Not found")
		return nil
	}

	if _, err := os.Stat(pathFile); os.IsNotExist(err) {
		fmt.Println("⚠️  Not found")
		return nil
	}

	// Remove PATH installation
	if err := os.Remove(pathFile); err != nil {
		fmt.Println("❌ FAILED")
		return err
	}

	fmt.Println("✅ SUCCESS")
	fmt.Printf("   Removed: %s\n", pathFile)
	return nil
}

func removeUserData(force bool) error {
	fmt.Print("🗑️  Removing database and conversation data... ")

	// Get user data directory
	dataDir, err := getUserDataDir()
	if err != nil {
		fmt.Println("❌ ERROR")
		return fmt.Errorf("failed to locate user data directory: %w", err)
	}

	// Check if data directory exists
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		fmt.Println("⚠️  Not found")
		fmt.Println("   No data directory found (nothing to remove)")
		return nil
	}

	// Final confirmation for data deletion (if not forced)
	if !force {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("\n   Data directory: %s\n", dataDir)
		fmt.Print("   🚨 Final confirmation - DELETE ALL DATA? (yes/no): ")

		response, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read confirmation: %w", err)
		}

		response = strings.TrimSpace(strings.ToLower(response))
		if response != "yes" {
			fmt.Println("❌ CANCELLED")
			fmt.Println("   Data deletion cancelled by user")
			return nil
		}
		fmt.Print("🗑️  Removing database and conversation data... ")
	}

	// Remove data directory
	if err := os.RemoveAll(dataDir); err != nil {
		fmt.Println("❌ FAILED")
		return fmt.Errorf("failed to remove data directory: %w", err)
	}

	fmt.Println("✅ SUCCESS")
	fmt.Printf("   Removed: %s\n", dataDir)
	return nil
}

func getUserDataDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".context-extender"), nil
}

func showUninstallComplete(keepData bool) {
	fmt.Println()
	fmt.Println("🎉 Uninstall Complete!")
	fmt.Println("======================")
	fmt.Println()
	fmt.Println("✅ Claude Code hooks removed")
	fmt.Println("✅ System PATH installation removed")

	if keepData {
		fmt.Println("📦 Your conversation data has been preserved")
		fmt.Println()
		fmt.Println("💡 To reinstall Context-Extender:")
		fmt.Println("   context-extender install    # Interactive wizard")
		fmt.Println("   context-extender configure  # Quick setup")
		fmt.Println()
		fmt.Println("   Your existing data will be automatically detected.")
	} else {
		fmt.Println("🗑️  All conversation data removed")
		fmt.Println()
		fmt.Println("💡 To reinstall Context-Extender:")
		fmt.Println("   context-extender install    # Fresh installation")
		fmt.Println("   context-extender configure  # Quick setup")
	}

	fmt.Println()
	fmt.Println("Thank you for using Context-Extender! 🙏")
}