package base

import (
	"fmt"

	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/completion"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/project"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/stack"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/types"
)

// Registry manages command handlers
type Registry struct {
	handlers map[string]types.CommandHandler
}

// NewRegistry creates a new handler registry
func NewRegistry() *Registry {
	registry := &Registry{
		handlers: make(map[string]types.CommandHandler),
	}
	registry.registerDefaultHandlers()
	return registry
}

// RegisterHandler registers a handler for a command
func (r *Registry) RegisterHandler(name string, handler types.CommandHandler) {
	r.handlers[name] = handler
}

// GetHandler retrieves a handler by name
func (r *Registry) GetHandler(name string) (types.CommandHandler, error) {
	handler, exists := r.handlers[name]
	if !exists {
		return nil, fmt.Errorf("handler not found for command: %s", name)
	}
	return handler, nil
}

// GetAllHandlers returns all registered handlers
func (r *Registry) GetAllHandlers() map[string]types.CommandHandler {
	return r.handlers
}

// HasHandler checks if a handler exists for the given command
func (r *Registry) HasHandler(name string) bool {
	_, exists := r.handlers[name]
	return exists
}

// registerDefaultHandlers registers all default command handlers
func (r *Registry) registerDefaultHandlers() {
	r.RegisterHandler("up", stack.NewUpHandler())
	r.RegisterHandler("down", stack.NewDownHandler())
	r.RegisterHandler("restart", stack.NewRestartHandler())
	r.RegisterHandler("status", stack.NewStatusHandler())
	r.RegisterHandler("logs", stack.NewLogsHandler())
	r.RegisterHandler("exec", stack.NewExecHandler())
	r.RegisterHandler("connect", stack.NewConnectHandler())
	r.RegisterHandler("cleanup", stack.NewCleanupHandler())
	r.RegisterHandler("deps", project.NewDepsHandler())
	r.RegisterHandler("conflicts", project.NewConflictsHandler())
	r.RegisterHandler("services", project.NewServicesHandler())
	r.RegisterHandler("init", project.NewInitHandler())
	r.RegisterHandler("doctor", project.NewDoctorHandler())
	r.RegisterHandler("validate", project.NewValidateHandler())
	r.RegisterHandler("version", project.NewVersionHandler())
	r.RegisterHandler("completion", completion.NewCompletionHandler())
}
