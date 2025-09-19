package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"context-extender/internal/storage"
)

var storageCmd = &cobra.Command{
	Use:   "storage",
	Short: "Manage storage directories and configuration",
	Long:  "Commands for managing conversation storage directories, checking status, and performing maintenance.",
	Run: func(cmd *cobra.Command, args []string) {
		// Default action: show storage status
		showStorageStatus()
	},
}

var storageStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show storage directory status and information",
	Long:  "Display information about storage directories, usage, and accessibility.",
	Run: func(cmd *cobra.Command, args []string) {
		showStorageStatus()
	},
}

var storageInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize storage directories",
	Long:  "Create and validate storage directory structure for conversation data.",
	Run: func(cmd *cobra.Command, args []string) {
		initializeStorage()
	},
}

var storageCleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean temporary files from storage",
	Long:  "Remove temporary files and perform basic storage maintenance.",
	Run: func(cmd *cobra.Command, args []string) {
		cleanStorage()
	},
}

var storageUsageCmd = &cobra.Command{
	Use:   "usage",
	Short: "Show storage usage statistics",
	Long:  "Display detailed storage usage information and statistics.",
	Run: func(cmd *cobra.Command, args []string) {
		showStorageUsage()
	},
}

// Command flags
var (
	storagePath   string
	jsonOutput    bool
	verboseOutput bool
)

func init() {
	// Add subcommands
	storageCmd.AddCommand(storageStatusCmd)
	storageCmd.AddCommand(storageInitCmd)
	storageCmd.AddCommand(storageCleanCmd)
	storageCmd.AddCommand(storageUsageCmd)

	// Add flags
	storageCmd.PersistentFlags().StringVar(&storagePath, "path", "", "Custom storage path (overrides default)")
	storageCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")
	storageCmd.PersistentFlags().BoolVar(&verboseOutput, "verbose", false, "Verbose output")

	// Add to root command
	rootCmd.AddCommand(storageCmd)
}

// showStorageStatus displays current storage status and information
func showStorageStatus() {
	fmt.Println("📁 Context-Extender Storage Status")
	fmt.Println("==================================")

	// Create storage manager
	var config *storage.StorageConfig
	if storagePath != "" {
		config = &storage.StorageConfig{CustomPath: storagePath}
	}

	sm, err := storage.NewStorageManager(config)
	if err != nil {
		fmt.Printf("❌ Failed to initialize storage manager: %v\n", err)
		os.Exit(1)
	}

	// Get storage information
	info := sm.GetStorageInfo()

	if jsonOutput {
		output, err := json.MarshalIndent(info, "", "  ")
		if err != nil {
			fmt.Printf("❌ Failed to marshal JSON: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(string(output))
		return
	}

	// Display storage information
	fmt.Printf("🖥️  Platform: %s\n", info.Platform)
	fmt.Printf("📂 Base Directory: %s\n", info.BaseDir)
	fmt.Printf("💬 Conversations: %s\n", info.ConversationsDir)
	fmt.Printf("⚙️  Configuration: %s\n", info.ConfigDir)
	fmt.Printf("📝 Logs: %s\n", info.LogsDir)

	// Check if directories exist and are accessible
	fmt.Println("\n🔍 Directory Status:")
	checkDirectoryStatus("Base", info.BaseDir)
	checkDirectoryStatus("Conversations", info.ConversationsDir)
	checkDirectoryStatus("Configuration", info.ConfigDir)
	checkDirectoryStatus("Logs", info.LogsDir)

	// Validate storage access
	fmt.Println("\n✅ Access Validation:")
	if err := sm.ValidateStorageAccess(); err != nil {
		fmt.Printf("❌ Storage validation failed: %v\n", err)
	} else {
		fmt.Println("✅ All storage directories are accessible and writable")
	}
}

// checkDirectoryStatus checks and displays the status of a directory
func checkDirectoryStatus(name, path string) {
	if info, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("   ⚠️  %s: Not created\n", name)
		} else {
			fmt.Printf("   ❌ %s: Error - %v\n", name, err)
		}
	} else if !info.IsDir() {
		fmt.Printf("   ❌ %s: Not a directory\n", name)
	} else {
		fmt.Printf("   ✅ %s: Ready\n", name)
	}
}

