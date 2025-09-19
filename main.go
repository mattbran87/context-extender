package main

import "context-extender/cmd"

// Build information (set by ldflags during build)
var (
	version   = "dev"
	buildDate = "unknown"
	gitCommit = "unknown"
)

func main() {
	// Set build information for the CLI
	cmd.SetBuildInfo(version, buildDate, gitCommit)
	cmd.Execute()
}