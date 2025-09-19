package hooks

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"context-extender/internal/config"
)

// ContextExtenderHooks defines the hooks that context-extender needs to install
type ContextExtenderHooks struct {
	SessionStart    config.HookEntry
	UserPromptSubmit config.HookEntry
	Stop            config.HookEntry
	SessionEnd      config.HookEntry
}

// GetContextExtenderHooks returns the hook configuration for context-extender
func GetContextExtenderHooks() (*ContextExtenderHooks, error) {
	// Get the path to the context-extender binary
	execPath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("failed to get executable path: %w", err)
	}

	// Use absolute path to ensure hooks work regardless of working directory
	execPath, err = filepath.Abs(execPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	hooks := &ContextExtenderHooks{
		SessionStart: config.HookEntry{
			Matcher: "",
			Hooks: []config.HookConfig{
				{
					Type:    "command",
					Command: fmt.Sprintf("%s capture --event=session-start", execPath),
					Timeout: 30,
				},
			},
		},
		UserPromptSubmit: config.HookEntry{
			Matcher: "",
			Hooks: []config.HookConfig{
				{
					Type:    "command",
					Command: fmt.Sprintf("%s capture --event=user-prompt", execPath),
					Timeout: 30,
				},
			},
		},
		Stop: config.HookEntry{
			Matcher: "",
			Hooks: []config.HookConfig{
				{
					Type:    "command",
					Command: fmt.Sprintf("%s capture --event=claude-response", execPath),
					Timeout: 30,
				},
			},
		},
		SessionEnd: config.HookEntry{
			Matcher: "",
			Hooks: []config.HookConfig{
				{
					Type:    "command",
					Command: fmt.Sprintf("%s capture --event=session-end", execPath),
					Timeout: 30,
				},
			},
		},
	}

	return hooks, nil
}

// InstallHooks installs context-extender hooks into Claude Code settings
func InstallHooks() error {
	// Validate Claude Code installation
	if err := config.ValidateClaudeInstallation(); err != nil {
		return fmt.Errorf("Claude Code validation failed: %w", err)
	}

	// Read existing settings
	settings, err := config.ReadClaudeSettings()
	if err != nil {
		return fmt.Errorf("failed to read Claude settings: %w", err)
	}

	// Get context-extender hooks
	contextHooks, err := GetContextExtenderHooks()
	if err != nil {
		return fmt.Errorf("failed to get context-extender hooks: %w", err)
	}

	// Initialize hooks map if it doesn't exist
	if settings.Hooks == nil {
		settings.Hooks = make(map[string][]config.HookEntry)
	}

	// Install each hook type
	hookTypes := map[string]config.HookEntry{
		"SessionStart":     contextHooks.SessionStart,
		"UserPromptSubmit": contextHooks.UserPromptSubmit,
		"Stop":             contextHooks.Stop,
		"SessionEnd":       contextHooks.SessionEnd,
	}

	for hookType, hookEntry := range hookTypes {
		if err := installSingleHook(settings, hookType, hookEntry); err != nil {
			return fmt.Errorf("failed to install %s hook: %w", hookType, err)
		}
	}

	// Write updated settings
	if err := config.WriteClaudeSettings(settings); err != nil {
		return fmt.Errorf("failed to write Claude settings: %w", err)
	}

	return nil
}

