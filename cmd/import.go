package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"context-extender/internal/database"
	"context-extender/internal/importer"
	"github.com/spf13/cobra"
)

var (
	importVerbose      bool
	importDryRun       bool
	importSkipExisting bool
	importMaxFiles     int
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import Claude conversations to database",
	Long: `Import existing Claude conversation files (JSONL format) into the Context Extender database.

This command can automatically discover and import all Claude conversations from your system,
or import specific files/directories.`,
}

var importAutoCmd = &cobra.Command{
	Use:   "auto",
	Short: "Automatically import all Claude conversations",
	Long: `Automatically discover and import all Claude conversation files from standard locations.

This will search for Claude conversation files in:
- ~/.claude/projects/
- ~/Library/Application Support/Claude/projects/ (macOS)
- ~/AppData/Roaming/Claude/projects/ (Windows)

The import process will:
1. Find all JSONL conversation files
2. Parse each conversation
3. Import messages, sessions, and metadata to the database
4. Track imported files to avoid duplicates`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Initialize database
		config := database.DefaultConfig()
		if err := database.Initialize(config); err != nil {
			return fmt.Errorf("failed to initialize database: %w", err)
		}

		// Run migrations
		if err := database.RunMigrations(); err != nil {
			return fmt.Errorf("failed to run migrations: %w", err)
		}

		fmt.Println("🔍 Searching for Claude conversation files...")

		// Find all Claude conversations
		files, err := importer.FindClaudeConversations()
		if err != nil {
			return fmt.Errorf("failed to find conversations: %w", err)
		}

		if len(files) == 0 {
			fmt.Println("No Claude conversation files found.")
			fmt.Println("\nTried locations:")
			homeDir, _ := os.UserHomeDir()
			fmt.Printf("  - %s\n", filepath.Join(homeDir, ".claude", "projects"))
			fmt.Printf("  - %s\n", filepath.Join(homeDir, "Library", "Application Support", "Claude", "projects"))
			fmt.Printf("  - %s\n", filepath.Join(homeDir, "AppData", "Roaming", "Claude", "projects"))
			return nil
		}

		fmt.Printf("Found %d conversation file(s)\n\n", len(files))

		// Show preview
		fmt.Println("Files to import:")
		for i, file := range files {
			if i >= 5 && !importVerbose {
				fmt.Printf("  ... and %d more\n", len(files)-5)
				break
			}
			fileInfo, _ := os.Stat(file)
			size := float64(fileInfo.Size()) / 1024 / 1024 // MB
			fmt.Printf("  - %s (%.2f MB)\n", filepath.Base(filepath.Dir(file)), size)
		}

		if importDryRun {
			fmt.Println("\n[DRY RUN MODE] No changes will be made")
		} else {
			fmt.Print("\nProceed with import? (y/N): ")
			var response string
			fmt.Scanln(&response)
			if response != "y" && response != "Y" {
				fmt.Println("Import cancelled")
				return nil
			}
		}

		// Create import manager
		options := importer.ImportOptions{
			Verbose:      importVerbose,
			DryRun:       importDryRun,
			SkipExisting: importSkipExisting,
			MaxFiles:     importMaxFiles,
		}
		manager := importer.NewImportManager(options)

		// Perform import
		fmt.Println("\n📥 Starting import...")
		result, err := manager.ImportAllClaude()
		if err != nil {
			return fmt.Errorf("import failed: %w", err)
		}

		// Display results
		fmt.Println("\n📊 Import Results:")
		fmt.Printf("  Total Files:      %d\n", result.TotalFiles)
		fmt.Printf("  Successful:       %d\n", result.SuccessfulFiles)
		fmt.Printf("  Failed:           %d\n", result.FailedFiles)
		fmt.Printf("  Skipped:          %d\n", result.SkippedFiles)
		fmt.Printf("  Duration:         %s\n", result.ImportDuration.Round(time.Second))

		if len(result.Errors) > 0 && importVerbose {
			fmt.Println("\n❌ Errors:")
			for _, err := range result.Errors {
				fmt.Printf("  - %s\n", err)
			}
		}

		if !importDryRun && result.SuccessfulFiles > 0 {
			fmt.Println("\n✅ Import completed successfully!")
			fmt.Printf("   Imported %d conversation(s) to the database\n", result.SuccessfulFiles)
		}

		return nil
	},
}

var importFileCmd = &cobra.Command{
	Use:   "file [path]",
	Short: "Import a specific Claude JSONL file",
	Long: `Import a specific Claude conversation file into the database.

Example:
  context-extender import file ~/.claude/projects/my-project/conversation.jsonl`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]

		// Check if file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return fmt.Errorf("file not found: %s", filePath)
		}

		// Initialize database
		config := database.DefaultConfig()
		if err := database.Initialize(config); err != nil {
			return fmt.Errorf("failed to initialize database: %w", err)
		}

		if err := database.RunMigrations(); err != nil {
			return fmt.Errorf("failed to run migrations: %w", err)
		}

		// Create import manager
		options := importer.ImportOptions{
			Verbose:      importVerbose,
			DryRun:       importDryRun,
			SkipExisting: importSkipExisting,
		}
		manager := importer.NewImportManager(options)

		fmt.Printf("Importing file: %s\n", filePath)

		// Import the file
		if err := manager.ImportFile(filePath); err != nil {
			return fmt.Errorf("failed to import file: %w", err)
		}

		if !importDryRun {
			fmt.Println("✅ File imported successfully!")
		}

		return nil
	},
}

