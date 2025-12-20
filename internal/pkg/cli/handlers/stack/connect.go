package stack

import (
	"context"
	"fmt"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/ci"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/spf13/cobra"
)

// ConnectHandler handles the connect command
type ConnectHandler struct{}

// NewConnectHandler creates a new connect handler
func NewConnectHandler() *ConnectHandler {
	return &ConnectHandler{}
}

// ValidateArgs validates the command arguments
func (h *ConnectHandler) ValidateArgs(args []string) error {
	if len(args) < 1 {
		return pkgerrors.NewValidationError("service-name", "service name is required", nil)
	}
	return nil
}

// GetRequiredFlags returns required flags for this command
func (h *ConnectHandler) GetRequiredFlags() []string {
	return []string{}
}

// Handle executes the connect command
func (h *ConnectHandler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand) error {
	if err := h.ValidateArgs(args); err != nil {
		return err
	}

	ciFlags := ci.GetFlags(cmd)
	serviceName := args[0]

	if !ciFlags.Quiet {
		base.Output.Header(core.MsgStack_connecting_to, serviceName)
	}

	setup, cleanup, err := SetupCoreCommand(ctx, base)
	if err != nil {
		return err
	}
	defer cleanup()

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

	return docker.NewComposeBuilder().
		Project(setup.Config.Project.Name).
		File(docker.DockerComposeFilePath).
		Exec(serviceName, command...).
		Run()
}

// getConnectionCommand returns the appropriate connection command for the service
func (h *ConnectHandler) getConnectionCommand(serviceName, database, user, host string, port int, _ bool) ([]string, error) {
	manager, err := GetServicesManager()
	if err != nil {
		return nil, pkgerrors.NewServiceError(ComponentServiceManager, ActionCreateManager, err)
	}

	service, err := manager.GetService(serviceName)
	if err != nil {
		return nil, pkgerrors.NewValidationErrorf(pkgerrors.FieldServiceName, "unsupported service type: %s", serviceName)
	}

	if service.Service.Connection == nil {
		return nil, pkgerrors.NewValidationErrorf(pkgerrors.FieldServiceName, "service %s does not support connections", serviceName)
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
