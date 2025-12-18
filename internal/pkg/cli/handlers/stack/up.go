package stack

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/ci"
	"github.com/otto-nation/otto-stack/internal/pkg/compose"
	"github.com/otto-nation/otto-stack/internal/pkg/env"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/logger"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/spf13/cobra"
)

const (
	// DefaultTimeoutSeconds is the default timeout for operations
	DefaultTimeoutSeconds = 30
)

// UpHandler handles the up command
type UpHandler struct {
	dockerOpsManager *DockerOperationsManager
	stateManager     *StateManager
	initManager      *InitContainerManager
}

// NewUpHandler creates a new up handler
func NewUpHandler() *UpHandler {
	return &UpHandler{
		dockerOpsManager: NewDockerOperationsManager(),
		stateManager:     NewStateManager(),
		initManager:      NewInitContainerManager(),
	}
}

// Handle executes the up command
func (h *UpHandler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand) error {
	// Check initialization first

	// Check initialization first, before any output
	setup, cleanup, err := SetupCoreCommand(ctx, base)
	if err != nil {
		return err // Return error directly without logging or headers
	}
	defer cleanup()

	// Start operation logging only after initialization check passes
	logger.Info(logger.LogMsgStartingOperation, logger.LogFieldOperation, logger.OperationStackUp, logger.LogFieldServices, args)

	ciFlags := ci.GetFlags(cmd)

	if ciFlags.DryRun {
		return h.handleDryRun(args, setup, base)
	}
	defer func() {
		if r := recover(); r != nil {
			logger.Error(logger.LogMsgOperationFailed, logger.LogFieldOperation, logger.OperationStackUp, logger.LogFieldError, fmt.Errorf("panic: %v", r))
			panic(r)
		}
	}()

	base.Output.Header("%s", core.MsgStarting)
	logger.Info(logger.LogMsgServiceAction, logger.LogFieldAction, logger.ActionStart, logger.LogFieldService, "stack", logger.LogFieldServices, args)

	// Parse all flags with validation
	flags, err := core.ParseUpFlags(cmd)
	if err != nil {
		logger.Error(logger.LogMsgOperationFailed, logger.LogFieldOperation, logger.OperationStackUp, logger.LogFieldError, err)
		return pkgerrors.NewValidationError("flags", ActionParseFlags, err)
	}

	// Parse timeout from string to duration
	timeoutSecs := h.parseTimeoutSeconds(flags.Timeout)

	// Clean usage with no repetitive error handling
	options := docker.StartOptions{
		Build:          flags.Build,
		ForceRecreate:  flags.ForceRecreate,
		Detach:         flags.Detach,
		Timeout:        time.Duration(timeoutSecs) * time.Second,
		NoDeps:         flags.NoDeps,
		ResolveDeps:    flags.ResolveDeps,
		CheckConflicts: flags.CheckConflicts,
		RemoveOrphans:  flags.ForceRecreate, // Auto-remove orphans when force recreating
	}

	// Determine and resolve services
	serviceNames, filteredServices, err := h.resolveServices(args, setup)
	if err != nil {
		return err
	}

	// Check for config changes
	configHash, err := h.stateManager.GetConfigHash(setup.Config)
	if err != nil {
		return pkgerrors.NewConfigError("", "failed to calculate config hash", err)
	}

	previousState, err := h.stateManager.LoadState()
	if err != nil {
		// Restart operation
		previousState = &StackState{}
	}

	configChanged := previousState.ConfigHash != configHash
	h.handleConfigChange(ctx, setup, previousState.Services, filteredServices, base, configChanged)

	// Generate and write compose file
	if err := h.generateComposeFile(serviceNames, setup); err != nil {
		return err
	}

	// Generate .env.generated file
	if err := h.generateEnvFile(filteredServices, setup.Config.Project.Name); err != nil {
		base.Output.Warning("Failed to generate .env file: %v", err)
	}

	// Execute Docker operations
	if err := h.dockerOpsManager.ExecuteOperations(ctx, setup, filteredServices, configHash, options, base); err != nil {
		return err
	}

	base.Output.Success(core.MsgStartSuccess)
	logger.Info(logger.LogMsgOperationCompleted, logger.LogFieldOperation, logger.OperationStackUp)
	return nil
}

// executeDockerOperations starts services and runs init containers

