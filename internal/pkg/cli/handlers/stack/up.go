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
	"github.com/otto-nation/otto-stack/internal/pkg/logger"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
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
		base.Output.Info("%s", core.MsgDry_run_showing_what_would_happen)
		base.Output.Info(core.MsgDry_run_would_start_services, fmt.Sprintf("%v", args))
		base.Output.Info(core.MsgDry_run_would_use_config, filepath.Join(core.OttoStackDir, core.ConfigFileName))
		return nil
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
		return err
	}

	// Parse timeout from string to duration
	timeoutSecs := 30 // default
	if flags.Timeout != "" {
		if parsed, err := strconv.Atoi(flags.Timeout); err == nil {
			timeoutSecs = parsed
		}
	}

	// Clean usage with no repetitive error handling
	options := docker.StartOptions{
		Build:          flags.Build,
		ForceRecreate:  flags.ForceRecreate,
		Detach:         flags.Detach,
		Timeout:        time.Duration(timeoutSecs) * time.Second,
		NoDeps:         flags.NoDeps,
		ResolveDeps:    flags.ResolveDeps,
		CheckConflicts: flags.CheckConflicts,
	}

	// Determine services to start
	serviceNames := args
	if len(serviceNames) == 0 {
		serviceNames = setup.Config.Stack.Enabled
	}

	// Filter services to only include container services
	serviceUtils := services.NewServiceUtils()
	filteredServices, err := serviceUtils.ResolveServices(serviceNames)
	if err != nil {
		return fmt.Errorf(core.MsgStack_failed_resolve_services, err)
	}

	// Check for config changes
	configHash, err := h.getConfigHash(setup.Config)
	if err != nil {
		return fmt.Errorf(core.MsgStack_failed_calculate_hash, err)
	}

	previousState, err := h.loadState()
	if err != nil {
		// Restart operation
		previousState = &StackState{}
	}

	configChanged := previousState.ConfigHash != configHash
	if configChanged {
		// Restart operation

		// Clean up removed services
		if err := h.cleanupRemovedServices(ctx, setup, previousState.Services, filteredServices, base); err != nil {
			base.Output.Warning("Failed to clean up removed services: %v", err)
		}
	}

	// Generate compose file
	generator, err := compose.NewGenerator(setup.Config.Project.Name, services.ServicesDir)
	if err != nil {
		return fmt.Errorf(core.MsgStack_failed_create_generator, err)
	}

	composeFile, err := generator.Generate(serviceNames)
	if err != nil {
		return fmt.Errorf(core.MsgStack_failed_generate_compose, err)
	}

	// Ensure otto-stack directory exists
	if err := os.MkdirAll(core.OttoStackDir, core.DirPermReadWriteExec); err != nil {
		return fmt.Errorf(core.MsgStack_failed_create_directory, err)
	}

	// Write compose file
	composeData, err := yaml.Marshal(composeFile)
	if err != nil {
		return fmt.Errorf(core.MsgStack_failed_marshal_compose, err)
	}

	composePath := docker.DockerComposeFilePath
	if err := os.WriteFile(composePath, composeData, core.FilePermReadWrite); err != nil {
		return fmt.Errorf(core.MsgStack_failed_write_compose, err)
	}

	// Start services
	if err := setup.DockerClient.ComposeUp(ctx, setup.Config.Project.Name, filteredServices, options); err != nil {
		return fmt.Errorf(core.MsgStack_failed_start_services, err)
	}

	// Save new state
	newState := &StackState{
		ConfigHash: configHash,
		Services:   filteredServices,
	}
	if err := h.saveState(newState); err != nil {
		base.Output.Warning("Failed to save state: %v", err)
	}

	base.Output.Success(core.MsgStartSuccess)
	logger.Info(logger.LogMsgOperationCompleted, logger.LogFieldOperation, logger.OperationStackUp)
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
	statePath := filepath.Join(core.OttoStackDir, core.StateFileName)
	data, err := os.ReadFile(statePath)
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
	statePath := filepath.Join(core.OttoStackDir, core.StateFileName)
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}
	return os.WriteFile(statePath, data, core.FilePermReadWrite)
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
