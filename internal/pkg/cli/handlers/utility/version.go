package utility

import (
	"context"
	"log/slog"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/logger"
	"github.com/otto-nation/otto-stack/internal/pkg/messages"
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
	flags, err := core.ParseVersionFlags(cmd)
	if err != nil {
		return pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, pkgerrors.FieldFlags, messages.ValidationFailedParseFlags, err)
	}

	if flags.CheckUpdates {
		return h.handleCheckUpdates(ctx, cmd, args, base)
	}

	return h.handleVersionDisplay(flags, base)
}

// handleCheckUpdates handles the --check-updates flag
func (h *Handler) handleCheckUpdates(_ context.Context, _ *cobra.Command, _ []string, base *base.BaseCommand) error {
	base.Output.Header(messages.VersionCheckingUpdates)

	currentVersion := h.versionDisplayManager.GetCurrentVersion()
	base.Output.Info(messages.VersionCurrentInfo, currentVersion)

	checker := version.NewUpdateChecker(currentVersion)
	release, hasUpdate, err := checker.CheckForUpdates()
	if err != nil {
		return pkgerrors.NewSystemError(pkgerrors.ErrCodeOperationFail, messages.ValidationFailedCheckUpdates, err)
	}

	if !hasUpdate {
		base.Output.Success(messages.SuccessLatestVersion)
		return nil
	}

	base.Output.Warning(messages.VersionUpdateAvailable, currentVersion, release.TagName)
	base.Output.Info(messages.VersionReleaseInfo, release.Name)
	base.Output.Info("")
	base.Output.Info(messages.VersionUpdateInfo)
	base.Output.Info(messages.VersionInstallScript, core.GitHubOrg, core.GitHubRepo)
	base.Output.Info(messages.VersionManualDownload, core.GitHubOrg, core.GitHubRepo)

	return nil
}

// handleVersionDisplay handles the default version display
func (h *Handler) handleVersionDisplay(flags *core.VersionFlags, _ *base.BaseCommand) error {
	currentVersion := h.versionDisplayManager.GetCurrentVersion()

	if flags.Full {
		return h.versionDisplayManager.DisplayFull(currentVersion, flags.Format)
	}

	return h.versionDisplayManager.DisplayBasic(currentVersion, flags.Format)
}
