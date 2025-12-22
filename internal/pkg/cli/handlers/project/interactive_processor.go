package project

import (
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
)

// InteractiveProcessor handles interactive mode processing
type InteractiveProcessor struct {
	handler *InitHandler
}

// Process runs interactive prompts to gather project details
func (p *InteractiveProcessor) Process(flags any, base *base.BaseCommand) (string, []string, map[string]bool, map[string]bool, error) {
	projectName, err := p.handler.promptManager.PromptForProjectDetails()
	if err != nil {
		return "", nil, nil, nil, pkgerrors.NewValidationError(pkgerrors.FieldProjectName, MsgFailedToGetProjectDetails, err)
	}

	services, validation, advanced, err := p.handler.serviceSelectionManager.RunWorkflow(p.handler, base)
	if err != nil {
		return "", nil, nil, nil, err
	}

	return projectName, services, validation, advanced, nil
}
