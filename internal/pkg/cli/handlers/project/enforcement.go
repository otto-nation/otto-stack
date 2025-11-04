package project

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/utils"
	"github.com/otto-nation/otto-stack/internal/pkg/constants"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/otto-nation/otto-stack/internal/pkg/version"
	"github.com/spf13/cobra"
)

// EnforcementHandler handles version enforcement commands
type EnforcementHandler struct {
	enforcer *version.VersionEnforcer
	notifier *version.UpdateNotifier
}

// NewEnforcementHandler creates a new enforcement handler
func NewEnforcementHandler(manager version.VersionManager) *EnforcementHandler {
	policy := version.EnforcementPolicy{
		StrictMode:       version.DefaultStrictMode,
		AllowDrift:       version.DefaultAllowDrift,
		MaxDriftDuration: version.DefaultMaxDriftDuration,
		AutoSync:         version.DefaultAutoSync,
	}

	enforcer := version.NewVersionEnforcer(manager, policy)

	// Use main config file for notifications
	configPath := filepath.Join(constants.OttoStackDir, constants.ConfigFileName)
	notifier := version.NewUpdateNotifier(manager, configPath)

	return &EnforcementHandler{
		enforcer: enforcer,
		notifier: notifier,
	}
}

// HandleCheck handles the version check command
func (h *EnforcementHandler) HandleCheck(ctx context.Context, cmd *cobra.Command, args []string, base *types.BaseCommand) error {
	flags := utils.GetCIFlags(cmd)
	projectPath := "."

	if len(args) > 0 {
		projectPath = args[0]
	}

	result, err := h.enforcer.CheckCompliance(projectPath)
	if err != nil {
		utils.HandleError(flags, fmt.Errorf(constants.Messages[constants.MsgErrors_compliance_check_failed], err))
		return nil
	}

	if flags.JSON {
		return json.NewEncoder(os.Stdout).Encode(result)
	}

	h.displayComplianceResult(result, base)

	if result.ExitCode != constants.ExitSuccess {
		os.Exit(result.ExitCode)
	}

	return nil
}

// HandleEnforce handles the version enforce command
func (h *EnforcementHandler) HandleEnforce(ctx context.Context, cmd *cobra.Command, args []string, base *types.BaseCommand) error {
	flags := utils.GetCIFlags(cmd)
	projectPath := "."

	if len(args) > 0 {
		projectPath = args[0]
	}

	err := h.enforcer.EnforceCompliance(projectPath)
	if err != nil {
		utils.HandleError(flags, fmt.Errorf(constants.Messages[constants.MsgErrors_enforcement_failed], err))
		return nil
	}

	if !flags.Quiet {
		base.Output.Success("%s", constants.Messages[constants.MsgSuccess_version_compliance_enforced])
	}

	return nil
}

// HandleDrift handles the version drift command
func (h *EnforcementHandler) HandleDrift(ctx context.Context, cmd *cobra.Command, args []string, base *types.BaseCommand) error {
	flags := utils.GetCIFlags(cmd)

	drifts, err := h.enforcer.DetectAllDrift()
	if err != nil {
		utils.HandleError(flags, fmt.Errorf(constants.Messages[constants.MsgErrors_drift_detection_failed], err))
		return nil
	}

	if flags.JSON {
		return json.NewEncoder(os.Stdout).Encode(drifts)
	}

	h.displayDriftResults(drifts, base)
	return nil
}

// HandleNotify handles update notification commands
func (h *EnforcementHandler) HandleNotify(ctx context.Context, cmd *cobra.Command, args []string, base *types.BaseCommand) error {
	if len(args) == 0 {
		return h.handleNotifyCheck(cmd, base)
	}

	switch args[0] {
	case constants.Messages[constants.MsgCommands_check]:
		return h.handleNotifyCheck(cmd, base)
	case constants.Messages[constants.MsgCommands_config]:
		return h.handleNotifyConfig(cmd, args[1:], base)
	case constants.Messages[constants.MsgCommands_suppress]:
		return h.handleNotifySuppress(cmd, args[1:], base)
	default:
		return fmt.Errorf(constants.Messages[constants.MsgErrors_unknown_notify_command], args[0])
	}
}

