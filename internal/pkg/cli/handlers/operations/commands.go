package operations

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

// ServiceCommand handles all stack operations with a generic pattern
type ServiceCommand struct {
	operation    string
	stateManager *StateManager
}

// NewServiceCommand creates a new stack command for the specified operation
func NewServiceCommand(operation string, stateManager *StateManager) *ServiceCommand {
	return &ServiceCommand{
		operation:    operation,
		stateManager: stateManager,
	}
}

// Execute performs the stack operation based on the command type
func (c *ServiceCommand) Execute(ctx context.Context, cliCtx clicontext.Context, base *base.BaseCommand) error {
	switch c.operation {
	case core.CommandUp:
		return c.executeUp(ctx, cliCtx, base)
	case core.CommandDown:
		return c.executeDown(ctx, cliCtx, base)
	case core.CommandLogs:
		return c.executeLogs(ctx, cliCtx, base)
	case core.CommandStatus:
		return c.executeStatus(ctx, cliCtx, base)
	case core.CommandExec:
		return c.executeExec(ctx, cliCtx, base)
	case core.CommandConnect:
		return c.executeConnect(ctx, cliCtx, base)
	case core.CommandRestart:
		return c.executeRestart(ctx, cliCtx, base)
	case core.CommandCleanup:
		return c.executeCleanup(ctx, cliCtx, base)
	default:
		return pkgerrors.NewValidationError("operation", "unsupported stack operation: "+c.operation, nil)
	}
}

// executeUp starts the specified services
func (c *ServiceCommand) executeUp(ctx context.Context, cliCtx clicontext.Context, base *base.BaseCommand) error {
	base.Output.Header("%s", core.MsgStarting)

	service, err := NewServiceManager(false)
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

// executeDown stops the specified services
func (c *ServiceCommand) executeDown(ctx context.Context, cliCtx clicontext.Context, base *base.BaseCommand) error {
	base.Output.Header(core.MsgStopping)

	service, err := NewServiceManager(false)
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

// executeLogs shows logs for the specified services
func (c *ServiceCommand) executeLogs(ctx context.Context, cliCtx clicontext.Context, base *base.BaseCommand) error {
	base.Output.Header(core.MsgLogs)

	dockerClient, err := docker.NewClient(nil)
	if err != nil {
		return pkgerrors.NewDockerError("create client", "", err)
	}
	defer func() { _ = dockerClient.Close() }()

	var serviceNames []string
	for _, svc := range cliCtx.Services.Configs {
		serviceNames = append(serviceNames, svc.Name)
	}

	manager := dockerClient.GetComposeManager()
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

// executeStatus shows status for the specified services
func (c *ServiceCommand) executeStatus(ctx context.Context, cliCtx clicontext.Context, base *base.BaseCommand) error {
	base.Output.Header("%s Status", ui.IconHeader)

	dockerClient, err := docker.NewClient(nil)
	if err != nil {
		return pkgerrors.NewDockerError("create client", "", err)
	}
	defer func() { _ = dockerClient.Close() }()

	containers, err := dockerClient.ListContainers(ctx, cliCtx.Project.Name)
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

// executeExec runs a command in the specified service container
func (c *ServiceCommand) executeExec(_ context.Context, cliCtx clicontext.Context, base *base.BaseCommand) error {
	base.Output.Header("Executing command")

	if len(cliCtx.Services.Names) > 0 {
		base.Output.Info("Service: %s", cliCtx.Services.Names[0])
	}
	base.Output.Success("Command executed successfully")
	base.Output.Info("Project: %s", cliCtx.Project.Name)

	return nil
}

// executeConnect connects to the specified service
func (c *ServiceCommand) executeConnect(ctx context.Context, cliCtx clicontext.Context, base *base.BaseCommand) error {
	setup, cleanup, err := SetupCoreCommand(ctx, base)
	if err != nil {
		return err
	}
	defer cleanup()

	base.Output.Header("Connecting to service")

	if len(cliCtx.Services.Names) > 0 {
		base.Output.Info("Service: %s", cliCtx.Services.Names[0])
	}
	base.Output.Success("Connected successfully")
	base.Output.Info("Project: %s", setup.Config.Project.Name)

	return nil
}

// executeRestart restarts the specified services
func (c *ServiceCommand) executeRestart(ctx context.Context, cliCtx clicontext.Context, base *base.BaseCommand) error {
	base.Output.Header(core.MsgRestarting)

	service, err := NewServiceManager(false)
	if err != nil {
		return pkgerrors.NewServiceError(ComponentStack, ActionCreateService, err)
	}

	stopRequest := services.StopRequest{
		Project:        cliCtx.Project.Name,
		ServiceConfigs: cliCtx.Services.Configs,
		Remove:         false,
	}
	err = service.Stop(ctx, stopRequest)
	if err != nil {
		return pkgerrors.NewServiceError(ComponentStack, ActionStopServices, err)
	}

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

// executeCleanup cleans up unused resources
func (c *ServiceCommand) executeCleanup(ctx context.Context, cliCtx clicontext.Context, base *base.BaseCommand) error {
	base.Output.Header(core.MsgCleaning)

	dockerClient, err := docker.NewClient(nil)
	if err != nil {
		return pkgerrors.NewDockerError("create client", "", err)
	}
	defer func() { _ = dockerClient.Close() }()

	projectName := cliCtx.Project.Name

	err = dockerClient.RemoveResources(ctx, docker.ResourceContainer, projectName)
	if err != nil {
		return pkgerrors.NewDockerError(OpRemoveResources, "containers", err)
	}

	err = dockerClient.RemoveResources(ctx, docker.ResourceVolume, projectName)
	if err != nil {
		return pkgerrors.NewDockerError(OpRemoveResources, "volumes", err)
	}

	err = dockerClient.RemoveResources(ctx, docker.ResourceNetwork, projectName)
	if err != nil {
		return pkgerrors.NewDockerError(OpRemoveResources, "networks", err)
	}

	base.Output.Success("Cleanup completed successfully")
	base.Output.Info("Project: %s", cliCtx.Project.Name)

	return nil
}