var importDirCmd = &cobra.Command{
	Use:   "dir [path]",
	Short: "Import all JSONL files from a directory",
	Long: `Import all Claude conversation files from a specified directory.

This will recursively search for all .jsonl files in the directory and import them.

Example:
  context-extender import dir ~/.claude/projects/`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dirPath := args[0]

		// Check if directory exists
		if info, err := os.Stat(dirPath); os.IsNotExist(err) || !info.IsDir() {
			return fmt.Errorf("directory not found: %s", dirPath)
		}

		// Initialize database
		config := database.DefaultConfig()
		if err := database.Initialize(config); err != nil {
			return fmt.Errorf("failed to initialize database: %w", err)
		}

		if err := database.RunMigrations(); err != nil {
			return fmt.Errorf("failed to run migrations: %w", err)
		}

		// Create import manager
		options := importer.ImportOptions{
			Verbose:      importVerbose,
			DryRun:       importDryRun,
			SkipExisting: importSkipExisting,
			MaxFiles:     importMaxFiles,
		}
		manager := importer.NewImportManager(options)

		fmt.Printf("Importing from directory: %s\n", dirPath)

		// Import the directory
		result, err := manager.ImportDirectory(dirPath)
		if err != nil {
			return fmt.Errorf("failed to import directory: %w", err)
		}

		// Display results
		fmt.Println("\n📊 Import Results:")
		fmt.Printf("  Total Files:      %d\n", result.TotalFiles)
		fmt.Printf("  Successful:       %d\n", result.SuccessfulFiles)
		fmt.Printf("  Failed:           %d\n", result.FailedFiles)
		fmt.Printf("  Duration:         %s\n", result.ImportDuration.Round(time.Second))

		if !importDryRun && result.SuccessfulFiles > 0 {
			fmt.Println("\n✅ Directory imported successfully!")
		}

		return nil
	},
}

