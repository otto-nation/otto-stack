package project

import (
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	clicontext "github.com/otto-nation/otto-stack/internal/pkg/cli/context"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	svc "github.com/otto-nation/otto-stack/internal/pkg/services"
)

// InteractiveProcessor handles interactive mode processing
type InteractiveProcessor struct {
	handler *InitHandler
}

// Process runs interactive prompts to gather project details
func (p *InteractiveProcessor) Process(flags any, base *base.BaseCommand) (clicontext.Context, error) {
	projectName, err := p.handler.promptManager.PromptForProjectDetails()
	if err != nil {
		return clicontext.Context{}, pkgerrors.NewValidationError(pkgerrors.FieldProjectName, MsgFailedToGetProjectDetails, err)
	}

	serviceConfigs, validation, advanced, err := p.handler.serviceSelectionManager.RunWorkflow(p.handler, base)
	if err != nil {
		return clicontext.Context{}, err
	}

	// Extract original service names from configs for interactive mode
	originalServiceNames := svc.ExtractServiceNames(serviceConfigs)

	ctx := clicontext.NewBuilder().
		WithProject(projectName, "").
		WithServices(originalServiceNames, serviceConfigs).
		WithValidation(validation).
		WithAdvanced(advanced).
		WithRuntime(false, true, false).
		Build()

	return ctx, nil
}
