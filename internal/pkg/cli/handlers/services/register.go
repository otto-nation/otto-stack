package services

import (
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers"
	cliTypes "github.com/otto-nation/otto-stack/internal/pkg/cli/types"
)

var servicesHandlers = map[string]func() cliTypes.CommandHandler{
	"services":  func() cliTypes.CommandHandler { return NewServicesHandler() },
	"deps":      func() cliTypes.CommandHandler { return NewDepsHandler() },
	"conflicts": func() cliTypes.CommandHandler { return NewConflictsHandler() },
}

func init() {
	handlers.Register("services", func(name string) cliTypes.CommandHandler {
		if factory, exists := servicesHandlers[name]; exists {
			return factory()
		}
		return nil
	})
}
