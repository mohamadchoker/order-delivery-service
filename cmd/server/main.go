package main

import (
	"fmt"
	"os"
)

// Version information - set via ldflags during build
var (
	version   = "dev"
	buildDate = "unknown"
	gitCommit = "unknown"
)

func main() {
	// Create and initialize application
	application, err := NewApp(version, buildDate, gitCommit)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize application: %v\n", err)
		os.Exit(1)
	}

	// Run application (blocks until shutdown signal)
	if err := application.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Application error: %v\n", err)
		os.Exit(1)
	}
}
