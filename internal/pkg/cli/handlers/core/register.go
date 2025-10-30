package core

import (
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers"
	cliTypes "github.com/otto-nation/otto-stack/internal/pkg/cli/types"
)

func init() {
	handlers.Register(handlers.PackageCore, func(name string) cliTypes.CommandHandler {
		handlerMap := map[string]func() cliTypes.CommandHandler{
			getCommandName("up"):      func() cliTypes.CommandHandler { return NewUpHandler() },
			getCommandName("down"):    func() cliTypes.CommandHandler { return NewDownHandler() },
			getCommandName("restart"): func() cliTypes.CommandHandler { return NewRestartHandler() },
			getCommandName("status"):  func() cliTypes.CommandHandler { return NewStatusHandler() },
			getCommandName("logs"):    func() cliTypes.CommandHandler { return NewLogsHandler() },
			getCommandName("exec"):    func() cliTypes.CommandHandler { return NewExecHandler() },
			getCommandName("connect"): func() cliTypes.CommandHandler { return NewConnectHandler() },
			getCommandName("cleanup"): func() cliTypes.CommandHandler { return NewCleanupHandler() },
		}

		if factory, exists := handlerMap[name]; exists {
			return factory()
		}
		return nil
	})
}

func getCommandName(cmd string) string {
	return cmd // Could be enhanced to derive from YAML
}
