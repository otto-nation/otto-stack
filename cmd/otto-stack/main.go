package main

import (
	"fmt"
	"os"

	"github.com/otto-nation/otto-stack/internal/cli"
	"github.com/otto-nation/otto-stack/internal/pkg/logger"
)

func main() {
	// Initialize logger with default config
	if err := logger.Init(logger.DefaultConfig()); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Close()

	if err := cli.ExecuteFactory(); err != nil {
		os.Exit(1)
	}
}
