package init

import (
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers"
	cliTypes "github.com/otto-nation/otto-stack/internal/pkg/cli/types"
)

func init() {
	handlers.Register(handlers.PackageInit, func(name string) cliTypes.CommandHandler {
		return NewInitHandler()
	})
}
