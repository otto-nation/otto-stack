package project

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/logger"
	"github.com/otto-nation/otto-stack/internal/pkg/ui"
	"github.com/otto-nation/otto-stack/internal/pkg/version"
	"github.com/spf13/cobra"
)

// Handler handles version command with update functionality
type Handler struct {
	enforcement           *EnforcementHandler
	output                *ui.Output
	logger                *slog.Logger
	versionDisplayManager *VersionDisplayManager
}

// NewVersionHandler creates a new version handler
func NewVersionHandler() *Handler {
	return &Handler{
		enforcement:           NewEnforcementHandler(nil), // Can handle nil
		output:                ui.NewOutput(),
		logger:                logger.GetLogger(),
		versionDisplayManager: NewVersionDisplayManager(),
	}
}

// ValidateArgs validates version command arguments
func (h *Handler) ValidateArgs(args []string) error {
	// Version command doesn't require arguments
	return nil
}

// GetRequiredFlags returns required flags for version command
func (h *Handler) GetRequiredFlags() []string {
	return []string{} // No required flags
}

// Handle handles the version command with update checking
func (h *Handler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand) error {
	// Parse all flags with validation - single line!
	flags, err := core.ParseVersionFlags(cmd)
	if err != nil {
		return err
	}

	if flags.CheckUpdates {
		return h.handleCheckUpdates(ctx, cmd, args, base)
	}

	// Default version display behavior
	return h.handleVersionDisplay(ctx, cmd, args, base)
}

// handleCheckUpdates handles the --check-updates flag
func (h *Handler) handleCheckUpdates(_ context.Context, _ *cobra.Command, _ []string, _ *base.BaseCommand) error {
	h.output.Header("%s", core.MsgVersion_checking_updates)

	// Get current version (this should come from build-time ldflags)
	currentVersion := h.versionDisplayManager.GetCurrentVersion()
	h.output.Info(core.MsgVersion_current_info, currentVersion)

	// Check for updates
	checker := version.NewUpdateChecker(currentVersion)
	release, hasUpdate, err := checker.CheckForUpdates()
	if err != nil {
		return fmt.Errorf(core.MsgErrors_failed_check_updates, err)
	}

	if !hasUpdate {
		h.output.Success("%s", core.MsgSuccess_latest_version)
		return nil
	}

	h.output.Warning(core.MsgVersion_update_available, currentVersion, release.TagName)
	h.output.Info(core.MsgVersion_release_info, release.Name)

	// Show update instructions
	h.output.Info("")
	h.output.Info("%s", core.MsgVersion_update_info)
	h.output.Info(core.MsgVersion_install_script,
		core.GitHubOrg, core.GitHubRepo)
	h.output.Info(core.MsgVersion_manual_download,
		core.GitHubOrg, core.GitHubRepo)

	return nil
}

// handleVersionDisplay handles the default version display
func (h *Handler) handleVersionDisplay(_ context.Context, cmd *cobra.Command, _ []string, _ *base.BaseCommand) error {
	// Parse all flags with validation - single line!
	flags, err := core.ParseVersionFlags(cmd)
	if err != nil {
		return err
	}

	currentVersion := h.versionDisplayManager.GetCurrentVersion()

	if flags.Full {
		return h.versionDisplayManager.DisplayFull(currentVersion, flags.Format)
	}

	return h.versionDisplayManager.DisplayBasic(currentVersion, flags.Format)
}
