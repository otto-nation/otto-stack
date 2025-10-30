package standard

import (
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/completion"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/doctor"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/validate"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/version"
	cliTypes "github.com/otto-nation/otto-stack/internal/pkg/cli/types"
)

var standardHandlers = map[string]func() cliTypes.CommandHandler{
	"doctor":     func() cliTypes.CommandHandler { return doctor.NewDoctorHandler() },
	"completion": func() cliTypes.CommandHandler { return completion.NewCompletionHandler() },
	"validate":   func() cliTypes.CommandHandler { return validate.NewValidateHandler() },
	"version":    func() cliTypes.CommandHandler { return version.NewVersionHandler() },
}

func init() {
	handlers.Register(handlers.PackageStandard, func(name string) cliTypes.CommandHandler {
		if factory, exists := standardHandlers[name]; exists {
			return factory()
		}
		return nil
	})
}