func (h *EnforcementHandler) handleNotifyCheck(cmd *cobra.Command, base *types.BaseCommand) error {
	flags := utils.GetCIFlags(cmd)

	notification, err := h.notifier.CheckForUpdates()
	if err != nil {
		utils.HandleError(flags, fmt.Errorf(constants.Messages[constants.MsgErrors_update_check_failed], err))
		return nil
	}

	if flags.JSON {
		return json.NewEncoder(os.Stdout).Encode(notification)
	}

	if notification == nil {
		if !flags.Quiet {
			base.Output.Success("%s", constants.Messages[constants.MsgSuccess_no_updates_available])
		}
	} else {
		base.Output.Warning(constants.Messages[constants.MsgWarnings_update_available],
			notification.CurrentVersion, notification.LatestVersion, notification.Severity)
		base.Output.Info("   %s", notification.Message)
	}

	return nil
}

func (h *EnforcementHandler) handleNotifyConfig(cmd *cobra.Command, args []string, base *types.BaseCommand) error {
	flags := utils.GetCIFlags(cmd)

	if len(args) == 0 {
		// Show current config
		config := h.notifier.GetConfig()
		if flags.JSON {
			return json.NewEncoder(os.Stdout).Encode(config)
		}

		h.displayNotificationConfig(config, base)
		return nil
	}

	// Set config values
	config := h.notifier.GetConfig()

	for _, arg := range args {
		parts := strings.SplitN(arg, "=", constants.KeyValueParts)
		if len(parts) != constants.KeyValueParts {
			return fmt.Errorf(constants.Messages[constants.MsgErrors_invalid_config_format], arg)
		}

		key, value := parts[0], parts[1]

		switch key {
		case constants.Messages[constants.MsgConfig_keys_enabled]:
			config.Enabled = value == constants.BoolTrue
		case constants.Messages[constants.MsgConfig_keys_frequency]:
			config.Frequency = value
		case constants.Messages[constants.MsgConfig_keys_min_severity]:
			config.MinSeverity = value
		case constants.Messages[constants.MsgConfig_keys_auto_check]:
			config.AutoCheck = value == "true"
		case constants.Messages[constants.MsgConfig_keys_show_prerelease]:
			config.ShowPrerelease = value == "true"
		default:
			return fmt.Errorf(constants.Messages[constants.MsgErrors_unknown_config_key], key)
		}
	}

	err := h.notifier.SetConfig(config)
	if err != nil {
		utils.HandleError(flags, fmt.Errorf(constants.Messages[constants.MsgErrors_failed_save_config], err))
		return nil
	}

	if !flags.Quiet {
		base.Output.Success("%s", constants.Messages[constants.MsgSuccess_notification_config_updated])
	}

	return nil
}

func (h *EnforcementHandler) handleNotifySuppress(cmd *cobra.Command, args []string, base *types.BaseCommand) error {
	flags := utils.GetCIFlags(cmd)

	if len(args) == 0 {
		return fmt.Errorf("%s", constants.Messages[constants.MsgErrors_suppress_duration_required])
	}

	duration, err := parseDuration(args[0])
	if err != nil {
		utils.HandleError(flags, fmt.Errorf(constants.Messages[constants.MsgErrors_invalid_duration], err))
		return nil
	}

	err = h.notifier.SuppressNotifications(duration)
	if err != nil {
		utils.HandleError(flags, fmt.Errorf(constants.Messages[constants.MsgErrors_failed_suppress_notifications], err))
		return nil
	}

	if !flags.Quiet {
		base.Output.Success(constants.Messages[constants.MsgSuccess_notifications_suppressed], duration)
	}

	return nil
}

