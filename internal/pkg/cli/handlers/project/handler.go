package project

import (
	"context"
	"fmt"

	"github.com/otto-nation/otto-stack/internal/pkg/cli/types"
	"github.com/otto-nation/otto-stack/internal/pkg/constants"
	"github.com/otto-nation/otto-stack/internal/pkg/ui"
	"github.com/otto-nation/otto-stack/internal/pkg/version"
	"github.com/spf13/cobra"
)

// Handler handles version command with update functionality
type Handler struct {
	enforcement *EnforcementHandler
	output      *ui.Output
}

// NewVersionHandler creates a new version handler
func NewVersionHandler() *Handler {
	return &Handler{
		enforcement: NewEnforcementHandler(nil), // Can handle nil
		output:      ui.NewOutput(),
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
func (h *Handler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *types.BaseCommand) error {
	// Check if --check-updates flag is set
	checkUpdates, _ := cmd.Flags().GetBool("check-updates")

	if checkUpdates {
		return h.handleCheckUpdates(ctx, cmd, args, base)
	}

	// Default version display behavior
	return h.handleVersionDisplay(ctx, cmd, args, base)
}

// handleCheckUpdates handles the --check-updates flag
func (h *Handler) handleCheckUpdates(ctx context.Context, cmd *cobra.Command, args []string, base *types.BaseCommand) error {
	h.output.Header("🔍 Checking for Updates")

	// Get current version (this should come from build-time ldflags)
	currentVersion := getCurrentVersion()
	h.output.Info("Current version: %s", currentVersion)

	// Check for updates
	checker := version.NewUpdateChecker(currentVersion)
	release, hasUpdate, err := checker.CheckForUpdates()
	if err != nil {
		return fmt.Errorf("failed to check for updates: %w", err)
	}

	if !hasUpdate {
		h.output.Success("✅ You are running the latest version")
		return nil
	}

	h.output.Warning("⚠️  Update available: %s → %s", currentVersion, release.TagName)
	h.output.Info("Release: %s", release.Name)

	// Show update instructions
	h.output.Info("")
	h.output.Info("To update:")
	h.output.Info("  • Using install script: curl -fsSL https://raw.githubusercontent.com/%s/%s/main/scripts/install.sh | bash",
		constants.GitHubOrg, constants.GitHubRepo)
	if isBrewInstalled() {
		h.output.Info("  • Using Homebrew: brew upgrade %s", constants.AppName)
	}
	h.output.Info("  • Manual download: https://github.com/%s/%s/releases/latest",
		constants.GitHubOrg, constants.GitHubRepo)

	return nil
}

// handleVersionDisplay handles the default version display
func (h *Handler) handleVersionDisplay(ctx context.Context, cmd *cobra.Command, args []string, base *types.BaseCommand) error {
	full, _ := cmd.Flags().GetBool("full")
	format, _ := cmd.Flags().GetString("format")

	currentVersion := getCurrentVersion()

	if full {
		return h.displayFullVersion(currentVersion, format)
	}

	return h.displayBasicVersion(currentVersion, format)
}

// displayBasicVersion displays basic version information
func (h *Handler) displayBasicVersion(version, format string) error {
	switch format {
	case "json":
		fmt.Printf(`{"version":"%s","app":"%s"}%s`, version, constants.AppNameTitle, "\n")
	case "yaml":
		fmt.Printf("version: %s\napp: %s\n", version, constants.AppNameTitle)
	default:
		fmt.Printf("%s version %s\n", constants.AppNameTitle, version)
	}
	return nil
}

// displayFullVersion displays detailed version information
func (h *Handler) displayFullVersion(version, format string) error {
	buildInfo := getBuildInfo()

	switch format {
	case "json":
		fmt.Printf(`{"version":"%s","app":"%s","build":%s}%s`,
			version, constants.AppNameTitle, buildInfo.JSON(), "\n")
	case "yaml":
		fmt.Printf("version: %s\napp: %s\n%s", version, constants.AppNameTitle, buildInfo.YAML())
	default:
		h.output.Header("%s Version Information", constants.AppNameTitle)
		h.output.Info("Version: %s", version)
		h.output.Info("Build Date: %s", buildInfo.Date)
		h.output.Info("Git Commit: %s", buildInfo.Commit)
		h.output.Info("Go Version: %s", buildInfo.GoVersion)
		h.output.Info("Platform: %s/%s", buildInfo.OS, buildInfo.Arch)
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

// isBrewInstalled checks if Homebrew is available
func isBrewInstalled() bool {
	// Simple check - could be enhanced
	return false // Placeholder
}
