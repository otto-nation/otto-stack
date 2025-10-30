package stack

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/otto-nation/otto-stack/internal/core/docker"
	cliTypes "github.com/otto-nation/otto-stack/internal/pkg/cli/types"
	"github.com/otto-nation/otto-stack/internal/pkg/constants"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/otto-nation/otto-stack/internal/pkg/ui"
	"github.com/otto-nation/otto-stack/internal/pkg/utils"
	"github.com/spf13/cobra"
)

// RestartHandler handles the restart command
type RestartHandler struct{}

// NewRestartHandler creates a new restart handler
func NewRestartHandler() *RestartHandler {
	return &RestartHandler{}
}

// Handle executes the restart command
func (h *RestartHandler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *cliTypes.BaseCommand) error {
	ui.Header(constants.MsgRestarting)

	// Check if otto-stack is initialized
	configPath := filepath.Join(constants.DevStackDir, constants.ConfigFileName)
	if !utils.FileExists(configPath) {
		return errors.New(constants.ErrNotInitialized)
	}

	// Load project configuration
	cfg, err := LoadProjectConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Create Docker client
	logger := base.Logger.(loggerAdapter)
	dockerClient, err := docker.NewClient(logger.SlogLogger())
	if err != nil {
		return fmt.Errorf("failed to create Docker client: %w", err)
	}
	defer func() {
		if err := dockerClient.Close(); err != nil {
			base.Logger.Error("Failed to close Docker client", "error", err)
		}
	}()

	// Parse flags
	timeout, _ := cmd.Flags().GetInt("timeout")
	build, _ := cmd.Flags().GetBool("build")

	// Determine services to restart
	serviceNames := args
	if len(serviceNames) == 0 {
		serviceNames = cfg.Stack.Enabled
	}

	// Stop services first
	constants.SendMessage(constants.Message{Level: constants.LevelInfo, Content: "Stopping services..."})
	stopOptions := types.StopOptions{
		Timeout: timeout,
	}
	if err := dockerClient.Containers().Stop(ctx, cfg.Project.Name, serviceNames, stopOptions); err != nil {
		return fmt.Errorf("failed to stop services: %w", err)
	}

	// Start services
	constants.SendMessage(constants.Message{Level: constants.LevelInfo, Content: "Starting services..."})
	startOptions := types.StartOptions{
		Build: build,
	}
	if err := dockerClient.Containers().Start(ctx, cfg.Project.Name, serviceNames, startOptions); err != nil {
		return fmt.Errorf("failed to start services: %w", err)
	}

	ui.Success(constants.MsgRestartSuccess)
	constants.SendMessage(constants.Message{Level: constants.LevelInfo, Content: "Run '%s' to check service status"}, constants.AppName+" status")
	return nil
}

// ValidateArgs validates the command arguments
func (h *RestartHandler) ValidateArgs(args []string) error {
	return nil
}

// GetRequiredFlags returns required flags for this command
func (h *RestartHandler) GetRequiredFlags() []string {
	return []string{}
}