func (h *EnforcementHandler) displayComplianceResult(result *version.EnforcementResult, base *types.BaseCommand) {
	if result.Compliant {
		base.Output.Success("%s", constants.Messages[constants.MsgSuccess_version_compliance_satisfied])
		return
	}

	base.Output.Error("%s", constants.Messages[constants.MsgErrors_version_compliance_failed])
	base.Output.Error("   %s", result.Message)

	if result.Drift != nil {
		base.Output.Error(constants.Messages[constants.MsgDrift_required_version], result.Drift.RequiredVersion)
		base.Output.Error(constants.Messages[constants.MsgDrift_active_version], result.Drift.ActiveVersion)
		base.Output.Error(constants.Messages[constants.MsgDrift_drift_type], result.Drift.DriftType, result.Drift.Severity)
		base.Output.Error(constants.Messages[constants.MsgDrift_duration], result.Drift.DriftDuration)
	}

	base.Output.Error(constants.Messages[constants.MsgDrift_action], result.Action)
}

func (h *EnforcementHandler) displayDriftResults(drifts []version.DriftDetection, base *types.BaseCommand) {
	if len(drifts) == 0 {
		base.Output.Success("%s", constants.Messages[constants.MsgSuccess_no_version_drift])
		return
	}

	base.Output.Warning(constants.Messages[constants.MsgErrors_version_drift_detected], len(drifts))
	base.Output.Info("")

	for _, drift := range drifts {
		base.Output.Info(constants.Messages[constants.MsgDrift_project_path], drift.ProjectPath)
		base.Output.Info(constants.Messages[constants.MsgDrift_required_version], drift.RequiredVersion)
		base.Output.Info(constants.Messages[constants.MsgDrift_active_version], drift.ActiveVersion)
		base.Output.Info(constants.Messages[constants.MsgDrift_drift_type], drift.DriftType, drift.Severity)
		base.Output.Info(constants.Messages[constants.MsgDrift_duration], drift.DriftDuration)
		base.Output.Info("")
	}
}

func (h *EnforcementHandler) displayNotificationConfig(config version.NotificationConfig, base *types.BaseCommand) {
	base.Output.Info("%s", constants.Messages[constants.MsgNotifications_config_header])
	base.Output.Info(constants.Messages[constants.MsgNotifications_enabled], config.Enabled)
	base.Output.Info(constants.Messages[constants.MsgNotifications_frequency], config.Frequency)
	base.Output.Info(constants.Messages[constants.MsgNotifications_min_severity], config.MinSeverity)
	base.Output.Info(constants.Messages[constants.MsgNotifications_auto_check], config.AutoCheck)
	base.Output.Info(constants.Messages[constants.MsgNotifications_show_prerelease], config.ShowPrerelease)
	base.Output.Info(constants.Messages[constants.MsgNotifications_last_check], config.LastCheck.Format(time.RFC3339))

	if !config.SuppressedUntil.IsZero() {
		base.Output.Info(constants.Messages[constants.MsgNotifications_suppressed_until], config.SuppressedUntil.Format(time.RFC3339))
	}
}

func parseDuration(s string) (time.Duration, error) {
	// Handle simple formats like "1d", "2w", "3h"
	if len(s) < constants.KeyValueParts {
		return 0, fmt.Errorf("%s", constants.Messages[constants.MsgErrors_invalid_duration_format])
	}

	unit := s[len(s)-1:]
	valueStr := s[:len(s)-1]

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return time.ParseDuration(s) // Fall back to standard parsing
	}

	switch unit {
	case constants.Messages[constants.MsgDuration_units_days]:
		return time.Duration(value) * 24 * time.Hour, nil
	case constants.Messages[constants.MsgDuration_units_weeks]:
		return time.Duration(value) * 7 * 24 * time.Hour, nil
	case constants.Messages[constants.MsgDuration_units_hours]:
		return time.Duration(value) * time.Hour, nil
	case constants.Messages[constants.MsgDuration_units_minutes]:
		return time.Duration(value) * time.Minute, nil
	default:
		return time.ParseDuration(s)
	}
}
