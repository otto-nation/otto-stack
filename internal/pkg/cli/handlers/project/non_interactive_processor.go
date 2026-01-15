package project

import (
	"strings"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	clicontext "github.com/otto-nation/otto-stack/internal/pkg/cli/context"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	svc "github.com/otto-nation/otto-stack/internal/pkg/services"
)

// NonInteractiveProcessor handles non-interactive mode processing
type NonInteractiveProcessor struct {
	handler *InitHandler
}

// Process validates and processes flags for non-interactive mode
func (p *NonInteractiveProcessor) Process(flags any, base *base.BaseCommand) (clicontext.Context, error) {
	initFlags := flags.(*core.InitFlags)

	if initFlags.Services == "" {
		return clicontext.Context{}, pkgerrors.NewValidationError("services", "services flag is required in non-interactive mode", nil)
	}

	if initFlags.ProjectName == "" {
		return clicontext.Context{}, pkgerrors.NewValidationError("project-name", "project name is required in non-interactive mode", nil)
	}

	serviceNames := parseServices(initFlags.Services)

	// Convert service names to ServiceConfigs at entry point
	serviceConfigs, err := svc.ResolveUpServices(serviceNames, nil)
	if err != nil {
		return clicontext.Context{}, err
	}

	if err := p.handler.validateServiceConfigs(serviceConfigs); err != nil {
		return clicontext.Context{}, err
	}

	ctx := clicontext.NewBuilder().
		WithProject(initFlags.ProjectName, "").
		WithServices(serviceNames, serviceConfigs).
		WithValidation(getDefaultValidation()).
		WithAdvanced(map[string]bool{}).
		WithRuntime(initFlags.Force, false, false).
		Build()

	return ctx, nil
}

func parseServices(servicesStr string) []string {
	services := strings.Split(servicesStr, ",")
	for i := range services {
		services[i] = strings.TrimSpace(services[i])
	}
	return services
}

func getDefaultValidation() map[string]bool {
	validation := make(map[string]bool)
	for key := range ValidationRegistry {
		validation[key] = true
	}
	return validation
}
