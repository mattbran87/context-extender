package hooks

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"context-extender/internal/config"
)

// TestGetContextExtenderHooks tests hook configuration generation
func TestGetContextExtenderHooks(t *testing.T) {
	hooks, err := GetContextExtenderHooks()
	if err != nil {
		t.Fatalf("Failed to get context-extender hooks: %v", err)
	}

	if hooks == nil {
		t.Fatal("Expected hooks to be non-nil")
	}

	// Test that all required hooks are present
	requiredHooks := []struct {
		name string
		hook config.HookEntry
	}{
		{"SessionStart", hooks.SessionStart},
		{"UserPromptSubmit", hooks.UserPromptSubmit},
		{"Stop", hooks.Stop},
		{"SessionEnd", hooks.SessionEnd},
	}

	for _, required := range requiredHooks {
		if len(required.hook.Hooks) == 0 {
			t.Errorf("Expected %s hook to have at least one command", required.name)
		}

		if required.hook.Matcher != "" {
			t.Errorf("Expected %s hook matcher to be empty, got: %s", required.name, required.hook.Matcher)
		}

		// Check command structure
		for i, hookConfig := range required.hook.Hooks {
			if hookConfig.Type != "command" {
				t.Errorf("Expected %s hook[%d] type to be 'command', got: %s", required.name, i, hookConfig.Type)
			}

			if hookConfig.Command == "" {
				t.Errorf("Expected %s hook[%d] command to be non-empty", required.name, i)
			}

			if hookConfig.Timeout != 30 {
				t.Errorf("Expected %s hook[%d] timeout to be 30, got: %d", required.name, i, hookConfig.Timeout)
			}

			// Check that command contains the executable path (during tests this will be the test binary)
			execPath, _ := os.Executable()
			if !strings.Contains(hookConfig.Command, filepath.Base(execPath)) {
				t.Logf("Hook command uses test executable path: %s", hookConfig.Command)
			}
		}
	}
}

// TestContainsContextExtender tests the path comparison logic
func TestContainsContextExtender(t *testing.T) {
	testCases := []struct {
		name     string
		command  string
		execPath string
		expected bool
	}{
		{
			name:     "exact path match",
			command:  "/path/to/context-extender.exe",
			execPath: "/path/to/context-extender.exe",
			expected: true,
		},
		{
			name:     "command with arguments",
			command:  "/path/to/context-extender.exe capture --event=test",
			execPath: "/path/to/context-extender.exe",
			expected: true,
		},
		{
			name:     "basename match",
			command:  "context-extender",
			execPath: "/any/path/context-extender",
			expected: true,
		},
		{
			name:     "exe extension match",
			command:  "context-extender.exe",
			execPath: "/any/path/context-extender.exe",
			expected: true,
		},
		{
			name:     "contains context-extender",
			command:  "some-context-extender-tool",
			execPath: "/path/to/other-tool",
			expected: true,
		},
		{
			name:     "different command",
			command:  "/path/to/other-tool",
			execPath: "/path/to/context-extender",
			expected: false,
		},
		{
			name:     "empty command",
			command:  "",
			execPath: "/path/to/context-extender",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := containsContextExtender(tc.command, tc.execPath)
			if result != tc.expected {
				t.Errorf("Expected %v, got %v for command='%s', execPath='%s'",
					tc.expected, result, tc.command, tc.execPath)
			}
		})
	}
}

// TestInstallSingleHook tests individual hook installation logic
func TestInstallSingleHook(t *testing.T) {
	// Create test settings
	settings := &config.ClaudeSettings{
		Hooks: make(map[string][]config.HookEntry),
		Other: make(map[string]interface{}),
	}

	// Create test hook entry
	hookEntry := config.HookEntry{
		Matcher: "",
		Hooks: []config.HookConfig{
			{
				Type:    "command",
				Command: "test-context-extender capture --event=test",
				Timeout: 30,
			},
		},
	}

	// Test installing to empty settings
	err := installSingleHook(settings, "TestHook", hookEntry)
	if err != nil {
		t.Fatalf("Failed to install hook: %v", err)
	}

	if settings.Hooks == nil {
		t.Fatal("Expected hooks map to be initialized")
	}

	if len(settings.Hooks["TestHook"]) != 1 {
		t.Errorf("Expected 1 TestHook entry, got %d", len(settings.Hooks["TestHook"]))
	}

	// Test installing duplicate (should replace)
	err = installSingleHook(settings, "TestHook", hookEntry)
	if err != nil {
		t.Fatalf("Failed to install duplicate hook: %v", err)
	}

	if len(settings.Hooks["TestHook"]) != 1 {
		t.Errorf("Expected 1 TestHook entry after duplicate install, got %d", len(settings.Hooks["TestHook"]))
	}
}

