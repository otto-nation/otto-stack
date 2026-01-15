package lifecycle

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/common"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/utils"
)

// DownHandler handles the down command
type DownHandler struct {
	stateManager *common.StateManager
}

// NewDownHandler creates a new down handler
func NewDownHandler() *DownHandler {
	return &DownHandler{
		stateManager: common.NewStateManager(),
	}
}

// Handle executes the down command
func (h *DownHandler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand) error {
	return utils.ExecuteLifecycleCommand(ctx, cmd, args, base, core.CommandDown, h.stateManager)
}

// ValidateArgs validates the command arguments
func (h *DownHandler) ValidateArgs(args []string) error {
	return nil
}

// GetRequiredFlags returns required flags for this command
func (h *DownHandler) GetRequiredFlags() []string {
	return []string{}
}
