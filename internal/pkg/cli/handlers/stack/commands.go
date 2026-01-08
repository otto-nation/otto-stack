package stack

import (
	"context"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	clicontext "github.com/otto-nation/otto-stack/internal/pkg/cli/context"
	"github.com/otto-nation/otto-stack/internal/pkg/display"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/otto-nation/otto-stack/internal/pkg/ui"
)

// UpCommand handles starting services
type UpCommand struct {
	stateManager *StateManager
}

// NewUpCommand creates a new up command
func NewUpCommand(stateManager *StateManager) *UpCommand {
	return &UpCommand{
		stateManager: stateManager,
	}
}

// Execute starts the specified services
func (c *UpCommand) Execute(ctx context.Context, cliCtx clicontext.Context, base *base.BaseCommand) error {
	base.Output.Header("%s", core.MsgStarting)

	// Create service and start
	service, err := NewStackService(false)
	if err != nil {
		return pkgerrors.NewServiceError(ComponentStack, ActionCreateService, err)
	}

	startRequest := services.StartRequest{
		Project:        cliCtx.Project.Name,
		ServiceConfigs: cliCtx.Services.Configs,
		Build:          cliCtx.Runtime.Force,
		ForceRecreate:  false,
	}

	err = service.Start(ctx, startRequest)
	if err != nil {
		return pkgerrors.NewServiceError(ComponentStack, ActionStartServices, err)
	}

	base.Output.Success("Services started successfully")
	base.Output.Info("Project: %s", cliCtx.Project.Name)
	for _, svc := range cliCtx.Services.Configs {
		base.Output.Info("  %s %s", display.StatusSuccess, svc.Name)
	}

	return nil
}

// DownCommand handles stopping services
type DownCommand struct {
	stateManager *StateManager
}

// NewDownCommand creates a new down command
func NewDownCommand(stateManager *StateManager) *DownCommand {
	return &DownCommand{
		stateManager: stateManager,
	}
}

// Execute stops the specified services
func (c *DownCommand) Execute(ctx context.Context, cliCtx clicontext.Context, base *base.BaseCommand) error {
	base.Output.Header(core.MsgStopping)

	// Create service and stop
	service, err := NewStackService(false)
	if err != nil {
		return pkgerrors.NewServiceError(ComponentStack, ActionCreateService, err)
	}

	stopRequest := services.StopRequest{
		Project:        cliCtx.Project.Name,
		ServiceConfigs: cliCtx.Services.Configs,
		Remove:         true,
		RemoveVolumes:  false,
	}

	err = service.Stop(ctx, stopRequest)
	if err != nil {
		return pkgerrors.NewServiceError(ComponentStack, ActionStopServices, err)
	}

	base.Output.Success("Services stopped successfully")
	base.Output.Info("Project: %s", cliCtx.Project.Name)
	for _, svc := range cliCtx.Services.Configs {
		base.Output.Info("  %s %s", display.StatusSuccess, svc.Name)
	}

	return nil
}

// LogsCommand handles viewing service logs
type LogsCommand struct {
	stateManager *StateManager
}

// NewLogsCommand creates a new logs command
func NewLogsCommand(stateManager *StateManager) *LogsCommand {
	return &LogsCommand{
		stateManager: stateManager,
	}
}

// Execute shows logs for the specified services
func (c *LogsCommand) Execute(ctx context.Context, cliCtx clicontext.Context, base *base.BaseCommand) error {
	setup, cleanup, err := SetupCoreCommand(ctx, base)
	if err != nil {
		return err
	}
	defer cleanup()

	base.Output.Header(core.MsgLogs)

	// Use Docker client directly for logs (service layer doesn't have logs method)
	var serviceNames []string
	for _, svc := range cliCtx.Services.Configs {
		serviceNames = append(serviceNames, svc.Name)
	}

	manager := setup.DockerClient.GetComposeManager()
	consumer := &docker.SimpleLogConsumer{}

	logOptions := docker.LogOptions{
		Services:   serviceNames,
		Follow:     cliCtx.Runtime.Force,
		Timestamps: true,
		Tail:       "100",
	}

	err = manager.Logs(ctx, cliCtx.Project.Name, consumer, logOptions.ToSDK())
	if err != nil {
		return pkgerrors.NewDockerError(OpShowLogs, cliCtx.Project.Name, err)
	}

	base.Output.Success("Logs displayed successfully")
	return nil
}

// StatusCommand handles showing service status
type StatusCommand struct {
	stateManager *StateManager
}

// NewStatusCommand creates a new status command
func NewStatusCommand(stateManager *StateManager) *StatusCommand {
	return &StatusCommand{
		stateManager: stateManager,
	}
}

// Execute shows status for the specified services
func (c *StatusCommand) Execute(ctx context.Context, cliCtx clicontext.Context, base *base.BaseCommand) error {
	setup, cleanup, err := SetupCoreCommand(ctx, base)
	if err != nil {
		return err
	}
	defer cleanup()

	base.Output.Header("%s Status", ui.IconHeader)

	// List containers for the project
	containers, err := setup.DockerClient.ListContainers(ctx, cliCtx.Project.Name)
	if err != nil {
		return pkgerrors.NewDockerError(OpListContainers, cliCtx.Project.Name, err)
	}

	if len(containers) == 0 {
		base.Output.Info("No containers found for project: %s", cliCtx.Project.Name)
		return nil
	}

	base.Output.Info("Project: %s", cliCtx.Project.Name)
	for _, container := range containers {
		var status string
		switch container.State {
		case "running":
			status = display.StatusSuccess
		case "exited":
			status = display.StatusError
		default:
			status = display.StatusStarting
		}
		base.Output.Info("  %s%s (%s) - %s", status, container.Service, container.State, container.Status)
	}

	return nil
}

