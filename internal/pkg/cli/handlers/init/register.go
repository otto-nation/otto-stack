package init

import (
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers"
	cliTypes "github.com/otto-nation/otto-stack/internal/pkg/cli/types"
)

var initHandlers = map[string]func() cliTypes.CommandHandler{
	"init": func() cliTypes.CommandHandler { return NewInitHandler() },
}

func init() {
	handlers.Register("init", func(name string) cliTypes.CommandHandler {
		if factory, exists := initHandlers[name]; exists {
			return factory()
		}
		return nil
	})
}