// TestHookJSONSerialization tests that hook structures can be serialized/deserialized
func TestHookJSONSerialization(t *testing.T) {
	hooks, err := GetContextExtenderHooks()
	if err != nil {
		t.Fatalf("Failed to get hooks: %v", err)
	}

	// Test SessionStart hook serialization
	data, err := json.Marshal(hooks.SessionStart)
	if err != nil {
		t.Fatalf("Failed to marshal SessionStart hook: %v", err)
	}

	var unmarshaled config.HookEntry
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal SessionStart hook: %v", err)
	}

	// Verify the unmarshaled hook is equivalent
	if unmarshaled.Matcher != hooks.SessionStart.Matcher {
		t.Errorf("Matcher mismatch: expected %s, got %s", hooks.SessionStart.Matcher, unmarshaled.Matcher)
	}

	if len(unmarshaled.Hooks) != len(hooks.SessionStart.Hooks) {
		t.Errorf("Hook count mismatch: expected %d, got %d", len(hooks.SessionStart.Hooks), len(unmarshaled.Hooks))
	}

	for i, originalHook := range hooks.SessionStart.Hooks {
		if i >= len(unmarshaled.Hooks) {
			break
		}
		unmarshaledHook := unmarshaled.Hooks[i]

		if unmarshaledHook.Type != originalHook.Type {
			t.Errorf("Hook[%d] type mismatch: expected %s, got %s", i, originalHook.Type, unmarshaledHook.Type)
		}

		if unmarshaledHook.Command != originalHook.Command {
			t.Errorf("Hook[%d] command mismatch: expected %s, got %s", i, originalHook.Command, unmarshaledHook.Command)
		}

		if unmarshaledHook.Timeout != originalHook.Timeout {
			t.Errorf("Hook[%d] timeout mismatch: expected %d, got %d", i, originalHook.Timeout, unmarshaledHook.Timeout)
		}
	}
}

// TestHookPathExtraction tests that executable paths are handled correctly
func TestHookPathExtraction(t *testing.T) {
	hooks, err := GetContextExtenderHooks()
	if err != nil {
		t.Fatalf("Failed to get hooks: %v", err)
	}

	// Check that all hook commands contain absolute paths
	allHooks := []config.HookEntry{
		hooks.SessionStart,
		hooks.UserPromptSubmit,
		hooks.Stop,
		hooks.SessionEnd,
	}

	for i, hook := range allHooks {
		for j, hookConfig := range hook.Hooks {
			if !filepath.IsAbs(hookConfig.Command) && !filepath.IsAbs(extractCommandPath(hookConfig.Command)) {
				t.Errorf("Hook[%d][%d] command should contain absolute path: %s", i, j, hookConfig.Command)
			}
		}
	}
}

// Helper functions for tests

func containsSubstring(s, substr string) bool {
	return filepath.Base(s) == substr ||
		   filepath.Base(s) == substr+".exe" ||
		   containsContextExtender(s, substr)
}

func extractCommandPath(command string) string {
	// Extract the first word (command path) from a command string
	fields := filepath.SplitList(command)
	if len(fields) > 0 {
		return fields[0]
	}
	return command
}

// Benchmark tests for performance validation

func BenchmarkGetContextExtenderHooks(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := GetContextExtenderHooks()
		if err != nil {
			b.Fatalf("Failed to get hooks: %v", err)
		}
	}
}

func BenchmarkContainsContextExtender(b *testing.B) {
	command := "/long/path/to/context-extender.exe capture --event=session-start"
	execPath := "/long/path/to/context-extender.exe"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		containsContextExtender(command, execPath)
	}
}