var importHistoryCmd = &cobra.Command{
	Use:   "history",
	Short: "Show import history",
	Long:  `Display the history of imported Claude conversation files.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Initialize database
		config := database.DefaultConfig()
		if err := database.Initialize(config); err != nil {
			return fmt.Errorf("failed to initialize database: %w", err)
		}

		// Get import history
		records, err := importer.GetImportHistory()
		if err != nil {
			return fmt.Errorf("failed to get import history: %w", err)
		}

		if len(records) == 0 {
			fmt.Println("No import history found.")
			return nil
		}

		fmt.Println("📜 Import History:")
		fmt.Println(strings.Repeat("-", 80))

		for _, record := range records {
			// Parse time for better display
			importTime := record.ImportedAt
			if t, err := time.Parse(time.RFC3339, record.ImportedAt); err == nil {
				importTime = t.Format("2006-01-02 15:04:05")
			}

			fmt.Printf("File: %s\n", filepath.Base(record.FilePath))
			fmt.Printf("  Path:     %s\n", record.FilePath)
			fmt.Printf("  Imported: %s\n", importTime)
			fmt.Printf("  Messages: %d\n", record.EventCount)
			if record.Checksum != "" {
				fmt.Printf("  Checksum: %s\n", record.Checksum[:12]+"...")
			}
			fmt.Println()
		}

		fmt.Printf("Total imported files: %d\n", len(records))

		return nil
	},
}

var importWizardCmd = &cobra.Command{
	Use:   "wizard",
	Short: "Interactive import wizard",
	Long:  `Launch an interactive wizard to guide you through importing Claude conversations.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🧙 Claude Conversation Import Wizard")
		fmt.Println("====================================")
		fmt.Println()

		// Initialize database
		fmt.Print("Initializing database... ")
		config := database.DefaultConfig()
		if err := database.Initialize(config); err != nil {
			fmt.Println("❌")
			return fmt.Errorf("failed to initialize database: %w", err)
		}

		if err := database.RunMigrations(); err != nil {
			fmt.Println("❌")
			return fmt.Errorf("failed to run migrations: %w", err)
		}
		fmt.Println("✅")

		// Search for conversations
		fmt.Print("Searching for Claude conversations... ")
		files, err := importer.FindClaudeConversations()
		if err != nil {
			fmt.Println("❌")
			return fmt.Errorf("failed to find conversations: %w", err)
		}
		fmt.Printf("✅ Found %d file(s)\n\n", len(files))

		if len(files) == 0 {
			fmt.Println("No Claude conversation files found on your system.")
			fmt.Print("\nWould you like to import from a custom location? (y/N): ")

			var response string
			fmt.Scanln(&response)

			if response == "y" || response == "Y" {
				fmt.Print("Enter the path to import from: ")
				var customPath string
				fmt.Scanln(&customPath)

				// Import custom path
				options := importer.ImportOptions{
					Verbose:      true,
					SkipExisting: true,
				}
				manager := importer.NewImportManager(options)

				if info, err := os.Stat(customPath); err == nil && info.IsDir() {
					result, err := manager.ImportDirectory(customPath)
					if err != nil {
						return fmt.Errorf("import failed: %w", err)
					}
					fmt.Printf("\n✅ Imported %d file(s) successfully!\n", result.SuccessfulFiles)
				} else {
					if err := manager.ImportFile(customPath); err != nil {
						return fmt.Errorf("import failed: %w", err)
					}
					fmt.Println("\n✅ File imported successfully!")
				}
			}
			return nil
		}

		// Show summary by project
		projectFiles := make(map[string][]string)
		for _, file := range files {
			project := filepath.Base(filepath.Dir(file))
			projectFiles[project] = append(projectFiles[project], file)
		}

		fmt.Println("📁 Found conversations in these projects:")
		for project, pFiles := range projectFiles {
			// Clean up project name
			cleanName := importer.GetProjectName(project)
			totalSize := int64(0)
			for _, f := range pFiles {
				if info, err := os.Stat(f); err == nil {
					totalSize += info.Size()
				}
			}
			sizeMB := float64(totalSize) / 1024 / 1024
			fmt.Printf("  • %s (%d conversation(s), %.2f MB)\n", cleanName, len(pFiles), sizeMB)
		}

		// Import options
		fmt.Println("\n📝 Import Options:")
		fmt.Println("  1. Import all conversations")
		fmt.Println("  2. Import specific project")
		fmt.Println("  3. Skip duplicate imports")
		fmt.Println("  4. Cancel")

		fmt.Print("\nSelect option (1-4): ")
		var choice string
		fmt.Scanln(&choice)

		options := importer.ImportOptions{
			Verbose:      true,
			SkipExisting: true,
		}
		manager := importer.NewImportManager(options)

		switch choice {
		case "1":
			fmt.Println("\n📥 Importing all conversations...")
			result, err := manager.ImportAllClaude()
			if err != nil {
				return fmt.Errorf("import failed: %w", err)
			}
			fmt.Printf("\n✅ Import complete! Successfully imported %d/%d files\n",
				result.SuccessfulFiles, result.TotalFiles)

		case "2":
			fmt.Print("\nEnter project name to import: ")
			var projectName string
			fmt.Scanln(&projectName)

			// Find matching project
			var matchedFiles []string
			for project, pFiles := range projectFiles {
				if strings.Contains(strings.ToLower(project), strings.ToLower(projectName)) {
					matchedFiles = append(matchedFiles, pFiles...)
				}
			}

			if len(matchedFiles) == 0 {
				fmt.Printf("No project found matching '%s'\n", projectName)
				return nil
			}

			fmt.Printf("Importing %d file(s) from project...\n", len(matchedFiles))
			successCount := 0
			for _, file := range matchedFiles {
				if err := manager.ImportFile(file); err != nil {
					fmt.Printf("  ❌ Failed: %s\n", filepath.Base(file))
				} else {
					successCount++
					fmt.Printf("  ✅ Imported: %s\n", filepath.Base(file))
				}
			}
			fmt.Printf("\n✅ Import complete! Successfully imported %d/%d files\n",
				successCount, len(matchedFiles))

		case "3":
			options.SkipExisting = true
			manager = importer.NewImportManager(options)
			fmt.Println("\n📥 Importing new conversations only...")
			result, err := manager.ImportAllClaude()
			if err != nil {
				return fmt.Errorf("import failed: %w", err)
			}
			fmt.Printf("\n✅ Import complete! Imported %d new files, skipped %d existing\n",
				result.SuccessfulFiles, result.SkippedFiles)

		case "4":
			fmt.Println("Import cancelled")
			return nil

		default:
			fmt.Println("Invalid option")
			return nil
		}

		return nil
	},
}

func init() {
	// Add flags
	importCmd.PersistentFlags().BoolVarP(&importVerbose, "verbose", "v", false, "Verbose output")
	importCmd.PersistentFlags().BoolVar(&importDryRun, "dry-run", false, "Preview import without making changes")
	importCmd.PersistentFlags().BoolVar(&importSkipExisting, "skip-existing", true, "Skip already imported files")
	importCmd.PersistentFlags().IntVar(&importMaxFiles, "max-files", 0, "Maximum number of files to import (0=unlimited)")

	// Add subcommands
	importCmd.AddCommand(importAutoCmd)
	importCmd.AddCommand(importFileCmd)
	importCmd.AddCommand(importDirCmd)
	importCmd.AddCommand(importHistoryCmd)
	importCmd.AddCommand(importWizardCmd)

	rootCmd.AddCommand(importCmd)
}