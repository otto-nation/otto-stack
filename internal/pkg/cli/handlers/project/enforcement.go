package project

import (
	"context"

	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/constants"
	"github.com/otto-nation/otto-stack/internal/pkg/version"
	"github.com/spf13/cobra"
)

// EnforcementHandler handles simple version validation
type EnforcementHandler struct{}

// NewEnforcementHandler creates a new enforcement handler
func NewEnforcementHandler(_ any) *EnforcementHandler {
	return &EnforcementHandler{}
}

// HandleCheck handles basic version validation
func (h *EnforcementHandler) HandleCheck(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand) error {
	projectPath := "."
	if len(args) > 0 {
		projectPath = args[0]
	}

	err := version.ValidateProjectVersion(projectPath)
	if err != nil {
		base.Output.Error(constants.MsgErrors_version_compliance_failed)
		return err
	}

	base.Output.Success(constants.MsgSuccess_version_compliance_satisfied)
	return nil
}

// HandleEnforce handles version enforcement (simplified to just validation)
func (h *EnforcementHandler) HandleEnforce(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand) error {
	return h.HandleCheck(ctx, cmd, args, base)
}
