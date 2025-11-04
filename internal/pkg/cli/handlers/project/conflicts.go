package project

import (
	"context"

	"github.com/otto-nation/otto-stack/internal/pkg/constants"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/spf13/cobra"
)

// ConflictsHandler handles the conflicts command
type ConflictsHandler struct{}

// NewConflictsHandler creates a new conflicts handler
func NewConflictsHandler() *ConflictsHandler {
	return &ConflictsHandler{}
}

// Handle executes the conflicts command
func (h *ConflictsHandler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *types.BaseCommand) error {
	base.Output.Header("%s", constants.Messages[constants.MsgConflicts_header])
	base.Output.Info("%s", constants.Messages[constants.MsgConflicts_not_implemented])
	return nil
}

// ValidateArgs validates the command arguments
func (h *ConflictsHandler) ValidateArgs(args []string) error {
	return nil
}

// GetRequiredFlags returns required flags for this command
func (h *ConflictsHandler) GetRequiredFlags() []string {
	return []string{}
}
