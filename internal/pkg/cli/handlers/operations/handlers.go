package operations

import (
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/stack"
)

// NewStatusHandler creates a new status handler
func NewStatusHandler() base.CommandHandler {
	return stack.NewStatusHandler()
}

// NewLogsHandler creates a new logs handler
func NewLogsHandler() base.CommandHandler {
	return stack.NewLogsHandler()
}

// NewExecHandler creates a new exec handler
func NewExecHandler() base.CommandHandler {
	return stack.NewExecHandler()
}

// NewConnectHandler creates a new connect handler
func NewConnectHandler() base.CommandHandler {
	return stack.NewConnectHandler()
}
