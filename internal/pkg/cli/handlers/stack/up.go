package stack

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/utils"
	cliTypes "github.com/otto-nation/otto-stack/internal/pkg/cli/types"
	"github.com/otto-nation/otto-stack/internal/pkg/compose"
	"github.com/otto-nation/otto-stack/internal/pkg/constants"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/otto-nation/otto-stack/internal/pkg/ui"
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
func (h *UpHandler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *cliTypes.BaseCommand) error {
	ui.Header(constants.MsgStarting)

	setup, cleanup, err := SetupCoreCommand(ctx, base)
	if err != nil {
		return err
	}
	defer cleanup()

	// Parse flags
	build, _ := cmd.Flags().GetBool("build")
	forceRecreate, _ := cmd.Flags().GetBool("force-recreate")

	options := types.StartOptions{
		Build:         build,
		ForceRecreate: forceRecreate,
		Detach:        true,
	}

	// Determine services to start
	serviceNames := args
	if len(serviceNames) == 0 {
		serviceNames = setup.Config.Stack.Enabled
	}

	// Filter services to only include container services
	serviceUtils := utils.NewServiceUtils()
	filteredServices, err := serviceUtils.ResolveServices(serviceNames)
	if err != nil {
		return fmt.Errorf("failed to resolve services: %w", err)
	}

	// Check for config changes
	configHash, err := h.getConfigHash(setup.Config)
	if err != nil {
		return fmt.Errorf("failed to calculate config hash: %w", err)
	}

	previousState, err := h.loadState()
	if err != nil {
		constants.SendMessage(constants.Message{Level: constants.LevelInfo, Content: "No previous state found, performing fresh setup"})
		previousState = &StackState{}
	}

	configChanged := previousState.ConfigHash != configHash
	if configChanged {
		constants.SendMessage(constants.Message{Level: constants.LevelInfo, Content: "Configuration changes detected, updating stack..."})

		// Clean up removed services
		if err := h.cleanupRemovedServices(ctx, setup, previousState.Services, filteredServices); err != nil {
			ui.Warning("Failed to clean up removed services: %v", err)
		}
	}

	// Generate compose file
	generator, err := compose.NewGenerator(setup.Config.Project.Name, constants.ServicesDir)
	if err != nil {
		return fmt.Errorf("failed to create compose generator: %w", err)
	}

	composeFile, err := generator.Generate(serviceNames)
	if err != nil {
		return fmt.Errorf("failed to generate compose file: %w", err)
	}

	// Ensure otto-stack directory exists
	if err := os.MkdirAll(constants.DevStackDir, 0755); err != nil {
		return fmt.Errorf("failed to create otto-stack directory: %w", err)
	}

	// Write compose file
	composeData, err := yaml.Marshal(composeFile)
	if err != nil {
		return fmt.Errorf("failed to marshal compose file: %w", err)
	}

	composePath := constants.DockerComposeFile
	if err := os.WriteFile(composePath, composeData, 0644); err != nil {
		return fmt.Errorf("failed to write compose file: %w", err)
	}

	// Start services
	if err := setup.DockerClient.Containers().Start(ctx, setup.Config.Project.Name, filteredServices, options); err != nil {
		return fmt.Errorf("failed to start services: %w", err)
	}

	// Save new state
	newState := &StackState{
		ConfigHash: configHash,
		Services:   filteredServices,
	}
	if err := h.saveState(newState); err != nil {
		ui.Warning("Failed to save state: %v", err)
	}

	ui.Success(constants.MsgStartSuccess)
	constants.SendMessage(constants.Message{Level: constants.LevelInfo, Content: "Run '%s' to check service status"}, constants.AppName+" status")
	return nil
}

// ValidateArgs validates the command arguments
func (h *UpHandler) ValidateArgs(args []string) error {
	return nil
}

// GetRequiredFlags returns required flags for this command
func (h *UpHandler) GetRequiredFlags() []string {
	return []string{}
}

// getConfigHash calculates hash of current config
func (h *UpHandler) getConfigHash(config *ProjectConfig) (string, error) {
	data, err := yaml.Marshal(config)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash), nil
}

// loadState loads previous stack state
func (h *UpHandler) loadState() (*StackState, error) {
	statePath := filepath.Join(constants.DevStackDir, constants.StateFileName)
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
	statePath := filepath.Join(constants.DevStackDir, constants.StateFileName)
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}
	return os.WriteFile(statePath, data, 0644)
}

// cleanupRemovedServices removes services no longer in config
func (h *UpHandler) cleanupRemovedServices(ctx context.Context, setup *CoreSetup, oldServices, newServices []string) error {
	removedServices := h.findRemovedServices(oldServices, newServices)
	if len(removedServices) == 0 {
		return nil
	}

	ui.Info("Removing services: %v", removedServices)
	return setup.DockerClient.Containers().Stop(ctx, setup.Config.Project.Name, removedServices, types.StopOptions{
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
