package handlers

import (
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/middleware"
)

type HandlerFactory func(string) base.CommandHandler

var registry = make(map[string]HandlerFactory)

func Register(packageName string, factory HandlerFactory) {
	registry[packageName] = factory
}

func Get(packageName, commandName string) base.CommandHandler {
	if factory, exists := registry[packageName]; exists {
		handler := factory(commandName)
		if handler == nil {
			return nil
		}
		return applyMiddleware(packageName, handler)
	}
	return nil
}

// applyMiddleware wraps handlers with the appropriate middleware chain.
// Lifecycle and operations commands get execution-context detection and
// project setup in addition to logging.
func applyMiddleware(packageName string, h base.CommandHandler) base.CommandHandler {
	switch packageName {
	case "lifecycle", "operations":
		return middleware.Chain(h,
			middleware.Logging(),
			middleware.WithExecContext(),
			middleware.WithProjectSetup(),
		)
	default:
		return middleware.Chain(h, middleware.Logging())
	}
}
