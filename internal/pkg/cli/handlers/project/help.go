package project

import (
	"context"
	"fmt"

	"github.com/otto-nation/otto-stack/internal/pkg/cli/types"
	"github.com/spf13/cobra"
)

// HelpHandler handles help command
type HelpHandler struct{}

// NewHelpHandler creates a new help handler
func NewHelpHandler() *HelpHandler {
	return &HelpHandler{}
}

// ValidateArgs validates help command arguments
func (h *HelpHandler) ValidateArgs(args []string) error {
	return nil
}

// GetRequiredFlags returns required flags for help command
func (h *HelpHandler) GetRequiredFlags() []string {
	return []string{}
}

// Handle handles the help command
func (h *HelpHandler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *types.BaseCommand) error {
	if len(args) > 0 {
		// Show help for specific command
		return cmd.Root().Help()
	}

	// Show general help
	fmt.Println("otto-stack - Development stack management tool")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  otto-stack [command]")
	fmt.Println("")
	fmt.Println("Available Commands:")
	fmt.Println("  up        Start services")
	fmt.Println("  down      Stop services")
	fmt.Println("  status    Show service status")
	fmt.Println("  init      Initialize project")
	fmt.Println("  services  List available services")
	fmt.Println("")
	fmt.Println("Use \"otto-stack [command] --help\" for more information about a command.")

	return nil
}
