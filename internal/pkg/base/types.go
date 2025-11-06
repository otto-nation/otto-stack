package base

import (
	"context"

	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/logger"
	"github.com/spf13/cobra"
)

// Output interface for command output
type Output interface {
	Success(msg string, args ...any)
	Error(msg string, args ...any)
	Warning(msg string, args ...any)
	Info(msg string, args ...any)
	Header(msg string, args ...any)
	Muted(msg string, args ...any)
}

// ServiceManagerAdapter interface for service management operations
type ServiceManagerAdapter interface {
	StartServices(ctx context.Context, serviceNames []string, options docker.StartOptions) error
	StopServices(ctx context.Context, serviceNames []string, options docker.StopOptions) error
	GetServiceStatus(ctx context.Context, serviceNames []string) ([]docker.DockerServiceStatus, error)
}

// BaseCommand provides common dependencies for command handlers
type BaseCommand struct {
	Manager ServiceManagerAdapter
	Logger  logger.Adapter
	Output  Output
}

// CommandHandler interface for command handlers
type CommandHandler interface {
	Handle(ctx context.Context, cmd *cobra.Command, args []string, base *BaseCommand) error
}
