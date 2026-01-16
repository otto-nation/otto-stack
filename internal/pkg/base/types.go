package base

import (
	"context"
	"io"

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
	Writer() io.Writer
}

// BaseCommand provides common dependencies for command handlers
type BaseCommand struct {
	Logger logger.Adapter
	Output Output
}

// GetVerbose extracts verbose flag from command
func (b *BaseCommand) GetVerbose(cmd *cobra.Command) bool {
	verbose, _ := cmd.Flags().GetBool("verbose")
	return verbose
}

// CommandHandler interface for command handlers
type CommandHandler interface {
	Handle(ctx context.Context, cmd *cobra.Command, args []string, base *BaseCommand) error
}
