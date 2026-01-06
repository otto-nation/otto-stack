package project

import (
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	svc "github.com/otto-nation/otto-stack/internal/pkg/services"
)

// InteractiveProcessor handles interactive mode processing
type InteractiveProcessor struct {
	handler *InitHandler
}

// Process runs interactive prompts to gather project details
func (p *InteractiveProcessor) Process(flags any, base *base.BaseCommand) (string, []string, []svc.ServiceConfig, map[string]bool, map[string]bool, error) {
	projectName, err := p.handler.promptManager.PromptForProjectDetails()
	if err != nil {
		return "", nil, nil, nil, nil, pkgerrors.NewValidationError(pkgerrors.FieldProjectName, MsgFailedToGetProjectDetails, err)
	}

	serviceConfigs, validation, advanced, err := p.handler.serviceSelectionManager.RunWorkflow(p.handler, base)
	if err != nil {
		return "", nil, nil, nil, nil, err
	}

	// Extract original service names from configs for interactive mode
	// Note: serviceConfigs from interactive selection are already user-selected services (no dependencies resolved yet)
	originalServiceNames := svc.ExtractServiceNames(serviceConfigs)

	return projectName, originalServiceNames, serviceConfigs, validation, advanced, nil
}
