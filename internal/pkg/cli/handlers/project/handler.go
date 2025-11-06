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
	enforcement *EnforcementHandler
	output      *ui.Output
	logger      *slog.Logger
}

// NewVersionHandler creates a new version handler
func NewVersionHandler() *Handler {
	return &Handler{
		enforcement: NewEnforcementHandler(nil), // Can handle nil
		output:      ui.NewOutput(),
		logger:      logger.GetLogger(),
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
	currentVersion := getCurrentVersion()
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

	currentVersion := getCurrentVersion()

	if flags.Full {
		return h.displayFullVersion(currentVersion, flags.Format)
	}

	return h.displayBasicVersion(currentVersion, flags.Format)
}

// displayBasicVersion displays basic version information
func (h *Handler) displayBasicVersion(version, format string) error {
	switch format {
	case "json":
		fmt.Printf(`{"version":"%s","app":"%s"}%s`, version, core.AppNameTitle, "\n")
	case "yaml":
		fmt.Printf("version: %s\napp: %s\n", version, core.AppNameTitle)
	default:
		fmt.Printf("%s version %s\n", core.AppNameTitle, version)
	}
	return nil
}

// displayFullVersion displays detailed version information
func (h *Handler) displayFullVersion(version, format string) error {
	buildInfo := getBuildInfo()
	h.logger.Info("Displaying full version", logger.LogFieldVersion, version, logger.LogFieldFormat, format, logger.LogFieldBuildInfo, buildInfo)

	switch format {
	case "json":
		fmt.Printf(`{"version":"%s","app":"%s","build":%s}%s`,
			version, core.AppNameTitle, buildInfo.JSON(), "\n")
	case "yaml":
		fmt.Printf("version: %s\napp: %s\n%s", version, core.AppNameTitle, buildInfo.YAML())
	default:
		h.output.Header("%s Version Information", core.AppNameTitle)
		h.output.Info(core.MsgVersion_version_label, version)
		h.output.Info(core.MsgVersion_build_date, buildInfo.Date)
		h.output.Info(core.MsgVersion_git_commit, buildInfo.Commit)
		h.output.Info(core.MsgVersion_go_version, buildInfo.GoVersion)
		h.output.Info(core.MsgVersion_platform, buildInfo.OS, buildInfo.Arch)
	}
	return nil
}

// BuildInfo contains build-time information
type BuildInfo struct {
	Date      string `json:"date" yaml:"date"`
	Commit    string `json:"commit" yaml:"commit"`
	GoVersion string `json:"go_version" yaml:"go_version"`
	OS        string `json:"os" yaml:"os"`
	Arch      string `json:"arch" yaml:"arch"`
}

// JSON returns JSON representation
func (b BuildInfo) JSON() string {
	return fmt.Sprintf(`{"date":"%s","commit":"%s","go_version":"%s","os":"%s","arch":"%s"}`,
		b.Date, b.Commit, b.GoVersion, b.OS, b.Arch)
}

// YAML returns YAML representation
func (b BuildInfo) YAML() string {
	return fmt.Sprintf("build:\n  date: %s\n  commit: %s\n  go_version: %s\n  os: %s\n  arch: %s\n",
		b.Date, b.Commit, b.GoVersion, b.OS, b.Arch)
}

// Helper functions

// getCurrentVersion returns the current version (set at build time)
func getCurrentVersion() string {
	// Use the main version package which has proper build-time injection
	return version.GetShortVersion()
}

// getBuildInfo returns build information (set at build time)
func getBuildInfo() BuildInfo {
	mainBuildInfo := version.GetBuildInfo()
	return BuildInfo{
		Date:      mainBuildInfo.BuildDate,
		Commit:    mainBuildInfo.GitCommit,
		GoVersion: mainBuildInfo.GoVersion,
		OS:        mainBuildInfo.Platform,
		Arch:      mainBuildInfo.Arch,
	}
}
