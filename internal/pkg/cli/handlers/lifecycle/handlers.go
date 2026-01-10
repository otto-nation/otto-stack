package lifecycle

import (
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/stack"
)

// NewUpHandler creates a new up handler
func NewUpHandler() base.CommandHandler {
	return stack.NewUpHandler()
}

// NewDownHandler creates a new down handler
func NewDownHandler() base.CommandHandler {
	return stack.NewDownHandler()
}

// NewRestartHandler creates a new restart handler
func NewRestartHandler() base.CommandHandler {
	return stack.NewRestartHandler()
}

// NewCleanupHandler creates a new cleanup handler
func NewCleanupHandler() base.CommandHandler {
	return stack.NewCleanupHandler()
}
