package stack

import (
	"context"
	"crypto/sha256"
	"encoding/json"
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
	"github.com/otto-nation/otto-stack/internal/pkg/config"
	"github.com/otto-nation/otto-stack/internal/pkg/env"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/logger"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

const (
	// DefaultTimeoutSeconds is the default timeout for operations
	DefaultTimeoutSeconds = 30
)

// StackState tracks the current state of the stack
type StackState struct {
	ConfigHash string   `json:"config_hash"`
	Services   []string `json:"services"`
}

// UpHandler handles the up command
type UpHandler struct{}

// NewUpHandler creates a new up handler
func NewUpHandler() *UpHandler {
	return &UpHandler{}
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
	configHash, err := h.getConfigHash(setup.Config)
	if err != nil {
		return pkgerrors.NewConfigError("", "failed to calculate config hash", err)
	}

	previousState, err := h.loadState()
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
	if err := h.executeDockerOperations(ctx, setup, filteredServices, configHash, options, base); err != nil {
		return err
	}

	base.Output.Success(core.MsgStartSuccess)
	logger.Info(logger.LogMsgOperationCompleted, logger.LogFieldOperation, logger.OperationStackUp)
	return nil
}

// executeDockerOperations starts services and runs init containers
func (h *UpHandler) executeDockerOperations(ctx context.Context, setup *CoreSetup, filteredServices []string, configHash string, options docker.StartOptions, base *base.BaseCommand) error {
	// Start services first
	if err := setup.DockerClient.ComposeUp(ctx, setup.Config.Project.Name, filteredServices, options); err != nil {
		return pkgerrors.NewServiceError(ComponentStack, ActionStartServices, err)
	}

	// Run init containers after main services are started
	if err := h.runInitContainers(ctx, setup, filteredServices, base); err != nil {
		base.Output.Warning("Failed to run init containers: %v", err)
	}

	// Save new state
	newState := &StackState{
		ConfigHash: configHash,
		Services:   filteredServices,
	}
	if err := h.saveState(newState); err != nil {
		base.Output.Warning("Failed to save state: %v", err)
	}

	return nil
}

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
func (h *UpHandler) getConfigHash(config *config.Config) (string, error) {
	data, err := yaml.Marshal(config)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash), nil
}

// loadState loads previous stack state
func (h *UpHandler) loadState() (*StackState, error) {
	data, err := os.ReadFile(core.StateFilePath)
	if err != nil {
		return nil, err
	}

	var state StackState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}
	return &state, nil
}

// saveState saves current stack state
func (h *UpHandler) saveState(state *StackState) error {
	if err := os.MkdirAll(filepath.Dir(core.StateFilePath), core.PermReadWriteExec); err != nil {
		return pkgerrors.NewServiceError(ComponentStack, ActionCreateDirectory, err)
	}
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}
	return os.WriteFile(core.StateFilePath, data, core.PermReadWrite)
}

// cleanupRemovedServices removes services no longer in config
func (h *UpHandler) cleanupRemovedServices(ctx context.Context, setup *CoreSetup, oldServices, newServices []string, base *base.BaseCommand) error {
	removedServices := h.findRemovedServices(oldServices, newServices)
	if len(removedServices) == 0 {
		return nil
	}

	base.Output.Info(core.MsgStack_removing_services, fmt.Sprintf("%v", removedServices))
	return setup.DockerClient.ComposeDown(ctx, setup.Config.Project.Name, docker.StopOptions{
		Remove:        true,
		RemoveVolumes: true,
	})
}

// findRemovedServices compares old vs new service lists
func (h *UpHandler) findRemovedServices(oldServices, newServices []string) []string {
	newServiceMap := make(map[string]bool)
	for _, service := range newServices {
		newServiceMap[service] = true
	}

	var removed []string
	for _, service := range oldServices {
		if !newServiceMap[service] {
			removed = append(removed, service)
		}
	}
	return removed
}

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

// runInitContainers discovers and runs initialization containers
func (h *UpHandler) runInitContainers(ctx context.Context, setup *CoreSetup, resolvedServices []string, base *base.BaseCommand) error {
	initManager := NewInitContainerManager()
	return initManager.DiscoverAndRun(ctx, setup, resolvedServices, base)
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
	if err := h.cleanupRemovedServices(ctx, setup, oldServices, newServices, base); err != nil {
		base.Output.Warning("Failed to clean up removed services: %v", err)
	}
}
