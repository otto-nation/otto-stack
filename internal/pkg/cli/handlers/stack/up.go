package stack

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/ci"
	"github.com/otto-nation/otto-stack/internal/pkg/compose"
	"github.com/otto-nation/otto-stack/internal/pkg/env"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/logger"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/otto-nation/otto-stack/internal/pkg/validation"
	"github.com/spf13/cobra"
)

const (
	// DefaultTimeoutSeconds is the default timeout for operations
	DefaultTimeoutSeconds = 30
)

// UpHandler handles the up command
type UpHandler struct {
	stateManager *StateManager
}

// NewUpHandler creates a new up handler
func NewUpHandler() *UpHandler {
	return &UpHandler{
		stateManager: NewStateManager(),
	}
}

// Handle executes the up command
func (h *UpHandler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand) error {
	setup, cleanup, err := h.prepare(ctx, cmd, args, base)
	if err != nil {
		return err
	}
	defer cleanup()

	serviceConfigs, err := h.resolveServices(args, setup, base, cmd)
	if err != nil {
		return err
	}

	if len(serviceConfigs) == 0 {
		base.Output.Warning("No services available to start")
		base.Output.Success("Stack operation completed (no services started)")
		return nil
	}

	if err := h.generateFiles(serviceConfigs, setup, base); err != nil {
		return err
	}

	if err := h.startServices(ctx, serviceConfigs, setup, base, cmd); err != nil {
		return err
	}

	h.startInitContainers(ctx, serviceConfigs, setup, base)
	base.Output.Success(core.MsgStartSuccess)
	logger.Info(logger.LogMsgOperationCompleted, logger.LogFieldOperation, logger.OperationStackUp)
	return nil
}

func (h *UpHandler) prepare(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand) (*CoreSetup, func(), error) {
	setup, cleanup, err := SetupCoreCommand(ctx, base)
	if err != nil {
		return nil, nil, err
	}

	logger.Info(logger.LogMsgStartingOperation, logger.LogFieldOperation, logger.OperationStackUp, logger.LogFieldServices, args)

	if ciFlags := ci.GetFlags(cmd); ciFlags.DryRun {
		return setup, cleanup, h.handleDryRun(args, setup, base)
	}

	base.Output.Header("%s", core.MsgStarting)
	return setup, cleanup, nil
}

func (h *UpHandler) resolveServices(args []string, setup *CoreSetup, base *base.BaseCommand, cmd *cobra.Command) ([]services.ServiceConfig, error) {
	flags, err := core.ParseUpFlags(cmd)
	if err != nil {
		return nil, pkgerrors.NewValidationError("flags", ActionParseFlags, err)
	}

	serviceConfigs, err := ResolveServiceConfigs(args, setup)
	if err != nil {
		return nil, err
	}

	if !flags.SkipConflicts {
		if err := h.checkPortConflictsForConfigs(serviceConfigs, base); err != nil {
			return nil, err
		}
	} else {
		serviceConfigs = h.skipConflictedServiceConfigs(serviceConfigs, base)
	}

	return serviceConfigs, nil
}

func (h *UpHandler) checkPortConflictsForConfigs(serviceConfigs []services.ServiceConfig, base *base.BaseCommand) error {
	return checkPortConflictsForConfigs(serviceConfigs, base)
}

func (h *UpHandler) skipConflictedServiceConfigs(serviceConfigs []services.ServiceConfig, base *base.BaseCommand) []services.ServiceConfig {
	conflicts := collectPortConflictsFromConfigs(serviceConfigs)
	if len(conflicts) == 0 {
		return serviceConfigs
	}

	conflictedServices := make(map[string]bool)
	for _, conflict := range conflicts {
		conflictedServices[conflict.ServiceName] = true
	}

	var availableConfigs []services.ServiceConfig
	var conflictedNames []string
	for _, config := range serviceConfigs {
		if conflictedServices[config.Name] {
			conflictedNames = append(conflictedNames, config.Name)
		} else {
			availableConfigs = append(availableConfigs, config)
		}
	}

	if len(conflictedNames) > 0 {
		base.Output.Warning("Skipping conflicted services: %v", conflictedNames)
	}
	return availableConfigs
}