// generateComposeFile creates and writes the docker-compose file
func (h *UpHandler) generateComposeFile(serviceNames []string, setup *CoreSetup) error {
	manager, err := GetServicesManager()
	if err != nil {
		return pkgerrors.NewServiceError(ComponentServiceManager, ActionGetManager, err)
	}

	generator, err := compose.NewGenerator(setup.Config.Project.Name, services.ServicesDir, manager)
	if err != nil {
		return pkgerrors.NewServiceError(ComponentStack, ActionCreateGenerator, err)
	}

	composeData, err := generator.GenerateYAML(serviceNames)
	if err != nil {
		return pkgerrors.NewServiceError(ComponentStack, ActionGenerateCompose, err)
	}

	if err := os.MkdirAll(core.OttoStackDir, core.PermReadWriteExec); err != nil {
		return pkgerrors.NewServiceError(ComponentStack, ActionCreateDirectory, err)
	}

	composePath := docker.DockerComposeFilePath
	if err := os.MkdirAll(filepath.Dir(composePath), core.PermReadWriteExec); err != nil {
		return pkgerrors.NewServiceError(ComponentStack, ActionCreateDirectory, err)
	}

	return os.WriteFile(composePath, composeData, core.PermReadWrite)
}

// resolveServices determines and filters services to start
func (h *UpHandler) resolveServices(args []string, setup *CoreSetup) ([]string, []string, error) {
	serviceNames := args
	if len(serviceNames) == 0 {
		serviceNames = setup.Config.Stack.Enabled
	}

	serviceUtils := services.NewServiceUtils()
	filteredServices, err := serviceUtils.ResolveServices(serviceNames)
	if err != nil {
		return nil, nil, pkgerrors.NewServiceError(ComponentStack, ActionResolveServices, err)
	}

	return serviceNames, filteredServices, nil
}

// handleDryRun processes dry run mode
func (h *UpHandler) handleDryRun(args []string, setup *CoreSetup, base *base.BaseCommand) error { //nolint:unparam
	base.Output.Info("%s", core.MsgDry_run_showing_what_would_happen)

	serviceNames := args
	if len(serviceNames) == 0 {
		serviceNames = setup.Config.Stack.Enabled
	}

	base.Output.Info(core.MsgDry_run_would_start_services, fmt.Sprintf("%v", serviceNames))
	base.Output.Info(core.MsgDry_run_would_use_config, filepath.Join(core.OttoStackDir, core.ConfigFileName))
	return nil
}

// ValidateArgs validates the command arguments
func (h *UpHandler) ValidateArgs(args []string) error {
	// Service names are optional - if none provided, all enabled services are used
	return nil
}

// GetRequiredFlags returns required flags for this command
func (h *UpHandler) GetRequiredFlags() []string {
	// No flags are strictly required for the up command
	return []string{}
}

// getConfigHash calculates hash of current config

// loadState loads previous stack state

// saveState saves current stack state

// cleanupRemovedServices removes services no longer in config

// findRemovedServices compares old vs new service lists

// generateEnvFile generates .env.generated file with resolved services
func (h *UpHandler) generateEnvFile(services []string, projectName string) error {
	manager, err := GetServicesManager()
	if err != nil {
		return pkgerrors.NewServiceError(ComponentServiceManager, ActionGetManager, err)
	}

	envContent, err := env.Generate(projectName, services, manager)
	if err != nil {
		return pkgerrors.NewServiceError(ComponentStack, ActionGenerateEnv, err)
	}

	if err := os.MkdirAll(filepath.Dir(core.EnvGeneratedFilePath), core.PermReadWriteExec); err != nil {
		return pkgerrors.NewServiceError(ComponentStack, ActionCreateDirectory, err)
	}
	return os.WriteFile(core.EnvGeneratedFilePath, envContent, core.PermReadWrite)
}

// parseTimeoutSeconds parses timeout string or returns default
func (h *UpHandler) parseTimeoutSeconds(timeoutStr string) int {
	if timeoutStr == "" {
		return DefaultTimeoutSeconds
	}

	parsed, err := strconv.Atoi(timeoutStr)
	if err != nil {
		return DefaultTimeoutSeconds
	}
	return parsed
}

// handleConfigChange handles configuration changes and cleanup
func (h *UpHandler) handleConfigChange(ctx context.Context, setup *CoreSetup, oldServices, newServices []string, base *base.BaseCommand, configChanged bool) {
	if !configChanged {
		return
	}

	// Clean up removed services
	if err := h.dockerOpsManager.CleanupRemovedServices(ctx, setup, oldServices, newServices, base); err != nil {
		base.Output.Warning("Failed to clean up removed services: %v", err)
	}
}
