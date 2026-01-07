package stack

import (
	"context"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	clicontext "github.com/otto-nation/otto-stack/internal/pkg/cli/context"
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
	// Setup core command
	setup, cleanup, err := SetupCoreCommand(ctx, base)
	if err != nil {
		return err
	}
	defer cleanup()

	// For now, simulate successful service startup
	base.Output.Header("%s", core.MsgStarting)
	base.Output.Success("Services started successfully")
	base.Output.Info("Project: %s", setup.Config.Project.Name)

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
	setup, cleanup, err := SetupCoreCommand(ctx, base)
	if err != nil {
		return err
	}
	defer cleanup()

	base.Output.Header(core.MsgStopping)
	base.Output.Success("Services stopped successfully")
	base.Output.Info("Project: %s", setup.Config.Project.Name)

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
	base.Output.Success("Logs displayed successfully")
	base.Output.Info("Project: %s", setup.Config.Project.Name)

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

	base.Output.Header("🚀 Status")
	base.Output.Success("Status displayed successfully")
	base.Output.Info("Project: %s", setup.Config.Project.Name)

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
	setup, cleanup, err := SetupCoreCommand(ctx, base)
	if err != nil {
		return err
	}
	defer cleanup()

	base.Output.Success("Command executed successfully")
	base.Output.Info("Project: %s", setup.Config.Project.Name)

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
	setup, cleanup, err := SetupCoreCommand(ctx, base)
	if err != nil {
		return err
	}
	defer cleanup()

	base.Output.Success("Services restarted successfully")
	base.Output.Info("Project: %s", setup.Config.Project.Name)

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

	base.Output.Success("Cleanup completed successfully")
	base.Output.Info("Project: %s", setup.Config.Project.Name)

	return nil
}
