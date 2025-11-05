package project

import (
	"context"
	"fmt"

	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/utils"
	"github.com/otto-nation/otto-stack/internal/pkg/config"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/spf13/cobra"
)

// ValidateHandler handles the validate command
type ValidateHandler struct{}

// NewValidateHandler creates a new validate handler
func NewValidateHandler() *ValidateHandler {
	return &ValidateHandler{}
}

// Handle executes the validate command
func (h *ValidateHandler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *types.BaseCommand) error {
	// Check initialization first
	if err := utils.CheckInitialization(); err != nil {
		return err
	}

	flags := utils.GetCIFlags(cmd)

	// Load configuration
	_, err := config.LoadConfig()
	if err != nil {
		utils.HandleError(flags, fmt.Errorf("failed to load configuration: %w", err))
		return nil
	}

	if !flags.Quiet {
		base.Output.Success("Configuration is valid")
	}

	return nil
}

// ValidateArgs validates the command arguments
func (h *ValidateHandler) ValidateArgs(args []string) error {
	return nil
}

// GetRequiredFlags returns required flags for this command
func (h *ValidateHandler) GetRequiredFlags() []string {
	return []string{}
}
