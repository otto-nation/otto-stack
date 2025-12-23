package stack

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"

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
	stateManager  *StateManager
	fileGenerator *services.FileGenerator
}

// NewUpHandler creates a new up handler
func NewUpHandler() *UpHandler {
	return &UpHandler{
		stateManager:  NewStateManager(),
		fileGenerator: services.NewFileGenerator(),
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

	// Check if verbose logging is enabled
	verbose, _ := cmd.Flags().GetBool("verbose")
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
	_ = timeoutSecs // Keep for potential future use

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
	base.Output.Info("%s", core.MsgProcess_generating_compose)
	if err := h.generateComposeFile(serviceNames, setup); err != nil {
		return err
	}

	// Generate .env.generated file
	base.Output.Info("%s", core.MsgProcess_generating_env)
	if err := h.generateEnvFile(filteredServices, setup.Config.Project.Name); err != nil {
		base.Output.Warning("Failed to generate .env file: %v", err)
	}

	// Check for port conflicts and resolve them if skip-conflicts is enabled
	if flags.SkipConflicts {
		filteredServices = h.skipConflictedServices(filteredServices, base)
	} else {
		if err := checkPortConflicts(filteredServices, base); err != nil {
			return err
		}
	}

	// Check if any services remain after conflict resolution
	if len(filteredServices) == 0 {
		base.Output.Warning("No services available to start")
		base.Output.Success("Stack operation completed (no services started)")
		return nil
	}

	// Execute Docker operations using new stack service
	base.Output.Info("Starting services: %v", filteredServices)

	// Create stack service
	stackService, err := NewStackService(verbose)
	if err != nil {
		return err
	}

	// Create start request (no characteristics for now)
	startRequest := services.StartRequest{
		Project:       setup.Config.Project.Name,
		Services:      filteredServices,
		Build:         flags.Build,
		ForceRecreate: flags.ForceRecreate,
	}

	if err := stackService.Start(ctx, startRequest); err != nil {
		return err
	}
	if verbose {
		base.Output.Info("Docker Compose completed successfully")
	}
	base.Output.Success("Services started successfully")

	// Run init containers after services are up
	base.Output.Info("Running initialization containers...")
	if err := h.runInitContainers(ctx, filteredServices, setup, base); err != nil {
		base.Output.Warning("Failed to run init containers: %v", err)
	}

	base.Output.Success(core.MsgStartSuccess)
	logger.Info(logger.LogMsgOperationCompleted, logger.LogFieldOperation, logger.OperationStackUp)
	return nil
}

// runInitContainers handles the initialization container execution
func (h *UpHandler) runInitContainers(ctx context.Context, filteredServices []string, setup *CoreSetup, base *base.BaseCommand) error {
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

	manager, err := GetServicesManager()
	if err != nil {
		return err
	}

	serviceConfigs := make(map[string]*services.ServiceConfig)
	for _, serviceName := range filteredServices {
		if config, err := manager.GetService(serviceName); err == nil {
			serviceConfigs[serviceName] = config
		}
	}

	if err := initManager.RunInitContainers(ctx, serviceConfigs, setup.Config.Project.Name); err != nil {
		return err
	}

	base.Output.Success("Initialization completed")
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

	// Generate compose file even in dry run mode
	if err := h.generateComposeFile(serviceNames, setup); err != nil {
		return err
	}
	base.Output.Success("Generated %s", docker.DockerComposeFileName)

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

	removedServices := findRemovedServices(oldServices, newServices)
	if len(removedServices) == 0 {
		return
	}

	base.Output.Info("Removing services no longer in configuration: %v", removedServices)

	// Use stack service for cleanup
	stackService, err := NewStackService(false)
	if err != nil {
		base.Output.Warning("Failed to create stack service: %v", err)
		return
	}

	// Create stop request for removed services
	stopRequest := services.StopRequest{
		Project:  setup.Config.Project.Name,
		Services: removedServices,
		Remove:   true, // Remove containers for cleanup
	}

	if err := stackService.Stop(ctx, stopRequest); err != nil {
		base.Output.Warning("Failed to remove services: %v", err)
	}
}

// findRemovedServices identifies services that were removed from configuration
func findRemovedServices(oldServices, newServices []string) []string {
	newServiceSet := make(map[string]bool)
	for _, service := range newServices {
		newServiceSet[service] = true
	}

	var removed []string
	for _, service := range oldServices {
		if !newServiceSet[service] {
			removed = append(removed, service)
		}
	}
	return removed
}

// skipConflictedServices filters out services with port conflicts
func (h *UpHandler) skipConflictedServices(services []string, base *base.BaseCommand) []string {
	conflicts := collectPortConflicts(services)
	if len(conflicts) == 0 {
		return services
	}

	conflictedServices := make(map[string]bool)
	for _, conflict := range conflicts {
		conflictedServices[conflict.ServiceName] = true
	}

	var availableServices []string
	for _, service := range services {
		if !conflictedServices[service] {
			availableServices = append(availableServices, service)
		}
	}

	base.Output.Warning("Skipping %d services due to port conflicts: %v", len(conflictedServices), getKeys(conflictedServices))
	base.Output.Info("Starting %d available services: %v", len(availableServices), availableServices)

	return availableServices
}

// getKeys returns keys from a map[string]bool
func getKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
