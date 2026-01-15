package base

import (
	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/ui"
	"github.com/spf13/cobra"
)

// NewBaseCommand creates a properly initialized BaseCommand from cobra command
func NewBaseCommand(cmd *cobra.Command) *BaseCommand {
	// Extract persistent global flags (available on all commands)
	quiet, _ := cmd.Flags().GetBool(core.FlagQuiet)

	// Extract conditional global flags (may not exist on all commands)
	noColor := false
	if flag := cmd.Flags().Lookup(core.FlagNoColor); flag != nil {
		noColor, _ = cmd.Flags().GetBool(core.FlagNoColor)
	}

	// Create output with proper configuration
	output := ui.NewOutput()
	output.Quiet = quiet
	output.NoColor = noColor

	return &BaseCommand{
		Output: output,
	}
}
