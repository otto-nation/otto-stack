package command

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/otto-nation/otto-stack/internal/pkg/base"
	clicontext "github.com/otto-nation/otto-stack/internal/pkg/cli/context"
)

// Command represents a CLI command that can be executed
type Command interface {
	Execute(ctx context.Context, cliCtx clicontext.Context, base *base.BaseCommand) error
}

// Middleware represents middleware that can wrap command execution
type Middleware interface {
	Execute(ctx context.Context, cliCtx clicontext.Context, base *base.BaseCommand, next Command) error
}

// Handler wraps a command with middleware chain
type Handler struct {
	command     Command
	middlewares []Middleware
}

// NewHandler creates a new command handler with middleware
func NewHandler(cmd Command, middlewares ...Middleware) *Handler {
	return &Handler{
		command:     cmd,
		middlewares: middlewares,
	}
}

// Execute runs the command through the middleware chain
func (h *Handler) Execute(ctx context.Context, cliCtx clicontext.Context, base *base.BaseCommand) error {
	return h.executeWithMiddleware(ctx, cliCtx, base, 0)
}

func (h *Handler) executeWithMiddleware(ctx context.Context, cliCtx clicontext.Context, base *base.BaseCommand, index int) error {
	if index >= len(h.middlewares) {
		return h.command.Execute(ctx, cliCtx, base)
	}

	next := &nextCommand{
		handler: h,
		index:   index + 1,
	}

	return h.middlewares[index].Execute(ctx, cliCtx, base, next)
}

type nextCommand struct {
	handler *Handler
	index   int
}

func (n *nextCommand) Execute(ctx context.Context, cliCtx clicontext.Context, base *base.BaseCommand) error {
	return n.handler.executeWithMiddleware(ctx, cliCtx, base, n.index)
}

// CobraAdapter adapts the command pattern to work with cobra.Command
type CobraAdapter struct {
	handler *Handler
}

// NewCobraAdapter creates a new cobra adapter
func NewCobraAdapter(handler *Handler) *CobraAdapter {
	return &CobraAdapter{handler: handler}
}

// Handle adapts the command pattern to cobra's signature
func (a *CobraAdapter) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand) error {
	// For now, create empty context - this will be populated by processors
	cliCtx := clicontext.Context{}
	return a.handler.Execute(ctx, cliCtx, base)
}
