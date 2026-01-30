package main

import (
	"fmt"
	"os"

	"github.com/otto-nation/otto-stack/cmd/codegen"
)

const (
	ConfigDir        = "internal/config"
	MessagesYAMLPath = ConfigDir + "/messages.yaml"
	OutputDir        = "internal/pkg/messages"
	OutputPath       = OutputDir + "/messages_generated.go"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Ensure output directory exists
	if err := codegen.EnsureDir(OutputDir); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate messages
	if err := codegen.GenerateMessages(MessagesYAMLPath, OutputPath); err != nil {
		return fmt.Errorf("failed to generate messages: %w", err)
	}

	fmt.Printf("✅ Generated %s\n", OutputPath)
	return nil
}
