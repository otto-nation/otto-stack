package core

import (
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers"
	cliTypes "github.com/otto-nation/otto-stack/internal/pkg/cli/types"
)

var coreHandlers = map[string]func() cliTypes.CommandHandler{
	"up":      func() cliTypes.CommandHandler { return NewUpHandler() },
	"down":    func() cliTypes.CommandHandler { return NewDownHandler() },
	"restart": func() cliTypes.CommandHandler { return NewRestartHandler() },
	"status":  func() cliTypes.CommandHandler { return NewStatusHandler() },
	"logs":    func() cliTypes.CommandHandler { return NewLogsHandler() },
	"exec":    func() cliTypes.CommandHandler { return NewExecHandler() },
	"connect": func() cliTypes.CommandHandler { return NewConnectHandler() },
	"cleanup": func() cliTypes.CommandHandler { return NewCleanupHandler() },
}

func init() {
	handlers.Register("core", func(name string) cliTypes.CommandHandler {
		if factory, exists := coreHandlers[name]; exists {
			return factory()
		}
		return nil
	})
}
