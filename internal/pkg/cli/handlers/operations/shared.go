package operations

import (
	"github.com/otto-nation/otto-stack/internal/pkg/cli/command"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/common"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/middleware"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
)

// ResolveServiceConfigs resolves services to ServiceConfigs using consistent logic across handlers
func ResolveServiceConfigs(args []string, setup *common.CoreSetup) ([]types.ServiceConfig, error) {
	if len(args) > 0 {
		// Resolve specific services from args
		serviceConfigs, err := services.ResolveUpServices(args, setup.Config)
		return serviceConfigs, err
	}
	// Use enabled services from config
	serviceConfigs, err := services.ResolveUpServices(setup.Config.Stack.Enabled, setup.Config)
	return serviceConfigs, err
}

// CreateStandardMiddlewareChain creates the standard middleware chain used by all stack handlers
func CreateStandardMiddlewareChain() (validationMiddleware, loggingMiddleware command.Middleware) {
	return middleware.NewInitializationMiddleware(), middleware.NewLoggingMiddleware()
}
