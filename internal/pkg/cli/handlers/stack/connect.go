package stack

import (
	"context"
	"fmt"

	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/utils"
	"github.com/otto-nation/otto-stack/internal/pkg/constants"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/spf13/cobra"
)

// ConnectHandler handles the connect command
type ConnectHandler struct{}

// NewConnectHandler creates a new connect handler
func NewConnectHandler() *ConnectHandler {
	return &ConnectHandler{}
}

// Handle executes the connect command
func (h *ConnectHandler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *types.BaseCommand) error {
	// Get CI-friendly flags
	ciFlags := utils.GetCIFlags(cmd)

	if len(args) < 1 {
		utils.HandleError(ciFlags, fmt.Errorf("%s", constants.Messages[constants.MsgErrors_requires_service_name]))
		return nil
	}

	if !ciFlags.Quiet {
		base.Output.Header(constants.Messages[constants.MsgStack_connecting_to], args[0])
	}

	setup, cleanup, err := SetupCoreCommand(ctx, base)
	if err != nil {
		utils.HandleError(ciFlags, err)
		return nil
	}
	defer cleanup()

	serviceName := args[0]

	// Parse all flags with validation - single line!
	flags, err := constants.ParseConnectFlags(cmd)
	if err != nil {
		return err
	}

	// Create connection command based on service type
	command, err := h.getConnectionCommand(serviceName, flags.Database, flags.User, flags.Host, flags.Port, flags.ReadOnly)
	if err != nil {
		utils.HandleError(ciFlags, err)
		return nil
	}

	// Execute connection command
	options := types.ExecOptions{
		Interactive: true,
		TTY:         true,
	}

	return setup.DockerClient.Containers().Exec(ctx, setup.Config.Project.Name, serviceName, command, options)
}

// getConnectionCommand returns the appropriate connection command for the service
func (h *ConnectHandler) getConnectionCommand(serviceName, database, user, host string, port int, _ bool) ([]string, error) {
	config, err := services.GetServiceConnectionConfig(serviceName)
	if err != nil {
		return nil, fmt.Errorf("unsupported service type: %s", serviceName)
	}

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
	if config.HostFlag != "" && host != "" && host != constants.ServiceLocalhost {
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
		return fmt.Errorf("%s", constants.Messages[constants.MsgErrors_requires_service_name])
	}
	return nil
}

// GetRequiredFlags returns required flags for this command
func (h *ConnectHandler) GetRequiredFlags() []string {
	return []string{}
}
