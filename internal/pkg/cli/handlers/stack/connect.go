package stack

import (
	"context"
	"fmt"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/ci"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/otto-nation/otto-stack/internal/pkg/validation"
	"github.com/spf13/cobra"
)

// ConnectHandler handles the connect command
type ConnectHandler struct{}

// NewConnectHandler creates a new connect handler
func NewConnectHandler() *ConnectHandler {
	return &ConnectHandler{}
}

// Handle executes the connect command
func (h *ConnectHandler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand) error {
	// Get CI-friendly flags
	ciFlags := ci.GetFlags(cmd)

	if len(args) < 1 {
		return fmt.Errorf("%s", core.Messages[core.MsgErrors_requires_service_name])
	}

	// Check initialization first
	if err := validation.CheckInitialization(); err != nil {
		return err
	}

	if !ciFlags.Quiet {
		base.Output.Header(core.MsgStack_connecting_to, args[0])
	}

	setup, cleanup, err := SetupCoreCommand(ctx, base)
	if err != nil {
		return err // Return error directly for clean output
	}
	defer cleanup()

	if !ciFlags.Quiet {
		base.Output.Header(core.MsgStack_connecting_to, args[0])
	}
	defer cleanup()

	serviceName := args[0]

	// Parse all flags with validation - single line!
	flags, err := core.ParseConnectFlags(cmd)
	if err != nil {
		return err
	}

	// Create connection command based on service type
	command, err := h.getConnectionCommand(serviceName, flags.Database, flags.User, flags.Host, flags.Port, flags.ReadOnly)
	if err != nil {
		ci.HandleError(ciFlags, err)
		return nil
	}

	// Execute connection command
	options := docker.ExecOptions{
		Interactive: true,
		TTY:         true,
	}

	dockerArgs := []string{"compose", "-f", docker.DockerComposeFile, "-p", setup.Config.Project.Name, "exec"}
	if options.User != "" {
		dockerArgs = append(dockerArgs, "--user", options.User)
	}
	if options.WorkingDir != "" {
		dockerArgs = append(dockerArgs, "--workdir", options.WorkingDir)
	}
	dockerArgs = append(dockerArgs, serviceName)
	dockerArgs = append(dockerArgs, command...)

	return setup.DockerClient.RunCommand(ctx, dockerArgs...)
}

// getConnectionCommand returns the appropriate connection command for the service
func (h *ConnectHandler) getConnectionCommand(serviceName, database, user, host string, port int, _ bool) ([]string, error) {
	manager, err := services.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create service manager: %w", err)
	}

	service, err := manager.GetServiceV2(serviceName)
	if err != nil {
		return nil, fmt.Errorf("unsupported service type: %s", serviceName)
	}

	if service.Service.Connection == nil {
		return nil, fmt.Errorf("service %s does not support connections", serviceName)
	}

	config := service.Service.Connection

	cmd := []string{config.Client}

	// Add user flag if specified or use default
	if config.UserFlag != "" {
		if user != "" {
			cmd = append(cmd, config.UserFlag, user)
		} else if config.DefaultUser != "" {
			cmd = append(cmd, config.UserFlag, config.DefaultUser)
		}
	}

	// Add database flag if specified and supported
	if config.DBFlag != "" && database != "" {
		cmd = append(cmd, config.DBFlag, database)
	}

	// Add host flag if specified and not localhost
	if config.HostFlag != "" && host != "" && host != services.ServiceLocalhost {
		cmd = append(cmd, config.HostFlag, host)
	}

	// Add port flag if specified
	if config.PortFlag != "" && port > 0 {
		cmd = append(cmd, config.PortFlag, fmt.Sprintf("%d", port))
	}

	// Add extra flags (like MySQL password prompt)
	cmd = append(cmd, config.ExtraFlags...)

	return cmd, nil
}

// ValidateArgs validates the command arguments
func (h *ConnectHandler) ValidateArgs(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("%s", core.Messages[core.MsgErrors_requires_service_name])
	}
	return nil
}

// GetRequiredFlags returns required flags for this command
func (h *ConnectHandler) GetRequiredFlags() []string {
	return []string{}
}
