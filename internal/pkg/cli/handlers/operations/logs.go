package operations

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/common"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
)

// LogsHandler handles the logs command
type LogsHandler struct{}

// NewLogsHandler creates a new logs handler
func NewLogsHandler() *LogsHandler {
	return &LogsHandler{}
}

// Handle executes the logs command
func (h *LogsHandler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand) error {
	base.Output.Header(core.MsgLifecycle_logs)

	setup, cleanup, err := common.SetupCoreCommand(ctx, base)
	if err != nil {
		return err
	}
	defer cleanup()

	serviceConfigs, err := common.ResolveServiceConfigs(args, setup)
	if err != nil {
		return err
	}

	serviceNames := make([]string, 0, len(serviceConfigs))
	for _, svc := range serviceConfigs {
		serviceNames = append(serviceNames, svc.Name)
	}

	logOptions := docker.LogOptions{
		Services:   serviceNames,
		Follow:     true,
		Timestamps: true,
		Tail:       core.DefaultLogTailLines,
	}

	stackService, err := common.NewServiceManager(false)
	if err != nil {
		return pkgerrors.NewServiceError("stack", "create service", err)
	}

	manager := stackService.DockerClient.GetComposeManager()
	consumer := &docker.SimpleLogConsumer{}

	if err := manager.Logs(ctx, setup.Config.Project.Name, consumer, logOptions.ToSDK()); err != nil {
		return pkgerrors.NewDockerError("show logs", setup.Config.Project.Name, err)
	}

	return nil
}

// ValidateArgs validates the command arguments
func (h *LogsHandler) ValidateArgs(args []string) error {
	return nil
}

// GetRequiredFlags returns required flags for this command
func (h *LogsHandler) GetRequiredFlags() []string {
	return []string{}
}