// initializeStorage creates and validates storage directory structure
func initializeStorage() {
	fmt.Println("🚀 Initializing Context-Extender Storage")
	fmt.Println("=======================================")

	// Create storage manager
	var config *storage.StorageConfig
	if storagePath != "" {
		config = &storage.StorageConfig{CustomPath: storagePath}
		fmt.Printf("📁 Using custom storage path: %s\n", storagePath)
	}

	sm, err := storage.NewStorageManager(config)
	if err != nil {
		fmt.Printf("❌ Failed to initialize storage manager: %v\n", err)
		os.Exit(1)
	}

	// Display what will be created
	info := sm.GetStorageInfo()
	fmt.Printf("🖥️  Platform: %s\n", info.Platform)
	fmt.Printf("📂 Base Directory: %s\n", info.BaseDir)

	if verboseOutput {
		fmt.Println("\n📋 Directories to create:")
		fmt.Printf("   📂 %s\n", info.BaseDir)
		fmt.Printf("   💬 %s\n", info.ConversationsDir)
		fmt.Printf("   ⚙️  %s\n", info.ConfigDir)
		fmt.Printf("   📝 %s\n", info.LogsDir)
	}

	fmt.Println("\n🔨 Creating storage structure...")

	// Ensure storage structure
	if err := sm.EnsureStorageStructure(); err != nil {
		fmt.Printf("❌ Failed to create storage structure: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Storage structure created successfully")

	// Validate access
	fmt.Println("\n🔍 Validating storage access...")
	if err := sm.ValidateStorageAccess(); err != nil {
		fmt.Printf("❌ Storage validation failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Storage validation completed successfully")
	fmt.Println("\n🎉 Storage initialization complete! Ready to capture conversations.")
}

// cleanStorage performs storage cleanup operations
func cleanStorage() {
	fmt.Println("🧹 Cleaning Context-Extender Storage")
	fmt.Println("===================================")

	// Create storage manager
	var config *storage.StorageConfig
	if storagePath != "" {
		config = &storage.StorageConfig{CustomPath: storagePath}
	}

	sm, err := storage.NewStorageManager(config)
	if err != nil {
		fmt.Printf("❌ Failed to initialize storage manager: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("🔍 Scanning for temporary files...")

	// Perform cleanup
	if err := sm.CleanupStorage(); err != nil {
		fmt.Printf("❌ Cleanup failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Temporary files cleaned successfully")

	if verboseOutput {
		// Show updated usage after cleanup
		fmt.Println("\n📊 Updated storage usage:")
		showStorageUsageInternal(sm)
	}

	fmt.Println("\n🎉 Storage cleanup complete!")
}

// showStorageUsage displays storage usage statistics
func showStorageUsage() {
	fmt.Println("📊 Context-Extender Storage Usage")
	fmt.Println("=================================")

	// Create storage manager
	var config *storage.StorageConfig
	if storagePath != "" {
		config = &storage.StorageConfig{CustomPath: storagePath}
	}

	sm, err := storage.NewStorageManager(config)
	if err != nil {
		fmt.Printf("❌ Failed to initialize storage manager: %v\n", err)
		os.Exit(1)
	}

	showStorageUsageInternal(sm)
}

// showStorageUsageInternal displays storage usage with provided storage manager
func showStorageUsageInternal(sm *storage.StorageManager) {
	// Get storage usage
	usage, err := sm.GetStorageUsage()
	if err != nil {
		fmt.Printf("❌ Failed to calculate storage usage: %v\n", err)
		return
	}

	if jsonOutput {
		output, err := json.MarshalIndent(usage, "", "  ")
		if err != nil {
			fmt.Printf("❌ Failed to marshal JSON: %v\n", err)
			return
		}
		fmt.Println(string(output))
		return
	}

	// Display usage information
	fmt.Printf("💬 Conversations: %s\n", formatBytes(usage.ConversationsSize))
	fmt.Printf("⚙️  Configuration: %s\n", formatBytes(usage.ConfigSize))
	fmt.Printf("📝 Logs: %s\n", formatBytes(usage.LogsSize))
	fmt.Printf("📊 Total: %s\n", formatBytes(usage.TotalSize))

	if verboseOutput {
		info := sm.GetStorageInfo()
		fmt.Printf("\n📁 Directory Details:\n")
		fmt.Printf("   Conversations: %s\n", info.ConversationsDir)
		fmt.Printf("   Configuration: %s\n", info.ConfigDir)
		fmt.Printf("   Logs: %s\n", info.LogsDir)
	}
}

// formatBytes formats bytes into human-readable string
func formatBytes(bytes int64) string {
	if bytes == 0 {
		return "0 B"
	}

	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}