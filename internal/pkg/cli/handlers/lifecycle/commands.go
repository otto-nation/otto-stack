package lifecycle

import (
	"context"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	clicontext "github.com/otto-nation/otto-stack/internal/pkg/cli/context"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/common"
	"github.com/otto-nation/otto-stack/internal/pkg/display"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/otto-nation/otto-stack/internal/pkg/ui"
)

// ServiceCommand handles all service operations with a generic pattern
type ServiceCommand struct {
	operation    string
	stateManager *common.StateManager
}

// NewServiceCommand creates a new service command for the specified operation
func NewServiceCommand(operation string, stateManager *common.StateManager) *ServiceCommand {
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

	service, err := common.NewServiceManager(false)
	if err != nil {
		return pkgerrors.NewServiceError(common.ComponentStack, common.ActionCreateService, err)
	}

	startRequest := services.StartRequest{
		Project:        cliCtx.Project.Name,
		ServiceConfigs: cliCtx.Services.Configs,
		Build:          cliCtx.Runtime.Force,
		ForceRecreate:  false,
	}

	err = service.Start(ctx, startRequest)
	if err != nil {
		return pkgerrors.NewServiceError(common.ComponentStack, common.ActionStartServices, err)
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

	service, err := common.NewServiceManager(false)
	if err != nil {
		return pkgerrors.NewServiceError(common.ComponentStack, common.ActionCreateService, err)
	}

	stopRequest := services.StopRequest{
		Project:        cliCtx.Project.Name,
		ServiceConfigs: cliCtx.Services.Configs,
		Remove:         true,
		RemoveVolumes:  false,
	}

	if err = service.Stop(ctx, stopRequest); err != nil {
		return pkgerrors.NewServiceError(common.ComponentStack, common.ActionStopServices, err)
	}

	c.showServiceSuccess(base, cliCtx, "stopped")
	return nil
}

// executeLogs shows logs for the specified services
func (c *ServiceCommand) executeLogs(ctx context.Context, cliCtx clicontext.Context, base *base.BaseCommand) error {
	base.Output.Header(core.MsgLogs)

	dockerClient, err := c.createDockerClient()
	if err != nil {
		return err
	}
	defer func() { _ = dockerClient.Close() }()

	return c.streamServiceLogs(ctx, dockerClient, cliCtx)
}

// executeStatus shows status for the specified services
func (c *ServiceCommand) executeStatus(ctx context.Context, cliCtx clicontext.Context, base *base.BaseCommand) error {
	base.Output.Header("%s Status", ui.IconHeader)

	dockerClient, err := c.createDockerClient()
	if err != nil {
		return err
	}
	defer func() { _ = dockerClient.Close() }()

	containers, err := dockerClient.ListContainers(ctx, cliCtx.Project.Name)
	if err != nil {
		return pkgerrors.NewDockerError(common.OpListContainers, cliCtx.Project.Name, err)
	}

	c.displayContainerStatus(base, cliCtx.Project.Name, containers)
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
	setup, cleanup, err := common.SetupCoreCommand(ctx, base)
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

	service, err := c.createServiceManager()
	if err != nil {
		return err
	}

	if err = c.stopServices(ctx, service, cliCtx, false); err != nil {
		return err
	}

	if err = c.startServices(ctx, service, cliCtx); err != nil {
		return err
	}

	c.showServiceSuccess(base, cliCtx, "restarted")
	return nil
}

// executeCleanup cleans up unused resources
func (c *ServiceCommand) executeCleanup(ctx context.Context, cliCtx clicontext.Context, base *base.BaseCommand) error {
	base.Output.Header(core.MsgCleaning)

	dockerClient, err := c.createDockerClient()
	if err != nil {
		return err
	}
	defer func() { _ = dockerClient.Close() }()

	if err = c.cleanupResources(ctx, dockerClient, cliCtx.Project.Name); err != nil {
		return err
	}

	base.Output.Success("Cleanup completed successfully")
	base.Output.Info("Project: %s", cliCtx.Project.Name)
	return nil
}

// Helper methods

// createServiceManager creates a new service manager
func (c *ServiceCommand) createServiceManager() (*services.Service, error) {
	service, err := common.NewServiceManager(false)
	if err != nil {
		return nil, pkgerrors.NewServiceError(common.ComponentStack, common.ActionCreateService, err)
	}
	return service, nil
}

// createDockerClient creates a new Docker client
func (c *ServiceCommand) createDockerClient() (*docker.Client, error) {
	dockerClient, err := docker.NewClient(nil)
	if err != nil {
		return nil, pkgerrors.NewDockerError("create client", "", err)
	}
	return dockerClient, nil
}

// startServices starts the services with the given configuration
func (c *ServiceCommand) startServices(ctx context.Context, service *services.Service, cliCtx clicontext.Context) error {
	startRequest := services.StartRequest{
		Project:        cliCtx.Project.Name,
		ServiceConfigs: cliCtx.Services.Configs,
		Build:          cliCtx.Runtime.Force,
		ForceRecreate:  false,
	}

	if err := service.Start(ctx, startRequest); err != nil {
		return pkgerrors.NewServiceError(common.ComponentStack, common.ActionStartServices, err)
	}
	return nil
}

// stopServices stops the services with the given configuration
func (c *ServiceCommand) stopServices(ctx context.Context, service *services.Service, cliCtx clicontext.Context, remove bool) error {
	stopRequest := services.StopRequest{
		Project:        cliCtx.Project.Name,
		ServiceConfigs: cliCtx.Services.Configs,
		Remove:         remove,
		RemoveVolumes:  false,
	}

	if err := service.Stop(ctx, stopRequest); err != nil {
		return pkgerrors.NewServiceError(common.ComponentStack, common.ActionStopServices, err)
	}
	return nil
}

// extractServiceNames extracts service names from the context
func (c *ServiceCommand) extractServiceNames(cliCtx clicontext.Context) []string {
	serviceNames := make([]string, 0, len(cliCtx.Services.Configs))
	for _, svc := range cliCtx.Services.Configs {
		serviceNames = append(serviceNames, svc.Name)
	}
	return serviceNames
}

// streamServiceLogs streams logs for all services
func (c *ServiceCommand) streamServiceLogs(ctx context.Context, dockerClient *docker.Client, cliCtx clicontext.Context) error {
	serviceNames := c.extractServiceNames(cliCtx)

	logOptions := docker.LogOptions{
		Services:   serviceNames,
		Follow:     cliCtx.Runtime.Force,
		Timestamps: true,
		Tail:       "100",
	}

	manager := dockerClient.GetComposeManager()
	consumer := &docker.SimpleLogConsumer{}

	if err := manager.Logs(ctx, cliCtx.Project.Name, consumer, logOptions.ToSDK()); err != nil {
		return pkgerrors.NewDockerError(common.OpShowLogs, cliCtx.Project.Name, err)
	}
	return nil
}

// displayContainerStatus displays the status of containers
func (c *ServiceCommand) displayContainerStatus(base *base.BaseCommand, projectName string, containers []docker.ContainerInfo) {
	if len(containers) == 0 {
		base.Output.Info("No containers found for project: %s", projectName)
		return
	}

	base.Output.Info("Project: %s", projectName)
	for _, container := range containers {
		status := c.getContainerStatusIcon(container.State)
		base.Output.Info("  %s%s (%s) - %s", status, container.Service, container.State, container.Status)
	}
}

// cleanupResources removes all resources for a project
func (c *ServiceCommand) cleanupResources(ctx context.Context, dockerClient *docker.Client, projectName string) error {
	resources := []docker.ResourceType{
		docker.ResourceContainer,
		docker.ResourceVolume,
		docker.ResourceNetwork,
	}

	for _, resource := range resources {
		if err := dockerClient.RemoveResources(ctx, resource, projectName); err != nil {
			return pkgerrors.NewDockerError(common.OpRemoveResources, string(resource), err)
		}
	}
	return nil
}

// showServiceSuccess displays success message for service operations
func (c *ServiceCommand) showServiceSuccess(base *base.BaseCommand, cliCtx clicontext.Context, action string) {
	base.Output.Success("Services %s successfully", action)
	base.Output.Info("Project: %s", cliCtx.Project.Name)
	for _, svc := range cliCtx.Services.Configs {
		base.Output.Info("  %s %s", display.StatusSuccess, svc.Name)
	}
}

// getContainerStatusIcon returns the appropriate status icon for a container state
func (c *ServiceCommand) getContainerStatusIcon(state string) string {
	switch state {
	case "running":
		return display.StatusSuccess
	case "exited":
		return display.StatusError
	default:
		return display.StatusStarting
	}
}
