package project

import (
	"strings"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
)

// NonInteractiveProcessor handles non-interactive mode processing
type NonInteractiveProcessor struct {
	handler *InitHandler
}

// Process validates and processes flags for non-interactive mode
func (p *NonInteractiveProcessor) Process(flags any, base *base.BaseCommand) (string, []string, map[string]bool, map[string]bool, error) {
	initFlags := flags.(*core.InitFlags)

	if initFlags.Services == "" {
		return "", nil, nil, nil, pkgerrors.NewValidationError("services", "services flag is required in non-interactive mode", nil)
	}

	if initFlags.ProjectName == "" {
		return "", nil, nil, nil, pkgerrors.NewValidationError("project-name", "project name is required in non-interactive mode", nil)
	}

	services := parseServices(initFlags.Services)
	if err := p.handler.validateServices(services); err != nil {
		return "", nil, nil, nil, err
	}

	validation := getDefaultValidation()
	advanced := map[string]bool{}

	return initFlags.ProjectName, services, validation, advanced, nil
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