// installSingleHook installs a single hook, avoiding duplicates
func installSingleHook(settings *config.ClaudeSettings, hookType string, hookEntry config.HookEntry) error {
	// Get existing hooks for this type
	existingHooks := settings.Hooks[hookType]

	// Check if our hook is already installed
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	execPath, err = filepath.Abs(execPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Remove any existing context-extender hooks for this type
	var filteredHooks []config.HookEntry
	for _, entry := range existingHooks {
		hasContextExtender := false
		for _, hook := range entry.Hooks {
			if containsContextExtender(hook.Command, execPath) {
				hasContextExtender = true
				break
			}
		}
		if !hasContextExtender {
			filteredHooks = append(filteredHooks, entry)
		}
	}

	// Add our hook
	filteredHooks = append(filteredHooks, hookEntry)

	// Update settings
	settings.Hooks[hookType] = filteredHooks

	return nil
}

// containsContextExtender checks if a command refers to context-extender
func containsContextExtender(command, execPath string) bool {
	// Normalize paths for comparison
	cleanCommand := filepath.Clean(command)
	cleanExecPath := filepath.Clean(execPath)

	// Check for exact executable path match
	if cleanCommand == cleanExecPath {
		return true
	}

	// Extract just the command part (before any arguments)
	commandParts := strings.Fields(command)
	if len(commandParts) > 0 {
		commandPath := filepath.Clean(commandParts[0])
		if commandPath == cleanExecPath {
			return true
		}
	}

	// Check for context-extender in command
	commandBase := filepath.Base(cleanCommand)
	execBase := filepath.Base(cleanExecPath)

	return commandBase == execBase ||
		   commandBase == "context-extender" ||
		   commandBase == "context-extender.exe" ||
		   strings.Contains(command, "context-extender")
}

// UninstallHooks removes context-extender hooks from Claude Code settings
func UninstallHooks() error {
	// Read existing settings
	settings, err := config.ReadClaudeSettings()
	if err != nil {
		return fmt.Errorf("failed to read Claude settings: %w", err)
	}

	if settings.Hooks == nil {
		// No hooks to uninstall
		return nil
	}

	// Get executable path for comparison
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	execPath, err = filepath.Abs(execPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Remove context-extender hooks from each hook type
	hookTypes := []string{"SessionStart", "UserPromptSubmit", "Stop", "SessionEnd"}

	for _, hookType := range hookTypes {
		existingHooks := settings.Hooks[hookType]
		var filteredHooks []config.HookEntry

		for _, entry := range existingHooks {
			var filteredCommands []config.HookConfig
			for _, hook := range entry.Hooks {
				if !containsContextExtender(hook.Command, execPath) {
					filteredCommands = append(filteredCommands, hook)
				}
			}

			// Only keep the entry if it has remaining commands
			if len(filteredCommands) > 0 {
				entry.Hooks = filteredCommands
				filteredHooks = append(filteredHooks, entry)
			}
		}

		// Update or remove the hook type
		if len(filteredHooks) > 0 {
			settings.Hooks[hookType] = filteredHooks
		} else {
			delete(settings.Hooks, hookType)
		}
	}

	// Write updated settings
	if err := config.WriteClaudeSettings(settings); err != nil {
		return fmt.Errorf("failed to write Claude settings: %w", err)
	}

	return nil
}

// IsInstalled checks if context-extender hooks are currently installed
func IsInstalled() (bool, error) {
	// Read existing settings
	settings, err := config.ReadClaudeSettings()
	if err != nil {
		return false, fmt.Errorf("failed to read Claude settings: %w", err)
	}

	if settings.Hooks == nil {
		return false, nil
	}

	// Get executable path for comparison
	execPath, err := os.Executable()
	if err != nil {
		return false, fmt.Errorf("failed to get executable path: %w", err)
	}

	execPath, err = filepath.Abs(execPath)
	if err != nil {
		return false, fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Check if any of our required hooks are installed
	hookTypes := []string{"SessionStart", "UserPromptSubmit", "Stop", "SessionEnd"}
	installedCount := 0

	for _, hookType := range hookTypes {
		existingHooks := settings.Hooks[hookType]
		for _, entry := range existingHooks {
			for _, hook := range entry.Hooks {
				if containsContextExtender(hook.Command, execPath) {
					installedCount++
					break
				}
			}
		}
	}

	// Consider installed if at least 3 of 4 hooks are present
	// (allows for some flexibility in case a hook type is missing)
	return installedCount >= 3, nil
}

// GetInstallationStatus returns detailed information about hook installation
func GetInstallationStatus() (map[string]bool, error) {
	status := make(map[string]bool)

	// Read existing settings
	settings, err := config.ReadClaudeSettings()
	if err != nil {
		return status, fmt.Errorf("failed to read Claude settings: %w", err)
	}

	// Get executable path for comparison
	execPath, err := os.Executable()
	if err != nil {
		return status, fmt.Errorf("failed to get executable path: %w", err)
	}

	execPath, err = filepath.Abs(execPath)
	if err != nil {
		return status, fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Check each hook type
	hookTypes := []string{"SessionStart", "UserPromptSubmit", "Stop", "SessionEnd"}

	for _, hookType := range hookTypes {
		status[hookType] = false

		if settings.Hooks != nil {
			existingHooks := settings.Hooks[hookType]
			for _, entry := range existingHooks {
				for _, hook := range entry.Hooks {
					if containsContextExtender(hook.Command, execPath) {
						status[hookType] = true
						break
					}
				}
				if status[hookType] {
					break
				}
			}
		}
	}

	return status, nil
}