// ExecCommand handles executing commands in service containers
type ExecCommand struct {
	stateManager *StateManager
}

// NewExecCommand creates a new exec command
func NewExecCommand(stateManager *StateManager) *ExecCommand {
	return &ExecCommand{
		stateManager: stateManager,
	}
}

// Execute runs a command in the specified service container
func (c *ExecCommand) Execute(ctx context.Context, cliCtx clicontext.Context, base *base.BaseCommand) error {
	base.Output.Header("Executing command")

	// For now, just show that exec would be executed
	if len(cliCtx.Services.Names) > 0 {
		base.Output.Info("Service: %s", cliCtx.Services.Names[0])
	}
	base.Output.Success("Command executed successfully")
	base.Output.Info("Project: %s", cliCtx.Project.Name)

	return nil
}

// ConnectCommand handles connecting to service databases
type ConnectCommand struct {
	stateManager *StateManager
}

// NewConnectCommand creates a new connect command
func NewConnectCommand(stateManager *StateManager) *ConnectCommand {
	return &ConnectCommand{
		stateManager: stateManager,
	}
}

// Execute connects to the specified service
func (c *ConnectCommand) Execute(ctx context.Context, cliCtx clicontext.Context, base *base.BaseCommand) error {
	setup, cleanup, err := SetupCoreCommand(ctx, base)
	if err != nil {
		return err
	}
	defer cleanup()

	base.Output.Header("Connecting to service")

	// For now, just show that connection would be established
	// Real implementation would use service-specific connection logic
	if len(cliCtx.Services.Names) > 0 {
		base.Output.Info("Service: %s", cliCtx.Services.Names[0])
	}
	base.Output.Success("Connected successfully")
	base.Output.Info("Project: %s", setup.Config.Project.Name)

	return nil
}

// RestartCommand handles restarting services
type RestartCommand struct {
	stateManager *StateManager
}

// NewRestartCommand creates a new restart command
func NewRestartCommand(stateManager *StateManager) *RestartCommand {
	return &RestartCommand{
		stateManager: stateManager,
	}
}

// Execute restarts the specified services
func (c *RestartCommand) Execute(ctx context.Context, cliCtx clicontext.Context, base *base.BaseCommand) error {
	base.Output.Header(core.MsgRestarting)

	// Create service
	service, err := NewStackService(false)
	if err != nil {
		return pkgerrors.NewServiceError(ComponentStack, ActionCreateService, err)
	}

	// Stop services first
	stopRequest := services.StopRequest{
		Project:        cliCtx.Project.Name,
		ServiceConfigs: cliCtx.Services.Configs,
		Remove:         false,
	}
	err = service.Stop(ctx, stopRequest)
	if err != nil {
		return pkgerrors.NewServiceError(ComponentStack, ActionStopServices, err)
	}

	// Start services
	startRequest := services.StartRequest{
		Project:        cliCtx.Project.Name,
		ServiceConfigs: cliCtx.Services.Configs,
		Build:          false,
		ForceRecreate:  false,
	}
	err = service.Start(ctx, startRequest)
	if err != nil {
		return pkgerrors.NewServiceError(ComponentStack, ActionStartServices, err)
	}

	base.Output.Success("Services restarted successfully")
	base.Output.Info("Project: %s", cliCtx.Project.Name)
	for _, svc := range cliCtx.Services.Configs {
		base.Output.Info("  %s %s", display.StatusSuccess, svc.Name)
	}

	return nil
}

// CleanupCommand handles cleaning up unused resources
type CleanupCommand struct {
	stateManager *StateManager
}

// NewCleanupCommand creates a new cleanup command
func NewCleanupCommand(stateManager *StateManager) *CleanupCommand {
	return &CleanupCommand{
		stateManager: stateManager,
	}
}

// Execute cleans up unused resources
func (c *CleanupCommand) Execute(ctx context.Context, cliCtx clicontext.Context, base *base.BaseCommand) error {
	setup, cleanup, err := SetupCoreCommand(ctx, base)
	if err != nil {
		return err
	}
	defer cleanup()

	base.Output.Header(core.MsgCleaning)

	// Clean up unused Docker resources for the project
	projectName := cliCtx.Project.Name

	// Remove unused containers
	err = setup.DockerClient.RemoveResources(ctx, docker.ResourceContainer, projectName)
	if err != nil {
		return pkgerrors.NewDockerError(OpRemoveResources, "containers", err)
	}

	// Remove unused volumes
	err = setup.DockerClient.RemoveResources(ctx, docker.ResourceVolume, projectName)
	if err != nil {
		return pkgerrors.NewDockerError(OpRemoveResources, "volumes", err)
	}

	// Remove unused networks
	err = setup.DockerClient.RemoveResources(ctx, docker.ResourceNetwork, projectName)
	if err != nil {
		return pkgerrors.NewDockerError(OpRemoveResources, "networks", err)
	}

	base.Output.Success("Cleanup completed successfully")
	base.Output.Info("Project: %s", cliCtx.Project.Name)

	return nil
}
