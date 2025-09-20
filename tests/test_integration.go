package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	fmt.Println("🧪 Integration Test Suite - Cycle 4 Fixes")
	fmt.Println("==========================================")

	binary := "../context-extender-v1.0.1.exe"
	passed := 0
	total := 0

	// Test 1: Version command
	total++
	fmt.Printf("\n%d. Testing version command... ", total)
	if testCommand(binary, "version") {
		fmt.Println("✅ PASS")
		passed++
	} else {
		fmt.Println("❌ FAIL")
	}

	// Test 2: Help command
	total++
	fmt.Printf("%d. Testing help command... ", total)
	if testCommand(binary, "--help") {
		fmt.Println("✅ PASS")
		passed++
	} else {
		fmt.Println("❌ FAIL")
	}

	// Test 3: Database init
	total++
	fmt.Printf("%d. Testing database init... ", total)
	if testCommand(binary, "database", "init") {
		fmt.Println("✅ PASS")
		passed++
	} else {
		fmt.Println("❌ FAIL")
	}

	// Test 4: Database status
	total++
	fmt.Printf("%d. Testing database status... ", total)
	output, err := runCommand(binary, "database", "status")
	if err == nil && strings.Contains(output, "Pure Go SQLite") && strings.Contains(output, "CGO Required: false") {
		fmt.Println("✅ PASS")
		passed++
	} else {
		fmt.Println("❌ FAIL")
	}

	// Test 5: Capture command exists
	total++
	fmt.Printf("%d. Testing capture command exists... ", total)
	if testCommand(binary, "capture", "--help") {
		fmt.Println("✅ PASS")
		passed++
	} else {
		fmt.Println("❌ FAIL")
	}

	// Test 6: Capture session-start
	total++
	fmt.Printf("%d. Testing capture session-start... ", total)
	sessionID := fmt.Sprintf("test-session-%d", time.Now().Unix())
	cmd := exec.Command(binary, "capture", "--event=session-start")
	cmd.Env = append(os.Environ(), fmt.Sprintf("CLAUDE_SESSION_ID=%s", sessionID))
	outputBytes, err := cmd.CombinedOutput()
	output = string(outputBytes)
	if err == nil && strings.Contains(output, "started") {
		fmt.Println("✅ PASS")
		passed++
	} else {
		fmt.Println("❌ FAIL")
		fmt.Printf("    Error: %v\n", err)
		fmt.Printf("    Output: %s\n", output)
	}

	// Test 7: Capture user-prompt
	total++
	fmt.Printf("%d. Testing capture user-prompt... ", total)
	cmd = exec.Command(binary, "capture", "--event=user-prompt", "--data=Test message")
	cmd.Env = append(os.Environ(), fmt.Sprintf("CLAUDE_SESSION_ID=%s", sessionID))
	outputBytes, err = cmd.CombinedOutput()
	output = string(outputBytes)
	if err == nil && strings.Contains(output, "captured") {
		fmt.Println("✅ PASS")
		passed++
	} else {
		fmt.Println("❌ FAIL")
		fmt.Printf("    Error: %v\n", err)
		fmt.Printf("    Output: %s\n", output)
	}

	// Test 8: GraphQL stats
	total++
	fmt.Printf("%d. Testing GraphQL stats... ", total)
	output, err = runCommand(binary, "graphql", "stats")
	if err == nil && strings.Contains(output, "Sessions:") && strings.Contains(output, "Conversations:") {
		fmt.Println("✅ PASS")
		passed++
	} else {
		fmt.Println("❌ FAIL")
		fmt.Printf("    Error: %v\n", err)
	}

	// Test 9: Storage status
	total++
	fmt.Printf("%d. Testing storage status... ", total)
	if testCommand(binary, "storage", "status") {
		fmt.Println("✅ PASS")
		passed++
	} else {
		fmt.Println("❌ FAIL")
	}

	// Test 10: Query list
	total++
	fmt.Printf("%d. Testing query list... ", total)
	if testCommand(binary, "query", "list") {
		fmt.Println("✅ PASS")
		passed++
	} else {
		fmt.Println("❌ FAIL")
	}

	// Summary
	fmt.Printf("\n🎯 Test Results: %d/%d passed (%.1f%%)\n", passed, total, float64(passed)/float64(total)*100)

	if passed == total {
		fmt.Println("🎉 ALL TESTS PASSED - Ready for v1.0.1 release!")
		os.Exit(0)
	} else {
		fmt.Printf("❌ %d tests failed - Needs more work\n", total-passed)
		os.Exit(1)
	}
}

// testCommand runs a command and returns true if it exits successfully
func testCommand(binary string, args ...string) bool {
	_, err := runCommand(binary, args...)
	return err == nil
}

// runCommand runs a command and returns output and error
func runCommand(binary string, args ...string) (string, error) {
	cmd := exec.Command(binary, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}