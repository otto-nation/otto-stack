package stack

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/command"
	clicontext "github.com/otto-nation/otto-stack/internal/pkg/cli/context"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
)

// ConnectHandler handles the connect command
type ConnectHandler struct {
	stateManager *StateManager
}

// NewConnectHandler creates a new connect handler
func NewConnectHandler() *ConnectHandler {
	return &ConnectHandler{
		stateManager: NewStateManager(),
	}
}

// ValidateArgs validates the command arguments
func (h *ConnectHandler) ValidateArgs(args []string) error {
	if len(args) < 1 {
		return pkgerrors.NewValidationError(pkgerrors.FieldServiceName, "service name is required", nil)
	}
	return nil
}

// GetRequiredFlags returns required flags for this command
func (h *ConnectHandler) GetRequiredFlags() []string {
	return []string{}
}

// Handle executes the connect command
func (h *ConnectHandler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand) error {
	// Create command and middleware chain
	connectCommand := NewConnectCommand(h.stateManager)
	validationMiddleware, loggingMiddleware := CreateStandardMiddlewareChain()

	handler := command.NewHandler(connectCommand, loggingMiddleware, validationMiddleware)

	// For now, create empty context - will be enhanced with flag processing
	cliCtx := clicontext.Context{}

	// Execute through command pattern
	return handler.Execute(ctx, cliCtx, base)
}

// getConnectionCommand - legacy method for tests
//
//nolint:unused,unparam
func (h *ConnectHandler) getConnectionCommand(serviceName, database, user, host string, port int, _ bool) ([]string, error) {
	// Simplified implementation for test compatibility
	switch serviceName {
	case "postgres":
		cmd := []string{"psql", "-U"}
		if user != "" {
			cmd = append(cmd, user)
		} else {
			cmd = append(cmd, "postgres")
		}
		if database != "" {
			cmd = append(cmd, "-d", database)
		}
		return cmd, nil
	default:
		return nil, pkgerrors.NewValidationError(pkgerrors.FieldServiceName, MsgUnsupportedService+": "+serviceName, nil)
	}
}