func (h *UpHandler) generateFiles(serviceConfigs []services.ServiceConfig, setup *CoreSetup, base *base.BaseCommand) error {
	// State management
	configHash, err := h.stateManager.GetConfigHash(setup.Config)
	if err != nil {
		return err
	}
	previousState, _ := h.stateManager.LoadState()
	if previousState == nil {
		previousState = &StackState{}
	}

	configChanged := previousState.ConfigHash != configHash

	// TODO: MAJOR REFACTOR NEEDED - This entire state management section needs to be redesigned
	// Current issues:
	// 1. Mixed ServiceConfig/string-based state storage
	// 2. Incomplete config change detection
	// 3. No proper cleanup of removed services
	// 4. State format migration needed
	// Should be redesigned to:
	// - Store ServiceConfigs in state consistently
	// - Implement proper service lifecycle management
	// - Handle config changes with proper cleanup
	// - Migrate existing string-based state files
	_ = configChanged

	// Generate compose file
	base.Output.Info("%s", core.MsgProcess_generating_compose)
	generator, err := compose.NewGenerator(setup.Config.Project.Name, "", nil)
	if err != nil {
		return err
	}
	if err := generator.GenerateFromServiceConfigs(serviceConfigs, setup.Config.Project.Name); err != nil {
		return err
	}

	// Generate env file
	base.Output.Info("%s", core.MsgProcess_generating_env)
	if err := env.GenerateFile(setup.Config.Project.Name, serviceConfigs, core.EnvGeneratedFilePath); err != nil {
		base.Output.Warning("Failed to generate .env file: %v", err)
	}

	return nil
}

func (h *UpHandler) startServices(ctx context.Context, serviceConfigs []services.ServiceConfig, setup *CoreSetup, base *base.BaseCommand, cmd *cobra.Command) error {
	flags, _ := core.ParseUpFlags(cmd)
	verbose, _ := cmd.Flags().GetBool("verbose")

	stackService, err := NewStackService(verbose)
	if err != nil {
		return err
	}
	visibleServiceNames := services.ExtractVisibleServiceNames(serviceConfigs)
	base.Output.Info("Starting services: %v", visibleServiceNames)

	return stackService.Start(ctx, services.StartRequest{
		Project:        setup.Config.Project.Name,
		ServiceConfigs: serviceConfigs,
		Build:          flags.Build,
		ForceRecreate:  flags.ForceRecreate,
	})
}

func (h *UpHandler) startInitContainers(ctx context.Context, serviceConfigs []services.ServiceConfig, setup *CoreSetup, base *base.BaseCommand) {
	if err := h.runInitContainers(ctx, serviceConfigs, setup, base); err != nil {
		base.Output.Warning("Init containers failed: %v", err)
	}
	base.Output.Success("Services started successfully")
}

// runInitContainers handles the initialization container execution
func (h *UpHandler) runInitContainers(ctx context.Context, serviceConfigs []services.ServiceConfig, setup *CoreSetup, base *base.BaseCommand) error {
	client, err := docker.NewClient(slog.Default())
	if err != nil {
		return err
	}
	defer func() {
		_ = client.Close()
	}()

	initManager, err := NewServiceInitManager()
	if err != nil {
		return err
	}

	// Convert ServiceConfigs to map for init manager API compatibility
	// Note: Init manager API is already ServiceConfig-based, this conversion is just for the map format
	serviceConfigMap := make(map[string]*services.ServiceConfig)
	for _, config := range serviceConfigs {
		configCopy := config // Create a copy to avoid pointer issues
		serviceConfigMap[config.Name] = &configCopy
	}

	if err := initManager.RunInitContainers(ctx, serviceConfigMap, setup.Config.Project.Name); err != nil {
		return err
	}

	base.Output.Success("Initialization completed")
	return nil
}

// handleDryRun processes dry run mode
func (h *UpHandler) handleDryRun(args []string, setup *CoreSetup, base *base.BaseCommand) error { //nolint:unparam
	base.Output.Info("%s", core.MsgDry_run_showing_what_would_happen)

	// Resolve services for dry run
	serviceConfigs, err := ResolveServiceConfigs(args, setup)
	if err != nil {
		return err
	}

	// Generate compose file even in dry run mode
	generator, err := compose.NewGenerator(setup.Config.Project.Name, "", nil)
	if err != nil {
		return err
	}
	if err := generator.GenerateFromServiceConfigs(serviceConfigs, setup.Config.Project.Name); err != nil {
		return err
	}
	base.Output.Success("Generated %s", docker.DockerComposeFileName)

	// Display ServiceConfigs directly in dry run output
	serviceNames := services.ExtractServiceNames(serviceConfigs)
	base.Output.Info(core.MsgDry_run_would_start_services, fmt.Sprintf("%v", serviceNames))
	base.Output.Info(core.MsgDry_run_would_use_config, filepath.Join(core.OttoStackDir, core.ConfigFileName))
	return nil
}

// ValidateArgs validates the command arguments
func (h *UpHandler) ValidateArgs(args []string) error {
	return validation.ValidateUpArgs(args)
}

// GetRequiredFlags returns required flags for this command
func (h *UpHandler) GetRequiredFlags() []string {
	// No flags are strictly required for the up command
	return []string{}
}
